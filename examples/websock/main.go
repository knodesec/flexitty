package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/knodesec/flexitty/examples/websock/broker"
)

var (
	LISTENIP   = "0.0.0.0"
	LISTENPORT = "8080"
	upgrader   = websocket.Upgrader{}
)

func main() {
	log.SetFlags(0)

	log.Printf("Websocket TTY Example\n")

	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/new", NewSessionHandler)
	r.HandleFunc("/session/{uuid}", TTYPageHandler)
	r.HandleFunc("/session/{uuid}/ws", WebSockHandler)
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	log.Printf("Starting server on %s:%s\n", LISTENIP, LISTENPORT)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", LISTENIP, LISTENPORT), r))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<h1>This is <b>working!</b> :)</h1><br>")
	fmt.Fprintf(w, "<a href=/new>New Session</a><br>")
}

func NewSessionHandler(w http.ResponseWriter, r *http.Request) {
	uuid := broker.NewSession()
	log.Printf("Created new session %s\n", uuid)
	http.Redirect(w, r, fmt.Sprintf("/session/%s", uuid), http.StatusMovedPermanently)
}

func TTYPageHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	if exists := broker.SessionExists(vars["uuid"]); !exists {
		log.Printf("Session %s does not exist\n", vars["uuid"])
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No such session: %s", vars["uuid"])
		return
	}

	ttyTemplate, err := template.ParseFiles("tty.html")
	if err != nil {
		panic(err)
	}

	err = ttyTemplate.ExecuteTemplate(w, "tty.html", "")
	if err != nil {
		panic(err)
	}
}

func WebSockHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error occured")
		return
	}

	err = broker.AddWS(vars["uuid"], c)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error occured")
	}
}
