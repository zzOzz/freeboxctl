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

// tasksCmd represents the downloads command
var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fbx := freebox.GetInstance()
		stats, err := fbx.FSTasks()
		if err != nil {
			logrus.Fatalf("fbx.FSTasks(): %v", err)
		}
		b, err := json.Marshal(stats)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))
	},
}

func init() {
	getCmd.AddCommand(tasksCmd)
}
