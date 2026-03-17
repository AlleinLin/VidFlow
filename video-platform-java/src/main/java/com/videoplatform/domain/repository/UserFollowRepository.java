package com.videoplatform.domain.repository;

import com.videoplatform.domain.entity.UserFollow;
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
public interface UserFollowRepository extends JpaRepository<UserFollow, Long> {
    
    Optional<UserFollow> findByFollowerIdAndFollowingId(Long followerId, Long followingId);
    
    boolean existsByFollowerIdAndFollowingId(Long followerId, Long followingId);
    
    @Query("SELECT uf.following FROM UserFollow uf WHERE uf.follower.id = :userId")
    Page<UserFollow> findByFollowerId(@Param("userId") Long userId, Pageable pageable);
    
    @Query("SELECT uf.follower FROM UserFollow uf WHERE uf.following.id = :userId")
    Page<UserFollow> findByFollowingId(@Param("userId") Long userId, Pageable pageable);
    
    @Query("SELECT uf.following.id FROM UserFollow uf WHERE uf.follower.id = :userId")
    List<Long> findFollowingIdsByFollowerId(@Param("userId") Long userId);
    
    @Modifying
    @Query("DELETE FROM UserFollow uf WHERE uf.follower.id = :followerId AND uf.following.id = :followingId")
    void deleteByFollowerIdAndFollowingId(@Param("followerId") Long followerId, @Param("followingId") Long followingId);
    
    @Query("SELECT COUNT(uf) FROM UserFollow uf WHERE uf.follower.id = :userId")
    Long countByFollowerId(@Param("userId") Long userId);
    
    @Query("SELECT COUNT(uf) FROM UserFollow uf WHERE uf.following.id = :userId")
    Long countByFollowingId(@Param("userId") Long userId);
}
