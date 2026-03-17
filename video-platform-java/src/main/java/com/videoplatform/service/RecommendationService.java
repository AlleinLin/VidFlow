package com.videoplatform.service;

import com.videoplatform.domain.entity.User;
import com.videoplatform.domain.entity.UserFollow;
import com.videoplatform.domain.entity.Video;
import com.videoplatform.domain.repository.UserFollowRepository;
import com.videoplatform.domain.repository.UserRepository;
import com.videoplatform.domain.repository.VideoRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageImpl;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.*;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
public class RecommendationService {

    private final VideoRepository videoRepository;
    private final UserFollowRepository userFollowRepository;
    private final UserRepository userRepository;

    public List<Long> getHotRecommendations(int limit) {
        Pageable pageable = PageRequest.of(0, limit);
        List<Video> videos = videoRepository.findTopByOrderByScoreDesc(pageable);
        return videos.stream().map(Video::getId).collect(Collectors.toList());
    }

    public List<Long> getPersonalizedRecommendations(Long userId, int limit) {
        List<Long> followingIds = userFollowRepository.findFollowingIdsByFollowerId(userId);
        
        if (followingIds.isEmpty()) {
            return getHotRecommendations(limit);
        }
        
        Pageable pageable = PageRequest.of(0, limit);
        List<Video> videos = videoRepository.findByUserIdInOrderByCreatedAtDesc(followingIds, pageable);
        
        if (videos.size() < limit) {
            List<Long> existingIds = videos.stream().map(Video::getId).collect(Collectors.toList());
            List<Video> hotVideos = videoRepository.findTopByOrderByScoreDesc(PageRequest.of(0, limit - videos.size()));
            for (Video v : hotVideos) {
                if (!existingIds.contains(v.getId())) {
                    videos.add(v);
                }
            }
        }
        
        return videos.stream().map(Video::getId).collect(Collectors.toList());
    }

    public List<Long> getSimilarVideos(Long videoId, int limit) {
        Optional<Video> videoOpt = videoRepository.findById(videoId);
        if (videoOpt.isEmpty()) {
            return getHotRecommendations(limit);
        }
        
        Video video = videoOpt.get();
        Long categoryId = video.getCategory() != null ? video.getCategory().getId() : null;
        
        if (categoryId == null) {
            return getHotRecommendations(limit);
        }
        
        Pageable pageable = PageRequest.of(0, limit + 1);
        List<Video> similarVideos = videoRepository.findByCategoryIdAndIdNot(categoryId, videoId, pageable);
        
        return similarVideos.stream()
            .limit(limit)
            .map(Video::getId)
            .collect(Collectors.toList());
    }

    public List<Long> getFollowingFeed(Long userId, int limit) {
        List<Long> followingIds = userFollowRepository.findFollowingIdsByFollowerId(userId);
        
        if (followingIds.isEmpty()) {
            return getHotRecommendations(limit);
        }
        
        Pageable pageable = PageRequest.of(0, limit);
        List<Video> videos = videoRepository.findByUserIdInOrderByCreatedAtDesc(followingIds, pageable);
        
        return videos.stream().map(Video::getId).collect(Collectors.toList());
    }
}
