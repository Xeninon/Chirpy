package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", healthHandler)
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func healthHandler(w http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
