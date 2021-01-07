/*
Package cmd version
Copyright © 2021 Evan

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os/exec"

	"github.com/buger/jsonparser"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "View the version of tkctl and tidb-operator",
	Long:  `View the version of tkctl and tidb-operator`,
	Run: func(cmd *cobra.Command, args []string) {
		podShell := exec.Command("kubectl", "get", "pods", "-A", "-o", "json")
		podsList, _ := podShell.Output()
		var tidbControllerManagerImage string
		var tidbSchedulerImage string
		jsonparser.ArrayEach(podsList, func(pod []byte, dataType jsonparser.ValueType, offset int, err error) {
			serviceAccountName, _ := jsonparser.GetString(pod, "spec", "serviceAccountName")
			if serviceAccountName == "tidb-controller-manager" || serviceAccountName == "tidb-scheduler" {
				jsonparser.ArrayEach(pod, func(container []byte, dataType jsonparser.ValueType, offset int, err error) {
					imageName, _ := jsonparser.GetString(container, "image")
					containerName, _ := jsonparser.GetString(container, "name")
					if containerName == "tidb-operator" && serviceAccountName == "tidb-controller-manager" {
						tidbControllerManagerImage = imageName
					}
					if containerName == "tidb-scheduler" {
						tidbSchedulerImage = imageName
					}
				}, "spec", "containers")
			}
		}, "items")

		fmt.Println("Welcome to tkctl-golang: " + TKCTLVERSION)
		fmt.Println("TiDB Controller Manager Version: " + tidbControllerManagerImage)
		fmt.Println("TiDB Scheduler Version: " + tidbSchedulerImage)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
