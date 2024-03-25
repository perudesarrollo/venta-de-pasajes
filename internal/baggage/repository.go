package baggage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"venta-de-pasajes/config"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BaggageRepository representa un repositorio para el manejo del equipaje
type BaggageRepository struct {
	config *config.Config // Configuración de MongoDB
	client *mongo.Client
}

// NewRepository inicializa y retorna una nueva instancia de BaggageRepository
func NewRepository(cfg *config.Config) (*BaggageRepository, error) {
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

	log.Println("Conexión a MongoDB establecida, para el servicio de equipaje")

	return &BaggageRepository{
		config: cfg,
		client: client,
	}, nil
}

// CreateReservation crea una nueva reserva de equipaje en la base de datos y devuelve su ID
func (r *BaggageRepository) CreateReservation(reservation *BaggageReservation) (string, error) {
	// Generar un nuevo ID único UUID
	reservation.ID = uuid.New().String()

	// Contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Colección de reservas de equipaje
	collection := r.client.Database(r.config.MongoDB.DatabaseName).Collection(r.config.MongoDB.BaggageReservationsCollection)

	// Insertar la reserva de equipaje en la base de datos
	_, err := collection.InsertOne(ctx, reservation)
	if err != nil {
		return "erro al registrar reserva de equipaje", err
	}

	// Retornar el ID generado
	return reservation.ID, nil
}

// Asignar el precio del tipo de equipaje a la reserva
func (r *BaggageRepository) assignBaggageTypePrice(reservation *BaggageReservation) error {
	baggageType, err := r.getBaggageTypeByName(reservation.Type)
	if err != nil {
		return err
	}
	if baggageType == nil {
		return errors.New("tipo de equipaje no encontrado")
	}

	// Asignar el precio del tipo de equipaje a la reserva
	reservation.Price = baggageType.Price
	return nil
}

// Insertar la reserva en la base de datos
func (r *BaggageRepository) insertReservation(reservation *BaggageReservation) error {
	// Contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Colección de reservas de equipaje
	collection := r.client.Database(r.config.MongoDB.DatabaseName).Collection(r.config.MongoDB.BaggageReservationsCollection)

	// Generar un ID único UUID para la reserva
	reservation.ID = uuid.New().String()

	// Insertar la reserva en la base de datos
	_, err := collection.InsertOne(ctx, reservation)
	if err != nil {
		return err
	}

	return nil
}

// Obtener todos los tipos de equipaje
func (r *BaggageRepository) getAllBaggageTypes() ([]*BaggageType, error) {
	// Llamar a la función helper para realizar la búsqueda en la colección de equipajes
	return r.findBaggageTypes(context.Background(), bson.M{})
}

// Obtener el tipo de equipaje por nombre
func (r *BaggageRepository) getBaggageTypeByName(name string) (*BaggageType, error) {
	// Llamar a la función helper para realizar la búsqueda en la colección de equipajes
	filter := bson.M{"name": name}
	return r.findSingleBaggageType(context.Background(), filter)
}

// Función helper para buscar múltiples tipos de equipaje
func (r *BaggageRepository) findBaggageTypes(ctx context.Context, filter bson.M) ([]*BaggageType, error) {
	// Colección de tipos de equipaje
	collection := r.client.Database(r.config.MongoDB.DatabaseName).Collection(r.config.MongoDB.BaggageTypesCollection)

	// Realizar la búsqueda de todos los tipos de equipaje
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Variable para almacenar los tipos de equipaje encontrados
	var baggageTypes []*BaggageType

	// Iterar sobre los resultados del cursor
	for cursor.Next(ctx) {
		var baggageType BaggageType
		if err := cursor.Decode(&baggageType); err != nil {
			return nil, err
		}
		baggageTypes = append(baggageTypes, &baggageType)
	}

	// Verificar si hubo algún error durante la iteración
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return baggageTypes, nil
}

