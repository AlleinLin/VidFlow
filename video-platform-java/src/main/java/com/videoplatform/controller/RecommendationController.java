package com.videoplatform.controller;

import com.videoplatform.service.RecommendationService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/recommendations")
@RequiredArgsConstructor
public class RecommendationController {

    private final RecommendationService recommendationService;

    @GetMapping("/hot")
    public ResponseEntity<?> getHotRecommendations(
            @RequestParam(defaultValue = "20") int limit) {
        List<Long> videoIds = recommendationService.getHotRecommendations(limit);
        return ResponseEntity.ok(Map.of(
            "video_ids", videoIds,
            "type", "hot"
        ));
    }

    @GetMapping("/personalized")
    public ResponseEntity<?> getPersonalizedRecommendations(
            @AuthenticationPrincipal UserDetails userDetails,
            @RequestParam(defaultValue = "20") int limit) {
        Long userId = Long.parseLong(userDetails.getUsername());
        List<Long> videoIds = recommendationService.getPersonalizedRecommendations(userId, limit);
        return ResponseEntity.ok(Map.of(
            "video_ids", videoIds,
            "type", "personalized"
        ));
    }

    @GetMapping("/similar/{videoId}")
    public ResponseEntity<?> getSimilarVideos(
            @PathVariable Long videoId,
            @RequestParam(defaultValue = "10") int limit) {
        List<Long> videoIds = recommendationService.getSimilarVideos(videoId, limit);
        return ResponseEntity.ok(Map.of(
            "video_ids", videoIds,
            "type", "similar"
        ));
    }

    @GetMapping("/following")
    public ResponseEntity<?> getFollowingFeed(
            @AuthenticationPrincipal UserDetails userDetails,
            @RequestParam(defaultValue = "20") int limit) {
        Long userId = Long.parseLong(userDetails.getUsername());
        List<Long> videoIds = recommendationService.getFollowingFeed(userId, limit);
        return ResponseEntity.ok(Map.of(
            "video_ids", videoIds,
            "type", "following"
        ));
    }
}
