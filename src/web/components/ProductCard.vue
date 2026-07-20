<template>
  <NuxtLink :to="'/san-pham/' + product.id" class="product-card product-card__link">
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

      <p v-if="descriptionExcerpt" class="product-card__description">
        {{ descriptionExcerpt }}
      </p>

      <div v-if="product.dinh_dang && product.dinh_dang.length" class="product-card__formats">
        <span
          v-for="fmt in product.dinh_dang"
          :key="fmt"
          class="product-card__format-badge"
        >
          {{ fmt }}
        </span>
      </div>

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
  </NuxtLink>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { SanPhamSo } from '~/composables/useCatalog'

const props = defineProps<{
  product: SanPhamSo
}>()

const descriptionExcerpt = computed(() => {
  if (!props.product.mo_ta) return ''
  if (props.product.mo_ta.length <= 120) return props.product.mo_ta
  return props.product.mo_ta.slice(0, 120) + '...'
})

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
.product-card__link {
  display: flex;
  flex-direction: column;
  text-decoration: none;
  color: inherit;
}
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

.product-card__description {
  font-size: 0.8rem;
  line-height: 1.5;
  color: var(--color-text-secondary);
  margin: 0;
}

.product-card__formats {
  display: flex;
  flex-wrap: wrap;
  gap: 0.35rem;
}

.product-card__format-badge {
  display: inline-block;
  padding: 0.15rem 0.5rem;
  border-radius: 4px;
  background: var(--color-badge-bg, #e8e8ed);
  color: var(--color-badge-text, #555);
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: lowercase;
  line-height: 1.4;
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
