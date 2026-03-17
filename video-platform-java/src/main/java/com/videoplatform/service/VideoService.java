package com.videoplatform.service;

import com.videoplatform.domain.entity.Video;
import com.videoplatform.domain.repository.VideoRepository;
import com.videoplatform.dto.video.CreateVideoRequest;
import com.videoplatform.dto.video.VideoListResponse;
import com.videoplatform.dto.video.VideoResponse;
import com.videoplatform.exception.BusinessException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.HashMap;
import java.util.Map;

@Service
@RequiredArgsConstructor
@Slf4j
public class VideoService {

    private final VideoRepository videoRepository;
    private final KafkaTemplate<String, Object> kafkaTemplate;

    @Transactional
    public VideoResponse create(Long userId, CreateVideoRequest request) {
        Video video = Video.builder()
                .user(com.videoplatform.domain.entity.User.builder().id(userId).build())
                .title(request.getTitle())
                .description(request.getDescription())
                .status(Video.VideoStatus.UPLOADING)
                .visibility(Video.VideoVisibility.valueOf(request.getVisibility() != null ? request.getVisibility() : "PUBLIC"))
                .build();

        if (request.getCategoryId() != null) {
            video.setCategory(com.videoplatform.domain.entity.Category.builder().id(request.getCategoryId()).build());
        }

        video = videoRepository.save(video);

        log.info("Video created: videoId={}, userId={}", video.getId(), userId);

        return mapToResponse(video);
    }

    @Cacheable(value = "video:meta", key = "#id")
    @Transactional(readOnly = true)
    public VideoResponse getById(Long id) {
        Video video = videoRepository.findById(id)
                .orElseThrow(() -> new BusinessException("Video not found"));
        return mapToResponse(video);
    }

    @CacheEvict(value = "video:meta", key = "#id")
    @Transactional
    public VideoResponse update(Long id, Long userId, String title, String description, String visibility) {
        Video video = videoRepository.findById(id)
                .orElseThrow(() -> new BusinessException("Video not found"));

        if (!video.getUser().getId().equals(userId)) {
            throw new BusinessException("Unauthorized");
        }

        if (title != null) {
            video.setTitle(title);
        }
        if (description != null) {
            video.setDescription(description);
        }
        if (visibility != null) {
            video.setVisibility(Video.VideoVisibility.valueOf(visibility));
        }

        video = videoRepository.save(video);
        return mapToResponse(video);
    }

    @CacheEvict(value = "video:meta", key = "#id")
    @Transactional
    public void delete(Long id, Long userId) {
        Video video = videoRepository.findById(id)
                .orElseThrow(() -> new BusinessException("Video not found"));

        if (!video.getUser().getId().equals(userId)) {
            throw new BusinessException("Unauthorized");
        }

        video.setStatus(Video.VideoStatus.DELETED);
        videoRepository.save(video);

        log.info("Video deleted: videoId={}", id);
    }

    @Transactional(readOnly = true)
    public VideoListResponse list(String keyword, Long categoryId, int page, int size) {
        Pageable pageable = PageRequest.of(page - 1, size, Sort.by(Sort.Direction.DESC, "createdAt"));

        Page<Video> videoPage;
        if (keyword != null && !keyword.isBlank()) {
            videoPage = videoRepository.searchByKeyword(Video.VideoStatus.PUBLISHED, keyword, pageable);
        } else if (categoryId != null) {
            videoPage = videoRepository.findByCategoryAndStatus(categoryId, Video.VideoStatus.PUBLISHED, pageable);
        } else {
            videoPage = videoRepository.findByStatusAndVisibility(Video.VideoStatus.PUBLISHED, Video.VideoVisibility.PUBLIC, pageable);
        }

        return VideoListResponse.builder()
                .videos(videoPage.getContent().stream().map(this::mapToResponse).toList())
                .total(videoPage.getTotalElements())
                .page(page)
                .pageSize(size)
                .build();
    }

    @Transactional
    public void incrementViewCount(Long id) {
        videoRepository.incrementViewCount(id);
    }

    @Transactional
    public void publish(Long id, Long userId) {
        Video video = videoRepository.findById(id)
                .orElseThrow(() -> new BusinessException("Video not found"));

        if (!video.getUser().getId().equals(userId)) {
            throw new BusinessException("Unauthorized");
        }

        videoRepository.updateStatusAndPublishedAt(id, Video.VideoStatus.PUBLISHED, LocalDateTime.now());

        Map<String, Object> event = new HashMap<>();
        event.put("videoId", id);
        event.put("userId", userId);
        event.put("timestamp", System.currentTimeMillis());
        kafkaTemplate.send("video.published", String.valueOf(id), event);

        log.info("Video published: videoId={}", id);
    }

    private VideoResponse mapToResponse(Video video) {
        return VideoResponse.builder()
                .id(video.getId())
                .userId(video.getUser().getId())
                .title(video.getTitle())
                .description(video.getDescription())
                .status(video.getStatus().name())
                .visibility(video.getVisibility().name())
                .durationSeconds(video.getDurationSeconds())
                .thumbnailUrl(video.getThumbnailUrl())
                .categoryId(video.getCategory() != null ? video.getCategory().getId() : null)
                .viewCount(video.getViewCount())
                .likeCount(video.getLikeCount())
                .commentCount(video.getCommentCount())
                .publishedAt(video.getPublishedAt())
                .createdAt(video.getCreatedAt())
                .build();
    }
}
