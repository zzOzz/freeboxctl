// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zzOzz/freeboxctl/freebox"
	"regexp"
)

var validArgs = []string{}
var downloadFile = false

// downloadsCmd represents the downloads command
var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fbx := freebox.GetInstance()
		if len(args) > 0 {
			var validBase64 = regexp.MustCompile(`^[-A-Za-z0-9+=]{1,50}|=[^=]|={3,}$`)
			var file = ""
			if (validBase64.MatchString(args[0])) {
				file = args[0]
			} else {
				file = base64.StdEncoding.EncodeToString([]byte(args[0]))
			}
			if downloadFile {
				stats, err := fbx.DownloadFile(file)
				if err != nil {
					logrus.Fatalf("fbx.Files(): %v", err)
				}
				fmt.Print(string(stats))
			} else {
				stats, err := fbx.Files(file)
				if err != nil {
					logrus.Fatalf("fbx.Files(): %v", err)
				}
				b, err := json.Marshal(stats)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println(string(b))
			}

		} else {

			stats, err := fbx.Files("")
			if err != nil {
				logrus.Fatalf("fbx.Files(): %v", err)
			}
			b, err := json.Marshal(stats)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(b))
		}

	},
	ValidArgs: validArgs,
}

func generatePathCompletion(file freebox.FileResult) {
	logrus.Info("tree - ", file.Name)
	//files, err := freebox.GetInstance().Files(base64.StdEncoding.EncodeToString([]byte(file.Name)))
	files, err := freebox.GetInstance().Files(file.Path)
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range *files {
		// element is the element from someSlice for where we are
		//fmt.Print("ici:" + file.Name)
		decoded, _ := base64.StdEncoding.DecodeString(file.Path)
		if file.Name != "." && file.Name != ".." {
			filesCmd.ValidArgs = append(filesCmd.ValidArgs, "/" + string(decoded))
			if len(*files) > 1 {
				generatePathCompletion(file)
			}
		}
	}
}

func init() {
	filesCmd.Flags().BoolVarP(&downloadFile, "download","d", false, "download file")
	getCmd.AddCommand(filesCmd)
}
