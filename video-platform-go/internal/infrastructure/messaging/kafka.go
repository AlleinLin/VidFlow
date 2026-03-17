package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/video-platform/go/internal/config"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(cfg *config.KafkaConfig) (*KafkaProducer, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    cfg.Topic,
		Balancer: &kafka.LeastBytes{},
		Async:    cfg.Async,
	}

	return &KafkaProducer{writer: writer}, nil
}

func (p *KafkaProducer) SendMessage(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: data,
	})
}

func (p *KafkaProducer) SendMessages(ctx context.Context, messages []kafka.Message) error {
	return p.writer.WriteMessages(ctx, messages...)
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(cfg *config.KafkaConfig, groupID string) (*KafkaConsumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		Topic:    cfg.Topic,
		GroupID:  groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	return &KafkaConsumer{reader: reader}, nil
}

func (c *KafkaConsumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.ReadMessage(ctx)
}

func (c *KafkaConsumer) FetchMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.FetchMessage(ctx)
}

func (c *KafkaConsumer) CommitMessages(ctx context.Context, msgs ...kafka.Message) error {
	return c.reader.CommitMessages(ctx, msgs...)
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}

type MessageHandler func(ctx context.Context, msg kafka.Message) error

func (c *KafkaConsumer) Consume(ctx context.Context, handler MessageHandler) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := c.FetchMessage(ctx)
			if err != nil {
				return fmt.Errorf("failed to fetch message: %w", err)
			}

			if err := handler(ctx, msg); err != nil {
				return fmt.Errorf("failed to handle message: %w", err)
			}

			if err := c.CommitMessages(ctx, msg); err != nil {
				return fmt.Errorf("failed to commit message: %w", err)
			}
		}
	}
}

type Event struct {
	Type      string          `json:"type"`
	Timestamp int64           `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

const (
	EventTypeVideoUploaded    = "video.uploaded"
	EventTypeVideoTranscoded  = "video.transcoded"
	EventTypeVideoPublished   = "video.published"
	EventTypeVideoDeleted     = "video.deleted"
	EventTypeUserRegistered   = "user.registered"
	EventTypeUserFollowed     = "user.followed"
	EventTypeCommentCreated   = "comment.created"
	EventTypeLikeCreated      = "like.created"
	EventTypeDanmakuCreated   = "danmaku.created"
	EventTypeWatchProgress    = "watch.progress"
	EventTypeRecommendRequest = "recommend.request"
)

type VideoUploadedEvent struct {
	VideoID     int64  `json:"video_id"`
	UserID      int64  `json:"user_id"`
	Title       string `json:"title"`
	Filename    string `json:"filename"`
	StorageKey  string `json:"storage_key"`
}

type VideoTranscodedEvent struct {
	VideoID     int64  `json:"video_id"`
	Resolution  string `json:"resolution"`
	Format      string `json:"format"`
	StorageKey  string `json:"storage_key"`
}

type UserFollowedEvent struct {
	FollowerID  int64 `json:"follower_id"`
	FollowingID int64 `json:"following_id"`
}

type CommentCreatedEvent struct {
	CommentID int64 `json:"comment_id"`
	VideoID   int64 `json:"video_id"`
	UserID    int64 `json:"user_id"`
	Content   string `json:"content"`
}

type WatchProgressEvent struct {
	UserID    int64   `json:"user_id"`
	VideoID   int64   `json:"video_id"`
	Progress  float64 `json:"progress"`
	Timestamp int64   `json:"timestamp"`
}
