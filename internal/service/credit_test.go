package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tucredito/backend-api/internal/domain"
	"github.com/tucredito/backend-api/internal/decision"
	"github.com/tucredito/backend-api/internal/event"
	repomocks "github.com/tucredito/backend-api/internal/repository/mocks"
	"go.uber.org/zap"
)

func TestCreditService_CreateSync_InvalidInput(t *testing.T) {
	log, _ := zap.NewDevelopment()
	creditRepo := &repomocks.CreditRepository{}
	clientRepo := &repomocks.ClientRepository{}
	bankRepo := &repomocks.BankRepository{}
	publisher := event.NewMockPublisher()
	engine := decision.NewRuleEngine()
	engine.RegisterRule(decision.PaymentRangeRule{})

	svc := NewCreditService(creditRepo, clientRepo, bankRepo, nil, publisher, engine, log)
	defer svc.Shutdown()

	_, err := svc.CreateSync(context.Background(), domain.CreateCreditInput{
		ClientID:   "c1",
		BankID:     "b1",
		MinPayment: 100,
		MaxPayment: 50,
		TermMonths: 12,
	})
	require.Error(t, err)
	assert.Equal(t, ErrInvalidInput, err)
}

func TestCreditService_CreateSync_ClientNotFound(t *testing.T) {
	log, _ := zap.NewDevelopment()
	creditRepo := &repomocks.CreditRepository{}
	clientRepo := &repomocks.ClientRepository{}
	clientRepo.GetByIDFunc = func(ctx context.Context, id string) (*domain.Client, error) {
		return nil, nil
	}
	bankRepo := &repomocks.BankRepository{}
	publisher := event.NewMockPublisher()
	engine := decision.NewRuleEngine()
	engine.RegisterRule(decision.PaymentRangeRule{})

	svc := NewCreditService(creditRepo, clientRepo, bankRepo, nil, publisher, engine, log)
	defer svc.Shutdown()

	_, err := svc.CreateSync(context.Background(), domain.CreateCreditInput{
		ClientID:   "c1",
		BankID:     "b1",
		MinPayment: 100,
		MaxPayment: 500,
		TermMonths: 12,
		CreditType: domain.CreditTypeAuto,
	})
	require.Error(t, err)
	assert.Equal(t, ErrClientNotFound, err)
}

func TestCreditService_CreateSync_Success(t *testing.T) {
	log, _ := zap.NewDevelopment()
	client := &domain.Client{ID: "c1", FullName: "Test", Email: "a@b.com", Country: "US", BirthDate: time.Now()}
	bank := &domain.Bank{ID: "b1", Name: "Bank", Type: domain.BankTypePrivate}
	credit := &domain.Credit{
		ID: "cr1", ClientID: "c1", BankID: "b1",
		MinPayment: 100, MaxPayment: 500, TermMonths: 12,
		CreditType: domain.CreditTypeAuto, Status: domain.CreditStatusPending,
		CreatedAt: time.Now(),
	}
	creditRepo := &repomocks.CreditRepository{}
	creditRepo.CreateFunc = func(ctx context.Context, input domain.CreateCreditInput) (*domain.Credit, error) {
		return credit, nil
	}
	creditRepo.UpdateStatusFunc = func(ctx context.Context, id string, status domain.CreditStatus) (*domain.Credit, error) {
		c := *credit
		c.Status = status
		return &c, nil
	}
	clientRepo := &repomocks.ClientRepository{}
	clientRepo.GetByIDFunc = func(ctx context.Context, id string) (*domain.Client, error) {
		return client, nil
	}
	bankRepo := &repomocks.BankRepository{}
	bankRepo.GetByIDFunc = func(ctx context.Context, id string) (*domain.Bank, error) {
		return bank, nil
	}
	publisher := event.NewMockPublisher()
	engine := decision.NewRuleEngine()
	engine.RegisterRule(decision.PaymentRangeRule{})

	svc := NewCreditService(creditRepo, clientRepo, bankRepo, nil, publisher, engine, log)
	defer svc.Shutdown()

	out, err := svc.CreateSync(context.Background(), domain.CreateCreditInput{
		ClientID:   "c1",
		BankID:     "b1",
		MinPayment: 100,
		MaxPayment: 500,
		TermMonths: 12,
		CreditType: domain.CreditTypeAuto,
	})
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, "cr1", out.ID)
	events := publisher.Events()
	assert.GreaterOrEqual(t, len(events), 1)
}

func TestCreditService_Reenable(t *testing.T) {
	reenabled := &domain.Credit{ID: "cr1", ClientID: "c1", BankID: "b1", IsActive: true}
	creditRepo := &repomocks.CreditRepository{}
	creditRepo.SetActiveFunc = func(ctx context.Context, id string) (*domain.Credit, error) {
		if id == "cr1" {
			return reenabled, nil
		}
		return nil, nil
	}
	clientRepo := &repomocks.ClientRepository{}
	bankRepo := &repomocks.BankRepository{}
	publisher := event.NewMockPublisher()
	engine := decision.NewRuleEngine()
	log, _ := zap.NewDevelopment()
	svc := NewCreditService(creditRepo, clientRepo, bankRepo, nil, publisher, engine, log)
	defer svc.Shutdown()

	got, err := svc.Reenable(context.Background(), "cr1")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.True(t, got.IsActive)
}

func TestCreditService_Reenable_NotFound(t *testing.T) {
	creditRepo := &repomocks.CreditRepository{}
	creditRepo.SetActiveFunc = func(ctx context.Context, id string) (*domain.Credit, error) {
		return nil, nil
	}
	clientRepo := &repomocks.ClientRepository{}
	bankRepo := &repomocks.BankRepository{}
	publisher := event.NewMockPublisher()
	engine := decision.NewRuleEngine()
	log, _ := zap.NewDevelopment()
	svc := NewCreditService(creditRepo, clientRepo, bankRepo, nil, publisher, engine, log)
	defer svc.Shutdown()

	got, err := svc.Reenable(context.Background(), "none")
	require.NoError(t, err)
	require.Nil(t, got)
}
