package main

import (
	"github.com/lomtom/image_push/pkg/harbor"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/upload", uploadHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.FormValue("address") == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	skipTls := false
	skipTlsParam := r.FormValue("skipTls")
	if skipTlsParam != "" {
		skipTls, err = strconv.ParseBool(skipTlsParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	chunkSize := 0
	chunkSizeParam := r.FormValue("chunkSize")
	if chunkSizeParam != "" {
		chunkSize, err = strconv.Atoi(chunkSizeParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	c, err := harbor.NewConfig(r.FormValue("address"), r.FormValue("username"), r.FormValue("password"), r.FormValue("project"), fileHeader.Filename, skipTls, chunkSize, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = c.Push(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
