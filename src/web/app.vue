<template>
  <div id="app">
    <nav class="app-nav" v-if="mounted">
      <div class="nav-inner">
        <NuxtLink to="/" class="nav-brand">Sàn Sản phẩm số</NuxtLink>
        <div class="nav-links">
          <NuxtLink to="/" class="nav-link">Danh mục</NuxtLink>
          <template v-if="isLoggedIn">
            <span class="nav-user">{{ state.account?.ten || state.account?.email }}</span>
            <button @click="handleLogout" class="nav-link nav-btn">Đăng xuất</button>
          </template>
          <template v-else>
            <NuxtLink to="/dang-nhap" class="nav-link nav-btn">Đăng nhập</NuxtLink>
            <NuxtLink to="/dang-ky" class="nav-link nav-btn nav-btn--primary">Đăng ký</NuxtLink>
          </template>
        </div>
      </div>
    </nav>
    <NuxtPage />
  </div>
</template>

<script setup lang="ts">
const { isLoggedIn, state, dangXuat } = useAuth()
const mounted = ref(false)

onMounted(() => {
  mounted.value = true
})

async function handleLogout() {
  await dangXuat()
  navigateTo('/')
}
</script>

<style>
:root {
  --bg-primary: #fff;
  --bg-secondary: #f5f5f5;
  --text-primary: #1a1a1a;
  --text-secondary: #666;
  --border: #d1d5db;
  --accent: #2563eb;
  --accent-hover: #1d4ed8;
}

* {
  box-sizing: border-box;
}

body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  color: var(--text-primary);
  background: var(--bg-secondary);
}

a {
  color: var(--accent);
  text-decoration: none;
}
</style>

<style scoped>
.app-nav {
  background: var(--bg-primary);
  border-bottom: 1px solid var(--border);
  position: sticky;
  top: 0;
  z-index: 100;
}

.nav-inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nav-brand {
  font-size: 1.125rem;
  font-weight: 700;
  color: var(--text-primary);
  text-decoration: none;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.nav-link {
  font-size: 0.875rem;
  color: var(--text-secondary);
  text-decoration: none;
  padding: 0.375rem 0.75rem;
  border-radius: 6px;
  transition: background 0.15s, color 0.15s;
}

.nav-link:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.nav-btn {
  background: transparent;
  border: 1px solid var(--border);
  cursor: pointer;
  font-family: inherit;
}

.nav-btn--primary {
  background: var(--accent);
  color: #fff;
  border-color: var(--accent);
}

.nav-btn--primary:hover {
  background: var(--accent-hover);
  color: #fff;
}

.nav-user {
  font-size: 0.875rem;
  color: var(--text-primary);
  font-weight: 600;
}
</style>
