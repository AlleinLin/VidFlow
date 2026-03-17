package com.videoplatform.domain.repository;

import com.videoplatform.domain.entity.Comment;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface CommentRepository extends JpaRepository<Comment, Long> {

    Page<Comment> findByVideoIdAndStatusAndParentIsNull(
        Long videoId,
        Comment.CommentStatus status,
        Pageable pageable
    );

    List<Comment> findByRootIdAndStatusOrderByCreatedAtAsc(
        Long rootId,
        Comment.CommentStatus status,
        Pageable pageable
    );

    long countByVideoIdAndStatus(Long videoId, Comment.CommentStatus status);

    @Modifying
    @Query("UPDATE Comment c SET c.likeCount = c.likeCount + 1 WHERE c.id = :id")
    void incrementLikeCount(@Param("id") Long id);

    @Modifying
    @Query("UPDATE Comment c SET c.likeCount = GREATEST(c.likeCount - 1, 0) WHERE c.id = :id")
    void decrementLikeCount(@Param("id") Long id);
}
