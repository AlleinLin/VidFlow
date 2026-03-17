<script lang="ts">
  import { cn } from '$lib/utils';
  import { Button } from '$lib/components/ui/button';
  import { Sun, Moon, Upload, User, Menu } from 'lucide-svelte';
  import { toggleMode } from 'mode-watcher';
  import { goto } from '$app/navigation';

  export let isAuthenticated: boolean = false;
  export let className: string = '';

  let mobileMenuOpen = false;

  $: navClasses = cn('border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60', className);

  async function handleLogout() {
    document.cookie = 'access_token=; path=/; max-age=0';
    document.cookie = 'refresh_token=; path=/; max-age=0';
    await goto('/login');
  }
</script>

<nav class={navClasses}>
  <div class="container flex h-16 items-center justify-between">
    <div class="flex items-center gap-6">
      <a href="/" class="flex items-center gap-2">
        <svg class="h-8 w-8 text-primary" viewBox="0 0 24 24" fill="currentColor">
          <path d="M23 7l-7 5 7 5V7z"/>
          <rect x="1" y="5" width="15" height="14" rx="2" ry="2"/>
        </svg>
        <span class="font-bold text-xl">VideoHub</span>
      </a>
      
      <div class="hidden md:flex items-center gap-4">
        <a href="/" class="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors">
          Home
        </a>
        <a href="/trending" class="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors">
          Trending
        </a>
        <a href="/categories" class="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors">
          Categories
        </a>
      </div>
    </div>
    
    <div class="flex items-center gap-2">
      <Button variant="ghost" size="icon" on:click={() => toggleMode()}>
        <Sun class="h-5 w-5 rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
        <Moon class="absolute h-5 w-5 rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
        <span class="sr-only">Toggle theme</span>
      </Button>
      
      {#if isAuthenticated}
        <Button variant="ghost" size="icon" on:click={() => goto('/upload')}>
          <Upload class="h-5 w-5" />
        </Button>
        <Button variant="ghost" size="icon" on:click={() => goto('/profile')}>
          <User class="h-5 w-5" />
        </Button>
        <Button variant="outline" on:click={handleLogout}>
          Sign out
        </Button>
      {:else}
        <Button on:click={() => goto('/login')}>
          Sign in
        </Button>
      {/if}
      
      <Button variant="ghost" size="icon" className="md:hidden" on:click={() => mobileMenuOpen = !mobileMenuOpen}>
        <Menu class="h-5 w-5" />
      </Button>
    </div>
  </div>
  
  {#if mobileMenuOpen}
    <div class="md:hidden border-t py-4">
      <div class="container flex flex-col gap-2">
        <a href="/" class="text-sm font-medium py-2">Home</a>
        <a href="/trending" class="text-sm font-medium py-2">Trending</a>
        <a href="/categories" class="text-sm font-medium py-2">Categories</a>
      </div>
    </div>
  {/if}
</nav>
