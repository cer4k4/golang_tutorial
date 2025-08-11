package domain

type Category struct {
	Id          uint   `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	FatherId    uint   `json:"father_id" db:"father_id"`
}
