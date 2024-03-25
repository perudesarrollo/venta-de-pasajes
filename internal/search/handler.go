package search

import (
	"encoding/json"
	"log"
	"net/http"
)

// SearchHandler maneja las solicitudes relacionadas con la búsqueda y reserva de rutas
type SearchHandler struct {
	repo SearchRepository // Cambiado de *SearchRepository
}

// NewSearchHandler crea una nueva instancia de SearchHandler
func NewSearchHandler(repo SearchRepository) *SearchHandler { // Cambiado de *SearchRepository
	return &SearchHandler{
		repo: repo,
	}
}

// SearchRoutesHandler maneja las solicitudes para buscar rutas disponibles entre un origen y un destino
func (h *SearchHandler) SearchRoutesHandler(w http.ResponseWriter, r *http.Request) {
	origin := r.URL.Query().Get("origin")
	destination := r.URL.Query().Get("destination")

	log.Println("origen %s\n", origin)
	log.Println("destino %s\n", destination)

	routes, err := h.repo.FindRoutes(r.Context(), origin, destination) // Pasamos el contexto
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(routes)
}

// ReserveRouteHandler maneja las solicitudes para reservar una ruta
func (h *SearchHandler) ReserveRouteHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		RouteID string `json:"route_id"`
		UserID  string `json:"user_id"`
		Seats   int    `json:"seats"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	reservation, err := h.repo.ReserveRoute(r.Context(), requestBody.RouteID, requestBody.UserID, requestBody.Seats) // Pasamos el contexto
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservation)
}

// MigrateDBHandler maneja las solicitudes para migrar la base de datos
func (h *SearchHandler) MigrateDBHandler() error {
	err := h.repo.MigrateDB()
	if err != nil {
		log.Printf("Error en migración: %s\n", err)
		return err
	}
	log.Println("Migración de base de datos completada exitosamente")
	return nil
}
