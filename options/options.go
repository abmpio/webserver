package options

import (
	"fmt"
	"sync"

	"github.com/abmpio/abmp/pkg/log"
	"github.com/abmpio/configurationx"
)

type Options struct {
	Log *LogOptions `json:"log"`
}

type LogOptions struct {
	EnableLogWhenHandlerError bool `json:"enableLogWhenHandlerError"`
	EnableLogRequest          bool `json:"enableLogRequest"`
}

func newDefaultLogOptions() *LogOptions {
	return &LogOptions{
		EnableLogWhenHandlerError: true,
		EnableLogRequest:          true,
	}
}

func (o *Options) Normalize() *Options {
	if o.Log == nil {
		o.Log = newDefaultLogOptions()
	}
	return o
}

const (
	ConfigurationKey string = "webserver"
)

var (
	_option *Options
	_once   sync.Once
)

func GetOptions() *Options {
	_once.Do(func() {
		_option = &Options{}
		if err := configurationx.GetInstance().UnmarshFromKey(ConfigurationKey, _option); err != nil {
			err = fmt.Errorf("无效的配置文件,%s", err)
			log.Logger.Error(err.Error())
			panic(err)
		}
		_option.Normalize()
	})
	return _option
}
