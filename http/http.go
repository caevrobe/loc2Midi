package http

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var dest chan string

var upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Receive(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			return
		}
		dest <- string(msg) // msg is coords in format: coord_x, coord_y
	}
}

func noCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=0")
		h.ServeHTTP(w, r)
	})
}

func ListenAndServe(port int, destination chan string) {
	dest = destination

	r := mux.NewRouter()
	r.HandleFunc("/socket", Receive)
	r.Use(noCache)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	srv := &http.Server{
		Handler: handlers.LoggingHandler(os.Stdout, r),
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
	}

	panic(srv.ListenAndServe())
}
