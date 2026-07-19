<template>
  <article class="product-card">
    <div class="product-card__image">
      <img
        v-if="product.anh_demo"
        :src="product.anh_demo"
        :alt="product.ten"
        loading="lazy"
      />
      <div v-else class="product-card__placeholder">
        <span>Không có ảnh</span>
      </div>
    </div>

    <div class="product-card__body">
      <span class="product-card__category">{{ product.danh_muc }}</span>
      <h3 class="product-card__title">{{ product.ten }}</h3>

      <div class="product-card__meta">
        <span
          class="product-card__price"
          :class="{ 'product-card__price--free': product.gia.mien_phi }"
        >
          {{ product.gia.mien_phi ? 'Miễn phí' : `${formatPrice(product.gia.so_xu ?? 0)} Xu` }}
        </span>

        <span v-if="product.so_luong_danh_gia > 0" class="product-card__rating">
          <span class="product-card__stars">{{ renderStars(product.diem_danh_gia) }}</span>
          {{ product.diem_danh_gia.toFixed(1) }}
          <span class="product-card__rating-count">({{ product.so_luong_danh_gia }})</span>
        </span>
        <span v-else class="product-card__rating product-card__rating--none">
          Chưa có đánh giá
        </span>
      </div>
    </div>
  </article>
</template>

<script setup lang="ts">
import type { SanPhamSo } from '~/composables/useCatalog'

defineProps<{
  product: SanPhamSo
}>()

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
.product-card {
  background: var(--color-surface);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  transition: transform 0.2s, box-shadow 0.2s;
}

.product-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.product-card__image {
  width: 100%;
  height: 180px;
  overflow: hidden;
  background: #e8e8ed;
}

.product-card__image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.product-card__placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--color-text-secondary);
  font-size: 0.875rem;
}

.product-card__body {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  flex: 1;
}

.product-card__category {
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--color-accent);
}

.product-card__title {
  font-size: 1rem;
  font-weight: 600;
  line-height: 1.4;
  color: var(--color-text);
}

.product-card__meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: auto;
}

.product-card__price {
  font-size: 0.875rem;
  font-weight: 700;
  color: var(--color-paid);
}

.product-card__price--free {
  color: var(--color-free);
}

.product-card__rating {
  font-size: 0.8rem;
  color: var(--color-text-secondary);
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.product-card__rating--none {
  font-style: italic;
}

.product-card__stars {
  color: var(--color-star);
}

.product-card__rating-count {
  font-size: 0.75rem;
}
</style>
