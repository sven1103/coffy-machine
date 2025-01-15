package storage

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type EventRepository interface {
	SaveAll([]EventEntry)
	LoadAll(aggregateID string) ([]EventEntry, error)
}

type EventEntry struct {
	ID          int
	AggregateID string
	EventType   string
	EventData   []byte
}

type eventRepositoryImpl struct {
	db *gorm.DB
}

func (r *eventRepositoryImpl) SaveAll(events []EventEntry) {
	r.db.Create(events)
}

func (r *eventRepositoryImpl) LoadAll(aggregateID string) ([]EventEntry, error) {
	users := make([]EventEntry, 0)
	result := r.db.Where("aggregate_id = ?", aggregateID).Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("error loading events: %w", result.Error)
	}
	return users, nil
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
