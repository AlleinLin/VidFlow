import { API_BASE_URL } from '$lib/config';
import type { VideoListResponse, VideoResponse } from '$lib/types';

export const videoApi = {
  async list(token?: string, page = 1, pageSize = 20): Promise<VideoListResponse> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json'
    };
    
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
    
    const response = await fetch(
      `${API_BASE_URL}/api/v1/videos?page=${page}&page_size=${pageSize}`,
      { headers }
    );
    
    if (!response.ok) {
      throw new Error('Failed to fetch videos');
    }
    
    return response.json();
  },
  
  async get(id: number, token?: string): Promise<VideoResponse> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json'
    };
    
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
    
    const response = await fetch(`${API_BASE_URL}/api/v1/videos/${id}`, {
      headers
    });
    
    if (!response.ok) {
      throw new Error('Failed to fetch video');
    }
    
    return response.json();
  },
  
  async search(query: string, page = 1, pageSize = 20): Promise<VideoListResponse> {
    const response = await fetch(
      `${API_BASE_URL}/api/v1/videos?keyword=${encodeURIComponent(query)}&page=${page}&page_size=${pageSize}`
    );
    
    if (!response.ok) {
      throw new Error('Search failed');
    }
    
    return response.json();
  },
  
  async create(token: string, data: { title: string; description?: string; category_id?: number; visibility?: string }): Promise<VideoResponse> {
    const response = await fetch(`${API_BASE_URL}/api/v1/videos`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(data)
    });
    
    if (!response.ok) {
      throw new Error('Failed to create video');
    }
    
    return response.json();
  },
  
  async delete(token: string, id: number): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/api/v1/videos/${id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to delete video');
    }
  },
  
  async publish(token: string, id: number): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/api/v1/videos/${id}/publish`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to publish video');
    }
  }
};
