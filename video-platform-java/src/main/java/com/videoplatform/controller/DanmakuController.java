package com.videoplatform.controller;

import com.videoplatform.domain.entity.Danmaku;
import com.videoplatform.service.DanmakuService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/interactions/danmakus")
@RequiredArgsConstructor
public class DanmakuController {

    private final DanmakuService danmakuService;

    @PostMapping
    public ResponseEntity<?> createDanmaku(
            @AuthenticationPrincipal UserDetails userDetails,
            @RequestBody DanmakuRequest request) {
        Long userId = Long.parseLong(userDetails.getUsername());
        Danmaku danmaku = danmakuService.createDanmaku(
            userId,
            request.getVideoId(),
            request.getContent(),
            request.getPositionSeconds(),
            request.getStyle(),
            request.getColor(),
            request.getFontSize()
        );
        return ResponseEntity.ok(danmaku);
    }

    @GetMapping("/video/{videoId}")
    public ResponseEntity<?> getDanmakus(
            @PathVariable Long videoId,
            @RequestParam(required = false) Double start,
            @RequestParam(required = false) Double end) {
        List<Danmaku> danmakus = danmakuService.getDanmakusByVideo(videoId, start, end);
        return ResponseEntity.ok(Map.of("danmakus", danmakus));
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<?> deleteDanmaku(@PathVariable Long id) {
        danmakuService.deleteDanmaku(id);
        return ResponseEntity.noContent().build();
    }

    public static class DanmakuRequest {
        private Long videoId;
        private String content;
        private Double positionSeconds;
        private Danmaku.DanmakuStyle style;
        private String color;
        private Integer fontSize;

        public Long getVideoId() { return videoId; }
        public void setVideoId(Long videoId) { this.videoId = videoId; }
        public String getContent() { return content; }
        public void setContent(String content) { this.content = content; }
        public Double getPositionSeconds() { return positionSeconds; }
        public void setPositionSeconds(Double positionSeconds) { this.positionSeconds = positionSeconds; }
        public Danmaku.DanmakuStyle getStyle() { return style; }
        public void setStyle(Danmaku.DanmakuStyle style) { this.style = style; }
        public String getColor() { return color; }
        public void setColor(String color) { this.color = color; }
        public Integer getFontSize() { return fontSize; }
        public void setFontSize(Integer fontSize) { this.fontSize = fontSize; }
    }
}
