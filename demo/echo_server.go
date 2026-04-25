package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	logFile, err := os.Create("./demo/logs/on_rebuild_webhook.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/rebuild", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_, _ = logFile.WriteString(r.Method + " " + r.URL.Path + "\n")
		_, _ = logFile.WriteString(string(body) + "\n")
		_ = logFile.Sync()
		w.WriteHeader(http.StatusNoContent)
	})

	srv := &http.Server{
		Addr:    "127.0.0.1:9999",
		Handler: mux,
	}
	log.Fatal(srv.ListenAndServe())
}
