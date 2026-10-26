package main

import (
	"fmt"
	"net/http"
)

func (c *apiConfig) RequestNum(w http.ResponseWriter, r *http.Request) {
	hits := fmt.Sprintf("Hits: %d", c.fileserverHits.Load())
	w.Write([]byte(hits))
}

func (c *apiConfig) Reset(w http.ResponseWriter, r *http.Request) {
	c.fileserverHits.Store(0)
}
