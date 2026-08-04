import { useEffect, useRef, useState } from 'react'
import { Paperclip, Pencil, Plus, Trash2, Upload } from 'lucide-react'
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

interface Expense {
  id: string
  employee_id?: string
  category_id?: string
  payment_method_id?: string
  expense_date?: string
  description: string
  original_amount?: string | number
  original_currency?: string
  merchant_name?: string
  observation?: string
  status: string
  created_at?: string
}

interface Category {
  id: string
  code: string
  name: string
  description?: string | null
  requires_receipt?: boolean
  is_active?: boolean
  sort_order?: number
}

interface Travel {
  id: string
  title: string
  purpose?: string | null
  origin: string
  destination: string
  departure_date?: string
  return_date?: string
  estimated_budget?: string | number | null
  currency?: string
  notes?: string | null
  status: string
  employee_id?: string
}

interface ExpenseReport {
  id: string
  title: string
  description?: string | null
  total_amount?: string | number
  currency?: string
  status: string
  employee_id?: string
}

interface Advance {
  id: string
  requested_amount?: string | number
  approved_amount?: string | number | null
  currency?: string
  request_date?: string
  status: string
  employee_id?: string
}

interface Policy {
  id: string
  name: string
  description?: string | null
  version?: number
  is_active?: boolean
  created_at?: string
}

interface PolicyRule {
  id: string
  policy_id: string
  category_id?: string | null
  employee_category?: string | null
  max_amount?: string | number | null
  currency?: string | null
  requires_receipt?: boolean
  requires_approval?: boolean
  priority?: number
  is_active?: boolean
}

interface Receipt {
  id: string
  filename: string
  mime_type?: string
  size?: number
  uploaded_at?: string
}

interface SelectOption {
  value: string
  label: string
}

interface WorkflowMeta {
  label: string
  cls: string
  needsInput?: boolean
  inputType?: 'text' | 'number'
  inputLabel?: string
  placeholder?: string
}

const statusStyles: Record<string, { label: string; cls: string }> = {
  draft: { label: 'Borrador', cls: 'bg-slate-100 text-slate-600' },
  requested: { label: 'Solicitado', cls: 'bg-blue-50 text-blue-700' },
  submitted: { label: 'Enviado', cls: 'bg-blue-50 text-blue-700' },
  pending: { label: 'Pendiente', cls: 'bg-amber-50 text-amber-700' },
  pending_approval: { label: 'Pendiente', cls: 'bg-amber-50 text-amber-700' },
  approved: { label: 'Aprobado', cls: 'bg-emerald-50 text-emerald-700' },
  rejected: { label: 'Rechazado', cls: 'bg-red-50 text-red-700' },
  observed: { label: 'Observado', cls: 'bg-violet-50 text-violet-700' },
  paid: { label: 'Pagado', cls: 'bg-teal-50 text-teal-700' },
  settled: { label: 'Liquidado', cls: 'bg-teal-50 text-teal-700' },
  completed: { label: 'Completado', cls: 'bg-emerald-50 text-emerald-700' },
  cancelled: { label: 'Cancelado', cls: 'bg-slate-100 text-slate-500' },
  canceled: { label: 'Cancelado', cls: 'bg-slate-100 text-slate-500' },
  active: { label: 'Activo', cls: 'bg-emerald-50 text-emerald-700' },
  inactive: { label: 'Inactivo', cls: 'bg-slate-100 text-slate-600' },
}

