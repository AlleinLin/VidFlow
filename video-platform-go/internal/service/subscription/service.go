package subscription

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionTier string

const (
	TierFree      SubscriptionTier = "free"
	TierBasic     SubscriptionTier = "basic"
	TierPremium   SubscriptionTier = "premium"
	TierEnterprise SubscriptionTier = "enterprise"
)

type SubscriptionStatus string

const (
	StatusActive    SubscriptionStatus = "active"
	StatusExpired   SubscriptionStatus = "expired"
	StatusCancelled SubscriptionStatus = "cancelled"
	StatusPaused    SubscriptionStatus = "paused"
)

type SubscriptionFeature string

const (
	FeatureHD          SubscriptionFeature = "hd"
	Feature4K          SubscriptionFeature = "4k"
	FeatureNoAds       SubscriptionFeature = "no_ads"
	FeatureDownload    SubscriptionFeature = "download"
	FeatureExclusive   SubscriptionFeature = "exclusive"
	FeaturePriority    SubscriptionFeature = "priority"
	FeatureUnlimited   SubscriptionFeature = "unlimited"
)

type SubscriptionPlan struct {
	ID          int64              `json:"id"`
	Name        string             `json:"name"`
	Tier        SubscriptionTier   `json:"tier"`
	Price       int64              `json:"price"`
	Currency    string             `json:"currency"`
	Duration    int                `json:"duration"`
	Features    []SubscriptionFeature `json:"features"`
	Description string             `json:"description"`
	IsActive    bool               `json:"is_active"`
	CreatedAt   time.Time          `json:"created_at"`
}

type Subscription struct {
	ID          int64              `json:"id"`
	UserID      int64              `json:"user_id"`
	PlanID      int64              `json:"plan_id"`
	Status      SubscriptionStatus `json:"status"`
	StartedAt   time.Time          `json:"started_at"`
	ExpiresAt   time.Time          `json:"expires_at"`
	CancelledAt *time.Time         `json:"cancelled_at,omitempty"`
	AutoRenew   bool               `json:"auto_renew"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type UserBenefits struct {
	Tier           SubscriptionTier   `json:"tier"`
	Features       []SubscriptionFeature `json:"features"`
	MaxResolution  string             `json:"max_resolution"`
	DownloadLimit  int                `json:"download_limit"`
	StorageLimit   int64              `json:"storage_limit"`
	NoAds          bool               `json:"no_ads"`
	Priority       bool               `json:"priority"`
	ExpiresAt      *time.Time         `json:"expires_at,omitempty"`
}

type SubscriptionRepository interface {
	CreatePlan(ctx context.Context, plan *SubscriptionPlan) error
	GetPlanByID(ctx context.Context, id int64) (*SubscriptionPlan, error)
	GetActivePlans(ctx context.Context) ([]*SubscriptionPlan, error)
	UpdatePlan(ctx context.Context, plan *SubscriptionPlan) error
	
	CreateSubscription(ctx context.Context, sub *Subscription) error
	GetSubscriptionByID(ctx context.Context, id int64) (*Subscription, error)
	GetActiveSubscriptionByUserID(ctx context.Context, userID int64) (*Subscription, error)
	GetSubscriptionsByUserID(ctx context.Context, userID int64) ([]*Subscription, error)
	UpdateSubscription(ctx context.Context, sub *Subscription) error
	CancelSubscription(ctx context.Context, id int64, reason string) error
	RenewSubscription(ctx context.Context, id int64) error
}

type subscriptionRepository struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepository(pool *pgxpool.Pool) SubscriptionRepository {
	return &subscriptionRepository{pool: pool}
}

func (r *subscriptionRepository) CreatePlan(ctx context.Context, plan *SubscriptionPlan) error {
	featuresJSON, _ := json.Marshal(plan.Features)
	query := `
		INSERT INTO subscription_plans (name, tier, price, currency, duration, features, description, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		RETURNING id, created_at
	`
	return r.pool.QueryRow(ctx, query,
		plan.Name, plan.Tier, plan.Price, plan.Currency, plan.Duration,
		featuresJSON, plan.Description, plan.IsActive,
	).Scan(&plan.ID, &plan.CreatedAt)
}

func (r *subscriptionRepository) GetPlanByID(ctx context.Context, id int64) (*SubscriptionPlan, error) {
	query := `
		SELECT id, name, tier, price, currency, duration, features, description, is_active, created_at
		FROM subscription_plans WHERE id = $1
	`
	
	var plan SubscriptionPlan
	var featuresJSON []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&plan.ID, &plan.Name, &plan.Tier, &plan.Price, &plan.Currency,
		&plan.Duration, &featuresJSON, &plan.Description, &plan.IsActive, &plan.CreatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	
	json.Unmarshal(featuresJSON, &plan.Features)
	return &plan, nil
}

func (r *subscriptionRepository) GetActivePlans(ctx context.Context) ([]*SubscriptionPlan, error) {
	query := `
		SELECT id, name, tier, price, currency, duration, features, description, is_active, created_at
		FROM subscription_plans WHERE is_active = true ORDER BY price
	`
	
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var plans []*SubscriptionPlan
	for rows.Next() {
		var plan SubscriptionPlan
		var featuresJSON []byte
		if err := rows.Scan(
			&plan.ID, &plan.Name, &plan.Tier, &plan.Price, &plan.Currency,
			&plan.Duration, &featuresJSON, &plan.Description, &plan.IsActive, &plan.CreatedAt,
		); err != nil {
			return nil, err
		}
		json.Unmarshal(featuresJSON, &plan.Features)
		plans = append(plans, &plan)
	}
	
	return plans, nil
}

func (r *subscriptionRepository) UpdatePlan(ctx context.Context, plan *SubscriptionPlan) error {
	featuresJSON, _ := json.Marshal(plan.Features)
	query := `
		UPDATE subscription_plans SET name = $2, tier = $3, price = $4, currency = $5,
		duration = $6, features = $7, description = $8, is_active = $9
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query,
		plan.ID, plan.Name, plan.Tier, plan.Price, plan.Currency,
		plan.Duration, featuresJSON, plan.Description, plan.IsActive,
	)
	return err
}

