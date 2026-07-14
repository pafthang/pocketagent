<script lang="ts">
  import { agentsStore } from "$lib/stores";
  import * as Select from "$lib/components/ui/select";
  import { Bot } from "@lucide/svelte";

  let agents = $derived(agentsStore.agents);
  let selectedId = $derived(agentsStore.selectedAgentId ?? "");
  let selectedLabel = $derived(agentsStore.selectedAgent?.name ?? "Select agent");

  function onValueChange(value: string) {
    agentsStore.selectAgent(value || null);
  }
</script>

{#if agents.length > 0}
  <Select.Root type="single" value={selectedId} onValueChange={onValueChange}>
    <Select.Trigger class="h-7 max-w-[160px] gap-1.5 border-border/60 bg-muted/30 text-[11px]">
      <Bot class="h-3 w-3 shrink-0 text-muted-foreground" />
      <span class="truncate">{selectedLabel}</span>
    </Select.Trigger>
    <Select.Content>
      {#each agents as agent (agent.id)}
        <Select.Item value={agent.id} label={agent.name}>
          <span class="truncate">{agent.name}</span>
          {#if agent.model}
            <span class="ml-1 text-[10px] text-muted-foreground">({agent.model})</span>
          {/if}
        </Select.Item>
      {/each}
    </Select.Content>
  </Select.Root>
{/if}