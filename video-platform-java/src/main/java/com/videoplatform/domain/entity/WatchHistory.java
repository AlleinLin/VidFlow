package com.videoplatform.domain.entity;

import jakarta.persistence.*;
import lombok.*;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;

import java.time.LocalDateTime;

@Entity
@Table(name = "watch_history", indexes = {
    @Index(name = "idx_watch_history_user_id", columnList = "user_id"),
    @Index(name = "idx_watch_history_video_id", columnList = "video_id"),
    @Index(name = "idx_watch_history_watched_at", columnList = "watched_at")
}, uniqueConstraints = {
    @UniqueConstraint(name = "uk_watch_history_user_video", columnNames = {"user_id", "video_id"})
})
@EntityListeners(AuditingEntityListener.class)
@Getter
@Setter
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class WatchHistory {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "user_id", nullable = false)
    private User user;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "video_id", nullable = false)
    private Video video;

    @Column(name = "watch_duration", nullable = false)
    @Builder.Default
    private Long watchDuration = 0L;

    @Column(name = "watch_progress", nullable = false)
    @Builder.Default
    private Double watchProgress = 0.0;

    @Column(name = "last_position", nullable = false)
    @Builder.Default
    private Double lastPosition = 0.0;

    @Column(nullable = false)
    @Builder.Default
    private Boolean completed = false;

    @Column(name = "watched_at", nullable = false)
    private LocalDateTime watchedAt;

    @CreatedDate
    @Column(name = "created_at", nullable = false, updatable = false)
    private LocalDateTime createdAt;

    @LastModifiedDate
    @Column(name = "updated_at", nullable = false)
    private LocalDateTime updatedAt;
}
