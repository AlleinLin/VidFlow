export interface User {
  id: number;
  username: string;
  email: string;
  display_name: string;
  avatar_url?: string;
  bio?: string;
  role: string;
  status: string;
  follower_count: number;
  following_count: number;
  created_at: string;
}

export interface VideoResponse {
  id: number;
  user_id: number;
  title: string;
  description?: string;
  status: string;
  visibility: string;
  duration_seconds?: number;
  thumbnail_url?: string;
  category_id?: number;
  view_count: number;
  like_count: number;
  comment_count: number;
  published_at?: string;
  created_at: string;
}

export interface VideoListResponse {
  videos: VideoResponse[];
  total: number;
  page: number;
  page_size: number;
}

export interface Comment {
  id: number;
  video_id: number;
  user_id: number;
  parent_id?: number;
  root_id?: number;
  content: string;
  like_count: number;
  status: string;
  created_at: string;
}

export interface TokenPair {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  token_type: string;
}

export interface LoginResponse {
  token: TokenPair;
  user: User;
}

export interface ApiResponse<T> {
  data?: T;
  error?: string;
}
