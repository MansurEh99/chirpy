package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	apiCfg := apiConfig {}
	apiCfg.fileserverHits.Store(0)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	//mux.Handle("/assets/logo.png", http.FileServer(http.Dir("/assets")))
	mux.HandleFunc("GET /healthz", ready)
	mux.HandleFunc("GET /metrics", apiCfg.RequestNum)
	mux.HandleFunc("POST /reset", apiCfg.Reset)
	server := http.Server {
		Addr: ":8080",
		Handler: mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
