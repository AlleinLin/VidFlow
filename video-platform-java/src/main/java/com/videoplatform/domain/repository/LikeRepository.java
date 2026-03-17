package com.videoplatform.domain.repository;

import com.videoplatform.domain.entity.Like;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface LikeRepository extends JpaRepository<Like, Long> {

    Optional<Like> findByUserIdAndTargetIdAndType(Long userId, Long targetId, Like.LikeType type);

    boolean existsByUserIdAndTargetIdAndType(Long userId, Long targetId, Like.LikeType type);

    long countByTargetIdAndType(Long targetId, Like.LikeType type);

    void deleteByUserIdAndTargetIdAndType(Long userId, Long targetId, Like.LikeType type);
}