func (r *subscriptionRepository) CreateSubscription(ctx context.Context, sub *Subscription) error {
	query := `
		INSERT INTO subscriptions (user_id, plan_id, status, started_at, expires_at, auto_renew, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query,
		sub.UserID, sub.PlanID, sub.Status, sub.StartedAt, sub.ExpiresAt, sub.AutoRenew,
	).Scan(&sub.ID, &sub.CreatedAt, &sub.UpdatedAt)
}

func (r *subscriptionRepository) GetSubscriptionByID(ctx context.Context, id int64) (*Subscription, error) {
	query := `
		SELECT id, user_id, plan_id, status, started_at, expires_at, cancelled_at, auto_renew, created_at, updated_at
		FROM subscriptions WHERE id = $1
	`
	
	var sub Subscription
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&sub.ID, &sub.UserID, &sub.PlanID, &sub.Status, &sub.StartedAt,
		&sub.ExpiresAt, &sub.CancelledAt, &sub.AutoRenew, &sub.CreatedAt, &sub.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	
	return &sub, nil
}

func (r *subscriptionRepository) GetActiveSubscriptionByUserID(ctx context.Context, userID int64) (*Subscription, error) {
	query := `
		SELECT id, user_id, plan_id, status, started_at, expires_at, cancelled_at, auto_renew, created_at, updated_at
		FROM subscriptions 
		WHERE user_id = $1 AND status = 'active' AND expires_at > NOW()
		ORDER BY expires_at DESC
		LIMIT 1
	`
	
	var sub Subscription
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&sub.ID, &sub.UserID, &sub.PlanID, &sub.Status, &sub.StartedAt,
		&sub.ExpiresAt, &sub.CancelledAt, &sub.AutoRenew, &sub.CreatedAt, &sub.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	
	return &sub, nil
}

func (r *subscriptionRepository) GetSubscriptionsByUserID(ctx context.Context, userID int64) ([]*Subscription, error) {
	query := `
		SELECT id, user_id, plan_id, status, started_at, expires_at, cancelled_at, auto_renew, created_at, updated_at
		FROM subscriptions WHERE user_id = $1 ORDER BY created_at DESC
	`
	
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var subs []*Subscription
	for rows.Next() {
		var sub Subscription
		if err := rows.Scan(
			&sub.ID, &sub.UserID, &sub.PlanID, &sub.Status, &sub.StartedAt,
			&sub.ExpiresAt, &sub.CancelledAt, &sub.AutoRenew, &sub.CreatedAt, &sub.UpdatedAt,
		); err != nil {
			return nil, err
		}
		subs = append(subs, &sub)
	}
	
	return subs, nil
}

func (r *subscriptionRepository) UpdateSubscription(ctx context.Context, sub *Subscription) error {
	query := `
		UPDATE subscriptions SET status = $2, expires_at = $3, auto_renew = $4, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, sub.ID, sub.Status, sub.ExpiresAt, sub.AutoRenew)
	return err
}

