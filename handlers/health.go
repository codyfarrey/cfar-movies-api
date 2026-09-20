package handlers

import (
	"fmt"
	"net/http"
)

func HealthHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(rw, "OK")
}
