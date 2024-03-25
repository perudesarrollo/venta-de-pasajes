package search

import (
	"time"
)

// Route representa una ruta disponible para la reserva de pasajes
type Route struct {
	ID          string    `json:"id,omitempty" bson:"_id,omitempty"`
	Origin      string    `json:"origin" bson:"origin"`
	OriginCode  string    `json:"originCode" bson:"originCode"`
	Destination string    `json:"destination" bson:"destination"`
	DestCode    string    `json:"destCode" bson:"destCode"`
	Departure   time.Time `json:"departure" bson:"departure"`
	Arrival     time.Time `json:"arrival" bson:"arrival"`
	Seats       int       `json:"seats" bson:"seats"`
	Price       float64   `json:"price" bson:"price"`
}

// Reservation representa una reserva de pasajes realizada por un usuario
type Reservation struct {
	ID         string  `json:"id,omitempty" bson:"_id,omitempty"`
	RouteID    string  `json:"route_id"`
	UserID     string  `json:"user_id"`
	Seats      int     `json:"seats"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"` // Puede ser "confirmado", "pendiente", "cancelado", etc.
}
