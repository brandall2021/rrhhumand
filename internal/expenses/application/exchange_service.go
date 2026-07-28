package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

type ExchangeRepository interface {
	Create(ctx context.Context, rate *domain.ExchangeRate) error
	GetLatest(ctx context.Context, fromCurrency, toCurrency string) (*domain.ExchangeRate, error)
	GetByDate(ctx context.Context, fromCurrency, toCurrency string, date time.Time) (*domain.ExchangeRate, error)
}

type ExchangeService struct {
	exchangeRepo ExchangeRepository
}

func NewExchangeService(exchangeRepo ExchangeRepository) *ExchangeService {
	return &ExchangeService{exchangeRepo: exchangeRepo}
}

func (s *ExchangeService) CreateRate(ctx context.Context, companyID uuid.UUID, rate *domain.ExchangeRate) (*domain.ExchangeRate, error) {
	const op = "CreateRate"
	now := time.Now()
	rate.ID = uuid.New()
	rate.CompanyID = companyID
	rate.CreatedAt = now
	if err := s.exchangeRepo.Create(ctx, rate); err != nil {
		return nil, svcErr(op, err)
	}
	return rate, nil
}

func (s *ExchangeService) GetLatestRate(ctx context.Context, fromCurrency, toCurrency string) (*decimal.Decimal, error) {
	const op = "GetLatestRate"
	rate, err := s.exchangeRepo.GetLatest(ctx, fromCurrency, toCurrency)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return &rate.Rate, nil
}

func (s *ExchangeService) Convert(ctx context.Context, amount decimal.Decimal, fromCurrency, toCurrency string, date time.Time) (*decimal.Decimal, error) {
	const op = "Convert"
	var rate *domain.ExchangeRate
	var err error

	if date.IsZero() {
		rate, err = s.exchangeRepo.GetLatest(ctx, fromCurrency, toCurrency)
	} else {
		rate, err = s.exchangeRepo.GetByDate(ctx, fromCurrency, toCurrency, date)
	}
	if err != nil {
		return nil, svcErr(op, err)
	}

	converted := amount.Mul(rate.Rate)
	return &converted, nil
}
