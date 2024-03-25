package baggage

// Baggage representa un objeto de equipaje
type Baggage struct {
	Quantity int    `json:"quantity" bson:"quantity"`
	Type     string `json:"type" bson:"type"`
}

// BaggageReservation representa la información de reserva de equipaje
type BaggageReservation struct {
	ID            string    `json:"id,omitempty" bson:"_id,omitempty"`
	ReservationID string    `json:"reservation_id" bson:"reservation_id"`
	Weight        float64   `json:"weight" bson:"weight"`
	Price         float64   `json:"price" bson:"price"`
	Type          string    `json:"type" bson:"type"`
	Baggage       []Baggage `json:"baggage" bson:"baggage"`
}

// BaggageType representa los tipos de equipaje disponibles
type BaggageType struct {
	ID    string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string  `json:"name" bson:"name"`
	Price float64 `json:"price" bson:"price"`
}
