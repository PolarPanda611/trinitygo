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
	newProjectCmd = &cobra.Command{
		Use:   "newhttp",
		Short: "New HTTP project which will be implemented the Trinity GO",
		Long:  `This command will generate the basic fold structure `,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) != 1 {
				return errors.New("[project_name] empty")
			}
			projectName := args[0]
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
				fmt.Println(sysPath + filepath)
				pf, err = os.Create(sysPath + filepath)
				if err != nil {
					return err
				}
				tmpl, err := template.New(projectName).Parse(content)
				if err = tmpl.Execute(pf, pData); err != nil {
					return err
				}
			}
			exec.Command("gofmt", "-w", sysPath).Output()
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(newProjectCmd)
}
