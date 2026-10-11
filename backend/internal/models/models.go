package models

type User struct {
	ID           int    `db:"id" json:"id"`
	Username     string `db:"username" json:"username"`
	Role         string `db:"role" json:"role"`
	PasswordHash string `db:"password_hash" json:"-"`
}

type Car struct {
	ID        int    `db:"id" json:"id"`
	Num       string `db:"num" json:"num"`
	Color     string `db:"color" json:"color"`
	Mark      string `db:"mark" json:"mark"`
	IsForeign int    `db:"is_foreign" json:"is_foreign"`
}

type Master struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Service struct {
	ID          int     `db:"id" json:"id"`
	Name        string  `db:"name" json:"name"`
	CostOur     float64 `db:"cost_our" json:"cost_our"`
	CostForeign float64 `db:"cost_foreign" json:"cost_foreign"`
}

type Work struct {
	ID        int    `db:"id" json:"id"`
	DateWork  string `db:"date_work" json:"date_work"`
	MasterID  int    `db:"master_id" json:"master_id"`
	CarID     int    `db:"car_id" json:"car_id"`
	ServiceID int    `db:"service_id" json:"service_id"`
}
