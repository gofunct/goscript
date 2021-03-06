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

package protoc

import (
	kitlog "github.com/go-kit/kit/log"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// protocCmd represents the protoc command
var ProtocCmd = &cobra.Command{
	Use:   "protoc",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		if err := Grpc(dir); err != nil {
			log.Fatalln("failed to execute command", err)
		}
	},
}

func init() {
	ProtocCmd.Flags().StringVar(&dir, "dir", ".", "path to directory containing protobuf files")
	logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	logger = kitlog.With(logger, "time", kitlog.DefaultTimestampUTC, "exec", kitlog.DefaultCaller, "dir", dir)
	log.SetOutput(kitlog.NewStdlibAdapter(logger))

}

func Grpc(d string) error {

	return filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalln(err)
		}
		// skip vendor directory
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		// find all protobuf files
		if filepath.Ext(path) == ".proto" {
			// args
			args := []string{
				"--go_out=plugins=grpc:.",
				path,
			}
			cmd := exec.Command("protoc", args...)
			log.Print("starting command")
			cmd.Env = os.Environ()
			if err := cmd.Run(); err != nil {
				return err
			}
		}
		return nil
	})
}
