<script lang="ts">
  import { ChevronDown, Layers } from "@lucide/svelte";
  import { agentsStore, chatStore, connectionStore, explorerStore, sessionStore } from "$lib/stores";
  import { defaultExplorerSource } from "$lib/filesystem";

  let open = $state(false);

  let activeSpace = $derived(
    connectionStore.spaces.find((s) => s.id === connectionStore.spaceId) ?? null,
  );

  function toggle() {
    open = !open;
  }

  async function select(id: string) {
    if (id === connectionStore.spaceId) {
      open = false;
      return;
    }
    connectionStore.selectSpace(id);
    open = false;
    chatStore.clearMessages();
    sessionStore.setActiveSession(null);
    await Promise.all([agentsStore.loadAgents(), sessionStore.loadSessions()]);
    await explorerStore.switchSource(defaultExplorerSource());
  }

  function handleClickOutside(e: MouseEvent) {
    const target = e.target as HTMLElement;
    if (!target.closest("[data-space-switcher]")) {
      open = false;
    }
  }
</script>

<svelte:window onclick={handleClickOutside} />

{#if connectionStore.spaces.length > 0}
  <div class="relative" data-space-switcher>
    <button
      type="button"
      onclick={toggle}
      class="flex max-w-[200px] items-center gap-1.5 rounded-md border border-border/60 bg-muted/30 px-2 py-1 text-[11px] text-foreground transition-colors hover:bg-accent"
      aria-haspopup="listbox"
      aria-expanded={open}
    >
      <Layers class="h-3 w-3 shrink-0 text-muted-foreground" strokeWidth={1.75} />
      <span class="truncate">{activeSpace?.name ?? "Select space"}</span>
      <ChevronDown class="h-3 w-3 shrink-0 text-muted-foreground" strokeWidth={1.75} />
    </button>

    {#if open}
      <ul
        class="absolute right-0 top-full z-50 mt-1 min-w-[180px] overflow-hidden rounded-md border border-border bg-popover py-1 shadow-md"
        role="listbox"
      >
        {#each connectionStore.spaces as space (space.id)}
          <li>
            <button
              type="button"
              role="option"
              aria-selected={space.id === connectionStore.spaceId}
              onclick={() => select(space.id)}
              class={[
                "flex w-full flex-col px-3 py-2 text-left text-[11px] transition-colors hover:bg-accent",
                space.id === connectionStore.spaceId ? "bg-accent/50 font-medium" : "",
              ].join(" ")}
            >
              <span class="text-foreground">{space.name}</span>
              {#if space.is_system}
                <span class="text-[10px] text-muted-foreground">system</span>
              {/if}
            </button>
          </li>
        {/each}
      </ul>
    {/if}
  </div>
{/if}