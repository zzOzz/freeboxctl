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
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zzOzz/freeboxctl/freebox"
	"regexp"
)

var validArgs = []string{}

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
		// fbx := freebox.New()
		fbx := freebox.GetInstance()

		//err := fbx.Connect()
		//if err != nil {
		//	logrus.Fatalf("fbx.Connect(): %v", err)
		//}
		//
		//err = fbx.Authorize()
		//if err != nil {
		//	logrus.Fatalf("fbx.Authorize(): %v", err)
		//}
		//
		//err = fbx.Login()
		//if err != nil {
		//	logrus.Fatalf("fbx.Login(): %v", err)
		//}

		var validBase64 = regexp.MustCompile(`^[-A-Za-z0-9+=]{1,50}|=[^=]|={3,}$`)
		var file = ""
		if (validBase64.MatchString(args[0])) {
			file = args[0]
		} else {
			file = base64.StdEncoding.EncodeToString([]byte(args[0]))
		}
		stats, err := fbx.DownloadFile(file)
		if err != nil {
			logrus.Fatalf("fbx.Files(): %v", err)
		}
		fmt.Print(string(stats))
		//b, err := json.Marshal(stats)
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//fmt.Println(string(b))
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

	//files, err := freebox.GetInstance().Files(base64.StdEncoding.EncodeToString([]byte("")))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//for _, file := range *files {
	//	// element is the element from someSlice for where we are
	//	//fmt.Print("ici:" + file.Name)
	//	filesCmd.ValidArgs = append(filesCmd.ValidArgs, file.Name)
	//}
	// var rootPath = freebox.FileResult{Name:""}
	// generatePathCompletion(rootPath)
	getCmd.AddCommand(filesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
