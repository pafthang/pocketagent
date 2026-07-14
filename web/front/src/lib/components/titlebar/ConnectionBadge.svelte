<script lang="ts">
  import { connectionStore } from "$lib/stores";

  let label = $derived(
    connectionStore.isOffline
      ? "Offline"
      : connectionStore.isConnected
        ? "Connected"
        : "Disconnected",
  );

  let colorClass = $derived(
    connectionStore.isOffline
      ? "bg-amber-500/20 text-amber-600"
      : connectionStore.isConnected
        ? "bg-emerald-500/20 text-emerald-600"
        : "bg-muted text-muted-foreground",
  );
</script>

<span class="rounded-full px-2 py-0.5 text-[10px] font-medium {colorClass}">
  {label}
</span>