package store

// Type defines the possible types for Scopes
type Type uint

const (
	Unknown Type = 0
	Global  Type = 1
	Part    Type = 2
)

func (s Type) String() string {
	return [...]string{
		"unknown",
		"global",
		"part",
	}[s]
}

var Map = map[string]Type{
	Global.String(): Global,
	Part.String():   Part,
}
