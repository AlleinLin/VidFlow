package com.videoplatform.service;

import com.videoplatform.domain.entity.User;
import com.videoplatform.domain.entity.UserFollow;
import com.videoplatform.domain.repository.UserFollowRepository;
import com.videoplatform.domain.repository.UserRepository;
import com.videoplatform.security.JwtTokenProvider;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.HashMap;
import java.util.Map;
import java.util.Optional;

@Service
@RequiredArgsConstructor
public class UserService {

    private final UserRepository userRepository;
    private final UserFollowRepository userFollowRepository;
    private final PasswordEncoder passwordEncoder;
    private final JwtTokenProvider jwtTokenProvider;

    @Transactional
    public User register(String username, String email, String password, String displayName) {
        if (userRepository.existsByUsername(username)) {
            throw new RuntimeException("Username already exists");
        }
        if (userRepository.existsByEmail(email)) {
            throw new RuntimeException("Email already exists");
        }

        User user = User.builder()
            .username(username)
            .email(email)
            .passwordHash(passwordEncoder.encode(password))
            .displayName(displayName != null ? displayName : username)
            .role(User.UserRole.USER)
            .status(User.UserStatus.ACTIVE)
            .followerCount(0L)
            .followingCount(0L)
            .build();

        return userRepository.save(user);
    }

    @Transactional
    public Map<String, Object> login(String username, String password) {
        User user = userRepository.findByUsername(username)
            .orElseThrow(() -> new RuntimeException("User not found"));

        if (!passwordEncoder.matches(password, user.getPasswordHash())) {
            throw new RuntimeException("Invalid password");
        }

        user.setLastLoginAt(LocalDateTime.now());
        userRepository.save(user);

        String accessToken = jwtTokenProvider.generateToken(user.getId().toString());
        String refreshToken = jwtTokenProvider.generateRefreshToken(user.getId().toString());

        Map<String, Object> response = new HashMap<>();
        response.put("user", user);
        response.put("access_token", accessToken);
        response.put("refresh_token", refreshToken);
        response.put("token_type", "Bearer");
        response.put("expires_in", 86400);

        return response;
    }

    public User getUserById(Long id) {
        return userRepository.findById(id)
            .orElseThrow(() -> new RuntimeException("User not found"));
    }

    @Transactional
    public User updateProfile(Long userId, String displayName, String bio, String avatarUrl) {
        User user = getUserById(userId);
        if (displayName != null) user.setDisplayName(displayName);
        if (bio != null) user.setBio(bio);
        if (avatarUrl != null) user.setAvatarUrl(avatarUrl);
        return userRepository.save(user);
    }

    @Transactional
    public void followUser(Long followerId, Long followingId) {
        if (followerId.equals(followingId)) {
            throw new RuntimeException("Cannot follow yourself");
        }

        if (userFollowRepository.existsByFollowerIdAndFollowingId(followerId, followingId)) {
            return;
        }

        User follower = getUserById(followerId);
        User following = getUserById(followingId);

        UserFollow follow = UserFollow.builder()
            .follower(follower)
            .following(following)
            .build();
        userFollowRepository.save(follow);

        follower.setFollowingCount(follower.getFollowingCount() + 1);
        following.setFollowerCount(following.getFollowerCount() + 1);
        userRepository.save(follower);
        userRepository.save(following);
    }

    @Transactional
    public void unfollowUser(Long followerId, Long followingId) {
        if (!userFollowRepository.existsByFollowerIdAndFollowingId(followerId, followingId)) {
            return;
        }

        userFollowRepository.deleteByFollowerIdAndFollowingId(followerId, followingId);

        User follower = getUserById(followerId);
        User following = getUserById(followingId);

        follower.setFollowingCount(Math.max(0, follower.getFollowingCount() - 1));
        following.setFollowerCount(Math.max(0, following.getFollowerCount() - 1));
        userRepository.save(follower);
        userRepository.save(following);
    }

    public Page<User> getFollowers(Long userId, int page, int size) {
        Pageable pageable = PageRequest.of(page, size);
        Page<UserFollow> follows = userFollowRepository.findByFollowingId(userId, pageable);
        return follows.map(UserFollow::getFollower);
    }

    public Page<User> getFollowing(Long userId, int page, int size) {
        Pageable pageable = PageRequest.of(page, size);
        Page<UserFollow> follows = userFollowRepository.findByFollowerId(userId, pageable);
        return follows.map(UserFollow::getFollowing);
    }

    public Map<String, Object> refreshToken(String refreshToken) {
        Long userId = jwtTokenProvider.getUserIdFromToken(refreshToken);
        User user = getUserById(userId);

        String newAccessToken = jwtTokenProvider.generateToken(user.getId().toString());
        String newRefreshToken = jwtTokenProvider.generateRefreshToken(user.getId().toString());

        Map<String, Object> response = new HashMap<>();
        response.put("access_token", newAccessToken);
        response.put("refresh_token", newRefreshToken);
        response.put("token_type", "Bearer");
        response.put("expires_in", 86400);

        return response;
    }
}
