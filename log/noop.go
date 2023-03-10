package log

type Noop struct {
}

func NewNoop() Logger {
	return &Noop{}
}

func (l *Noop) Fields(fields map[string]interface{}) Logger {
	return l
}
func (*Noop) Log(level Level, v ...interface{})                 {}
func (*Noop) Logf(level Level, format string, v ...interface{}) {}
