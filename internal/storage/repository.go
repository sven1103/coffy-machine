package storage

import (
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type EventRepository interface {
	SaveAll([]EventEntry) error
	LoadAll(aggregateID string) ([]EventEntry, error)
	FetchByEventType(event string) ([]EventEntry, error)
}

type EventEntry struct {
	ID          int
	AggregateID string
	EventType   string
	Date        time.Time
	EventData   []byte
}

type eventRepositoryImpl struct {
	db *gorm.DB
}

func (r *eventRepositoryImpl) SaveAll(events []EventEntry) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = errors.New("saving events failed")
		}
	}()
	r.db.Create(events)
	return nil
}

func (r *eventRepositoryImpl) LoadAll(aggregateID string) ([]EventEntry, error) {
	users := make([]EventEntry, 0)
	result := r.db.Where("aggregate_id = ?", aggregateID).Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("error loading events: %w", result.Error)
	}
	return users, nil
}

func (r *eventRepositoryImpl) FetchByEventType(t string) ([]EventEntry, error) {
	events := make([]EventEntry, 0)
	result := r.db.Where("event_type LIKE ?", t).Find(&events)
	if result.Error != nil {
		return nil, fmt.Errorf("error fetching accounts: %w", result.Error)
	}
	return events, nil
}

func CreateEventRepository(storage string) (EventRepository, error) {
	db, err := gorm.Open(sqlite.Open(storage), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&EventEntry{}); err != nil {
		return nil, err
	}
	return &eventRepositoryImpl{db}, nil
}
