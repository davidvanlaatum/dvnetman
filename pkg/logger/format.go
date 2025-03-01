package logger

import (
	"dvnetman/pkg/utils"
	"fmt"
	"github.com/fatih/color"
	"time"
)

type LogKeyFormatter interface {
	FormatLogKeyValue() interface{}
}

func ExpandLogKeyValue(value interface{}) interface{} {
	for {
		switch v := value.(type) {
		case LogKeyFormatter:
			value = v.FormatLogKeyValue()
		case error:
			return fmt.Sprintf("%+v", v)
		default:
			return v
		}
	}
}

type Formatter interface {
	Format(data *EventData) string
}

type ConsoleFormatter struct {
	color  bool
	caller bool
}

func NewConsoleFormatter() *ConsoleFormatter {
	return &ConsoleFormatter{caller: true}
}

func (f *ConsoleFormatter) EnableColor(color bool) *ConsoleFormatter {
	f.color = color
	return f
}

func (f *ConsoleFormatter) DisableCaller() *ConsoleFormatter {
	f.caller = false
	return f
}

func (f *ConsoleFormatter) Format(data *EventData) string {
	c := data.Level.ColorString()
	key := color.New(color.FgCyan)
	t := color.New(color.FgWhite)
	d := color.New(color.FgHiWhite)
	if f.color {
		c.EnableColor()
		key.EnableColor()
		t.EnableColor()
		d.EnableColor()
	}

	s := ""
	s = t.Sprint(data.Time.Format(time.RFC3339Nano))
	if f.caller {
		s += " "
		s += data.File
	}
	s += " "
	s += c.Sprint(data.Level.String())
	s += " "
	s += d.Sprint(data.Message)
	for _, v := range utils.MapSortedByKey(data.Keys, utils.MapSortedByKeyString) {
		s += " "
		s += key.Sprint(v.Key + "=")
		s += d.Sprintf("%v", ExpandLogKeyValue(v.Value))
	}
	return s
}
