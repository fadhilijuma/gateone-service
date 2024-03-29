package patientgrp

import (
	"github.com/fadhilijuma/gateone-service/business/core/crud/delegate"
	"github.com/fadhilijuma/gateone-service/business/core/crud/patient"
	"github.com/fadhilijuma/gateone-service/business/core/crud/patient/stores/patientdb"
	"github.com/fadhilijuma/gateone-service/business/core/crud/user"
	"github.com/fadhilijuma/gateone-service/business/core/crud/user/stores/usercache"
	"github.com/fadhilijuma/gateone-service/business/core/crud/user/stores/userdb"
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
	prdCore := patient.NewCore(cfg.Log, usrCore, cfg.Delegate, patientdb.NewStore(cfg.Log, cfg.DB))

	authen := mid.Authenticate(cfg.Auth)
	ruleAny := mid.Authorize(cfg.Auth, auth.RuleAny)
	ruleUserOnly := mid.Authorize(cfg.Auth, auth.RuleUserOnly)
	ruleAdminOrSubject := mid.AuthorizePatient(cfg.Auth, auth.RuleAdminOrSubject, prdCore)

	hdl := new(prdCore, usrCore)
	app.Handle(http.MethodGet, version, "/patients", hdl.query, authen, ruleAny)
	app.Handle(http.MethodGet, version, "/patients/{patient_id}", hdl.queryByID, authen, ruleAdminOrSubject)
	app.Handle(http.MethodPost, version, "/patients", hdl.create, authen, ruleUserOnly)
	app.Handle(http.MethodPut, version, "/patients/{patient_id}", hdl.update, authen, ruleAdminOrSubject)
	app.Handle(http.MethodDelete, version, "/patients/{patient_id}", hdl.delete, authen, ruleAdminOrSubject)
}
