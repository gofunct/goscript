// Copyright © 2019 Coleman Word <coleman.word@gofunct.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package aws

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(dbCmd)
}

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		Service()
	},
}

func Service() {
	log.SetFlags(0)
	log.SetPrefix("aws/provision_db: ")
	host := flag.String("host", "", "hostname of database")
	securityGroup := flag.String("security_group", "", "database security group")
	database := flag.String("database", "", "name of database to provision")
	password := flag.String("password", "", "root password on database")
	schema := flag.String("schema", "", "path to .sql file defining the database schema")
	flag.Parse()
	missing := false
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() == "" {
			log.Printf("Required flag -%s not set.", f.Name)
			missing = true
		}
	})
	if missing {
		os.Exit(64)
	}
	if err := provisionDb(*host, *securityGroup, *database, *password, *schema); err != nil {
		log.Fatal(err)
	}
}

func provisionDb(DbHost, securityGroupID, DbName, DbPassword, schemaPath string) error {
	const mySQLImage = "mysql:5.6"

	// Pull the necessary Docker images.
	log.Print("Downloading Docker images...")
	if _, err := run("docker", "pull", mySQLImage); err != nil {
		return err
	}

	// Create a temporary directory to hold the certificates.
	// We resolve all symlinks to avoid Docker on Mac issues, see
	// https://github.com/google/go-cloud/issues/110.
	tempdir, err := ioutil.TempDir("", "guestbook-ca")
	if err != nil {
		return fmt.Errorf("creating temp dir for certs: %v", err)
	}
	defer os.RemoveAll(tempdir)
	tempdir, err = filepath.EvalSymlinks(tempdir)
	if err != nil {
		return fmt.Errorf("evaluating any symlinks: %v", err)
	}
	resp, err := http.Get("https://s3.amazonaws.com/rds-downloads/rds-ca-2015-root.pem")
	if err != nil {
		return fmt.Errorf("fetching pem file: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("response status code is %d, want 200", resp.StatusCode)
	}
	defer resp.Body.Close()
	caPath := filepath.Join(tempdir, "rds-ca.pem")
	caFile, err := os.Create(caPath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(caFile, resp.Body); err != nil {
		return fmt.Errorf("copying response to file: %v", err)
	}

	log.Print("Adding a temporary ingress rule")
	if _, err := run("aws", "ec2", "authorize-security-group-ingress", "--group-id", securityGroupID, "--protocol=tcp", "--port=3306", "--cidr=0.0.0.0/0"); err != nil {
		return err
	}
	defer func() {
		log.Print("Removing ingress rule...")
		if _, err := run("aws", "ec2", "revoke-security-group-ingress", "--group-id", securityGroupID, "--protocol=tcp", "--port=3306", "--cidr=0.0.0.0/0"); err != nil {
			log.Print(err)
		}
	}()
	log.Printf("Added ingress rule to %s for port 3306", securityGroupID)

	// Send schema.
	log.Print("Sending schema to database...")
	schema, err := os.Open(schemaPath)
	if err != nil {
		return err
	}
	defer schema.Close()
	mySQLCmd := fmt.Sprintf(`mysql -h'%s' -uroot -p'%s' --ssl-ca=/ca/rds-ca.pem '%s'`, DbHost, DbPassword, DbName)
	connect := exec.Command("docker", "run", "--rm", "--interactive", "--volume", tempdir+":/ca", mySQLImage, "sh", "-c", mySQLCmd)
	connect.Stdin = schema
	connect.Stderr = os.Stderr
	if err := connect.Run(); err != nil {
		return fmt.Errorf("running %v: %v", connect.Args, err)
	}

	return nil
}

func run(args ...string) (stdout string, err error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, os.Environ()...)
	stdoutb, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("running %v: %v", cmd.Args, err)
	}
	return strings.TrimSpace(string(stdoutb)), nil
}
