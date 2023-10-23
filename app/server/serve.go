package server

import (
	"fmt"
	"net/http"

	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
	"github.com/abibby/eztvrss/config"
	"github.com/abibby/eztvrss/database"
	"github.com/abibby/eztvrss/routes"
)

func Serve() error {
	r := router.New()

	r.Use(request.WithDB(database.DB))

	routes.InitRoutes(r)

	return http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r)
}
