package healthcheck

import (
	"github.com/abmpio/abmp/pkg/log"
	"github.com/abmpio/app"
)

func init() {

	log.Logger.Info("webserver.healthcheck starter init")
	app.RegisterStartupAction(healthcheckStartup)
}
