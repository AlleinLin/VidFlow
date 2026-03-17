-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    avatar_url VARCHAR(500),
    bio TEXT,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    follower_count BIGINT NOT NULL DEFAULT 0,
    following_count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP
);

-- Create indexes for users
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);

-- Create user_follows table
CREATE TABLE IF NOT EXISTS user_follows (
    follower_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (follower_id, following_id)
);

-- Create videos table
CREATE TABLE IF NOT EXISTS videos (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'UPLOADING',
    visibility VARCHAR(20) NOT NULL DEFAULT 'public',
    duration_seconds INTEGER NOT NULL DEFAULT 0,
    original_filename VARCHAR(255),
    storage_key VARCHAR(500),
    thumbnail_url VARCHAR(500),
    category_id INTEGER,
    view_count BIGINT NOT NULL DEFAULT 0,
    like_count BIGINT NOT NULL DEFAULT 0,
    comment_count INTEGER NOT NULL DEFAULT 0,
    published_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for videos
CREATE INDEX idx_videos_user_id ON videos(user_id);
CREATE INDEX idx_videos_status ON videos(status);
CREATE INDEX idx_videos_visibility ON videos(visibility);
CREATE INDEX idx_videos_category_id ON videos(category_id);
CREATE INDEX idx_videos_created_at ON videos(created_at DESC);
CREATE INDEX idx_videos_view_count ON videos(view_count DESC);

-- Create video_transcode_tasks table
CREATE TABLE IF NOT EXISTS video_transcode_tasks (
    id BIGSERIAL PRIMARY KEY,
    video_id BIGINT NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    resolution VARCHAR(20) NOT NULL,
    format VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    output_storage_key VARCHAR(500),
    error_message TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for transcode tasks
CREATE INDEX idx_transcode_tasks_video_id ON video_transcode_tasks(video_id);
CREATE INDEX idx_transcode_tasks_status ON video_transcode_tasks(status);

-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    parent_id INTEGER REFERENCES categories(id)
);

-- Create video_tags table
CREATE TABLE IF NOT EXISTS video_tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

-- Create video_tag_relations table
CREATE TABLE IF NOT EXISTS video_tag_relations (
    video_id BIGINT NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES video_tags(id) ON DELETE CASCADE,
    PRIMARY KEY (video_id, tag_id)
);

-- Create comments table
CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    video_id BIGINT NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id BIGINT REFERENCES comments(id) ON DELETE CASCADE,
    root_id BIGINT REFERENCES comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    like_count INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'visible',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for comments
CREATE INDEX idx_comments_video_id ON comments(video_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_root_id ON comments(root_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_id);

-- Create likes table
CREATE TABLE IF NOT EXISTS likes (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    video_id BIGINT NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, video_id)
);

-- Create favorites table
CREATE TABLE IF NOT EXISTS favorites (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    video_id BIGINT NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, video_id)
);

-- Create danmakus table
CREATE TABLE IF NOT EXISTS danmakus (
    id BIGSERIAL PRIMARY KEY,
    video_id BIGINT NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content VARCHAR(100) NOT NULL,
    position_seconds DOUBLE PRECISION NOT NULL,
    style VARCHAR(20) NOT NULL DEFAULT 'scroll',
    color VARCHAR(7) NOT NULL DEFAULT '#FFFFFF',
    font_size INTEGER NOT NULL DEFAULT 24,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for danmakus
CREATE INDEX idx_danmakus_video_id ON danmakus(video_id);
CREATE INDEX idx_danmakus_position ON danmakus(video_id, position_seconds);

-- Create watch_history table
CREATE TABLE IF NOT EXISTS watch_history (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    video_id BIGINT NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    progress_seconds INTEGER NOT NULL DEFAULT 0,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, video_id)
);

-- Create indexes for watch history
CREATE INDEX idx_watch_history_user_id ON watch_history(user_id);
CREATE INDEX idx_watch_history_updated_at ON watch_history(updated_at DESC);

-- Insert default categories
INSERT INTO categories (name, description) VALUES
    ('电影', '电影相关视频'),
    ('电视剧', '电视剧相关视频'),
    ('综艺', '综艺节目'),
    ('动漫', '动画和漫画相关'),
    ('游戏', '游戏相关视频'),
    ('音乐', '音乐相关视频'),
    ('教育', '教育类视频'),
    ('科技', '科技类视频'),
    ('生活', '生活类视频'),
    ('娱乐', '娱乐类视频')
ON CONFLICT (name) DO NOTHING;
