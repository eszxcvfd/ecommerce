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

          <!-- Detailed description (only when present) -->
          <div v-if="product.mo_ta_chi_tiet" class="product-detail__description-detailed">
            <p class="product-detail__description-detailed-text">{{ product.mo_ta_chi_tiet }}</p>
          </div>

          <div class="product-detail__price">
            <span v-if="product.gia.mien_phi" class="product-detail__price-free">Miễn phí</span>
            <span v-else class="product-detail__price-paid">
              {{ formatPrice(product.gia.so_xu || 0) }} <span class="product-detail__xu-label">Xu</span>
            </span>
          </div>

          <!-- License (only when set) -->
          <div v-if="product.giay_phep" class="product-detail__license">
            <span class="product-detail__license-label">{{ product.giay_phep }}</span>
          </div>

          <!-- Rating -->
          <div v-if="product.diem_danh_gia > 0" class="product-detail__rating">
            <span class="product-detail__stars">{{ renderStars(product.diem_danh_gia) }}</span>
            <span class="product-detail__rating-count">
              ({{ product.so_luong_danh_gia }} đánh giá)
            </span>
          </div>

          <!-- Dates: creation + publish (if present) -->
          <div class="product-detail__dates">
            <span class="product-detail__date">Đăng ngày {{ formatDate(product.ngay_tao) }}</span>
            <span v-if="product.ngay_dang" class="product-detail__date-separator"> • </span>
            <span v-if="product.ngay_dang" class="product-detail__publish-date">Công bố ngày {{ formatDate(product.ngay_dang) }}</span>
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

          <!-- File table (only when files exist) -->
          <div v-if="product.tep && product.tep.length" class="product-detail__file-table">
            <h3 class="product-detail__file-table-title">Danh sách tệp</h3>
            <table class="product-detail__file-table-content">
              <thead>
                <tr>
                  <th>Tên tệp</th>
                  <th>Định dạng</th>
                  <th>Dung lượng</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(file, idx) in product.tep" :key="idx">
                  <td class="product-detail__file-name">{{ file.ten_tep }}</td>
                  <td><span class="format-badge">{{ file.dinh_dang.toUpperCase() }}</span></td>
                  <td class="product-detail__file-size">{{ formatFileSize(file.dung_luong_bytes) }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Seller line (only when set) -->
          <div v-if="product.nguoi_ban_hien_thi" class="product-detail__seller">
            <span>Người đăng: {{ product.nguoi_ban_hien_thi }}</span>
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

      <!-- Recommended products section -->
      <div v-if="recommendations.length > 0" class="product-detail__recommendations">
        <h2 class="product-detail__recommendations-title">Sản phẩm đề xuất</h2>
        <div class="product-detail__recommendations-grid">
          <ProductCard
            v-for="rec in recommendations"
            :key="rec.id"
            :product="rec"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useSanPhamDetail } from '~/composables/useCatalog'
import ProductCard from '~/components/ProductCard.vue'

const route = useRoute()
const id = computed(() => route.params.id as string)

const { product, recommendations, error, loaded, notFound, refresh } = useSanPhamDetail(id.value)

const hasError = computed(() => !!error.value && !notFound.value)
const fullyLoaded = computed(() => loaded.value)

function formatPrice(xu: number): string {
  return xu.toLocaleString('vi-VN')
}

function formatDate(dateStr: string): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString('vi-VN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  const k = 1024
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  const size = parseFloat((bytes / Math.pow(k, i)).toFixed(i > 0 ? 1 : 0))
  return `${size.toLocaleString('vi-VN')} ${units[i]}`
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
  padding: 24px 16px;
  font-family: system-ui, -apple-system, sans-serif;
  color: #1a1a2e;
}

.product-detail__loading,
.product-detail__not-found,
.product-detail__error {
  text-align: center;
  padding: 48px 16px;
}

.product-detail__back-link {
  display: inline-block;
  margin-bottom: 16px;
  color: #4361ee;
  text-decoration: none;
  font-weight: 500;
}

.product-detail__back-link:hover {
  text-decoration: underline;
}

.product-detail__main {
  display: flex;
  gap: 32px;
  flex-wrap: wrap;
}

.product-detail__image {
  flex: 0 0 380px;
  max-width: 100%;
}

