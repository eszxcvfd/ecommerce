<template>
  <div class="product-detail">
    <!-- Loading state -->
    <div v-if="!fullyLoaded && !notFound && !hasError" class="product-detail__loading">
      <p>Đang tải...</p>
    </div>

    <!-- Error / not-found -->
    <div v-else-if="notFound" class="product-detail__not-found">
      <h2>Không tìm thấy</h2>
      <p>Sản phẩm không tồn tại hoặc đã bị ẩn.</p>
      <a href="/" class="product-detail__back-link">Quay lại danh sách</a>
    </div>

    <div v-else-if="hasError" class="product-detail__error">
      <p>Không thể tải thông tin sản phẩm.</p>
      <button class="product-detail__retry-btn" @click="refresh()">Thử lại</button>
    </div>

    <!-- Product detail -->
    <div v-else-if="product" class="product-detail__content">
      <a href="/" class="product-detail__back-link">&larr; Quay lại danh sách</a>

      <div class="product-detail__main">
        <!-- Demo image -->
        <div class="product-detail__image">
          <img
            v-if="product.anh_demo"
            :src="product.anh_demo"
            :alt="product.ten"
          />
          <div v-else class="product-detail__image-placeholder">
            <span>Không có ảnh</span>
          </div>
        </div>

        <!-- Info -->
        <div class="product-detail__info">
          <span class="product-detail__category">{{ product.danh_muc }}</span>
          <h1 class="product-detail__title">{{ product.ten }}</h1>

          <p v-if="product.mo_ta" class="product-detail__description">{{ product.mo_ta }}</p>
          <div class="product-detail__price">
            <span v-if="product.gia.mien_phi" class="product-detail__price-free">Miễn phí</span>
            <span v-else class="product-detail__price-paid">
              {{ formatPrice(product.gia.so_xu || 0) }} <span class="product-detail__xu-label">Xu</span>
            </span>
          </div>

          <!-- Rating -->
          <div v-if="product.diem_danh_gia > 0" class="product-detail__rating">
            <span class="product-detail__stars">{{ renderStars(product.diem_danh_gia) }}</span>
            <span class="product-detail__rating-count">
              ({{ product.so_luong_danh_gia }} đánh giá)
            </span>
          </div>

          <!-- Download count -->
          <div class="product-detail__downloads">
            <span>{{ product.so_luot_tai.toLocaleString('vi-VN') }} lượt tải</span>
          </div>

          <!-- Formats -->
          <div v-if="product.dinh_dang && product.dinh_dang.length" class="product-detail__formats">
            <h3 class="product-detail__formats-title">Định dạng tệp</h3>
            <div class="product-detail__formats-list">
              <span
                v-for="fmt in product.dinh_dang"
                :key="fmt"
                class="format-badge"
              >{{ fmt.toUpperCase() }}</span>
            </div>
          </div>

          <!-- Action -->
          <div class="product-detail__action">
            <span v-if="product.gia.mien_phi" class="product-detail__action-btn product-detail__action-btn--free">
              Tải xuống miễn phí
            </span>
            <span v-else class="product-detail__action-btn product-detail__action-btn--paid">
              Mua với {{ formatPrice(product.gia.so_xu || 0) }} Xu
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useSanPhamDetail } from '~/composables/useCatalog'

const route = useRoute()
const id = computed(() => route.params.id as string)

const { product, error, loaded, notFound, refresh } = useSanPhamDetail(id.value)

const hasError = computed(() => !!error.value && !notFound.value)
const fullyLoaded = computed(() => loaded.value)

function formatPrice(xu: number): string {
  return xu.toLocaleString('vi-VN')
}

function renderStars(rating: number): string {
  const full = Math.floor(rating)
  const half = rating - full >= 0.5 ? 1 : 0
  return '★'.repeat(full) + (half ? '½' : '') + '☆'.repeat(5 - full - half)
}
</script>

<style scoped>
.product-detail {
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.product-detail__loading,
.product-detail__not-found,
.product-detail__error {
  text-align: center;
  padding: 4rem 1rem;
  color: var(--color-text-secondary);
}

.product-detail__not-found h2 {
  font-size: 1.5rem;
  margin-bottom: 0.5rem;
  color: var(--color-text);
}

.product-detail__back-link {
  display: inline-block;
  margin-top: 1rem;
  color: var(--color-accent);
  font-size: 0.95rem;
}

.product-detail__retry-btn {
  margin-top: 0.75rem;
  padding: 0.5rem 1.5rem;
  border-radius: 8px;
  border: 1px solid var(--color-accent);
  background: transparent;
  color: var(--color-accent);
  cursor: pointer;
  font-family: inherit;
  font-size: 0.9rem;
}

.product-detail__main {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

@media (min-width: 768px) {
  .product-detail__main {
    flex-direction: row;
    align-items: flex-start;
  }
}

/* Image */
.product-detail__image {
  flex: 0 0 auto;
  width: 100%;
  max-width: 400px;
  border-radius: var(--radius);
  overflow: hidden;
  background: var(--color-surface);
  box-shadow: var(--shadow);
}

.product-detail__image img {
  width: 100%;
  height: auto;
  display: block;
}

.product-detail__image-placeholder {
  width: 100%;
  aspect-ratio: 4 / 3;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg);
  color: var(--color-text-secondary);
  font-size: 3rem;
  font-weight: 700;
}

/* Info */
.product-detail__info {
  flex: 1;
  min-width: 0;
}

.product-detail__category {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 999px;
  background: var(--color-bg);
  font-size: 0.85rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  margin-bottom: 0.75rem;
}

.product-detail__title {
  font-size: 1.75rem;
  font-weight: 700;
  line-height: 1.3;
  margin-bottom: 1rem;
}

.product-detail__description {
  color: var(--color-text-secondary);
  font-size: 1rem;
  line-height: 1.6;
  margin-bottom: 1.5rem;
}

/* Price */
.product-detail__price {
  margin-bottom: 1rem;
}

.product-detail__price-free {
  display: inline-block;
  padding: 0.35rem 1rem;
  border-radius: 8px;
  background: var(--color-free);
  color: white;
  font-weight: 600;
  font-size: 1rem;
}

.product-detail__price-paid {
  display: inline-block;
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-paid);
}

.product-detail__xu-label {
  font-size: 1rem;
  font-weight: 500;
}

/* Rating */
.product-detail__rating {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
}

.product-detail__stars {
  color: var(--color-star);
  font-size: 1.1rem;
}

.product-detail__rating-count {
  font-size: 0.85rem;
  color: var(--color-text-secondary);
}

/* Downloads */
.product-detail__downloads {
  font-size: 0.9rem;
  color: var(--color-text-secondary);
  margin-bottom: 1.25rem;
}

/* Formats */
.product-detail__formats {
  margin-bottom: 1.5rem;
}

.product-detail__formats-title {
  font-size: 0.9rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  color: var(--color-text-secondary);
}

.product-detail__formats-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.format-badge {
  display: inline-block;
  padding: 0.3rem 0.75rem;
  border-radius: 6px;
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  font-size: 0.8rem;
  font-weight: 600;
  font-family: monospace;
  color: var(--color-text-secondary);
}

/* Action */
.product-detail__action {
  margin-top: 1rem;
}

.product-detail__action-btn {
  display: inline-block;
  padding: 0.75rem 2rem;
  border-radius: 8px;
  font-weight: 600;
  font-size: 1rem;
  cursor: default;
}

.product-detail__action-btn--free {
  background: var(--color-free);
  color: white;
}

.product-detail__action-btn--paid {
  background: var(--color-accent);
  color: white;
}
</style>
