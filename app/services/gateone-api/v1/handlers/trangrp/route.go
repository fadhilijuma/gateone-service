package trangrp

import (
	"github.com/fadhilijuma/gateone-service/business/core/crud/delegate"
	"github.com/fadhilijuma/gateone-service/business/core/crud/product"
	"github.com/fadhilijuma/gateone-service/business/core/crud/product/stores/productdb"
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
	prdCore := product.NewCore(cfg.Log, usrCore, cfg.Delegate, productdb.NewStore(cfg.Log, cfg.DB))

	authen := mid.Authenticate(cfg.Auth)
	tran := mid.ExecuteInTransaction(cfg.Log, sqldb.NewBeginner(cfg.DB))

	hdl := new(usrCore, prdCore)
	app.Handle(http.MethodPost, version, "/tranexample", hdl.create, authen, tran)
}
