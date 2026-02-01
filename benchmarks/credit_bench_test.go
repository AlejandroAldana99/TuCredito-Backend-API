package benchmarks

import (
	"context"
	"testing"
	"time"

	"github.com/tucredito/backend-api/internal/decision"
	"github.com/tucredito/backend-api/internal/domain"
	"github.com/tucredito/backend-api/internal/event"
	"github.com/tucredito/backend-api/internal/repository/mocks"
	"github.com/tucredito/backend-api/internal/service"
	"go.uber.org/zap"
)

func setupCreditServiceBench(_ *testing.B) (service.CreditService, domain.CreateCreditInput) {
	log, _ := zap.NewDevelopment()
	client := &domain.Client{ID: "c1", FullName: "Test", Email: "a@b.com", Country: "US", BirthDate: time.Now()}
	bank := &domain.Bank{ID: "b1", Name: "Bank", Type: domain.BankTypePrivate}
	credit := &domain.Credit{
		ID: "cr1", ClientID: "c1", BankID: "b1",
		MinPayment: 100, MaxPayment: 500, TermMonths: 12,
		CreditType: domain.CreditTypeAuto, Status: domain.CreditStatusPending,
		CreatedAt: time.Now(),
	}
	creditRepo := &mocks.CreditRepository{}
	creditRepo.CreateFunc = func(ctx context.Context, input domain.CreateCreditInput) (*domain.Credit, error) {
		return credit, nil
	}
	creditRepo.GetByIDFunc = func(ctx context.Context, id string) (*domain.Credit, error) {
		return credit, nil
	}
	creditRepo.UpdateFunc = func(ctx context.Context, id string, input domain.UpdateCreditInput) (*domain.Credit, error) {
		c := *credit
		c.MinPayment = input.MinPayment
		c.MaxPayment = input.MaxPayment
		c.TermMonths = input.TermMonths
		c.Status = input.Status
		return &c, nil
	}
	creditRepo.UpdateStatusFunc = func(ctx context.Context, id string, status domain.CreditStatus) (*domain.Credit, error) {
		c := *credit
		c.Status = status
		return &c, nil
	}
	creditRepo.SetInactiveFunc = func(ctx context.Context, id string) (*domain.Credit, error) {
		c := *credit
		c.IsActive = false
		return &c, nil
	}
	creditRepo.ListFunc = func(ctx context.Context, limit, offset int) ([]*domain.Credit, error) {
		return []*domain.Credit{credit}, nil
	}
	creditRepo.ListByClientIDFunc = func(ctx context.Context, clientID string, limit, offset int) ([]*domain.Credit, error) {
		return []*domain.Credit{credit}, nil
	}
	clientRepo := &mocks.ClientRepository{}
	clientRepo.GetByIDFunc = func(ctx context.Context, id string) (*domain.Client, error) { return client, nil }
	bankRepo := &mocks.BankRepository{}
	bankRepo.GetByIDFunc = func(ctx context.Context, id string) (*domain.Bank, error) { return bank, nil }
	publisher := event.NewMockPublisher()
	engine := decision.NewRuleEngine()
	engine.RegisterRule(decision.PaymentRangeRule{})

	svc := service.NewCreditService(creditRepo, clientRepo, bankRepo, nil, publisher, engine, log)
	return svc, domain.CreateCreditInput{
		ClientID: "c1", BankID: "b1",
		MinPayment: 100, MaxPayment: 500, TermMonths: 12,
		CreditType: domain.CreditTypeAuto,
	}
}

func BenchmarkCreditService_CreateSync(b *testing.B) {
	svc, input := setupCreditServiceBench(b)
	defer svc.Shutdown()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.CreateSync(ctx, input)
	}
}

func BenchmarkCreditService_Create(b *testing.B) {
	svc, input := setupCreditServiceBench(b)
	defer svc.Shutdown()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.Create(ctx, input)
	}
}

func BenchmarkCreditService_GetByID(b *testing.B) {
	svc, _ := setupCreditServiceBench(b)
	defer svc.Shutdown()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.GetByID(ctx, "cr1")
	}
}

func BenchmarkCreditService_Update(b *testing.B) {
	svc, _ := setupCreditServiceBench(b)
	defer svc.Shutdown()
	ctx := context.Background()
	input := domain.UpdateCreditInput{
		MinPayment: 150, MaxPayment: 600, TermMonths: 24,
		Status: domain.CreditStatusApproved,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.Update(ctx, "cr1", input)
	}
}

func BenchmarkCreditService_UpdateStatus(b *testing.B) {
	svc, _ := setupCreditServiceBench(b)
	defer svc.Shutdown()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.UpdateStatus(ctx, "cr1", domain.CreditStatusApproved)
	}
}

func BenchmarkCreditService_Delete(b *testing.B) {
	svc, _ := setupCreditServiceBench(b)
	defer svc.Shutdown()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.Delete(ctx, "cr1")
	}
}

func BenchmarkCreditService_List(b *testing.B) {
	svc, _ := setupCreditServiceBench(b)
	defer svc.Shutdown()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.List(ctx, 10, 0)
	}
}

func BenchmarkCreditService_ListByClientID(b *testing.B) {
	svc, _ := setupCreditServiceBench(b)
	defer svc.Shutdown()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.ListByClientID(ctx, "c1", 10, 0)
	}
}
