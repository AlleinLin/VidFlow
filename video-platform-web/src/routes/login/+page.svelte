<script lang="ts">
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
  import { goto } from '$app/navigation';
  import { authApi } from '$lib/api/auth';

  let email = '';
  let password = '';
  let loading = false;
  let error = '';

  async function handleLogin(e: Event) {
    e.preventDefault();
    loading = true;
    error = '';
    
    try {
      const response = await authApi.login(email, password);
      document.cookie = `access_token=${response.token.access_token}; path=/; max-age=${response.token.expires_in}`;
      document.cookie = `refresh_token=${response.token.refresh_token}; path=/; max-age=604800; httpOnly`;
      await goto('/');
    } catch (err) {
      error = 'Invalid email or password';
    } finally {
      loading = false;
    }
  }
</script>

<svelte:head>
  <title>Sign In - VideoHub</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-50 to-indigo-100 dark:from-gray-900 dark:to-gray-800">
  <Card class="w-full max-w-md">
    <CardHeader class="space-y-1">
      <CardTitle class="text-2xl font-bold text-center">Sign in</CardTitle>
      <CardDescription class="text-center">
        Enter your credentials to access your account
      </CardDescription>
    </CardHeader>
    <CardContent>
      <form on:submit={handleLogin} class="space-y-4">
        {#if error}
          <div class="p-3 text-sm text-red-600 bg-red-50 rounded-md dark:bg-red-900/20 dark:text-red-400">
            {error}
          </div>
        {/if}
        
        <div class="space-y-2">
          <Label for="email">Email</Label>
          <Input
            id="email"
            type="email"
            placeholder="name@example.com"
            bind:value={email}
            required
          />
        </div>
        
        <div class="space-y-2">
          <Label for="password">Password</Label>
          <Input
            id="password"
            type="password"
            bind:value={password}
            required
          />
        </div>
        
        <Button type="submit" class="w-full" disabled={loading}>
          {#if loading}
            <svg class="animate-spin -ml-1 mr-3 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Signing in...
          {:else}
            Sign in
          {/if}
        </Button>
      </form>
      
      <div class="mt-4 text-center text-sm">
        Don't have an account?
        <a href="/register" class="text-primary hover:underline font-medium">
          Sign up
        </a>
      </div>
    </CardContent>
  </Card>
</div>
