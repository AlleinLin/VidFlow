<script lang="ts">
  import { cn } from '$lib/utils';

  export let className: string = '';

  $: classes = cn('text-sm text-muted-foreground', className);
</script>

<p class={classes} {...$$restProps}>
  <slot />
</p>
