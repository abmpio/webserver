package healthcheck

import (
	"fmt"

	"github.com/abmpio/app"
)

func init() {

	if app.IsServerMode() {
		fmt.Println("webserver.healthcheck starter init")
	}
	app.RegisterStartupAction(healthcheckStartup)
}
