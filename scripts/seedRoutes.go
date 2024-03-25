package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"venta-de-pasajes/config"
	"venta-de-pasajes/internal/search"
)

func main() {
	// Obtener la configuración desde el paquete config
	cfg := config.NewConfig()

	// Configurar cliente de MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.MongoDB.MongoURL))
	if err != nil {
		log.Fatal(err)
	}

	// Conectar al servidor de MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Seleccionar la base de datos y la colección
	db := client.Database(cfg.MongoDB.DatabaseName)
	routesCollection := db.Collection(cfg.MongoDB.RoutesCollection)

	// Lista de ciudades y sus códigos en el Perú
	peruCities := map[string]string{
		"Lima":     "LIM",
		"Cusco":    "CUZ",
		"Arequipa": "AQP",
		"Trujillo": "TRU",
		"Iquitos":  "IQT",
		"Piura":    "PIU",
		"Tacna":    "TCQ",
		"Pucallpa": "PCL",
	}

	// Generar todas las combinaciones de pares de ciudades
	var routes []search.Route
	generatedRoutes := make(map[string]bool)

	for originCity, originCode := range peruCities {
		for destinationCity, destinationCode := range peruCities {
			if originCity != destinationCity {
				routeID := uuid.New().String() // Generar un ID único UUID v4

				// Generar claves para verificar la duplicación
				forwardRouteKey := originCity + "-" + destinationCity
				backwardRouteKey := destinationCity + "-" + originCity

				// Verificar si ya se generó la ruta o su inversa
				if !generatedRoutes[forwardRouteKey] && !generatedRoutes[backwardRouteKey] {
					// Crear la ruta y agregarla a la lista
					route := search.Route{
						ID:          routeID,
						Origin:      originCity,
						OriginCode:  originCode,
						Destination: destinationCity,
						DestCode:    destinationCode,
						Departure:   time.Now().Add(24 * time.Hour),
						Arrival:     time.Now().Add(26 * time.Hour),
						Seats:       100,
						Price:       50.0,
					}
					routes = append(routes, route)

					// Marcar la ruta y su inversa como generadas
					generatedRoutes[forwardRouteKey] = true
					generatedRoutes[backwardRouteKey] = true
				}
			}
		}
	}

	// Insertar todas las rutas en la base de datos
	for _, route := range routes {
		_, err := routesCollection.InsertOne(ctx, route)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Datos de ejemplo insertados correctamente.")
}
