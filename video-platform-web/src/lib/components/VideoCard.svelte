<script lang="ts">
  import { cn } from '$lib/utils';
  import type { VideoResponse } from '$lib/types';

  export let video: VideoResponse;
  export let className: string = '';

  $: linkClasses = cn('group block', className);

  function formatViewCount(count: number): string {
    if (count >= 1000000) {
      return (count / 1000000).toFixed(1) + 'M';
    } else if (count >= 1000) {
      return (count / 1000).toFixed(1) + 'K';
    }
    return count.toString();
  }

  function timeAgo(date: string): string {
    const seconds = Math.floor((new Date().getTime() - new Date(date).getTime()) / 1000);
    
    if (seconds < 60) return 'just now';
    if (seconds < 3600) return Math.floor(seconds / 60) + ' minutes ago';
    if (seconds < 86400) return Math.floor(seconds / 3600) + ' hours ago';
    if (seconds < 2592000) return Math.floor(seconds / 86400) + ' days ago';
    if (seconds < 31536000) return Math.floor(seconds / 2592000) + ' months ago';
    return Math.floor(seconds / 31536000) + ' years ago';
  }
</script>

<a href="/video/{video.id}" class={linkClasses}>
  <div class="relative aspect-video rounded-lg overflow-hidden bg-muted">
    {#if video.thumbnail_url}
      <img
        src={video.thumbnail_url}
        alt={video.title}
        class="w-full h-full object-cover transition-transform group-hover:scale-105"
        loading="lazy"
      />
    {:else}
      <div class="w-full h-full flex items-center justify-center bg-gradient-to-br from-purple-500/20 to-indigo-500/20">
        <svg class="h-12 w-12 text-muted-foreground" viewBox="0 0 24 24" fill="currentColor">
          <path d="M23 7l-7 5 7 5V7z"/>
          <rect x="1" y="5" width="15" height="14" rx="2" ry="2"/>
        </svg>
      </div>
    {/if}
    
    {#if video.duration_seconds}
      <div class="absolute bottom-2 right-2 bg-black/80 text-white text-xs px-1.5 py-0.5 rounded">
        {Math.floor(video.duration_seconds / 60)}:{String(video.duration_seconds % 60).padStart(2, '0')}
      </div>
    {/if}
  </div>
  
  <div class="mt-3 space-y-1">
    <h3 class="font-medium line-clamp-2 group-hover:text-primary transition-colors">
      {video.title}
    </h3>
    
    <div class="flex items-center gap-2 text-sm text-muted-foreground">
      <span>{formatViewCount(video.view_count)} views</span>
      <span>•</span>
      <span>{timeAgo(video.created_at)}</span>
    </div>
  </div>
</a>
