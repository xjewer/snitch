package simplelog

type Logger interface {
	Println(...interface{})
	Print(...interface{})
	Printf(string, ...interface{})
}

type NoopLogger struct {}

func (n *NoopLogger) Println(...interface{}) {}
func (n *NoopLogger) Print(...interface{}) {}
func (n *NoopLogger) Printf(string, ...interface{}) {}
