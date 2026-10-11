package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TheEmfield/autoservice/backend/internal/database"
	"github.com/TheEmfield/autoservice/backend/internal/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	user, err := database.Authenticate(creds.Username, creds.Password)
	if err != nil {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func GetCars(w http.ResponseWriter, r *http.Request) {
	cars, err := database.GetAllCars()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

func CreateCar(w http.ResponseWriter, r *http.Request) {
	var car models.Car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	if err := database.AddCar(car); err != nil {
		http.Error(w, "Ошибка добавления (возможно, дубликат номера): "+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func DeleteCar(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	var id int
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if err := database.DeleteCar(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// --- MASTERS ---
func GetMasters(w http.ResponseWriter, r *http.Request) {
	masters, err := database.GetAllMasters()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(masters)
}

func CreateMaster(w http.ResponseWriter, r *http.Request) {
	var m models.Master
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	if err := database.AddMaster(m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func DeleteMaster(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	if err := database.DeleteMaster(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// --- SERVICES ---
func GetServices(w http.ResponseWriter, r *http.Request) {
	services, err := database.GetAllServices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func CreateService(w http.ResponseWriter, r *http.Request) {
	var s models.Service
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	if err := database.AddService(s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func DeleteService(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	if err := database.DeleteService(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// --- WORKS ---
func GetWorks(w http.ResponseWriter, r *http.Request) {
	works, err := database.GetAllWorks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(works)
}

func CreateWork(w http.ResponseWriter, r *http.Request) {
	var wk models.Work
	if err := json.NewDecoder(r.Body).Decode(&wk); err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	if err := database.AddWork(wk); err != nil {
		http.Error(w, "Ошибка БД (проверьте триггеры): "+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func DeleteWork(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	if err := database.DeleteWork(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// --- REPORTS ---
func ReportCostByType(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	reports, err := database.GetCostByType(start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=report_cost.txt")
	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintf(w, "ОТЧЕТ: Стоимость обслуживания по типам авто\n")
	fmt.Fprintf(w, "Период: %s - %s\n", start, end)
	fmt.Fprintf(w, "--------------------------------------------------\n")
	for _, rep := range reports {
		fmt.Fprintf(w, "Тип: %-15s | Общая стоимость: %.2f руб.\n", rep.CarType, rep.TotalCost)
	}
}

func ReportTopMasters(w http.ResponseWriter, r *http.Request) {
	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")

	year, _ := strconv.Atoi(yearStr)
	month, _ := strconv.Atoi(monthStr)

	reports, err := database.GetTopMasters(year, month)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=top_masters.txt")
	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintf(w, "ОТЧЕТ: Топ-5 мастеров за %d-%02d\n", year, month)
	fmt.Fprintf(w, "--------------------------------------------------\n")
	for _, rep := range reports {
		fmt.Fprintf(w, "Мастер: %-30s | Работ: %d\n", rep.MasterName, rep.WorksCount)
	}
}
