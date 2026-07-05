package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ============================================================
// 1. МОДЕЛИ ДАННЫХ
// ============================================================

type Route struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"` // novice / experienced
	CapacityCap int    `json:"capacity_cap"`
	DurationMin int    `json:"duration_min"`
}

type Instructor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Slot struct {
	ID                string      `json:"id"`
	StartAt           time.Time   `json:"start_at"`
	Route             Route       `json:"route"`
	Instructor        Instructor  `json:"instructor"`
	TotalSeats        int         `json:"total_seats"`
	FreeSeats         int         `json:"free_seats"`
	FreeRentalBoards  int         `json:"free_rental_boards"`
	Price             int         `json:"price"`
	RentalPrice       int         `json:"rental_price"`
	MeetingPoint      string      `json:"meeting_point"`
	MeetingPointLat   float64     `json:"meeting_point_lat"`
	MeetingPointLng   float64     `json:"meeting_point_lng"`
	Status            string      `json:"status"` // scheduled / cancelled
}

type Booking struct {
	ID                 string     `json:"id"`
	SlotID             string     `json:"slot_id"`
	ClientID           string     `json:"client_id"`
	SeatsCount         int        `json:"seats_count"`
	RentalCount        int        `json:"rental_count"`
	Status             string     `json:"status"` // active / cancelled / late_cancel / club_cancelled
	PriceTotal         int        `json:"price_total"`
	CreatedAt          time.Time  `json:"created_at"`
	CancelledAt        *time.Time `json:"cancelled_at,omitempty"`
	CancellationReason string     `json:"cancellation_reason,omitempty"`
	Slot               Slot       `json:"slot,omitempty"`
}

type Client struct {
	ID        string    `json:"id"`
	Phone     string    `json:"phone"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ============================================================
// 2. ХРАНИЛИЩЕ (В ПАМЯТИ)
// ============================================================

var (
	clients     = map[string]Client{}
	slots       = map[string]Slot{}
	bookings    = map[string]Booking{}
	otpCodes    = map[string]string{} // phone -> code
	idempotency = map[string]string{} // key -> booking_id
)

func init() {
	// Инструкторы
	instructors := map[string]Instructor{
		"inst-1": {ID: "inst-1", Name: "Александр Иванов"},
		"inst-2": {ID: "inst-2", Name: "Мария Петрова"},
	}

	// Маршруты
	routes := map[string]Route{
		"route-1": {
			ID:          "route-1",
			Name:        "Скалодром «Высота»",
			Description: "Новичковый маршрут с инструктажем",
			Type:        "novice",
			CapacityCap: 8,
			DurationMin: 90,
		},
		"route-2": {
			ID:          "route-2",
			Name:        "Трасса с веревкой",
			Description: "Для опытных скалолазов",
			Type:        "experienced",
			CapacityCap: 12,
			DurationMin: 120,
		},
	}

	// Слоты
	now := time.Now().Truncate(time.Hour)
	slots = map[string]Slot{
		"slot-1": {
			ID:               "slot-1",
			StartAt:          now.Add(2 * time.Hour),
			Route:            routes["route-1"],
			Instructor:       instructors["inst-1"],
			TotalSeats:       8,
			FreeSeats:        5,
			FreeRentalBoards: 4,
			Price:            1500,
			RentalPrice:      500,
			MeetingPoint:     "Скалодром «Вертикаль», ул. Спортивная, д. 10",
			MeetingPointLat:  59.978,
			MeetingPointLng:  30.262,
			Status:           "scheduled",
		},
		"slot-2": {
			ID:               "slot-2",
			StartAt:          now.Add(4 * time.Hour),
			Route:            routes["route-2"],
			Instructor:       instructors["inst-2"],
			TotalSeats:       12,
			FreeSeats:        8,
			FreeRentalBoards: 6,
			Price:            2000,
			RentalPrice:      500,
			MeetingPoint:     "Скалодром «Вертикаль», ул. Спортивная, д. 10",
			MeetingPointLat:  59.978,
			MeetingPointLng:  30.262,
			Status:           "scheduled",
		},
	}

	// Тестовый клиент
	clients["client-1"] = Client{
		ID:        "client-1",
		Phone:     "+79123456789",
		Name:      "Тестовый Клиент",
		CreatedAt: time.Now(),
	}
	otpCodes["+79123456789"] = "123456"
}

// ============================================================
// 3. JWT
// ============================================================

var jwtSecret = []byte("climbing-gym-secret-key-change-me")

func generateToken(clientID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"client_id": clientID,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString(jwtSecret)
}

func validateToken(tokenString string) (string, error) {
	if tokenString == "" {
		return "", fmt.Errorf("no token")
	}
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}
	clientID, ok := claims["client_id"].(string)
	if !ok {
		return "", fmt.Errorf("client_id not found")
	}
	return clientID, nil
}

// ============================================================
// 4. МИДЛВАРЫ
// ============================================================

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"code":"unauthorized","message":"Требуется авторизация"}`, http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		clientID, err := validateToken(token)
		if err != nil {
			http.Error(w, `{"code":"unauthorized","message":"Неверный токен"}`, http.StatusUnauthorized)
			return
		}
		r.Header.Set("X-Client-ID", clientID)
		next(w, r)
	}
}

func jsonResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func errorResponse(w http.ResponseWriter, code, message string, status int) {
	jsonResponse(w, map[string]string{"code": code, "message": message}, status)
}

// ============================================================
// 5. ХЕНДЛЕРЫ
// ============================================================

func handleVerifyCode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, "bad_request", "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	expected, ok := otpCodes[req.Phone]
	if !ok || expected != req.Code {
		errorResponse(w, "invalid_code", "Неверный код", http.StatusBadRequest)
		return
	}

	client, ok := clients[req.Phone]
	if !ok {
		client = Client{
			ID:        uuid.New().String(),
			Phone:     req.Phone,
			Name:      "",
			CreatedAt: time.Now(),
		}
		clients[req.Phone] = client
	}

	token, err := generateToken(client.ID)
	if err != nil {
		errorResponse(w, "internal_error", "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]interface{}{
		"tokens": map[string]interface{}{
			"access_token":  token,
			"refresh_token": token,
			"token_type":    "Bearer",
			"expires_in":    86400,
		},
		"client": client,
		"is_new": client.Name == "",
	}, http.StatusOK)
}

func handleListSlots(w http.ResponseWriter, r *http.Request) {
	items := []Slot{}
	for _, s := range slots {
		items = append(items, s)
	}
	jsonResponse(w, map[string]interface{}{
		"items": items,
		"meta": map[string]int{
			"limit":  20,
			"offset": 0,
			"total":  len(items),
		},
	}, http.StatusOK)
}

