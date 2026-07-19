export interface Gia {
  mien_phi: boolean
  so_xu?: number
}

export interface SanPhamSo {
  id: string
  ten: string
  anh_demo: string
  gia: Gia
  danh_muc: string
  diem_danh_gia: number
  so_luong_danh_gia: number
}

interface DanhMucResponse {
  danh_muc: string[]
}

interface SanPhamResponse {
  san_pham: SanPhamSo[]
}

export interface CatalogResult<T> {
  data: Ref<T>
  error: Ref<unknown>
  loaded: Ref<boolean>
}

function useCatalogFetch<T>(path: string): CatalogResult<T> {
  const { data, error } = useFetch<T>(path)
  const loaded = computed(() => data.value !== null && data.value !== undefined)
  return { data: data as Ref<T>, error, loaded }
}

/**
 * Fetch the list of all six categories.
 */
export function useDanhMuc(): CatalogResult<string[]> {
  const result = useCatalogFetch<DanhMucResponse>('/api/v1/danh-muc')
  return {
    data: computed(() => result.data.value?.danh_muc ?? []),
    error: result.error,
    loaded: result.loaded,
  }
}

/**
 * Fetch the list of approved/public products.
 */
export function useSanPham(): CatalogResult<SanPhamSo[]> {
  const result = useCatalogFetch<SanPhamResponse>('/api/v1/san-pham')
  return {
    data: computed(() => result.data.value?.san_pham ?? []),
    error: result.error,
    loaded: result.loaded,
  }
}
