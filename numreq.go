package main

import (
	"fmt"
	"net/http"
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
	c.fileserverHits.Store(0)
}
