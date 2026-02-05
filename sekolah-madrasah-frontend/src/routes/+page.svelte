<script lang="ts">
  import { isAuthenticated } from "$lib";
</script>

<svelte:head>
  <title>Sekolah Madrasah - School Management System</title>
</svelte:head>

<div class="min-h-[calc(100vh-4rem)] flex flex-col">
  <!-- Hero Section -->
  <section class="flex-1 flex items-center justify-center px-4">
    <div class="max-w-4xl mx-auto text-center">
      <h1
        class="text-5xl md:text-6xl font-bold mb-6 bg-gradient-to-r from-primary-400 to-primary-600 bg-clip-text text-transparent"
      >
        Sekolah Madrasah
      </h1>
      <p class="text-xl text-slate-400 mb-8 max-w-2xl mx-auto">
        Platform modern untuk manajemen sekolah dan madrasah. Kelola data siswa, 
        kurikulum, nilai, dan administrasi dengan mudah.
      </p>

      <div class="flex flex-col sm:flex-row gap-4 justify-center">
        {#if $isAuthenticated}
          <a href="/dashboard" class="btn-primary text-lg px-8 py-3">
            Go to Dashboard
          </a>
        {:else}
          <a href="/auth/register" class="btn-primary text-lg px-8 py-3">
            Daftar Sekarang
          </a>
          <a href="/auth/login" class="btn-secondary text-lg px-8 py-3">
            Login
          </a>
        {/if}
      </div>
    </div>
  </section>

  <!-- Features Section -->
  <section class="py-16 px-4 bg-slate-800/50">
    <div class="max-w-6xl mx-auto">
      <h2 class="text-3xl font-bold text-center mb-12">Fitur Utama</h2>
      <div class="grid md:grid-cols-3 gap-8">
        <div class="card text-center">
          <div class="text-4xl mb-4">🏫</div>
          <h3 class="text-xl font-semibold mb-2">Manajemen Sekolah</h3>
          <p class="text-slate-400">
            Kelola data sekolah dan unit pendidikan secara terpusat.
          </p>
        </div>
        <div class="card text-center">
          <div class="text-4xl mb-4">👥</div>
          <h3 class="text-xl font-semibold mb-2">Guru & Staff</h3>
          <p class="text-slate-400">
            Manajemen data guru, staff, dan akses pengguna dengan role.
          </p>
        </div>
        <div class="card text-center">
          <div class="text-4xl mb-4">🎓</div>
          <h3 class="text-xl font-semibold mb-2">Siswa</h3>
          <p class="text-slate-400">
            Pencatatan data siswa yang lengkap dan terintegrasi.
          </p>
        </div>
      </div>
    </div>
  </section>
</div>

