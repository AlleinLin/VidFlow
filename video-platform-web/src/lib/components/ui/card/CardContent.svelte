<script lang="ts">
  import { cn } from '$lib/utils';

  export let className: string = '';

  $: classes = cn('p-6 pt-0', className);
</script>

<div class={classes} {...$$restProps}>
  <slot />
</div>
