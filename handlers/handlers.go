package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to Container From Scratch (CFS)!")
}

func GetContainerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	result, err := json.Marshal(containers)
	if err != nil {
		fmt.Errorf("Error marshaling json: %v ", err)
	}
	w.Write(result)

}
func PostContainerHandler(w http.ResponseWriter, r *http.Request) {
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
