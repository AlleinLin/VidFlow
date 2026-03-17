package com.videoplatform.domain.repository;

import com.videoplatform.domain.entity.WatchHistory;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;

@Repository
public interface WatchHistoryRepository extends JpaRepository<WatchHistory, Long> {
    
    Optional<WatchHistory> findByUserIdAndVideoId(Long userId, Long videoId);
    
    Page<WatchHistory> findByUserIdOrderByWatchedAtDesc(Long userId, Pageable pageable);
    
    @Query("SELECT wh FROM WatchHistory wh WHERE wh.user.id = :userId AND wh.completed = false AND wh.watchProgress > 0.05 ORDER BY wh.watchedAt DESC")
    List<WatchHistory> findContinueWatching(@Param("userId") Long userId, Pageable pageable);
    
    @Modifying
    @Query("DELETE FROM WatchHistory wh WHERE wh.user.id = :userId AND wh.video.id = :videoId")
    void deleteByUserIdAndVideoId(@Param("userId") Long userId, @Param("videoId") Long videoId);
    
    @Modifying
    @Query("DELETE FROM WatchHistory wh WHERE wh.user.id = :userId")
    void deleteByUserId(@Param("userId") Long userId);
    
    @Query("SELECT COUNT(wh) FROM WatchHistory wh WHERE wh.user.id = :userId")
    Long countByUserId(@Param("userId") Long userId);
}
