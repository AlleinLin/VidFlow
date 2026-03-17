package com.videoplatform.controller;

import com.videoplatform.dto.video.CreateVideoRequest;
import com.videoplatform.dto.video.VideoListResponse;
import com.videoplatform.dto.video.VideoResponse;
import com.videoplatform.service.VideoService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/v1/videos")
@RequiredArgsConstructor
@Tag(name = "Videos", description = "Video management endpoints")
public class VideoController {

    private final VideoService videoService;

    @GetMapping
    @Operation(summary = "List videos")
    public ResponseEntity<VideoListResponse> list(
            @RequestParam(required = false) String keyword,
            @RequestParam(required = false) Long categoryId,
            @RequestParam(defaultValue = "1") int page,
            @RequestParam(defaultValue = "20") int size
    ) {
        return ResponseEntity.ok(videoService.list(keyword, categoryId, page, size));
    }

    @GetMapping("/{id}")
    @Operation(summary = "Get video by ID")
    public ResponseEntity<VideoResponse> get(@PathVariable Long id) {
        VideoResponse video = videoService.getById(id);
        videoService.incrementViewCount(id);
        return ResponseEntity.ok(video);
    }

    @PostMapping
    @Operation(summary = "Create a new video")
    public ResponseEntity<VideoResponse> create(
            @RequestAttribute("userId") Long userId,
            @Valid @RequestBody CreateVideoRequest request
    ) {
        return ResponseEntity.ok(videoService.create(userId, request));
    }

    @PutMapping("/{id}")
    @Operation(summary = "Update video")
    public ResponseEntity<VideoResponse> update(
            @PathVariable Long id,
            @RequestAttribute("userId") Long userId,
            @RequestBody CreateVideoRequest request
    ) {
        return ResponseEntity.ok(videoService.update(id, userId, request.getTitle(), request.getDescription(), request.getVisibility()));
    }

    @DeleteMapping("/{id}")
    @Operation(summary = "Delete video")
    public ResponseEntity<Void> delete(
            @PathVariable Long id,
            @RequestAttribute("userId") Long userId
    ) {
        videoService.delete(id, userId);
        return ResponseEntity.ok().build();
    }

    @PostMapping("/{id}/publish")
    @Operation(summary = "Publish video")
    public ResponseEntity<Void> publish(
            @PathVariable Long id,
            @RequestAttribute("userId") Long userId
    ) {
        videoService.publish(id, userId);
        return ResponseEntity.ok().build();
    }
}
