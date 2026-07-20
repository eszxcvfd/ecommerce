<template>
  <div class="catalog">
    <header class="catalog__header">
      <h1 class="catalog__title">Sàn Sản phẩm số</h1>
      <p class="catalog__subtitle">
        Khám phá các Sản phẩm số thiết kế và kỹ thuật đã được duyệt
      </p>
    </header>

    <!-- Search and filters bar -->
    <div v-if="fullyLoaded" class="catalog__controls">
      <div class="catalog__search">
        <input
          v-model="rawQ"
          type="search"
          placeholder="Tìm kiếm sản phẩm..."
          class="catalog__search-input"
          aria-label="Tìm kiếm sản phẩm"
        />
        <button
          v-if="rawQ || danhMuc || dinhDang || sort"
          class="catalog__reset-btn"
          @click="reset"
          aria-label="Đặt lại bộ lọc"
        >
          Đặt lại
        </button>
      </div>

      <div class="catalog__filters">
        <!-- Category filter -->
        <div class="catalog__category-grid">
          <button
            class="catalog__filter-btn"
            :class="{ 'catalog__filter-btn--active': danhMuc === '' }"
            @click="danhMuc = ''"
          >
            Tất cả
          </button>
          <button
            v-for="dm in danhMucList"
            :key="dm"
            class="catalog__filter-btn"
            :class="{ 'catalog__filter-btn--active': danhMuc === dm }"
            @click="danhMuc = danhMuc === dm ? '' : dm"
          >
            {{ dm }}
          </button>
        </div>

        <div class="catalog__filter-row">
          <!-- Format filter -->
          <select v-model="dinhDang" class="catalog__select" aria-label="Lọc theo định dạng">
            <option value="">Tất cả định dạng</option>
            <option v-for="f in formatList" :key="f" :value="f">{{ f }}</option>
          </select>

          <!-- Sort -->
          <select v-model="sort" class="catalog__select" aria-label="Sắp xếp">
            <option value="newest">Mới nhất</option>
            <option value="popular">Phổ biến</option>
            <option value="price_asc">Giá: Thấp đến cao</option>
            <option value="price_desc">Giá: Cao đến thấp</option>
            <option value="rating">Đánh giá</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Error state -->
    <div v-if="hasError && fullyLoaded" class="catalog__error">
      <p>Không thể tải dữ liệu. Vui lòng thử lại sau.</p>
      <button class="catalog__retry-btn" @click="retry">Thử lại</button>
    </div>

    <!-- Loading state -->
    <div v-else-if="!fullyLoaded" class="catalog__loading">
      <p>Đang tải...</p>
    </div>

    <!-- Products section -->
    <section v-if="fullyLoaded && !hasError" class="catalog__products">
      <h2 class="catalog__section-title">Sản phẩm số</h2>

      <div v-if="!loaded" class="catalog__loading">
        <p>Đang tải...</p>
      </div>

      <div v-else-if="products.length === 0" class="catalog__empty">
        <p>Không tìm thấy Sản phẩm số nào phù hợp.</p>
        <p v-if="q || danhMuc || dinhDang || sort" class="catalog__empty-hint">
          Thử thay đổi bộ lọc hoặc
          <button class="catalog__reset-link" @click="reset">đặt lại</button>
          để xem tất cả sản phẩm.
        </p>
      </div>

      <div v-else class="catalog__product-grid">
        <ProductCard
          v-for="sp in products"
          :key="sp.id"
          :product="sp"
        />
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { useDanhMuc, useDinhDang, useCatalogSearch } from '~/composables/useCatalog'

const { data: danhMucList } = await useDanhMuc()
const { data: formatList } = await useDinhDang()

const {
  rawQ, q, danhMuc, dinhDang, sort, reset, refresh,
  products, error, loaded,
} = useCatalogSearch()

const hasError = computed(() => !!error.value)
const fullyLoaded = computed(() => danhMucList.value.length > 0)

function retry() {
  refresh()
}
</script>

<style scoped>
.catalog {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.catalog__header {
  text-align: center;
  margin-bottom: 2rem;
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

.catalog__empty-hint {
  margin-top: 0.5rem;
  font-size: 0.9rem;
}

.catalog__reset-link {
  background: none;
  border: none;
  color: var(--color-primary, #2563eb);
  cursor: pointer;
  text-decoration: underline;
  font-family: inherit;
  font-size: inherit;
}

.catalog__section-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 1rem;
}

/* Controls */
.catalog__controls {
  margin-bottom: 2rem;
}

.catalog__search {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.catalog__search-input {
  flex: 1;
  padding: 0.65rem 1rem;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  font-size: 1rem;
  font-family: inherit;
  background: transparent;
  color: inherit;
}

.catalog__reset-btn {
  padding: 0.5rem 1rem;
  border-radius: 8px;
  border: 1px solid var(--color-border);
  background: transparent;
  cursor: pointer;
  font-family: inherit;
  font-size: 0.9rem;
  color: inherit;
  white-space: nowrap;
}

.catalog__filters {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.catalog__category-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.catalog__filter-btn {
  padding: 0.4rem 1rem;
  border-radius: 999px;
  border: 1px solid var(--color-border);
  font-size: 0.85rem;
  font-weight: 500;
  cursor: pointer;
  background: transparent;
  color: inherit;
  font-family: inherit;
  transition: all 0.15s;
}

.catalog__filter-btn--active {
  background: var(--color-primary, #2563eb);
  color: white;
  border-color: var(--color-primary, #2563eb);
}

.catalog__filter-row {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.catalog__select {
  padding: 0.5rem 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  font-size: 0.9rem;
  font-family: inherit;
  background: transparent;
  color: inherit;
  min-width: 160px;
}

.catalog__retry-btn {
  margin-top: 0.75rem;
  padding: 0.5rem 1.5rem;
  border-radius: 8px;
  border: 1px solid var(--color-primary, #2563eb);
  background: transparent;
  color: var(--color-primary, #2563eb);
  cursor: pointer;
  font-family: inherit;
  font-size: 0.9rem;
}

/* Products */
.catalog__product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 1.5rem;
}
</style>
