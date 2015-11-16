package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var community = make(map[int]life)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/{citizenId}", handler)
	log.Fatal(http.ListenAndServe(":3000", router))
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Printf("Error : %s\n", err)
		return
	}

	for {
		var data citizenData
		messageType, incomingMessage, readErr := conn.ReadMessage()
		// fmt.Printf("HERE : %d\n", messageType)

		unmarshalErr := json.Unmarshal(incomingMessage, &data)
		citizen := findOrCreateCitizen(data)
		citizenLife, needNewLife := community[citizen.Id]
		startLeakOfLife(citizenLife)
		if needNewLife {
			newLife := life{conn, citizen}
			community[citizen.Id] = newLife
			startLeakOfLife(newLife)
		} else {
			citizenLife.citizen = citizen
			citizenLife.connection = conn
		}

		citizenJson, jsonErr := citizen.toJson()

		if readErr != nil {
			fmt.Printf("readErr : %s\n", readErr)
			return
		}

		if unmarshalErr != nil {
			fmt.Printf("unmarshalErr : %s\n", unmarshalErr)
			return
		}

		if jsonErr != nil {
			fmt.Printf("json Error : %s\n", jsonErr)
			return
		}

		writeErr := conn.WriteMessage(messageType, citizenJson)

		if writeErr != nil {
			fmt.Printf("write Error : %s\n", writeErr)
			return
		}
	}
}

func startLeakOfLife(newLife life) {
	fmt.Printf("stuff\n")
	ticker := time.NewTicker(1000 * time.Millisecond)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				citizen := newLife.citizen
				citizen.Excitation -= 1
				citizenJson, jsonErr := citizen.toJson()
				fmt.Printf("json %s\n", citizenJson)
				if jsonErr == nil {
					newLife.connection.WriteMessage(1, citizenJson)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func findOrCreateCitizen(data citizenData) Citizen {
	life, present := community[data.Id]
	if present {
		newCitizen := life.citizen.update(data)
		return newCitizen
	}
	newCitizen := CreateCitizen(data)
	return newCitizen
}

type citizenData struct {
	Id         int
	Excitation int
	Momentum   int
}

type life struct {
	connection *websocket.Conn
	citizen    Citizen
}
