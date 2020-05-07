package server

import (
	"io/ioutil"
	"logserver/config"
	"logserver/util"
	"net/http"
	"time"
)

type HTTPServer struct {
	conn *http.Server
	Log
}

func (s *HTTPServer) Listen() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		//w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		//w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if req.Method == "POST" {
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				util.Printf("HTTP Read error:%s", err)
				return
			}

			s.logs <- b
		}
	})

	s.conn = &http.Server{
		Addr:           config.Server.Address,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	defer s.Close()

	go s.Write()

	err := s.conn.ListenAndServe()
	util.Fatal(err)
}

func (s *HTTPServer) Close() {
	err := s.conn.Close()
	util.Fatal(err)
}
