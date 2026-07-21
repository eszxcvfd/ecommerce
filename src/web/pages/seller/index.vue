<template>
  <div class="seller-page">
    <header class="seller-page__header">
      <h1>Quản lý Sản phẩm nháp</h1>
      <NuxtLink to="/seller/tao-moi" class="btn btn--primary">Tạo bản nháp mới</NuxtLink>
    </header>

    <!-- Loading state -->
    <div v-if="loading" class="seller-page__loading">Đang tải...</div>

    <!-- Error state -->
    <div v-else-if="error" class="seller-page__error">
      <p>{{ error }}</p>
      <button @click="loadDrafts" class="btn">Thử lại</button>
    </div>

    <!-- Empty state -->
    <div v-else-if="drafts.length === 0" class="seller-page__empty">
      <p>Chưa có bản nháp nào. Hãy tạo bản nháp đầu tiên!</p>
    </div>

    <!-- Drafts list -->
    <div v-else class="seller-page__list">
      <div v-for="draft in drafts" :key="draft.id" class="draft-card">
        <div class="draft-card__info">
          <h3 class="draft-card__title">{{ draft.ten }}</h3>
          <p class="draft-card__meta">
            <span>Danh mục: {{ draft.danh_muc }}</span>
            <span v-if="draft.gia.mien_phi">Miễn phí</span>
            <span v-else>{{ draft.gia.so_xu?.toLocaleString() }} Xu</span>
            <span>{{ draft.tep?.length || 0 }} tệp</span>
          </p>
          <p v-if="draft.mo_ta" class="draft-card__desc">{{ draft.mo_ta }}</p>
        </div>
        <div class="draft-card__actions">
          <NuxtLink :to="`/seller/sua/${draft.id}`" class="btn btn--small">Sửa</NuxtLink>
          <button @click="handleDelete(draft.id)" class="btn btn--small btn--danger">Xóa</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { SanPhamSo } from '~/composables/useCatalog'
const router = useRouter()
const { listDrafts, deleteDraft } = useSeller()

const drafts = ref<SanPhamSo[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

onMounted(() => {
  loadDrafts()
})

async function loadDrafts() {
  loading.value = true
  error.value = null
  try {
    drafts.value = await listDrafts()
  } catch (e: unknown) {
    const fetchErr = e as { response?: { status?: number }; data?: { message?: string } }
    if (fetchErr?.response?.status === 401) {
      router.push('/dang-nhap')
      return
    }
    error.value = 'Không thể tải danh sách bản nháp'
  } finally {
    loading.value = false
  }
}

async function handleDelete(id: string) {
  if (!confirm('Xóa bản nháp này?')) return
  try {
    await deleteDraft(id)
    drafts.value = drafts.value.filter((d) => d.id !== id)
  } catch {
    alert('Không thể xóa bản nháp')
  }
}
</script>

<style scoped>
.seller-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem;
}

.seller-page__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.seller-page__header h1 {
  margin: 0;
  font-size: 1.5rem;
}

.seller-page__loading,
.seller-page__error,
.seller-page__empty {
  text-align: center;
  padding: 3rem 1rem;
  color: #666;
}

.seller-page__error .btn {
  margin-top: 1rem;
}

.draft-card {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  background: var(--bg-primary);
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 1.25rem;
  margin-bottom: 1rem;
}

.draft-card__info {
  flex: 1;
  min-width: 0;
}

.draft-card__title {
  margin: 0 0 0.5rem;
  font-size: 1.1rem;
}

.draft-card__meta {
  display: flex;
  gap: 1rem;
  font-size: 0.85rem;
  color: var(--text-secondary);
  margin: 0 0 0.5rem;
}

.draft-card__desc {
  font-size: 0.9rem;
  color: var(--text-secondary);
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.draft-card__actions {
  display: flex;
  gap: 0.5rem;
  margin-left: 1rem;
  flex-shrink: 0;
}

.btn {
  display: inline-block;
  padding: 0.5rem 1rem;
  border: 1px solid var(--border);
  border-radius: 6px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 0.9rem;
  cursor: pointer;
  text-decoration: none;
  transition: background 0.15s;
}

.btn:hover {
  background: var(--bg-secondary);
}

.btn--primary {
  background: var(--accent);
  color: #fff;
  border-color: var(--accent);
}

.btn--primary:hover {
  background: var(--accent-hover);
}

.btn--small {
  padding: 0.35rem 0.75rem;
  font-size: 0.85rem;
}

.btn--danger {
  color: #dc2626;
}

.btn--danger:hover {
  background: #fef2f2;
}
</style>
