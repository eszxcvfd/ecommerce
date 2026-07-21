export interface Gia {
  mien_phi: boolean
  so_xu?: number
}

export interface Tep {
  ten_tep: string
  dinh_dang: string
  dung_luong_bytes: number
}

export interface SanPhamSo {
  id: string
  ten: string
  mo_ta: string
  mo_ta_chi_tiet: string
  anh_demo: string
  gia: Gia
  danh_muc: string
  dinh_dang: string[]
  diem_danh_gia: number
  so_luong_danh_gia: number
  ngay_tao: string
  ngay_dang: string
  so_luot_tai: number
  giay_phep: string
  nguoi_ban_hien_thi: string
  tep: Tep[]
}

export interface SanPhamChiTietResponse {
  san_pham: SanPhamSo
  san_pham_de_xuat: SanPhamSo[]
}

interface DanhMucResponse {
  danh_muc: string[]
}

interface DinhDangResponse {
  dinh_dang: string[]
}

interface SanPhamResponse {
  san_pham: SanPhamSo[]
}

export interface CatalogResult<T> {
  data: Ref<T | null>
  error: Ref<unknown>
  loaded: Ref<boolean>
  refresh: () => Promise<void>
}

export interface CatalogSearchState {
  q: Ref<string>
  danhMuc: Ref<string>
  dinhDang: Ref<string>
  sort: Ref<string>
  results: Ref<SanPhamSo[]>
  loading: Ref<boolean>
  error: Ref<unknown>
}

/**
 * Fetch the list of all six categories.
 */
export function useDanhMuc(): CatalogResult<string[]> {
  const { data, error, refresh } = useFetch<DanhMucResponse>('/api/v1/danh-muc', {
    watch: false,
  })

  return {
    data: computed(() => data.value?.danh_muc ?? null),
    error,
    loaded: computed(() => data.value !== null),
    refresh,
  }
}

/**
 * Fetch the list of available formats.
 */
export function useDinhDang(): CatalogResult<string[]> {
  const { data, error, refresh } = useFetch<DinhDangResponse>('/api/v1/dinh-dang', {
    watch: false,
  })

  return {
    data: computed(() => data.value?.dinh_dang ?? null),
    error,
    loaded: computed(() => data.value !== null),
    refresh,
  }
}

/**
 * Build a query string from reactive search state, omitting empty values.
 */
function buildQueryString(q: string, danhMuc: string, dinhDang: string, sort: string): string {
  const params = new URLSearchParams()
  if (q) params.set('q', q)
  if (danhMuc) params.set('danh_muc', danhMuc)
  if (dinhDang) params.set('dinh_dang', dinhDang)
  if (sort) params.set('sort', sort)
  return params.toString()
}

/**
 * Search state factory with debounced search text.
 * Returns reactive values that sync with URL query params.
 */
export function useCatalogSearch() {
  const route = useRoute()
  const router = useRouter()

  const q = ref(decodeURIComponent((route.query.q as string) || ''))
  const danhMuc = ref((route.query.danh_muc as string) || '')
  const dinhDang = ref((route.query.dinh_dang as string) || '')
  const sort = ref((route.query.sort as string) || '')
  const results = ref<SanPhamSo[]>([])
  const loading = ref(false)
  const error = ref<unknown>(null)

  // Debounced search text
  const debouncedQ = ref(q.value)
  let debounceTimer: number | undefined
  // Fetch products — extracted so it can be called directly on refresh
  async function fetchProducts() {
    loading.value = true
    error.value = null
    const qs = buildQueryString(debouncedQ.value, danhMuc.value, dinhDang.value, sort.value)
    const url = qs ? `/api/v1/san-pham?${qs}` : '/api/v1/san-pham'
    try {
      const res = await $fetch<SanPhamResponse>(url)
      results.value = res.san_pham
    } catch (e) {
      error.value = e
      results.value = []
    } finally {
      loading.value = false
    }
  }
  
  // Fetch on filter change
  watch(
    [debouncedQ, danhMuc, dinhDang, sort],
    fetchProducts,
    { immediate: true },
  )

  // Reset all filters
  function resetAll() {
    q.value = ''
    danhMuc.value = ''
    dinhDang.value = ''
    sort.value = ''
  }

  return {
    rawQ: q,
    q,
    danhMuc,
    dinhDang,
    sort,
    reset: resetAll,
    refresh: fetchProducts,
    products: results,
    error,
    loaded: loading,
  }
}

/**
 * Fetch a single product by ID from the detail API.
 */
export function useSanPhamDetail(id: string) {
  const { data, error, refresh } = useFetch<SanPhamChiTietResponse>(`/api/v1/san-pham/${id}`, {
    watch: false,
  })

  const product = computed(() => data.value?.san_pham ?? null)
  const recommendations = computed(() => data.value?.san_pham_de_xuat ?? [])
  const notFound = computed(() => {
    if (!error.value) return false
    const err = error.value
    if (err && typeof err === 'object' && 'statusCode' in err) {
      const withCode = err as Record<string, unknown>
      return typeof withCode.statusCode === 'number' && withCode.statusCode === 404
    }
    if (err && typeof err === 'object' && 'status' in err) {
      const withStatus = err as Record<string, unknown>
      return typeof withStatus.status === 'number' && withStatus.status === 404
    }
    return false
  })
  return {
    product,
    recommendations,
    error,
    loaded: computed(() => data.value !== null && data.value !== undefined),
    notFound,
    refresh,
  }
}
