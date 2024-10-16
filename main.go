package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

func updateService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := vars["servicename"] // Obtener el nombre del servicio de la ruta
	log.Println("Recibida peticion para actualizar el servicio:", serviceName)

	err := updateServiceWithNewImage(serviceName)
	if err != nil {
		log.Println("Error en el handler de actualización del servicio:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Servicio actualizado correctamente con la nueva imagen.")
}

func updateServiceWithNewImage(serviceName string) error {
	// Crear cliente Docker usando el socket por defecto
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		log.Println("Error al conectar con Docker:", err.Error())
		return fmt.Errorf("Error al conectar con Docker: %w", err)
	}
	defer cli.Close()

	// Obtener el servicio por nombre
	service, _, err := cli.ServiceInspectWithRaw(context.Background(), serviceName, types.ServiceInspectOptions{})
	if err != nil {
		log.Println("Error al obtener detalles del servicio:", err.Error())
		return fmt.Errorf("Error al obtener detalles del servicio: %w", err)
	}

	// Hacer pull de la imagen actual del servicio
	imageName := service.Spec.TaskTemplate.ContainerSpec.Image
	imagePullResponse, err := cli.ImagePull(context.Background(), imageName, image.PullOptions{})
	if err != nil {
		log.Println("Error al hacer pull de la imagen:", err.Error())
		return fmt.Errorf("Error al hacer pull de la imagen: %w", err)
	}
	defer imagePullResponse.Close()

	// Leer la respuesta para completar el pull
	if _, err := io.Copy(io.Discard, imagePullResponse); err != nil {
		log.Println("Error al leer la respuesta del pull de la imagen:", err.Error())
		return fmt.Errorf("Error al leer la respuesta del pull de la imagen: %w", err)
	}

	// Incrementar ForceUpdate para forzar la actualización de los contenedores
	service.Spec.TaskTemplate.ForceUpdate++

	// Actualizar el servicio
	_, err = cli.ServiceUpdate(context.Background(), service.ID, service.Version, service.Spec, types.ServiceUpdateOptions{})
	if err != nil {
		log.Println("Error al actualizar el servicio:", err.Error())
		return fmt.Errorf("Error al actualizar el servicio: %w", err)
	}

	log.Println("Servicio actualizado correctamente con la nueva imagen:", imageName)
	return nil
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/update/{servicename}", updateService).Methods("POST") // Define la ruta

	log.Println("Servidor escuchando en el puerto 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
