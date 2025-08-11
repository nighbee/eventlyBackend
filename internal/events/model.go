package events

import "time"

type Event struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Location    string    `gorm:"size:255;not null" json:"location"`
	StartsAt    time.Time `json:"starts_at"`
	Price       int64     `json:"price"`
	Capacity    int       `json:"capacity"`
}
