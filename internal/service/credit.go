package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tucredito/backend-api/internal/cache"
	"github.com/tucredito/backend-api/internal/decision"
	"github.com/tucredito/backend-api/internal/domain"
	"github.com/tucredito/backend-api/internal/event"
	"github.com/tucredito/backend-api/internal/metrics"
	"github.com/tucredito/backend-api/internal/repository"
	"go.uber.org/zap"
)

var (
	ErrClientNotFound = errors.New("client not found")
	ErrBankNotFound   = errors.New("bank not found")
	ErrInvalidInput   = errors.New("invalid input")
)

const (
	cacheTTLSeconds      = 300
	creditCacheKeyPrefix = "credit:"
	workerPoolSize       = 10
)

type CreditService struct {
	creditRepo repository.CreditRepository
	clientRepo repository.ClientRepository
	bankRepo   repository.BankRepository
	cache      cache.Cache
	publisher  event.Publisher
	engine     decision.Engine
	log        *zap.Logger
	jobCh      chan creditJob
	done       chan struct{}
	wg         sync.WaitGroup
}

type creditJob struct {
	ctx    context.Context
	input  domain.CreateCreditInput
	result chan creditResult
}

type creditResult struct {
	credit *domain.Credit
	err    error
}

// Creates a CreditService and starts the worker pool
func NewCreditService(
	creditRepo repository.CreditRepository,
	clientRepo repository.ClientRepository,
	bankRepo repository.BankRepository,
	cache cache.Cache,
	publisher event.Publisher,
	engine decision.Engine,
	log *zap.Logger,
) *CreditService {
	s := &CreditService{
		creditRepo: creditRepo,
		clientRepo: clientRepo,
		bankRepo:   bankRepo,
		cache:      cache,
		publisher:  publisher,
		engine:     engine,
		log:        log,
		jobCh:      make(chan creditJob, 100),
		done:       make(chan struct{}),
	}
	for i := 0; i < workerPoolSize; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}
	return s
}

// Processes credit creation jobs from the channel (concurrent processing)
func (s *CreditService) worker(id int) {
	defer s.wg.Done()
	for {
		select {
		case <-s.done:
			return
		case job := <-s.jobCh:
			credit, err := s.createCredit(job.ctx, job.input)
			job.result <- creditResult{credit: credit, err: err}
		}
	}
}

