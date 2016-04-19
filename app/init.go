package app

import (
	"fmt"
	"time"
	"github.com/revel/revel"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	revel.OnAppStart(InitDB)
	revel.OnAppStart(Pinger.StartPolling)
}

var DB *sql.DB
var Pinger *ServerPinger = NewServerPinger()

func InitDB() {
	DB, _ = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", "root", "password", "ismyhostup"))

	// Keep at least one connection open at all times
	go func() {
		for {
			DB.Ping()
			time.Sleep(15 * time.Second)
		}
	}()
}

var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("Connection", "keep-alive")
	fc[0](c, fc[1:]) // Execute the next filter stage.
}
