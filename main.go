package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
	"github.com/gorilla/mux"
	"github.com/gunjan5/container-from-scratch/cmd"
)

const (
	IP   = "127.0.0.1"
	PORT = ":1337"
)

var containers = []Container{Container{"run", "TinyCore", "ls"}}

type Container struct {
	State   string
	Image   string
	Command string
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", index)
	router.HandleFunc("/run", getContainerHandler).Methods("GET")
	router.HandleFunc("/run", postContainerHandler).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(IP+PORT, nil)

	app := makeCmd()

	if len(os.Args) > 0 {
		app.Run(os.Args)
	}

}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to Container From Scratch (CFS)!")
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

	result, err := json.Marshal(containers)
	if err != nil {
		fmt.Errorf("Error marshaling json: %v ", err)
	}
	w.Write(result)

}

func makeCmd() *cli.App {

	app := cli.NewApp()
	app.Name = "CFS"
	app.Usage = "sudo ./cfs <action_command> <OS_image> <command_to_run_inside_the_container>"
	app.Version = "0.0.2"

	fmt.Println(os.Args)

	app.Commands = []cli.Command{
		// {
		// 	Name:        "server",
		// 	ShortName:   "s",
		// 	Description: "Start the REST server for CFS",
		// 	Action:      cmd.Serve,
		// },
		{
			Name:        "run",
			ShortName:   "r",
			Description: "run a container with a task",
			Action:      cmd.Run,
		},
		{
			Name:        "newroot",
			ShortName:   "n",
			Description: "Chroot. Not meant for direct usage",
			Action:      cmd.NewRoot,
		},
		{
			Name:        "child",
			ShortName:   "c",
			Description: "child process called by run, not meant for direct usage",
			Action:      cmd.Child,
		},
	}

	return app

}
