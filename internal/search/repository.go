package search

import (
	"context"
)

// SearchRepository define la interfaz para el acceso a datos del módulo de búsqueda
type SearchRepository interface {
	FindRoutes(ctx context.Context, origin, destination string) ([]*Route, error)
	ReserveRoute(ctx context.Context, routeID, userID string, seats int) (*Reservation, error)
	GetReservationByID(ctx context.Context, reservationID string) (*Reservation, error)
	MigrateDB() error
}
