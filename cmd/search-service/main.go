package main

import (
	"log"
	"net/http"

	"venta-de-pasajes/config"
	"venta-de-pasajes/internal/search"
	"venta-de-pasajes/internal/search/repository"
)

func main() {
	// Configurar la carga de configuración
	cfg := config.NewConfig()

	// Inicializar el repositorio de búsqueda para MongoDB
	searchRepo, err := repository.NewMongoDBRepository(cfg)
	if err != nil {
		log.Fatalf("Error al inicializar el repositorio de MongoDB: %v", err)
	}

	// Inicializar el manejador de búsqueda
	searchHandler := search.NewSearchHandler(searchRepo)

	// Configurar rutas de búsqueda
	http.HandleFunc("/search", searchHandler.SearchRoutesHandler)
	http.HandleFunc("/reserve", searchHandler.ReserveRouteHandler)

	// Realizar la migración de la base de datos
	// err = searchHandler.MigrateDBHandler()
	// if err != nil {
	// 	log.Fatalf("Error al migrar la base de datos: %v", err)
	// }

	// Configurar el servidor HTTP para que escuche en un puerto específico
	serverAddr := ":8080" // Puerto al que HAProxy redirigirá las solicitudes
	server := &http.Server{
		Addr:    serverAddr,
		Handler: nil, // Usa el enrutador predeterminado (http.DefaultServeMux)
	}

	// Iniciar el servidor HTTP
	log.Printf("Servidor iniciado en el puerto %s...", serverAddr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error al iniciar el servidor HTTP: %v", err)
	}
}
