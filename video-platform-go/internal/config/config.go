package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Scylla   ScyllaConfig
	Kafka    KafkaConfig
	JWT      JWTConfig
	Storage  StorageConfig
	Log      LogConfig
	Metrics  MetricsConfig
	Tracing  TracingConfig
}

type ServerConfig struct {
	Name         string        `mapstructure:"name"`
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	Mode         string        `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Driver          string        `mapstructure:"driver"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode,
	)
}

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Addr         string `mapstructure:"addr"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

func (c *RedisConfig) Address() string {
	if c.Addr != "" {
		return c.Addr
	}
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type ScyllaConfig struct {
	Hosts    []string `mapstructure:"hosts"`
	Keyspace string   `mapstructure:"keyspace"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
	GroupID string   `mapstructure:"group_id"`
	Async   bool     `mapstructure:"async"`
	Topics  Topics   `mapstructure:"topics"`
}

type Topics struct {
	VideoUpload        string `mapstructure:"video_upload"`
	VideoTranscode     string `mapstructure:"video_transcode"`
	VideoAudit         string `mapstructure:"video_audit"`
	UserAction         string `mapstructure:"user_action"`
	Notification       string `mapstructure:"notification"`
	Interaction        string `mapstructure:"interaction"`
}

type JWTConfig struct {
	Secret          string        `mapstructure:"secret"`
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
	Issuer          string        `mapstructure:"issuer"`
}

type StorageConfig struct {
	Provider string `mapstructure:"provider"`
	Endpoint string `mapstructure:"endpoint"`
	Region   string `mapstructure:"region"`
	Bucket   string `mapstructure:"bucket"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	UseSSL   bool   `mapstructure:"use_ssl"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
	Port    int    `mapstructure:"port"`
}

type TracingConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Endpoint string `mapstructure:"endpoint"`
	Sampler  float64 `mapstructure:"sampler"`
}

func Load(configPath string) (*Config, error) {
	v := viper.New()
	
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	
	v.AutomaticEnv()
	v.SetEnvPrefix("VIDEO")
	
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	return &cfg, nil
}

func LoadFromEnv() (*Config, error) {
	v := viper.New()
	
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", "30s")
	v.SetDefault("server.write_timeout", "30s")
	v.SetDefault("server.idle_timeout", "60s")
	v.SetDefault("server.mode", "debug")
	
	v.SetDefault("database.driver", "postgres")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.max_open_conns", 100)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime", "1h")
	v.SetDefault("database.conn_max_idle_time", "10m")
	
	v.SetDefault("redis.addr", "localhost:6379")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.pool_size", 100)
	v.SetDefault("redis.min_idle_conns", 10)
	
	v.SetDefault("jwt.access_token_ttl", "15m")
	v.SetDefault("jwt.refresh_token_ttl", "168h")
	v.SetDefault("jwt.issuer", "video-platform")
	
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.output", "stdout")
	
	v.SetDefault("metrics.enabled", true)
	v.SetDefault("metrics.path", "/metrics")
	v.SetDefault("metrics.port", 9090)
	
	v.SetDefault("tracing.enabled", false)
	v.SetDefault("tracing.sampler", 0.1)
	
	v.AutomaticEnv()
	v.SetEnvPrefix("VIDEO")
	
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	return &cfg, nil
}