// Función helper para buscar un solo tipo de equipaje
func (r *BaggageRepository) findSingleBaggageType(ctx context.Context, filter bson.M) (*BaggageType, error) {
	// Colección de tipos de equipaje
	collection := r.client.Database(r.config.MongoDB.DatabaseName).Collection(r.config.MongoDB.BaggageTypesCollection)

	// Realizar la búsqueda del tipo de equipaje por nombre
	var baggageType BaggageType
	err := collection.FindOne(ctx, filter).Decode(&baggageType)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Tipo de equipaje no encontrado
		}
		return nil, err
	}

	return &baggageType, nil
}

// AddBaggageToReservation agrega equipaje a una reserva existente
func (r *BaggageRepository) AddBaggageToReservation(reservationID, baggageType string, quantity int) (*BaggageReservation, error) {
	// Verificar si la reserva existe
	reservation, err := r.getReservationByID(reservationID)
	if err != nil {
		return nil, err
	}
	if reservation == nil {
		return nil, errors.New("reserva no encontrada")
	}

	// Calcular el precio total del equipaje
	baggagePrice, err := r.calculateBaggagePrice(baggageType, quantity)
	if err != nil {
		return nil, err
	}

	// Actualizar la reserva con el equipaje agregado
	err = r.updateReservationWithBaggage(reservationID, baggageType, quantity, baggagePrice)
	if err != nil {
		return nil, err
	}

	// Obtener la reserva actualizada
	reservationUpdate, err := r.getReservationByID(reservationID)
	if err != nil {
		return nil, err
	}

	return reservationUpdate, nil
}

// Obtener la reserva por su ID
func (r *BaggageRepository) getReservationByID(reservationID string) (*BaggageReservation, error) {
	// Contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Colección de reservas de equipaje
	collection := r.client.Database(r.config.MongoDB.DatabaseName).Collection(r.config.MongoDB.BaggageReservationsCollection)

	// Realizar la búsqueda de la reserva por su ID
	var reservation BaggageReservation
	err := collection.FindOne(ctx, bson.M{"_id": reservationID}).Decode(&reservation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// Reserva no encontrada
			return nil, fmt.Errorf("reserva con ID %s no encontrada", reservationID)
		}
		// Otro tipo de error
		return nil, fmt.Errorf("error al buscar la reserva con ID %s: %v", reservationID, err)
	}

	return &reservation, nil
}

// Actualizar la reserva con el equipaje agregado
func (r *BaggageRepository) updateReservationWithBaggage(reservationID, baggageType string, quantity int, baggagePrice float64) error {
	// Contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Colección de reservas de equipaje
	collection := r.client.Database(r.config.MongoDB.DatabaseName).Collection(r.config.MongoDB.BaggageReservationsCollection)

	// Actualizar la reserva con el equipaje agregado
	filter := bson.M{"_id": reservationID}

	// Obtener la reserva existente
	var existingReservation BaggageReservation
	err := collection.FindOne(ctx, filter).Decode(&existingReservation)
	if err != nil {
		return err
	}

	// Combinar el equipaje existente con el nuevo equipaje
	updatedBaggage := append(existingReservation.Baggage, Baggage{Quantity: quantity, Type: baggageType})

	update := bson.M{
		"$set": bson.M{"baggage": updatedBaggage},
		"$inc": bson.M{
			"price":  baggagePrice,
			"weight": float64(quantity), // Agregar el peso del equipaje
		},
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	log.Printf("Equipaje agregado a la reserva %s: %d unidades de tipo %s", reservationID, quantity, baggageType)
	return nil
}

// CalculateBaggagePrice calcula el precio total del equipaje en función del tipo y la cantidad
func (r *BaggageRepository) calculateBaggagePrice(baggageType string, quantity int) (float64, error) {
	ctx := context.Background()
	collection := r.client.Database(r.config.MongoDB.DatabaseName).Collection(r.config.MongoDB.BaggageTypesCollection)

	var baggageTypeData BaggageType
	err := collection.FindOne(ctx, bson.M{"name": baggageType}).Decode(&baggageTypeData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Tipo de equipaje %s no encontrado", baggageType)
			return 0, nil
		}
		return 0, err
	}

	totalPrice := float64(quantity) * baggageTypeData.Price
	return totalPrice, nil
}
