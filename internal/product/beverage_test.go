package product

import "testing"

func TestNewBeverage(t *testing.T) {
	b, err := NewBeverage("Black Coffee", 0.25)
	if err != nil {
		t.Errorf("NewBeverage() error = %v", err)
		return
	}
	if b == nil {
		t.Errorf("NewBeverage() returned nil")
		return
	}
	if b.BeverageType != "Black Coffee" {
		t.Errorf("BeverageType() should have returned Black Coffee")
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
	case BeverageCreated:
		return
	default:
		t.Errorf("Events() should have returned a BeverageCreated, but returned %v", eventType)
	}
}

func TestPriceUpdate(t *testing.T) {
	oldPrice := 0.25
	newPrice := 0.40

	bev, err := NewBeverage("Black Coffee", oldPrice)
	if err != nil {
		t.Errorf("NewBeverage() error = %v", err)
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
