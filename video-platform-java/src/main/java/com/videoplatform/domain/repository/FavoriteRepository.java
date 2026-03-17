package com.videoplatform.domain.repository;

import com.videoplatform.domain.entity.Favorite;
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
public interface FavoriteRepository extends JpaRepository<Favorite, Long> {
    
    Optional<Favorite> findByUserIdAndVideoId(Long userId, Long videoId);
    
    boolean existsByUserIdAndVideoId(Long userId, Long videoId);
    
    Page<Favorite> findByUserIdOrderByCreatedAtDesc(Long userId, Pageable pageable);
    
    @Query("SELECT f.video.id FROM Favorite f WHERE f.user.id = :userId")
    List<Long> findVideoIdsByUserId(@Param("userId") Long userId);
    
    @Modifying
    @Query("DELETE FROM Favorite f WHERE f.user.id = :userId AND f.video.id = :videoId")
    void deleteByUserIdAndVideoId(@Param("userId") Long userId, @Param("videoId") Long videoId);
    
    @Query("SELECT COUNT(f) FROM Favorite f WHERE f.video.id = :videoId")
    Long countByVideoId(@Param("videoId") Long videoId);
}
