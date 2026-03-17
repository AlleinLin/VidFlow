-- Add migration script here
-- Initial schema migration

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    avatar_url VARCHAR(500),
    bio TEXT,
    role VARCHAR(20) NOT NULL DEFAULT 'User',
    status VARCHAR(20) NOT NULL DEFAULT 'Active',
    follower_count BIGINT NOT NULL DEFAULT 0,
    following_count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMPTZ
);

-- Videos table
CREATE TABLE IF NOT EXISTS videos (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'Uploading',
    visibility VARCHAR(20) NOT NULL DEFAULT 'Public',
    duration_seconds INTEGER,
    original_filename VARCHAR(255),
    storage_key VARCHAR(500),
    thumbnail_url VARCHAR(500),
    category_id BIGINT,
    view_count BIGINT NOT NULL DEFAULT 0,
    like_count BIGINT NOT NULL DEFAULT 0,
    comment_count INTEGER NOT NULL DEFAULT 0,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Comments table
CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    video_id BIGINT NOT NULL REFERENCES videos(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    parent_id BIGINT REFERENCES comments(id),
    root_id BIGINT REFERENCES comments(id),
    content TEXT NOT NULL,
    like_count INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'visible',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Likes table
CREATE TABLE IF NOT EXISTS likes (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    target_id BIGINT NOT NULL,
    type VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, target_id, type)
);

-- Favorites table
CREATE TABLE IF NOT EXISTS favorites (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    video_id BIGINT NOT NULL REFERENCES videos(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, video_id)
);

-- User follows table
CREATE TABLE IF NOT EXISTS user_follows (
    id BIGSERIAL PRIMARY KEY,
    follower_id BIGINT NOT NULL REFERENCES users(id),
    following_id BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(follower_id, following_id)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_videos_user_id ON videos(user_id);
CREATE INDEX IF NOT EXISTS idx_videos_status ON videos(status);
CREATE INDEX IF NOT EXISTS idx_videos_created_at ON videos(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_comments_video_id ON comments(video_id);
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);
CREATE INDEX IF NOT EXISTS idx_likes_user_target ON likes(user_id, target_id);
CREATE INDEX IF NOT EXISTS idx_favorites_user_id ON favorites(user_id);
