package api

import (
	"fmt"
	"net/http"

	"github.com/av-belyakov/thehivehook_go_package/internal/datamodels"
)

func (app *datamodels.Application) backgroundTask(r *http.Request, fn func() error) {
	app.wg.Add(1)

	go func() {
		defer app.wg.Done()

		defer func() {
			err := recover()
			if err != nil {
				app.reportServerError(r, fmt.Errorf("%s", err))
			}
		}()

		err := fn()
		if err != nil {
			app.reportServerError(r, err)
		}
	}()
}
