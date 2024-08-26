package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/av-belyakov/thehivehook_go_package/internal/datamodels"
)

func (app *datamodels.Application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.NotFound(app.notFound)
	mux.MethodNotAllowed(app.methodNotAllowed)

	mux.Use(app.recoverPanic)

	mux.Get("/status", app.status)

	return mux
}
