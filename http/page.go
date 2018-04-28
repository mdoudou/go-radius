package http

import (
	"io"
	"net/http"
)

func configPageRoutes() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello, I'm GoRadius (｡A｡)")
	})

}
