package notification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type NotificationType string

const (
	NotificationTypeSystem      NotificationType = "system"
	NotificationTypeLike        NotificationType = "like"
	NotificationTypeComment     NotificationType = "comment"
	NotificationTypeFollow      NotificationType = "follow"
	NotificationTypeMention     NotificationType = "mention"
	NotificationTypeReply       NotificationType = "reply"
	NotificationTypeVideoReady  NotificationType = "video_ready"
	NotificationTypeVideoAudit  NotificationType = "video_audit"
	NotificationTypeSubscription NotificationType = "subscription"
	NotificationTypePayment     NotificationType = "payment"
)

type Notification struct {
	ID          int64           `json:"id"`
	UserID      int64           `json:"user_id"`
	Type        NotificationType `json:"type"`
	Title       string          `json:"title"`
	Content     string          `json:"content"`
	Data        json.RawMessage `json:"data,omitempty"`
	IsRead      bool            `json:"is_read"`
	ReadAt      *time.Time      `json:"read_at,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

type NotificationPreference struct {
	UserID              int64 `json:"user_id"`
	SystemEnabled       bool  `json:"system_enabled"`
	LikeEnabled         bool  `json:"like_enabled"`
	CommentEnabled      bool  `json:"comment_enabled"`
	FollowEnabled       bool  `json:"follow_enabled"`
	MentionEnabled      bool  `json:"mention_enabled"`
	ReplyEnabled        bool  `json:"reply_enabled"`
	VideoReadyEnabled   bool  `json:"video_ready_enabled"`
	SubscriptionEnabled bool  `json:"subscription_enabled"`
	PaymentEnabled      bool  `json:"payment_enabled"`
	EmailEnabled        bool  `json:"email_enabled"`
	PushEnabled         bool  `json:"push_enabled"`
}

type NotificationRepository interface {
	Create(ctx context.Context, notification *Notification) error
	GetByID(ctx context.Context, id int64) (*Notification, error)
	GetByUserID(ctx context.Context, userID int64, page, pageSize int, unreadOnly bool) ([]*Notification, int64, error)
	MarkAsRead(ctx context.Context, id int64) error
	MarkAllAsRead(ctx context.Context, userID int64) error
	Delete(ctx context.Context, id int64) error
	DeleteAll(ctx context.Context, userID int64) error
	GetUnreadCount(ctx context.Context, userID int64) (int64, error)
	
	GetPreference(ctx context.Context, userID int64) (*NotificationPreference, error)
	UpsertPreference(ctx context.Context, pref *NotificationPreference) error
}

type notificationRepository struct {
	pool  *pgxpool.Pool
	redis *redis.Client
}

func NewNotificationRepository(pool *pgxpool.Pool, redis *redis.Client) NotificationRepository {
	return &notificationRepository{pool: pool, redis: redis}
}

func (r *notificationRepository) Create(ctx context.Context, notification *Notification) error {
	query := `
		INSERT INTO notifications (user_id, type, title, content, data, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`
	
	return r.pool.QueryRow(ctx, query,
		notification.UserID, notification.Type, notification.Title, notification.Content,
		notification.Data, false, time.Now(),
	).Scan(&notification.ID, &notification.CreatedAt)
}

func (r *notificationRepository) GetByID(ctx context.Context, id int64) (*Notification, error) {
	query := `
		SELECT id, user_id, type, title, content, data, is_read, read_at, created_at
		FROM notifications WHERE id = $1
	`
	
	var n Notification
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&n.ID, &n.UserID, &n.Type, &n.Title, &n.Content, &n.Data,
		&n.IsRead, &n.ReadAt, &n.CreatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	
	return &n, nil
}

func (r *notificationRepository) GetByUserID(ctx context.Context, userID int64, page, pageSize int, unreadOnly bool) ([]*Notification, int64, error) {
	whereClause := "WHERE user_id = $1"
	args := []interface{}{userID}
	argNum := 2
	
	if unreadOnly {
		whereClause += fmt.Sprintf(" AND is_read = $%d", argNum)
		args = append(args, false)
		argNum++
	}
	
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM notifications %s", whereClause)
	var total int64
	r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT id, user_id, type, title, content, data, is_read, read_at, created_at
		FROM notifications %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)
	
	args = append(args, pageSize, offset)
	
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var notifications []*Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(
			&n.ID, &n.UserID, &n.Type, &n.Title, &n.Content, &n.Data,
			&n.IsRead, &n.ReadAt, &n.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		notifications = append(notifications, &n)
	}
	
	return notifications, total, nil
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id int64) error {
	query := `UPDATE notifications SET is_read = true, read_at = $2 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, time.Now())
	return err
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID int64) error {
	query := `UPDATE notifications SET is_read = true, read_at = $2 WHERE user_id = $1 AND is_read = false`
	_, err := r.pool.Exec(ctx, query, userID, time.Now())
	return err
}

