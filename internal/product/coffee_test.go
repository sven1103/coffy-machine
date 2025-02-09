package product

import (
	"coffy/internal/event"
	"testing"
)

func TestNewCoffee(t *testing.T) {
	b, err := NewCoffee("Black Coffee", 0.25)
	if err != nil {
		t.Errorf("NewCoffee() error = %v", err)
		return
	}
	if b == nil {
		t.Errorf("NewCoffee() returned nil")
		return
	}
	if b.Type != "Black Coffee" {
		t.Errorf("Type() should have returned Black Coffee")
		return
	}
	if b.Price() != 0.25 {
		t.Errorf("Price() should have returned '0.25'")
		return
	}
	if len(b.Events()) != 1 {
		t.Errorf("Events() should have returned 1 event")
		return
	}
	switch eventType := b.Events()[0].(type) {
	case CoffeeCreated:
		return
	default:
		t.Errorf("Events() should have returned a CoffeeCreated, but returned %v", eventType)
	}
}

func TestPriceUpdate(t *testing.T) {
	oldPrice := 0.25
	newPrice := 0.40

	bev, err := NewCoffee("Black Coffee", oldPrice)
	if err != nil {
		t.Errorf("NewCoffee() error = %v", err)
		return
	}
	if err := bev.ChangePrice(newPrice, "Global warming!"); err != nil {
		t.Errorf("ChangePrice() error = %v", err)
		return
	}
	if bev.Price() != newPrice {
		t.Errorf("Price() should have returned '%.2f'", newPrice)
	}
}

func TestPriceUpdatePositiveOnly(t *testing.T) {
	oldPrice := 0.25
	newPrice := -0.2

	bev, err := NewCoffee("Black Coffee", oldPrice)
	if err != nil {
		t.Errorf("NewCoffee() error = %v", err)
		return
	}
	if err := bev.ChangePrice(newPrice, "Global warming!"); err == nil {
		t.Errorf("ChangePrice() should have retured an error, since the new price was negative")
		return
	}
}

func TestCVA(t *testing.T) {
	bev, _ := NewCoffee("Black Coffee", 0.25)
	bev.Clear()

	err := bev.SetCuppingScore(95)
	if err != nil {
		t.Errorf("SetCuppingScore() error = %v", err)
		return
	}

	switch eType := bev.Events()[0].(type) {
	case CvaProvided:
		bev.Clear()
		eType.Value = 89
		list := make([]event.Event, 0)
		list = append(list, eType)
		err := bev.Load(list)
		if err != nil {
			t.Errorf("Load() error = %v", err)
		}
		if bev.cva.Value != 89 {
			t.Errorf("cva.Value should have been 89")
		}
	default:
		t.Errorf("Events() should have returned a 'CvaProvided' but was %T", eType)
		return
	}
}

func TestDetails(t *testing.T) {
	c, _ := NewCoffee("Black Coffee", 0.25)
	c.Clear()

	d := Details{"Kirinyaga County", "Kenianischer Waldkaffee", "", nil}

	err := c.UpdateDetails(d)

	if err != nil {
		t.Errorf("UpdateDetails() error = %v", err)
	}

	if !(c.Details().Origin == "Kirinyaga County") {
		t.Errorf("Details should have been Kirinyaga County")
	}

	switch c.Events()[0].(type) {
	case DetailsUpdated:
		break
	default:
		t.Errorf("Events() should have returned a DetailsUpdated, but returned %v", c.Events()[0])
	}

}
