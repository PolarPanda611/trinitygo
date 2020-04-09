package runtime

// RuntimeKey runtime key
type RuntimeKey interface {
	GetKeyName() string
	GetRequired() bool
	GetDefaultValue() string
}

type runtimeKey struct {
	Key          string
	Required     bool
	DefaultValue func() string
}

func (r *runtimeKey) GetKeyName() string { return r.Key }
func (r *runtimeKey) GetRequired() bool  { return r.Required }
func (r *runtimeKey) GetDefaultValue() string {
	if r.DefaultValue == nil {
		return ""
	}
	return r.DefaultValue()
}

// NewRuntimeKey Register new runtime key
func NewRuntimeKey(key string, required bool, defaultValue func() string) RuntimeKey {
	return &runtimeKey{
		Key:          key,
		Required:     required,
		DefaultValue: defaultValue,
	}
}
