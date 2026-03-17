package com.videoplatform.security;

import io.jsonwebtoken.*;
import io.jsonwebtoken.security.Keys;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import javax.crypto.SecretKey;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.Date;

@Component
public class JwtTokenProvider {

    private final SecretKey secretKey;
    private final Duration accessTokenTtl;
    private final Duration refreshTokenTtl;
    private final String issuer;

    public JwtTokenProvider(
            @Value("${jwt.secret}") String secret,
            @Value("${jwt.access-token-ttl:15m}") Duration accessTokenTtl,
            @Value("${jwt.refresh-token-ttl:168h}") Duration refreshTokenTtl,
            @Value("${jwt.issuer:video-platform}") String issuer
    ) {
        this.secretKey = Keys.hmacShaKeyFor(secret.getBytes(StandardCharsets.UTF_8));
        this.accessTokenTtl = accessTokenTtl;
        this.refreshTokenTtl = refreshTokenTtl;
        this.issuer = issuer;
    }

    public String generateAccessToken(Long userId, String username, String role) {
        return generateToken(userId, username, role, accessTokenTtl);
    }

    public String generateRefreshToken(Long userId, String username, String role) {
        return generateToken(userId, username, role, refreshTokenTtl);
    }

    private String generateToken(Long userId, String username, String role, Duration ttl) {
        Date now = new Date();
        Date expiryDate = new Date(now.getTime() + ttl.toMillis());

        return Jwts.builder()
                .subject(String.valueOf(userId))
                .claim("userId", userId)
                .claim("username", username)
                .claim("role", role)
                .issuer(issuer)
                .issuedAt(now)
                .expiration(expiryDate)
                .signWith(secretKey)
                .compact();
    }

    public TokenPair generateTokenPair(Long userId, String username, String role) {
        return TokenPair.builder()
                .accessToken(generateAccessToken(userId, username, role))
                .refreshToken(generateRefreshToken(userId, username, role))
                .expiresIn(accessTokenTtl.getSeconds())
                .tokenType("Bearer")
                .build();
    }

    public boolean validateToken(String token) {
        try {
            Jwts.parser()
                    .verifyWith(secretKey)
                    .build()
                    .parseSignedClaims(token);
            return true;
        } catch (JwtException | IllegalArgumentException e) {
            return false;
        }
    }

    public Long getUserIdFromToken(String token) {
        Claims claims = Jwts.parser()
                .verifyWith(secretKey)
                .build()
                .parseSignedClaims(token)
                .getPayload();
        return claims.get("userId", Long.class);
    }

    public String getUsernameFromToken(String token) {
        Claims claims = Jwts.parser()
                .verifyWith(secretKey)
                .build()
                .parseSignedClaims(token)
                .getPayload();
        return claims.get("username", String.class);
    }

    public String getRoleFromToken(String token) {
        Claims claims = Jwts.parser()
                .verifyWith(secretKey)
                .build()
                .parseSignedClaims(token)
                .getPayload();
        return claims.get("role", String.class);
    }

    @lombok.Data
    @lombok.Builder
    public static class TokenPair {
        private String accessToken;
        private String refreshToken;
        private long expiresIn;
        private String tokenType;
    }
}
