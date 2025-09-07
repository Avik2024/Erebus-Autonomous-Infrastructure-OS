package version

import (
	"fmt"
	"net/http"
)

const Version = "v0.0.1"

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, Version)
}