// Enqueues credit creation and returns the result
func (s *CreditService) Create(ctx context.Context, input domain.CreateCreditInput) (*domain.Credit, error) {
	if input.ClientID == "" || input.BankID == "" || input.MaxPayment < input.MinPayment || input.TermMonths <= 0 {
		return nil, ErrInvalidInput
	}

	resultCh := make(chan creditResult, 1)

	select {
	case s.jobCh <- creditJob{ctx: ctx, input: input, result: resultCh}:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	select {
	case res := <-resultCh:
		return res.credit, res.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Runs validations, eligibility, persistence, cache, and events
func (s *CreditService) createCredit(ctx context.Context, input domain.CreateCreditInput) (*domain.Credit, error) {
	client, err := s.clientRepo.GetByID(ctx, input.ClientID)
	if err != nil {
		return nil, err
	}

	if client == nil {
		return nil, ErrClientNotFound
	}

	bank, err := s.bankRepo.GetByID(ctx, input.BankID)
	if err != nil {
		return nil, err
	}

	if bank == nil {
		return nil, ErrBankNotFound
	}

	eligibilityInput := &decision.EligibilityInput{
		Client:     client,
		Bank:       bank,
		MinPayment: input.MinPayment,
		MaxPayment: input.MaxPayment,
		TermMonths: input.TermMonths,
		CreditType: input.CreditType,
	}

	result, err := s.engine.Evaluate(ctx, eligibilityInput)
	if err != nil {
		return nil, err
	}

	credit, err := s.creditRepo.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	metrics.IncCreditsCreated()

	if result != nil && result.Approved {
		credit, err = s.creditRepo.UpdateStatus(ctx, credit.ID, domain.CreditStatusApproved)
		if err != nil {
			s.log.Warn("failed to update credit status to approved", zap.Error(err), zap.String("credit_id", credit.ID))
		} else {
			metrics.IncCreditsApproved()
			_ = s.emitCreditApproved(ctx, credit)
		}
	}

	_ = s.emitCreditCreated(ctx, credit)
	s.cacheCredit(ctx, credit)

	return credit, nil
}

func (s *CreditService) cacheCredit(ctx context.Context, c *domain.Credit) {
	if s.cache == nil {
		return
	}

	key := creditCacheKeyPrefix + c.ID
	if cacheWithJSON, ok := s.cache.(interface {
		SetJSON(context.Context, string, interface{}, int) error
	}); ok {
		_ = cacheWithJSON.SetJSON(ctx, key, c, cacheTTLSeconds)
	}
}

// Returns a credit by ID, checking cache first
func (s *CreditService) GetByID(ctx context.Context, id string) (*domain.Credit, error) {
	if s.cache != nil {
		if cacheWithJSON, ok := s.cache.(interface {
			GetJSON(context.Context, string, interface{}) error
		}); ok {
			var c domain.Credit
			if err := cacheWithJSON.GetJSON(ctx, creditCacheKeyPrefix+id, &c); err == nil && c.ID != "" {
				return &c, nil
			}
		}
	}

	return s.creditRepo.GetByID(ctx, id)
}

// Updates credit status and emits domain events
func (s *CreditService) UpdateStatus(ctx context.Context, id string, status domain.CreditStatus) (*domain.Credit, error) {
	credit, err := s.creditRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		return nil, err
	}

	if credit == nil {
		return nil, nil
	}

	switch status {
	case domain.CreditStatusApproved:
		metrics.IncCreditsApproved()
		_ = s.emitCreditApproved(ctx, credit)
	case domain.CreditStatusRejected:
		metrics.IncCreditsRejected()
		_ = s.emitCreditRejected(ctx, credit)
	}

	if s.cache != nil {
		_ = s.cache.Delete(ctx, creditCacheKeyPrefix+id)
	}

	return credit, nil
}

// List returns credits with pagination
func (s *CreditService) List(ctx context.Context, limit, offset int) ([]*domain.Credit, error) {
	return s.creditRepo.List(ctx, limit, offset)
}

// Returns credits for a client
func (s *CreditService) ListByClientID(ctx context.Context, clientID string, limit, offset int) ([]*domain.Credit, error) {
	return s.creditRepo.ListByClientID(ctx, clientID, limit, offset)
}

// Stops the worker pool gracefully
func (s *CreditService) Shutdown() {
	close(s.done)
	s.wg.Wait()
}

func (s *CreditService) emitCreditCreated(ctx context.Context, c *domain.Credit) error {
	payload := domain.CreditCreatedPayload{
		CreditID:   c.ID,
		ClientID:   c.ClientID,
		BankID:     c.BankID,
		CreditType: c.CreditType,
		Status:     c.Status,
		CreatedAt:  c.CreatedAt,
	}

	payloadBytes, err := domain.MarshalPayload(payload)
	if err != nil {
		return err
	}

	evt := &domain.DomainEvent{
		ID:         uuid.New().String(),
		Type:       domain.EventCreditCreated,
		Payload:    payloadBytes,
		OccurredAt: time.Now().UTC(),
	}

	return s.publisher.Publish(ctx, evt)
}

func (s *CreditService) emitCreditApproved(ctx context.Context, c *domain.Credit) error {
	payload := domain.CreditApprovedPayload{
		CreditID:   c.ID,
		ClientID:   c.ClientID,
		BankID:     c.BankID,
		ApprovedAt: time.Now().UTC(),
	}

	payloadBytes, err := domain.MarshalPayload(payload)
	if err != nil {
		return err
	}

	evt := &domain.DomainEvent{
		ID:         uuid.New().String(),
		Type:       domain.EventCreditApproved,
		Payload:    payloadBytes,
		OccurredAt: time.Now().UTC(),
	}

	return s.publisher.Publish(ctx, evt)
}

func (s *CreditService) emitCreditRejected(ctx context.Context, c *domain.Credit) error {
	payload := domain.CreditRejectedPayload{
		CreditID:   c.ID,
		ClientID:   c.ClientID,
		BankID:     c.BankID,
		RejectedAt: time.Now().UTC(),
	}

	payloadBytes, err := domain.MarshalPayload(payload)
	if err != nil {
		return err
	}

	evt := &domain.DomainEvent{
		ID:         uuid.New().String(),
		Type:       domain.EventCreditRejected,
		Payload:    payloadBytes,
		OccurredAt: time.Now().UTC(),
	}

	return s.publisher.Publish(ctx, evt)
}

// Creates a credit synchronously (no worker pool)
func (s *CreditService) CreateSync(ctx context.Context, input domain.CreateCreditInput) (*domain.Credit, error) {
	if input.ClientID == "" || input.BankID == "" || input.MaxPayment < input.MinPayment || input.TermMonths <= 0 {
		return nil, ErrInvalidInput
	}

	return s.createCredit(ctx, input)
}

// Runs client/bank fetch and eligibility concurrently
func (s *CreditService) ValidateEligibilityConcurrent(ctx context.Context, input domain.CreateCreditInput) (*decision.EligibilityResult, error) {
	var client *domain.Client
	var bank *domain.Bank
	var clientErr, bankErr error
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		client, clientErr = s.clientRepo.GetByID(ctx, input.ClientID)
	}()

	go func() {
		defer wg.Done()
		bank, bankErr = s.bankRepo.GetByID(ctx, input.BankID)
	}()

	wg.Wait()
	if clientErr != nil {
		return nil, clientErr
	}

	if bankErr != nil {
		return nil, bankErr
	}

	if client == nil {
		return nil, ErrClientNotFound
	}

	if bank == nil {
		return nil, ErrBankNotFound
	}

	eligibilityInput := &decision.EligibilityInput{
		Client:     client,
		Bank:       bank,
		MinPayment: input.MinPayment,
		MaxPayment: input.MaxPayment,
		TermMonths: input.TermMonths,
		CreditType: input.CreditType,
	}

	return s.engine.Evaluate(ctx, eligibilityInput)
}
