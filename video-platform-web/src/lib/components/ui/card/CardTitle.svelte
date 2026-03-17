<script lang="ts">
  import { cn } from '$lib/utils';

  export let className: string = '';

  $: classes = cn(
    'text-2xl font-semibold leading-none tracking-tight',
    className
  );
</script>

<h3 class={classes} {...$$restProps}>
  <slot />
</h3>
