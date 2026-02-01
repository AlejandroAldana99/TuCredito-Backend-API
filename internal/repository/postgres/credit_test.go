package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tucredito/backend-api/internal/domain"
	"github.com/tucredito/backend-api/internal/repository/postgres"
)

func TestCreditRepository_Create(t *testing.T) {
	pool := testDBPool(t)
	defer pool.Close()
	ctx := context.Background()
	clientRepo := postgres.NewClientRepository(pool)
	bankRepo := postgres.NewBankRepository(pool)
	creditRepo := postgres.NewCreditRepository(pool)

	client, err := clientRepo.Create(ctx, domain.CreateClientInput{
		FullName: "Credit Test Client", Email: uniqueClientEmail(t), BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), Country: "US",
	})
	require.NoError(t, err)
	require.NotNil(t, client)
	defer deleteClient(t, pool, client.ID)
	bank, err := bankRepo.Create(ctx, domain.CreateBankInput{Name: "Credit Test Bank", Type: domain.BankTypePrivate})
	require.NoError(t, err)
	require.NotNil(t, bank)
	defer deleteBank(t, pool, bank.ID)

	input := domain.CreateCreditInput{
		ClientID:   client.ID,
		BankID:     bank.ID,
		MinPayment: 100,
		MaxPayment: 500,
		TermMonths: 12,
		CreditType: domain.CreditTypeAuto,
	}
	credit, err := creditRepo.Create(ctx, input)
	require.NoError(t, err)
	require.NotNil(t, credit)
	defer deleteCredit(t, pool, credit.ID)
	assert.NotEmpty(t, credit.ID)
	assert.Equal(t, client.ID, credit.ClientID)
	assert.Equal(t, bank.ID, credit.BankID)
	assert.Equal(t, 100.0, credit.MinPayment)
	assert.Equal(t, 500.0, credit.MaxPayment)
	assert.Equal(t, 12, credit.TermMonths)
	assert.Equal(t, domain.CreditTypeAuto, credit.CreditType)
	assert.Equal(t, domain.CreditStatusPending, credit.Status)
	assert.True(t, credit.IsActive)
}

func TestCreditRepository_GetByID(t *testing.T) {
	pool := testDBPool(t)
	defer pool.Close()
	ctx := context.Background()
	clientRepo := postgres.NewClientRepository(pool)
	bankRepo := postgres.NewBankRepository(pool)
	creditRepo := postgres.NewCreditRepository(pool)

	client, err := clientRepo.Create(ctx, domain.CreateClientInput{FullName: "C", Email: uniqueClientEmail(t), BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), Country: "US"})
	require.NoError(t, err)
	require.NotNil(t, client)
	defer deleteClient(t, pool, client.ID)
	bank, err := bankRepo.Create(ctx, domain.CreateBankInput{Name: "B", Type: domain.BankTypePrivate})
	require.NoError(t, err)
	require.NotNil(t, bank)
	defer deleteBank(t, pool, bank.ID)
	created, err := creditRepo.Create(ctx, domain.CreateCreditInput{
		ClientID: client.ID, BankID: bank.ID, MinPayment: 200, MaxPayment: 600, TermMonths: 24, CreditType: domain.CreditTypeMortgage,
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	defer deleteCredit(t, pool, created.ID)

	got, err := creditRepo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, domain.CreditTypeMortgage, got.CreditType)
}

