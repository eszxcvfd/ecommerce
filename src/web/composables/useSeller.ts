import type { SanPhamSo } from './useCatalog'

export interface TepInput {
  ten_tep: string
  dinh_dang: string
  dung_luong_bytes: number
}

export interface DraftInput {
  ten: string
  mo_ta: string
  mo_ta_chi_tiet: string
  anh_demo: string
  mien_phi: boolean
  so_xu: number
  danh_muc: string
  giay_phep: string
  tep: TepInput[]
}

export interface DraftUpdateInput {
  ten?: string
  mo_ta?: string
  mo_ta_chi_tiet?: string
  anh_demo?: string
  mien_phi?: boolean
  so_xu?: number
  danh_muc?: string
  giay_phep?: string
  tep?: TepInput[]
}

export interface SanPhamResponse {
  san_pham: SanPhamSo
}

export interface SanPhamListResponse {
  san_pham: SanPhamSo[]
}

// Known file formats and their display names
export const ALLOWED_FORMATS: Record<string, string> = {
  dwg: 'AutoCAD',
  dxf: 'AutoCAD',
  skp: 'SketchUp',
  rvt: 'Revit',
  rfa: 'Revit',
  max: '3ds Max',
  '3ds': '3ds Max',
  psd: 'Photoshop',
  ai: 'Illustrator',
  eps: 'Illustrator',
}

export function useSeller() {
  const { state: authState } = useAuth()
  const token = computed(() => authState.value.token)

  const headers = () => ({
    'Content-Type': 'application/json',
    ...(token.value ? { Authorization: `Bearer ${token.value}` } : {}),
  })

  async function createDraft(input: DraftInput): Promise<SanPhamSo> {
    const res: SanPhamResponse = await $fetch('/api/v1/seller/san-pham', {
      method: 'POST',
      headers: headers(),
      body: input,
    })
    return res.san_pham
  }

  async function listDrafts(): Promise<SanPhamSo[]> {
    const res: SanPhamListResponse = await $fetch('/api/v1/seller/san-pham', {
      headers: headers(),
    })
    return res.san_pham
  }

  async function getDraft(id: string): Promise<SanPhamSo> {
    const res: SanPhamResponse = await $fetch(`/api/v1/seller/san-pham/${id}`, {
      headers: headers(),
    })
    return res.san_pham
  }

  async function updateDraft(id: string, input: DraftUpdateInput): Promise<SanPhamSo> {
    const res: SanPhamResponse = await $fetch(`/api/v1/seller/san-pham/${id}`, {
      method: 'PUT',
      headers: headers(),
      body: input,
    })
    return res.san_pham
  }

  async function deleteDraft(id: string): Promise<void> {
    await $fetch(`/api/v1/seller/san-pham/${id}`, {
      method: 'DELETE',
      headers: headers(),
    })
  }

  return {
    createDraft,
    listDrafts,
    getDraft,
    updateDraft,
    deleteDraft,
  }
}