func handleCreateBooking(w http.ResponseWriter, r *http.Request) {
	clientID := r.Header.Get("X-Client-ID")
	if clientID == "" {
		errorResponse(w, "unauthorized", "Не авторизован", http.StatusUnauthorized)
		return
	}

	var req struct {
		SlotID      string `json:"slot_id"`
		SeatsCount  int    `json:"seats_count"`
		RentalCount int    `json:"rental_count"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, "bad_request", "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	idempKey := r.Header.Get("Idempotency-Key")
	if idempKey != "" {
		if existingBookingID, ok := idempotency[idempKey]; ok {
			if booking, ok := bookings[existingBookingID]; ok {
				jsonResponse(w, booking, http.StatusCreated)
				return
			}
		}
	}

	slot, ok := slots[req.SlotID]
	if !ok {
		errorResponse(w, "not_found", "Слот не найден", http.StatusNotFound)
		return
	}
	if slot.Status == "cancelled" {
		errorResponse(w, "slot_cancelled", "Тренировка отменена", http.StatusGone)
		return
	}
	if slot.FreeSeats < req.SeatsCount {
		errorResponse(w, "slot_full", fmt.Sprintf("Свободно только %d мест", slot.FreeSeats), http.StatusConflict)
		return
	}
	if slot.FreeRentalBoards < req.RentalCount {
		errorResponse(w, "rental_unavailable", fmt.Sprintf("Доступно только %d комплектов", slot.FreeRentalBoards), http.StatusConflict)
		return
	}
	if req.SeatsCount < 1 || req.SeatsCount > 3 {
		errorResponse(w, "bad_request", "Количество мест должно быть от 1 до 3", http.StatusBadRequest)
		return
	}
	if req.RentalCount < 0 || req.RentalCount > req.SeatsCount {
		errorResponse(w, "bad_request", "Прокат не может превышать количество мест", http.StatusBadRequest)
		return
	}

	bookingID := uuid.New().String()
	priceTotal := slot.Price*req.SeatsCount + slot.RentalPrice*req.RentalCount

	booking := Booking{
		ID:          bookingID,
		SlotID:      req.SlotID,
		ClientID:    clientID,
		SeatsCount:  req.SeatsCount,
		RentalCount: req.RentalCount,
		Status:      "active",
		PriceTotal:  priceTotal,
		CreatedAt:   time.Now(),
		Slot:        slot,
	}
	bookings[bookingID] = booking

	slot.FreeSeats -= req.SeatsCount
	slot.FreeRentalBoards -= req.RentalCount
	slots[req.SlotID] = slot

	if idempKey != "" {
		idempotency[idempKey] = bookingID
	}

	jsonResponse(w, booking, http.StatusCreated)
}

func handleListBookings(w http.ResponseWriter, r *http.Request) {
	clientID := r.Header.Get("X-Client-ID")
	if clientID == "" {
		errorResponse(w, "unauthorized", "Не авторизован", http.StatusUnauthorized)
		return
	}

	items := []Booking{}
	for _, b := range bookings {
		if b.ClientID == clientID {
			if slot, ok := slots[b.SlotID]; ok {
				b.Slot = slot
			}
			items = append(items, b)
		}
	}

	jsonResponse(w, map[string]interface{}{
		"items": items,
		"meta": map[string]int{
			"limit":  20,
			"offset": 0,
			"total":  len(items),
		},
	}, http.StatusOK)
}

func handleCancelBooking(w http.ResponseWriter, r *http.Request) {
	clientID := r.Header.Get("X-Client-ID")
	if clientID == "" {
		errorResponse(w, "unauthorized", "Не авторизован", http.StatusUnauthorized)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/bookings/")
	path = strings.TrimSuffix(path, "/cancel")
	bookingID := path

	booking, ok := bookings[bookingID]
	if !ok {
		errorResponse(w, "not_found", "Бронь не найдена", http.StatusNotFound)
		return
	}
	if booking.ClientID != clientID {
		errorResponse(w, "forbidden", "Это не ваша бронь", http.StatusForbidden)
		return
	}
	if booking.Status != "active" {
		errorResponse(w, "already_cancelled", "Бронь уже отменена", http.StatusConflict)
		return
	}

	slot, ok := slots[booking.SlotID]
	if !ok {
		errorResponse(w, "not_found", "Слот не найден", http.StatusNotFound)
		return
	}
	if time.Now().After(slot.StartAt) {
		errorResponse(w, "slot_started", "Тренировка уже началась", http.StatusUnprocessableEntity)
		return
	}

	hoursUntilStart := slot.StartAt.Sub(time.Now()).Hours()
	now := time.Now()

	if hoursUntilStart >= 2 {
		booking.Status = "cancelled"
		slot.FreeSeats += booking.SeatsCount
		slot.FreeRentalBoards += booking.RentalCount
		slots[booking.SlotID] = slot
	} else {
		booking.Status = "late_cancel"
	}

	booking.CancelledAt = &now
	bookings[bookingID] = booking

	jsonResponse(w, booking, http.StatusOK)
}

// ============================================================
// 6. ЗАПУСК
// ============================================================

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/verify-code", handleVerifyCode)
	mux.HandleFunc("GET /slots", authMiddleware(handleListSlots))
	mux.HandleFunc("POST /bookings", authMiddleware(handleCreateBooking))
	mux.HandleFunc("GET /bookings", authMiddleware(handleListBookings))
	mux.HandleFunc("POST /bookings/{id}/cancel", authMiddleware(handleCancelBooking))

	handler := corsMiddleware(mux)

	log.Println("🚀 Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Idempotency-Key")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}