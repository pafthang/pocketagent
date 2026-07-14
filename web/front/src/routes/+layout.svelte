<script lang="ts">
  import "../styles/global.css";
  import type { Snippet } from "svelte";
  import { onMount, onDestroy } from "svelte";
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import { ModeWatcher } from "mode-watcher";
  import { Toaster } from "$lib/components/ui/sonner";
  import { Provider as TooltipProvider } from "$lib/components/ui/tooltip";
  import AppShell from "$lib/components/AppShell.svelte";
  import TitleBar from "$lib/components/titlebar/TitleBar.svelte";
  import MobileHeader from "$lib/components/MobileHeader.svelte";
  import {
    initializeStores,
    activityStore,
    connectionStore,
    platformStore,
    sessionStore,
    chatStore,
    settingsStore,
    uiStore,
  } from "$lib/stores";
  import {
    registerHotkeys,
    unregisterHotkeys,
    setupTrayListeners,
    cleanupTrayListeners,
    requestNotificationPermission,
    notifyAgentComplete,
    emitSessionSwitch,
    emitChatSync,
    emitSettingsUpdate,
    onSidePanelReady,
    onChatSync,
    disposeAllBridgeListeners,
  } from "$lib/platform";

  let { children }: { children: Snippet } = $props();

  let isOnboarding = $derived(page.url.pathname.startsWith("/onboarding"));
  let isLogin = $derived(page.url.pathname.startsWith("/login"));

  function handleToggleSidebar() {
    if (platformStore.isDesktop) {
      uiStore.toggleSidebar();
    } else {
      uiStore.toggleDrawer();
    }
  }

  let bootError = $state<string | null>(null);

  let prevWorking = $state(false);
  $effect(() => {
    const working = activityStore.isAgentWorking;
    if (prevWorking && !working) {
      const latest = activityStore.latestEntry;
      notifyAgentComplete(latest?.content ?? "Task completed");
    }
    prevWorking = working;
  });

  let prevSessionId = $state<string | null>(null);
  $effect(() => {
    const sid = sessionStore.activeSessionId;
    if (sid && sid !== prevSessionId && connectionStore.isAuthenticated) {
      emitSessionSwitch(sid);
    }
    prevSessionId = sid;
  });

  let prevStreaming = $state(false);
  $effect(() => {
    const streaming = chatStore.isStreaming;
    if (prevStreaming && !streaming && connectionStore.isAuthenticated) {
      emitChatSync({
        messages: chatStore.messages,
        streaming: false,
        streamingContent: "",
      });
    }
    prevStreaming = streaming;
  });

  $effect(() => {
    const settings = settingsStore.settings;
    if (settings && connectionStore.isAuthenticated) {
      emitSettingsUpdate(settings);
    }
  });

  $effect(() => {
    if (!connectionStore.ready) return;
    if (!connectionStore.isAuthenticated && !isLogin) {
      goto("/login");
    }
  });

  function finishSetup() {
    const pathname = page.url.pathname;
    const onboarded = localStorage.getItem("pocketagent_onboarded");
    if (!onboarded && !pathname.startsWith("/onboarding")) {
      goto("/onboarding");
    }

    requestNotificationPermission();

    if (!platformStore.isNativeMobile) {
      registerHotkeys({});
      setupTrayListeners({
        onNavigate: (path: string) => goto(path),
      });

      onSidePanelReady(() => {
        if (sessionStore.activeSessionId) {
          emitSessionSwitch(sessionStore.activeSessionId);
        }
        emitChatSync({
          messages: chatStore.messages,
          streaming: chatStore.isStreaming,
          streamingContent: chatStore.streamingContent,
        });
        if (settingsStore.settings) {
          emitSettingsUpdate(settingsStore.settings);
        }
      });

      onChatSync((payload) => {
        if (!chatStore.isStreaming) {
          chatStore.loadHistory(payload.messages);
        }
      });
    }
  }

  async function bootstrap() {
    bootError = null;
    try {
      const authenticated = await connectionStore.bootstrap();
      if (authenticated) {
        await initializeStores();
        finishSetup();
      }
    } catch (err) {
      bootError = err instanceof Error ? err.message : "Failed to start";
      console.error("[layout] bootstrap failed:", err);
    }
  }

  onMount(() => {
    bootstrap();
  });

  onDestroy(() => {
    unregisterHotkeys();
    cleanupTrayListeners();
    disposeAllBridgeListeners();
  });
</script>

<ModeWatcher />
<div>
  <TooltipProvider>
    {#if isLogin}
      {@render children()}
    {:else if !connectionStore.ready || connectionStore.loading}
      <div class="flex h-dvh items-center justify-center bg-background">
        <div class="flex flex-col items-center gap-3">
          <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary border-t-transparent"></div>
          <p class="text-sm text-muted-foreground">Loading...</p>
        </div>
      </div>
    {:else if bootError}
      <div class="flex h-dvh items-center justify-center bg-background">
        <div class="flex flex-col items-center gap-4 text-center">
          <p class="text-sm text-destructive">{bootError}</p>
          <button
            onclick={bootstrap}
            class="rounded-md bg-primary px-4 py-2 text-sm text-primary-foreground hover:opacity-90"
          >
            Retry
          </button>
        </div>
      </div>
    {:else if !connectionStore.isAuthenticated}
      <!-- redirect to /login via $effect -->
      <div class="flex h-dvh items-center justify-center bg-background">
        <p class="text-sm text-muted-foreground">Redirecting to sign in...</p>
      </div>
    {:else}
      <div class="flex h-dvh w-screen flex-col bg-background">
        {#if platformStore.isNativeMobile}
          <MobileHeader />
        {:else}
          {@const isAppReady = !isOnboarding}
          <TitleBar onToggleSidebar={isAppReady ? handleToggleSidebar : undefined} showTabs={isAppReady} />
        {/if}

        {#if isOnboarding}
          <div class="flex flex-1 items-center justify-center overflow-hidden">
            {@render children()}
          </div>
        {:else}
          <AppShell>
            {@render children()}
          </AppShell>
        {/if}
      </div>
    {/if}
    {#if !isLogin}
      <Toaster />
    {/if}
  </TooltipProvider>
</div>