package database

import (
	"fmt"

	"github.com/TheEmfield/autoservice/backend/internal/config"
	"github.com/TheEmfield/autoservice/backend/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var DB *sqlx.DB

func Connect(cfg *config.Config) error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.DB, cfg.Database.SSLMode)

	var err error
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		return err
	}
	return DB.Ping()
}

func Authenticate(username, password string) (*models.User, error) {
	var user models.User
	err := DB.Get(&user, "SELECT id, username, role, password_hash FROM users WHERE username=$1", username)
	if err != nil {
		return nil, fmt.Errorf("пользователь не найден")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("неверный пароль")
	}

	return &user, nil
}

func GetAllCars() ([]models.Car, error) {
	var cars []models.Car
	err := DB.Select(&cars, "SELECT id, num, color, mark, is_foreign FROM cars ORDER BY id")
	return cars, err
}
func AddCar(c models.Car) error {
	_, err := DB.NamedExec(`INSERT INTO cars (num, color, mark, is_foreign) VALUES (:num, :color, :mark, :is_foreign)`, c)
	return err
}
func DeleteCar(id int) error {
	_, err := DB.Exec("DELETE FROM cars WHERE id=$1", id)
	return err
}

func GetAllMasters() ([]models.Master, error) {
	var masters []models.Master
	err := DB.Select(&masters, "SELECT id, name FROM masters ORDER BY id")
	return masters, err
}
func AddMaster(m models.Master) error {
	_, err := DB.NamedExec(`INSERT INTO masters (name) VALUES (:name)`, m)
	return err
}
func DeleteMaster(id int) error {
	_, err := DB.Exec("DELETE FROM masters WHERE id=$1", id)
	return err
}

func GetAllServices() ([]models.Service, error) {
	var services []models.Service
	err := DB.Select(&services, "SELECT id, name, cost_our, cost_foreign FROM services ORDER BY id")
	return services, err
}
func AddService(s models.Service) error {
	_, err := DB.NamedExec(`INSERT INTO services (name, cost_our, cost_foreign) VALUES (:name, :cost_our, :cost_foreign)`, s)
	return err
}
func DeleteService(id int) error {
	_, err := DB.Exec("DELETE FROM services WHERE id=$1", id)
	return err
}

func GetAllWorks() ([]models.Work, error) {
	var works []models.Work
	query := `SELECT w.id, w.date_work, w.master_id, w.car_id, w.service_id FROM works w ORDER BY w.date_work DESC`
	err := DB.Select(&works, query)
	return works, err
}
func AddWork(w models.Work) error {
	_, err := DB.NamedExec(`INSERT INTO works (date_work, master_id, car_id, service_id) VALUES (:date_work, :master_id, :car_id, :service_id)`, w)
	return err
}
func DeleteWork(id int) error {
	_, err := DB.Exec("DELETE FROM works WHERE id=$1", id)
	return err
}

type CostReport struct {
	CarType    string  `db:"car_type"`
	TotalCost  float64 `db:"total_cost"`
	WorksCount int64   `db:"works_count"`
}

func GetCostByType(startDate, endDate string) ([]CostReport, error) {
	var reports []CostReport
	query := `SELECT * FROM get_cost_by_car_type($1::date, $2::date)`
	err := DB.Select(&reports, query, startDate, endDate)
	return reports, err
}

type TopMasterReport struct {
	MasterName string `db:"master_name"`
	WorksCount int64  `db:"works_count"`
	CarsCount  int64  `db:"cars_count"`
}

func GetTopMasters(year, month int) ([]TopMasterReport, error) {
	var reports []TopMasterReport
	query := `SELECT * FROM get_top_5_masters_by_month($1, $2)`
	err := DB.Select(&reports, query, year, month)
	return reports, err
}

func UpdateCar(id int, c models.Car) error {
	_, err := DB.NamedExec(`UPDATE cars SET num=:num, color=:color, mark=:mark, is_foreign=:is_foreign WHERE id=:id`,
		map[string]interface{}{"id": id, "num": c.Num, "color": c.Color, "mark": c.Mark, "is_foreign": c.IsForeign})
	return err
}

func UpdateMaster(id int, m models.Master) error {
	_, err := DB.NamedExec(`UPDATE masters SET name=:name WHERE id=:id`,
		map[string]interface{}{"id": id, "name": m.Name})
	return err
}

func UpdateService(id int, s models.Service) error {
	_, err := DB.NamedExec(`UPDATE services SET name=:name, cost_our=:cost_our, cost_foreign=:cost_foreign WHERE id=:id`,
		map[string]interface{}{"id": id, "name": s.Name, "cost_our": s.CostOur, "cost_foreign": s.CostForeign})
	return err
}

func UpdateWork(id int, w models.Work) error {
	_, err := DB.NamedExec(`UPDATE works SET date_work=:date_work, master_id=:master_id, car_id=:car_id, service_id=:service_id WHERE id=:id`,
		map[string]interface{}{"id": id, "date_work": w.DateWork, "master_id": w.MasterID, "car_id": w.CarID, "service_id": w.ServiceID})
	return err
}
