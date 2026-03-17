package com.videoplatform.controller;

import com.videoplatform.domain.entity.User;
import com.videoplatform.domain.entity.UserFollow;
import com.videoplatform.service.UserService;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.web.bind.annotation.*;

import java.util.Map;

@RestController
@RequestMapping("/api/v1/users")
@RequiredArgsConstructor
public class UserController {

    private final UserService userService;

    @GetMapping("/me")
    public ResponseEntity<?> getCurrentUser(@AuthenticationPrincipal UserDetails userDetails) {
        Long userId = Long.parseLong(userDetails.getUsername());
        User user = userService.getUserById(userId);
        return ResponseEntity.ok(user);
    }

    @PutMapping("/me")
    public ResponseEntity<?> updateProfile(
            @AuthenticationPrincipal UserDetails userDetails,
            @RequestBody UpdateProfileRequest request) {
        Long userId = Long.parseLong(userDetails.getUsername());
        User user = userService.updateProfile(userId, request.getDisplayName(), request.getBio(), request.getAvatarUrl());
        return ResponseEntity.ok(user);
    }

    @GetMapping("/{id}")
    public ResponseEntity<?> getUserById(@PathVariable Long id) {
        User user = userService.getUserById(id);
        return ResponseEntity.ok(Map.of(
            "id", user.getId(),
            "username", user.getUsername(),
            "display_name", user.getDisplayName(),
            "avatar_url", user.getAvatarUrl(),
            "bio", user.getBio(),
            "follower_count", user.getFollowerCount(),
            "following_count", user.getFollowingCount(),
            "created_at", user.getCreatedAt()
        ));
    }

    @PostMapping("/{id}/follow")
    public ResponseEntity<?> followUser(
            @AuthenticationPrincipal UserDetails userDetails,
            @PathVariable Long id) {
        Long userId = Long.parseLong(userDetails.getUsername());
        userService.followUser(userId, id);
        return ResponseEntity.ok(Map.of("message", "Followed successfully"));
    }

    @DeleteMapping("/{id}/follow")
    public ResponseEntity<?> unfollowUser(
            @AuthenticationPrincipal UserDetails userDetails,
            @PathVariable Long id) {
        Long userId = Long.parseLong(userDetails.getUsername());
        userService.unfollowUser(userId, id);
        return ResponseEntity.ok(Map.of("message", "Unfollowed successfully"));
    }

    @GetMapping("/{id}/followers")
    public ResponseEntity<?> getFollowers(
            @PathVariable Long id,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size) {
        Page<User> followers = userService.getFollowers(id, page, size);
        return ResponseEntity.ok(Map.of(
            "users", followers.getContent(),
            "total", followers.getTotalElements(),
            "page", page,
            "page_size", size
        ));
    }

    @GetMapping("/{id}/following")
    public ResponseEntity<?> getFollowing(
            @PathVariable Long id,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size) {
        Page<User> following = userService.getFollowing(id, page, size);
        return ResponseEntity.ok(Map.of(
            "users", following.getContent(),
            "total", following.getTotalElements(),
            "page", page,
            "page_size", size
        ));
    }

    public static class UpdateProfileRequest {
        private String displayName;
        private String bio;
        private String avatarUrl;

        public String getDisplayName() { return displayName; }
        public void setDisplayName(String displayName) { this.displayName = displayName; }
        public String getBio() { return bio; }
        public void setBio(String bio) { this.bio = bio; }
        public String getAvatarUrl() { return avatarUrl; }
        public void setAvatarUrl(String avatarUrl) { this.avatarUrl = avatarUrl; }
    }
}
