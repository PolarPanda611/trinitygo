package runtime

// RuntimeKey runtime key
type RuntimeKey interface {
	GetKeyName() string
	GetRequired() bool
	GetDefaultValue() string
	IsLog() bool
}

type runtimeKey struct {
	key          string
	required     bool
	defaultValue func() string
	islog        bool
}

func (r *runtimeKey) GetKeyName() string { return r.key }
func (r *runtimeKey) GetRequired() bool  { return r.required }
func (r *runtimeKey) IsLog() bool        { return r.islog }
func (r *runtimeKey) GetDefaultValue() string {
	if r.defaultValue == nil {
		return ""
	}
	return r.defaultValue()
}

// NewRuntimeKey Register new runtime key
// @islog if the runtime key been logged and been transported
// @required is false , the runtime key will use the
// newValueFunc to generate a new value
// usage : trace_id
// newValueFunc : func() string { return uuid.New().String() })
// to generate new trace_id
// p.s : the key should be lower case , because the grpc meta data will
// transfer all the key too lower case , if you use the upcase you will not
// find your runtime key in metadata
func NewRuntimeKey(key string, required bool, newValueFunc func() string, islog bool) RuntimeKey {
	return &runtimeKey{
		key:          key,
		required:     required,
		defaultValue: newValueFunc,
		islog:        islog,
	}
}
