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
	//IP address of the server
	IP = "127.0.0.1"
	//PORT of the server
	PORT = ":1337"
)

//CID is Container ID UUID
type CID string

//TODO: add new struct for history with timestamp (maybe some composition?)

//Container structure
type Container struct {
	ID      CID    `json:"id"`
	State   string `json:"state"`
	Image   string `json:"image"`
	Command string `json:"command"`
}

var containers = map[CID]Container{}
var history = []Container{}

//MakeServer creates the server mux and register handlers
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
	fmt.Fprintln(w, `
  ___  ____  ___       ___  _____  _  _  ____   __    ____  _  _  ____  ____    ____  ____  _____  __  __ 
 / __)( ___)/ __)()   / __)(  _  )( \( )(_  _) /__\  (_  _)( \( )( ___)(  _ \  ( ___)(  _ \(  _  )(  \/  )
( (__  )__) \__ \    ( (__  )(_)(  )  (   )(  /(__)\  _)(_  )  (  )__)  )   /   )__)  )   / )(_)(  )    ( 
 \___)(__)  (___/()   \___)(_____)(_)\_) (__)(__)(__)(____)(_)\_)(____)(_)\_)  (__)  (_)\_)(_____)(_/\/\_)
 ___   ___  ____    __   ____  ___  _   _ 
/ __) / __)(  _ \  /__\ (_  _)/ __)( )_( )
\__ \( (__  )   / /(__)\  )( ( (__  ) _ ( 
(___/ \___)(_)\_)(__)(__)(__) \___)(_) (_)

		`)

	fmt.Fprintln(w, "METHODS: ")
	fmt.Fprintf(w, "(GET) %s%s/containers\n", IP, PORT)
	fmt.Fprintf(w, "(GET) %s%s/history\n", IP, PORT)
	fmt.Fprintf(w, "(POST) %s%s/run\n", IP, PORT)
	fmt.Fprintln(w, "\n\n\nJSON structure examples:")
	fmt.Fprintln(w, `
	//Run a new container
	{
		"state": "run",
		"image": "BusyBox",
		"command": "pwd"
	}

	//Stop a running container with it's Container ID
	{
		"id": "e7887770-da8e-43db-9ca1-69526d144d7c",
		"state": "stop"
	}
		`)

	fmt.Fprintln(w, "\n\nCURL call example:")
	fmt.Fprintln(w, `curl -H "Content-Type: application/json" -X POST -d '{"state":"run","image":"TinyCore","command":"ls"}' http://localhost:1337/run`)
	fmt.Fprintln(w, `curl -H "Content-Type: application/json" -X POST -d '{"id":"d78347b9-d7c1-4e22-b2fc-782c8111cfcb","state":"stop"}' http://localhost:1337/run`)

}

func getContainerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	result, err := json.Marshal(containers)
	if err != nil {
		fmt.Printf("ERROR: Error marshaling json: %v\n", err)
	}
	w.Write(result)

}

func getHistoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	result, err := json.Marshal(history)
	if err != nil {
		fmt.Printf("ERROR: Error marshaling json: %v\n", err)
	}
	w.Write(result)

}

func postContainerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var c Container
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("ERROR: %v while reading request body: %v\n", err, r.Body)
	}

	json.Unmarshal(body, &c)

	switch c.State {
	case "run":
		u := uuid.NewV4()
		c.ID = CID(u.String())
		err = container.Run([]string{c.Image, c.Command})
		if err != nil {
			fmt.Printf("ERROR: Error starting the container: %v\n", err)
			c.State = "Stopped: ERROR"
			//break
		}
		c.State = "Running"
		containers[c.ID] = c

	case "stop":
		//TODO: need to implement this properly
		//how do you even stop a container
		// _, err = uuid.FromString(string(c.ID))
		// if err != nil {
		// 	fmt.Fprintln(w, "Put some proper container ID in your request yo!")
		// 	return
		// }
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
		fmt.Printf("ERROR: Error marshaling json: %v\n", err)
	}

	w.Write(result)

}
