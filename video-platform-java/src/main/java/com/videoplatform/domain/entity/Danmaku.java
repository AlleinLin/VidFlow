package com.videoplatform.domain.entity;

import jakarta.persistence.*;
import lombok.*;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;

import java.time.LocalDateTime;

@Entity
@Table(name = "danmakus", indexes = {
    @Index(name = "idx_danmakus_video_id", columnList = "video_id"),
    @Index(name = "idx_danmakus_position", columnList = "video_id, position_seconds")
})
@EntityListeners(AuditingEntityListener.class)
@Getter
@Setter
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class Danmaku {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "video_id", nullable = false)
    private Video video;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "user_id", nullable = false)
    private User user;

    @Column(nullable = false, length = 500)
    private String content;

    @Column(name = "position_seconds", nullable = false)
    private Double positionSeconds;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false, length = 20)
    @Builder.Default
    private DanmakuStyle style = DanmakuStyle.SCROLL;

    @Column(length = 10)
    @Builder.Default
    private String color = "#FFFFFF";

    @Column(name = "font_size")
    @Builder.Default
    private Integer fontSize = 24;

    @CreatedDate
    @Column(name = "created_at", nullable = false, updatable = false)
    private LocalDateTime createdAt;

    public enum DanmakuStyle {
        SCROLL, TOP, BOTTOM
    }
}
