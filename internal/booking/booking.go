package booking

import (
	"time"
)

// Booking represents a user's reservation for a specific seat at an event.
// Uses a composite unique index on (event_id, seat_number, active) to prevent
// double booking of the same seat while allowing historical records after cancel.
type Booking struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UserID     uint `gorm:"not null;index" json:"user_id"`
	EventID    uint `gorm:"not null;index;uniqueIndex:uniq_event_seat_active" json:"event_id"`
	SeatNumber int  `gorm:"not null;uniqueIndex:uniq_event_seat_active" json:"seat_number"`
	Active     bool `gorm:"not null;default:true;uniqueIndex:uniq_event_seat_active" json:"active"`

	// Optional payment simulation fields
	Paid   bool   `gorm:"not null;default:false" json:"paid"`
	Status string `gorm:"size:32;not null;default:reserved" json:"status"`
}
