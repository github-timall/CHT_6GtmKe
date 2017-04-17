package main

import (
	"net/http"
	"log"
	"html/template"
	"github.com/gorilla/websocket"
)

var (
	upgrader  = websocket.Upgrader{}
)

type (
	Message struct {
		AuthorID string `json:"author_id"`
		Text string `json:"text"`
	}
)

func indexHandler (w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/index.html")
	if err != nil {
		log.Println("Error parse index: ", err.Error())
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Println("Error executr index: ", err.Error())
	}
}


func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error serveWs: ", err.Error())
	}
	defer ws.Close()

	for {
		var message Message

		err := ws.ReadJSON(&message)
		if err != nil {
			log.Println("Error read json", err)
		}

		log.Printf("%+v\n", message)

		err = ws.WriteJSON(message)
		if err != nil {
			log.Println("Error write json", err)
		}
	}
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/ws", serveWs)

	err := http.ListenAndServe(":9090", nil)

	if err != nil {
		log.Println("ListenAndServe: ", err.Error())
	}
}
