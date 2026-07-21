export interface TaiKhoanPublic {
  id: string
  email: string
  ten: string
  vai_tro: string
  created_at: string
}

export interface HOSoBan {
  id: string
  tai_khoan_id: string
  trang_thai: string
  created_at: string
  updated_at: string
}

export interface DangKyRequest {
  email: string
  password: string
  ten: string
}

export interface DangNhapRequest {
  email: string
  password: string
}

export interface DangNhapResponse {
  tai_khoan: TaiKhoanPublic
  token: string
}

interface AuthState {
  account: TaiKhoanPublic | null
  token: string | null
  loading: boolean
  error: string | null
}

/**
 * Auth composable providing login, logout, registration, and session management.
 */
export const useAuth = () => {
  const state = useState<AuthState>('auth', () => ({
    account: null,
    token: null,
    loading: false,
    error: null,
  }))

  // Initialize from stored session
  if (process.client) {
    const stored = localStorage.getItem('auth_token')
    const storedAccount = localStorage.getItem('auth_account')
    if (stored) {
      state.value.token = stored
      // Try to fetch current account
      fetchMe(stored).catch(() => {
        // Token expired, clear
        clearAuth()
      })
    }
    if (storedAccount && state.value.account === null) {
      try {
        state.value.account = JSON.parse(storedAccount)
      } catch { /* ignore */ }
    }
  }

  async function fetchMe(token: string): Promise<TaiKhoanPublic> {
    const { data, error } = await useFetch<TaiKhoanPublic>('/api/v1/tai-khoan/me', {
      headers: { Authorization: `Bearer ${token}` },
    })
    if (error.value) {
      throw new Error('Token không hợp lệ')
    }
    return data.value!
  }

  async function dangKy(req: DangKyRequest): Promise<TaiKhoanPublic> {
    state.value.loading = true
    state.value.error = null
    try {
      const { data, error } = await useFetch<{ tai_khoan: TaiKhoanPublic }>('/api/v1/dang-ky', {
        method: 'POST',
        body: req,
        headers: { 'Content-Type': 'application/json' },
      })
      if (error.value) {
        const msg = error.value.data?.message || 'Đăng ký thất bại'
        state.value.error = msg
        throw new Error(msg)
      }
      return data.value!.tai_khoan
    } finally {
      state.value.loading = false
    }
  }

  async function dangNhap(req: DangNhapRequest): Promise<void> {
    state.value.loading = true
    state.value.error = null
    try {
      const { data, error } = await useFetch<DangNhapResponse>('/api/v1/dang-nhap', {
        method: 'POST',
        body: req,
        headers: { 'Content-Type': 'application/json' },
      })
      if (error.value) {
        const msg = error.value.data?.message || 'Đăng nhập thất bại'
        state.value.error = msg
        throw new Error(msg)
      }
      state.value.account = data.value!.tai_khoan
      state.value.token = data.value!.token
      if (process.client) {
        localStorage.setItem('auth_token', data.value!.token)
        localStorage.setItem('auth_account', JSON.stringify(data.value!.tai_khoan))
      }
    } finally {
      state.value.loading = false
    }
  }

  async function dangXuat(): Promise<void> {
    if (state.value.token) {
      try {
        await useFetch('/api/v1/dang-xuat', {
          method: 'POST',
          headers: { Authorization: `Bearer ${state.value.token}` },
        })
      } catch { /* ignore */ }
    }
    clearAuth()
  }

  function clearAuth(): void {
    state.value.account = null
    state.value.token = null
    if (process.client) {
      localStorage.removeItem('auth_token')
      localStorage.removeItem('auth_account')
    }
  }

  async function kichHoatBan(): Promise<HOSoBan> {
    const { data, error } = await useFetch<HOSoBan>('/api/v1/ho-so-nguoi-ban', {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${state.value.token}`,
        'Content-Type': 'application/json',
      },
    })
    if (error.value) {
      throw new Error(error.value.data?.message || 'Kích hoạt thất bại')
    }
    return data.value!
  }

  const isLoggedIn = computed(() => !!state.value.token && !!state.value.account)
  const isAdmin = computed(() => state.value.account?.vai_tro === 'admin')

  return {
    state: readonly(state),
    isLoggedIn,
    isAdmin,
    dangKy,
    dangNhap,
    dangXuat,
    kichHoatBan,
    fetchMe,
  }
}
