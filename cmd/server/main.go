package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type Request struct {
	Message string
}

type Server struct{}

func (Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	var data Request
	if err := json.Unmarshal(body, &data); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(data.Message))
}

func main() {
	var s Server

	err := http.ListenAndServe("localhost:57309", s)
	if err != nil {
		panic(err)
	}
}
