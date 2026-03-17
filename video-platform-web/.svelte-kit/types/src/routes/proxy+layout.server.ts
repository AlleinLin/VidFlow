// @ts-nocheck
import type { LayoutServerLoad } from './$types';

export const load = async ({ cookies }: Parameters<LayoutServerLoad>[0]) => {
  const token = cookies.get('access_token');
  
  return {
    isAuthenticated: !!token
  };
};
