package com.videoplatform.service;

import com.videoplatform.domain.entity.User;
import com.videoplatform.domain.entity.Video;
import com.videoplatform.domain.repository.UserRepository;
import com.videoplatform.domain.repository.VideoRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageImpl;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
public class SearchService {

    private final VideoRepository videoRepository;
    private final UserRepository userRepository;

    public Page<VideoSearchResult> searchVideos(String query, int page, int size) {
        Pageable pageable = PageRequest.of(page, size);
        Page<Video> videos = videoRepository.searchByKeyword(query, pageable);
        
        List<VideoSearchResult> results = videos.getContent().stream()
            .map(this::toVideoSearchResult)
            .collect(Collectors.toList());
        
        return new PageImpl<>(results, pageable, videos.getTotalElements());
    }

    public Page<UserSearchResult> searchUsers(String query, int page, int size) {
        Pageable pageable = PageRequest.of(page, size);
        Page<User> users = userRepository.searchByKeyword(query, pageable);
        
        List<UserSearchResult> results = users.getContent().stream()
            .map(this::toUserSearchResult)
            .collect(Collectors.toList());
        
        return new PageImpl<>(results, pageable, users.getTotalElements());
    }

    public Map<String, Object> searchAll(String query, int page, int size) {
        int halfSize = Math.max(1, size / 2);
        
        Page<VideoSearchResult> videos = searchVideos(query, page, halfSize);
        Page<UserSearchResult> users = searchUsers(query, page, halfSize);
        
        List<SearchResult> results = new ArrayList<>();
        
        videos.getContent().forEach(v -> {
            SearchResult sr = new SearchResult();
            sr.setType("video");
            sr.setId(v.getId());
            sr.setTitle(v.getTitle());
            sr.setDescription(v.getDescription());
            sr.setScore(v.getScore());
            sr.setData(v);
            results.add(sr);
        });
        
        users.getContent().forEach(u -> {
            SearchResult sr = new SearchResult();
            sr.setType("user");
            sr.setId(u.getId());
            sr.setTitle(u.getDisplayName());
            sr.setDescription(u.getBio());
            sr.setScore(u.getScore());
            sr.setData(u);
            results.add(sr);
        });
        
        Map<String, Object> response = new HashMap<>();
        response.put("results", results);
        response.put("total", videos.getTotalElements() + users.getTotalElements());
        response.put("page", page);
        response.put("page_size", size);
        
        return response;
    }

    private VideoSearchResult toVideoSearchResult(Video video) {
        VideoSearchResult result = new VideoSearchResult();
        result.setId(video.getId());
        result.setTitle(video.getTitle());
        result.setDescription(video.getDescription());
        result.setThumbnailUrl(video.getThumbnailUrl());
        result.setViewCount(video.getViewCount());
        result.setLikeCount(video.getLikeCount());
        result.setDuration(video.getDurationSeconds());
        result.setPublishedAt(video.getPublishedAt());
        if (video.getUser() != null) {
            result.setAuthorId(video.getUser().getId());
            result.setAuthorName(video.getUser().getDisplayName());
        }
        result.setScore(1.0);
        return result;
    }

    private UserSearchResult toUserSearchResult(User user) {
        UserSearchResult result = new UserSearchResult();
        result.setId(user.getId());
        result.setUsername(user.getUsername());
        result.setDisplayName(user.getDisplayName());
        result.setAvatarUrl(user.getAvatarUrl());
        result.setBio(user.getBio());
        result.setFollowerCount(user.getFollowerCount());
        result.setScore(1.0);
        return result;
    }

    public static class VideoSearchResult {
        private Long id;
        private String title;
        private String description;
        private String thumbnailUrl;
        private Long viewCount;
        private Long likeCount;
        private Integer duration;
        private java.time.LocalDateTime publishedAt;
        private Long authorId;
        private String authorName;
        private Double score;

        public Long getId() { return id; }
        public void setId(Long id) { this.id = id; }
        public String getTitle() { return title; }
        public void setTitle(String title) { this.title = title; }
        public String getDescription() { return description; }
        public void setDescription(String description) { this.description = description; }
        public String getThumbnailUrl() { return thumbnailUrl; }
        public void setThumbnailUrl(String thumbnailUrl) { this.thumbnailUrl = thumbnailUrl; }
        public Long getViewCount() { return viewCount; }
        public void setViewCount(Long viewCount) { this.viewCount = viewCount; }
        public Long getLikeCount() { return likeCount; }
        public void setLikeCount(Long likeCount) { this.likeCount = likeCount; }
        public Integer getDuration() { return duration; }
        public void setDuration(Integer duration) { this.duration = duration; }
        public java.time.LocalDateTime getPublishedAt() { return publishedAt; }
        public void setPublishedAt(java.time.LocalDateTime publishedAt) { this.publishedAt = publishedAt; }
        public Long getAuthorId() { return authorId; }
        public void setAuthorId(Long authorId) { this.authorId = authorId; }
        public String getAuthorName() { return authorName; }
        public void setAuthorName(String authorName) { this.authorName = authorName; }
        public Double getScore() { return score; }
        public void setScore(Double score) { this.score = score; }
    }

    public static class UserSearchResult {
        private Long id;
        private String username;
        private String displayName;
        private String avatarUrl;
        private String bio;
        private Long followerCount;
        private Double score;

        public Long getId() { return id; }
        public void setId(Long id) { this.id = id; }
        public String getUsername() { return username; }
        public void setUsername(String username) { this.username = username; }
        public String getDisplayName() { return displayName; }
        public void setDisplayName(String displayName) { this.displayName = displayName; }
        public String getAvatarUrl() { return avatarUrl; }
        public void setAvatarUrl(String avatarUrl) { this.avatarUrl = avatarUrl; }
        public String getBio() { return bio; }
        public void setBio(String bio) { this.bio = bio; }
        public Long getFollowerCount() { return followerCount; }
        public void setFollowerCount(Long followerCount) { this.followerCount = followerCount; }
        public Double getScore() { return score; }
        public void setScore(Double score) { this.score = score; }
    }

    public static class SearchResult {
        private String type;
        private Long id;
        private String title;
        private String description;
        private Double score;
        private Object data;

        public String getType() { return type; }
        public void setType(String type) { this.type = type; }
        public Long getId() { return id; }
        public void setId(Long id) { this.id = id; }
        public String getTitle() { return title; }
        public void setTitle(String title) { this.title = title; }
        public String getDescription() { return description; }
        public void setDescription(String description) { this.description = description; }
        public Double getScore() { return score; }
        public void setScore(Double score) { this.score = score; }
        public Object getData() { return data; }
        public void setData(Object data) { this.data = data; }
    }
}
