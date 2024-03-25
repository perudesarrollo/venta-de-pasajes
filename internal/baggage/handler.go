package baggage

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

// BaggageHandler contiene los métodos HTTP para el manejo de equipaje
type BaggageHandler struct {
	repo *BaggageRepository
}

// NewBaggageHandler crea una nueva instancia de BaggageHandler con el repositorio proporcionado
func NewBaggageHandler(repo *BaggageRepository) *BaggageHandler {
	return &BaggageHandler{repo: repo}
}

// handleError es una función de utilidad para manejar errores y escribir una respuesta HTTP de error
func (h *BaggageHandler) handleError(w http.ResponseWriter, err error, status int) {
	http.Error(w, err.Error(), status)
}

func (h *BaggageHandler) CreateReservationBaggageHandler(w http.ResponseWriter, r *http.Request) {
	// Decodificar la solicitud JSON en una estructura de reserva de equipaje
	var reservation BaggageReservation
	err := json.NewDecoder(r.Body).Decode(&reservation)
	if err != nil {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}

	// Llamar a la función del repositorio para crear la reserva y obtener su ID
	reservationID, err := h.repo.CreateReservation(&reservation)
	if err != nil {
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	// Escribir la respuesta con el ID de la reserva creada
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		ReservationID string `json:"baggage_reservation_id"`
	}{
		ReservationID: reservationID,
	})
}

// GetBaggageTypesByNameBaggageHandler maneja la obtención de tipos de equipaje por nombre
func (h *BaggageHandler) GetBaggageTypesByNameBaggageHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el nombre del tipo de equipaje de los parámetros de la URL
	name := r.URL.Query().Get("name")

	// Si el nombre está vacío, obtiene todas las colecciones
	if name == "" {
		types, err := h.repo.getAllBaggageTypes()
		if err != nil {
			h.handleError(w, err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(types)
		return
	}

	// Llamar a la función del repositorio para obtener los tipos de equipaje por nombre
	types, err := h.repo.getBaggageTypeByName(name)
	if err != nil {
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	// Escribir la respuesta con los tipos de equipaje encontrados
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(types)
}

// parseQuantity convierte la cantidad de equipaje de string a int y maneja los errores
func (h *BaggageHandler) parseQuantity(quantityStr string) (int, error) {
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		return 0, err
	}
	return quantity, nil
}

// AddBaggageToReservationBaggageHandler maneja la adición de equipaje a una reserva existente
func (h *BaggageHandler) AddBaggageToReservationBaggageHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BaggageReservationID string `json:"baggage_reservation_id"`
		BaggageType          string `json:"baggage_type"`
		Quantity             int    `json:"quantity"`
	}

	log.Printf("Valor de BaggageReservationID: %s\n", req.BaggageReservationID)
	log.Printf("Valor de BaggageType: %s\n", req.BaggageType)
	log.Printf("Valor de Quantity: %s\n", req.Quantity)

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.handleError(w, errors.New("error al decodificar la solicitud"), http.StatusBadRequest)
		return
	}

	// Validar los campos de la solicitud
	if req.BaggageReservationID == "" || req.BaggageType == "" || req.Quantity <= 0 {
		h.handleError(w, errors.New("los campos baggage_reservation_id, baggage_type y quantity son obligatorios"), http.StatusBadRequest)
		return
	}

	// Llamar a la función del repositorio para agregar equipaje a la reserva
	insertedBaggage, err := h.repo.AddBaggageToReservation(req.BaggageReservationID, req.BaggageType, req.Quantity)
	if err != nil {
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	// Escribir la respuesta de éxito
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(insertedBaggage)
}

// CalculateBaggagePriceBaggageHandler maneja el cálculo del precio del equipaje
func (h *BaggageHandler) CalculateBaggagePriceBaggageHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener los parámetros de la solicitud
	baggageType := r.URL.Query().Get("baggage_type")
	quantityStr := r.URL.Query().Get("quantity")

	// Convertir la cantidad de equipaje a entero
	quantity, err := h.parseQuantity(quantityStr)
	if err != nil {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}

	// Llamar a la función del repositorio para calcular el precio del equipaje
	price, err := h.repo.calculateBaggagePrice(baggageType, quantity)
	if err != nil {
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	// Escribir la respuesta con el precio calculado
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"price": price})
}
