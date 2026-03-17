<script lang="ts">
  import { cn } from '$lib/utils';

  export let className: string = '';

  $: classes = cn(
    'text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70',
    className
  );
</script>

<label
  class={classes}
  {...$$restProps}
>
  <slot />
</label>
