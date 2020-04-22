package template

var (
	_templates map[string]string = make(map[string]string)
)

// Templates get all templates
func Templates() map[string]string {
	return _templates
}