package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/PolarPanda611/trinitygo/trinitygo/crudtemplate"
	"github.com/PolarPanda611/trinitygo/util"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

var (
	newCRUDCmd = &cobra.Command{
		Use:   "NewCrud",
		Short: "generate new crud code for model",
		Long:  `This command will generate the crud code for your model name , first characters should be Upper case `,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) != 1 {
				return errors.New("[ModelName] empty")
			}
			modelName := args[0]
			modelNameWithUnderscore := gorm.ToColumnName(modelName)
			currentPath, err := os.Getwd()
			if err != nil {
				return err
			}
			projectName := filepath.Base(currentPath)
			charRune := []rune(modelName)
			firstChar := charRune[0]
			if !unicode.IsUpper(firstChar) {
				return fmt.Errorf("first characters should be Upper case ")
			}
			pData := map[string]interface{}{
				"ProjectName":           projectName,
				"ModelName":             modelName,
				"ModelNamePrivate":      fmt.Sprint(strings.ToLower(modelName[0:1]) + modelName[1:]),
				"ModelNameToUnderscore": modelNameWithUnderscore,
			}
			fmt.Println(pData)
			templates := crudtemplate.Templates()
			for filepath, content := range templates {
				targetPathFmt := currentPath + filepath
				targetDir := path.Dir(targetPathFmt)
				fmt.Println(targetDir)
				if !util.CheckFileIsExist(targetDir) {
					return fmt.Errorf("%v not exist ", targetDir)
				}
				targetPath := fmt.Sprintf(targetPathFmt, modelNameWithUnderscore)
				var pf *os.File
				pf, err = os.Create(targetPath)
				if err != nil {
					return err
				}
				tmpl, err := template.New(modelName).Parse(content)
				if err != nil {
					return err
				}
				if err = tmpl.Execute(pf, pData); err != nil {
					return err
				}
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(newCRUDCmd)
}
