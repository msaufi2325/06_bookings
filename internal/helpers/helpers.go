package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/msaufi2325/06_bookings/internal/config"
)

var app *config.AppConfig

// NewHelpers sets up the app config for the helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func IsAuthenticated(r *http.Request) bool {
	exixts := app.Session.Exists(r.Context(), "user_id")
	return exixts
}
