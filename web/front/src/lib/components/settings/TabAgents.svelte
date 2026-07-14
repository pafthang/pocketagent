<script lang="ts">
  import type { Agent } from "@pocketagent/client";
  import { agentsStore } from "$lib/stores";
  import { Button } from "$lib/components/ui/button";
  import * as Dialog from "$lib/components/ui/dialog";
  import { Loader2, Plus, Trash2, Pencil, Bot } from "@lucide/svelte";

  let agents = $derived(agentsStore.agents);
  let isLoading = $derived(agentsStore.isLoading);

  let dialogOpen = $state(false);
  let editingAgent = $state<Agent | null>(null);
  let name = $state("");
  let description = $state("");
  let model = $state("");
  let systemPrompt = $state("");
  let saving = $state(false);
  let confirmDelete = $state(false);

  function openCreate() {
    editingAgent = null;
    name = "";
    description = "";
    model = "";
    systemPrompt = "";
    confirmDelete = false;
    dialogOpen = true;
  }

  function openEdit(agent: Agent) {
    editingAgent = agent;
    name = agent.name;
    description = agent.description ?? "";
    model = agent.model ?? "";
    systemPrompt = agent.system_prompt ?? "";
    confirmDelete = false;
    dialogOpen = true;
  }

  async function handleSave() {
    if (!name.trim()) return;
    saving = true;
    try {
      const input = {
        name: name.trim(),
        description: description.trim() || undefined,
        model: model.trim() || undefined,
        system_prompt: systemPrompt.trim() || undefined,
      };
      if (editingAgent) {
        await agentsStore.updateAgent(editingAgent.id, input);
      } else {
        await agentsStore.createAgent(input);
      }
      dialogOpen = false;
    } finally {
      saving = false;
    }
  }

  async function handleDelete() {
    if (!editingAgent) return;
    saving = true;
    try {
      const ok = await agentsStore.deleteAgent(editingAgent.id);
      if (ok) dialogOpen = false;
    } finally {
      saving = false;
    }
  }
</script>

<div class="flex h-full flex-col gap-4 overflow-hidden">
  <div class="flex items-center justify-between gap-3">
    <div>
      <h2 class="text-sm font-semibold text-foreground">Agents</h2>
      <p class="text-xs text-muted-foreground">
        Configure agents for your space. The selected agent is used in chat.
      </p>
    </div>
    <Button size="sm" onclick={openCreate}>
      <Plus class="mr-1 h-3.5 w-3.5" />
      New agent
    </Button>
  </div>

  {#if isLoading}
    <div class="flex items-center gap-2 text-sm text-muted-foreground">
      <Loader2 class="h-4 w-4 animate-spin" />
      Loading agents…
    </div>
  {:else if agents.length === 0}
    <div class="rounded-lg border border-dashed border-border p-8 text-center">
      <Bot class="mx-auto mb-2 h-8 w-8 text-muted-foreground" />
      <p class="text-sm text-muted-foreground">No agents yet. Create one to start chatting.</p>
      <Button class="mt-4" size="sm" onclick={openCreate}>Create agent</Button>
    </div>
  {:else}
    <div class="flex-1 space-y-2 overflow-y-auto">
      {#each agents as agent (agent.id)}
        <div
          class={[
            "flex w-full items-start gap-3 rounded-lg border px-3 py-3 transition-colors",
            agentsStore.selectedAgentId === agent.id ? "border-primary/40 bg-accent/30" : "border-border",
          ].join(" ")}
        >
          <button
            type="button"
            onclick={() => agentsStore.selectAgent(agent.id)}
            class="flex min-w-0 flex-1 items-start gap-3 text-left hover:opacity-90"
          >
            <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-muted">
              <Bot class="h-4 w-4 text-muted-foreground" />
            </div>
            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <span class="truncate text-sm font-medium">{agent.name}</span>
                {#if agentsStore.selectedAgentId === agent.id}
                  <span class="rounded bg-primary/10 px-1.5 py-0.5 text-[10px] text-primary">active</span>
                {/if}
              </div>
              {#if agent.model}
                <p class="truncate text-xs text-muted-foreground">{agent.model}</p>
              {/if}
              {#if agent.description}
                <p class="mt-0.5 line-clamp-2 text-xs text-muted-foreground">{agent.description}</p>
              {/if}
            </div>
          </button>
          <button
            type="button"
            onclick={() => openEdit(agent)}
            class="mt-1 rounded p-1 text-muted-foreground hover:bg-accent hover:text-foreground"
            aria-label="Edit agent"
          >
            <Pencil class="h-3.5 w-3.5" />
          </button>
        </div>
      {/each}
    </div>
  {/if}
</div>

<Dialog.Root bind:open={dialogOpen}>
  <Dialog.Content class="max-w-md">
    <Dialog.Header>
      <Dialog.Title>{editingAgent ? "Edit agent" : "New agent"}</Dialog.Title>
    </Dialog.Header>

    <div class="space-y-3 py-2">
      <label class="block space-y-1">
        <span class="text-xs font-medium text-muted-foreground">Name</span>
        <input
          bind:value={name}
          class="w-full rounded-md border border-border bg-background px-3 py-2 text-sm"
          placeholder="Research assistant"
        />
      </label>
      <label class="block space-y-1">
        <span class="text-xs font-medium text-muted-foreground">Model</span>
        <input
          bind:value={model}
          class="w-full rounded-md border border-border bg-background px-3 py-2 text-sm"
          placeholder="llama3.1"
        />
      </label>
      <label class="block space-y-1">
        <span class="text-xs font-medium text-muted-foreground">Description</span>
        <input
          bind:value={description}
          class="w-full rounded-md border border-border bg-background px-3 py-2 text-sm"
          placeholder="Optional"
        />
      </label>
      <label class="block space-y-1">
        <span class="text-xs font-medium text-muted-foreground">System prompt</span>
        <textarea
          bind:value={systemPrompt}
          rows={4}
          class="w-full resize-none rounded-md border border-border bg-background px-3 py-2 text-sm"
          placeholder="You are a helpful assistant…"
        ></textarea>
      </label>
    </div>

    <Dialog.Footer class="flex items-center justify-between gap-2">
      {#if editingAgent}
        <div>
          {#if confirmDelete}
            <Button variant="destructive" size="sm" disabled={saving} onclick={handleDelete}>
              Confirm delete
            </Button>
          {:else}
            <Button variant="ghost" size="sm" onclick={() => (confirmDelete = true)}>
              <Trash2 class="mr-1 h-3.5 w-3.5" />
              Delete
            </Button>
          {/if}
        </div>
      {:else}
        <div></div>
      {/if}
      <div class="flex gap-2">
        <Button variant="outline" onclick={() => (dialogOpen = false)}>Cancel</Button>
        <Button disabled={saving || !name.trim()} onclick={handleSave}>
          {#if saving}<Loader2 class="mr-1 h-3.5 w-3.5 animate-spin" />{/if}
          Save
        </Button>
      </div>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>