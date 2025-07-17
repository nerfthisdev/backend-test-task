package server

import (
	"fmt"
	"net/http"
)

func NewServer(port int) *http.Server {
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	return server
}
