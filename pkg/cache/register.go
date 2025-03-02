package cache

import (
	"context"
	"net/url"
	"sort"
)

var drivers = map[string]Driver{}

func Register(driver Driver) {
	if drivers[driver.Name()] != nil {
		panic("driver already registered")
	}
	drivers[driver.Name()] = driver
}

func GetDriver(name string) Driver {
	return drivers[name]
}

func GetDriverNames() []string {
	names := make([]string, 0, len(drivers))
	for name := range drivers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func NewPool(ctx context.Context, config string) (p Pool, err error) {
	u, err := url.Parse(config)
	if err != nil {
		return
	}
	d := GetDriver(u.Scheme)
	var c Cache
	if c, err = d.New(ctx, u); err != nil {
		return
	}
	var ok bool
	if p, ok = (interface{})(c).(Pool); ok {
		return
	}
	panic("not implemented")
}
