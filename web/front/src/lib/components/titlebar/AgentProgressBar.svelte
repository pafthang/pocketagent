<script lang="ts">
  import { activityStore } from "$lib/stores";
</script>

{#if activityStore.isAgentWorking}
  <div
    class="absolute bottom-0 left-0 right-0 h-0.5 overflow-hidden bg-primary/20"
    role="progressbar"
    aria-label="Agent working"
  >
    <div class="h-full w-1/3 animate-pulse bg-primary"></div>
  </div>
{/if}