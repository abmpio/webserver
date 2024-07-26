package app

import (
	"strings"
	"time"

	"github.com/abmpio/abmp/pkg/log"
	"github.com/abmpio/abmp/pkg/utils/validator"
	"github.com/abmpio/app"
	"github.com/abmpio/app/cli"
	"github.com/abmpio/app/host"
	"github.com/abmpio/configurationx"
	cors "github.com/abmpio/webserver/app/middleware/cors"
	errHandler "github.com/abmpio/webserver/app/middleware/err"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	requestLogger "github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func init() {
	app.Register(NewApplication)
}

func newIrisApplication() *iris.Application {
	app := iris.New()
	//错误封装
	app.Use(errHandler.New())
	app.Use(recover.New())
	app.Use(requestLogger.New(requestLogConfig()))
	if configurationx.GetInstance().Web != nil {
		cors.UseCors(app.APIBuilder, configurationx.GetInstance().Web.Cors)
	}

	//设置validator
	app.Validator = validator.Validate

	return app
}

func requestLogConfig() requestLogger.Config {
	c := requestLogger.DefaultConfig()
	c.AddSkipper(func(ctx *context.Context) bool {
		p := ctx.Path()
		return strings.HasPrefix(p, "/api/health/check")
	})
	return c
}

type Application struct {
	*iris.Application
	Address string

	isBuilded        bool
	irisConfigurator []iris.Configurator
	Err              error
}

type Configurator func(*Application)

func NewApplication() *Application {
	app := &Application{
		Application:      newIrisApplication(),
		irisConfigurator: make([]iris.Configurator, 0),
		isBuilded:        false,
	}

	return app
}

func (a *Application) Configure(configurators ...Configurator) *Application {
	return a
}

// build application environments
func (a *Application) Build(configurators ...Configurator) *Application {
	if a.isBuilded {
		return a
	}
	if a.Err != nil {
		return a
	}
	defer func() {
		a.isBuilded = true
	}()
	envHttp := host.GetHostEnvironment().GetEnvString(host.ENV_HTTP)
	if len(envHttp) > 0 {
		a.Address = envHttp
	} else {
		host.GetHostEnvironment().SetHttp(a.Address)
	}
	if len(a.Address) <= 0 {
		msg := "没有配置好app.http参数"
		log.Error(msg)
		panic(msg)
	}

	cli.GetHost().Application().ConfigureService()

	// a.pprofStartupAction()
	//运行启动项
	cli.GetHost().Application().RunStartup()

	//构建配置
	appConfigurators := make([]iris.Configurator, 0)
	for _, eachConfigurator := range configurators {
		if eachConfigurator == nil {
			continue
		}
		newAppConfigurator := func(irisApp *iris.Application) {
			eachConfigurator(a)
		}
		appConfigurators = append(appConfigurators, newAppConfigurator)
	}
	a.irisConfigurator = appConfigurators

	//设置启动消耗的时间
	startTime := host.GetHostEnvironment().GetEnv(host.ENV_StartTime).(time.Time)
	interval := time.Since(startTime)
	host.GetHostEnvironment().SetEnv(host.ENV_StartInterval, interval)

	return a
}

func (a *Application) Run(configurators ...Configurator) *Application {
	a.Build(configurators...)

	err := a.Application.Run(iris.Addr(a.Address), a.irisConfigurator...)
	a.Err = err
	return a
}

// func (a *Application) pprofStartupAction() {
// 	if app.HostApplication.SystemConfig().App.IsRunInCli {
// 		return
// 	}

// 	log.Logger.Debug("正在构建pprof路径组件,/debug/pprof...")
// 	a.Any("/debug/pprof/cmdline", iris.FromStd(pprof.Cmdline))
// 	a.Any("/debug/pprof/profile", iris.FromStd(pprof.Profile))
// 	a.Any("/debug/pprof/symbol", iris.FromStd(pprof.Symbol))
// 	a.Any("/debug/pprof/trace", iris.FromStd(pprof.Trace))
// 	a.Any("/debug/pprof/debug/pprof/{action:string}", requestPprof.New())

// 	httpValue := os.Getenv("app.http")
// 	advertiseHostValue := os.Getenv("app.advertisehost")
// 	if len(httpValue) > 0 {
// 		pprofPath := httpValue
// 		if len(advertiseHostValue) > 0 {
// 			pprofPath = strings.Replace(httpValue, "0.0.0.0", advertiseHostValue, 1)
// 		}
// 		log.Logger.Debug(fmt.Sprintf("已经构建好pprof路径组件,你可以通过 %s/debug/pprof 来访问pprof", pprofPath))
// 	}
// }
