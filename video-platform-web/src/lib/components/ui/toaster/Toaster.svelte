<script lang="ts">
  import { cn } from '$lib/utils';

  export let className: string = '';

  $: classes = cn('toast-container', className);
</script>

<div class={classes} {...$$restProps}>
  <slot />
</div>
