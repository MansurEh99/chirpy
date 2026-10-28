package main

import (
	"fmt"
	"net/http"
	"os"
)

func (c *apiConfig) RequestNum(w http.ResponseWriter, r *http.Request) {
	hits := fmt.Sprintf(`
	<html>
  		<body>
   			<h1>Welcome, Chirpy Admin</h1>
      		<p>Chirpy has been visited %d times!</p>
        </body>
	</html>
	`, c.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(hits))
}

func (c *apiConfig) Reset(w http.ResponseWriter, r *http.Request) {
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	c.fileserverHits.Store(0)
	ctx := r.Context()
	err := c.dbQueries.DeleteUser(ctx)
	if err != nil {
		fmt.Printf("Error deleting users: %v\n", err)
		return
	}
}
