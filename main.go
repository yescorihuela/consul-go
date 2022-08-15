package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/consul/api"
)

func main() {
	port := flag.Int("port", 8080, "port to address requests")
	flag.Parse()
	ServiceRegistryInConsul(*port)
	http.HandleFunc("/status", health)
	http.HandleFunc("/hello", hello)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	t := struct {
		Status string `json:"status"`
		Code   int    `json:"code"`
	}{
		Status: "ok",
		Code:   200,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	t := struct {
		Message string
	}{
		Message: "Hello world",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func ServiceRegistryInConsul(port int) {
	consul, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Println(err)
	}

	serviceId := fmt.Sprintf("first-service-demo_%d", port)
	address := "localhost"

	registration := &api.AgentServiceRegistration{
		ID:      serviceId,
		Name:    "proof-of-concept-consul",
		Port:    port,
		Address: address,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/status", address, port),
			Interval: "2s",
			Timeout:  "20s",
		},
	}

	err = consul.Agent().ServiceRegister(registration)
	if err != nil {
		log.Printf("Failed to register service: %s:%d", address, port)
	} else {
		log.Printf("Service registration has occured successfully at: %s:%d", address, port)
	}

}
