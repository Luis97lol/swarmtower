package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

func updateService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := vars["servicename"] // Obtener el nombre del servicio de la ruta
	log.Println("Recibida peticion para actualizar el servicio:", serviceName)
	// Crear cliente Docker usando el socket por defecto
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		log.Println("Error al conectar con Docker:", err.Error())
		http.Error(w, "Error al conectar con Docker", http.StatusInternalServerError)
		return
	}
	defer cli.Close()

	service, _, err := cli.ServiceInspectWithRaw(context.Background(), serviceName, types.ServiceInspectOptions{})
	if err != nil {
		log.Println("Error al inspeccionar el servicio:", err.Error())
		http.Error(w, "Error al inspeccionar el servicio", http.StatusInternalServerError)
		return
	}

	// Realiza el service update
	response, err := cli.ServiceUpdate(context.Background(), service.ID, service.Version, service.Spec, types.ServiceUpdateOptions{})
	if err != nil {
		log.Println("Error al actualizar el servicio:", err.Error())
		http.Error(w, "Error al actualizar el servicio", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Servicio actualizado: %v", response.Warnings)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/update/{servicename}", updateService).Methods("POST") // Define la ruta

	log.Println("Servidor escuchando en el puerto 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
