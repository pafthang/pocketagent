<script lang="ts">
  import { goto } from "$app/navigation";

  function complete() {
    localStorage.setItem("pocketagent_onboarded", "1");
    goto("/");
  }
</script>

<div class="flex max-w-md flex-col items-center gap-6 px-6 text-center">
  <span class="text-4xl">🐾</span>
  <div class="flex flex-col gap-2">
    <h1 class="text-xl font-semibold text-foreground">Welcome to PocketAgent</h1>
    <p class="text-sm text-muted-foreground">
      Web UI for managing AI agents, tasks, and spaces. Connect to your gate backend to get started.
    </p>
  </div>
  <button
    type="button"
    onclick={complete}
    class="rounded-md bg-primary px-6 py-2.5 text-sm font-medium text-primary-foreground hover:opacity-90"
  >
    Get started
  </button>
</div>