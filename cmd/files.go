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
	"fmt"
	"github.com/spf13/cobra"
	"github.com/Sirupsen/logrus"
	"encoding/json"
	"github.com/zzOzz/freeboxctl/freebox"
)

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
		stats, err := fbx.Files(args[0])
		if err != nil {
			logrus.Fatalf("fbx.Files(): %v", err)
		}
		//fmt.Println(stats)
		b, err := json.Marshal(stats)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))
	},
}

func init() {
	getCmd.AddCommand(filesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
