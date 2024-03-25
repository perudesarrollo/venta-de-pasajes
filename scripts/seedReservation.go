package main

import (
	"context"
	"log"
	"math/rand"
	"strconv"
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
	reservationsCollection := db.Collection(cfg.MongoDB.ReservationsCollection)

	// Generar datos de reserva de ejemplo
	var reservations []interface{}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		// Generar un ID único UUID v4 para la reserva
		id := uuid.New()

		// Generar un ID de ruta aleatorio
		routeID := "ROUTE" + strconv.Itoa(i+1)

		// Generar un ID de usuario aleatorio
		userID := "USER" + strconv.Itoa(i+1)

		// Generar un número de asientos aleatorio entre 1 y 5
		seats := rand.Intn(5) + 1

		// Generar un precio total aleatorio entre 50 y 100
		totalPrice := float64(rand.Intn(51) + 50)

		// Generar un estado aleatorio para la reserva
		var status string
		switch rand.Intn(3) {
		case 0:
			status = "confirmado"
		case 1:
			status = "pendiente"
		case 2:
			status = "cancelado"
		}

		log.Printf("Valor de ID: %s\n", id)

		// Crear la reserva utilizando el modelo Reservation del paquete internal/search
		reservation := search.Reservation{
			ID:         uuid.New().String(),
			RouteID:    routeID,
			UserID:     userID,
			Seats:      seats,
			TotalPrice: totalPrice,
			Status:     status,
		}

		// Agregar la reserva a la lista de reservas
		reservations = append(reservations, reservation)
	}

	// Insertar las reservas en la base de datos
	_, err = reservationsCollection.InsertMany(ctx, reservations)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Datos de reserva de ejemplo insertados correctamente.")
}
