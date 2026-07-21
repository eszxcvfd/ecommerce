<template>
  <div class="auth-page">
    <div class="auth-card">
      <h1 class="auth-title">Đăng ký</h1>
      <p class="auth-subtitle">Tạo tài khoản Sàn Sản phẩm số</p>

      <form @submit.prevent="handleDangKy" class="auth-form">
        <div class="form-field">
          <label for="ten">Họ tên</label>
          <input id="ten" v-model="ten" type="text" placeholder="Nguyễn Văn A" required />
        </div>

        <div class="form-field">
          <label for="email">Email</label>
          <input id="email" v-model="email" type="email" placeholder="your@email.com" required autocomplete="email" />
        </div>

        <div class="form-field">
          <label for="password">Mật khẩu</label>
          <input id="password" v-model="password" type="password" placeholder="••••••••" required autocomplete="new-password" />
        </div>

        <p v-if="errorMsg" class="auth-error">{{ errorMsg }}</p>

        <button type="submit" class="auth-btn" :disabled="loading">
          {{ loading ? 'Đang đăng ký...' : 'Đăng ký' }}
        </button>
      </form>

      <p class="auth-link">
        Đã có tài khoản?
        <NuxtLink to="/dang-nhap">Đăng nhập</NuxtLink>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
const { dangKy, isLoggedIn } = useAuth()

const ten = ref('')
const email = ref('')
const password = ref('')
const errorMsg = ref('')
const loading = ref(false)

if (isLoggedIn.value) {
  navigateTo('/')
}

async function handleDangKy() {
  errorMsg.value = ''
  loading.value = true
  try {
    await dangKy({ email: email.value, password: password.value, ten: ten.value })
    navigateTo('/dang-nhap')
  } catch (e: any) {
    errorMsg.value = e.message || 'Đăng ký thất bại'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-secondary, #f5f5f5);
  padding: 1rem;
}

.auth-card {
  background: var(--bg-primary, #fff);
  border-radius: 12px;
  box-shadow: 0 2px 16px rgba(0, 0, 0, 0.08);
  padding: 2.5rem;
  width: 100%;
  max-width: 400px;
}

.auth-title {
  font-size: 1.5rem;
  font-weight: 700;
  margin: 0 0 0.25rem;
  color: var(--text-primary, #1a1a1a);
}

.auth-subtitle {
  font-size: 0.875rem;
  color: var(--text-secondary, #666);
  margin: 0 0 2rem;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
}

.form-field label {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-primary, #1a1a1a);
}

.form-field input {
  padding: 0.625rem 0.75rem;
  border: 1px solid var(--border, #d1d5db);
  border-radius: 8px;
  font-size: 1rem;
  transition: border-color 0.2s;
}

.form-field input:focus {
  outline: none;
  border-color: var(--accent, #2563eb);
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.auth-error {
  color: #dc2626;
  font-size: 0.875rem;
  margin: 0;
  padding: 0.5rem 0.75rem;
  background: #fef2f2;
  border-radius: 6px;
}

.auth-btn {
  padding: 0.75rem;
  background: var(--accent, #2563eb);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}

.auth-btn:hover:not(:disabled) {
  background: var(--accent-hover, #1d4ed8);
}

.auth-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.auth-link {
  text-align: center;
  margin: 1.5rem 0 0;
  font-size: 0.875rem;
  color: var(--text-secondary, #666);
}

.auth-link a {
  color: var(--accent, #2563eb);
  text-decoration: none;
  font-weight: 600;
}

.auth-link a:hover {
  text-decoration: underline;
}
</style>
