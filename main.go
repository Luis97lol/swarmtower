package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

func removeDigest(imageName string) string {
	parts := strings.Split(imageName, "@")
	return parts[0] // Retorna solo la parte antes del digest
}

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

	// Obtener el servicio por nombre
	service, _, err := cli.ServiceInspectWithRaw(context.Background(), serviceName, types.ServiceInspectOptions{})
	if err != nil {
		log.Println("Error al obtener detalles del servicio:", err.Error())
		http.Error(w, "Error al obtener detalles del servicio", http.StatusInternalServerError)
		return
	}

	currentImageDigest := service.Spec.TaskTemplate.ContainerSpec.Image
	log.Println("Digest de la imagen actual:", currentImageDigest)

	// Comando para hacer pull de la imagen y actualizar el servicio
	pullCmd := exec.Command("docker", "service", "update", "--force", "--image", removeDigest(currentImageDigest), serviceName)

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
