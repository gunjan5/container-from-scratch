package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gunjan5/container-from-scratch/container"
)

const (
	IP   = "127.0.0.1"
	PORT = ":1337"
)

var containers = []Container{}

type Container struct {
	State   string `json:"state"`
	Image   string `json:"image"`
	Command string `json:"command"`
}

func MakeServer() {
	router := mux.NewRouter()
	router.HandleFunc("/", index)
	router.HandleFunc("/run", getContainerHandler).Methods("GET")
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
	fmt.Printf("CFS Server running at: %s%s", IP, PORT)

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
func postContainerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var c Container
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Errorf("Error: %v while reading request body: %v ", r.Body)
	}

	json.Unmarshal(body, &c)
	containers = append(containers, c)

	err = container.Run([]string{c.Image, c.Command})
	if err != nil {
		fmt.Errorf("Error starting the container: %v", err)
	}

	result, err := json.Marshal(containers)
	if err != nil {
		fmt.Errorf("Error marshaling json: %v ", err)
	}
	w.Write(result)

}
