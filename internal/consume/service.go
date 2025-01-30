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
	recipient string
	submitter string
	amount    float64
	purpose   string
	date      time.Time
}

func NewService(accounting *account.Accounting, product *product.Service) *Service {
	return &Service{accounting: accounting, product: product}
}

func (s *Service) Consume(accountID string, productID string, n int) (*Receipt, error) {
	// first fetch the account
	a, err := s.accounting.Find(accountID)
	if err != nil {
		return nil, errors.Join(errors.New("failed to lookup account"), err)
	}
	p, err := s.product.Find(productID)
	if err != nil {
		return nil, errors.Join(errors.New("failed to lookup product"), err)
	}

	err = a.ConsumeN(p.Price(), p.BeverageType, n)
	if err != nil {
		return nil, errors.Join(errors.New("failed to consume product"), err)
	}

	err = a.Consume(p.Price(), p.BeverageType)
	if err != nil {
		return nil, errors.Join(errors.New("failed to consume product"), err)
	}

	return &Receipt{
		recipient: recipient,
		submitter: a.Owner(),
		amount:    float64(n) * p.Price(),
		purpose:   fmt.Sprintf("consumption of '%s'", p.BeverageType),
		date:      time.Now()}, nil
}
