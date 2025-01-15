package main

import (
	"coffy/internal/account"
	"coffy/internal/storage"
	"encoding/json"
	"fmt"
)

func main() {
	repo, err := storage.CreateEventRepository("test.db")
	if err != nil {
		panic(err)
	}

	service := account.NewAccounting(&repo)

	acc, err := account.NewAccount("Sven")
	if err != nil {
		panic(err)
	}

	err = acc.Consume(0.35, "coffee cream")
	if err != nil {
		panic(err)
	}
	err = acc.Consume(0.50, "latte macciato")
	if err != nil {
		panic(err)
	}

	eventEntries := make([]storage.EventEntry, 0)
	for _, e := range acc.Events() {
		serial, err := json.Marshal(e)
		if err != nil {
			panic(err)
		}
		eventEntries = append(eventEntries, storage.EventEntry{AggregateID: e.AggregateID(), EventType: e.Type(), EventData: serial})
	}

	repo.SaveAll(eventEntries)

	result, err := service.Find(acc.ID())
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
	fmt.Printf("Account balance of owner '%s' is €%.2f\n", acc.Owner(), acc.Balance())

	for _, e := range acc.Events() {
		fmt.Println(e)
	}
}