const statusBadge = (status: string) => {
  const s = statusStyles[(status || '').toLowerCase()] ?? { label: status, cls: 'bg-slate-100 text-slate-600' }
  return <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${s.cls}`}>{s.label}</span>
}

const fmtDate = (s?: string) => (s ? s.slice(0, 10) : '-')

const fmtMoney = (v?: string | number) => {
  if (v == null || v === '') return '-'
  const n = parseFloat(String(v))
  return isNaN(n) ? '-' : n.toLocaleString()
}

const currencyOptions: SelectOption[] = [
  { value: 'ARS', label: 'ARS' },
  { value: 'USD', label: 'USD' },
  { value: 'EUR', label: 'EUR' },
]

const entityNoun: Record<string, string> = {
  expense: 'gasto',
  travel: 'viaje',
  report: 'reporte',
  advance: 'anticipo',
}

const workflowMeta: Record<string, WorkflowMeta> = {
  submit: { label: 'Enviar', cls: 'text-blue-600' },
  request: { label: 'Solicitar', cls: 'text-blue-600' },
  approve: { label: 'Aprobar', cls: 'text-emerald-600' },
  approve_advance: { label: 'Aprobar', cls: 'text-emerald-600', needsInput: true, inputType: 'number', inputLabel: 'Monto aprobado', placeholder: 'Opcional' },
  reject: { label: 'Rechazar', cls: 'text-red-600', needsInput: true, inputType: 'text', inputLabel: 'Motivo de rechazo', placeholder: 'Motivo obligatorio' },
  observe: { label: 'Observar', cls: 'text-violet-600', needsInput: true, inputType: 'text', inputLabel: 'Observación', placeholder: 'Observación obligatoria' },
  cancel: { label: 'Cancelar', cls: 'text-slate-600' },
  complete: { label: 'Completar', cls: 'text-teal-600' },
  pay: { label: 'Pagar', cls: 'text-teal-600' },
}

const allowedActions: Record<string, (status: string) => string[]> = {
  expense: (st) => {
    switch (st) {
      case 'draft': return ['submit', 'cancel']
      case 'submitted': return ['approve', 'reject', 'observe', 'cancel']
      case 'observed': return ['reject', 'cancel']
      case 'rejected': return ['cancel']
      default: return []
    }
  },
  travel: (st) => {
    switch (st) {
      case 'draft': return ['request', 'cancel']
      case 'requested': return ['approve', 'reject', 'cancel']
      case 'approved': return ['complete', 'cancel']
      case 'rejected': return ['cancel']
      default: return []
    }
  },
  report: (st) => {
    switch (st) {
      case 'draft': return ['submit']
      case 'submitted': return ['approve', 'reject', 'observe']
      case 'observed': return ['reject']
      default: return []
    }
  },
  advance: (st) => {
    switch (st) {
      case 'requested': return ['approve_advance', 'reject', 'cancel']
      case 'approved': return ['pay', 'cancel']
      case 'rejected': return ['cancel']
      default: return []
    }
  },
}

const emptyExpenseForm = {
  category_id: '',
  expense_date: '',
  description: '',
  original_amount: '',
  original_currency: 'ARS',
  payment_method_id: '',
  merchant_name: '',
  observation: '',
}

const emptyCategoryForm = {
  code: '',
  name: '',
  description: '',
  requires_receipt: false,
  is_active: true,
  sort_order: '0',
}

const emptyTravelForm = {
  title: '',
  purpose: '',
  origin: '',
  destination: '',
  departure_date: '',
  return_date: '',
  estimated_budget: '',
  currency: 'ARS',
  notes: '',
}

const emptyReportForm = {
  title: '',
  description: '',
  currency: 'ARS',
}

const emptyAdvanceForm = {
  requested_amount: '',
  currency: 'ARS',
  request_date: '',
}

const emptyPolicyForm = {
  name: '',
  description: '',
  is_active: true,
}

const emptyRuleForm = {
  category_id: '',
  employee_category: '',
  max_amount: '',
  currency: '',
  requires_receipt: false,
  requires_approval: false,
  priority: '1',
}

const expenseStatusOptions: SelectOption[] = [
  { value: '', label: 'Todas' },
  { value: 'DRAFT', label: 'Borrador' },
  { value: 'SUBMITTED', label: 'Enviado' },
  { value: 'APPROVED', label: 'Aprobado' },
  { value: 'REJECTED', label: 'Rechazado' },
  { value: 'OBSERVED', label: 'Observado' },
  { value: 'CANCELLED', label: 'Cancelado' },
]

export default function ExpensesPage() {
  const [error, setError] = useState('')

  // Gastos
  const [items, setItems] = useState<Expense[]>([])
  const [loading, setLoading] = useState(true)
  const [statusFilter, setStatusFilter] = useState('')
  const [categories, setCategories] = useState<SelectOption[]>([])
  const [categoryMap, setCategoryMap] = useState<Record<string, string>>({})
  const [paymentMethods, setPaymentMethods] = useState<SelectOption[]>([])
  const [employeesMap, setEmployeesMap] = useState<Record<string, string>>({})

  const [showExpenseModal, setShowExpenseModal] = useState(false)
  const [editingExpense, setEditingExpense] = useState<Expense | null>(null)
  const [expenseForm, setExpenseForm] = useState({ ...emptyExpenseForm })
  const [savingExpense, setSavingExpense] = useState(false)

  // Recibos
  const [showReceipts, setShowReceipts] = useState(false)
  const [receiptExpense, setReceiptExpense] = useState<Expense | null>(null)
  const [receipts, setReceipts] = useState<Receipt[]>([])
  const [receiptFile, setReceiptFile] = useState<File | null>(null)
  const [uploadingReceipt, setUploadingReceipt] = useState(false)
  const receiptInputRef = useRef<HTMLInputElement>(null)

  // Categorías
  const [cats, setCats] = useState<Category[]>([])
  const [catLoading, setCatLoading] = useState(true)
  const [showCatModal, setShowCatModal] = useState(false)
  const [editingCat, setEditingCat] = useState<Category | null>(null)
  const [catForm, setCatForm] = useState({ ...emptyCategoryForm })
  const [savingCat, setSavingCat] = useState(false)

  // Viajes
  const [travels, setTravels] = useState<Travel[]>([])
  const [travelLoading, setTravelLoading] = useState(true)
  const [showTravelModal, setShowTravelModal] = useState(false)
  const [editingTravel, setEditingTravel] = useState<Travel | null>(null)
  const [travelForm, setTravelForm] = useState({ ...emptyTravelForm })
  const [savingTravel, setSavingTravel] = useState(false)

  // Reportes
  const [reports, setReports] = useState<ExpenseReport[]>([])
  const [reportLoading, setReportLoading] = useState(true)
  const [showReportModal, setShowReportModal] = useState(false)
  const [reportForm, setReportForm] = useState({ ...emptyReportForm })
  const [savingReport, setSavingReport] = useState(false)

  // Anticipos
  const [advances, setAdvances] = useState<Advance[]>([])
  const [advanceLoading, setAdvanceLoading] = useState(true)
  const [showAdvanceModal, setShowAdvanceModal] = useState(false)
  const [advanceForm, setAdvanceForm] = useState({ ...emptyAdvanceForm })
  const [savingAdvance, setSavingAdvance] = useState(false)

  // Políticas
  const [policies, setPolicies] = useState<Policy[]>([])
  const [policyLoading, setPolicyLoading] = useState(true)
  const [showPolicyModal, setShowPolicyModal] = useState(false)
  const [editingPolicy, setEditingPolicy] = useState<Policy | null>(null)
  const [policyForm, setPolicyForm] = useState({ ...emptyPolicyForm })
  const [savingPolicy, setSavingPolicy] = useState(false)

  const [showRules, setShowRules] = useState(false)
  const [rulePolicy, setRulePolicy] = useState<Policy | null>(null)
  const [rules, setRules] = useState<PolicyRule[]>([])
  const [ruleForm, setRuleForm] = useState({ ...emptyRuleForm })
  const [savingRule, setSavingRule] = useState(false)

  // Acción de flujo
  const [pendingAction, setPendingAction] = useState<{ entityType: string; kind: string; id: string; meta: WorkflowMeta } | null>(null)
  const [actionInput, setActionInput] = useState('')
  const [actionBusy, setActionBusy] = useState(false)

  const fetchExpenses = async () => {
    setLoading(true)
    try {
      const params: Record<string, string> = { limit: '100' }
      if (statusFilter) params.status = statusFilter
      const res = await api.get('/expenses', { params })
      setItems(res.data.data ?? [])
      setError('')
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar gastos')
      setItems([])
    } finally {
      setLoading(false)
    }
  }

  const fetchSelects = async () => {
    try {
      const [cRes, pRes, eRes] = await Promise.all([
        api.get('/expenses/categories'),
        api.get('/expenses/payment-methods'),
        api.get('/employees', { params: { limit: '500' } }),
      ])
      const cList = cRes.data.data ?? []
      setCategories(cList.map((c: any) => ({ value: c.id, label: `${c.code ? c.code + ' - ' : ''}${c.name}` })))
      const cMap: Record<string, string> = {}
      cList.forEach((c: any) => { cMap[c.id] = c.name })
      setCategoryMap(cMap)
      setPaymentMethods((pRes.data.data ?? []).map((p: any) => ({ value: p.id, label: p.name })))
      const eMap: Record<string, string> = {}
      ;(eRes.data.data ?? []).forEach((em: any) => { eMap[em.id] = `${em.first_name} ${em.last_name}`.trim() })
      setEmployeesMap(eMap)
    } catch {}
  }

  const fetchCats = async () => {
    setCatLoading(true)
    try {
      const res = await api.get('/expenses/categories')
      setCats(res.data.data ?? [])
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar categorías')
      setCats([])
    } finally {
      setCatLoading(false)
    }
  }

  const fetchTravels = async () => {
    setTravelLoading(true)
    try {
      const res = await api.get('/expenses/travels', { params: { limit: '100' } })
      setTravels(res.data.data ?? [])
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar viajes')
      setTravels([])
    } finally {
      setTravelLoading(false)
    }
  }

  const fetchReports = async () => {
    setReportLoading(true)
    try {
      const res = await api.get('/expenses/reports', { params: { limit: '100' } })
      setReports(res.data.data ?? [])
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar reportes')
      setReports([])
    } finally {
      setReportLoading(false)
    }
  }

  const fetchAdvances = async () => {
    setAdvanceLoading(true)
    try {
      const res = await api.get('/expenses/advances', { params: { limit: '100' } })
      setAdvances(res.data.data ?? [])
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar anticipos')
      setAdvances([])
    } finally {
      setAdvanceLoading(false)
    }
  }

  const fetchPolicies = async () => {
    setPolicyLoading(true)
    try {
      const res = await api.get('/expenses/policies')
      setPolicies(res.data.data ?? [])
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar políticas')
      setPolicies([])
    } finally {
      setPolicyLoading(false)
    }
  }

  useEffect(() => { fetchExpenses() }, [statusFilter])
  useEffect(() => {
    fetchSelects()
    fetchCats()
    fetchTravels()
    fetchReports()
    fetchAdvances()
    fetchPolicies()
  }, [])

  const refreshByType = (entityType: string) => {
    if (entityType === 'expense') fetchExpenses()
    else if (entityType === 'travel') fetchTravels()
    else if (entityType === 'report') fetchReports()
    else if (entityType === 'advance') fetchAdvances()
  }

  // --- Flujo de trabajo ---
  const openWorkflow = (entityType: string, kind: string, id: string, meta: WorkflowMeta, initial = '') => {
    setPendingAction({ entityType, kind, id, meta })
    setActionInput(initial)
  }

  const runWorkflow = async () => {
    if (!pendingAction) return
    setActionBusy(true)
    try {
      const { entityType, kind, id } = pendingAction
      const base: Record<string, string> = {
        expense: `/expenses/${id}`,
        travel: `/expenses/travels/${id}`,
        report: `/expenses/reports/${id}`,
        advance: `/expenses/advances/${id}`,
      }
      const action = kind === 'approve_advance' ? 'approve' : kind
      const body: Record<string, any> = {}
      if (kind === 'reject') body.reason = actionInput
      else if (kind === 'observe') body.observation = actionInput
      else if (kind === 'approve_advance' && actionInput) body.approved_amount = parseFloat(actionInput)
      await api.post(`${base[entityType]}/${action}`, body)
      setPendingAction(null)
      setActionInput('')
      refreshByType(entityType)
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al ejecutar la acción')
    } finally {
      setActionBusy(false)
    }
  }

  const renderWorkflow = (entityType: string, entity: { id: string; status: string }, initFor?: (kind: string) => string) =>
    allowedActions[entityType]((entity.status || '').toLowerCase()).map((kind) => {
      const meta = workflowMeta[kind]
      return (
        <Button
          key={kind}
          variant="ghost"
          size="sm"
          className={meta.cls}
          title={meta.label}
          onClick={() => openWorkflow(entityType, kind, entity.id, meta, initFor ? initFor(kind) : '')}
        >
          {meta.label}
        </Button>
      )
    })

  // --- Gastos ---
  const openCreateExpense = () => {
    setEditingExpense(null)
    setExpenseForm({ ...emptyExpenseForm })
    setShowExpenseModal(true)
  }

  const openEditExpense = (e: Expense) => {
    setEditingExpense(e)
    setExpenseForm({
      category_id: e.category_id ?? '',
      expense_date: e.expense_date ? e.expense_date.slice(0, 10) : '',
      description: e.description,
      original_amount: e.original_amount != null ? String(e.original_amount) : '',
      original_currency: e.original_currency || 'ARS',
      payment_method_id: e.payment_method_id ?? '',
      merchant_name: e.merchant_name ?? '',
      observation: e.observation ?? '',
    })
    setShowExpenseModal(true)
  }

  const handleSaveExpense = async () => {
    setSavingExpense(true)
    try {
      const body: Record<string, any> = {
        category_id: expenseForm.category_id,
        expense_date: expenseForm.expense_date ? expenseForm.expense_date + 'T00:00:00Z' : null,
        description: expenseForm.description,
        original_amount: expenseForm.original_amount,
        original_currency: expenseForm.original_currency,
        total_amount: expenseForm.original_amount ? parseFloat(expenseForm.original_amount) : 0,
      }
      if (expenseForm.payment_method_id) body.payment_method_id = expenseForm.payment_method_id
      if (expenseForm.merchant_name) body.merchant_name = expenseForm.merchant_name
      if (expenseForm.observation) body.observation = expenseForm.observation
      if (editingExpense) {
        body.status = editingExpense.status
        await api.put(`/expenses/${editingExpense.id}`, body)
      } else {
        await api.post('/expenses', body)
      }
      setShowExpenseModal(false)
      fetchExpenses()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar gasto')
    } finally {
      setSavingExpense(false)
    }
  }

  const handleDeleteExpense = async (e: Expense) => {
    if (!confirm(`¿Eliminar el gasto "${e.description}"?`)) return
    try {
      await api.delete(`/expenses/${e.id}`)
      fetchExpenses()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al eliminar gasto')
    }
  }

  // --- Recibos ---
  const openReceipts = async (e: Expense) => {
    setReceiptExpense(e)
    setReceipts([])
    setReceiptFile(null)
    if (receiptInputRef.current) receiptInputRef.current.value = ''
    setShowReceipts(true)
    try {
      const res = await api.get(`/expenses/${e.id}/receipts`)
      setReceipts(res.data.data ?? [])
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al cargar recibos')
    }
  }

  const handleUploadReceipt = async () => {
    if (!receiptExpense || !receiptFile) return
    setUploadingReceipt(true)
    try {
      const fd = new FormData()
      fd.append('file', receiptFile)
      await api.post(`/expenses/${receiptExpense.id}/receipts`, fd, { headers: { 'Content-Type': null } })
      setReceiptFile(null)
      if (receiptInputRef.current) receiptInputRef.current.value = ''
      const res = await api.get(`/expenses/${receiptExpense.id}/receipts`)
      setReceipts(res.data.data ?? [])
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al subir recibo')
    } finally {
      setUploadingReceipt(false)
    }
  }

  const handleDeleteReceipt = async (r: Receipt) => {
    if (!receiptExpense) return
    if (!confirm(`¿Eliminar el recibo "${r.filename}"?`)) return
    try {
      await api.delete(`/expenses/${receiptExpense.id}/receipts/${r.id}`)
      const res = await api.get(`/expenses/${receiptExpense.id}/receipts`)
      setReceipts(res.data.data ?? [])
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al eliminar recibo')
    }
  }

  // --- Categorías ---
  const openCreateCat = () => {
    setEditingCat(null)
    setCatForm({ ...emptyCategoryForm })
    setShowCatModal(true)
  }

  const openEditCat = (c: Category) => {
    setEditingCat(c)
    setCatForm({
      code: c.code,
      name: c.name,
      description: c.description ?? '',
      requires_receipt: c.requires_receipt ?? false,
      is_active: c.is_active ?? true,
      sort_order: c.sort_order != null ? String(c.sort_order) : '0',
    })
    setShowCatModal(true)
  }

  const handleSaveCat = async () => {
    setSavingCat(true)
    try {
      const body: Record<string, any> = {
        code: catForm.code,
        name: catForm.name,
        requires_receipt: catForm.requires_receipt,
        is_active: catForm.is_active,
        sort_order: parseInt(catForm.sort_order || '0', 10),
      }
      if (catForm.description) body.description = catForm.description
      if (editingCat) {
        await api.put(`/expenses/categories/${editingCat.id}`, body)
      } else {
        await api.post('/expenses/categories', body)
      }
      setShowCatModal(false)
      fetchCats()
      fetchSelects()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar categoría')
    } finally {
      setSavingCat(false)
    }
  }

  const handleDeleteCat = async (c: Category) => {
    if (!confirm(`¿Eliminar la categoría "${c.name}"?`)) return
    try {
      await api.delete(`/expenses/categories/${c.id}`)
      fetchCats()
      fetchSelects()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al eliminar categoría')
    }
  }

  // --- Viajes ---
  const openCreateTravel = () => {
    setEditingTravel(null)
    setTravelForm({ ...emptyTravelForm })
    setShowTravelModal(true)
  }

  const openEditTravel = (t: Travel) => {
    setEditingTravel(t)
    setTravelForm({
      title: t.title,
      purpose: t.purpose ?? '',
      origin: t.origin,
      destination: t.destination,
      departure_date: t.departure_date ? t.departure_date.slice(0, 10) : '',
      return_date: t.return_date ? t.return_date.slice(0, 10) : '',
      estimated_budget: t.estimated_budget != null ? String(t.estimated_budget) : '',
      currency: t.currency || 'ARS',
      notes: t.notes ?? '',
    })
    setShowTravelModal(true)
  }

  const handleSaveTravel = async () => {
    setSavingTravel(true)
    try {
      const body: Record<string, any> = {
        title: travelForm.title,
        origin: travelForm.origin,
        destination: travelForm.destination,
        departure_date: travelForm.departure_date ? travelForm.departure_date + 'T00:00:00Z' : null,
        return_date: travelForm.return_date ? travelForm.return_date + 'T00:00:00Z' : null,
        currency: travelForm.currency,
      }
      if (travelForm.purpose) body.purpose = travelForm.purpose
      if (travelForm.notes) body.notes = travelForm.notes
      if (travelForm.estimated_budget) body.estimated_budget = travelForm.estimated_budget
      if (editingTravel) {
        await api.put(`/expenses/travels/${editingTravel.id}`, body)
      } else {
        await api.post('/expenses/travels', body)
      }
      setShowTravelModal(false)
      fetchTravels()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar viaje')
    } finally {
      setSavingTravel(false)
    }
  }

  // --- Reportes ---
  const openCreateReport = () => {
    setReportForm({ ...emptyReportForm })
    setShowReportModal(true)
  }

  const handleSaveReport = async () => {
    setSavingReport(true)
    try {
      const body: Record<string, any> = {
        title: reportForm.title,
        currency: reportForm.currency,
      }
      if (reportForm.description) body.description = reportForm.description
      await api.post('/expenses/reports', body)
      setShowReportModal(false)
      fetchReports()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar reporte')
    } finally {
      setSavingReport(false)
    }
  }

  // --- Anticipos ---
  const openCreateAdvance = () => {
    setAdvanceForm({ ...emptyAdvanceForm })
    setShowAdvanceModal(true)
  }

  const handleSaveAdvance = async () => {
    setSavingAdvance(true)
    try {
      const body: Record<string, any> = {
        requested_amount: advanceForm.requested_amount,
        currency: advanceForm.currency,
        request_date: advanceForm.request_date ? advanceForm.request_date + 'T00:00:00Z' : null,
      }
      await api.post('/expenses/advances', body)
      setShowAdvanceModal(false)
      fetchAdvances()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar anticipo')
    } finally {
      setSavingAdvance(false)
    }
  }

  // --- Políticas ---
  const openCreatePolicy = () => {
    setEditingPolicy(null)
    setPolicyForm({ ...emptyPolicyForm })
    setShowPolicyModal(true)
  }

  const openEditPolicy = (p: Policy) => {
    setEditingPolicy(p)
    setPolicyForm({
      name: p.name,
      description: p.description ?? '',
      is_active: p.is_active ?? true,
    })
    setShowPolicyModal(true)
  }

  const handleSavePolicy = async () => {
    setSavingPolicy(true)
    try {
      const body: Record<string, any> = { name: policyForm.name, is_active: policyForm.is_active }
      if (policyForm.description) body.description = policyForm.description
      if (editingPolicy) {
        await api.put(`/expenses/policies/${editingPolicy.id}`, body)
      } else {
        await api.post('/expenses/policies', body)
      }
      setShowPolicyModal(false)
      fetchPolicies()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar política')
    } finally {
      setSavingPolicy(false)
    }
  }

  // --- Reglas ---
  const openRules = async (p: Policy) => {
    setRulePolicy(p)
    setRules([])
    setRuleForm({ ...emptyRuleForm })
    setShowRules(true)
    try {
      const res = await api.get(`/expenses/policies/${p.id}/rules`)
      setRules(res.data.data ?? [])
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al cargar reglas')
    }
  }

  const handleSaveRule = async () => {
    if (!rulePolicy) return
    setSavingRule(true)
    try {
      const body: Record<string, any> = {
        policy_id: rulePolicy.id,
        requires_receipt: ruleForm.requires_receipt,
        requires_approval: ruleForm.requires_approval,
        is_active: true,
        priority: parseInt(ruleForm.priority || '1', 10),
      }
      if (ruleForm.category_id) body.category_id = ruleForm.category_id
      if (ruleForm.employee_category) body.employee_category = ruleForm.employee_category
      if (ruleForm.max_amount) body.max_amount = ruleForm.max_amount
      if (ruleForm.currency) body.currency = ruleForm.currency
      await api.post(`/expenses/policies/${rulePolicy.id}/rules`, body)
      setRuleForm({ ...emptyRuleForm })
      const res = await api.get(`/expenses/policies/${rulePolicy.id}/rules`)
      setRules(res.data.data ?? [])
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al crear regla')
    } finally {
      setSavingRule(false)
    }
  }

  const handleDeleteRule = async (r: PolicyRule) => {
    if (!rulePolicy) return
    if (!confirm('¿Eliminar esta regla?')) return
    try {
      await api.delete(`/expenses/policies/${rulePolicy.id}/rules/${r.id}`)
      const res = await api.get(`/expenses/policies/${rulePolicy.id}/rules`)
      setRules(res.data.data ?? [])
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al eliminar regla')
    }
  }

  const categoryName = (id?: string) => {
    if (!id) return '-'
    return categoryMap[id] ?? id.slice(0, 8)
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-slate-900">Gastos</h1>
      </div>

      {error && <div className="mb-4 p-3 bg-red-50 text-red-700 text-sm rounded-lg">{error}</div>}

      <Tabs defaultValue="expenses">
        <TabsList>
          <TabsTrigger value="expenses">Gastos</TabsTrigger>
          <TabsTrigger value="categories">Categorías</TabsTrigger>
          <TabsTrigger value="travels">Viajes</TabsTrigger>
          <TabsTrigger value="reports">Reportes</TabsTrigger>
          <TabsTrigger value="advances">Anticipos</TabsTrigger>
          <TabsTrigger value="policies">Políticas</TabsTrigger>
        </TabsList>

        {/* ---------------- Gastos ---------------- */}
        <TabsContent value="expenses">
          <div className="flex items-center gap-3 mb-4">
            <div className="w-56">
              <Select
                aria-label="Filtrar por estado"
                options={expenseStatusOptions}
                value={statusFilter}
                onChange={e => setStatusFilter(e.target.value)}
              />
            </div>
            <div className="ml-auto">
              <Button size="sm" onClick={openCreateExpense}><Plus size={16} className="mr-1" /> Nuevo gasto</Button>
            </div>
          </div>

          <Card>
            <CardContent className="p-0">
              {loading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
              : items.length === 0 && !error ? <div className="p-6 text-center text-slate-500">No hay gastos registrados</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Fecha</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Descripción</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Categoría</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Monto</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Moneda</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Solicitante</th>
                        <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                      </tr>
                    </thead>
                    <tbody>
                      {items.map(e => {
                        const st = (e.status || '').toLowerCase()
                        return (
                          <tr key={e.id} className="border-b border-slate-100 hover:bg-slate-50">
                            <td className="px-4 py-3 text-slate-600">{fmtDate(e.expense_date)}</td>
                            <td className="px-4 py-3 font-medium text-slate-900">{e.description}</td>
                            <td className="px-4 py-3 text-slate-600">{categoryName(e.category_id)}</td>
                            <td className="px-4 py-3 text-slate-900 font-medium">{fmtMoney(e.original_amount)}</td>
                            <td className="px-4 py-3 text-slate-600">{e.original_currency || '-'}</td>
                            <td className="px-4 py-3">{statusBadge(e.status)}</td>
                            <td className="px-4 py-3 text-slate-600">{e.employee_id ? (employeesMap[e.employee_id] || '-') : '-'}</td>
                            <td className="px-4 py-3 text-right whitespace-nowrap">
                              {renderWorkflow('expense', e)}
                              <Button variant="ghost" size="sm" className="text-slate-600" title="Recibos" onClick={() => openReceipts(e)}>
                                <Paperclip size={14} />
                              </Button>
                              {st === 'draft' && (
                                <Button variant="ghost" size="sm" title="Editar" onClick={() => openEditExpense(e)}>
                                  <Pencil size={14} />
                                </Button>
                              )}
                              {st === 'draft' && (
                                <Button variant="ghost" size="sm" className="text-red-500" title="Eliminar" onClick={() => handleDeleteExpense(e)}>
                                  <Trash2 size={14} />
                                </Button>
                              )}
                            </td>
                          </tr>
                        )
                      })}
                    </tbody>
                  </table>
                </div>}
            </CardContent>
          </Card>
        </TabsContent>

        {/* ---------------- Categorías ---------------- */}
        <TabsContent value="categories">
          <div className="flex items-center justify-between mb-4">
            <span className="text-sm text-slate-500">{cats.length} categorías</span>
            <Button size="sm" onClick={openCreateCat}><Plus size={16} className="mr-1" /> Nueva</Button>
          </div>

          <Card>
            <CardContent className="p-0">
              {catLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
              : cats.length === 0 ? <div className="p-6 text-center text-slate-500">No hay categorías</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Código</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Descripción</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Requiere recibo</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                        <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                      </tr>
                    </thead>
                    <tbody>
                      {cats.map(c => (
                        <tr key={c.id} className="border-b border-slate-100 hover:bg-slate-50">
                          <td className="px-4 py-3 text-slate-500">{c.code}</td>
                          <td className="px-4 py-3 font-medium text-slate-900">{c.name}</td>
                          <td className="px-4 py-3 text-slate-600">{c.description || '-'}</td>
                          <td className="px-4 py-3 text-slate-600">{c.requires_receipt ? 'Sí' : 'No'}</td>
                          <td className="px-4 py-3">{statusBadge(c.is_active ? 'active' : 'inactive')}</td>
                          <td className="px-4 py-3 text-right whitespace-nowrap">
                            <Button variant="ghost" size="sm" title="Editar" onClick={() => openEditCat(c)}>
                              <Pencil size={14} />
                            </Button>
                            <Button variant="ghost" size="sm" className="text-red-500" title="Eliminar" onClick={() => handleDeleteCat(c)}>
                              <Trash2 size={14} />
                            </Button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>}
            </CardContent>
          </Card>
        </TabsContent>

        {/* ---------------- Viajes ---------------- */}
        <TabsContent value="travels">
          <div className="flex items-center justify-between mb-4">
            <span className="text-sm text-slate-500">{travels.length} viajes</span>
            <Button size="sm" onClick={openCreateTravel}><Plus size={16} className="mr-1" /> Nuevo viaje</Button>
          </div>

          <Card>
            <CardContent className="p-0">
              {travelLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
              : travels.length === 0 ? <div className="p-6 text-center text-slate-500">No hay viajes registrados</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Título</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Ruta</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Salida</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Regreso</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Presupuesto</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                        <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                      </tr>
                    </thead>
                    <tbody>
                      {travels.map(t => {
                        const st = (t.status || '').toLowerCase()
                        return (
                          <tr key={t.id} className="border-b border-slate-100 hover:bg-slate-50">
                            <td className="px-4 py-3 font-medium text-slate-900">{t.title}</td>
                            <td className="px-4 py-3 text-slate-600">{t.origin} → {t.destination}</td>
                            <td className="px-4 py-3 text-slate-600">{fmtDate(t.departure_date)}</td>
                            <td className="px-4 py-3 text-slate-600">{fmtDate(t.return_date)}</td>
                            <td className="px-4 py-3 text-slate-900 font-medium">{fmtMoney(t.estimated_budget ?? 0)}</td>
                            <td className="px-4 py-3">{statusBadge(t.status)}</td>
                            <td className="px-4 py-3 text-right whitespace-nowrap">
                              {renderWorkflow('travel', t)}
                              {st === 'draft' && (
                                <Button variant="ghost" size="sm" title="Editar" onClick={() => openEditTravel(t)}>
                                  <Pencil size={14} />
                                </Button>
                              )}
                            </td>
                          </tr>
                        )
                      })}
                    </tbody>
                  </table>
                </div>}
            </CardContent>
          </Card>
        </TabsContent>

        {/* ---------------- Reportes ---------------- */}
        <TabsContent value="reports">
          <div className="flex items-center justify-between mb-4">
            <span className="text-sm text-slate-500">{reports.length} reportes</span>
            <Button size="sm" onClick={openCreateReport}><Plus size={16} className="mr-1" /> Nuevo reporte</Button>
          </div>

          <Card>
            <CardContent className="p-0">
              {reportLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
              : reports.length === 0 ? <div className="p-6 text-center text-slate-500">No hay reportes registrados</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Título</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Descripción</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Total</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Moneda</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                        <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                      </tr>
                    </thead>
                    <tbody>
                      {reports.map(r => (
                        <tr key={r.id} className="border-b border-slate-100 hover:bg-slate-50">
                          <td className="px-4 py-3 font-medium text-slate-900">{r.title}</td>
                          <td className="px-4 py-3 text-slate-600">{r.description || '-'}</td>
                          <td className="px-4 py-3 text-slate-900 font-medium">{fmtMoney(r.total_amount ?? 0)}</td>
                          <td className="px-4 py-3 text-slate-600">{r.currency || '-'}</td>
                          <td className="px-4 py-3">{statusBadge(r.status)}</td>
                          <td className="px-4 py-3 text-right whitespace-nowrap">
                            {renderWorkflow('report', r)}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>}
            </CardContent>
          </Card>
        </TabsContent>

        {/* ---------------- Anticipos ---------------- */}
        <TabsContent value="advances">
          <div className="flex items-center justify-between mb-4">
            <span className="text-sm text-slate-500">{advances.length} anticipos</span>
            <Button size="sm" onClick={openCreateAdvance}><Plus size={16} className="mr-1" /> Nuevo anticipo</Button>
          </div>

          <Card>
            <CardContent className="p-0">
              {advanceLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
              : advances.length === 0 ? <div className="p-6 text-center text-slate-500">No hay anticipos registrados</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Fecha</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Solicitado</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Aprobado</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Moneda</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                        <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                      </tr>
                    </thead>
                    <tbody>
                      {advances.map(a => (
                        <tr key={a.id} className="border-b border-slate-100 hover:bg-slate-50">
                          <td className="px-4 py-3 text-slate-600">{fmtDate(a.request_date)}</td>
                          <td className="px-4 py-3 text-slate-900 font-medium">{fmtMoney(a.requested_amount ?? 0)}</td>
                          <td className="px-4 py-3 text-slate-600">{a.approved_amount != null ? fmtMoney(a.approved_amount) : '-'}</td>
                          <td className="px-4 py-3 text-slate-600">{a.currency || '-'}</td>
                          <td className="px-4 py-3">{statusBadge(a.status)}</td>
                          <td className="px-4 py-3 text-right whitespace-nowrap">
                            {renderWorkflow('advance', a, (k) => k === 'approve_advance' ? String(a.requested_amount ?? '') : '')}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>}
            </CardContent>
          </Card>
        </TabsContent>

        {/* ---------------- Políticas ---------------- */}
        <TabsContent value="policies">
          <div className="flex items-center justify-between mb-4">
            <span className="text-sm text-slate-500">{policies.length} políticas</span>
            <Button size="sm" onClick={openCreatePolicy}><Plus size={16} className="mr-1" /> Nueva</Button>
          </div>

          <Card>
            <CardContent className="p-0">
              {policyLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
              : policies.length === 0 ? <div className="p-6 text-center text-slate-500">No hay políticas registradas</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Descripción</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Versión</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                        <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                      </tr>
                    </thead>
                    <tbody>
                      {policies.map(p => (
                        <tr key={p.id} className="border-b border-slate-100 hover:bg-slate-50">
                          <td className="px-4 py-3 font-medium text-slate-900">{p.name}</td>
                          <td className="px-4 py-3 text-slate-600">{p.description || '-'}</td>
                          <td className="px-4 py-3 text-slate-600">{p.version ?? '-'}</td>
                          <td className="px-4 py-3">{statusBadge(p.is_active ? 'active' : 'inactive')}</td>
                          <td className="px-4 py-3 text-right whitespace-nowrap">
                            <Button variant="ghost" size="sm" className="text-slate-600" title="Reglas" onClick={() => openRules(p)}>
                              Lista de reglas
                            </Button>
                            <Button variant="ghost" size="sm" title="Editar" onClick={() => openEditPolicy(p)}>
                              <Pencil size={14} />
                            </Button>
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

      {/* --- Dialogo nuevo/editar gasto --- */}
      <Dialog open={showExpenseModal} onOpenChange={setShowExpenseModal}>
        <DialogContent className="max-w-xl">
          <DialogHeader>
            <DialogTitle>{editingExpense ? 'Editar Gasto' : 'Nuevo Gasto'}</DialogTitle>
            <DialogDescription>
              {editingExpense ? 'Modificá los datos del gasto' : 'Completá los datos para registrar un nuevo gasto'}
            </DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2 col-span-2">
              <Label htmlFor="exp-desc">Descripción *</Label>
              <Input id="exp-desc" value={expenseForm.description} onChange={e => setExpenseForm({ ...expenseForm, description: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="exp-cat">Categoría *</Label>
              <Select id="exp-cat" options={categories} placeholder="Seleccionar..." value={expenseForm.category_id} onChange={e => setExpenseForm({ ...expenseForm, category_id: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="exp-date">Fecha *</Label>
              <Input id="exp-date" type="date" value={expenseForm.expense_date} onChange={e => setExpenseForm({ ...expenseForm, expense_date: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="exp-amount">Monto *</Label>
              <Input id="exp-amount" type="number" step="0.01" value={expenseForm.original_amount} onChange={e => setExpenseForm({ ...expenseForm, original_amount: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="exp-currency">Moneda *</Label>
              <Select id="exp-currency" options={currencyOptions} value={expenseForm.original_currency} onChange={e => setExpenseForm({ ...expenseForm, original_currency: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="exp-payment">Método de Pago</Label>
              <Select id="exp-payment" options={paymentMethods} placeholder="Seleccionar..." value={expenseForm.payment_method_id} onChange={e => setExpenseForm({ ...expenseForm, payment_method_id: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="exp-merchant">Comercio</Label>
              <Input id="exp-merchant" value={expenseForm.merchant_name} onChange={e => setExpenseForm({ ...expenseForm, merchant_name: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="exp-obs">Observación</Label>
              <Input id="exp-obs" value={expenseForm.observation} onChange={e => setExpenseForm({ ...expenseForm, observation: e.target.value })} />
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowExpenseModal(false)}>Cancelar</Button>
            <Button onClick={handleSaveExpense} disabled={savingExpense || !expenseForm.description || !expenseForm.category_id || !expenseForm.expense_date || !expenseForm.original_amount}>
              {savingExpense ? 'Guardando...' : editingExpense ? 'Guardar Cambios' : 'Crear Gasto'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* --- Dialogo recibos --- */}
      <Dialog open={showReceipts} onOpenChange={setShowReceipts}>
        <DialogContent className="max-w-xl">
          <DialogHeader>
            <DialogTitle>Recibos</DialogTitle>
            <DialogDescription>
              {receiptExpense ? `Gasto: ${receiptExpense.description}` : ''}
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-3">
            {receipts.length === 0 ? (
              <div className="p-4 text-center text-sm text-slate-500 border border-slate-100 rounded-lg">Sin recibos cargados</div>
            ) : (
              <div className="space-y-2">
                {receipts.map(r => (
                  <div key={r.id} className="flex items-center justify-between px-3 py-2 border border-slate-100 rounded-lg">
                    <div className="flex items-center gap-2 text-sm text-slate-700">
                      <Paperclip size={14} className="text-slate-400" />
                      <span className="font-medium">{r.filename}</span>
                      <span className="text-xs text-slate-400">{r.size ? `${Math.round(r.size / 1024)} KB` : ''}</span>
                      <span className="text-xs text-slate-400">{fmtDate(r.uploaded_at)}</span>
                    </div>
                    <Button variant="ghost" size="sm" className="text-red-500" title="Eliminar" onClick={() => handleDeleteReceipt(r)}>
                      <Trash2 size={14} />
                    </Button>
                  </div>
                ))}
              </div>
            )}

            <div className="flex items-center gap-3 border-t border-slate-100 pt-4">
              <Input
                ref={receiptInputRef}
                type="file"
                className="max-w-xs"
                onChange={e => setReceiptFile(e.target.files?.[0] ?? null)}
              />
              <Button size="sm" onClick={handleUploadReceipt} disabled={uploadingReceipt || !receiptFile}>
                <Upload size={14} className="mr-1" /> {uploadingReceipt ? 'Subiendo...' : 'Subir recibo'}
              </Button>
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowReceipts(false)}>Cerrar</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* --- Dialogo categorías --- */}
      <Dialog open={showCatModal} onOpenChange={setShowCatModal}>
        <DialogContent className="max-w-xl">
          <DialogHeader>
            <DialogTitle>{editingCat ? 'Editar Categoría' : 'Nueva Categoría'}</DialogTitle>
            <DialogDescription>
              {editingCat ? 'Modificá los datos de la categoría' : 'Completá los datos para registrar una nueva categoría'}
            </DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="cat-code">Código *</Label>
              <Input id="cat-code" value={catForm.code} onChange={e => setCatForm({ ...catForm, code: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cat-name">Nombre *</Label>
              <Input id="cat-name" value={catForm.name} onChange={e => setCatForm({ ...catForm, name: e.target.value })} required />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="cat-desc">Descripción</Label>
              <Input id="cat-desc" value={catForm.description} onChange={e => setCatForm({ ...catForm, description: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cat-sort">Orden</Label>
              <Input id="cat-sort" type="number" value={catForm.sort_order} onChange={e => setCatForm({ ...catForm, sort_order: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <label className="flex items-center gap-2 text-sm text-slate-700">
                <input type="checkbox" checked={catForm.requires_receipt} onChange={e => setCatForm({ ...catForm, requires_receipt: e.target.checked })} className="h-4 w-4 rounded border-slate-300 text-brand-600 focus:ring-brand-500" />
                Requiere recibo
              </label>
            </div>
            <div className="space-y-2 col-span-2">
              <label className="flex items-center gap-2 text-sm text-slate-700">
                <input type="checkbox" checked={catForm.is_active} onChange={e => setCatForm({ ...catForm, is_active: e.target.checked })} className="h-4 w-4 rounded border-slate-300 text-brand-600 focus:ring-brand-500" />
                Activa
              </label>
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowCatModal(false)}>Cancelar</Button>
            <Button onClick={handleSaveCat} disabled={savingCat || !catForm.code || !catForm.name}>
              {savingCat ? 'Guardando...' : editingCat ? 'Guardar Cambios' : 'Crear Categoría'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* --- Dialogo viajes --- */}
      <Dialog open={showTravelModal} onOpenChange={setShowTravelModal}>
        <DialogContent className="max-w-xl">
          <DialogHeader>
            <DialogTitle>{editingTravel ? 'Editar Viaje' : 'Nuevo Viaje'}</DialogTitle>
            <DialogDescription>
              {editingTravel ? 'Modificá los datos del viaje' : 'Completá los datos para registrar un nuevo viaje'}
            </DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2 col-span-2">
              <Label htmlFor="trv-title">Título *</Label>
              <Input id="trv-title" value={travelForm.title} onChange={e => setTravelForm({ ...travelForm, title: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="trv-origin">Origen *</Label>
              <Input id="trv-origin" value={travelForm.origin} onChange={e => setTravelForm({ ...travelForm, origin: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="trv-dest">Destino *</Label>
              <Input id="trv-dest" value={travelForm.destination} onChange={e => setTravelForm({ ...travelForm, destination: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="trv-dep">Salida *</Label>
              <Input id="trv-dep" type="date" value={travelForm.departure_date} onChange={e => setTravelForm({ ...travelForm, departure_date: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="trv-ret">Regreso *</Label>
              <Input id="trv-ret" type="date" value={travelForm.return_date} onChange={e => setTravelForm({ ...travelForm, return_date: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="trv-budget">Presupuesto estimado</Label>
              <Input id="trv-budget" type="number" step="0.01" value={travelForm.estimated_budget} onChange={e => setTravelForm({ ...travelForm, estimated_budget: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="trv-currency">Moneda</Label>
              <Select id="trv-currency" options={currencyOptions} value={travelForm.currency} onChange={e => setTravelForm({ ...travelForm, currency: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="trv-purpose">Propósito</Label>
              <Input id="trv-purpose" value={travelForm.purpose} onChange={e => setTravelForm({ ...travelForm, purpose: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="trv-notes">Notas</Label>
              <Input id="trv-notes" value={travelForm.notes} onChange={e => setTravelForm({ ...travelForm, notes: e.target.value })} />
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowTravelModal(false)}>Cancelar</Button>
            <Button onClick={handleSaveTravel} disabled={savingTravel || !travelForm.title || !travelForm.origin || !travelForm.destination || !travelForm.departure_date || !travelForm.return_date}>
              {savingTravel ? 'Guardando...' : editingTravel ? 'Guardar Cambios' : 'Crear Viaje'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* --- Dialogo reportes --- */}
      <Dialog open={showReportModal} onOpenChange={setShowReportModal}>
        <DialogContent className="max-w-xl">
          <DialogHeader>
            <DialogTitle>Nuevo Reporte</DialogTitle>
            <DialogDescription>Completá los datos para registrar un nuevo reporte de gastos</DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2 col-span-2">
              <Label htmlFor="rep-title">Título *</Label>
              <Input id="rep-title" value={reportForm.title} onChange={e => setReportForm({ ...reportForm, title: e.target.value })} required />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="rep-desc">Descripción</Label>
              <Input id="rep-desc" value={reportForm.description} onChange={e => setReportForm({ ...reportForm, description: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="rep-currency">Moneda *</Label>
              <Select id="rep-currency" options={currencyOptions} value={reportForm.currency} onChange={e => setReportForm({ ...reportForm, currency: e.target.value })} />
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowReportModal(false)}>Cancelar</Button>
            <Button onClick={handleSaveReport} disabled={savingReport || !reportForm.title}>
              {savingReport ? 'Guardando...' : 'Crear Reporte'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* --- Dialogo anticipos --- */}
      <Dialog open={showAdvanceModal} onOpenChange={setShowAdvanceModal}>
        <DialogContent className="max-w-xl">
          <DialogHeader>
            <DialogTitle>Nuevo Anticipo</DialogTitle>
            <DialogDescription>Completá los datos para solicitar un nuevo anticipo</DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="adv-amount">Monto solicitado *</Label>
              <Input id="adv-amount" type="number" step="0.01" value={advanceForm.requested_amount} onChange={e => setAdvanceForm({ ...advanceForm, requested_amount: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="adv-currency">Moneda *</Label>
              <Select id="adv-currency" options={currencyOptions} value={advanceForm.currency} onChange={e => setAdvanceForm({ ...advanceForm, currency: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="adv-date">Fecha de solicitud *</Label>
              <Input id="adv-date" type="date" value={advanceForm.request_date} onChange={e => setAdvanceForm({ ...advanceForm, request_date: e.target.value })} required />
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowAdvanceModal(false)}>Cancelar</Button>
            <Button onClick={handleSaveAdvance} disabled={savingAdvance || !advanceForm.requested_amount || !advanceForm.request_date}>
              {savingAdvance ? 'Guardando...' : 'Solicitar Anticipo'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* --- Dialogo políticas --- */}
      <Dialog open={showPolicyModal} onOpenChange={setShowPolicyModal}>
        <DialogContent className="max-w-xl">
          <DialogHeader>
            <DialogTitle>{editingPolicy ? 'Editar Política' : 'Nueva Política'}</DialogTitle>
            <DialogDescription>
              {editingPolicy ? 'Modificá los datos de la política' : 'Completá los datos para registrar una nueva política'}
            </DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2 col-span-2">
              <Label htmlFor="pol-name">Nombre *</Label>
              <Input id="pol-name" value={policyForm.name} onChange={e => setPolicyForm({ ...policyForm, name: e.target.value })} required />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="pol-desc">Descripción</Label>
              <Input id="pol-desc" value={policyForm.description} onChange={e => setPolicyForm({ ...policyForm, description: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <label className="flex items-center gap-2 text-sm text-slate-700">
                <input type="checkbox" checked={policyForm.is_active} onChange={e => setPolicyForm({ ...policyForm, is_active: e.target.checked })} className="h-4 w-4 rounded border-slate-300 text-brand-600 focus:ring-brand-500" />
                Activa
              </label>
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowPolicyModal(false)}>Cancelar</Button>
            <Button onClick={handleSavePolicy} disabled={savingPolicy || !policyForm.name}>
              {savingPolicy ? 'Guardando...' : editingPolicy ? 'Guardar Cambios' : 'Crear Política'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* --- Dialogo reglas de política --- */}
      <Dialog open={showRules} onOpenChange={setShowRules}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Reglas</DialogTitle>
            <DialogDescription>{rulePolicy ? `Política: ${rulePolicy.name}` : ''}</DialogDescription>
          </DialogHeader>

          <div className="space-y-4">
            <div>
              <div className="text-sm font-medium text-slate-900 mb-2">Reglas existentes ({rules.length})</div>
              {rules.length === 0 ? (
                <div className="p-4 text-center text-sm text-slate-500 border border-slate-100 rounded-lg">Sin reglas definidas</div>
              ) : (
                <div className="space-y-2">
                  {rules.map(r => (
                    <div key={r.id} className="flex items-center justify-between px-3 py-2 border border-slate-100 rounded-lg">
                      <div className="text-sm text-slate-700">
                        <span className="font-medium">{r.employee_category || (r.category_id ? categoryName(r.category_id) : 'Todas las categorías')}</span>
                        {r.max_amount != null && <span className="text-slate-500"> · máx {fmtMoney(r.max_amount)} {r.currency || ''}</span>}
                        {r.requires_receipt && <span className="text-slate-500"> · recibo</span>}
                        {r.requires_approval && <span className="text-slate-500"> · aprobación</span>}
                      </div>
                      <Button variant="ghost" size="sm" className="text-red-500" title="Eliminar" onClick={() => handleDeleteRule(r)}>
                        <Trash2 size={14} />
                      </Button>
                    </div>
                  ))}
                </div>
              )}
            </div>

            <div className="border-t border-slate-100 pt-4">
              <div className="text-sm font-medium text-slate-900 mb-3">Nueva regla</div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="rule-cat">Categoría</Label>
                  <Select id="rule-cat" options={categories} placeholder="Todas" value={ruleForm.category_id} onChange={e => setRuleForm({ ...ruleForm, category_id: e.target.value })} />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="rule-emp-cat">Categoría de empleado</Label>
                  <Input id="rule-emp-cat" value={ruleForm.employee_category} onChange={e => setRuleForm({ ...ruleForm, employee_category: e.target.value })} />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="rule-max">Monto máximo</Label>
                  <Input id="rule-max" type="number" step="0.01" value={ruleForm.max_amount} onChange={e => setRuleForm({ ...ruleForm, max_amount: e.target.value })} />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="rule-currency">Moneda</Label>
                  <Select id="rule-currency" options={currencyOptions} placeholder="Cualquiera" value={ruleForm.currency} onChange={e => setRuleForm({ ...ruleForm, currency: e.target.value })} />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="rule-priority">Prioridad</Label>
                  <Input id="rule-priority" type="number" value={ruleForm.priority} onChange={e => setRuleForm({ ...ruleForm, priority: e.target.value })} />
                </div>
                <div className="space-y-2 col-span-2">
                  <label className="flex items-center gap-2 text-sm text-slate-700">
                    <input type="checkbox" checked={ruleForm.requires_receipt} onChange={e => setRuleForm({ ...ruleForm, requires_receipt: e.target.checked })} className="h-4 w-4 rounded border-slate-300 text-brand-600 focus:ring-brand-500" />
                    Requiere recibo
                  </label>
                </div>
                <div className="space-y-2 col-span-2">
                  <label className="flex items-center gap-2 text-sm text-slate-700">
                    <input type="checkbox" checked={ruleForm.requires_approval} onChange={e => setRuleForm({ ...ruleForm, requires_approval: e.target.checked })} className="h-4 w-4 rounded border-slate-300 text-brand-600 focus:ring-brand-500" />
                    Requiere aprobación
                  </label>
                </div>
              </div>
              <div className="mt-4">
                <Button size="sm" onClick={handleSaveRule} disabled={savingRule}>
                  <Plus size={14} className="mr-1" /> {savingRule ? 'Creando...' : 'Agregar regla'}
                </Button>
              </div>
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowRules(false)}>Cerrar</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* --- Dialogo accion de flujo --- */}
      <Dialog open={!!pendingAction} onOpenChange={o => { if (!o) setPendingAction(null) }}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>{pendingAction ? `${pendingAction.meta.label} ${entityNoun[pendingAction.entityType]}` : ''}</DialogTitle>
            <DialogDescription>Confirmá la acción sobre el registro seleccionado</DialogDescription>
          </DialogHeader>

          {pendingAction?.meta.needsInput && (
            <div className="space-y-2">
              <Label htmlFor="action-input">{pendingAction.meta.inputLabel} *</Label>
              <Input
                id="action-input"
                type={pendingAction.meta.inputType === 'number' ? 'number' : 'text'}
                step={pendingAction.meta.inputType === 'number' ? '0.01' : undefined}
                placeholder={pendingAction.meta.placeholder}
                value={actionInput}
                onChange={e => setActionInput(e.target.value)}
              />
            </div>
          )}

          <DialogFooter>
            <Button variant="outline" onClick={() => setPendingAction(null)}>Cancelar</Button>
            <Button
              onClick={runWorkflow}
              disabled={actionBusy || (!!pendingAction?.meta.needsInput && !actionInput.trim())}
            >
              {actionBusy ? 'Procesando...' : pendingAction?.meta.label}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
