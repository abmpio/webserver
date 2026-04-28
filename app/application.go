package app

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/abmpio/abmp/pkg/log"
	"github.com/abmpio/abmp/pkg/utils/validator"
	"github.com/abmpio/app"
	"github.com/abmpio/app/cli"
	"github.com/abmpio/app/host"
	"github.com/abmpio/configurationx"
	jsonUtil "github.com/abmpio/libx/json"
	cors "github.com/abmpio/webserver/app/middleware/cors"
	errHandler "github.com/abmpio/webserver/app/middleware/err"
	recover "github.com/abmpio/webserver/app/middleware/recover"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	requestLogger "github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/requestid"
)

func init() {
	app.Register(NewApplication)
}

func newIrisApplication() *iris.Application {
	app := iris.New()
	app.Use(requestid.New())
	app.Use(requestLogger.New(requestLogConfig()))
	//错误封装
	app.Use(errHandler.New())
	app.Use(recover.New())
	if configurationx.GetInstance().Web != nil {
		cors.UseCors(app.APIBuilder, configurationx.GetInstance().Web.Cors)
	}

	//设置validator
	app.Validator = validator.Validate

	return app
}

func requestLogConfig() requestLogger.Config {
	c := requestLogger.DefaultConfig()
	// c.MessageContextKeys = []string{
	// 	"iris.context.id",
	// 	"userId",
	// }
	c.MessageHeaderKeys = []string{
		"tenant",
	}
	c.AddSkipper(func(ctx *context.Context) bool {
		p := ctx.Path()
		skipped := strings.HasPrefix(p, "/api/health/check")
		if skipped {
			return true
		}
		return _irisApplicationConfiguratorOptions.shouldSkipRequestLogPath(p)
	})
	logInfo := func(ctx *context.Context, latency time.Duration) string {
		// all except latency to string
		var status, ip, method, path string
		requestId := ctx.GetID()
		ip = ctx.RemoteAddr()
		method = ctx.Method()
		path = ctx.Request().URL.RequestURI()
		status = strconv.Itoa(ctx.GetStatusCode())
		userId := ctx.Values().Get("userId")
		headerMessage := jsonUtil.ObjectToJson(getHeaderMessages([]string{
			"tenant",
		}, ctx))
		line := fmt.Sprintf("%s %s,requestId:%v,userId:%v,status:%v,duration:%4v,ip:%s,header:%s",
			method,
			path,
			requestId,
			userId,
			status,
			latency,
			ip,
			headerMessage)
		return line
	}
	c.LogFuncCtx = func(ctx *context.Context, latency time.Duration) {
		// no new line, the framework's logger is responsible how to render each log.
		line := logInfo(ctx, latency)
		log.Logger.Info(line)
	}
	return c
}

func getHeaderMessages(keyList []string, ctx *context.Context) map[string]string {
	m := make(map[string]string)
	for _, key := range keyList {
		msg := ctx.GetHeader(key)
		m[key] = msg
	}
	return m
}

type irisApplicationConfiguratorOptions struct {
	// 用来跳过logger的path比较函数列表
	requestLogPathSkipped []func(path string) bool
}

var (
	_irisApplicationConfiguratorOptions = &irisApplicationConfiguratorOptions{
		requestLogPathSkipped: make([]func(path string) bool, 0),
	}
)

// 检测path对应的request log是否应该skip
func (o *irisApplicationConfiguratorOptions) shouldSkipRequestLogPath(path string) bool {
	if len(o.requestLogPathSkipped) <= 0 {
		return false
	}
	for _, eachFn := range o.requestLogPathSkipped {
		if eachFn == nil {
			continue
		}
		if eachFn(path) {
			return true
		}
	}
	return false
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

	// set host.ENV_Healthcheck,host.ENV_AdvertiseHost env value
	setHealthCheckEnv()
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

// AddRequestLogSkipper 添加进程级请求日志跳过规则。
// 调用方应在Application.Run前完成配置；请求处理阶段只读，运行期不要再修改该规则列表。
func (a *Application) AddRequestLogSkipper(pathFn func(path string) bool) {
	if pathFn == nil {
		return
	}
	_irisApplicationConfiguratorOptions.requestLogPathSkipped = append(_irisApplicationConfiguratorOptions.requestLogPathSkipped, pathFn)
}
