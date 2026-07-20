export interface Gia {
  mien_phi: boolean
  so_xu?: number
}

export interface SanPhamSo {
  id: string
  ten: string
  mo_ta: string
  anh_demo: string
  gia: Gia
  danh_muc: string
  dinh_dang: string[]
  diem_danh_gia: number
  so_luong_danh_gia: number
  ngay_tao: string
  so_luot_tai: number
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
  data: Ref<T>
  error: Ref<unknown>
  loaded: Ref<boolean>
}

export interface CatalogSearchState {
  q: Ref<string>
  danhMuc: Ref<string>
  dinhDang: Ref<string>
  sort: Ref<string>
}

/**
 * Fetch the list of all six categories.
 */
export function useDanhMuc(): CatalogResult<string[]> {
  const { data, error } = useFetch<DanhMucResponse>('/api/v1/danh-muc')
  return {
    data: computed(() => data.value?.danh_muc ?? []),
    error,
    loaded: computed(() => data.value !== null),
  }
}

/**
 * Fetch the list of available formats.
 */
export function useDinhDang(): CatalogResult<string[]> {
  const { data, error } = useFetch<DinhDangResponse>('/api/v1/dinh-dang')
  return {
    data: computed(() => data.value?.dinh_dang ?? []),
    error,
    loaded: computed(() => data.value !== null),
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

  // Initialize from URL query params
  const rawQ = ref((route.query.q as string) || '')
  const danhMuc = ref((route.query.danh_muc as string) || '')
  const dinhDang = ref((route.query.dinh_dang as string) || '')
  const sort = ref((route.query.sort as string) || '')

  // Debounced search text (delayed 250ms)
  const q = ref(rawQ.value)
  let debounceTimer: ReturnType<typeof setTimeout> | null = null
  watch(rawQ, (val) => {
    if (debounceTimer) clearTimeout(debounceTimer)
    debounceTimer = setTimeout(() => {
      q.value = val
    }, 250)
  })

  // Sync state to URL
  function syncUrl() {
    const params: Record<string, string> = {}
    if (q.value) params.q = q.value
    if (danhMuc.value) params.danh_muc = danhMuc.value
    if (dinhDang.value) params.dinh_dang = dinhDang.value
    if (sort.value) params.sort = sort.value
    router.replace({ query: params })
  }

  watch([q, danhMuc, dinhDang, sort], () => { syncUrl() })

  function reset() {
    rawQ.value = ''
    q.value = ''
    danhMuc.value = ''
    dinhDang.value = ''
    sort.value = ''
  }

  // Build API URL reactively
  const apiUrl = computed(() => {
    const qs = buildQueryString(q.value, danhMuc.value, dinhDang.value, sort.value)
    return `/api/v1/san-pham${qs ? '?' + qs : ''}`
  })

  // Fetch products
  const { data, error, refresh } = useFetch<SanPhamResponse>(apiUrl, { watch: false })

  // Re-fetch when URL changes
  watch(apiUrl, () => { refresh() })

  return {
    // Search state
    rawQ,
    q,
    danhMuc,
    dinhDang,
    sort,
    reset,
    // Results
    products: computed(() => data.value?.san_pham ?? []),
    error,
    loaded: computed(() => data.value !== null),
    // Re-fetch trigger
    refresh,
  }
}
