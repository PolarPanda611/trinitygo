package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	ttemplate "github.com/PolarPanda611/trinitygo/trinitygo/template"
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "trinitygo",
		Short: "New HTTP project which will be implemented the Trinity GO",
		Long:  `This command will generate the basic fold structure `,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			fmt.Println(args)
			if len(args) != 2 {
				return errors.New("[project_name] empty")
			}
			projectName := args[1]
			sysPath, err := filepath.Abs(projectName)
			if err != nil {
				return err
			}
			initHTTPFolder(sysPath)
			pData := map[string]interface{}{
				"PackageName": projectName,
				"VersionNum":  versionNum,
			}

			m := ttemplate.Templates()
			for filepath, content := range m {
				var pf *os.File
				pf, err = os.Create(sysPath + filepath)
				if err != nil {
					return err
				}
				tmpl, err := template.New(projectName).Parse(content)
				if err != nil {
					return err
				}
				if err = tmpl.Execute(pf, pData); err != nil {
					return err
				}
			}
			exec.Command("gofmt", "-w", sysPath).Output()
			return nil
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
