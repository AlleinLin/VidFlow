package com.videoplatform.domain.repository;

import com.videoplatform.domain.entity.Danmaku;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface DanmakuRepository extends JpaRepository<Danmaku, Long> {
    
    @Query("SELECT d FROM Danmaku d WHERE d.video.id = :videoId AND d.positionSeconds >= :startTime AND d.positionSeconds <= :endTime ORDER BY d.positionSeconds")
    List<Danmaku> findByVideoIdAndPositionRange(
        @Param("videoId") Long videoId,
        @Param("startTime") Double startTime,
        @Param("endTime") Double endTime
    );
    
    List<Danmaku> findByVideoIdOrderByPositionSeconds(Long videoId);
    
    @Query("SELECT COUNT(d) FROM Danmaku d WHERE d.video.id = :videoId")
    Long countByVideoId(@Param("videoId") Long videoId);
}
