<script lang="ts">
  import { goto } from "$app/navigation";
  import { connectionStore, initializeStores } from "$lib/stores";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Loader2 } from "@lucide/svelte";

  let mode = $state<"login" | "register">("login");
  let email = $state("");
  let password = $state("");
  let confirmPassword = $state("");
  let localError = $state<string | null>(null);
  let successMessage = $state<string | null>(null);
  let submitting = $state(false);

  async function handleSubmit(e: Event) {
    e.preventDefault();
    localError = null;
    successMessage = null;

    if (!email.trim() || !password) {
      localError = "Email and password are required.";
      return;
    }

    if (mode === "register" && password !== confirmPassword) {
      localError = "Passwords do not match.";
      return;
    }

    submitting = true;
    try {
      if (mode === "login") {
        await connectionStore.login(email.trim(), password);
        await initializeStores();
        const onboarded = localStorage.getItem("pocketagent_onboarded");
        goto(onboarded ? "/" : "/onboarding");
      } else {
        const res = await connectionStore.register(email.trim(), password);
        if (res.email_verification_required) {
          successMessage = "Account created. Check your email to verify, then sign in.";
          mode = "login";
          password = "";
          confirmPassword = "";
        } else {
          await connectionStore.login(email.trim(), password);
          await initializeStores();
          goto("/onboarding");
        }
      }
    } catch (err) {
      localError = connectionStore.error ?? (err instanceof Error ? err.message : "Request failed");
    } finally {
      submitting = false;
    }
  }
</script>

<div class="flex min-h-dvh items-center justify-center bg-background px-4">
  <div class="w-full max-w-sm">
    <div class="mb-8 flex flex-col items-center gap-2 text-center">
      <span class="text-4xl">🐾</span>
      <h1 class="text-xl font-semibold text-foreground">PocketAgent</h1>
      <p class="text-sm text-muted-foreground">
        {mode === "login" ? "Sign in to your workspace" : "Create an account"}
      </p>
    </div>

    <form onsubmit={handleSubmit} class="flex flex-col gap-4">
      <div class="flex flex-col gap-1.5">
        <label for="email" class="text-xs font-medium text-muted-foreground">Email</label>
        <Input
          id="email"
          type="email"
          autocomplete="email"
          bind:value={email}
          placeholder="you@example.com"
          disabled={submitting}
        />
      </div>

      <div class="flex flex-col gap-1.5">
        <label for="password" class="text-xs font-medium text-muted-foreground">Password</label>
        <Input
          id="password"
          type="password"
          autocomplete={mode === "login" ? "current-password" : "new-password"}
          bind:value={password}
          disabled={submitting}
        />
      </div>

      {#if mode === "register"}
        <div class="flex flex-col gap-1.5">
          <label for="confirm" class="text-xs font-medium text-muted-foreground">Confirm password</label>
          <Input
            id="confirm"
            type="password"
            autocomplete="new-password"
            bind:value={confirmPassword}
            disabled={submitting}
          />
        </div>
      {/if}

      {#if localError}
        <p class="text-sm text-destructive">{localError}</p>
      {/if}
      {#if successMessage}
        <p class="text-sm text-emerald-600">{successMessage}</p>
      {/if}

      <Button type="submit" class="w-full" disabled={submitting}>
        {#if submitting}
          <Loader2 class="mr-2 h-4 w-4 animate-spin" />
        {/if}
        {mode === "login" ? "Sign in" : "Create account"}
      </Button>
    </form>

    <p class="mt-6 text-center text-xs text-muted-foreground">
      {#if mode === "login"}
        No account?
        <button
          type="button"
          class="text-primary hover:underline"
          onclick={() => {
            mode = "register";
            localError = null;
          }}
        >
          Register
        </button>
      {:else}
        Already have an account?
        <button
          type="button"
          class="text-primary hover:underline"
          onclick={() => {
            mode = "login";
            localError = null;
          }}
        >
          Sign in
        </button>
      {/if}
    </p>

    <p class="mt-4 text-center text-[10px] text-muted-foreground">
      API: {connectionStore.backendUrl}
    </p>
  </div>
</div>