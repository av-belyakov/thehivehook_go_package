package api

import (
	"fmt"
	"net/http"

	"github.com/av-belyakov/thehivehook_go_package/internal/datamodels"
)

func (app *datamodels.Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
