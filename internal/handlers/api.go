package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"calendar-app/internal/database"
	"calendar-app/internal/models"

	"github.com/go-chi/chi/v5"
)

type APIHandler struct {
	db *database.DB
}

func NewAPIHandler(db *database.DB) *APIHandler {
	return &APIHandler{db: db}
}

func (h *APIHandler) GetHolidays(w http.ResponseWriter, r *http.Request) {
	holidays, err := h.db.GetHolidays()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}
	if holidays == nil {
		holidays = []models.Holiday{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(holidays)
}

func (h *APIHandler) AddHoliday(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var holiday models.Holiday
	if err := json.NewDecoder(r.Body).Decode(&holiday); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid request"})
		return
	}

	id, err := h.db.AddHoliday(holiday)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	holiday.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(holiday)
}

func (h *APIHandler) DeleteHoliday(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid id"})
		return
	}

	if err := h.db.DeleteHoliday(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *APIHandler) GetMenuItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.db.GetMenuItems()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *APIHandler) AddMenuItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var item models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid request"})
		return
	}

	if item.MealType != "breakfast" && item.MealType != "lunch" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "meal_type must be 'breakfast' or 'lunch'"})
		return
	}

	id, err := h.db.AddMenuItem(item.MealType, item.ItemName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	item.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *APIHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid id"})
		return
	}

	if err := h.db.DeleteMenuItem(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

type SaveCalendarRequest struct {
	SchoolName string          `json:"school_name"`
	Menus      json.RawMessage `json:"menus"`
}

func (h *APIHandler) GetCalendarData(w http.ResponseWriter, r *http.Request) {
	year, _ := strconv.Atoi(chi.URLParam(r, "year"))
	month, _ := strconv.Atoi(chi.URLParam(r, "month"))

	cd, err := h.db.GetCalendarData(year, month)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	if cd == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"status": "not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cd)
}

func (h *APIHandler) SaveCalendarData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	year, _ := strconv.Atoi(chi.URLParam(r, "year"))
	month, _ := strconv.Atoi(chi.URLParam(r, "month"))

	var body json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid request"})
		return
	}

	var req struct {
		SchoolName string `json:"school_name"`
	}
	json.Unmarshal(body, &req)

	menusJSON, _ := json.Marshal(body)

	cd := models.CalendarData{
		Year:       year,
		Month:      month,
		SchoolName: req.SchoolName,
		DataJSON:   string(menusJSON),
	}

	if err := h.db.SaveCalendarData(cd); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *APIHandler) ExportData(w http.ResponseWriter, r *http.Request) {
	data, err := h.db.ExportAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=calendar-export.json")
	json.NewEncoder(w).Encode(data)
}

func (h *APIHandler) ImportData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data models.ExportData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid JSON"})
		return
	}

	if err := h.db.ImportAll(&data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *APIHandler) ImportHolidays(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Holidays *[]models.Holiday `json:"holidays"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid holidays JSON"})
		return
	}
	if body.Holidays == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "holidays must be an array"})
		return
	}
	for _, holiday := range *body.Holidays {
		maxDate := time.Date(2000, time.Month(holiday.Month+2), 0, 0, 0, 0, 0, time.UTC).Day()
		if holiday.Month < 0 || holiday.Month > 11 || holiday.Date < 1 || holiday.Date > maxDate || holiday.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "each holiday needs a valid month, date, and name"})
			return
		}
	}
	if err := h.db.ImportHolidays(*body.Holidays); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"imported": len(*body.Holidays)})
}

func (h *APIHandler) ImportMenuItems(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Items *[]models.MenuItem `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid menu items JSON"})
		return
	}
	if body.Items == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "items must be an array"})
		return
	}
	for _, item := range *body.Items {
		if (item.MealType != "breakfast" && item.MealType != "lunch") || item.ItemName == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "each menu item needs meal_type breakfast/lunch and item_name"})
			return
		}
	}
	if err := h.db.ImportMenuItems(*body.Items); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"imported": len(*body.Items)})
}

func (h *APIHandler) ImportCalendarData(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CalendarData *[]models.CalendarData `json:"calendar_data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid calendar data JSON"})
		return
	}
	if body.CalendarData == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "calendar_data must be an array"})
		return
	}
	for _, data := range *body.CalendarData {
		if data.Year < 1 || data.Month < 0 || data.Month > 11 || data.DataJSON == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "each calendar record needs year, month 0-11, and data_json"})
			return
		}
		var raw interface{}
		if err := json.Unmarshal([]byte(data.DataJSON), &raw); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "calendar data_json must contain valid JSON"})
			return
		}
		object, ok := raw.(map[string]interface{})
		if !ok || object["menus"] == nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "calendar data_json must contain a menus array"})
			return
		}
		if _, ok := object["menus"].([]interface{}); !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "calendar data_json menus must be an array"})
			return
		}
	}
	if err := h.db.ImportCalendarData(*body.CalendarData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"imported": len(*body.CalendarData)})
}
