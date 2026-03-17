package payment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSuccess   PaymentStatus = "success"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

type PaymentMethod string

const (
	PaymentMethodAlipay    PaymentMethod = "alipay"
	PaymentMethodWechat    PaymentMethod = "wechat"
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodPayPal    PaymentMethod = "paypal"
	PaymentMethodBalance   PaymentMethod = "balance"
)

type PaymentType string

const (
	PaymentTypeSubscription PaymentType = "subscription"
	PaymentTypeTip          PaymentType = "tip"
	PaymentTypePurchase     PaymentType = "purchase"
)

type Payment struct {
	ID             int64         `json:"id"`
	UserID         int64         `json:"user_id"`
	OrderID        string        `json:"order_id"`
	Type           PaymentType   `json:"type"`
	Method         PaymentMethod `json:"method"`
	Amount         int64         `json:"amount"`
	Currency       string        `json:"currency"`
	Status         PaymentStatus `json:"status"`
	Description    string        `json:"description"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
	TransactionID  string        `json:"transaction_id,omitempty"`
	PaidAt         *time.Time    `json:"paid_at,omitempty"`
	RefundedAt     *time.Time    `json:"refunded_at,omitempty"`
	RefundAmount   int64         `json:"refund_amount,omitempty"`
	RefundReason   string        `json:"refund_reason,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type PaymentOrder struct {
	ID          string      `json:"id"`
	UserID      int64       `json:"user_id"`
	Type        PaymentType `json:"type"`
	Amount      int64       `json:"amount"`
	Currency    string      `json:"currency"`
	Description string      `json:"description"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	ExpiresAt   time.Time   `json:"expires_at"`
	CreatedAt   time.Time   `json:"created_at"`
}

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *Payment) error
	GetPaymentByID(ctx context.Context, id int64) (*Payment, error)
	GetPaymentByOrderID(ctx context.Context, orderID string) (*Payment, error)
	GetPaymentsByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*Payment, int64, error)
	UpdatePaymentStatus(ctx context.Context, id int64, status PaymentStatus, transactionID string) error
	RefundPayment(ctx context.Context, id int64, amount int64, reason string) error
	
	CreateOrder(ctx context.Context, order *PaymentOrder) error
	GetOrderByID(ctx context.Context, orderID string) (*PaymentOrder, error)
	DeleteOrder(ctx context.Context, orderID string) error
}

type paymentRepository struct {
	pool *pgxpool.Pool
}

func NewPaymentRepository(pool *pgxpool.Pool) PaymentRepository {
	return &paymentRepository{pool: pool}
}

func (r *paymentRepository) CreatePayment(ctx context.Context, payment *Payment) error {
	query := `
		INSERT INTO payments (user_id, order_id, type, method, amount, currency, status, description, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	
	return r.pool.QueryRow(ctx, query,
		payment.UserID, payment.OrderID, payment.Type, payment.Method,
		payment.Amount, payment.Currency, payment.Status, payment.Description, payment.Metadata,
	).Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)
}

func (r *paymentRepository) GetPaymentByID(ctx context.Context, id int64) (*Payment, error) {
	query := `
		SELECT id, user_id, order_id, type, method, amount, currency, status, description, metadata,
			   transaction_id, paid_at, refunded_at, refund_amount, refund_reason, created_at, updated_at
		FROM payments WHERE id = $1
	`
	
	var p Payment
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.UserID, &p.OrderID, &p.Type, &p.Method, &p.Amount, &p.Currency,
		&p.Status, &p.Description, &p.Metadata, &p.TransactionID, &p.PaidAt,
		&p.RefundedAt, &p.RefundAmount, &p.RefundReason, &p.CreatedAt, &p.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	
	return &p, nil
}

