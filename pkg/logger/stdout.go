package logger

type StdoutDriver struct {
	formatter Formatter
}

func (s *StdoutDriver) Log(data *EventData) {
	msg := s.formatter.Format(data)
	println(msg)
}

func (s *StdoutDriver) HelperFuncs() []func() {
	return nil
}

var _ Driver = (*StdoutDriver)(nil)

func NewStdoutDriver(formatter Formatter) *StdoutDriver {
	return &StdoutDriver{formatter: formatter}
}
