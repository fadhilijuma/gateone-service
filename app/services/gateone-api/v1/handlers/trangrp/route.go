package trangrp

import (
	"github.com/fadhilijuma/gateone-service/business/core/crud/delegate"
	"github.com/fadhilijuma/gateone-service/business/core/crud/patient"
	"github.com/fadhilijuma/gateone-service/business/core/crud/patient/stores/patientdb"
	"github.com/fadhilijuma/gateone-service/business/core/crud/user"
	"github.com/fadhilijuma/gateone-service/business/core/crud/user/stores/usercache"
	"github.com/fadhilijuma/gateone-service/business/core/crud/user/stores/userdb"
	"github.com/fadhilijuma/gateone-service/business/data/sqldb"
	"github.com/fadhilijuma/gateone-service/business/web/v1/auth"
	"github.com/fadhilijuma/gateone-service/business/web/v1/mid"
	"github.com/fadhilijuma/gateone-service/foundation/logger"
	"github.com/fadhilijuma/gateone-service/foundation/web"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log      *logger.Logger
	Delegate *delegate.Delegate
	Auth     *auth.Auth
	DB       *sqlx.DB
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	usrCore := user.NewCore(cfg.Log, cfg.Delegate, usercache.NewStore(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB)))
	pnCore := patient.NewCore(cfg.Log, usrCore, cfg.Delegate, patientdb.NewStore(cfg.Log, cfg.DB))

	authen := mid.Authenticate(cfg.Auth)
	tran := mid.ExecuteInTransaction(cfg.Log, sqldb.NewBeginner(cfg.DB))

	hdl := new(usrCore, pnCore)
	app.Handle(http.MethodPost, version, "/tranexample", hdl.create, authen, tran)
}
