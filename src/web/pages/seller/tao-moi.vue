<template>
  <div class="seller-form-page">
    <h1>Tạo bản nháp Sản phẩm số</h1>

    <div v-if="error" class="form-error">{{ error }}</div>
    <div v-if="success" class="form-success">Đã tạo bản nháp thành công!</div>

    <form @submit.prevent="handleSubmit" class="seller-form">
      <div class="form-group">
        <label for="ten">Tên sản phẩm *</label>
        <input id="ten" v-model="form.ten" required class="form-input" placeholder="VD: Bản vẽ nhà phố 3 tầng" />
      </div>

      <div class="form-group">
        <label for="mo_ta">Mô tả ngắn</label>
        <textarea id="mo_ta" v-model="form.mo_ta" class="form-input" rows="2" placeholder="Tóm tắt ngắn về sản phẩm"></textarea>
      </div>

      <div class="form-group">
        <label for="mo_ta_chi_tiet">Mô tả chi tiết</label>
        <textarea id="mo_ta_chi_tiet" v-model="form.mo_ta_chi_tiet" class="form-input" rows="4" placeholder="Mô tả đầy đủ về sản phẩm"></textarea>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label for="danh_muc">Danh mục *</label>
          <select id="danh_muc" v-model="form.danh_muc" required class="form-input">
            <option value="" disabled>Chọn danh mục</option>
            <option v-for="dm in danhMucList" :key="dm" :value="dm">{{ dm }}</option>
          </select>
        </div>

        <div class="form-group">
          <label for="giay_phep">Giấy phép</label>
          <input id="giay_phep" v-model="form.giay_phep" class="form-input" placeholder="VD: Giấy phép tiêu chuẩn" />
        </div>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label for="so_xu">Giá (Xu)</label>
          <input id="so_xu" v-model.number="form.so_xu" type="number" min="0" class="form-input" />
        </div>

        <div class="form-group checkbox-group">
          <label>
            <input id="mien_phi" v-model="form.mien_phi" type="checkbox" />
            Miễn phí
          </label>
        </div>
      </div>

      <div class="form-group">
        <label for="anh_demo">Ảnh demo (URL)</label>
        <input id="anh_demo" v-model="form.anh_demo" class="form-input" placeholder="https://example.com/demo.jpg" />
      </div>

      <!-- File entries -->
      <div class="form-section">
        <h2>Tệp sản phẩm *</h2>
        <p class="form-help">Hỗ trợ các định dạng: AutoCAD (dwg, dxf), SketchUp (skp), Revit (rvt, rfa), 3ds Max (max, 3ds), Photoshop (psd), Illustrator (ai, eps)</p>

        <div v-for="(tep, idx) in form.tep" :key="idx" class="file-entry">
          <div class="file-entry__row">
            <div class="form-group">
              <label :for="'tep_ten_' + idx">Tên tệp</label>
              <input :id="'tep_ten_' + idx" v-model="tep.ten_tep" class="form-input" placeholder="VD: mat-bang.dwg" />
            </div>
            <div class="form-group">
              <label :for="'tep_dd_' + idx">Định dạng</label>
              <select :id="'tep_dd_' + idx" v-model="tep.dinh_dang" class="form-input">
                <option value="" disabled>Chọn</option>
                <option v-for="(label, fmt) in ALLOWED_FORMATS" :key="fmt" :value="fmt">{{ label }} (.{{ fmt }})</option>
              </select>
            </div>
            <div class="form-group">
              <label :for="'tep_size_' + idx">Dung lượng (bytes)</label>
              <input :id="'tep_size_' + idx" v-model.number="tep.dung_luong_bytes" type="number" min="1" class="form-input" />
            </div>
            <button type="button" class="btn btn--danger btn--small file-entry__remove" @click="removeFile(idx)">Xóa</button>
          </div>
        </div>

        <button type="button" class="btn" @click="addFile">+ Thêm tệp</button>
      </div>

      <div class="form-actions">
        <NuxtLink to="/seller" class="btn">Hủy</NuxtLink>
        <button type="submit" class="btn btn--primary" :disabled="submitting">
          {{ submitting ? 'Đang lưu...' : 'Tạo bản nháp' }}
        </button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ALLOWED_FORMATS } from '~/composables/useSeller'
import type { TepInput, DraftInput } from '~/composables/useSeller'

const router = useRouter()
const { createDraft } = useSeller()
const { data: danhMucList } = await useDanhMuc()

const form = reactive<DraftInput>({
  ten: '',
  mo_ta: '',
  mo_ta_chi_tiet: '',
  anh_demo: '',
  mien_phi: false,
  so_xu: 0,
  danh_muc: '',
  giay_phep: '',
  tep: [],
})

const error = ref<string | null>(null)
const success = ref(false)
const submitting = ref(false)

function addFile() {
  form.tep.push({ ten_tep: '', dinh_dang: '', dung_luong_bytes: 0 })
}

function removeFile(idx: number) {
  form.tep.splice(idx, 1)
}

async function handleSubmit() {
  error.value = null
  success.value = false
  submitting.value = true

  try {
    await createDraft(form)
    success.value = true
    setTimeout(() => router.push('/seller'), 1500)
  } catch {
    error.value = 'Không thể tạo bản nháp'
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.seller-form-page {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
}

.seller-form-page h1 {
  margin: 0 0 1.5rem;
  font-size: 1.5rem;
}

.form-error {
  background: #fef2f2;
  color: #dc2626;
  padding: 0.75rem 1rem;
  border-radius: 6px;
  margin-bottom: 1rem;
}

.form-success {
  background: #f0fdf4;
  color: #16a34a;
  padding: 0.75rem 1rem;
  border-radius: 6px;
  margin-bottom: 1rem;
}

.seller-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-row {
  display: flex;
  gap: 1rem;
}

.form-row .form-group {
  flex: 1;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.form-group label {
  font-weight: 500;
  font-size: 0.9rem;
}

.form-input {
  padding: 0.6rem 0.75rem;
  border: 1px solid var(--border);
  border-radius: 6px;
  font-size: 0.9rem;
  font-family: inherit;
  background: var(--bg-primary);
}

.form-input:focus {
  outline: none;
  border-color: var(--accent);
  box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.15);
}

textarea.form-input {
  resize: vertical;
}

.checkbox-group {
  justify-content: flex-end;
}

.checkbox-group label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
}

.form-section {
  border-top: 1px solid var(--border);
  padding-top: 1.25rem;
}

.form-section h2 {
  font-size: 1.1rem;
  margin: 0 0 0.5rem;
}

.form-help {
  font-size: 0.85rem;
  color: var(--text-secondary);
  margin: 0 0 1rem;
}

.file-entry {
  margin-bottom: 0.75rem;
}

.file-entry__row {
  display: flex;
  gap: 0.75rem;
  align-items: flex-end;
}

.file-entry__row .form-group {
  flex: 1;
  min-width: 0;
}

.file-entry__remove {
  margin-bottom: 0.1rem;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  padding-top: 1rem;
  border-top: 1px solid var(--border);
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

.btn--primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
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

.file-entry__row .btn--small {
  margin-bottom: 0.15rem;
}
</style>
