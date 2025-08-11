package booking

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateIfSeatFree(ctx context.Context, b *Booking) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(b).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *Repository) ListByUser(userID uint) ([]Booking, error) {
	var list []Booking
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&list).Error
	return list, err
}

func (r *Repository) Cancel(userID, bookingID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var b Booking
		if err := tx.First(&b, "id = ? AND user_id = ? AND active = ?", bookingID, userID, true).Error; err != nil {
			return err
		}
		b.Active = false
		b.Status = "canceled"
		return tx.Save(&b).Error
	})
}

// ValidateSeat checks if seat number is within event capacity
func (r *Repository) ValidateSeat(eventID uint, seatNumber int) error {
	var event struct {
		Capacity int
	}
	if err := r.db.Table("events").Select("capacity").Where("id = ?", eventID).First(&event).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("event not found")
		}
		return err
	}
	if seatNumber < 1 || seatNumber > event.Capacity {
		return errors.New("seat number out of range")
	}
	return nil
}
