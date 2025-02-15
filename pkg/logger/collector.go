package logger

type Collector struct {
	logs []EventData
}

func (c *Collector) HelperFuncs() []func() {
	return nil
}

func (c *Collector) Log(data *EventData) {
	c.logs = append(c.logs, *data)
}

func (c *Collector) Logs() []EventData {
	return c.logs
}

func NewCollector() *Collector {
	return &Collector{}
}

var _ Driver = (*Collector)(nil)
