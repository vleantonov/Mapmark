package domain

type Mark struct {
	Id          int      `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Description string   `json:"description" db:"description"`
	Image       *string  `json:"image" db:"image"`
	Lat         *float64 `json:"lat" db:"lat"`
	Lng         *float64 `json:"lng" db:"lng"`

	CreatedAT string `json:"created_at"`
	UpdatedAT string `json:"updated_at"`
}
