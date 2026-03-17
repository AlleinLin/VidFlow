package cdn

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CDNProvider string

const (
	CDNProviderAliyun   CDNProvider = "aliyun"
	CDNProviderTencent  CDNProvider = "tencent"
	CDNProviderQiniu    CDNProvider = "qiniu"
	CDNProviderAWS      CDNProvider = "aws"
	CDNProviderCustom   CDNProvider = "custom"
)

type CDNConfig struct {
	Provider    CDNProvider
	Domain      string
	AccessKey   string
	SecretKey   string
	Bucket      string
	Region      string
	Endpoint    string
	EnableHTTPS bool
}

type CDNNode struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Region   string   `json:"region"`
	Endpoint string   `json:"endpoint"`
	Weight   int      `json:"weight"`
	Healthy  bool     `json:"healthy"`
}

type CacheRule struct {
	ID          int64    `json:"id"`
	PathPattern string   `json:"path_pattern"`
	TTL         int      `json:"ttl"`
	CacheKey    string   `json:"cache_key"`
	EnableGzip  bool     `json:"enable_gzip"`
	EnableSSL   bool     `json:"enable_ssl"`
}

type URLSigner struct {
	SecretKey string
	ExpireTime int
}

type CDNService interface {
	GetSignedURL(ctx context.Context, objectKey string, expire time.Duration) (string, error)
	GetVideoURL(ctx context.Context, videoID int64, resolution string) (string, error)
	GetThumbnailURL(ctx context.Context, videoID int64) (string, error)
	RefreshCache(ctx context.Context, urls []string) error
	Prefetch(ctx context.Context, urls []string) error
	GetNodeList(ctx context.Context) ([]*CDNNode, error)
	GetStatistics(ctx context.Context, startTime, endTime time.Time) (*CDNStatistics, error)
}

type CDNStatistics struct {
	Bandwidth      int64   `json:"bandwidth"`
	Traffic        int64   `json:"traffic"`
	RequestCount   int64   `json:"request_count"`
	HitRate        float64 `json:"hit_rate"`
	StatusCodeDist map[int]int64 `json:"status_code_dist"`
	TopURLs        []string `json:"top_urls"`
}

type CDNRepository interface {
	GetCDNConfig(ctx context.Context, provider CDNProvider) (*CDNConfig, error)
	GetCacheRules(ctx context.Context) ([]*CacheRule, error)
	SaveStatistics(ctx context.Context, stats *CDNStatistics) error
	GetStatistics(ctx context.Context, startTime, endTime time.Time) ([]*CDNStatistics, error)
}

type cdnRepository struct {
	pool *pgxpool.Pool
}

func NewCDNRepository(pool *pgxpool.Pool) CDNRepository {
	return &cdnRepository{pool: pool}
}

func (r *cdnRepository) GetCDNConfig(ctx context.Context, provider CDNProvider) (*CDNConfig, error) {
	return &CDNConfig{
		Provider:    provider,
		Domain:      "cdn.example.com",
		EnableHTTPS: true,
	}, nil
}

func (r *cdnRepository) GetCacheRules(ctx context.Context) ([]*CacheRule, error) {
	return []*CacheRule{
		{PathPattern: "/video/*", TTL: 86400, EnableGzip: true, EnableSSL: true},
		{PathPattern: "/thumbnail/*", TTL: 604800, EnableGzip: true, EnableSSL: true},
	}, nil
}

func (r *cdnRepository) SaveStatistics(ctx context.Context, stats *CDNStatistics) error {
	return nil
}

func (r *cdnRepository) GetStatistics(ctx context.Context, startTime, endTime time.Time) ([]*CDNStatistics, error) {
	return []*CDNStatistics{}, nil
}

type cdnService struct {
	repo   CDNRepository
	config *CDNConfig
	signer *URLSigner
	nodes  []*CDNNode
}

func NewCDNService(repo CDNRepository, config *CDNConfig) CDNService {
	return &cdnService{
		repo:   repo,
		config: config,
		signer: &URLSigner{
			SecretKey:  config.SecretKey,
			ExpireTime: 3600,
		},
		nodes: []*CDNNode{
			{ID: "node1", Name: "Beijing", Region: "cn-north", Endpoint: "cdn1.example.com", Weight: 100, Healthy: true},
			{ID: "node2", Name: "Shanghai", Region: "cn-east", Endpoint: "cdn2.example.com", Weight: 100, Healthy: true},
			{ID: "node3", Name: "Guangzhou", Region: "cn-south", Endpoint: "cdn3.example.com", Weight: 80, Healthy: true},
		},
	}
}

func (s *cdnService) GetSignedURL(ctx context.Context, objectKey string, expire time.Duration) (string, error) {
	expireTime := time.Now().Add(expire).Unix()
	
	path := "/" + objectKey
	stringToSign := fmt.Sprintf("%s%d", path, expireTime)
	
	h := sha256.New()
	h.Write([]byte(stringToSign))
	h.Write([]byte(s.signer.SecretKey))
	signature := hex.EncodeToString(h.Sum(nil))[:32]
	
	scheme := "http"
	if s.config.EnableHTTPS {
		scheme = "https"
	}
	
	signedURL := fmt.Sprintf("%s://%s%s?sign=%s&t=%d",
		scheme, s.config.Domain, path, signature, expireTime)
	
	return signedURL, nil
}

func (s *cdnService) GetVideoURL(ctx context.Context, videoID int64, resolution string) (string, error) {
	objectKey := fmt.Sprintf("video/%d/%s/playlist.m3u8", videoID, resolution)
	return s.GetSignedURL(ctx, objectKey, 2*time.Hour)
}

