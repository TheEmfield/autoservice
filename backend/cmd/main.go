package main

import (
	"log"
	"net/http"

	"github.com/TheEmfield/autoservice/backend/internal/config"
	"github.com/TheEmfield/autoservice/backend/internal/database"
	"github.com/TheEmfield/autoservice/backend/internal/handlers"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal("Ошибка загрузки конфига: ", err)
	}

	if err := database.Connect(cfg); err != nil {
		log.Fatal("Ошибка подключения к БД: ", err)
	}
	defer database.DB.Close()

	http.HandleFunc("/api/login", handlers.Login)

	http.HandleFunc("/api/cars", methodRouter(handlers.GetCars, handlers.CreateCar, handlers.DeleteCar))
	http.HandleFunc("/api/masters", methodRouter(handlers.GetMasters, handlers.CreateMaster, handlers.DeleteMaster))
	http.HandleFunc("/api/services", methodRouter(handlers.GetServices, handlers.CreateService, handlers.DeleteService))
	http.HandleFunc("/api/works", methodRouter(handlers.GetWorks, handlers.CreateWork, handlers.DeleteWork))

	http.HandleFunc("/api/reports/cost", handlers.ReportCostByType)
	http.HandleFunc("/api/reports/top-masters", handlers.ReportTopMasters)

	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(":"+cfg.App.Port, nil))
}

func methodRouter(get, post, delete http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			get(w, r)
		case http.MethodPost:
			post(w, r)
		case http.MethodDelete:
			delete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
