<script lang="ts">
  import { PanelLeft } from "@lucide/svelte";
  import ConnectionBadge from "./ConnectionBadge.svelte";
  import AgentProgressBar from "./AgentProgressBar.svelte";
  import SpaceSwitcher from "$lib/components/SpaceSwitcher.svelte";
  import AgentSelector from "$lib/components/chat/AgentSelector.svelte";

  let {
    onToggleSidebar,
    showTabs = true,
  }: {
    onToggleSidebar?: () => void;
    showTabs?: boolean;
  } = $props();
</script>

<header
  class="relative flex w-full shrink-0 items-center border-b border-border/50 px-3"
  style="height: 40px;"
  data-no-select
>
  <div class="flex min-w-0 flex-1 items-center gap-2">
    {#if onToggleSidebar}
      <button
        type="button"
        onclick={onToggleSidebar}
        class="flex h-8 w-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
        aria-label="Toggle sidebar"
      >
        <PanelLeft class="h-4 w-4" strokeWidth={1.75} />
      </button>
    {/if}
    <span class="text-lg leading-none">🐾</span>
    <span class="text-[13px] font-semibold text-foreground">PocketAgent</span>
  </div>

  {#if showTabs}
    <div class="flex items-center gap-2">
      <AgentSelector />
      <SpaceSwitcher />
      <ConnectionBadge />
    </div>
  {/if}

  <AgentProgressBar />
</header>