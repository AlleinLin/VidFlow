// @ts-nocheck
import type { PageServerLoad } from './$types';
import { videoApi } from '$lib/api/video';

export const load = async ({ cookies }: Parameters<PageServerLoad>[0]) => {
  const token = cookies.get('access_token');
  
  try {
    const response = await videoApi.list(token);
    return {
      videos: response.videos,
      isAuthenticated: !!token
    };
  } catch (err) {
    return {
      videos: [],
      isAuthenticated: !!token
    };
  }
};
