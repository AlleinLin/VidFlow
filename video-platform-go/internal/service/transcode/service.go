package transcode

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/video-platform/go/internal/domain/video"
	"github.com/video-platform/go/pkg/logger"
)

type TranscodeConfig struct {
	OutputDir       string
	FFmpegPath      string
	FFprobePath     string
	MaxWorkers      int
	SegmentDuration int
	Resolutions     []string
}

type TranscodeTask struct {
	VideoID    int64
	InputPath  string
	OutputPath string
	Resolution string
	Format     string
	Status     string
	Error      error
	Done       chan struct{}
}

type TranscodeResult struct {
	VideoID         int64
	Resolution      string
	OutputPath      string
	DurationSeconds int
	Error           error
}

type VideoInfo struct {
	Duration      float64
	Width         int
	Height        int
	Codec         string
	BitRate       int
	FrameRate     float64
	HasAudio      bool
	AudioCodec    string
	AudioBitRate  int
}

type Service interface {
	Transcode(ctx context.Context, videoID int64, inputPath string) ([]*TranscodeResult, error)
	GetVideoInfo(ctx context.Context, inputPath string) (*VideoInfo, error)
	GenerateThumbnail(ctx context.Context, inputPath, outputPath string, timeOffset float64) error
}

type service struct {
	config TranscodeConfig
	taskCh chan *TranscodeTask
}

func NewService(config TranscodeConfig) Service {
	if config.MaxWorkers <= 0 {
		config.MaxWorkers = 4
	}
	if config.SegmentDuration <= 0 {
		config.SegmentDuration = 6
	}
	if config.FFmpegPath == "" {
		config.FFmpegPath = "ffmpeg"
	}
	if config.FFprobePath == "" {
		config.FFprobePath = "ffprobe"
	}
	if len(config.Resolutions) == 0 {
		config.Resolutions = []string{"240p", "480p", "720p", "1080p"}
	}

	s := &service{
		config: config,
		taskCh: make(chan *TranscodeTask, config.MaxWorkers*2),
	}

	for i := 0; i < config.MaxWorkers; i++ {
		go s.worker()
	}

	return s
}

func (s *service) worker() {
	for task := range s.taskCh {
		s.processTask(task)
	}
}

func (s *service) processTask(task *TranscodeTask) {
	defer close(task.Done)

	ctx := context.Background()
	logger.Info(ctx, "Starting transcode task", "video_id", task.VideoID, "resolution", task.Resolution)

	result := s.transcodeResolution(ctx, task.InputPath, task.OutputPath, task.Resolution)

	if result.Error != nil {
		task.Error = result.Error
		task.Status = "FAILED"
		logger.Error(ctx, "Transcode task failed", "video_id", task.VideoID, "resolution", task.Resolution, "error", result.Error)
	} else {
		task.Status = "COMPLETED"
		logger.Info(ctx, "Transcode task completed", "video_id", task.VideoID, "resolution", task.Resolution)
	}
}

func (s *service) Transcode(ctx context.Context, videoID int64, inputPath string) ([]*TranscodeResult, error) {
	videoInfo, err := s.GetVideoInfo(ctx, inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get video info: %w", err)
	}

	outputBase := filepath.Join(s.config.OutputDir, strconv.FormatInt(videoID, 10))
	if err := os.MkdirAll(outputBase, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	masterPlaylistPath := filepath.Join(outputBase, "master.m3u8")
	masterPlaylist := "#EXTM3U\n#EXT-X-VERSION:3\n"

	var tasks []*TranscodeTask
	var results []*TranscodeResult

	for _, res := range s.config.Resolutions {
		resConfig, ok := video.ResolutionConfigs[res]
		if !ok {
			continue
		}

		if videoInfo.Height < resConfig.Height {
			continue
		}

		outputDir := filepath.Join(outputBase, res)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create resolution directory: %w", err)
		}

		task := &TranscodeTask{
			VideoID:    videoID,
			InputPath:  inputPath,
			OutputPath: outputDir,
			Resolution: res,
			Format:     "HLS",
			Status:     "PENDING",
			Done:       make(chan struct{}),
		}

		tasks = append(tasks, task)
		s.taskCh <- task

		bandwidth := resConfig.BitrateH264 * 1000
		masterPlaylist += fmt.Sprintf(
			"#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%dx%d\n%s/playlist.m3u8\n",
			bandwidth, resConfig.Width, resConfig.Height, res,
		)
	}

	for _, task := range tasks {
		<-task.Done
		result := &TranscodeResult{
			VideoID:    task.VideoID,
			Resolution: task.Resolution,
			OutputPath: task.OutputPath,
			Error:      task.Error,
		}
		if task.Error == nil {
			result.DurationSeconds = int(videoInfo.Duration)
		}
		results = append(results, result)
	}

	for _, result := range results {
		if result.Error != nil {
			return results, fmt.Errorf("one or more transcode tasks failed")
		}
	}

	if err := os.WriteFile(masterPlaylistPath, []byte(masterPlaylist), 0644); err != nil {
		return nil, fmt.Errorf("failed to write master playlist: %w", err)
	}

	return results, nil
}

