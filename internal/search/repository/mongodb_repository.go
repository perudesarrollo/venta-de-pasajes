// venta-de-pasajes/internal/search/repository/mongodb_repository.go

package repository

import (
	"context"
	"errors"
	"log"
	"reflect"
	"time"

	"venta-de-pasajes/config"
	"venta-de-pasajes/internal/search"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBRepository es una implementación de SearchRepository para MongoDB
type MongoDBRepository struct {
	config *config.Config // Configuración de MongoDB
	client *mongo.Client
}

// NewMongoDBRepository crea una nueva instancia de MongoDBRepository
func NewMongoDBRepository(cfg *config.Config) (*MongoDBRepository, error) {
	// Configurar cliente de MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.MongoDB.MongoURL))
	if err != nil {
		return nil, err
	}

	// Conectar al servidor de MongoDB
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	log.Println("Conexión a MongoDB establecida")

	return &MongoDBRepository{
		config: cfg,
		client: client,
	}, nil
}

// Implementación de los métodos de la interfaz SearchRepository
func (r *MongoDBRepository) FindRoutes(ctx context.Context, origin, destination string) ([]*search.Route, error) {
	// Contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Colección de rutas
	collection := r.client.Database(r.config.MongoDB.DatabaseName).Collection(r.config.MongoDB.RoutesCollection)

	// Filtro para la búsqueda de rutas
	filter := bson.M{"origin": origin, "destination": destination, "seats": bson.M{"$gt": 0}}

	log.Printf("filter: %s\n", filter)

	// Realizar la búsqueda de rutas en la colección
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Variable para almacenar las rutas encontradas
	var routes []*search.Route

	// Iterar sobre los resultados del cursor
	for cursor.Next(ctx) {
		var route search.Route
		if err := cursor.Decode(&route); err != nil {
			return nil, err
		}
		routes = append(routes, &route)
	}

	// Verificar si hubo algún error durante la iteración
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return routes, nil
}

func (r *MongoDBRepository) ReserveRoute(ctx context.Context, routeID, userID string, seats int) (*search.Reservation, error) {
	// Contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Colección de rutas
	collection := r.client.Database(r.config.MongoDB.DatabaseName).Collection(r.config.MongoDB.RoutesCollection)

	// Obtener la ruta
	var route search.Route
	err := collection.FindOne(ctx, bson.M{"_id": routeID, "seats": bson.M{"$gte": seats}}).Decode(&route)
	if err != nil {
		return nil, errors.New("ruta no encontrada o no hay suficientes asientos disponibles")
	}

	// Actualizar la cantidad de asientos disponibles
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": routeID},
		bson.M{"$inc": bson.M{"seats": -seats}},
	)
	if err != nil {
		return nil, err
	}

	// Colección de reservas
	reservationsCollection := r.client.Database(r.config.MongoDB.DatabaseName).Collection(r.config.MongoDB.ReservationsCollection)

	// Generar un ID único UUID v4 para la reserva
	reservationID := uuid.New().String()

	// Crear la reserva
	reservation := &search.Reservation{
		ID:         reservationID,
		RouteID:    routeID,
		UserID:     userID,
		Seats:      seats,
		TotalPrice: float64(seats) * route.Price,
		Status:     "confirmado", // Se puede cambiar según la lógica de negocio
	}

	// Insertar la reserva en la colección de reservas
	_, err = reservationsCollection.InsertOne(ctx, reservation)
	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (r *MongoDBRepository) GetReservationByID(ctx context.Context, reservationID string) (*search.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.client.Database(r.config.MongoDB.DatabaseName).Collection(r.config.MongoDB.ReservationsCollection)

	var reservation search.Reservation
	err := collection.FindOne(ctx, bson.M{"_id": reservationID}).Decode(&reservation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Reserva con ID %s no encontrada", reservationID)
			return nil, nil
		}
		return nil, err
	}

	return &reservation, nil
}

// MigrateDB crea o migra las colecciones necesarias en la base de datos MongoDB
func (r *MongoDBRepository) MigrateDB() error {
	// Crear o migrar colección para el modelo Route
	if err := r.createIndexIfNotExists(
		r.config.MongoDB.RoutesCollection,
		bson.D{{Key: "origin", Value: 1}, {Key: "destination", Value: 1}},
	); err != nil {
		return err
	}

	// Crear o migrar colección para el modelo Reservation
	if err := r.createIndexIfNotExists(
		r.config.MongoDB.ReservationsCollection,
		bson.D{{Key: "route_id", Value: 1}},
	); err != nil {
		return err
	}

	return nil
}

// createIndexIfNotExists verifica si el índice dado existe en la colección dada y lo crea si no existe
func (r *MongoDBRepository) createIndexIfNotExists(collectionName string, keys bson.D) error {
	// Obtener el nombre de la base de datos desde la configuración
	dbName := r.config.MongoDB.DatabaseName

	// Obtener la colección
	collection := r.client.Database(dbName).Collection(collectionName)

	// Verificar si el índice ya existe
	indexExists, err := r.indexExists(collection, keys)
	if err != nil {
		return err
	}

	// Si el índice no existe, crearlo
	if !indexExists {
		_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
			Keys:    keys,
			Options: options.Index().SetUnique(true),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// Función para verificar si el índice ya existe
func (r *MongoDBRepository) indexExists(collection *mongo.Collection, keys bson.D) (bool, error) {
	ctx := context.Background()
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		return false, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var index bson.D
		if err := cursor.Decode(&index); err != nil {
			return false, err
		}
		if reflect.DeepEqual(index, keys) {
			return true, nil
		}
	}

	if err := cursor.Err(); err != nil {
		return false, err
	}

	return false, nil
}
