package com.videoplatform.controller;

import com.videoplatform.dto.auth.LoginRequest;
import com.videoplatform.dto.auth.LoginResponse;
import com.videoplatform.dto.auth.RegisterRequest;
import com.videoplatform.dto.auth.TokenPair;
import com.videoplatform.dto.user.UserProfileResponse;
import com.videoplatform.service.UserService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/v1/auth")
@RequiredArgsConstructor
@Tag(name = "Authentication", description = "Authentication endpoints")
public class AuthController {

    private final UserService userService;

    @PostMapping("/register")
    @Operation(summary = "Register a new user")
    public ResponseEntity<UserProfileResponse> register(@Valid @RequestBody RegisterRequest request) {
        return ResponseEntity.ok(userService.register(request));
    }

    @PostMapping("/login")
    @Operation(summary = "Login user")
    public ResponseEntity<LoginResponse> login(@Valid @RequestBody LoginRequest request) {
        return ResponseEntity.ok(userService.login(request));
    }

    @PostMapping("/refresh")
    @Operation(summary = "Refresh access token")
    public ResponseEntity<TokenPair> refreshToken(@RequestBody TokenPair request) {
        return ResponseEntity.ok(userService.refreshToken(request.getRefreshToken()));
    }

    @PostMapping("/logout")
    @Operation(summary = "Logout user")
    public ResponseEntity<Void> logout(
            @RequestAttribute("userId") Long userId,
            @RequestBody TokenPair request
    ) {
        userService.logout(userId, request.getRefreshToken());
        return ResponseEntity.ok().build();
    }
}
