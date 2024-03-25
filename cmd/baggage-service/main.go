package main

import (
	"log"
	"net/http"

	"venta-de-pasajes/config"
	"venta-de-pasajes/internal/baggage"
)

func main() {
	// Configurar la carga de configuración
	cfg := config.NewConfig()

	// Inicializar el repositorio de equipaje
	baggageRepo, err := baggage.NewRepository(cfg)
	if err != nil {
		log.Fatal("Error al crear el repositorio de equipaje: ", err)
	}

	// Inicializar el manejador de equipaje
	baggageHandler := baggage.NewBaggageHandler(baggageRepo)

	// Configurar rutas de equipaje
	http.HandleFunc("/baggage/reserve", baggageHandler.AddBaggageToReservationBaggageHandler)
	http.HandleFunc("/baggage/types", baggageHandler.GetBaggageTypesByNameBaggageHandler)
	http.HandleFunc("/baggage/add", baggageHandler.CreateReservationBaggageHandler)

	// Configurar el servidor HTTP para que escuche en un puerto específico
	serverAddr := ":8081" // Puerto al que HAProxy redirigirá las solicitudes
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