func (s *cdnService) GetThumbnailURL(ctx context.Context, videoID int64) (string, error) {
	objectKey := fmt.Sprintf("thumbnail/%d.jpg", videoID)
	return s.GetSignedURL(ctx, objectKey, 24*time.Hour)
}

func (s *cdnService) RefreshCache(ctx context.Context, urls []string) error {
	return nil
}

func (s *cdnService) Prefetch(ctx context.Context, urls []string) error {
	return nil
}

func (s *cdnService) GetNodeList(ctx context.Context) ([]*CDNNode, error) {
	return s.nodes, nil
}

func (s *cdnService) GetStatistics(ctx context.Context, startTime, endTime time.Time) (*CDNStatistics, error) {
	return &CDNStatistics{
		Bandwidth:    1024 * 1024 * 1024,
		Traffic:      10 * 1024 * 1024 * 1024,
		RequestCount: 1000000,
		HitRate:      0.95,
		StatusCodeDist: map[int]int64{
			200: 950000,
			304: 30000,
			404: 10000,
			500: 10000,
		},
	}, nil
}

type StorageService interface {
	Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error
	Download(ctx context.Context, key string) (io.ReadCloser, int64, error)
	Delete(ctx context.Context, key string) error
	GetSignedURL(ctx context.Context, key string, expire time.Duration) (string, error)
	Exists(ctx context.Context, key string) (bool, error)
	GetMetadata(ctx context.Context, key string) (map[string]string, error)
}

type storageService struct {
	config *CDNConfig
	client *http.Client
}

func NewStorageService(config *CDNConfig) StorageService {
	return &storageService{
		config: config,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *storageService) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	return nil
}

func (s *storageService) Download(ctx context.Context, key string) (io.ReadCloser, int64, error) {
	return nil, 0, errors.New("not implemented")
}

func (s *storageService) Delete(ctx context.Context, key string) error {
	return nil
}

func (s *storageService) GetSignedURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	expireTime := time.Now().Add(expire).Unix()
	
	objectPath := "/" + s.config.Bucket + "/" + key
	stringToSign := fmt.Sprintf("GET\n\n\n%d\n%s", expireTime, objectPath)
	
	h := sha256.New()
	h.Write([]byte(stringToSign))
	h.Write([]byte(s.config.SecretKey))
	signature := hex.EncodeToString(h.Sum(nil))
	
	scheme := "https"
	if !s.config.EnableHTTPS {
		scheme = "http"
	}
	
	signedURL := fmt.Sprintf("%s://%s/%s?OSSAccessKeyId=%s&Expires=%d&Signature=%s",
		scheme, s.config.Endpoint, key, s.config.AccessKey, expireTime, signature)
	
	return signedURL, nil
}

func (s *storageService) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func (s *storageService) GetMetadata(ctx context.Context, key string) (map[string]string, error) {
	return make(map[string]string), nil
}

type VideoStorageService struct {
	storage StorageService
	cdn     CDNService
}

func NewVideoStorageService(storage StorageService, cdn CDNService) *VideoStorageService {
	return &VideoStorageService{
		storage: storage,
		cdn:     cdn,
	}
}

func (s *VideoStorageService) UploadVideo(ctx context.Context, videoID int64, resolution string, reader io.Reader, size int64) error {
	key := fmt.Sprintf("video/%d/%s/video.mp4", videoID, resolution)
	return s.storage.Upload(ctx, key, reader, size, "video/mp4")
}

func (s *VideoStorageService) UploadHLS(ctx context.Context, videoID int64, resolution string, playlist []byte, segments map[string][]byte) error {
	playlistKey := fmt.Sprintf("video/%d/%s/playlist.m3u8", videoID, resolution)
	if err := s.storage.Upload(ctx, playlistKey, strings.NewReader(string(playlist)), int64(len(playlist)), "application/vnd.apple.mpegurl"); err != nil {
		return err
	}
	
	for segName, segData := range segments {
		segKey := fmt.Sprintf("video/%d/%s/%s", videoID, resolution, segName)
		if err := s.storage.Upload(ctx, segKey, strings.NewReader(string(segData)), int64(len(segData)), "video/MP2T"); err != nil {
			return err
		}
	}
	
	return nil
}

func (s *VideoStorageService) GetVideoURL(ctx context.Context, videoID int64, resolution string) (string, error) {
	return s.cdn.GetVideoURL(ctx, videoID, resolution)
}

func (s *VideoStorageService) UploadThumbnail(ctx context.Context, videoID int64, reader io.Reader, size int64) error {
	key := fmt.Sprintf("thumbnail/%d.jpg", videoID)
	return s.storage.Upload(ctx, key, reader, size, "image/jpeg")
}

func (s *VideoStorageService) GetThumbnailURL(ctx context.Context, videoID int64) (string, error) {
	return s.cdn.GetThumbnailURL(ctx, videoID)
}

func (s *VideoStorageService) DeleteVideo(ctx context.Context, videoID int64) error {
	resolutions := []string{"240p", "480p", "720p", "1080p", "4k"}
	
	for _, res := range resolutions {
		key := fmt.Sprintf("video/%d/%s", videoID, res)
		s.storage.Delete(ctx, key)
	}
	
	s.storage.Delete(ctx, fmt.Sprintf("thumbnail/%d.jpg", videoID))
	
	return nil
}

func ParseURL(rawURL string) (*url.URL, error) {
	return url.Parse(rawURL)
}

func JoinURL(base string, parts ...string) string {
	u, err := url.Parse(base)
	if err != nil {
		return base
	}
	
	u.Path = path.Join(u.Path, path.Join(parts...))
	return u.String()
}

func ToJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
