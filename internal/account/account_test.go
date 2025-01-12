package account

import "testing"

func TestCoffyConsumed(t *testing.T) {
	costs := 0.25
	coffee := "coffee cream"
	event := newCoffyConsumed("123", coffee, costs)

	if event.AggregateID() != "123" {
		t.Errorf("Event returned ID should be '123', got '%s'", event.AggregateID())
	}
	if event.CoffyType != coffee {
		t.Errorf("CoffyType should be '%s', got '%s'", coffee, event.CoffyType)
	}
	if event.Costs != costs {
		t.Errorf("Costs should be %f, got %f", costs, event.Costs)
	}
}

func TestIncomingPayment(t *testing.T) {
	amount := 5.00
	reason := "debt balance"
	event := newIncomingPayment("123", amount, reason)

	if event.AggregateID() != "123" {
		t.Errorf("Event returned ID should be '123', got '%s'", event.AggregateID())
		return
	}
	if event.Amount != amount {
		t.Errorf("Incoming payment amount should be %f, got %f", amount, event.Amount)
		return
	}
	if event.Reason != reason {
		t.Errorf("Incoming payment reason should be '%s', got '%s'", reason, event.Reason)
		return
	}
}

func TestAccountCreation(t *testing.T) {
	a, err := newAccount("Coffy")
	if err != nil {
		t.Errorf("Error creating new account: %s", err.Error())
		return
	}
	if a.Owner() != "Coffy" {
		t.Errorf("Owner should be 'Coffy', got '%s'", a.Owner())
	}
	events := a.Events()
	if len(events) != 1 {
		t.Errorf("Events should be 1, got %d", len(events))
		return
	}
	if events[0].Type() != "AccountCreated" {
		t.Errorf("Event Type should be 'AccountCreated', got '%s'", events[0].Type())
	}
}
