package com.videoplatform.service;

import com.videoplatform.domain.entity.Favorite;
import com.videoplatform.domain.entity.Like;
import com.videoplatform.domain.entity.Video;
import com.videoplatform.domain.repository.FavoriteRepository;
import com.videoplatform.domain.repository.LikeRepository;
import com.videoplatform.domain.repository.VideoRepository;
import com.videoplatform.dto.video.VideoResponse;
import com.videoplatform.exception.BusinessException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Sort;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Service
@RequiredArgsConstructor
@Slf4j
public class InteractionService {

    private final LikeRepository likeRepository;
    private final FavoriteRepository favoriteRepository;
    private final VideoRepository videoRepository;
    private final KafkaTemplate<String, Object> kafkaTemplate;

    @Transactional
    public void likeVideo(Long userId, Long videoId) {
        if (likeRepository.existsByUserIdAndTargetIdAndType(userId, videoId, Like.LikeType.VIDEO)) {
            throw new BusinessException("Already liked");
        }

        Like like = Like.builder()
                .user(com.videoplatform.domain.entity.User.builder().id(userId).build())
                .targetId(videoId)
                .type(Like.LikeType.VIDEO)
                .build();
        likeRepository.save(like);

        videoRepository.incrementLikeCount(videoId);

        publishInteractionEvent(userId, videoId, "like");
    }

    @Transactional
    public void unlikeVideo(Long userId, Long videoId) {
        if (!likeRepository.existsByUserIdAndTargetIdAndType(userId, videoId, Like.LikeType.VIDEO)) {
            throw new BusinessException("Not liked");
        }

        likeRepository.deleteByUserIdAndTargetIdAndType(userId, videoId, Like.LikeType.VIDEO);
        videoRepository.decrementLikeCount(videoId);

        publishInteractionEvent(userId, videoId, "unlike");
    }

    @Transactional(readOnly = true)
    public boolean isVideoLiked(Long userId, Long videoId) {
        return likeRepository.existsByUserIdAndTargetIdAndType(userId, videoId, Like.LikeType.VIDEO);
    }

    @Transactional
    public void favoriteVideo(Long userId, Long videoId) {
        if (favoriteRepository.existsByUserIdAndVideoId(userId, videoId)) {
            throw new BusinessException("Already favorited");
        }

        Favorite favorite = Favorite.builder()
                .user(com.videoplatform.domain.entity.User.builder().id(userId).build())
                .video(Video.builder().id(videoId).build())
                .build();
        favoriteRepository.save(favorite);

        publishInteractionEvent(userId, videoId, "favorite");
    }

    @Transactional
    public void unfavoriteVideo(Long userId, Long videoId) {
        if (!favoriteRepository.existsByUserIdAndVideoId(userId, videoId)) {
            throw new BusinessException("Not favorited");
        }

        favoriteRepository.deleteByUserIdAndVideoId(userId, videoId);
        publishInteractionEvent(userId, videoId, "unfavorite");
    }

    @Transactional(readOnly = true)
    public List<VideoResponse> getLikedVideos(Long userId, int page, int size) {
        PageRequest pageable = PageRequest.of(page - 1, size, Sort.by(Sort.Direction.DESC, "createdAt"));
        Page<Like> likes = likeRepository.findByUserIdAndType(userId, Like.LikeType.VIDEO, pageable);

        return likes.getContent().stream()
                .map(like -> videoRepository.findById(like.getTargetId()))
                .filter(opt -> opt.isPresent())
                .map(opt -> mapToResponse(opt.get()))
                .toList();
    }

    @Transactional(readOnly = true)
    public List<VideoResponse> getFavoriteVideos(Long userId, int page, int size) {
        PageRequest pageable = PageRequest.of(page - 1, size, Sort.by(Sort.Direction.DESC, "createdAt"));
        Page<Favorite> favorites = favoriteRepository.findByUserId(userId, pageable);

        return favorites.getContent().stream()
                .map(fav -> mapToResponse(fav.getVideo()))
                .toList();
    }

    private void publishInteractionEvent(Long userId, Long videoId, String action) {
        Map<String, Object> event = new HashMap<>();
        event.put("userId", userId);
        event.put("videoId", videoId);
        event.put("action", action);
        event.put("timestamp", System.currentTimeMillis());
        kafkaTemplate.send("user.interaction", String.valueOf(videoId), event);
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
