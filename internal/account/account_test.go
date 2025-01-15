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
	a, err := NewAccount("Coffy")
	if err != nil {
		t.Errorf("Error creating new Account: %s", err.Error())
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

func TestAccountConsumption(t *testing.T) {
	a, err := NewAccount("Coffy")
	if err != nil {
		t.Errorf("Error creating new Account: %s", err.Error())
	}
	if err := a.Consume(0.25, "coffee cream"); err != nil {
		t.Errorf("Error consuming Account: %s", err.Error())
	}
	if a.Balance() != -0.25 {
		t.Errorf("Balance should be %f, got %f", -0.25, a.Balance())
	}
}

func TestAccountConsumptionMalicious(t *testing.T) {
	a, err := NewAccount("Coffy")
	if err != nil {
		t.Errorf("Error creating new Account: %s", err.Error())
	}
	// We try to consume a coffee with negative price to cheat our balance
	if err := a.Consume(-0.25, "coffee cream"); err == nil {
		t.Errorf("Expected error, got none")
	}
}

func TestAccountPayment(t *testing.T) {
	a, err := NewAccount("Coffy")
	if err != nil {
		t.Errorf("Error creating new Account: %s", err.Error())
	}
	amount := 5.00
	reason := "debt balance"
	if err := a.Pay(amount, reason); err != nil {
		t.Errorf("Error paying Account: %s", err.Error())
	}
	if a.Balance() != amount {
		t.Errorf("Balance should be %.2f, got %.2f", amount, a.Balance())
	}
}

func TestAccountPaymentMalicious(t *testing.T) {
	a, err := NewAccount("Coffy")
	if err != nil {
		t.Errorf("Error creating new Account: %s", err.Error())
	}
	amount := -5.00 // negative values for payment are not allowed
	reason := "debt balance"
	if err := a.Pay(amount, reason); err == nil {
		t.Errorf("Expected error, got none")
	}
}