.product-detail__image img {
  width: 100%;
  height: auto;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.product-detail__image-placeholder {
  width: 100%;
  aspect-ratio: 4 / 3;
  background: #f0f0f5;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  font-size: 14px;
}

.product-detail__info {
  flex: 1;
  min-width: 280px;
}

.product-detail__category {
  display: inline-block;
  padding: 4px 10px;
  background: #eef0ff;
  color: #4361ee;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 8px;
}

.product-detail__title {
  font-size: 24px;
  font-weight: 700;
  margin: 0 0 12px;
  line-height: 1.3;
}

.product-detail__description {
  font-size: 15px;
  line-height: 1.6;
  color: #444;
  margin: 0 0 16px;
}

.product-detail__description-detailed {
  background: #f8f9ff;
  border-left: 3px solid #4361ee;
  padding: 12px 16px;
  margin-bottom: 16px;
  border-radius: 0 8px 8px 0;
}

.product-detail__description-detailed-text {
  font-size: 14px;
  line-height: 1.6;
  color: #555;
  margin: 0;
  white-space: pre-line;
}

.product-detail__price {
  margin-bottom: 12px;
}

.product-detail__price-free {
  font-size: 20px;
  font-weight: 700;
  color: #2ec4b6;
}

.product-detail__price-paid {
  font-size: 20px;
  font-weight: 700;
  color: #e63946;
}

.product-detail__xu-label {
  font-size: 14px;
  font-weight: 500;
  color: #888;
}

.product-detail__license {
  margin-bottom: 12px;
}

.product-detail__license-label {
  display: inline-block;
  padding: 4px 12px;
  background: #f0fdf4;
  color: #16a34a;
  border: 1px solid #bbf7d0;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
}

.product-detail__rating {
  margin-bottom: 8px;
}

.product-detail__stars {
  color: #f59e0b;
  font-size: 18px;
  letter-spacing: 2px;
}

.product-detail__rating-count {
  font-size: 13px;
  color: #666;
  margin-left: 4px;
}

.product-detail__dates {
  font-size: 13px;
  color: #777;
  margin-bottom: 8px;
}

.product-detail__date {
  white-space: nowrap;
}

.product-detail__date-separator {
  color: #ccc;
}

.product-detail__downloads {
  font-size: 13px;
  color: #666;
  margin-bottom: 16px;
}

.product-detail__formats {
  margin-bottom: 16px;
}

.product-detail__formats-title {
  font-size: 14px;
  font-weight: 600;
  margin: 0 0 8px;
}

.product-detail__formats-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.format-badge {
  display: inline-block;
  padding: 3px 10px;
  background: #eef0ff;
  color: #4361ee;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
  font-family: monospace;
}

.product-detail__file-table {
  margin-bottom: 16px;
}

.product-detail__file-table-title {
  font-size: 14px;
  font-weight: 600;
  margin: 0 0 8px;
}

.product-detail__file-table-content {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

.product-detail__file-table-content th {
  text-align: left;
  padding: 8px 12px;
  background: #f5f5fa;
  font-weight: 600;
  color: #555;
  border-bottom: 2px solid #e0e0e8;
}

.product-detail__file-table-content td {
  padding: 8px 12px;
  border-bottom: 1px solid #eee;
}

.product-detail__file-table-content tr:last-child td {
  border-bottom: none;
}

.product-detail__file-name {
  font-weight: 500;
  word-break: break-all;
}

.product-detail__file-size {
  color: #666;
  text-align: right;
  white-space: nowrap;
}

.product-detail__seller {
  margin-bottom: 16px;
  font-size: 14px;
  color: #555;
}

.product-detail__seller span {
  font-weight: 500;
}

.product-detail__action {
  margin-top: 20px;
}

.product-detail__action-btn {
  display: inline-block;
  padding: 12px 32px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}

.product-detail__action-btn--free {
  background: #2ec4b6;
  color: #fff;
}

.product-detail__action-btn--free:hover {
  background: #25a99d;
}

.product-detail__action-btn--paid {
  background: #e63946;
  color: #fff;
}

.product-detail__action-btn--paid:hover {
  background: #c1121f;
}

.product-detail__retry-btn {
  padding: 8px 24px;
  background: #4361ee;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
}

.product-detail__retry-btn:hover {
  background: #3a56d4;
}

/* ── Recommended products ── */
.product-detail__recommendations {
  margin-top: 48px;
  padding-top: 32px;
  border-top: 1px solid #e2e8f0;
}

.product-detail__recommendations-title {
  font-size: 20px;
  font-weight: 700;
  margin-bottom: 20px;
  color: #1e293b;
}

.product-detail__recommendations-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 20px;
}
</style>
