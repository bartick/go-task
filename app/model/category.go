package model

type Category struct {
	ID   uint64 `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
