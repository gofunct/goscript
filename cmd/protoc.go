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

package cmd

import (
	"context"
	"fmt"
	"github.com/gofunct/goscript/script"
	"github.com/gofunct/goscript/script/protoc"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	protoScrpt = protoc.Inject()
	protoPath  string
)

// protocCmd represents the protoc command
var protocCmd = &cobra.Command{
	Use:   "protoc",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ProtoArgs()
	},
}

func init() {
	rootCmd.AddCommand(protocCmd)
	protocCmd.Flags().StringVarP(&protoPath, "proto", "p", "proto", "path to protobuf files")
}

func ProtoArgs() error {
	return filepath.Walk(protoPath, func(path string, info os.FileInfo, err error) error {
		// skip vendor directory
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		// find all protobuf files
		if filepath.Ext(path) == ".proto" {
			// args
			args := []string{
				"-I=.",
				fmt.Sprintf("-I=%s", filepath.Join(goPath, "src")),
				fmt.Sprintf("-I=%s", filepath.Join(goPath, "src", "github.com", "gogo", "protobuf", "protobuf")),
				fmt.Sprintf("--proto_path=%s", filepath.Join(goPath, "src", "github.com")),
				"--gogofaster_out=Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:.",
				path,
			}
			return protoScrpt.Run(context.Background(), script.WithArgs(args...))
		}
		return nil
	})
}