func (r *subscriptionRepository) CancelSubscription(ctx context.Context, id int64, reason string) error {
	now := time.Now()
	query := `
		UPDATE subscriptions SET status = 'cancelled', cancelled_at = $2, auto_renew = false, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, id, now)
	return err
}

func (r *subscriptionRepository) RenewSubscription(ctx context.Context, id int64) error {
	sub, err := r.GetSubscriptionByID(ctx, id)
	if err != nil {
		return err
	}
	
	plan, err := r.GetPlanByID(ctx, sub.PlanID)
	if err != nil {
		return err
	}
	
	newExpiresAt := sub.ExpiresAt.AddDate(0, 0, plan.Duration)
	query := `
		UPDATE subscriptions SET expires_at = $2, status = 'active', updated_at = NOW()
		WHERE id = $1
	`
	_, err = r.pool.Exec(ctx, query, id, newExpiresAt)
	return err
}

type SubscriptionService interface {
	GetPlans(ctx context.Context) ([]*SubscriptionPlan, error)
	GetPlan(ctx context.Context, planID int64) (*SubscriptionPlan, error)
	Subscribe(ctx context.Context, userID, planID int64, autoRenew bool) (*Subscription, error)
	GetActiveSubscription(ctx context.Context, userID int64) (*Subscription, error)
	GetSubscriptionHistory(ctx context.Context, userID int64) ([]*Subscription, error)
	CancelSubscription(ctx context.Context, userID, subscriptionID int64) error
	RenewSubscription(ctx context.Context, userID, subscriptionID int64) error
	GetUserBenefits(ctx context.Context, userID int64) (*UserBenefits, error)
	CheckFeatureAccess(ctx context.Context, userID int64, feature SubscriptionFeature) (bool, error)
}

type subscriptionService struct {
	repo SubscriptionRepository
}

func NewSubscriptionService(repo SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) GetPlans(ctx context.Context) ([]*SubscriptionPlan, error) {
	return s.repo.GetActivePlans(ctx)
}

func (s *subscriptionService) GetPlan(ctx context.Context, planID int64) (*SubscriptionPlan, error) {
	return s.repo.GetPlanByID(ctx, planID)
}

func (s *subscriptionService) Subscribe(ctx context.Context, userID, planID int64, autoRenew bool) (*Subscription, error) {
	plan, err := s.repo.GetPlanByID(ctx, planID)
	if err != nil {
		return nil, err
	}
	
	if plan == nil || !plan.IsActive {
		return nil, errors.New("plan not found or inactive")
	}
	
	activeSub, _ := s.repo.GetActiveSubscriptionByUserID(ctx, userID)
	if activeSub != nil {
		return nil, errors.New("user already has an active subscription")
	}
	
	now := time.Now()
	expiresAt := now.AddDate(0, 0, plan.Duration)
	
	sub := &Subscription{
		UserID:    userID,
		PlanID:    planID,
		Status:    StatusActive,
		StartedAt: now,
		ExpiresAt: expiresAt,
		AutoRenew: autoRenew,
	}
	
	if err := s.repo.CreateSubscription(ctx, sub); err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}
	
	return sub, nil
}

func (s *subscriptionService) GetActiveSubscription(ctx context.Context, userID int64) (*Subscription, error) {
	return s.repo.GetActiveSubscriptionByUserID(ctx, userID)
}

func (s *subscriptionService) GetSubscriptionHistory(ctx context.Context, userID int64) ([]*Subscription, error) {
	return s.repo.GetSubscriptionsByUserID(ctx, userID)
}

func (s *subscriptionService) CancelSubscription(ctx context.Context, userID, subscriptionID int64) error {
	sub, err := s.repo.GetSubscriptionByID(ctx, subscriptionID)
	if err != nil {
		return err
	}
	
	if sub == nil || sub.UserID != userID {
		return errors.New("subscription not found")
	}
	
	return s.repo.CancelSubscription(ctx, subscriptionID, "")
}

func (s *subscriptionService) RenewSubscription(ctx context.Context, userID, subscriptionID int64) error {
	sub, err := s.repo.GetSubscriptionByID(ctx, subscriptionID)
	if err != nil {
		return err
	}
	
	if sub == nil || sub.UserID != userID {
		return errors.New("subscription not found")
	}
	
	return s.repo.RenewSubscription(ctx, subscriptionID)
}

func (s *subscriptionService) GetUserBenefits(ctx context.Context, userID int64) (*UserBenefits, error) {
	sub, err := s.repo.GetActiveSubscriptionByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	benefits := &UserBenefits{
		Tier:          TierFree,
		Features:      []SubscriptionFeature{},
		MaxResolution: "480p",
		DownloadLimit: 0,
		StorageLimit:  100 * 1024 * 1024,
		NoAds:         false,
		Priority:      false,
	}
	
	if sub != nil {
		plan, err := s.repo.GetPlanByID(ctx, sub.PlanID)
		if err != nil {
			return nil, err
		}
		
		if plan != nil {
			benefits.Tier = plan.Tier
			benefits.Features = plan.Features
			benefits.ExpiresAt = &sub.ExpiresAt
			
			for _, feature := range plan.Features {
				switch feature {
				case FeatureHD:
					benefits.MaxResolution = "1080p"
				case Feature4K:
					benefits.MaxResolution = "4k"
				case FeatureNoAds:
					benefits.NoAds = true
				case FeatureDownload:
					benefits.DownloadLimit = 100
				case FeaturePriority:
					benefits.Priority = true
				case FeatureUnlimited:
					benefits.StorageLimit = -1
					benefits.DownloadLimit = -1
				}
			}
		}
	}
	
	return benefits, nil
}

func (s *subscriptionService) CheckFeatureAccess(ctx context.Context, userID int64, feature SubscriptionFeature) (bool, error) {
	benefits, err := s.GetUserBenefits(ctx, userID)
	if err != nil {
		return false, err
	}
	
	for _, f := range benefits.Features {
		if f == feature {
			return true, nil
		}
	}
	
	return false, nil
}
