package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gunjan5/container-from-scratch/container"
	"github.com/satori/go.uuid"
)

const (
	IP   = "127.0.0.1"
	PORT = ":1337"
)

type CID string

type Container struct {
	ID      CID    `json:"id"`
	State   string `json:"state"`
	Image   string `json:"image"`
	Command string `json:"command"`
}

var containers = map[CID]Container{}
var history = []Container{}

func MakeServer() {
	router := mux.NewRouter()
	router.HandleFunc("/", index)
	router.HandleFunc("/history", getHistoryHandler).Methods("GET")
	router.HandleFunc("/containers", getContainerHandler).Methods("GET")

	router.HandleFunc("/run", postContainerHandler).Methods("POST")

	http.Handle("/", router)

	fmt.Println(`
  ___  ____  ___       ___  _____  _  _  ____   __    ____  _  _  ____  ____    ____  ____  _____  __  __ 
 / __)( ___)/ __)()   / __)(  _  )( \( )(_  _) /__\  (_  _)( \( )( ___)(  _ \  ( ___)(  _ \(  _  )(  \/  )
( (__  )__) \__ \    ( (__  )(_)(  )  (   )(  /(__)\  _)(_  )  (  )__)  )   /   )__)  )   / )(_)(  )    ( 
 \___)(__)  (___/()   \___)(_____)(_)\_) (__)(__)(__)(____)(_)\_)(____)(_)\_)  (__)  (_)\_)(_____)(_/\/\_)
 ___   ___  ____    __   ____  ___  _   _ 
/ __) / __)(  _ \  /__\ (_  _)/ __)( )_( )
\__ \( (__  )   / /(__)\  )( ( (__  ) _ ( 
(___/ \___)(_)\_)(__)(__)(__) \___)(_) (_)

		`)
	fmt.Printf("CFS Server running at: %s%s \n\n", IP, PORT)

	http.ListenAndServe(IP+PORT, nil)

}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to Container From Scratch (CFS)!")
	fmt.Fprintf(w, "Go to %s%s/run (GET/POST) for more", IP, PORT)
}

func getContainerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	result, err := json.Marshal(containers)
	if err != nil {
		fmt.Errorf("Error marshaling json: %v ", err)
	}
	w.Write(result)

}

func getHistoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	result, err := json.Marshal(history)
	if err != nil {
		fmt.Errorf("Error marshaling json: %v ", err)
	}
	w.Write(result)

}

func postContainerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var c Container
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Errorf("Error: %v while reading request body: %v ", r.Body)
	}

	json.Unmarshal(body, &c)

	switch c.State {
	case "run":
		u := uuid.NewV4()
		c.ID = CID(u.String())
		err = container.Run([]string{c.Image, c.Command})
		if err != nil {
			fmt.Errorf("Error starting the container: %v", err)
			c.State = "Stopped: ERROR"
			//break
		}
		c.State = "Running"
		containers[c.ID] = c

	case "stop":
		//TODO: need to implement this properly
		//how do you even stop a container
		_, err = uuid.FromString(string(c.ID))
		if err != nil {
			fmt.Fprintln(w, "Put some proper container ID in your request yo!")
			return
		}
		_, ok := containers[c.ID]
		if !ok {
			fmt.Fprintln(w, "This container doesn't exist, check the Container ID")
			return
		}
		c.State = "Stopped"
		delete(containers, c.ID)

	default:
		panic("Unrecognized container state")

	}

	history = append(history, c)

	fmt.Println(c)

	result, err := json.Marshal(c)
	if err != nil {
		fmt.Errorf("Error marshaling json: %v ", err)
	}

	w.Write(result)

}
