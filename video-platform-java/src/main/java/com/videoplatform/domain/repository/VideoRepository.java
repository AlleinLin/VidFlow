package com.videoplatform.domain.repository;

import com.videoplatform.domain.entity.Video;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface VideoRepository extends JpaRepository<Video, Long> {
    
    Page<Video> findByUserIdOrderByCreatedAtDesc(Long userId, Pageable pageable);
    
    Page<Video> findByStatusOrderByViewCountDescCreatedAtDesc(Video.VideoStatus status, Pageable pageable);
    
    @Query("SELECT v FROM Video v WHERE v.status = 'PUBLISHED' AND v.visibility = 'PUBLIC' " +
           "AND (LOWER(v.title) LIKE LOWER(CONCAT('%', :keyword, '%')) " +
           "OR LOWER(v.description) LIKE LOWER(CONCAT('%', :keyword, '%')))")
    Page<Video> searchByKeyword(@Param("keyword") String keyword, Pageable pageable);
    
    @Query("SELECT v FROM Video v WHERE v.status = 'PUBLISHED' AND v.visibility = 'PUBLIC' " +
           "ORDER BY (v.viewCount * 1 + v.likeCount * 5 + v.commentCount * 10) DESC, v.createdAt DESC")
    List<Video> findTopByOrderByScoreDesc(Pageable pageable);
    
    @Query("SELECT v FROM Video v WHERE v.user.id IN :userIds AND v.status = 'PUBLISHED' AND v.visibility = 'PUBLIC' " +
           "ORDER BY v.createdAt DESC")
    List<Video> findByUserIdInOrderByCreatedAtDesc(@Param("userIds") List<Long> userIds, Pageable pageable);
    
    @Query("SELECT v FROM Video v WHERE v.category.id = :categoryId AND v.id != :videoId " +
           "AND v.status = 'PUBLISHED' AND v.visibility = 'PUBLIC' " +
           "ORDER BY v.viewCount DESC, v.createdAt DESC")
    List<Video> findByCategoryIdAndIdNot(@Param("categoryId") Long categoryId, @Param("videoId") Long videoId, Pageable pageable);
    
    @Modifying
    @Query("UPDATE Video v SET v.viewCount = v.viewCount + 1 WHERE v.id = :id")
    void incrementViewCount(@Param("id") Long id);
    
    @Modifying
    @Query("UPDATE Video v SET v.likeCount = v.likeCount + :delta WHERE v.id = :id")
    void incrementLikeCount(@Param("id") Long id, @Param("delta") int delta);
    
    @Modifying
    @Query("UPDATE Video v SET v.commentCount = v.commentCount + :delta WHERE v.id = :id")
    void incrementCommentCount(@Param("id") Long id, @Param("delta") int delta);
    
    List<Video> findByStatus(Video.VideoStatus status);
}
