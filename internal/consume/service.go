package consume

import (
	"coffy/internal/account"
	"coffy/internal/product"
	"errors"
	"fmt"
	"time"
)

const recipient string = "Coffy - Consume Service"

type Service struct {
	accounting *account.Accounting
	product    *product.Service
}

type Receipt struct {
	Recipient string    `json:"recipient"`
	Submitter string    `json:"submitter"`
	Amount    float64   `json:"amount"`
	Purpose   string    `json:"purpose"`
	Date      time.Time `json:"date"`
}

func NewService(accounting *account.Accounting, product *product.Service) *Service {
	return &Service{accounting: accounting, product: product}
}

func (s *Service) Consume(accountID string, productID string, n int) (*Receipt, error) {
	// first fetch the account
	a, err := s.accounting.Find(accountID)
	if err != nil {
		return nil, errors.Join(ErrorAccountNotFound, err)
	}
	p, err := s.product.Find(productID)
	if err != nil {
		return nil, errors.Join(ErrorProductNotFound, err)
	}

	if err = s.accounting.Consume(accountID, p.Price(), p.Type, n); err != nil {
		return nil, errors.Join(errors.New("failed to consume product"), err)
	}

	return &Receipt{
		Recipient: recipient,
		Submitter: a.Owner(),
		Amount:    float64(n) * p.Price(),
		Purpose:   fmt.Sprintf("consumption of '%s'", p.Type),
		Date:      time.Now()}, nil
}

var ErrorProductNotFound = errors.New("product not found")
var ErrorAccountNotFound = errors.New("account not found")
