package models

type Permission struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}
