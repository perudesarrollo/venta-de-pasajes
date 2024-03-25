package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"venta-de-pasajes/config"
	"venta-de-pasajes/internal/baggage"
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
	baggageTypesCollection := db.Collection(cfg.MongoDB.BaggageTypesCollection)

	// Tipos de equipaje
	baggageTypes := []*baggage.BaggageType{
		{ID: uuid.New().String(), Name: "Maleta pequeña", Price: 10.0},
		{ID: uuid.New().String(), Name: "Maleta mediana", Price: 20.0},
		{ID: uuid.New().String(), Name: "Maleta grande", Price: 30.0},
	}

	// Insertar los tipos de equipaje en la base de datos
	for _, bt := range baggageTypes {
		_, err := baggageTypesCollection.InsertOne(ctx, bt)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Tipos de equipaje insertados correctamente.")
}
