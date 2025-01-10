package pprof

import (
	"fmt"
	"os"
	"strings"

	"github.com/abmpio/abmp/pkg/log"
	"github.com/abmpio/app"
	webapp "github.com/abmpio/webserver/app"

	"net/http/pprof"

	"github.com/kataras/iris/v12"
	requestPprof "github.com/kataras/iris/v12/middleware/pprof"
)

func init() {
	fmt.Println("webserver.pprof starter init")

	app.RegisterStartupAction(pprofStartupAction)
}

func pprofStartupAction(webApp *webapp.Application) app.IStartupAction {
	return app.NewStartupAction(func() {
		if app.HostApplication.SystemConfig().App.IsRunInCli {
			return
		}

		log.Logger.Debug("正在构建pprof路径组件,/debug/pprof...")
		webApp.Any("/debug/pprof/cmdline", iris.FromStd(pprof.Cmdline))
		webApp.Any("/debug/pprof/profile", iris.FromStd(pprof.Profile))
		webApp.Any("/debug/pprof/symbol", iris.FromStd(pprof.Symbol))
		webApp.Any("/debug/pprof/trace", iris.FromStd(pprof.Trace))
		webApp.Any("/debug/pprof/debug/pprof/{action:string}", requestPprof.New())

		httpValue := os.Getenv("app.http")
		advertiseHostValue := os.Getenv("app.advertisehost")
		if len(httpValue) > 0 {
			pprofPath := httpValue
			if len(advertiseHostValue) > 0 {
				pprofPath = strings.Replace(httpValue, "0.0.0.0", advertiseHostValue, 1)
			}
			log.Logger.Debug(fmt.Sprintf("已经构建好pprof路径组件,你可以通过 %s/debug/pprof 来访问pprof", pprofPath))
		}
	})
}
