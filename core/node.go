package core

type Node interface {
	Name() string
	Type() NodeType
	// ToShare returns the node as a shareable string
	ToShare() string
	// ToClash returns the node as a Clash config string
	ToClash() string
	// FromShare parses a shareable string and initializes the node with its values.
	// It returns an error if the string is not in a valid format or if any step
	// of the initialization process fails.
	FromShare(string) error
	// FromClash parses a Clash config string and initializes the node with its values.
	// It returns an error if the string is not in a valid format or if any step
	// of the initialization process fails.
	FromClash([]byte) error
}
