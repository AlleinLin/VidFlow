package com.videoplatform.controller;

import com.videoplatform.domain.entity.WatchHistory;
import com.videoplatform.service.PlaybackService;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/playback")
@RequiredArgsConstructor
public class PlaybackController {

    private final PlaybackService playbackService;

    @PostMapping("/progress")
    public ResponseEntity<?> updateProgress(
            @AuthenticationPrincipal UserDetails userDetails,
            @RequestBody ProgressRequest request) {
        Long userId = Long.parseLong(userDetails.getUsername());
        playbackService.updateProgress(
            userId, 
            request.getVideoId(), 
            request.getPosition(), 
            request.getDuration(), 
            request.getWatchDuration()
        );
        return ResponseEntity.ok(Map.of("message", "Progress updated successfully"));
    }

    @GetMapping("/progress/{videoId}")
    public ResponseEntity<?> getProgress(
            @AuthenticationPrincipal UserDetails userDetails,
            @PathVariable Long videoId) {
        Long userId = Long.parseLong(userDetails.getUsername());
        return ResponseEntity.ok(playbackService.getProgress(userId, videoId));
    }

    @GetMapping("/history")
    public ResponseEntity<?> getWatchHistory(
            @AuthenticationPrincipal UserDetails userDetails,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size) {
        Long userId = Long.parseLong(userDetails.getUsername());
        Page<WatchHistory> history = playbackService.getWatchHistory(userId, page, size);
        return ResponseEntity.ok(Map.of(
            "histories", history.getContent(),
            "total", history.getTotalElements(),
            "page", page,
            "page_size", size
        ));
    }

    @GetMapping("/continue-watching")
    public ResponseEntity<?> getContinueWatching(
            @AuthenticationPrincipal UserDetails userDetails,
            @RequestParam(defaultValue = "10") int limit) {
        Long userId = Long.parseLong(userDetails.getUsername());
        List<WatchHistory> history = playbackService.getContinueWatching(userId, limit);
        return ResponseEntity.ok(history);
    }

    @DeleteMapping("/history/{videoId}")
    public ResponseEntity<?> deleteWatchHistory(
            @AuthenticationPrincipal UserDetails userDetails,
            @PathVariable Long videoId) {
        Long userId = Long.parseLong(userDetails.getUsername());
        playbackService.deleteWatchHistory(userId, videoId);
        return ResponseEntity.noContent().build();
    }

    @DeleteMapping("/history")
    public ResponseEntity<?> clearWatchHistory(@AuthenticationPrincipal UserDetails userDetails) {
        Long userId = Long.parseLong(userDetails.getUsername());
        playbackService.clearWatchHistory(userId);
        return ResponseEntity.noContent().build();
    }

    public static class ProgressRequest {
        private Long videoId;
        private Double position;
        private Double duration;
        private Long watchDuration;

        public Long getVideoId() { return videoId; }
        public void setVideoId(Long videoId) { this.videoId = videoId; }
        public Double getPosition() { return position; }
        public void setPosition(Double position) { this.position = position; }
        public Double getDuration() { return duration; }
        public void setDuration(Double duration) { this.duration = duration; }
        public Long getWatchDuration() { return watchDuration; }
        public void setWatchDuration(Long watchDuration) { this.watchDuration = watchDuration; }
    }
}