func (s *service) transcodeResolution(ctx context.Context, inputPath, outputPath, resolution string) *TranscodeResult {
	resConfig, ok := video.ResolutionConfigs[resolution]
	if !ok {
		return &TranscodeResult{Error: fmt.Errorf("unknown resolution: %s", resolution)}
	}

	playlistPath := filepath.Join(outputPath, "playlist.m3u8")
	segmentPattern := filepath.Join(outputPath, "segment_%04d.ts")

	args := []string{
		"-i", inputPath,
		"-c:v", "libx264",
		"-preset", "fast",
		"-b:v", fmt.Sprintf("%dk", resConfig.BitrateH264),
		"-s", fmt.Sprintf("%dx%d", resConfig.Width, resConfig.Height),
		"-c:a", "aac",
		"-b:a", "128k",
		"-hls_time", strconv.Itoa(s.config.SegmentDuration),
		"-hls_list_size", "0",
		"-hls_segment_filename", segmentPattern,
		"-f", "hls",
		playlistPath,
	}

	cmd := exec.CommandContext(ctx, s.config.FFmpegPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &TranscodeResult{
			Error: fmt.Errorf("ffmpeg failed: %w, output: %s", err, string(output)),
		}
	}

	return &TranscodeResult{OutputPath: outputPath}
}

func (s *service) GetVideoInfo(ctx context.Context, inputPath string) (*VideoInfo, error) {
	args := []string{
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		inputPath,
	}

	cmd := exec.CommandContext(ctx, s.config.FFprobePath, args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe failed: %w", err)
	}

	return parseVideoInfo(output)
}

func (s *service) GenerateThumbnail(ctx context.Context, inputPath, outputPath string, timeOffset float64) error {
	args := []string{
		"-ss", strconv.FormatFloat(timeOffset, 'f', 2, 64),
		"-i", inputPath,
		"-vframes", "1",
		"-q:v", "2",
		"-y",
		outputPath,
	}

	cmd := exec.CommandContext(ctx, s.config.FFmpegPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg thumbnail failed: %w, output: %s", err, string(output))
	}

	return nil
}

type probeData struct {
	Format struct {
		Duration string `json:"duration"`
		BitRate  string `json:"bit_rate"`
	} `json:"format"`
	Streams []struct {
		CodecName   string `json:"codec_name"`
		CodecType   string `json:"codec_type"`
		Width       int    `json:"width"`
		Height      int    `json:"height"`
		BitRate     string `json:"bit_rate"`
		RFrameRate  string `json:"r_frame_rate"`
	} `json:"streams"`
}

func parseVideoInfo(data []byte) (*VideoInfo, error) {
	var pd probeData
	if err := json.Unmarshal(data, &pd); err != nil {
		return nil, fmt.Errorf("failed to parse ffprobe output: %w", err)
	}

	info := &VideoInfo{}

	if pd.Format.Duration != "" {
		info.Duration, _ = strconv.ParseFloat(pd.Format.Duration, 64)
	}
	if pd.Format.BitRate != "" {
		info.BitRate, _ = strconv.Atoi(pd.Format.BitRate)
	}

	for _, stream := range pd.Streams {
		if stream.CodecType == "video" && info.Width == 0 {
			info.Width = stream.Width
			info.Height = stream.Height
			info.Codec = stream.CodecName
			if stream.BitRate != "" {
				info.BitRate, _ = strconv.Atoi(stream.BitRate)
			}
			if stream.RFrameRate != "" {
				parts := strings.Split(stream.RFrameRate, "/")
				if len(parts) == 2 {
					num, _ := strconv.ParseFloat(parts[0], 64)
					den, _ := strconv.ParseFloat(parts[1], 64)
					if den != 0 {
						info.FrameRate = num / den
					}
				}
			}
		}
		if stream.CodecType == "audio" {
			info.HasAudio = true
			info.AudioCodec = stream.CodecName
			if stream.BitRate != "" {
				info.AudioBitRate, _ = strconv.Atoi(stream.BitRate)
			}
		}
	}

	return info, nil
}