func TestCreditRepository_Update(t *testing.T) {
	pool := testDBPool(t)
	defer pool.Close()
	ctx := context.Background()
	clientRepo := postgres.NewClientRepository(pool)
	bankRepo := postgres.NewBankRepository(pool)
	creditRepo := postgres.NewCreditRepository(pool)

	client, err := clientRepo.Create(ctx, domain.CreateClientInput{FullName: "C", Email: uniqueClientEmail(t), BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), Country: "US"})
	require.NoError(t, err)
	require.NotNil(t, client)
	defer deleteClient(t, pool, client.ID)
	bank, err := bankRepo.Create(ctx, domain.CreateBankInput{Name: "B", Type: domain.BankTypePrivate})
	require.NoError(t, err)
	require.NotNil(t, bank)
	defer deleteBank(t, pool, bank.ID)
	created, err := creditRepo.Create(ctx, domain.CreateCreditInput{
		ClientID: client.ID, BankID: bank.ID, MinPayment: 100, MaxPayment: 500, TermMonths: 12, CreditType: domain.CreditTypeAuto,
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	defer deleteCredit(t, pool, created.ID)

	updated, err := creditRepo.Update(ctx, created.ID, domain.UpdateCreditInput{
		MinPayment: 150, MaxPayment: 550, TermMonths: 18, Status: domain.CreditStatusApproved,
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, 150.0, updated.MinPayment)
	assert.Equal(t, 550.0, updated.MaxPayment)
	assert.Equal(t, 18, updated.TermMonths)
	assert.Equal(t, domain.CreditStatusApproved, updated.Status)
}

func TestCreditRepository_SetInactive(t *testing.T) {
	pool := testDBPool(t)
	defer pool.Close()
	ctx := context.Background()
	clientRepo := postgres.NewClientRepository(pool)
	bankRepo := postgres.NewBankRepository(pool)
	creditRepo := postgres.NewCreditRepository(pool)

	client, err := clientRepo.Create(ctx, domain.CreateClientInput{FullName: "C", Email: uniqueClientEmail(t), BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), Country: "US"})
	require.NoError(t, err)
	require.NotNil(t, client)
	defer deleteClient(t, pool, client.ID)
	bank, err := bankRepo.Create(ctx, domain.CreateBankInput{Name: "B", Type: domain.BankTypePrivate})
	require.NoError(t, err)
	require.NotNil(t, bank)
	defer deleteBank(t, pool, bank.ID)
	created, err := creditRepo.Create(ctx, domain.CreateCreditInput{
		ClientID: client.ID, BankID: bank.ID, MinPayment: 100, MaxPayment: 500, TermMonths: 12, CreditType: domain.CreditTypeAuto,
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	defer deleteCredit(t, pool, created.ID)

	softDeleted, err := creditRepo.SetInactive(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, softDeleted)
	assert.False(t, softDeleted.IsActive)
}

func TestCreditRepository_SetActive(t *testing.T) {
	pool := testDBPool(t)
	defer pool.Close()
	ctx := context.Background()
	clientRepo := postgres.NewClientRepository(pool)
	bankRepo := postgres.NewBankRepository(pool)
	creditRepo := postgres.NewCreditRepository(pool)

	client, err := clientRepo.Create(ctx, domain.CreateClientInput{FullName: "C", Email: uniqueClientEmail(t), BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), Country: "US"})
	require.NoError(t, err)
	require.NotNil(t, client)
	defer deleteClient(t, pool, client.ID)
	bank, err := bankRepo.Create(ctx, domain.CreateBankInput{Name: "B", Type: domain.BankTypePrivate})
	require.NoError(t, err)
	require.NotNil(t, bank)
	defer deleteBank(t, pool, bank.ID)
	created, err := creditRepo.Create(ctx, domain.CreateCreditInput{
		ClientID: client.ID, BankID: bank.ID, MinPayment: 100, MaxPayment: 500, TermMonths: 12, CreditType: domain.CreditTypeAuto,
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	defer deleteCredit(t, pool, created.ID)

	_, err = creditRepo.SetInactive(ctx, created.ID)
	require.NoError(t, err)

	reenabled, err := creditRepo.SetActive(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, reenabled)
	assert.True(t, reenabled.IsActive)
	assert.Equal(t, created.ID, reenabled.ID)
}

func TestCreditRepository_SetActive_NotFound(t *testing.T) {
	pool := testDBPool(t)
	defer pool.Close()
	ctx := context.Background()
	creditRepo := postgres.NewCreditRepository(pool)

	got, err := creditRepo.SetActive(ctx, "00000000-0000-0000-0000-000000000000")
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestCreditRepository_List(t *testing.T) {
	pool := testDBPool(t)
	defer pool.Close()
	ctx := context.Background()
	clientRepo := postgres.NewClientRepository(pool)
	bankRepo := postgres.NewBankRepository(pool)
	creditRepo := postgres.NewCreditRepository(pool)

	client, err := clientRepo.Create(ctx, domain.CreateClientInput{FullName: "C", Email: uniqueClientEmail(t), BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), Country: "US"})
	require.NoError(t, err)
	require.NotNil(t, client)
	defer deleteClient(t, pool, client.ID)
	bank, err := bankRepo.Create(ctx, domain.CreateBankInput{Name: "B", Type: domain.BankTypePrivate})
	require.NoError(t, err)
	require.NotNil(t, bank)
	defer deleteBank(t, pool, bank.ID)
	credit, err := creditRepo.Create(ctx, domain.CreateCreditInput{
		ClientID: client.ID, BankID: bank.ID, MinPayment: 100, MaxPayment: 500, TermMonths: 12, CreditType: domain.CreditTypeAuto,
	})
	require.NoError(t, err)
	require.NotNil(t, credit)
	defer deleteCredit(t, pool, credit.ID)

	list, err := creditRepo.List(ctx, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)
}

func TestCreditRepository_ListByClientID(t *testing.T) {
	pool := testDBPool(t)
	defer pool.Close()
	ctx := context.Background()
	clientRepo := postgres.NewClientRepository(pool)
	bankRepo := postgres.NewBankRepository(pool)
	creditRepo := postgres.NewCreditRepository(pool)

	client, err := clientRepo.Create(ctx, domain.CreateClientInput{FullName: "ListByClient", Email: uniqueClientEmail(t), BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), Country: "US"})
	require.NoError(t, err)
	require.NotNil(t, client)
	defer deleteClient(t, pool, client.ID)
	bank, err := bankRepo.Create(ctx, domain.CreateBankInput{Name: "ListByClient Bank", Type: domain.BankTypePrivate})
	require.NoError(t, err)
	require.NotNil(t, bank)
	defer deleteBank(t, pool, bank.ID)

	credit, err := creditRepo.Create(ctx, domain.CreateCreditInput{
		ClientID: client.ID, BankID: bank.ID, MinPayment: 100, MaxPayment: 500, TermMonths: 12, CreditType: domain.CreditTypeAuto,
	})
	require.NoError(t, err)
	require.NotNil(t, credit)
	defer deleteCredit(t, pool, credit.ID)

	list, err := creditRepo.ListByClientID(ctx, client.ID, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)
	assert.Equal(t, client.ID, list[0].ClientID)
}
