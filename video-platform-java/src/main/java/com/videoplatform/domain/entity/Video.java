package com.videoplatform.domain.entity;

import jakarta.persistence.*;
import lombok.*;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;

import java.time.LocalDateTime;

@Entity
@Table(name = "videos", indexes = {
    @Index(name = "idx_videos_user_id", columnList = "user_id"),
    @Index(name = "idx_videos_status", columnList = "status"),
    @Index(name = "idx_videos_category_id", columnList = "category_id"),
    @Index(name = "idx_videos_published_at", columnList = "published_at")
})
@EntityListeners(AuditingEntityListener.class)
@Getter
@Setter
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class Video {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "user_id", nullable = false)
    private User user;

    @Column(nullable = false, length = 200)
    private String title;

    @Column(columnDefinition = "TEXT")
    private String description;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false, length = 20)
    private VideoStatus status;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false, length = 20)
    @Builder.Default
    private VideoVisibility visibility = VideoVisibility.PUBLIC;

    @Column(name = "duration_seconds")
    private Integer durationSeconds;

    @Column(name = "original_filename", length = 500)
    private String originalFilename;

    @Column(name = "storage_key", columnDefinition = "TEXT")
    private String storageKey;

    @Column(name = "thumbnail_url", columnDefinition = "TEXT")
    private String thumbnailUrl;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "category_id")
    private Category category;

    @Column(name = "view_count", nullable = false)
    @Builder.Default
    private Long viewCount = 0L;

    @Column(name = "like_count", nullable = false)
    @Builder.Default
    private Long likeCount = 0L;

    @Column(name = "comment_count", nullable = false)
    @Builder.Default
    private Integer commentCount = 0;

    @Column(name = "published_at")
    private LocalDateTime publishedAt;

    @CreatedDate
    @Column(name = "created_at", nullable = false, updatable = false)
    private LocalDateTime createdAt;

    @LastModifiedDate
    @Column(name = "updated_at", nullable = false)
    private LocalDateTime updatedAt;

    public enum VideoStatus {
        UPLOADING, TRANSCODING, AUDITING, PUBLISHED, REJECTED, DELETED, HIDDEN
    }

    public enum VideoVisibility {
        PUBLIC, FOLLOWERS_ONLY, PRIVATE
    }
}
