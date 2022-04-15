package simplelog

// Logger interface, log.Logger satisfy this interface
type Logger interface {
	Println(...interface{})
	Print(...interface{})
	Printf(string, ...interface{})
}

// NoopLogger is a no-op Logger implementation
type NoopLogger struct{}

// Println is a no-op Println implementation
func (n *NoopLogger) Println(...interface{}) {}

// Print is a no-op Print implementation
func (n *NoopLogger) Print(...interface{}) {}

// Printf is a no-op Printf implementation
func (n *NoopLogger) Printf(string, ...interface{}) {}
