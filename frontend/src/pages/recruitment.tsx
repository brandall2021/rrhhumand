import { useEffect, useState } from 'react'
import {
  Plus,
  Pencil,
  Rocket,
  XCircle,
  CheckCircle2,
  Send,
  MoveRight,
  Search,
  LayoutDashboard,
  Briefcase,
  FileText,
  Users,
  ClipboardList,
  CalendarClock,
  Handshake,
} from 'lucide-react'
import api from '@/lib/api'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select } from '@/components/ui/select'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'

interface Position {
  id: string
  title: string
  department_id?: string | null
  location_id?: string | null
  employment_type?: string | null
  work_mode?: string | null
  description?: string | null
  requirements?: string | null
  responsibilities?: string | null
  benefits?: string | null
  salary_min?: number | null
  salary_max?: number | null
  currency?: string | null
  vacancies?: number
  vacancies_filled?: number
  status: string
  created_at?: string
}

interface Posting {
  id: string
  position_id?: string
  title: string
  description?: string
  requirements?: string | null
  employment_type?: string | null
  work_mode?: string | null
  location?: string | null
  salary_min?: number | null
  salary_max?: number | null
  currency?: string | null
  closing_at?: string | null
  is_public?: boolean
  external_url?: string | null
  status: string
  created_at?: string
}

interface Candidate {
  id: string
  first_name: string
  last_name: string
  email: string
  phone?: string | null
  phone_country_code?: string | null
  document_type?: string | null
  document_number?: string | null
  birth_date?: string | null
  location?: string | null
  nationality?: string | null
  gender?: string | null
  linkedin_url?: string | null
  portfolio_url?: string | null
  github_url?: string | null
  personal_website?: string | null
  current_company?: string | null
  current_position?: string | null
  notice_period?: number | null
  salary_expectation_min?: number | null
  salary_expectation_max?: number | null
  salary_currency?: string | null
  availability?: string | null
  source?: string | null
  source_detail?: string | null
  notes?: string | null
  status: string
  created_at?: string
}

interface Application {
  id: string
  candidate_id: string
  posting_id: string
  current_stage_id?: string | null
  status: string
  score?: number | null
  applied_at?: string
  created_at?: string
}

interface Interview {
  id: string
  application_id: string
  interview_type: string
  title?: string | null
  scheduled_at?: string | null
  duration_minutes?: number | null
  meeting_url?: string | null
  meeting_password?: string | null
  location?: string | null
  instructions?: string | null
  status: string
  notes?: string | null
  cancel_reason?: string | null
  created_at?: string
}

interface Offer {
  id: string
  application_id: string
  position_title: string
  department_id?: string | null
  offer_type: string
  start_date?: string | null
  employment_type?: string | null
  work_mode?: string | null
  salary_amount?: number | null
  salary_currency?: string | null
  salary_period?: string | null
  variable_compensation?: string | null
  benefits_summary?: string | null
  equity_terms?: string | null
  conditions?: string | null
  notes?: string | null
  response_deadline?: string | null
  status: string
  created_at?: string
}

interface DashboardStats {
  open_requisitions?: number
  total_candidates?: number
  applications_this_week?: number
  pending_offers?: number
  hires_this_month?: number
  total_interviews?: number
  avg_time_to_hire?: number
}

interface FunnelStage {
  stage: string
  count: number
}

interface SelectOption {
  value: string
  label: string
}

