<template>
  <div class="catalog">
    <header class="catalog__header">
      <h1 class="catalog__title">Sàn Sản phẩm số</h1>
      <p class="catalog__subtitle">
        Khám phá các Sản phẩm số thiết kế và kỹ thuật đã được duyệt
      </p>
    </header>

    <!-- Error state -->
    <div v-if="hasError" class="catalog__error">
      <p>Không thể tải danh mục. Vui lòng thử lại sau.</p>
    </div>

    <!-- Loading state -->
    <div v-else-if="!fullyLoaded" class="catalog__loading">
      <p>Đang tải...</p>
    </div>

    <!-- Categories section -->
    <section v-if="fullyLoaded" class="catalog__categories">
      <h2 class="catalog__section-title">Danh mục</h2>
      <div class="catalog__category-grid">
        <div
          v-for="dm in danhMuc"
          :key="dm"
          class="catalog__category-card"
        >
          {{ dm }}
        </div>
      </div>
    </section>

    <!-- Products section -->
    <section v-if="fullyLoaded" class="catalog__products">
      <h2 class="catalog__section-title">Sản phẩm số</h2>

      <div v-if="sanPham.length === 0" class="catalog__empty">
        <p>Không có Sản phẩm số nào.</p>
      </div>

      <div v-else class="catalog__product-grid">
        <ProductCard
          v-for="sp in sanPham"
          :key="sp.id"
          :product="sp"
        />
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { useDanhMuc, useSanPham } from '~/composables/useCatalog'

const { data: danhMuc, error: catError, loaded: catLoaded } = await useDanhMuc()
const { data: sanPham, error: prodError, loaded: prodLoaded } = await useSanPham()

const hasError = computed(() => !!(catError.value || prodError.value))
const fullyLoaded = computed(() => catLoaded.value && prodLoaded.value)
</script>

<style scoped>
.catalog {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.catalog__header {
  text-align: center;
  margin-bottom: 2.5rem;
}

.catalog__title {
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
}

.catalog__subtitle {
  color: var(--color-text-secondary);
  font-size: 1.1rem;
}

.catalog__error,
.catalog__loading,
.catalog__empty {
  text-align: center;
  padding: 3rem 1rem;
  color: var(--color-text-secondary);
}

.catalog__section-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 1rem;
}

.catalog__category-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
  margin-bottom: 2.5rem;
}

.catalog__category-card {
  padding: 0.5rem 1.25rem;
  border-radius: 999px;
  border: 1px solid var(--color-border);
  font-size: 0.9rem;
  font-weight: 500;
}

.catalog__product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 1.5rem;
}
</style>
