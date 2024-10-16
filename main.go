package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
)

func updateService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := vars["servicename"] // Obtener el nombre del servicio de la ruta
	log.Println("Recibida peticion para actualizar el servicio:", serviceName)

	// Comando para hacer pull de la imagen y actualizar el servicio
	pullCmd := exec.Command("docker", "service", "update", "--force", serviceName)

	// Ejecutar el comando y obtener la salida
	output, err := pullCmd.CombinedOutput()
	if err != nil {
		log.Printf("Error al actualizar el servicio: %v, salida: %s", err, output)
		http.Error(w, "Error al actualizar el servicio", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Servicio actualizado: %s", output)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/update/{servicename}", updateService).Methods("POST") // Define la ruta

	log.Println("Servidor escuchando en el puerto 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