const statusBadge = (status: string) => {
  const map: Record<string, { label: string; cls: string }> = {
    DRAFT: { label: 'Borrador', cls: 'bg-slate-100 text-slate-600' },
    PENDING_APPROVAL: { label: 'Pendiente aprobación', cls: 'bg-amber-50 text-amber-700' },
    APPROVED: { label: 'Aprobado', cls: 'bg-emerald-50 text-emerald-700' },
    SENT: { label: 'Enviado', cls: 'bg-blue-50 text-blue-700' },
    ACCEPTED: { label: 'Aceptado', cls: 'bg-teal-50 text-teal-700' },
    REJECTED: { label: 'Rechazado', cls: 'bg-red-50 text-red-700' },
    WITHDRAWN: { label: 'Retirado', cls: 'bg-slate-100 text-slate-500' },
    PUBLISHED: { label: 'Publicado', cls: 'bg-emerald-50 text-emerald-700' },
    CLOSED: { label: 'Cerrado', cls: 'bg-slate-100 text-slate-600' },
    CANCELLED: { label: 'Cancelado', cls: 'bg-red-50 text-red-700' },
    ACTIVE: { label: 'Activo', cls: 'bg-emerald-50 text-emerald-700' },
    FILLED: { label: 'Completada', cls: 'bg-teal-50 text-teal-700' },
    INACTIVE: { label: 'Inactivo', cls: 'bg-slate-100 text-slate-600' },
    BLACKLISTED: { label: 'Bloqueado', cls: 'bg-red-50 text-red-700' },
    HIRED: { label: 'Contratado', cls: 'bg-teal-50 text-teal-700' },
    NEW: { label: 'Nueva', cls: 'bg-slate-100 text-slate-600' },
    SCREENING: { label: 'Screening', cls: 'bg-blue-50 text-blue-700' },
    INTERVIEW: { label: 'Entrevista', cls: 'bg-violet-50 text-violet-700' },
    ASSESSMENT: { label: 'Evaluación', cls: 'bg-amber-50 text-amber-700' },
    OFFER: { label: 'Oferta', cls: 'bg-emerald-50 text-emerald-700' },
    ON_HOLD: { label: 'En espera', cls: 'bg-slate-100 text-slate-500' },
    SCHEDULED: { label: 'Agendada', cls: 'bg-blue-50 text-blue-700' },
    CONFIRMED: { label: 'Confirmada', cls: 'bg-emerald-50 text-emerald-700' },
    IN_PROGRESS: { label: 'En curso', cls: 'bg-amber-50 text-amber-700' },
    COMPLETED: { label: 'Completada', cls: 'bg-teal-50 text-teal-700' },
    NO_SHOW: { label: 'No asistió', cls: 'bg-red-50 text-red-700' },
    RESCHEDULED: { label: 'Reprogramada', cls: 'bg-slate-100 text-slate-600' },
    EXPIRED: { label: 'Expirada', cls: 'bg-slate-100 text-slate-500' },
    NEGOTIATING: { label: 'Negociando', cls: 'bg-violet-50 text-violet-700' },
  }
  const s = map[status] ?? { label: status, cls: 'bg-slate-100 text-slate-600' }
  return <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${s.cls}`}>{s.label}</span>
}

const fmtDateTime = (s?: string | null) => (s ? s.slice(0, 16) : '-')
const fmtMoney = (v?: number | null) => (v == null ? '-' : v.toLocaleString())

const currencyOptions: SelectOption[] = [
  { value: 'ARS', label: 'ARS' },
  { value: 'USD', label: 'USD' },
  { value: 'EUR', label: 'EUR' },
]

const employmentTypeOptions: SelectOption[] = [
  { value: 'FULL_TIME', label: 'Tiempo completo' },
  { value: 'PART_TIME', label: 'Medio tiempo' },
  { value: 'CONTRACTOR', label: 'Contratista' },
  { value: 'INTERNSHIP', label: 'Pasantía' },
]

const workModeOptions: SelectOption[] = [
  { value: 'REMOTE', label: 'Remoto' },
  { value: 'ONSITE', label: 'Presencial' },
  { value: 'HYBRID', label: 'Híbrido' },
]

const salaryPeriodOptions: SelectOption[] = [
  { value: 'MONTHLY', label: 'Mensual' },
  { value: 'ANNUAL', label: 'Anual' },
  { value: 'HOURLY', label: 'Por hora' },
]

const emptyPosForm = {
  title: '',
  department_id: '',
  location_id: '',
  employment_type: '',
  work_mode: '',
  description: '',
  requirements: '',
  responsibilities: '',
  benefits: '',
  salary_min: '',
  salary_max: '',
  currency: 'ARS',
  vacancies: '1',
}

const emptyPostingForm = {
  position_id: '',
  title: '',
  description: '',
  requirements: '',
  employment_type: '',
  work_mode: '',
  location: '',
  salary_min: '',
  salary_max: '',
  currency: 'ARS',
  closing_at: '',
  is_public: false,
  external_url: '',
}

const emptyCandidateForm = {
  first_name: '',
  last_name: '',
  email: '',
  phone: '',
  phone_country_code: '',
  document_type: '',
  document_number: '',
  birth_date: '',
  location: '',
  nationality: '',
  gender: '',
  linkedin_url: '',
  portfolio_url: '',
  github_url: '',
  personal_website: '',
  current_company: '',
  current_position: '',
  notice_period: '',
  salary_expectation_min: '',
  salary_expectation_max: '',
  salary_currency: 'ARS',
  availability: '',
  source: '',
  source_detail: '',
  notes: '',
}

const emptyApplicationForm = {
  candidate_id: '',
  posting_id: '',
}

const emptyInterviewForm = {
  application_id: '',
  interview_type: '',
  title: '',
  scheduled_at: '',
  duration_minutes: '',
  meeting_url: '',
  meeting_password: '',
  location: '',
  instructions: '',
}

const emptyOfferForm = {
  application_id: '',
  position_title: '',
  department_id: '',
  offer_type: '',
  start_date: '',
  employment_type: '',
  work_mode: '',
  salary_amount: '',
  salary_currency: 'ARS',
  salary_period: '',
  variable_compensation: '',
  benefits_summary: '',
  equity_terms: '',
  conditions: '',
  notes: '',
  response_deadline: '',
}

const textareaCls =
  'w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-brand-500'

const tabs = [
  { key: 'dashboard', label: 'Dashboard', icon: LayoutDashboard },
  { key: 'positions', label: 'Posiciones', icon: Briefcase },
  { key: 'postings', label: 'Publicaciones', icon: FileText },
  { key: 'candidates', label: 'Candidatos', icon: Users },
  { key: 'applications', label: 'Postulaciones', icon: ClipboardList },
  { key: 'interviews', label: 'Entrevistas', icon: CalendarClock },
  { key: 'offers', label: 'Ofertas', icon: Handshake },
] as const

export default function RecruitmentPage() {
  const [error, setError] = useState('')

  // Lookups
  const [departmentOptions, setDepartmentOptions] = useState<SelectOption[]>([])
  const [positionOptions, setPositionOptions] = useState<SelectOption[]>([])
  const [postingOptions, setPostingOptions] = useState<SelectOption[]>([])
  const [candidateOptions, setCandidateOptions] = useState<SelectOption[]>([])
  const [stageOptions, setStageOptions] = useState<SelectOption[]>([])
  const [reasonOptions, setReasonOptions] = useState<SelectOption[]>([])
  const [departmentMap, setDepartmentMap] = useState<Record<string, string>>({})
  const [positionMap, setPositionMap] = useState<Record<string, string>>({})
  const [candidateMap, setCandidateMap] = useState<Record<string, string>>({})
  const [postingMap, setPostingMap] = useState<Record<string, string>>({})
  const [stageMap, setStageMap] = useState<Record<string, string>>({})
  const [appMetaMap, setAppMetaMap] = useState<Record<string, { candidate_id: string; posting_id: string }>>({})

  // Dashboard
  const [stats, setStats] = useState<DashboardStats>({})
  const [funnel, setFunnel] = useState<FunnelStage[]>([])
  const [dashLoading, setDashLoading] = useState(true)

  // Posiciones
  const [positions, setPositions] = useState<Position[]>([])
  const [posLoading, setPosLoading] = useState(true)
  const [showPosModal, setShowPosModal] = useState(false)
  const [editingPos, setEditingPos] = useState<Position | null>(null)
  const [posForm, setPosForm] = useState({ ...emptyPosForm })
  const [savingPos, setSavingPos] = useState(false)

  // Publicaciones
  const [postings, setPostings] = useState<Posting[]>([])
  const [postLoading, setPostLoading] = useState(true)
  const [showPostModal, setShowPostModal] = useState(false)
  const [editingPost, setEditingPost] = useState<Posting | null>(null)
  const [postForm, setPostForm] = useState({ ...emptyPostingForm })
  const [savingPost, setSavingPost] = useState(false)

  // Candidatos
  const [candidates, setCandidates] = useState<Candidate[]>([])
  const [candLoading, setCandLoading] = useState(true)
  const [showCandModal, setShowCandModal] = useState(false)
  const [editingCand, setEditingCand] = useState<Candidate | null>(null)
  const [candForm, setCandForm] = useState({ ...emptyCandidateForm })
  const [savingCand, setSavingCand] = useState(false)
  const [candQuery, setCandQuery] = useState('')

  // Postulaciones
  const [applications, setApplications] = useState<Application[]>([])
  const [appLoading, setAppLoading] = useState(true)
  const [showAppModal, setShowAppModal] = useState(false)
  const [appForm, setAppForm] = useState({ ...emptyApplicationForm })
  const [savingApp, setSavingApp] = useState(false)
  const [showMoveStage, setShowMoveStage] = useState(false)
  const [movingApp, setMovingApp] = useState<Application | null>(null)
  const [moveForm, setMoveForm] = useState({ to_stage_id: '', reason: '' })
  const [showRejectApp, setShowRejectApp] = useState(false)
  const [rejectingApp, setRejectingApp] = useState<Application | null>(null)
  const [rejectForm, setRejectForm] = useState({ reason_id: '', reason_text: '' })
  const [actionBusy, setActionBusy] = useState(false)

  // Entrevistas
  const [interviews, setInterviews] = useState<Interview[]>([])
  const [ivLoading, setIvLoading] = useState(true)
  const [showIvModal, setShowIvModal] = useState(false)
  const [editingIv, setEditingIv] = useState<Interview | null>(null)
  const [ivForm, setIvForm] = useState({ ...emptyInterviewForm })
  const [savingIv, setSavingIv] = useState(false)

  // Ofertas
  const [offers, setOffers] = useState<Offer[]>([])
  const [offerLoading, setOfferLoading] = useState(true)
  const [showOfferModal, setShowOfferModal] = useState(false)
  const [editingOffer, setEditingOffer] = useState<Offer | null>(null)
  const [offerForm, setOfferForm] = useState({ ...emptyOfferForm })
  const [savingOffer, setSavingOffer] = useState(false)

  const fetchLookups = async () => {
    const [dRes, pRes, postRes, cRes, sRes, rRes] = await Promise.allSettled([
      api.get('/departments', { params: { limit: '100' } }),
      api.get('/recruitment/positions', { params: { limit: '100' } }),
      api.get('/recruitment/postings', { params: { limit: '100' } }),
      api.get('/recruitment/candidates', { params: { limit: '100' } }),
      api.get('/recruitment/settings/stages'),
      api.get('/recruitment/settings/rejection-reasons'),
    ])
    if (dRes.status === 'fulfilled') {
      const list = dRes.value.data.data ?? []
      setDepartmentOptions(list.map((d: any) => ({ value: d.id, label: d.name })))
      const m: Record<string, string> = {}
      list.forEach((d: any) => { m[d.id] = d.name })
      setDepartmentMap(m)
    }
    if (pRes.status === 'fulfilled') {
      const list = pRes.value.data.data ?? []
      setPositionOptions(list.map((p: any) => ({ value: p.id, label: p.title })))
      const m: Record<string, string> = {}
      list.forEach((p: any) => { m[p.id] = p.title })
      setPositionMap(m)
    }
    if (postRes.status === 'fulfilled') {
      const list = postRes.value.data.data ?? []
      setPostingOptions(list.map((p: any) => ({ value: p.id, label: p.title })))
      const m: Record<string, string> = {}
      list.forEach((p: any) => { m[p.id] = p.title })
      setPostingMap(m)
    }
    if (cRes.status === 'fulfilled') {
      const list = cRes.value.data.data ?? []
      setCandidateOptions(list.map((c: any) => ({ value: c.id, label: `${c.first_name} ${c.last_name}`.trim() })))
      const m: Record<string, string> = {}
      list.forEach((c: any) => { m[c.id] = `${c.first_name} ${c.last_name}`.trim() })
      setCandidateMap(m)
    }
    if (sRes.status === 'fulfilled') {
      const list = sRes.value.data.data ?? []
      setStageOptions(list.map((s: any) => ({ value: s.id, label: s.name })))
      const m: Record<string, string> = {}
      list.forEach((s: any) => { m[s.id] = s.name })
      setStageMap(m)
    }
    if (rRes.status === 'fulfilled') {
      const list = rRes.value.data.data ?? []
      setReasonOptions(list.map((r: any) => ({ value: r.id, label: r.name })))
    }
  }

  const fetchDashboard = async () => {
    setDashLoading(true)
    try {
      const [sRes, fRes] = await Promise.all([
        api.get('/recruitment/dashboard'),
        api.get('/recruitment/dashboard/funnel'),
      ])
      setStats(sRes.data.data ?? {})
      const stages = fRes.data.data?.stages ?? []
      setFunnel(Array.isArray(stages) ? stages : [])
    } catch {
      setStats({})
      setFunnel([])
    } finally {
      setDashLoading(false)
    }
  }

  const fetchPositions = async () => {
    setPosLoading(true)
    try {
      const res = await api.get('/recruitment/positions', { params: { limit: '100' } })
      setPositions(res.data.data ?? [])
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar posiciones')
      setPositions([])
    } finally {
      setPosLoading(false)
    }
  }

  const fetchPostings = async () => {
    setPostLoading(true)
    try {
      const res = await api.get('/recruitment/postings', { params: { limit: '100' } })
      setPostings(res.data.data ?? [])
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar publicaciones')
      setPostings([])
    } finally {
      setPostLoading(false)
    }
  }

  const fetchCandidates = async () => {
    setCandLoading(true)
    try {
      const res = candQuery.trim()
        ? await api.get('/recruitment/candidates/search', { params: { skills: candQuery.trim() } })
        : await api.get('/recruitment/candidates', { params: { limit: '100' } })
      setCandidates(res.data.data ?? [])
      setError('')
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar candidatos')
      setCandidates([])
    } finally {
      setCandLoading(false)
    }
  }

  const fetchApplications = async () => {
    setAppLoading(true)
    try {
      const res = await api.get('/recruitment/applications', { params: { limit: '100' } })
      const list = res.data.data ?? []
      setApplications(list)
      const m: Record<string, { candidate_id: string; posting_id: string }> = {}
      list.forEach((a: any) => { m[a.id] = { candidate_id: a.candidate_id, posting_id: a.posting_id } })
      setAppMetaMap(m)
      setError('')
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar postulaciones')
      setApplications([])
    } finally {
      setAppLoading(false)
    }
  }

  const fetchInterviews = async () => {
    setIvLoading(true)
    try {
      const res = await api.get('/recruitment/interviews', { params: { limit: '100' } })
      setInterviews(res.data.data ?? [])
      setError('')
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar entrevistas')
      setInterviews([])
    } finally {
      setIvLoading(false)
    }
  }

  const fetchOffers = async () => {
    setOfferLoading(true)
    try {
      const res = await api.get('/recruitment/offers', { params: { limit: '100' } })
      setOffers(res.data.data ?? [])
      setError('')
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar ofertas')
      setOffers([])
    } finally {
      setOfferLoading(false)
    }
  }

  useEffect(() => { fetchDashboard() }, [])
  useEffect(() => {
    fetchLookups()
    fetchPositions()
    fetchPostings()
    fetchCandidates()
    fetchApplications()
    fetchInterviews()
    fetchOffers()
  }, [])
  useEffect(() => { fetchCandidates() }, [candQuery])

  const runAction = async (url: string, confirmMsg?: string) => {
    if (confirmMsg && !confirm(confirmMsg)) return
    try {
      await api.post(url)
      fetchPositions()
      fetchPostings()
      fetchInterviews()
      fetchOffers()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al ejecutar acción')
    }
  }

  const lookupName = (m: Record<string, string>, id?: string | null) => {
    if (!id) return '-'
    return m[id] ?? id.slice(0, 8)
  }

  const candidateFor = (applicationId?: string) => {
    if (!applicationId) return '-'
    const meta = appMetaMap[applicationId]
    if (!meta) return '-'
    return lookupName(candidateMap, meta.candidate_id)
  }

  const applicationOptions: SelectOption[] = Object.entries(appMetaMap).map(([id, meta]) => ({
    value: id,
    label: `${lookupName(candidateMap, meta.candidate_id)} - ${lookupName(postingMap, meta.posting_id)}`,
  }))

  // ---------- Posiciones ----------
  const openCreatePos = () => {
    setEditingPos(null)
    setPosForm({ ...emptyPosForm })
    setShowPosModal(true)
  }

  const openEditPos = (p: Position) => {
    setEditingPos(p)
    setPosForm({
      title: p.title ?? '',
      department_id: p.department_id ?? '',
      location_id: p.location_id ?? '',
      employment_type: p.employment_type ?? '',
      work_mode: p.work_mode ?? '',
      description: p.description ?? '',
      requirements: p.requirements ?? '',
      responsibilities: p.responsibilities ?? '',
      benefits: p.benefits ?? '',
      salary_min: p.salary_min != null ? String(p.salary_min) : '',
      salary_max: p.salary_max != null ? String(p.salary_max) : '',
      currency: p.currency ?? 'ARS',
      vacancies: p.vacancies != null ? String(p.vacancies) : '1',
    })
    setShowPosModal(true)
  }

  const handleSavePos = async () => {
    setSavingPos(true)
    try {
      const f = posForm
      const body: Record<string, any> = { title: f.title, vacancies: parseInt(f.vacancies || '1', 10) }
      if (f.department_id) body.department_id = f.department_id
      if (f.location_id) body.location_id = f.location_id
      if (f.employment_type) body.employment_type = f.employment_type
      if (f.work_mode) body.work_mode = f.work_mode
      if (f.description) body.description = f.description
      if (f.requirements) body.requirements = f.requirements
      if (f.responsibilities) body.responsibilities = f.responsibilities
      if (f.benefits) body.benefits = f.benefits
      if (f.salary_min) body.salary_min = parseFloat(f.salary_min)
      if (f.salary_max) body.salary_max = parseFloat(f.salary_max)
      if (f.currency) body.currency = f.currency
      if (editingPos) await api.put(`/recruitment/positions/${editingPos.id}`, body)
      else await api.post('/recruitment/positions', body)
      setShowPosModal(false)
      fetchPositions()
      fetchLookups()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar posición')
    } finally {
      setSavingPos(false)
    }
  }

  // ---------- Publicaciones ----------
  const openCreatePost = () => {
    setEditingPost(null)
    setPostForm({ ...emptyPostingForm })
    setShowPostModal(true)
  }

  const openEditPost = (p: Posting) => {
    setEditingPost(p)
    setPostForm({
      position_id: p.position_id ?? '',
      title: p.title ?? '',
      description: p.description ?? '',
      requirements: p.requirements ?? '',
      employment_type: p.employment_type ?? '',
      work_mode: p.work_mode ?? '',
      location: p.location ?? '',
      salary_min: p.salary_min != null ? String(p.salary_min) : '',
      salary_max: p.salary_max != null ? String(p.salary_max) : '',
      currency: p.currency ?? 'ARS',
      closing_at: p.closing_at ? p.closing_at.slice(0, 10) : '',
      is_public: p.is_public ?? false,
      external_url: p.external_url ?? '',
    })
    setShowPostModal(true)
  }

  const handleSavePost = async () => {
    setSavingPost(true)
    try {
      const f = postForm
      const body: Record<string, any> = { title: f.title, is_public: f.is_public }
      if (f.position_id) body.position_id = f.position_id
      if (f.description) body.description = f.description
      if (f.requirements) body.requirements = f.requirements
      if (f.employment_type) body.employment_type = f.employment_type
      if (f.work_mode) body.work_mode = f.work_mode
      if (f.location) body.location = f.location
      if (f.salary_min) body.salary_min = parseFloat(f.salary_min)
      if (f.salary_max) body.salary_max = parseFloat(f.salary_max)
      if (f.currency) body.currency = f.currency
      if (f.closing_at) body.closing_at = f.closing_at + 'T00:00:00Z'
      if (f.external_url) body.external_url = f.external_url
      if (editingPost) await api.put(`/recruitment/postings/${editingPost.id}`, body)
      else await api.post('/recruitment/postings', body)
      setShowPostModal(false)
      fetchPostings()
      fetchLookups()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar publicación')
    } finally {
      setSavingPost(false)
    }
  }

  // ---------- Candidatos ----------
  const openCreateCand = () => {
    setEditingCand(null)
    setCandForm({ ...emptyCandidateForm })
    setShowCandModal(true)
  }

  const openEditCand = (c: Candidate) => {
    setEditingCand(c)
    setCandForm({
      first_name: c.first_name ?? '',
      last_name: c.last_name ?? '',
      email: c.email ?? '',
      phone: c.phone ?? '',
      phone_country_code: c.phone_country_code ?? '',
      document_type: c.document_type ?? '',
      document_number: c.document_number ?? '',
      birth_date: c.birth_date ? c.birth_date.slice(0, 10) : '',
      location: c.location ?? '',
      nationality: c.nationality ?? '',
      gender: c.gender ?? '',
      linkedin_url: c.linkedin_url ?? '',
      portfolio_url: c.portfolio_url ?? '',
      github_url: c.github_url ?? '',
      personal_website: c.personal_website ?? '',
      current_company: c.current_company ?? '',
      current_position: c.current_position ?? '',
      notice_period: c.notice_period != null ? String(c.notice_period) : '',
      salary_expectation_min: c.salary_expectation_min != null ? String(c.salary_expectation_min) : '',
      salary_expectation_max: c.salary_expectation_max != null ? String(c.salary_expectation_max) : '',
      salary_currency: c.salary_currency ?? 'ARS',
      availability: c.availability ?? '',
      source: c.source ?? '',
      source_detail: c.source_detail ?? '',
      notes: c.notes ?? '',
    })
    setShowCandModal(true)
  }

  const handleSaveCand = async () => {
    setSavingCand(true)
    try {
      const f = candForm
      const body: Record<string, any> = { first_name: f.first_name, last_name: f.last_name, email: f.email }
      if (f.phone) body.phone = f.phone
      if (f.phone_country_code) body.phone_country_code = f.phone_country_code
      if (f.document_type) body.document_type = f.document_type
      if (f.document_number) body.document_number = f.document_number
      if (f.birth_date) body.birth_date = f.birth_date + 'T00:00:00Z'
      if (f.location) body.location = f.location
      if (f.nationality) body.nationality = f.nationality
      if (f.gender) body.gender = f.gender
      if (f.linkedin_url) body.linkedin_url = f.linkedin_url
      if (f.portfolio_url) body.portfolio_url = f.portfolio_url
      if (f.github_url) body.github_url = f.github_url
      if (f.personal_website) body.personal_website = f.personal_website
      if (f.current_company) body.current_company = f.current_company
      if (f.current_position) body.current_position = f.current_position
      if (f.notice_period) body.notice_period = parseInt(f.notice_period, 10)
      if (f.salary_expectation_min) body.salary_expectation_min = parseFloat(f.salary_expectation_min)
      if (f.salary_expectation_max) body.salary_expectation_max = parseFloat(f.salary_expectation_max)
      if (f.salary_currency) body.salary_currency = f.salary_currency
      if (f.availability) body.availability = f.availability
      if (f.source) body.source = f.source
      if (f.source_detail) body.source_detail = f.source_detail
      if (f.notes) body.notes = f.notes
      if (editingCand) await api.put(`/recruitment/candidates/${editingCand.id}`, body)
      else await api.post('/recruitment/candidates', body)
      setShowCandModal(false)
      fetchCandidates()
      fetchLookups()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar candidato')
    } finally {
      setSavingCand(false)
    }
  }

  // ---------- Postulaciones ----------
  const openCreateApp = () => {
    setAppForm({ ...emptyApplicationForm })
    setShowAppModal(true)
  }

  const handleSaveApp = async () => {
    setSavingApp(true)
    try {
      const body: Record<string, any> = { candidate_id: appForm.candidate_id, posting_id: appForm.posting_id }
      await api.post('/recruitment/applications', body)
      setShowAppModal(false)
      fetchApplications()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al crear postulación')
    } finally {
      setSavingApp(false)
    }
  }

  const openMoveStage = (a: Application) => {
    setMovingApp(a)
    setMoveForm({ to_stage_id: a.current_stage_id ?? '', reason: '' })
    setShowMoveStage(true)
  }

  const handleMoveStage = async () => {
    if (!movingApp || !moveForm.to_stage_id) return
    setActionBusy(true)
    try {
      const body: Record<string, any> = { to_stage_id: moveForm.to_stage_id }
      if (moveForm.reason) body.reason = moveForm.reason
      await api.post(`/recruitment/applications/${movingApp.id}/stage`, body)
      setShowMoveStage(false)
      fetchApplications()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al mover etapa')
    } finally {
      setActionBusy(false)
    }
  }

  const openRejectApp = (a: Application) => {
    setRejectingApp(a)
    setRejectForm({ reason_id: '', reason_text: '' })
    setShowRejectApp(true)
  }

  const handleRejectApp = async () => {
    if (!rejectingApp) return
    setActionBusy(true)
    try {
      const body: Record<string, any> = {}
      if (rejectForm.reason_id) body.reason_id = rejectForm.reason_id
      if (rejectForm.reason_text) body.reason_text = rejectForm.reason_text
      await api.post(`/recruitment/applications/${rejectingApp.id}/reject`, body)
      setShowRejectApp(false)
      fetchApplications()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al rechazar postulación')
    } finally {
      setActionBusy(false)
    }
  }

  // ---------- Entrevistas ----------
  const openCreateIv = () => {
    setEditingIv(null)
    setIvForm({ ...emptyInterviewForm })
    setShowIvModal(true)
  }

  const openEditIv = (i: Interview) => {
    setEditingIv(i)
    setIvForm({
      application_id: i.application_id ?? '',
      interview_type: i.interview_type ?? '',
      title: i.title ?? '',
      scheduled_at: i.scheduled_at ? i.scheduled_at.slice(0, 16) : '',
      duration_minutes: i.duration_minutes != null ? String(i.duration_minutes) : '',
      meeting_url: i.meeting_url ?? '',
      meeting_password: i.meeting_password ?? '',
      location: i.location ?? '',
      instructions: i.instructions ?? '',
    })
    setShowIvModal(true)
  }

  const handleSaveIv = async () => {
    setSavingIv(true)
    try {
      const f = ivForm
      const body: Record<string, any> = {
        application_id: f.application_id,
        interview_type: f.interview_type,
      }
      if (f.title) body.title = f.title
      if (f.scheduled_at) body.scheduled_at = f.scheduled_at + ':00Z'
      if (f.duration_minutes) body.duration_minutes = parseInt(f.duration_minutes, 10)
      if (f.meeting_url) body.meeting_url = f.meeting_url
      if (f.meeting_password) body.meeting_password = f.meeting_password
      if (f.location) body.location = f.location
      if (f.instructions) body.instructions = f.instructions
      if (editingIv) await api.put(`/recruitment/interviews/${editingIv.id}`, body)
      else await api.post('/recruitment/interviews', body)
      setShowIvModal(false)
      fetchInterviews()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar entrevista')
    } finally {
      setSavingIv(false)
    }
  }

  // ---------- Ofertas ----------
  const openCreateOffer = () => {
    setEditingOffer(null)
    setOfferForm({ ...emptyOfferForm })
    setShowOfferModal(true)
  }

  const openEditOffer = (o: Offer) => {
    setEditingOffer(o)
    setOfferForm({
      application_id: o.application_id ?? '',
      position_title: o.position_title ?? '',
      department_id: o.department_id ?? '',
      offer_type: o.offer_type ?? '',
      start_date: o.start_date ? o.start_date.slice(0, 10) : '',
      employment_type: o.employment_type ?? '',
      work_mode: o.work_mode ?? '',
      salary_amount: o.salary_amount != null ? String(o.salary_amount) : '',
      salary_currency: o.salary_currency ?? 'ARS',
      salary_period: o.salary_period ?? '',
      variable_compensation: o.variable_compensation ?? '',
      benefits_summary: o.benefits_summary ?? '',
      equity_terms: o.equity_terms ?? '',
      conditions: o.conditions ?? '',
      notes: o.notes ?? '',
      response_deadline: o.response_deadline ? o.response_deadline.slice(0, 10) : '',
    })
    setShowOfferModal(true)
  }

  const handleSaveOffer = async () => {
    setSavingOffer(true)
    try {
      const f = offerForm
      const body: Record<string, any> = {
        application_id: f.application_id,
        position_title: f.position_title,
        offer_type: f.offer_type,
      }
      if (f.department_id) body.department_id = f.department_id
      if (f.start_date) body.start_date = f.start_date + 'T00:00:00Z'
      if (f.employment_type) body.employment_type = f.employment_type
      if (f.work_mode) body.work_mode = f.work_mode
      if (f.salary_amount) body.salary_amount = parseFloat(f.salary_amount)
      if (f.salary_currency) body.salary_currency = f.salary_currency
      if (f.salary_period) body.salary_period = f.salary_period
      if (f.variable_compensation) body.variable_compensation = f.variable_compensation
      if (f.benefits_summary) body.benefits_summary = f.benefits_summary
      if (f.equity_terms) body.equity_terms = f.equity_terms
      if (f.conditions) body.conditions = f.conditions
      if (f.notes) body.notes = f.notes
      if (f.response_deadline) body.response_deadline = f.response_deadline + 'T00:00:00Z'
      if (editingOffer) await api.put(`/recruitment/offers/${editingOffer.id}`, body)
      else await api.post('/recruitment/offers', body)
      setShowOfferModal(false)
      fetchOffers()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar oferta')
    } finally {
      setSavingOffer(false)
    }
  }

  const offerActions = (o: Offer) => {
    switch (o.status) {
      case 'DRAFT':
        return (
          <>
            <Button variant="ghost" size="sm" className="text-amber-600" title="Enviar a aprobación" onClick={() => runAction(`/recruitment/offers/${o.id}/submit`)}><Send size={14} /></Button>
            <Button variant="ghost" size="sm" title="Editar" onClick={() => openEditOffer(o)}><Pencil size={14} /></Button>
          </>
        )
      case 'PENDING_APPROVAL':
        return (
          <>
            <Button variant="ghost" size="sm" className="text-emerald-600" title="Aprobar" onClick={() => runAction(`/recruitment/offers/${o.id}/approve`)}><CheckCircle2 size={14} /></Button>
            <Button variant="ghost" size="sm" className="text-red-500" title="Rechazar" onClick={() => runAction(`/recruitment/offers/${o.id}/reject`, '¿Rechazar la oferta?')}><XCircle size={14} /></Button>
          </>
        )
      case 'APPROVED':
        return (
          <>
            <Button variant="ghost" size="sm" className="text-blue-600" title="Enviar oferta" onClick={() => runAction(`/recruitment/offers/${o.id}/send`)}><Rocket size={14} /></Button>
            <Button variant="ghost" size="sm" className="text-slate-500" title="Retirar" onClick={() => runAction(`/recruitment/offers/${o.id}/withdraw`, '¿Retirar la oferta?')}><XCircle size={14} /></Button>
          </>
        )
      case 'SENT':
        return (
          <>
            <Button variant="ghost" size="sm" className="text-emerald-600" title="Aceptar oferta" onClick={() => runAction(`/recruitment/offers/${o.id}/accept`, '¿Confirmar aceptación de la oferta?')}><CheckCircle2 size={14} /></Button>
            <Button variant="ghost" size="sm" className="text-red-500" title="Rechazar oferta" onClick={() => runAction(`/recruitment/offers/${o.id}/reject`, '¿Rechazar la oferta?')}><XCircle size={14} /></Button>
            <Button variant="ghost" size="sm" className="text-slate-500" title="Retirar" onClick={() => runAction(`/recruitment/offers/${o.id}/withdraw`, '¿Retirar la oferta?')}><XCircle size={14} /></Button>
          </>
        )
      default:
        return null
    }
  }

  const funnelMax = funnel.reduce((m, s) => Math.max(m, s.count), 0) || 1

  const statCards: { key: keyof DashboardStats; label: string; value?: string | number }[] = [
    { key: 'open_requisitions', label: 'Requisiciones abiertas' },
    { key: 'total_candidates', label: 'Candidatos' },
    { key: 'applications_this_week', label: 'Postulaciones (semana)' },
    { key: 'pending_offers', label: 'Ofertas pendientes' },
    { key: 'hires_this_month', label: 'Contrataciones (mes)' },
    { key: 'total_interviews', label: 'Entrevistas' },
  ]

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-slate-900">Reclutamiento</h1>
      </div>

      {error && <div className="mb-4 p-3 bg-red-50 text-red-700 text-sm rounded-lg">{error}</div>}

      <Tabs defaultValue="dashboard">
        <TabsList className="flex w-full justify-start overflow-x-auto">
          {tabs.map(t => (
            <TabsTrigger key={t.key} value={t.key} className="gap-1.5">
              <t.icon size={15} /> {t.label}
            </TabsTrigger>
          ))}
        </TabsList>

        {/* ---------- Dashboard ---------- */}
        <TabsContent value="dashboard">
          {dashLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div> : (
            <>
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                {statCards.map(card => (
                  <Card key={card.key}>
                    <CardContent className="p-6">
                      <p className="text-sm font-medium text-slate-500">{card.label}</p>
                      <p className="text-3xl font-bold text-slate-900 mt-1">
                        {stats[card.key] != null ? Number(stats[card.key]).toLocaleString() : '-'}
                      </p>
                    </CardContent>
                  </Card>
                ))}
              </div>

              <Card className="mt-4">
                <CardContent className="p-6">
                  <h2 className="text-sm font-medium text-slate-500 mb-4">Funnel de postulaciones</h2>
                  {funnel.length === 0 ? (
                    <p className="text-sm text-slate-500">Sin datos de funnel disponibles</p>
                  ) : (
                    <div className="space-y-3">
                      {funnel.map((s, idx) => (
                        <div key={`${s.stage}-${idx}`} className="flex items-center gap-3">
                          <span className="w-40 text-sm text-slate-600 truncate">{s.stage}</span>
                          <div className="flex-1 h-6 bg-slate-100 rounded-md overflow-hidden">
                            <div
                              className="h-full bg-brand-500 rounded-md transition-all"
                              style={{ width: `${Math.max(4, (s.count / funnelMax) * 100)}%` }}
                            />
                          </div>
                          <span className="w-10 text-right text-sm font-medium text-slate-900">{s.count}</span>
                        </div>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>
            </>
          )}
        </TabsContent>

        {/* ---------- Posiciones ---------- */}
        <TabsContent value="positions">
          <div className="flex items-center justify-between mb-4">
            <span className="text-sm text-slate-500">{positions.length} posiciones</span>
            <Button size="sm" onClick={openCreatePos}><Plus size={16} className="mr-1" /> Nueva posición</Button>
          </div>

          <Card>
            <CardContent className="p-0">
              {posLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
              : positions.length === 0 ? <div className="p-6 text-center text-slate-500">No hay posiciones</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Título</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Departamento</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Tipo</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Salario</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Vacantes</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                        <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                      </tr>
                    </thead>
                    <tbody>
                      {positions.map(p => (
                        <tr key={p.id} className="border-b border-slate-100 hover:bg-slate-50">
                          <td className="px-4 py-3 font-medium text-slate-900">{p.title}</td>
                          <td className="px-4 py-3 text-slate-600">{lookupName(departmentMap, p.department_id)}</td>
                          <td className="px-4 py-3 text-slate-600">{p.employment_type || '-'}</td>
                          <td className="px-4 py-3 text-slate-600">
                            {p.salary_min != null || p.salary_max != null
                              ? `${fmtMoney(p.salary_min)} - ${fmtMoney(p.salary_max)} ${p.currency || ''}`
                              : '-'}
                          </td>
                          <td className="px-4 py-3 text-slate-600">{p.vacancies ?? 0}</td>
                          <td className="px-4 py-3">{statusBadge(p.status)}</td>
                          <td className="px-4 py-3 text-right whitespace-nowrap">
                            <Button variant="ghost" size="sm" title="Editar" onClick={() => openEditPos(p)}><Pencil size={14} /></Button>
                            {p.status === 'ACTIVE' && (
                              <Button variant="ghost" size="sm" className="text-slate-500" title="Cerrar" onClick={() => runAction(`/recruitment/positions/${p.id}/close`, `¿Cerrar la posición "${p.title}"?`)}><XCircle size={14} /></Button>
                            )}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>}
            </CardContent>
          </Card>
        </TabsContent>

        {/* ---------- Publicaciones ---------- */}
        <TabsContent value="postings">
          <div className="flex items-center justify-between mb-4">
            <span className="text-sm text-slate-500">{postings.length} publicaciones</span>
            <Button size="sm" onClick={openCreatePost}><Plus size={16} className="mr-1" /> Nueva publicación</Button>
          </div>

          <Card>
            <CardContent className="p-0">
              {postLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
              : postings.length === 0 ? <div className="p-6 text-center text-slate-500">No hay publicaciones</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Título</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Posición</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Ubicación</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Salario</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                        <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                      </tr>
                    </thead>
                    <tbody>
                      {postings.map(p => (
                        <tr key={p.id} className="border-b border-slate-100 hover:bg-slate-50">
                          <td className="px-4 py-3 font-medium text-slate-900">{p.title}</td>
                          <td className="px-4 py-3 text-slate-600">{lookupName(positionMap, p.position_id)}</td>
                          <td className="px-4 py-3 text-slate-600">{p.location || '-'}</td>
                          <td className="px-4 py-3 text-slate-600">
                            {p.salary_min != null || p.salary_max != null
                              ? `${fmtMoney(p.salary_min)} - ${fmtMoney(p.salary_max)} ${p.currency || ''}`
                              : '-'}
                          </td>
                          <td className="px-4 py-3">{statusBadge(p.status)}</td>
                          <td className="px-4 py-3 text-right whitespace-nowrap">
                            <Button variant="ghost" size="sm" title="Editar" onClick={() => openEditPost(p)}><Pencil size={14} /></Button>
                            {p.status === 'DRAFT' && (
                              <Button variant="ghost" size="sm" className="text-emerald-600" title="Publicar" onClick={() => runAction(`/recruitment/postings/${p.id}/publish`)}><Rocket size={14} /></Button>
                            )}
                            {p.status === 'PUBLISHED' && (
                              <Button variant="ghost" size="sm" className="text-slate-500" title="Cerrar" onClick={() => runAction(`/recruitment/postings/${p.id}/close`, `¿Cerrar la publicación "${p.title}"?`)}><XCircle size={14} /></Button>
                            )}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>}
            </CardContent>
          </Card>
        </TabsContent>

        {/* ---------- Candidatos ---------- */}
        <TabsContent value="candidates">
          <div className="flex items-center gap-3 mb-4">
            <div className="relative w-72">
              <Search size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
              <Input
                className="pl-9"
                placeholder="Buscar por habilidades (ej: react, go)"
                value={candQuery}
                onChange={e => setCandQuery(e.target.value)}
              />
            </div>
            <div className="ml-auto">
              <Button size="sm" onClick={openCreateCand}><Plus size={16} className="mr-1" /> Nuevo candidato</Button>
            </div>
          </div>

          <Card>
            <CardContent className="p-0">
              {candLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
              : candidates.length === 0 ? <div className="p-6 text-center text-slate-500">No hay candidatos</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Email</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Empresa actual</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Posición actual</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                        <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                      </tr>
                    </thead>
                    <tbody>
                      {candidates.map(c => (
                        <tr key={c.id} className="border-b border-slate-100 hover:bg-slate-50">
                          <td className="px-4 py-3 font-medium text-slate-900">{c.first_name} {c.last_name}</td>
                          <td className="px-4 py-3 text-slate-600">{c.email}</td>
                          <td className="px-4 py-3 text-slate-600">{c.current_company || '-'}</td>
                          <td className="px-4 py-3 text-slate-600">{c.current_position || '-'}</td>
                          <td className="px-4 py-3">{statusBadge(c.status)}</td>
                          <td className="px-4 py-3 text-right whitespace-nowrap">
                            <Button variant="ghost" size="sm" title="Editar" onClick={() => openEditCand(c)}><Pencil size={14} /></Button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>}
            </CardContent>
          </Card>
        </TabsContent>

        {/* ---------- Postulaciones ---------- */}
        <TabsContent value="applications">
          <div className="flex items-center justify-between mb-4">
            <span className="text-sm text-slate-500">{applications.length} postulaciones</span>
            <Button size="sm" onClick={openCreateApp}><Plus size={16} className="mr-1" /> Nueva postulación</Button>
          </div>

          <Card>
            <CardContent className="p-0">
              {appLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
              : applications.length === 0 ? <div className="p-6 text-center text-slate-500">No hay postulaciones</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Candidato</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Publicación</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Etapa</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Score</th>
                        <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                      </tr>
                    </thead>
                    <tbody>
                      {applications.map(a => (
                        <tr key={a.id} className="border-b border-slate-100 hover:bg-slate-50">
                          <td className="px-4 py-3 font-medium text-slate-900">{lookupName(candidateMap, a.candidate_id)}</td>
                          <td className="px-4 py-3 text-slate-600">{lookupName(postingMap, a.posting_id)}</td>
                          <td className="px-4 py-3 text-slate-600">{lookupName(stageMap, a.current_stage_id)}</td>
                          <td className="px-4 py-3">{statusBadge(a.status)}</td>
                          <td className="px-4 py-3 text-slate-600">{a.score != null ? a.score : '-'}</td>
                          <td className="px-4 py-3 text-right whitespace-nowrap">
                            {!['REJECTED', 'WITHDRAWN', 'HIRED'].includes(a.status) && (
                              <Button variant="ghost" size="sm" className="text-blue-600" title="Mover etapa" onClick={() => openMoveStage(a)}><MoveRight size={14} /></Button>
                            )}
                            {!['REJECTED', 'WITHDRAWN', 'HIRED'].includes(a.status) && (
                              <Button variant="ghost" size="sm" className="text-red-500" title="Rechazar" onClick={() => openRejectApp(a)}><XCircle size={14} /></Button>
                            )}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>}
            </CardContent>
          </Card>
        </TabsContent>

        {/* ---------- Entrevistas ---------- */}
        <TabsContent value="interviews">
          <div className="flex items-center justify-between mb-4">
            <span className="text-sm text-slate-500">{interviews.length} entrevistas</span>
            <Button size="sm" onClick={openCreateIv}><Plus size={16} className="mr-1" /> Nueva entrevista</Button>
          </div>

          <Card>
            <CardContent className="p-0">
              {ivLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
              : interviews.length === 0 ? <div className="p-6 text-center text-slate-500">No hay entrevistas</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Candidato</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Tipo</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Fecha</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Duración</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Ubicación / Link</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                        <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                      </tr>
                    </thead>
                    <tbody>
                      {interviews.map(i => (
                        <tr key={i.id} className="border-b border-slate-100 hover:bg-slate-50">
                          <td className="px-4 py-3 font-medium text-slate-900">{candidateFor(i.application_id)}</td>
                          <td className="px-4 py-3 text-slate-600">{i.interview_type || '-'}</td>
                          <td className="px-4 py-3 text-slate-600">{fmtDateTime(i.scheduled_at)}</td>
                          <td className="px-4 py-3 text-slate-600">{i.duration_minutes != null ? `${i.duration_minutes} min` : '-'}</td>
                          <td className="px-4 py-3 text-slate-600">{i.location || i.meeting_url || '-'}</td>
                          <td className="px-4 py-3">{statusBadge(i.status)}</td>
                          <td className="px-4 py-3 text-right whitespace-nowrap">
                            {['SCHEDULED', 'CONFIRMED'].includes(i.status) && (
                              <Button variant="ghost" size="sm" title="Editar" onClick={() => openEditIv(i)}><Pencil size={14} /></Button>
                            )}
                            {['SCHEDULED', 'CONFIRMED', 'IN_PROGRESS'].includes(i.status) && (
                              <Button variant="ghost" size="sm" className="text-teal-600" title="Completar" onClick={() => runAction(`/recruitment/interviews/${i.id}/complete`)}><CheckCircle2 size={14} /></Button>
                            )}
                            {['SCHEDULED', 'CONFIRMED'].includes(i.status) && (
                              <Button variant="ghost" size="sm" className="text-red-500" title="Cancelar" onClick={() => runAction(`/recruitment/interviews/${i.id}/cancel`, '¿Cancelar la entrevista?')}><XCircle size={14} /></Button>
                            )}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>}
            </CardContent>
          </Card>
        </TabsContent>

        {/* ---------- Ofertas ---------- */}
        <TabsContent value="offers">
          <div className="flex items-center justify-between mb-4">
            <span className="text-sm text-slate-500">{offers.length} ofertas</span>
            <Button size="sm" onClick={openCreateOffer}><Plus size={16} className="mr-1" /> Nueva oferta</Button>
          </div>

          <Card>
            <CardContent className="p-0">
              {offerLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
              : offers.length === 0 ? <div className="p-6 text-center text-slate-500">No hay ofertas</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Posición</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Candidato</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Salario</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                        <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                      </tr>
                    </thead>
                    <tbody>
                      {offers.map(o => (
                        <tr key={o.id} className="border-b border-slate-100 hover:bg-slate-50">
                          <td className="px-4 py-3 font-medium text-slate-900">{o.position_title}</td>
                          <td className="px-4 py-3 text-slate-600">{candidateFor(o.application_id)}</td>
                          <td className="px-4 py-3 text-slate-600">
                            {o.salary_amount != null ? `${fmtMoney(o.salary_amount)} ${o.salary_currency || ''}` : '-'}
                          </td>
                          <td className="px-4 py-3">{statusBadge(o.status)}</td>
                          <td className="px-4 py-3 text-right whitespace-nowrap">
                            {offerActions(o)}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      {/* ---------- Modal Posición ---------- */}
      <Dialog open={showPosModal} onOpenChange={setShowPosModal}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>{editingPos ? 'Editar Posición' : 'Nueva Posición'}</DialogTitle>
            <DialogDescription>Completá los datos de la posición</DialogDescription>
          </DialogHeader>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2 col-span-2">
              <Label htmlFor="pos-title">Título *</Label>
              <Input id="pos-title" value={posForm.title} onChange={e => setPosForm({ ...posForm, title: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="pos-dept">Departamento</Label>
              <Select id="pos-dept" options={departmentOptions} placeholder="Seleccionar..." value={posForm.department_id} onChange={e => setPosForm({ ...posForm, department_id: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="pos-loc">Ubicación (ID)</Label>
              <Input id="pos-loc" value={posForm.location_id} onChange={e => setPosForm({ ...posForm, location_id: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="pos-emp">Tipo de contratación</Label>
              <Select id="pos-emp" options={employmentTypeOptions} placeholder="Seleccionar..." value={posForm.employment_type} onChange={e => setPosForm({ ...posForm, employment_type: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="pos-mode">Modalidad</Label>
              <Select id="pos-mode" options={workModeOptions} placeholder="Seleccionar..." value={posForm.work_mode} onChange={e => setPosForm({ ...posForm, work_mode: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="pos-min">Salario mínimo</Label>
              <Input id="pos-min" type="number" step="0.01" value={posForm.salary_min} onChange={e => setPosForm({ ...posForm, salary_min: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="pos-max">Salario máximo</Label>
              <Input id="pos-max" type="number" step="0.01" value={posForm.salary_max} onChange={e => setPosForm({ ...posForm, salary_max: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="pos-currency">Moneda</Label>
              <Select id="pos-currency" options={currencyOptions} value={posForm.currency} onChange={e => setPosForm({ ...posForm, currency: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="pos-vac">Vacantes</Label>
              <Input id="pos-vac" type="number" min="1" value={posForm.vacancies} onChange={e => setPosForm({ ...posForm, vacancies: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="pos-desc">Descripción</Label>
              <textarea id="pos-desc" value={posForm.description} onChange={e => setPosForm({ ...posForm, description: e.target.value })} rows={2} className={textareaCls} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="pos-req">Requerimientos</Label>
              <textarea id="pos-req" value={posForm.requirements} onChange={e => setPosForm({ ...posForm, requirements: e.target.value })} rows={2} className={textareaCls} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="pos-resp">Responsabilidades</Label>
              <textarea id="pos-resp" value={posForm.responsibilities} onChange={e => setPosForm({ ...posForm, responsibilities: e.target.value })} rows={2} className={textareaCls} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="pos-ben">Beneficios</Label>
              <textarea id="pos-ben" value={posForm.benefits} onChange={e => setPosForm({ ...posForm, benefits: e.target.value })} rows={2} className={textareaCls} />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowPosModal(false)}>Cancelar</Button>
            <Button onClick={handleSavePos} disabled={savingPos || !posForm.title}>
              {savingPos ? 'Guardando...' : editingPos ? 'Guardar Cambios' : 'Crear'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* ---------- Modal Publicación ---------- */}
      <Dialog open={showPostModal} onOpenChange={setShowPostModal}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>{editingPost ? 'Editar Publicación' : 'Nueva Publicación'}</DialogTitle>
            <DialogDescription>Completá los datos de la publicación</DialogDescription>
          </DialogHeader>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2 col-span-2">
              <Label htmlFor="post-title">Título *</Label>
              <Input id="post-title" value={postForm.title} onChange={e => setPostForm({ ...postForm, title: e.target.value })} required />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="post-pos">Posición</Label>
              <Select id="post-pos" options={positionOptions} placeholder="Seleccionar..." value={postForm.position_id} onChange={e => setPostForm({ ...postForm, position_id: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="post-desc">Descripción</Label>
              <textarea id="post-desc" value={postForm.description} onChange={e => setPostForm({ ...postForm, description: e.target.value })} rows={2} className={textareaCls} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="post-req">Requerimientos</Label>
              <textarea id="post-req" value={postForm.requirements} onChange={e => setPostForm({ ...postForm, requirements: e.target.value })} rows={2} className={textareaCls} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="post-emp">Tipo de contratación</Label>
              <Select id="post-emp" options={employmentTypeOptions} placeholder="Seleccionar..." value={postForm.employment_type} onChange={e => setPostForm({ ...postForm, employment_type: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="post-mode">Modalidad</Label>
              <Select id="post-mode" options={workModeOptions} placeholder="Seleccionar..." value={postForm.work_mode} onChange={e => setPostForm({ ...postForm, work_mode: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="post-loc">Ubicación</Label>
              <Input id="post-loc" value={postForm.location} onChange={e => setPostForm({ ...postForm, location: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="post-close">Cierre de publicación</Label>
              <Input id="post-close" type="date" value={postForm.closing_at} onChange={e => setPostForm({ ...postForm, closing_at: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="post-min">Salario mínimo</Label>
              <Input id="post-min" type="number" step="0.01" value={postForm.salary_min} onChange={e => setPostForm({ ...postForm, salary_min: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="post-max">Salario máximo</Label>
              <Input id="post-max" type="number" step="0.01" value={postForm.salary_max} onChange={e => setPostForm({ ...postForm, salary_max: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="post-currency">Moneda</Label>
              <Select id="post-currency" options={currencyOptions} value={postForm.currency} onChange={e => setPostForm({ ...postForm, currency: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="post-url">URL externa</Label>
              <Input id="post-url" value={postForm.external_url} onChange={e => setPostForm({ ...postForm, external_url: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <label className="flex items-center gap-2 text-sm text-slate-700">
                <input type="checkbox" checked={postForm.is_public} onChange={e => setPostForm({ ...postForm, is_public: e.target.checked })} className="h-4 w-4 rounded border-slate-300 text-brand-600 focus:ring-brand-500" />
                Publicación pública
              </label>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowPostModal(false)}>Cancelar</Button>
            <Button onClick={handleSavePost} disabled={savingPost || !postForm.title}>
              {savingPost ? 'Guardando...' : editingPost ? 'Guardar Cambios' : 'Crear'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* ---------- Modal Candidato ---------- */}
      <Dialog open={showCandModal} onOpenChange={setShowCandModal}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>{editingCand ? 'Editar Candidato' : 'Nuevo Candidato'}</DialogTitle>
            <DialogDescription>Completá los datos del candidato</DialogDescription>
          </DialogHeader>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="cand-first">Nombre *</Label>
              <Input id="cand-first" value={candForm.first_name} onChange={e => setCandForm({ ...candForm, first_name: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cand-last">Apellido *</Label>
              <Input id="cand-last" value={candForm.last_name} onChange={e => setCandForm({ ...candForm, last_name: e.target.value })} required />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="cand-email">Email *</Label>
              <Input id="cand-email" type="email" value={candForm.email} onChange={e => setCandForm({ ...candForm, email: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cand-phone">Teléfono</Label>
              <Input id="cand-phone" value={candForm.phone} onChange={e => setCandForm({ ...candForm, phone: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cand-cc">Código de país</Label>
              <Input id="cand-cc" value={candForm.phone_country_code} onChange={e => setCandForm({ ...candForm, phone_country_code: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cand-doctype">Tipo de documento</Label>
              <Input id="cand-doctype" value={candForm.document_type} onChange={e => setCandForm({ ...candForm, document_type: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cand-docnum">N° de documento</Label>
              <Input id="cand-docnum" value={candForm.document_number} onChange={e => setCandForm({ ...candForm, document_number: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cand-birth">Fecha de nacimiento</Label>
              <Input id="cand-birth" type="date" value={candForm.birth_date} onChange={e => setCandForm({ ...candForm, birth_date: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cand-loc">Ubicación</Label>
              <Input id="cand-loc" value={candForm.location} onChange={e => setCandForm({ ...candForm, location: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cand-nat">Nacionalidad</Label>
              <Input id="cand-nat" value={candForm.nationality} onChange={e => setCandForm({ ...candForm, nationality: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cand-company">Empresa actual</Label>
              <Input id="cand-company" value={candForm.current_company} onChange={e => setCandForm({ ...candForm, current_company: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cand-position">Posición actual</Label>
              <Input id="cand-position" value={candForm.current_position} onChange={e => setCandForm({ ...candForm, current_position: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cand-notice">Días de preaviso</Label>
              <Input id="cand-notice" type="number" min="0" value={candForm.notice_period} onChange={e => setCandForm({ ...candForm, notice_period: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cand-salmin">Expectativa salarial mínima</Label>
              <Input id="cand-salmin" type="number" step="0.01" value={candForm.salary_expectation_min} onChange={e => setCandForm({ ...candForm, salary_expectation_min: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cand-salmax">Expectativa salarial máxima</Label>
              <Input id="cand-salmax" type="number" step="0.01" value={candForm.salary_expectation_max} onChange={e => setCandForm({ ...candForm, salary_expectation_max: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cand-salcur">Moneda</Label>
              <Select id="cand-salcur" options={currencyOptions} value={candForm.salary_currency} onChange={e => setCandForm({ ...candForm, salary_currency: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cand-src">Fuente</Label>
              <Input id="cand-src" value={candForm.source} onChange={e => setCandForm({ ...candForm, source: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cand-srcdetail">Detalle de fuente</Label>
              <Input id="cand-srcdetail" value={candForm.source_detail} onChange={e => setCandForm({ ...candForm, source_detail: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="cand-notes">Notas</Label>
              <textarea id="cand-notes" value={candForm.notes} onChange={e => setCandForm({ ...candForm, notes: e.target.value })} rows={2} className={textareaCls} />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowCandModal(false)}>Cancelar</Button>
            <Button onClick={handleSaveCand} disabled={savingCand || !candForm.first_name || !candForm.last_name || !candForm.email}>
              {savingCand ? 'Guardando...' : editingCand ? 'Guardar Cambios' : 'Crear'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* ---------- Modal Postulación ---------- */}
      <Dialog open={showAppModal} onOpenChange={setShowAppModal}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>Nueva Postulación</DialogTitle>
            <DialogDescription>Asociá un candidato a una publicación</DialogDescription>
          </DialogHeader>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2 col-span-2">
              <Label htmlFor="app-cand">Candidato *</Label>
              <Select id="app-cand" options={candidateOptions} placeholder="Seleccionar..." value={appForm.candidate_id} onChange={e => setAppForm({ ...appForm, candidate_id: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="app-post">Publicación *</Label>
              <Select id="app-post" options={postingOptions} placeholder="Seleccionar..." value={appForm.posting_id} onChange={e => setAppForm({ ...appForm, posting_id: e.target.value })} />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowAppModal(false)}>Cancelar</Button>
            <Button onClick={handleSaveApp} disabled={savingApp || !appForm.candidate_id || !appForm.posting_id}>
              {savingApp ? 'Guardando...' : 'Crear'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* ---------- Modal Mover etapa ---------- */}
      <Dialog open={showMoveStage} onOpenChange={setShowMoveStage}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>Mover de etapa</DialogTitle>
            <DialogDescription>{movingApp ? `Candidato: ${lookupName(candidateMap, movingApp.candidate_id)}` : ''}</DialogDescription>
          </DialogHeader>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2 col-span-2">
              <Label htmlFor="mv-stage">Etapa destino *</Label>
              <Select id="mv-stage" options={stageOptions} placeholder="Seleccionar..." value={moveForm.to_stage_id} onChange={e => setMoveForm({ ...moveForm, to_stage_id: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="mv-reason">Motivo</Label>
              <Input id="mv-reason" value={moveForm.reason} onChange={e => setMoveForm({ ...moveForm, reason: e.target.value })} />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowMoveStage(false)}>Cancelar</Button>
            <Button onClick={handleMoveStage} disabled={actionBusy || !moveForm.to_stage_id}>
              {actionBusy ? 'Guardando...' : 'Mover'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* ---------- Modal Rechazar postulación ---------- */}
      <Dialog open={showRejectApp} onOpenChange={setShowRejectApp}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>Rechazar postulación</DialogTitle>
            <DialogDescription>{rejectingApp ? `Candidato: ${lookupName(candidateMap, rejectingApp.candidate_id)}` : ''}</DialogDescription>
          </DialogHeader>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2 col-span-2">
              <Label htmlFor="rj-reason">Motivo de rechazo</Label>
              <Select id="rj-reason" options={reasonOptions} placeholder="Seleccionar..." value={rejectForm.reason_id} onChange={e => setRejectForm({ ...rejectForm, reason_id: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="rj-text">Detalle</Label>
              <Input id="rj-text" value={rejectForm.reason_text} onChange={e => setRejectForm({ ...rejectForm, reason_text: e.target.value })} />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowRejectApp(false)}>Cancelar</Button>
            <Button onClick={handleRejectApp} disabled={actionBusy}>
              {actionBusy ? 'Guardando...' : 'Rechazar'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* ---------- Modal Entrevista ---------- */}
      <Dialog open={showIvModal} onOpenChange={setShowIvModal}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>{editingIv ? 'Editar Entrevista' : 'Nueva Entrevista'}</DialogTitle>
            <DialogDescription>Completá los datos de la entrevista</DialogDescription>
          </DialogHeader>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2 col-span-2">
              <Label htmlFor="iv-app">Postulación *</Label>
              <Select id="iv-app" options={applicationOptions} placeholder="Seleccionar..." value={ivForm.application_id} onChange={e => setIvForm({ ...ivForm, application_id: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="iv-type">Tipo de entrevista *</Label>
              <Input id="iv-type" value={ivForm.interview_type} onChange={e => setIvForm({ ...ivForm, interview_type: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="iv-title">Título</Label>
              <Input id="iv-title" value={ivForm.title} onChange={e => setIvForm({ ...ivForm, title: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="iv-when">Fecha y hora</Label>
              <Input id="iv-when" type="datetime-local" value={ivForm.scheduled_at} onChange={e => setIvForm({ ...ivForm, scheduled_at: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="iv-dur">Duración (minutos)</Label>
              <Input id="iv-dur" type="number" min="1" value={ivForm.duration_minutes} onChange={e => setIvForm({ ...ivForm, duration_minutes: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="iv-url">Link de reunión</Label>
              <Input id="iv-url" value={ivForm.meeting_url} onChange={e => setIvForm({ ...ivForm, meeting_url: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="iv-pass">Contraseña de reunión</Label>
              <Input id="iv-pass" value={ivForm.meeting_password} onChange={e => setIvForm({ ...ivForm, meeting_password: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="iv-loc">Ubicación</Label>
              <Input id="iv-loc" value={ivForm.location} onChange={e => setIvForm({ ...ivForm, location: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="iv-instr">Instrucciones</Label>
              <textarea id="iv-instr" value={ivForm.instructions} onChange={e => setIvForm({ ...ivForm, instructions: e.target.value })} rows={2} className={textareaCls} />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowIvModal(false)}>Cancelar</Button>
            <Button onClick={handleSaveIv} disabled={savingIv || !ivForm.application_id || !ivForm.interview_type}>
              {savingIv ? 'Guardando...' : editingIv ? 'Guardar Cambios' : 'Crear'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* ---------- Modal Oferta ---------- */}
      <Dialog open={showOfferModal} onOpenChange={setShowOfferModal}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>{editingOffer ? 'Editar Oferta' : 'Nueva Oferta'}</DialogTitle>
            <DialogDescription>Completá los datos de la oferta</DialogDescription>
          </DialogHeader>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2 col-span-2">
              <Label htmlFor="offer-app">Postulación *</Label>
              <Select id="offer-app" options={applicationOptions} placeholder="Seleccionar..." value={offerForm.application_id} onChange={e => setOfferForm({ ...offerForm, application_id: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="offer-title">Posición *</Label>
              <Input id="offer-title" value={offerForm.position_title} onChange={e => setOfferForm({ ...offerForm, position_title: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="offer-type">Tipo de oferta *</Label>
              <Input id="offer-type" value={offerForm.offer_type} onChange={e => setOfferForm({ ...offerForm, offer_type: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="offer-dept">Departamento</Label>
              <Select id="offer-dept" options={departmentOptions} placeholder="Seleccionar..." value={offerForm.department_id} onChange={e => setOfferForm({ ...offerForm, department_id: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="offer-start">Fecha de inicio</Label>
              <Input id="offer-start" type="date" value={offerForm.start_date} onChange={e => setOfferForm({ ...offerForm, start_date: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="offer-emp">Tipo de contratación</Label>
              <Select id="offer-emp" options={employmentTypeOptions} placeholder="Seleccionar..." value={offerForm.employment_type} onChange={e => setOfferForm({ ...offerForm, employment_type: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="offer-mode">Modalidad</Label>
              <Select id="offer-mode" options={workModeOptions} placeholder="Seleccionar..." value={offerForm.work_mode} onChange={e => setOfferForm({ ...offerForm, work_mode: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="offer-amount">Salario</Label>
              <Input id="offer-amount" type="number" step="0.01" value={offerForm.salary_amount} onChange={e => setOfferForm({ ...offerForm, salary_amount: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="offer-currency">Moneda</Label>
              <Select id="offer-currency" options={currencyOptions} value={offerForm.salary_currency} onChange={e => setOfferForm({ ...offerForm, salary_currency: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="offer-period">Periodicidad</Label>
              <Select id="offer-period" options={salaryPeriodOptions} placeholder="Seleccionar..." value={offerForm.salary_period} onChange={e => setOfferForm({ ...offerForm, salary_period: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="offer-deadline">Fecha límite de respuesta</Label>
              <Input id="offer-deadline" type="date" value={offerForm.response_deadline} onChange={e => setOfferForm({ ...offerForm, response_deadline: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="offer-variable">Compensación variable</Label>
              <Input id="offer-variable" value={offerForm.variable_compensation} onChange={e => setOfferForm({ ...offerForm, variable_compensation: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="offer-benefits">Beneficios</Label>
              <textarea id="offer-benefits" value={offerForm.benefits_summary} onChange={e => setOfferForm({ ...offerForm, benefits_summary: e.target.value })} rows={2} className={textareaCls} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="offer-equity">Condiciones de equity</Label>
              <textarea id="offer-equity" value={offerForm.equity_terms} onChange={e => setOfferForm({ ...offerForm, equity_terms: e.target.value })} rows={2} className={textareaCls} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="offer-conditions">Condiciones</Label>
              <textarea id="offer-conditions" value={offerForm.conditions} onChange={e => setOfferForm({ ...offerForm, conditions: e.target.value })} rows={2} className={textareaCls} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="offer-notes">Notas</Label>
              <textarea id="offer-notes" value={offerForm.notes} onChange={e => setOfferForm({ ...offerForm, notes: e.target.value })} rows={2} className={textareaCls} />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowOfferModal(false)}>Cancelar</Button>
            <Button onClick={handleSaveOffer} disabled={savingOffer || !offerForm.application_id || !offerForm.position_title || !offerForm.offer_type}>
              {savingOffer ? 'Guardando...' : editingOffer ? 'Guardar Cambios' : 'Crear'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
