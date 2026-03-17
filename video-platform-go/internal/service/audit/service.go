package audit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditStatus string

const (
	AuditStatusPending   AuditStatus = "pending"
	AuditStatusApproved  AuditStatus = "approved"
	AuditStatusRejected  AuditStatus = "rejected"
	AuditStatusReviewing AuditStatus = "reviewing"
)

type ContentType string

const (
	ContentTypeVideo    ContentType = "video"
	ContentTypeComment  ContentType = "comment"
	ContentTypeDanmaku  ContentType = "danmaku"
	ContentTypeUser     ContentType = "user"
	ContentTypeLive     ContentType = "live"
)

type AuditResult struct {
	ID          int64        `json:"id"`
	ContentType ContentType  `json:"content_type"`
	ContentID   int64        `json:"content_id"`
	Status      AuditStatus  `json:"status"`
	Reason      string       `json:"reason,omitempty"`
	Labels      []string     `json:"labels,omitempty"`
	Score       float64      `json:"score"`
	ReviewerID  *int64       `json:"reviewer_id,omitempty"`
	ReviewedAt  *time.Time   `json:"reviewed_at,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type ContentLabel string

const (
	LabelViolence   ContentLabel = "violence"
	LabelPorn       ContentLabel = "porn"
	LabelTerror     ContentLabel = "terror"
	LabelPolitical   ContentLabel = "political"
	LabelSpam       ContentLabel = "spam"
	LabelFraud      ContentLabel = "fraud"
	LabelGambling   ContentLabel = "gambling"
	LabelDrugs      ContentLabel = "drugs"
	LabelSensitive  ContentLabel = "sensitive"
	LabelCopyright  ContentLabel = "copyright"
)

type AuditRule struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	Type        ContentType  `json:"type"`
	Keywords    []string     `json:"keywords"`
	Labels      []ContentLabel `json:"labels"`
	Severity    int          `json:"severity"`
	AutoAction  string       `json:"auto_action"`
	IsActive    bool         `json:"is_active"`
	CreatedAt   time.Time    `json:"created_at"`
}

type AuditRepository interface {
	CreateAuditResult(ctx context.Context, result *AuditResult) error
	GetAuditResult(ctx context.Context, contentType ContentType, contentID int64) (*AuditResult, error)
	GetAuditResultByID(ctx context.Context, id int64) (*AuditResult, error)
	GetPendingAudits(ctx context.Context, contentType ContentType, page, pageSize int) ([]*AuditResult, int64, error)
	UpdateAuditStatus(ctx context.Context, id int64, status AuditStatus, reason string, reviewerID int64) error
	
	CreateRule(ctx context.Context, rule *AuditRule) error
	GetRules(ctx context.Context, contentType ContentType) ([]*AuditRule, error)
	UpdateRule(ctx context.Context, rule *AuditRule) error
	DeleteRule(ctx context.Context, id int64) error
}

type auditRepository struct {
	pool *pgxpool.Pool
}

func NewAuditRepository(pool *pgxpool.Pool) AuditRepository {
	return &auditRepository{pool: pool}
}

func (r *auditRepository) CreateAuditResult(ctx context.Context, result *AuditResult) error {
	labelsJSON, _ := json.Marshal(result.Labels)
	query := `
		INSERT INTO audit_results (content_type, content_id, status, reason, labels, score, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query,
		result.ContentType, result.ContentID, result.Status, result.Reason,
		labelsJSON, result.Score,
	).Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt)
}