func (r *paymentRepository) GetPaymentByOrderID(ctx context.Context, orderID string) (*Payment, error) {
	query := `
		SELECT id, user_id, order_id, type, method, amount, currency, status, description, metadata,
			   transaction_id, paid_at, refunded_at, refund_amount, refund_reason, created_at, updated_at
		FROM payments WHERE order_id = $1
	`
	
	var p Payment
	err := r.pool.QueryRow(ctx, query, orderID).Scan(
		&p.ID, &p.UserID, &p.OrderID, &p.Type, &p.Method, &p.Amount, &p.Currency,
		&p.Status, &p.Description, &p.Metadata, &p.TransactionID, &p.PaidAt,
		&p.RefundedAt, &p.RefundAmount, &p.RefundReason, &p.CreatedAt, &p.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	
	return &p, nil
}

func (r *paymentRepository) GetPaymentsByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*Payment, int64, error) {
	countQuery := `SELECT COUNT(*) FROM payments WHERE user_id = $1`
	var total int64
	r.pool.QueryRow(ctx, countQuery, userID).Scan(&total)
	
	offset := (page - 1) * pageSize
	query := `
		SELECT id, user_id, order_id, type, method, amount, currency, status, description, metadata,
			   transaction_id, paid_at, refunded_at, refund_amount, refund_reason, created_at, updated_at
		FROM payments WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.pool.Query(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var payments []*Payment
	for rows.Next() {
		var p Payment
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.OrderID, &p.Type, &p.Method, &p.Amount, &p.Currency,
			&p.Status, &p.Description, &p.Metadata, &p.TransactionID, &p.PaidAt,
			&p.RefundedAt, &p.RefundAmount, &p.RefundReason, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		payments = append(payments, &p)
	}
	
	return payments, total, nil
}

func (r *paymentRepository) UpdatePaymentStatus(ctx context.Context, id int64, status PaymentStatus, transactionID string) error {
	var paidAt *time.Time
	if status == PaymentStatusSuccess {
		now := time.Now()
		paidAt = &now
	}
	
	query := `UPDATE payments SET status = $2, transaction_id = $3, paid_at = $4, updated_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, status, transactionID, paidAt)
	return err
}

func (r *paymentRepository) RefundPayment(ctx context.Context, id int64, amount int64, reason string) error {
	now := time.Now()
	query := `
		UPDATE payments 
		SET status = 'refunded', refunded_at = $2, refund_amount = $3, refund_reason = $4, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, id, now, amount, reason)
	return err
}

func (r *paymentRepository) CreateOrder(ctx context.Context, order *PaymentOrder) error {
	query := `
		INSERT INTO payment_orders (id, user_id, type, amount, currency, description, metadata, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
	`
	
	_, err := r.pool.Exec(ctx, query,
		order.ID, order.UserID, order.Type, order.Amount, order.Currency,
		order.Description, order.Metadata, order.ExpiresAt,
	)
	return err
}

func (r *paymentRepository) GetOrderByID(ctx context.Context, orderID string) (*PaymentOrder, error) {
	query := `
		SELECT id, user_id, type, amount, currency, description, metadata, expires_at, created_at
		FROM payment_orders WHERE id = $1
	`
	
	var o PaymentOrder
	err := r.pool.QueryRow(ctx, query, orderID).Scan(
		&o.ID, &o.UserID, &o.Type, &o.Amount, &o.Currency,
		&o.Description, &o.Metadata, &o.ExpiresAt, &o.CreatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	
	return &o, nil
}

func (r *paymentRepository) DeleteOrder(ctx context.Context, orderID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM payment_orders WHERE id = $1`, orderID)
	return err
}

type PaymentGateway interface {
	CreatePayment(ctx context.Context, order *PaymentOrder) (string, error)
	QueryPayment(ctx context.Context, orderID string) (*PaymentStatus, error)
	Refund(ctx context.Context, orderID string, amount int64) error
	HandleCallback(ctx context.Context, data []byte) (*PaymentCallback, error)
}

type PaymentCallback struct {
	OrderID       string
	TransactionID string
	Status        PaymentStatus
	Amount        int64
}

type PaymentService interface {
	CreateOrder(ctx context.Context, userID int64, paymentType PaymentType, amount int64, currency, description string, metadata interface{}) (*PaymentOrder, error)
	InitiatePayment(ctx context.Context, orderID string, method PaymentMethod) (string, error)
	HandleCallback(ctx context.Context, method PaymentMethod, data []byte) error
	GetPayment(ctx context.Context, paymentID int64) (*Payment, error)
	GetUserPayments(ctx context.Context, userID int64, page, pageSize int) ([]*Payment, int64, error)
	RefundPayment(ctx context.Context, paymentID int64, amount int64, reason string) error
}

