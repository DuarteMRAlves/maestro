package message

// Instance specifies an interface for concrete messages. These messages
// can be sent and received in stages.
// An instance can have subfields, that can be updated.
type Instance interface {
	// Set updates the value of a given field.
	Set(Field, Instance) error
	// Get returns the current value of a field.
	Get(Field) (Instance, error)
}

// Type describes the underlying struct of Message.
type Type interface {
	Builder
	// Subfield returns the Type of a subfield.
	Subfield(Field) (Type, error)
	// Verifies whether these types are comparible.
	Compatible(Type) bool
}

// Field specifies the name of a field in a message in the pipeline.
type Field string

// IsUnspecified reports whether this field is "".
func (m Field) IsUnspecified() bool { return m == "" }

// Creates an empty instance of this type
type Builder interface {
	Build() Instance
}

// BuildFunc creates empty messages of a given type
type BuildFunc func() Instance

func (fn BuildFunc) Build() Instance { return fn() }