func (r *notificationRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM notifications WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *notificationRepository) DeleteAll(ctx context.Context, userID int64) error {
	query := `DELETE FROM notifications WHERE user_id = $1`
	_, err := r.pool.Exec(ctx, query, userID)
	return err
}

func (r *notificationRepository) GetUnreadCount(ctx context.Context, userID int64) (int64, error) {
	cacheKey := fmt.Sprintf("notification:unread:%d", userID)
	
	count, err := r.redis.Get(ctx, cacheKey).Int64()
	if err == nil {
		return count, nil
	}
	
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`
	var dbCount int64
	r.pool.QueryRow(ctx, query, userID).Scan(&dbCount)
	
	r.redis.Set(ctx, cacheKey, dbCount, 5*time.Minute)
	
	return dbCount, nil
}

func (r *notificationRepository) GetPreference(ctx context.Context, userID int64) (*NotificationPreference, error) {
	query := `
		SELECT user_id, system_enabled, like_enabled, comment_enabled, follow_enabled,
			   mention_enabled, reply_enabled, video_ready_enabled, subscription_enabled,
			   payment_enabled, email_enabled, push_enabled
		FROM notification_preferences WHERE user_id = $1
	`
	
	var pref NotificationPreference
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&pref.UserID, &pref.SystemEnabled, &pref.LikeEnabled, &pref.CommentEnabled,
		&pref.FollowEnabled, &pref.MentionEnabled, &pref.ReplyEnabled, &pref.VideoReadyEnabled,
		&pref.SubscriptionEnabled, &pref.PaymentEnabled, &pref.EmailEnabled, &pref.PushEnabled,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &NotificationPreference{
				UserID:            userID,
				SystemEnabled:     true,
				LikeEnabled:       true,
				CommentEnabled:    true,
				FollowEnabled:     true,
				MentionEnabled:    true,
				ReplyEnabled:      true,
				VideoReadyEnabled: true,
				SubscriptionEnabled: true,
				PaymentEnabled:    true,
				EmailEnabled:      true,
				PushEnabled:       true,
			}, nil
		}
		return nil, err
	}
	
	return &pref, nil
}

func (r *notificationRepository) UpsertPreference(ctx context.Context, pref *NotificationPreference) error {
	query := `
		INSERT INTO notification_preferences (
			user_id, system_enabled, like_enabled, comment_enabled, follow_enabled,
			mention_enabled, reply_enabled, video_ready_enabled, subscription_enabled,
			payment_enabled, email_enabled, push_enabled
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (user_id) DO UPDATE SET
			system_enabled = EXCLUDED.system_enabled,
			like_enabled = EXCLUDED.like_enabled,
			comment_enabled = EXCLUDED.comment_enabled,
			follow_enabled = EXCLUDED.follow_enabled,
			mention_enabled = EXCLUDED.mention_enabled,
			reply_enabled = EXCLUDED.reply_enabled,
			video_ready_enabled = EXCLUDED.video_ready_enabled,
			subscription_enabled = EXCLUDED.subscription_enabled,
			payment_enabled = EXCLUDED.payment_enabled,
			email_enabled = EXCLUDED.email_enabled,
			push_enabled = EXCLUDED.push_enabled
	`
	
	_, err := r.pool.Exec(ctx, query,
		pref.UserID, pref.SystemEnabled, pref.LikeEnabled, pref.CommentEnabled,
		pref.FollowEnabled, pref.MentionEnabled, pref.ReplyEnabled, pref.VideoReadyEnabled,
		pref.SubscriptionEnabled, pref.PaymentEnabled, pref.EmailEnabled, pref.PushEnabled,
	)
	
	return err
}

type NotificationService interface {
	Send(ctx context.Context, userID int64, notificationType NotificationType, title, content string, data interface{}) error
	GetNotifications(ctx context.Context, userID int64, page, pageSize int, unreadOnly bool) ([]*Notification, int64, error)
	MarkAsRead(ctx context.Context, userID, notificationID int64) error
	MarkAllAsRead(ctx context.Context, userID int64) error
	DeleteNotification(ctx context.Context, userID, notificationID int64) error
	GetUnreadCount(ctx context.Context, userID int64) (int64, error)
	GetPreference(ctx context.Context, userID int64) (*NotificationPreference, error)
	UpdatePreference(ctx context.Context, pref *NotificationPreference) error
	Broadcast(ctx context.Context, notificationType NotificationType, title, content string, data interface{}) error
}

type notificationService struct {
	repo NotificationRepository
	redis *redis.Client
}

func NewNotificationService(repo NotificationRepository, redis *redis.Client) NotificationService {
	return &notificationService{repo: repo, redis: redis}
}

func (s *notificationService) Send(ctx context.Context, userID int64, notificationType NotificationType, title, content string, data interface{}) error {
	pref, err := s.repo.GetPreference(ctx, userID)
	if err != nil {
		return err
	}
	
	if !s.isNotificationEnabled(pref, notificationType) {
		return nil
	}
	
	var dataBytes json.RawMessage
	if data != nil {
		dataBytes, _ = json.Marshal(data)
	}
	
	notification := &Notification{
		UserID:  userID,
		Type:    notificationType,
		Title:   title,
		Content: content,
		Data:    dataBytes,
	}
	
	if err := s.repo.Create(ctx, notification); err != nil {
		return err
	}
	
	s.invalidateUnreadCache(ctx, userID)
	
	if pref.PushEnabled {
		go s.sendPushNotification(userID, notification)
	}
	
	if pref.EmailEnabled {
		go s.sendEmailNotification(userID, notification)
	}
	
	return nil
}

func (s *notificationService) GetNotifications(ctx context.Context, userID int64, page, pageSize int, unreadOnly bool) ([]*Notification, int64, error) {
	return s.repo.GetByUserID(ctx, userID, page, pageSize, unreadOnly)
}

func (s *notificationService) MarkAsRead(ctx context.Context, userID, notificationID int64) error {
	if err := s.repo.MarkAsRead(ctx, notificationID); err != nil {
		return err
	}
	
	s.invalidateUnreadCache(ctx, userID)
	return nil
}

func (s *notificationService) MarkAllAsRead(ctx context.Context, userID int64) error {
	if err := s.repo.MarkAllAsRead(ctx, userID); err != nil {
		return err
	}
	
	s.invalidateUnreadCache(ctx, userID)
	return nil
}

func (s *notificationService) DeleteNotification(ctx context.Context, userID, notificationID int64) error {
	if err := s.repo.Delete(ctx, notificationID); err != nil {
		return err
	}
	
	s.invalidateUnreadCache(ctx, userID)
	return nil
}

func (s *notificationService) GetUnreadCount(ctx context.Context, userID int64) (int64, error) {
	return s.repo.GetUnreadCount(ctx, userID)
}

func (s *notificationService) GetPreference(ctx context.Context, userID int64) (*NotificationPreference, error) {
	return s.repo.GetPreference(ctx, userID)
}

func (s *notificationService) UpdatePreference(ctx context.Context, pref *NotificationPreference) error {
	return s.repo.UpsertPreference(ctx, pref)
}

func (s *notificationService) Broadcast(ctx context.Context, notificationType NotificationType, title, content string, data interface{}) error {
	return nil
}

func (s *notificationService) isNotificationEnabled(pref *NotificationPreference, notificationType NotificationType) bool {
	switch notificationType {
	case NotificationTypeSystem:
		return pref.SystemEnabled
	case NotificationTypeLike:
		return pref.LikeEnabled
	case NotificationTypeComment:
		return pref.CommentEnabled
	case NotificationTypeFollow:
		return pref.FollowEnabled
	case NotificationTypeMention:
		return pref.MentionEnabled
	case NotificationTypeReply:
		return pref.ReplyEnabled
	case NotificationTypeVideoReady:
		return pref.VideoReadyEnabled
	case NotificationTypeSubscription:
		return pref.SubscriptionEnabled
	case NotificationTypePayment:
		return pref.PaymentEnabled
	default:
		return true
	}
}

func (s *notificationService) invalidateUnreadCache(ctx context.Context, userID int64) {
	cacheKey := fmt.Sprintf("notification:unread:%d", userID)
	s.redis.Del(ctx, cacheKey)
}

func (s *notificationService) sendPushNotification(userID int64, notification *Notification) {
}

func (s *notificationService) sendEmailNotification(userID int64, notification *Notification) {
}