type paymentService struct {
	repo     PaymentRepository
	gateways map[PaymentMethod]PaymentGateway
}

func NewPaymentService(repo PaymentRepository, gateways map[PaymentMethod]PaymentGateway) PaymentService {
	return &paymentService{repo: repo, gateways: gateways}
}

func (s *paymentService) CreateOrder(ctx context.Context, userID int64, paymentType PaymentType, amount int64, currency, description string, metadata interface{}) (*PaymentOrder, error) {
	orderID := generateOrderID()
	
	var metadataBytes json.RawMessage
	if metadata != nil {
		metadataBytes, _ = json.Marshal(metadata)
	}
	
	order := &PaymentOrder{
		ID:          orderID,
		UserID:      userID,
		Type:        paymentType,
		Amount:      amount,
		Currency:    currency,
		Description: description,
		Metadata:    metadataBytes,
		ExpiresAt:   time.Now().Add(30 * time.Minute),
	}
	
	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	
	return order, nil
}

func (s *paymentService) InitiatePayment(ctx context.Context, orderID string, method PaymentMethod) (string, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return "", err
	}
	
	if order == nil {
		return "", errors.New("order not found")
	}
	
	if time.Now().After(order.ExpiresAt) {
		return "", errors.New("order has expired")
	}
	
	gateway, ok := s.gateways[method]
	if !ok {
		return "", errors.New("unsupported payment method")
	}
	
	paymentURL, err := gateway.CreatePayment(ctx, order)
	if err != nil {
		return "", err
	}
	
	payment := &Payment{
		UserID:      order.UserID,
		OrderID:     order.ID,
		Type:        order.Type,
		Method:      method,
		Amount:      order.Amount,
		Currency:    order.Currency,
		Status:      PaymentStatusPending,
		Description: order.Description,
		Metadata:    order.Metadata,
	}
	
	if err := s.repo.CreatePayment(ctx, payment); err != nil {
		return "", err
	}
	
	return paymentURL, nil
}

func (s *paymentService) HandleCallback(ctx context.Context, method PaymentMethod, data []byte) error {
	gateway, ok := s.gateways[method]
	if !ok {
		return errors.New("unsupported payment method")
	}
	
	callback, err := gateway.HandleCallback(ctx, data)
	if err != nil {
		return err
	}
	
	payment, err := s.repo.GetPaymentByOrderID(ctx, callback.OrderID)
	if err != nil {
		return err
	}
	
	if payment == nil {
		return errors.New("payment not found")
	}
	
	return s.repo.UpdatePaymentStatus(ctx, payment.ID, callback.Status, callback.TransactionID)
}

func (s *paymentService) GetPayment(ctx context.Context, paymentID int64) (*Payment, error) {
	return s.repo.GetPaymentByID(ctx, paymentID)
}

func (s *paymentService) GetUserPayments(ctx context.Context, userID int64, page, pageSize int) ([]*Payment, int64, error) {
	return s.repo.GetPaymentsByUserID(ctx, userID, page, pageSize)
}

func (s *paymentService) RefundPayment(ctx context.Context, paymentID int64, amount int64, reason string) error {
	payment, err := s.repo.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return err
	}
	
	if payment == nil {
		return errors.New("payment not found")
	}
	
	if payment.Status != PaymentStatusSuccess {
		return errors.New("payment cannot be refunded")
	}
	
	gateway, ok := s.gateways[payment.Method]
	if !ok {
		return errors.New("unsupported payment method")
	}
	
	if err := gateway.Refund(ctx, payment.OrderID, amount); err != nil {
		return err
	}
	
	return s.repo.RefundPayment(ctx, paymentID, amount, reason)
}

func generateOrderID() string {
	return fmt.Sprintf("ORD%d%d", time.Now().UnixNano(), time.Now().Nanosecond())
}
