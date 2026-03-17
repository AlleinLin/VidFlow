<script lang="ts">
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import VideoCard from '$lib/components/VideoCard.svelte';
  import Navigation from '$lib/components/Navigation.svelte';
  import { videoApi } from '$lib/api/video';
  import type { VideoResponse } from '$lib/types';

  export let data;

  let searchQuery = '';
  let videos: VideoResponse[] = data.videos || [];
  let loading = false;

  async function handleSearch() {
    if (!searchQuery.trim()) return;
    
    loading = true;
    try {
      const response = await videoApi.search(searchQuery);
      videos = response.videos;
    } catch (err) {
      console.error('Search failed:', err);
    } finally {
      loading = false;
    }
  }

  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Enter') {
      handleSearch();
    }
  }
</script>

<svelte:head>
  <title>VideoHub - Discover Amazing Videos</title>
</svelte:head>

<div class="min-h-screen bg-background">
  <Navigation isAuthenticated={data.isAuthenticated} />
  
  <main class="container py-8">
    <section class="mb-12">
      <div class="text-center mb-8">
        <h1 class="text-4xl font-bold tracking-tight mb-4">
          Discover Amazing Videos
        </h1>
        <p class="text-muted-foreground text-lg">
          Watch, create, and share content with millions of users worldwide
        </p>
      </div>
      
      <div class="max-w-xl mx-auto">
        <div class="flex gap-2">
          <Input
            type="search"
            placeholder="Search videos..."
            bind:value={searchQuery}
            on:keydown={handleKeydown}
            className="flex-1"
          />
          <Button on:click={handleSearch} disabled={loading}>
            Search
          </Button>
        </div>
      </div>
    </section>
    
    <section>
      <h2 class="text-2xl font-semibold mb-6">Trending Videos</h2>
      
      {#if videos.length === 0}
        <div class="text-center py-12 text-muted-foreground">
          No videos found. Be the first to upload!
        </div>
      {:else}
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {#each videos as video (video.id)}
            <VideoCard {video} />
          {/each}
        </div>
      {/if}
    </section>
  </main>
</div>
