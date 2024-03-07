// Package all binds all the routes into the specified app.
package all

import (
	"github.com/fadhilijuma/gateone-service/app/services/gateone-api/v1/handlers/checkgrp"
	"github.com/fadhilijuma/gateone-service/app/services/gateone-api/v1/handlers/conditiongrp"
	"github.com/fadhilijuma/gateone-service/app/services/gateone-api/v1/handlers/patientgrp"
	"github.com/fadhilijuma/gateone-service/app/services/gateone-api/v1/handlers/regiongrp"
	"github.com/fadhilijuma/gateone-service/app/services/gateone-api/v1/handlers/rolegrp"
	"github.com/fadhilijuma/gateone-service/app/services/gateone-api/v1/handlers/trangrp"
	"github.com/fadhilijuma/gateone-service/app/services/gateone-api/v1/handlers/usergrp"
	"github.com/fadhilijuma/gateone-service/business/web/v1/mux"
	"github.com/fadhilijuma/gateone-service/foundation/web"
)

// Routes constructs the add value which provides the implementation of
// of RouteAdder for specifying what routes to bind to this instance.
func Routes() add {
	return add{}
}

type add struct{}

// Add implements the RouterAdder interface.
func (add) Add(app *web.App, cfg mux.Config) {
	checkgrp.Routes(app, checkgrp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

	conditiongrp.Routes(app, conditiongrp.Config{
		Log:      cfg.Log,
		Delegate: cfg.Delegate,
		Auth:     cfg.Auth,
		DB:       cfg.DB,
	})
	patientgrp.Routes(app, patientgrp.Config{
		Log:      cfg.Log,
		Delegate: cfg.Delegate,
		Auth:     cfg.Auth,
		DB:       cfg.DB,
	})

	regiongrp.Routes(app, regiongrp.Config{
		Log:      cfg.Log,
		Delegate: cfg.Delegate,
		Auth:     cfg.Auth,
		DB:       cfg.DB,
	})
	rolegrp.Routes(app, rolegrp.Config{
		Log:      cfg.Log,
		Delegate: cfg.Delegate,
		Auth:     cfg.Auth,
		DB:       cfg.DB,
	})

	trangrp.Routes(app, trangrp.Config{
		Log:      cfg.Log,
		Delegate: cfg.Delegate,
		Auth:     cfg.Auth,
		DB:       cfg.DB,
	})

	usergrp.Routes(app, usergrp.Config{
		Log:      cfg.Log,
		Delegate: cfg.Delegate,
		Auth:     cfg.Auth,
		DB:       cfg.DB,
	})

}
