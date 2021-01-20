package models

type User struct {
	ID         int64  `json:"id" db:"id"`
	Surname    string `json:"surname" db:"surname"`
	Name       string `json:"name" db:"name"`
	Patronymic string `json:"patronymic" db:"patronymic"`
	Login      string `json:"login" db:"login"`
	Password   string `json:"password" db:"password"`
	Access     int64  `json:"access" db:"access"`
}
