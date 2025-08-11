package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewRepository(db *gorm.DB, redis *redis.Client) *Repository {
	return &Repository{db: db, redis: redis}
}

func (r *Repository) Create(event *Event) error {
	// Invalidate cache when creating new event
	if r.redis != nil {
		r.invalidateCache()
	}
	return r.db.Create(event).Error
}

func (r *Repository) Update(event *Event) error {
	// Invalidate cache when updating event
	if r.redis != nil {
		r.invalidateCache()
	}
	return r.db.Save(event).Error
}

func (r *Repository) Delete(id uint) error {
	// Invalidate cache when deleting event
	if r.redis != nil {
		r.invalidateCache()
	}
	return r.db.Delete(&Event{}, id).Error
}

func (r *Repository) GetByID(id uint) (*Event, error) {
	// Try cache first
	if r.redis != nil {
		if cached, err := r.getFromCache(id); err == nil {
			return cached, nil
		}
	}

	var ev Event
	err := r.db.First(&ev, id).Error
	if err != nil {
		return nil, err
	}

	// Cache the result
	if r.redis != nil {
		r.setInCache(id, &ev)
	}

	return &ev, nil
}

func (r *Repository) List() ([]Event, error) {
	// Try cache first
	if r.redis != nil {
		if cached, err := r.getListFromCache(); err == nil {
			return cached, nil
		}
	}

	var list []Event
	err := r.db.Order("starts_at asc").Find(&list).Error
	if err != nil {
		return nil, err
	}

	// Cache the result
	if r.redis != nil {
		r.setListInCache(list)
	}

	return list, nil
}

type ListFilters struct {
	Location string
	From     *time.Time
	To       *time.Time
}

func (r *Repository) ListFiltered(f ListFilters) ([]Event, error) {
	// For filtered results, we don't cache as they're dynamic
	q := r.db.Model(&Event{})
	if f.Location != "" {
		q = q.Where("location ILIKE ?", "%"+f.Location+"%")
	}
	if f.From != nil {
		q = q.Where("starts_at >= ?", *f.From)
	}
	if f.To != nil {
		q = q.Where("starts_at <= ?", *f.To)
	}
	var list []Event
	err := q.Order("starts_at asc").Find(&list).Error
	return list, err
}

// Cache methods
func (r *Repository) getFromCache(id uint) (*Event, error) {
	key := fmt.Sprintf("event:%d", id)
	data, err := r.redis.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, err
	}

	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *Repository) setInCache(id uint, event *Event) {
	key := fmt.Sprintf("event:%d", id)
	data, _ := json.Marshal(event)
	r.redis.Set(context.Background(), key, data, 30*time.Minute)
}

func (r *Repository) getListFromCache() ([]Event, error) {
	data, err := r.redis.Get(context.Background(), "events:list").Bytes()
	if err != nil {
		return nil, err
	}

	var list []Event
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) setListInCache(list []Event) {
	data, _ := json.Marshal(list)
	r.redis.Set(context.Background(), "events:list", data, 15*time.Minute)
}

func (r *Repository) invalidateCache() {
	// Delete all event-related cache keys
	ctx := context.Background()
	r.redis.Del(ctx, "events:list")
	// Note: In production, you might want to use pattern-based deletion
	// r.redis.Eval(ctx, "local keys = redis.call('keys', 'event:*') for i=1,#keys do redis.call('del', keys[i]) end", []string{})
}
