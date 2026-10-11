package models

type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Role     string `db:"role"`
}

type Car struct {
	ID        int    `db:"id"`
	Num       string `db:"num"`
	Color     string `db:"color"`
	Mark      string `db:"mark"`
	IsForeign int    `db:"is_foreign"`
}

type Master struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Service struct {
	ID          int     `db:"id"`
	Name        string  `db:"name"`
	CostOur     float64 `db:"cost_our"`
	CostForeign float64 `db:"cost_foreign"`
}

type Work struct {
	ID        int    `db:"id"`
	DateWork  string `db:"date_work"`
	MasterID  int    `db:"master_id"`
	CarID     int    `db:"car_id"`
	ServiceID int    `db:"service_id"`
}
