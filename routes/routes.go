package routes

import (
	"github.com/abibby/eztvrss/app/handlers"
	"github.com/abibby/eztvrss/database"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
)

func InitRoutes(r *router.Router) {
	r.Use(request.HandleErrors())
	r.Use(request.WithDB(database.DB))

	r.Get("/shows/{id}/{slug}", handlers.Shows)
	r.Get("/shows/{id}/{slug}/", handlers.Shows)

	// r.Handle("/", fileserver.WithFallback(resources.Content, "dist", "index.html", nil))
}
