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
// when the required is false , the runtime key will use the
// newValueFunc to generate a new value
// usage : trace_id
// newValueFunc : func() string { return uuid.New().String() })
// to generate new trace_id
// p.s : the key should be lower case , because the grpc meta data will
// transfer all the key too lower case , if you use the upcase you will not
// find your runtime key in metadata
func NewRuntimeKey(key string, required bool, newValueFunc func() string) RuntimeKey {
	return &runtimeKey{
		Key:          key,
		Required:     required,
		DefaultValue: newValueFunc,
	}
}
