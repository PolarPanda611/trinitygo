/**
 * @ Author: Daniel Tan
 * @ Date: 2020-04-22 09:17:01
 * @ LastEditTime: 2020-08-04 15:45:21
 * @ LastEditors: Daniel Tan
 * @ Description:
 * @ FilePath: /trinitygo/trinitygo/cmd/version.go
 * @
 */
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	versionNum = "v0.1.6"
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
