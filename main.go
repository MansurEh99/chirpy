package main

import _ "github.com/lib/pq"

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"os"
	"database/sql"

	"github.com/MansurEh99/chirpy/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries *database.Queries
	secretKey string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error accessing database: %v\n", err)
		os.Exit(1)
	}
	mux := http.NewServeMux()
	apiCfg := apiConfig {}
	apiCfg.fileserverHits.Store(0)
	apiCfg.dbQueries = database.New(db)
	apiCfg.secretKey = os.Getenv("SECRET_KEY")
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", ready)
	mux.HandleFunc("GET /admin/metrics", apiCfg.RequestNum)
	mux.HandleFunc("POST /admin/reset", apiCfg.Reset)
	mux.HandleFunc("POST /api/chirps", apiCfg.ChirpCreate)
	mux.HandleFunc("POST /api/users", apiCfg.UserCreate)
	mux.HandleFunc("POST /api/login", apiCfg.Login)
	mux.HandleFunc("POST /api/refresh", apiCfg.RefreshGet)
	mux.HandleFunc("POST /api/revoke", apiCfg.RefreshRevoke)
	mux.HandleFunc("GET /api/chirps", apiCfg.ChirpsGet)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.ChirpGet)
	server := http.Server {
		Addr: ":8080",
		Handler: mux,
	}
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