func (r *auditRepository) GetAuditResult(ctx context.Context, contentType ContentType, contentID int64) (*AuditResult, error) {
	query := `
		SELECT id, content_type, content_id, status, reason, labels, score, reviewer_id, reviewed_at, created_at, updated_at
		FROM audit_results WHERE content_type = $1 AND content_id = $2
		ORDER BY created_at DESC LIMIT 1
	`
	
	var result AuditResult
	var labelsJSON []byte
	err := r.pool.QueryRow(ctx, query, contentType, contentID).Scan(
		&result.ID, &result.ContentType, &result.ContentID, &result.Status, &result.Reason,
		&labelsJSON, &result.Score, &result.ReviewerID, &result.ReviewedAt,
		&result.CreatedAt, &result.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	
	json.Unmarshal(labelsJSON, &result.Labels)
	return &result, nil
}

func (r *auditRepository) GetAuditResultByID(ctx context.Context, id int64) (*AuditResult, error) {
	query := `
		SELECT id, content_type, content_id, status, reason, labels, score, reviewer_id, reviewed_at, created_at, updated_at
		FROM audit_results WHERE id = $1
	`
	
	var result AuditResult
	var labelsJSON []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&result.ID, &result.ContentType, &result.ContentID, &result.Status, &result.Reason,
		&labelsJSON, &result.Score, &result.ReviewerID, &result.ReviewedAt,
		&result.CreatedAt, &result.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	
	json.Unmarshal(labelsJSON, &result.Labels)
	return &result, nil
}

func (r *auditRepository) GetPendingAudits(ctx context.Context, contentType ContentType, page, pageSize int) ([]*AuditResult, int64, error) {
	whereClause := "WHERE status = 'pending'"
	args := []interface{}{}
	argNum := 1
	
	if contentType != "" {
		whereClause += fmt.Sprintf(" AND content_type = $%d", argNum)
		args = append(args, contentType)
		argNum++
	}
	
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM audit_results %s", whereClause)
	var total int64
	r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT id, content_type, content_id, status, reason, labels, score, reviewer_id, reviewed_at, created_at, updated_at
		FROM audit_results %s
		ORDER BY created_at ASC
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)
	
	args = append(args, pageSize, offset)
	
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var results []*AuditResult
	for rows.Next() {
		var result AuditResult
		var labelsJSON []byte
		if err := rows.Scan(
			&result.ID, &result.ContentType, &result.ContentID, &result.Status, &result.Reason,
			&labelsJSON, &result.Score, &result.ReviewerID, &result.ReviewedAt,
			&result.CreatedAt, &result.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		json.Unmarshal(labelsJSON, &result.Labels)
		results = append(results, &result)
	}
	
	return results, total, nil
}

func (r *auditRepository) UpdateAuditStatus(ctx context.Context, id int64, status AuditStatus, reason string, reviewerID int64) error {
	now := time.Now()
	query := `
		UPDATE audit_results SET status = $2, reason = $3, reviewer_id = $4, reviewed_at = $5, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, id, status, reason, reviewerID, now)
	return err
}

func (r *auditRepository) CreateRule(ctx context.Context, rule *AuditRule) error {
	keywordsJSON, _ := json.Marshal(rule.Keywords)
	labelsJSON, _ := json.Marshal(rule.Labels)
	query := `
		INSERT INTO audit_rules (name, type, keywords, labels, severity, auto_action, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		RETURNING id, created_at
	`
	return r.pool.QueryRow(ctx, query,
		rule.Name, rule.Type, keywordsJSON, labelsJSON, rule.Severity, rule.AutoAction, rule.IsActive,
	).Scan(&rule.ID, &rule.CreatedAt)
}

func (r *auditRepository) GetRules(ctx context.Context, contentType ContentType) ([]*AuditRule, error) {
	query := `SELECT id, name, type, keywords, labels, severity, auto_action, is_active, created_at FROM audit_rules WHERE is_active = true`
	args := []interface{}{}
	
	if contentType != "" {
		query += " AND type = $1"
		args = append(args, contentType)
	}
	
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var rules []*AuditRule
	for rows.Next() {
		var rule AuditRule
		var keywordsJSON, labelsJSON []byte
		if err := rows.Scan(
			&rule.ID, &rule.Name, &rule.Type, &keywordsJSON, &labelsJSON,
			&rule.Severity, &rule.AutoAction, &rule.IsActive, &rule.CreatedAt,
		); err != nil {
			return nil, err
		}
		json.Unmarshal(keywordsJSON, &rule.Keywords)
		json.Unmarshal(labelsJSON, &rule.Labels)
		rules = append(rules, &rule)
	}
	
	return rules, nil
}

func (r *auditRepository) UpdateRule(ctx context.Context, rule *AuditRule) error {
	keywordsJSON, _ := json.Marshal(rule.Keywords)
	labelsJSON, _ := json.Marshal(rule.Labels)
	query := `
		UPDATE audit_rules SET name = $2, type = $3, keywords = $4, labels = $5, severity = $6, auto_action = $7, is_active = $8
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query,
		rule.ID, rule.Name, rule.Type, keywordsJSON, labelsJSON,
		rule.Severity, rule.AutoAction, rule.IsActive,
	)
	return err
}

func (r *auditRepository) DeleteRule(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM audit_rules WHERE id = $1`, id)
	return err
}

type ContentProvider interface {
	GetContent(ctx context.Context, contentID int64) (string, error)
	UpdateStatus(ctx context.Context, contentID int64, status string) error
}

type AuditService interface {
	SubmitForAudit(ctx context.Context, contentType ContentType, contentID int64, content string) (*AuditResult, error)
	GetAuditResult(ctx context.Context, contentType ContentType, contentID int64) (*AuditResult, error)
	GetPendingAudits(ctx context.Context, contentType ContentType, page, pageSize int) ([]*AuditResult, int64, error)
	Approve(ctx context.Context, id int64, reviewerID int64) error
	Reject(ctx context.Context, id int64, reviewerID int64, reason string) error
	CreateRule(ctx context.Context, rule *AuditRule) error
	GetRules(ctx context.Context, contentType ContentType) ([]*AuditRule, error)
}

type auditService struct {
	repo            AuditRepository
	contentProviders map[ContentType]ContentProvider
}

func NewAuditService(repo AuditRepository, providers map[ContentType]ContentProvider) AuditService {
	return &auditService{repo: repo, contentProviders: providers}
}

func (s *auditService) SubmitForAudit(ctx context.Context, contentType ContentType, contentID int64, content string) (*AuditResult, error) {
	rules, err := s.repo.GetRules(ctx, contentType)
	if err != nil {
		return nil, err
	}
	
	var labels []string
	var maxSeverity int
	var score float64
	
	for _, rule := range rules {
		if s.matchesRule(content, rule) {
			for _, label := range rule.Labels {
				labels = append(labels, string(label))
			}
			if rule.Severity > maxSeverity {
				maxSeverity = rule.Severity
			}
		}
	}
	
	score = s.calculateScore(labels, maxSeverity)
	
	status := AuditStatusPending
	if maxSeverity >= 8 {
		status = AuditStatusRejected
	} else if len(labels) == 0 {
		status = AuditStatusApproved
	}
	
	result := &AuditResult{
		ContentType: contentType,
		ContentID:   contentID,
		Status:      status,
		Labels:      labels,
		Score:       score,
	}
	
	if err := s.repo.CreateAuditResult(ctx, result); err != nil {
		return nil, err
	}
	
	if status == AuditStatusApproved {
		if provider, ok := s.contentProviders[contentType]; ok {
			provider.UpdateStatus(ctx, contentID, "approved")
		}
	} else if status == AuditStatusRejected {
		if provider, ok := s.contentProviders[contentType]; ok {
			provider.UpdateStatus(ctx, contentID, "rejected")
		}
	}
	
	return result, nil
}

func (s *auditService) GetAuditResult(ctx context.Context, contentType ContentType, contentID int64) (*AuditResult, error) {
	return s.repo.GetAuditResult(ctx, contentType, contentID)
}

func (s *auditService) GetPendingAudits(ctx context.Context, contentType ContentType, page, pageSize int) ([]*AuditResult, int64, error) {
	return s.repo.GetPendingAudits(ctx, contentType, page, pageSize)
}

func (s *auditService) Approve(ctx context.Context, id int64, reviewerID int64) error {
	result, err := s.repo.GetAuditResultByID(ctx, id)
	if err != nil {
		return err
	}
	
	if result == nil {
		return errors.New("audit result not found")
	}
	
	if err := s.repo.UpdateAuditStatus(ctx, id, AuditStatusApproved, "", reviewerID); err != nil {
		return err
	}
	
	if provider, ok := s.contentProviders[result.ContentType]; ok {
		provider.UpdateStatus(ctx, result.ContentID, "approved")
	}
	
	return nil
}

func (s *auditService) Reject(ctx context.Context, id int64, reviewerID int64, reason string) error {
	result, err := s.repo.GetAuditResultByID(ctx, id)
	if err != nil {
		return err
	}
	
	if result == nil {
		return errors.New("audit result not found")
	}
	
	if err := s.repo.UpdateAuditStatus(ctx, id, AuditStatusRejected, reason, reviewerID); err != nil {
		return err
	}
	
	if provider, ok := s.contentProviders[result.ContentType]; ok {
		provider.UpdateStatus(ctx, result.ContentID, "rejected")
	}
	
	return nil
}

func (s *auditService) CreateRule(ctx context.Context, rule *AuditRule) error {
	return s.repo.CreateRule(ctx, rule)
}

func (s *auditService) GetRules(ctx context.Context, contentType ContentType) ([]*AuditRule, error) {
	return s.repo.GetRules(ctx, contentType)
}

func (s *auditService) matchesRule(content string, rule *AuditRule) bool {
	for _, keyword := range rule.Keywords {
		if containsIgnoreCase(content, keyword) {
			return true
		}
	}
	return false
}

func (s *auditService) calculateScore(labels []string, severity int) float64 {
	baseScore := float64(severity) * 10
	labelPenalty := float64(len(labels)) * 5
	return baseScore + labelPenalty
}

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr)
}
