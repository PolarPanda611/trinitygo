package template

func init() {
	_templates["/go.mod"] = genMod()
}

func genMod() string {
	return `

module {{.PackageName}}

go 1.13

require (
	github.com/PolarPanda611/trinitygo {{.VersionNum}}
	github.com/google/uuid v1.1.1
)

`
}
