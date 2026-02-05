<script lang="ts">
  import { auth } from "$lib";
  import { goto } from "$app/navigation";

  let email = "";
  let password = "";
  let fullName = "";
  let error = "";
  let success = "";
  let isLoading = false;

  async function handleRegister() {
    error = "";
    success = "";
    isLoading = true;

    try {
      await auth.register(email, password, fullName);
      success = "Registration successful! Redirecting to login...";
      setTimeout(() => goto("/auth/login"), 1500);
    } catch (err) {
      error = err instanceof Error ? err.message : "Registration failed";
    } finally {
      isLoading = false;
    }
  }
</script>

<svelte:head>
  <title>Register - Sekolah Madrasah</title>
</svelte:head>

<div class="min-h-[calc(100vh-4rem)] flex items-center justify-center px-4">
  <div class="w-full max-w-md">
    <div class="card">
      <h1 class="text-2xl font-bold text-center mb-6">Register</h1>

      {#if error}
        <div
          class="bg-red-500/20 border border-red-500 text-red-400 px-4 py-3 rounded-lg mb-4"
        >
          {error}
        </div>
      {/if}

      {#if success}
        <div
          class="bg-green-500/20 border border-green-500 text-green-400 px-4 py-3 rounded-lg mb-4"
        >
          {success}
        </div>
      {/if}

      <form on:submit|preventDefault={handleRegister} class="space-y-4">
        <div>
          <label
            for="fullName"
            class="block text-sm font-medium text-slate-300 mb-1"
          >
            Full Name
          </label>
          <input
            id="fullName"
            type="text"
            bind:value={fullName}
            class="input-field"
            placeholder="John Doe"
            required
          />
        </div>

        <div>
          <label
            for="email"
            class="block text-sm font-medium text-slate-300 mb-1"
          >
            Email
          </label>
          <input
            id="email"
            type="email"
            bind:value={email}
            class="input-field"
            placeholder="you@example.com"
            required
          />
        </div>

        <div>
          <label
            for="password"
            class="block text-sm font-medium text-slate-300 mb-1"
          >
            Password
          </label>
          <input
            id="password"
            type="password"
            bind:value={password}
            class="input-field"
            placeholder="••••••••"
            minlength="8"
            required
          />
          <p class="text-xs text-slate-500 mt-1">Minimum 8 characters</p>
        </div>

        <button
          type="submit"
          class="btn-primary w-full py-3"
          disabled={isLoading}
        >
          {isLoading ? "Creating account..." : "Register"}
        </button>
      </form>

      <p class="text-center text-slate-400 mt-6">
        Already have an account?
        <a href="/auth/login" class="text-primary-500 hover:underline">Login</a>
      </p>
    </div>
  </div>
</div>
