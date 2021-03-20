/*
 * @Author: Daniel TAN
 * @Description:
 * @Date: 2020-09-01 09:15:45
 * @LastEditTime: 2021-03-17 01:48:00
 * @LastEditors: Daniel TAN
 * @FilePath: /trinitygo/trinitygo/cmd/version.go
 */
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	versionNum = "v0.1.35"
)

var (
	versionCmd = &cobra.Command{
		Use:   "Version",
		Short: "Output current version number",
		Long:  `Output current version number`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			fmt.Println("trinitygo " + versionNum)
			return
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
