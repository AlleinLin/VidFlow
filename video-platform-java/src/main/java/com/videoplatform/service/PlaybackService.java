package com.videoplatform.service;

import com.videoplatform.domain.entity.Video;
import com.videoplatform.domain.entity.WatchHistory;
import com.videoplatform.domain.repository.VideoRepository;
import com.videoplatform.domain.repository.WatchHistoryRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
public class PlaybackService {

    private final WatchHistoryRepository watchHistoryRepository;
    private final VideoRepository videoRepository;

    @Transactional
    public WatchHistory updateProgress(Long userId, Long videoId, Double position, Double duration, Long watchDuration) {
        WatchHistory history = watchHistoryRepository
            .findByUserIdAndVideoId(userId, videoId)
            .orElse(WatchHistory.builder()
                .user(new User())
                .video(new Video())
                .build());
        
        history.getUser().setId(userId);
        history.getVideo().setId(videoId);
        
        double progress = duration > 0 ? position / duration : 0;
        boolean completed = progress >= 0.95;
        
        history.setWatchDuration(history.getWatchDuration() + watchDuration);
        history.setWatchProgress(progress);
        history.setLastPosition(position);
        history.setCompleted(completed);
        history.setWatchedAt(LocalDateTime.now());
        
        if (progress < 0.1) {
            Video video = videoRepository.findById(videoId).orElse(null);
            if (video != null) {
                video.setViewCount(video.getViewCount() + 1);
                videoRepository.save(video);
            }
        }
        
        return watchHistoryRepository.save(history);
    }

    public Optional<WatchHistory> getProgress(Long userId, Long videoId) {
        return watchHistoryRepository.findByUserIdAndVideoId(userId, videoId);
    }

    public Page<WatchHistory> getWatchHistory(Long userId, int page, int size) {
        Pageable pageable = PageRequest.of(page, size);
        return watchHistoryRepository.findByUserIdOrderByWatchedAtDesc(userId, pageable);
    }

    public List<WatchHistory> getContinueWatching(Long userId, int limit) {
        Pageable pageable = PageRequest.of(0, limit);
        return watchHistoryRepository.findContinueWatching(userId, pageable);
    }

    @Transactional
    public void deleteWatchHistory(Long userId, Long videoId) {
        watchHistoryRepository.deleteByUserIdAndVideoId(userId, videoId);
    }

    @Transactional
    public void clearWatchHistory(Long userId) {
        watchHistoryRepository.deleteByUserId(userId);
    }
}

class User {
    private Long id;
    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }
}

class Video {
    private Long id;
    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }
}
