package healthcheck

import (
	"fmt"

	"github.com/abmpio/app"
)

func init() {

	fmt.Println("webserver.healthcheck starter init")
	app.RegisterStartupAction(healthcheckStartup)
}
