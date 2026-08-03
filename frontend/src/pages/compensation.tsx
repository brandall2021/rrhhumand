import { useEffect, useState } from 'react'
import { Plus, Pencil } from 'lucide-react'
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

interface Structure {
  id: string
  name: string
  description?: string | null
  currency: string
  effective_from?: string
  effective_to?: string | null
  status: string
}

interface Grade {
  id: string
  structure_id: string
  code: string
  name: string
  sort_order?: number
  status: string
}

interface Band {
  id: string
  structure_id: string
  grade_id?: string | null
  code: string
  name: string
  minimum_amount?: number | string
  midpoint_amount?: number | string
  maximum_amount?: number | string
  currency: string
  status: string
}

interface Component {
  id: string
  code: string
  name: string
  description?: string | null
  component_type: string
  taxable: boolean
  recurring: boolean
  active: boolean
}

interface Adjustment {
  id: string
  employee_id: string
  adjustment_type: string
  value?: number | string
  currency: string
  reason: string
  effective_from?: string
  status: string
  notes?: string | null
}

interface Bonus {
  id: string
  employee_id: string
  bonus_plan_id?: string | null
  name: string
  bonus_type: string
  amount?: number | string
  currency: string
  period?: string | null
  reason?: string | null
  status: string
}

interface BonusPlan {
  id: string
  name: string
  description?: string | null
  period: string
  target_percentage?: number | string | null
  maximum_percentage?: number | string | null
  eligibility_rules?: string | null
  status: string
}

interface Benefit {
  id: string
  code: string
  name: string
  description?: string | null
  benefit_type: string
  provider?: string | null
  cost_amount?: number | string | null
  cost_currency: string
  frequency: string
  taxable: boolean
  active: boolean
}

interface Budget {
  id: string
  year: number
  department_id?: string | null
  budget_amount?: number | string
  committed_amount?: number | string
  spent_amount?: number | string
  currency: string
  status: string
}

interface DashboardStats {
  total_salary_cost?: number | string
  average_compensation?: number | string
  benefit_cost?: number | string
  total_bonuses?: number | string
  budget_total?: number | string
  budget_used?: number | string
  pending_proposals?: number
  employees_out_of_band?: number
  active_reviews?: number
  currency?: string
}

interface SelectOption {
  value: string
  label: string
}

const statusStyles: Record<string, { label: string; cls: string }> = {
  draft: { label: 'Borrador', cls: 'bg-slate-100 text-slate-600' },
  pending_approval: { label: 'Pendiente', cls: 'bg-amber-50 text-amber-700' },
  pending: { label: 'Pendiente', cls: 'bg-amber-50 text-amber-700' },
  active: { label: 'Activo', cls: 'bg-emerald-50 text-emerald-700' },
  inactive: { label: 'Inactivo', cls: 'bg-slate-100 text-slate-600' },
  approved: { label: 'Aprobado', cls: 'bg-emerald-50 text-emerald-700' },
  rejected: { label: 'Rechazado', cls: 'bg-red-50 text-red-700' },
  applied: { label: 'Aplicado', cls: 'bg-teal-50 text-teal-700' },
  closed: { label: 'Cerrado', cls: 'bg-slate-100 text-slate-500' },
}

const statusBadge = (status: string) => {
  const s = statusStyles[(status || '').toLowerCase()] ?? { label: status, cls: 'bg-slate-100 text-slate-600' }
  return <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${s.cls}`}>{s.label}</span>
}

const fmtDate = (s?: string) => (s ? s.slice(0, 10) : '-')

const fmtMoney = (v?: number | string) => {
  if (v == null || v === '') return '-'
  const n = parseFloat(String(v))
  return isNaN(n) ? '-' : n.toLocaleString('es-AR')
}

const unwrapList = (res: any): any[] => res?.data?.data ?? res?.data ?? []

const currencyOptions: SelectOption[] = [
  { value: 'USD', label: 'USD' },
  { value: 'ARS', label: 'ARS' },
  { value: 'EUR', label: 'EUR' },
]

const componentTypeOptions: SelectOption[] = [
  { value: 'salary', label: 'Salario' },
  { value: 'fixed', label: 'Fijo' },
  { value: 'variable', label: 'Variable' },
  { value: 'bonus', label: 'Bono' },
  { value: 'allowance', label: 'Asignación' },
  { value: 'deduction', label: 'Deducción' },
]

const adjustmentTypeOptions: SelectOption[] = [
  { value: 'percentage', label: 'Porcentaje' },
  { value: 'fixed_amount', label: 'Monto fijo' },
  { value: 'new_salary', label: 'Nuevo salario' },
]

const bonusTypeOptions: SelectOption[] = [
  { value: 'discretionary', label: 'Discrecional' },
  { value: 'performance', label: 'Desempeño' },
  { value: 'referral', label: 'Referido' },
  { value: 'retention', label: 'Retención' },
  { value: 'signing', label: 'Firma' },
  { value: 'other', label: 'Otro' },
]

const benefitTypeOptions: SelectOption[] = [
  { value: 'health', label: 'Salud' },
  { value: 'insurance', label: 'Seguro' },
  { value: 'meal', label: 'Comida' },
  { value: 'transport', label: 'Transporte' },
  { value: 'training', label: 'Capacitación' },
  { value: 'other', label: 'Otro' },
]

const frequencyOptions: SelectOption[] = [
  { value: 'monthly', label: 'Mensual' },
  { value: 'quarterly', label: 'Trimestral' },
  { value: 'annual', label: 'Anual' },
  { value: 'one_time', label: 'Único' },
]

const periodOptions: SelectOption[] = [
  { value: 'annual', label: 'Anual' },
  { value: 'semester', label: 'Semestral' },
  { value: 'quarterly', label: 'Trimestral' },
  { value: 'one_time', label: 'Único' },
]

const emptyStructForm = {
  name: '',
  description: '',
  currency: 'USD',
  effective_from: '',
  effective_to: '',
}

const emptyGradeForm = {
  code: '',
  name: '',
  sort_order: '',
}

const emptyBandForm = {
  grade_id: '',
  code: '',
  name: '',
  minimum_amount: '',
  midpoint_amount: '',
  maximum_amount: '',
  currency: 'USD',
}

const emptyCompForm = {
  code: '',
  name: '',
  description: '',
  component_type: 'fixed',
  taxable: true,
  recurring: false,
}

const emptyAssignCompForm = {
  employee_id: '',
  component_id: '',
  amount: '',
  currency: 'USD',
  effective_from: '',
}

const emptyAdjForm = {
  employee_id: '',
  adjustment_type: 'percentage',
  value: '',
  currency: 'USD',
  reason: '',
  effective_from: '',
  notes: '',
}

const emptyBonusForm = {
  employee_id: '',
  bonus_plan_id: '',
  name: '',
  bonus_type: 'discretionary',
  amount: '',
  currency: 'USD',
  period: '',
  reason: '',
}

const emptyPlanForm = {
  name: '',
  description: '',
  period: 'annual',
  target_percentage: '',
  maximum_percentage: '',
  eligibility_rules: '',
}

const emptyBenefitForm = {
  code: '',
  name: '',
  description: '',
  benefit_type: 'health',
  provider: '',
  cost_amount: '',
  cost_currency: 'USD',
  frequency: 'monthly',
  taxable: false,
}

const emptyAssignBenefitForm = {
  employee_id: '',
  benefit_id: '',
  effective_from: '',
  employee_cost: '',
  company_cost: '',
  currency: 'USD',
}

const emptyBudgetForm = {
  year: '',
  department_id: '',
  budget_amount: '',
  currency: 'USD',
}

const adjustmentTypeLabel: Record<string, string> = {
  percentage: 'Porcentaje',
  fixed_amount: 'Monto fijo',
  new_salary: 'Nuevo salario',
}

export default function CompensationPage() {
  const [error, setError] = useState('')

  // Selects
  const [employees, setEmployees] = useState<SelectOption[]>([])
  const [employeeMap, setEmployeeMap] = useState<Record<string, string>>({})
  const [departments, setDepartments] = useState<SelectOption[]>([])
  const [departmentMap, setDepartmentMap] = useState<Record<string, string>>({})

  // Dashboard
  const [stats, setStats] = useState<DashboardStats>({})
  const [statsLoading, setStatsLoading] = useState(true)

  // Estructuras
  const [structures, setStructures] = useState<Structure[]>([])
  const [structLoading, setStructLoading] = useState(true)
  const [showStructModal, setShowStructModal] = useState(false)
  const [editingStruct, setEditingStruct] = useState<Structure | null>(null)
  const [structForm, setStructForm] = useState({ ...emptyStructForm })
  const [savingStruct, setSavingStruct] = useState(false)

  // Contexto de estructura (grados/bandas)
  const [structureId, setStructureId] = useState('')

  // Grados
  const [grades, setGrades] = useState<Grade[]>([])
  const [gradeLoading, setGradeLoading] = useState(true)
  const [showGradeModal, setShowGradeModal] = useState(false)
  const [editingGrade, setEditingGrade] = useState<Grade | null>(null)
  const [gradeForm, setGradeForm] = useState({ ...emptyGradeForm })
  const [savingGrade, setSavingGrade] = useState(false)

  // Bandas
  const [bands, setBands] = useState<Band[]>([])
  const [bandLoading, setBandLoading] = useState(true)
  const [showBandModal, setShowBandModal] = useState(false)
  const [editingBand, setEditingBand] = useState<Band | null>(null)
  const [bandForm, setBandForm] = useState({ ...emptyBandForm })
  const [savingBand, setSavingBand] = useState(false)

  // Componentes
  const [components, setComponents] = useState<Component[]>([])
  const [compLoading, setCompLoading] = useState(true)
  const [showCompModal, setShowCompModal] = useState(false)
  const [compForm, setCompForm] = useState({ ...emptyCompForm })
  const [savingComp, setSavingComp] = useState(false)
  const [showAssignComp, setShowAssignComp] = useState(false)
  const [assignCompForm, setAssignCompForm] = useState({ ...emptyAssignCompForm })
  const [savingAssignComp, setSavingAssignComp] = useState(false)

  // Ajustes
  const [adjustments, setAdjustments] = useState<Adjustment[]>([])
  const [adjLoading, setAdjLoading] = useState(true)
  const [showAdjModal, setShowAdjModal] = useState(false)
  const [adjForm, setAdjForm] = useState({ ...emptyAdjForm })
  const [savingAdj, setSavingAdj] = useState(false)
  const [adjActionBusy, setAdjActionBusy] = useState('')

  // Bonos + planes
  const [bonuses, setBonuses] = useState<Bonus[]>([])
  const [bonusLoading, setBonusLoading] = useState(true)
  const [showBonusModal, setShowBonusModal] = useState(false)
  const [bonusForm, setBonusForm] = useState({ ...emptyBonusForm })
  const [savingBonus, setSavingBonus] = useState(false)
  const [bonusActionBusy, setBonusActionBusy] = useState('')
  const [bonusPlans, setBonusPlans] = useState<BonusPlan[]>([])
  const [showPlanModal, setShowPlanModal] = useState(false)
  const [planForm, setPlanForm] = useState({ ...emptyPlanForm })
  const [savingPlan, setSavingPlan] = useState(false)

  // Beneficios
  const [benefits, setBenefits] = useState<Benefit[]>([])
  const [benefitLoading, setBenefitLoading] = useState(true)
  const [showBenefitModal, setShowBenefitModal] = useState(false)
  const [editingBenefit, setEditingBenefit] = useState<Benefit | null>(null)
  const [benefitForm, setBenefitForm] = useState({ ...emptyBenefitForm })
  const [savingBenefit, setSavingBenefit] = useState(false)
  const [showAssignBenefit, setShowAssignBenefit] = useState(false)
  const [assignBenefitForm, setAssignBenefitForm] = useState({ ...emptyAssignBenefitForm })
  const [savingAssignBenefit, setSavingAssignBenefit] = useState(false)

  // Presupuestos
  const [budgets, setBudgets] = useState<Budget[]>([])
  const [budgetLoading, setBudgetLoading] = useState(true)
  const [showBudgetModal, setShowBudgetModal] = useState(false)
  const [budgetForm, setBudgetForm] = useState({ ...emptyBudgetForm })
  const [savingBudget, setSavingBudget] = useState(false)

  const employeeName = (id?: string | null) => {
    if (!id) return '-'
    return employeeMap[id] ?? id.slice(0, 8)
  }

  const departmentName = (id?: string | null) => {
    if (!id) return '-'
    return departmentMap[id] ?? id.slice(0, 8)
  }

  const fetchSelects = async () => {
    try {
      const [eRes, dRes] = await Promise.all([
        api.get('/employees', { params: { limit: '500' } }),
        api.get('/departments', { params: { limit: '200' } }),
      ])
      const eList = unwrapList(eRes)
      const eMap: Record<string, string> = {}
      eList.forEach((em: any) => { eMap[em.id] = `${em.first_name} ${em.last_name}`.trim() })
      setEmployeeMap(eMap)
      setEmployees(eList.map((em: any) => ({ value: em.id, label: eMap[em.id] || em.id })))
      const dMap: Record<string, string> = {}
      unwrapList(dRes).forEach((d: any) => { dMap[d.id] = d.name })
      setDepartmentMap(dMap)
      setDepartments(unwrapList(dRes).map((d: any) => ({ value: d.id, label: d.name })))
    } catch {}
  }

  const fetchStats = async () => {
    setStatsLoading(true)
    try {
      const res = await api.get('/compensation/dashboard')
      setStats(res.data?.data ?? res.data ?? {})
      setError('')
    } catch {
      setStats({})
    } finally {
      setStatsLoading(false)
    }
  }

  const fetchStructures = async () => {
    setStructLoading(true)
    try {
      const list = unwrapList(await api.get('/compensation/structures'))
      setStructures(list)
      setError('')
      setStructureId((prev) => prev || (list.length ? list[0].id : ''))
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar estructuras')
      setStructures([])
    } finally {
      setStructLoading(false)
    }
  }

  const fetchGrades = async () => {
    if (!structureId) return
    setGradeLoading(true)
    try {
      setGrades(unwrapList(await api.get(`/compensation/structures/${structureId}/grades`)))
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar grados')
      setGrades([])
    } finally {
      setGradeLoading(false)
    }
  }

  const fetchBands = async () => {
    if (!structureId) return
    setBandLoading(true)
    try {
      setBands(unwrapList(await api.get(`/compensation/structures/${structureId}/bands`)))
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar bandas')
      setBands([])
    } finally {
      setBandLoading(false)
    }
  }

  const fetchComponents = async () => {
    setCompLoading(true)
    try {
      setComponents(unwrapList(await api.get('/compensation/components')))
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar componentes')
      setComponents([])
    } finally {
      setCompLoading(false)
    }
  }

  const fetchAdjustments = async () => {
    setAdjLoading(true)
    try {
      setAdjustments(unwrapList(await api.get('/compensation/adjustments')))
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar ajustes')
      setAdjustments([])
    } finally {
      setAdjLoading(false)
    }
  }

  const fetchBonuses = async () => {
    setBonusLoading(true)
    try {
      setBonuses(unwrapList(await api.get('/compensation/bonuses')))
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar bonos')
      setBonuses([])
    } finally {
      setBonusLoading(false)
    }
  }

  const fetchBonusPlans = async () => {
    try {
      setBonusPlans(unwrapList(await api.get('/compensation/bonus-plans')))
    } catch {
      setBonusPlans([])
    }
  }

  const fetchBenefits = async () => {
    setBenefitLoading(true)
    try {
      setBenefits(unwrapList(await api.get('/compensation/benefits')))
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar beneficios')
      setBenefits([])
    } finally {
      setBenefitLoading(false)
    }
  }

  const fetchBudgets = async () => {
    setBudgetLoading(true)
    try {
      setBudgets(unwrapList(await api.get('/compensation/budgets')))
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar presupuestos')
      setBudgets([])
    } finally {
      setBudgetLoading(false)
    }
  }

  useEffect(() => {
    fetchSelects()
    fetchStats()
    fetchStructures()
    fetchComponents()
    fetchAdjustments()
    fetchBonuses()
    fetchBonusPlans()
    fetchBenefits()
    fetchBudgets()
  }, [])

  useEffect(() => {
    fetchGrades()
    fetchBands()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [structureId])

  // --- Estructuras ---
  const openCreateStruct = () => {
    setEditingStruct(null)
    setStructForm({ ...emptyStructForm })
    setShowStructModal(true)
  }

  const openEditStruct = (s: Structure) => {
    setEditingStruct(s)
    setStructForm({
      name: s.name,
      description: s.description ?? '',
      currency: s.currency || 'USD',
      effective_from: s.effective_from ? s.effective_from.slice(0, 10) : '',
      effective_to: s.effective_to ? s.effective_to.slice(0, 10) : '',
    })
    setShowStructModal(true)
  }

  const handleSaveStruct = async () => {
    setSavingStruct(true)
    try {
      const body: Record<string, any> = {
        name: structForm.name,
        currency: structForm.currency,
        effective_from: structForm.effective_from,
      }
      if (structForm.description) body.description = structForm.description
      if (structForm.effective_to) body.effective_to = structForm.effective_to
      if (editingStruct) {
        body.status = editingStruct.status
        await api.put(`/compensation/structures/${editingStruct.id}`, body)
      } else {
        await api.post('/compensation/structures', body)
      }
      setShowStructModal(false)
      fetchStructures()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar estructura')
    } finally {
      setSavingStruct(false)
    }
  }

  // --- Grados ---
  const openCreateGrade = () => {
    setEditingGrade(null)
    setGradeForm({ ...emptyGradeForm })
    setShowGradeModal(true)
  }

  const openEditGrade = (g: Grade) => {
    setEditingGrade(g)
    setGradeForm({
      code: g.code,
      name: g.name,
      sort_order: g.sort_order != null ? String(g.sort_order) : '',
    })
    setShowGradeModal(true)
  }

  const handleSaveGrade = async () => {
    setSavingGrade(true)
    try {
      const body: Record<string, any> = {
        code: gradeForm.code,
        name: gradeForm.name,
      }
      if (gradeForm.sort_order) body.sort_order = parseInt(gradeForm.sort_order, 10)
      if (editingGrade) {
        body.status = editingGrade.status
        await api.put(`/compensation/grades/${editingGrade.id}`, body)
      } else {
        if (!structureId) {
          alert('Seleccioná una estructura primero')
          return
        }
        await api.post(`/compensation/structures/${structureId}/grades`, body)
      }
      setShowGradeModal(false)
      fetchGrades()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar grado')
    } finally {
      setSavingGrade(false)
    }
  }

  // --- Bandas ---
  const openCreateBand = () => {
    setEditingBand(null)
    setBandForm({ ...emptyBandForm })
    setShowBandModal(true)
  }

  const openEditBand = (b: Band) => {
    setEditingBand(b)
    setBandForm({
      grade_id: b.grade_id ?? '',
      code: b.code,
      name: b.name,
      minimum_amount: b.minimum_amount != null ? String(b.minimum_amount) : '',
      midpoint_amount: b.midpoint_amount != null ? String(b.midpoint_amount) : '',
      maximum_amount: b.maximum_amount != null ? String(b.maximum_amount) : '',
      currency: b.currency || 'USD',
    })
    setShowBandModal(true)
  }

  const handleSaveBand = async () => {
    setSavingBand(true)
    try {
      const body: Record<string, any> = {
        code: bandForm.code,
        name: bandForm.name,
        minimum_amount: parseFloat(bandForm.minimum_amount),
        midpoint_amount: parseFloat(bandForm.midpoint_amount),
        maximum_amount: parseFloat(bandForm.maximum_amount),
        currency: bandForm.currency,
      }
      if (bandForm.grade_id) body.grade_id = bandForm.grade_id
      if (editingBand) {
        body.status = editingBand.status
        await api.put(`/compensation/bands/${editingBand.id}`, body)
      } else {
        if (!structureId) {
          alert('Seleccioná una estructura primero')
          return
        }
        await api.post(`/compensation/structures/${structureId}/bands`, body)
      }
      setShowBandModal(false)
      fetchBands()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar banda')
    } finally {
      setSavingBand(false)
    }
  }

  // --- Componentes ---
  const openCreateComp = () => {
    setCompForm({ ...emptyCompForm })
    setShowCompModal(true)
  }

  const handleSaveComp = async () => {
    setSavingComp(true)
    try {
      const body: Record<string, any> = {
        code: compForm.code,
        name: compForm.name,
        component_type: compForm.component_type,
        taxable: compForm.taxable,
        recurring: compForm.recurring,
      }
      if (compForm.description) body.description = compForm.description
      await api.post('/compensation/components', body)
      setShowCompModal(false)
      fetchComponents()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar componente')
    } finally {
      setSavingComp(false)
    }
  }

  const openAssignComp = () => {
    setAssignCompForm({ ...emptyAssignCompForm })
    setShowAssignComp(true)
  }

  const handleAssignComp = async () => {
    setSavingAssignComp(true)
    try {
      const body: Record<string, any> = {
        component_id: assignCompForm.component_id,
        amount: parseFloat(assignCompForm.amount),
        currency: assignCompForm.currency,
        effective_from: assignCompForm.effective_from,
      }
      await api.post(`/compensation/employees/${assignCompForm.employee_id}/components`, body)
      setShowAssignComp(false)
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al asignar componente')
    } finally {
      setSavingAssignComp(false)
    }
  }

  // --- Ajustes ---
  const openCreateAdj = () => {
    setAdjForm({ ...emptyAdjForm })
    setShowAdjModal(true)
  }

  const handleSaveAdj = async () => {
    setSavingAdj(true)
    try {
      const body: Record<string, any> = {
        employee_id: adjForm.employee_id,
        adjustment_type: adjForm.adjustment_type,
        value: parseFloat(adjForm.value),
        currency: adjForm.currency,
        reason: adjForm.reason,
        effective_from: adjForm.effective_from,
      }
      if (adjForm.notes) body.notes = adjForm.notes
      await api.post('/compensation/adjustments', body)
      setShowAdjModal(false)
      fetchAdjustments()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar ajuste')
    } finally {
      setSavingAdj(false)
    }
  }

  const runAdjustmentAction = async (a: Adjustment, action: 'approve' | 'reject' | 'apply') => {
    const label = { approve: 'aprobar', reject: 'rechazar', apply: 'aplicar' }[action]
    if (!confirm(`¿${label[0].toUpperCase() + label.slice(1)} el ajuste?`)) return
    setAdjActionBusy(a.id)
    try {
      await api.post(`/compensation/adjustments/${a.id}/${action}`)
      fetchAdjustments()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al ejecutar la acción')
    } finally {
      setAdjActionBusy('')
    }
  }

  const adjustmentActions = (a: Adjustment) => {
    const st = (a.status || '').toLowerCase()
    const actions: { action: 'approve' | 'reject' | 'apply'; label: string; cls: string }[] = []
    if (st === 'draft' || st === 'pending_approval' || st === 'pending') {
      actions.push({ action: 'approve', label: 'Aprobar', cls: 'text-emerald-600' })
      actions.push({ action: 'reject', label: 'Rechazar', cls: 'text-red-600' })
    }
    if (st === 'approved') {
      actions.push({ action: 'apply', label: 'Aplicar', cls: 'text-blue-600' })
    }
    return actions.map(b => (
      <Button
        key={b.action}
        variant="ghost"
        size="sm"
        className={b.cls}
        disabled={adjActionBusy === a.id}
        onClick={() => runAdjustmentAction(a, b.action)}
      >
        {b.label}
      </Button>
    ))
  }

  // --- Bonos ---
  const openCreateBonus = () => {
    setBonusForm({ ...emptyBonusForm })
    setShowBonusModal(true)
  }

  const handleSaveBonus = async () => {
    setSavingBonus(true)
    try {
      const body: Record<string, any> = {
        employee_id: bonusForm.employee_id,
        name: bonusForm.name,
        bonus_type: bonusForm.bonus_type,
        amount: parseFloat(bonusForm.amount),
        currency: bonusForm.currency,
      }
      if (bonusForm.bonus_plan_id) body.bonus_plan_id = bonusForm.bonus_plan_id
      if (bonusForm.period) body.period = bonusForm.period
      if (bonusForm.reason) body.reason = bonusForm.reason
      await api.post('/compensation/bonuses', body)
      setShowBonusModal(false)
      fetchBonuses()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar bono')
    } finally {
      setSavingBonus(false)
    }
  }

  const runBonusAction = async (b: Bonus, action: 'approve' | 'reject') => {
    const label = { approve: 'aprobar', reject: 'rechazar' }[action]
    if (!confirm(`¿${label[0].toUpperCase() + label.slice(1)} el bono "${b.name}"?`)) return
    setBonusActionBusy(b.id)
    try {
      await api.post(`/compensation/bonuses/${b.id}/${action}`)
      fetchBonuses()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al ejecutar la acción')
    } finally {
      setBonusActionBusy('')
    }
  }

  const bonusActions = (b: Bonus) => {
    const st = (b.status || '').toLowerCase()
    if (st !== 'draft' && st !== 'pending_approval' && st !== 'pending') return null
    return (
      <>
        <Button
          variant="ghost"
          size="sm"
          className="text-emerald-600"
          disabled={bonusActionBusy === b.id}
          onClick={() => runBonusAction(b, 'approve')}
        >
          Aprobar
        </Button>
        <Button
          variant="ghost"
          size="sm"
          className="text-red-600"
          disabled={bonusActionBusy === b.id}
          onClick={() => runBonusAction(b, 'reject')}
        >
          Rechazar
        </Button>
      </>
    )
  }

  // --- Planes de bonos ---
  const openCreatePlan = () => {
    setPlanForm({ ...emptyPlanForm })
    setShowPlanModal(true)
  }

  const handleSavePlan = async () => {
    setSavingPlan(true)
    try {
      const body: Record<string, any> = { name: planForm.name, period: planForm.period }
      if (planForm.description) body.description = planForm.description
      if (planForm.target_percentage) body.target_percentage = parseFloat(planForm.target_percentage)
      if (planForm.maximum_percentage) body.maximum_percentage = parseFloat(planForm.maximum_percentage)
      if (planForm.eligibility_rules) body.eligibility_rules = planForm.eligibility_rules
      await api.post('/compensation/bonus-plans', body)
      setShowPlanModal(false)
      fetchBonusPlans()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar plan de bonos')
    } finally {
      setSavingPlan(false)
    }
  }

  // --- Beneficios ---
  const openCreateBenefit = () => {
    setEditingBenefit(null)
    setBenefitForm({ ...emptyBenefitForm })
    setShowBenefitModal(true)
  }

  const openEditBenefit = (b: Benefit) => {
    setEditingBenefit(b)
    setBenefitForm({
      code: b.code,
      name: b.name,
      description: b.description ?? '',
      benefit_type: b.benefit_type || 'health',
      provider: b.provider ?? '',
      cost_amount: b.cost_amount != null ? String(b.cost_amount) : '',
      cost_currency: b.cost_currency || 'USD',
      frequency: b.frequency || 'monthly',
      taxable: b.taxable,
    })
    setShowBenefitModal(true)
  }

  const handleSaveBenefit = async () => {
    setSavingBenefit(true)
    try {
      const body: Record<string, any> = {
        code: benefitForm.code,
        name: benefitForm.name,
        benefit_type: benefitForm.benefit_type,
        frequency: benefitForm.frequency,
      }
      if (benefitForm.description) body.description = benefitForm.description
      if (benefitForm.provider) body.provider = benefitForm.provider
      if (benefitForm.cost_amount) body.cost_amount = parseFloat(benefitForm.cost_amount)
      if (benefitForm.cost_currency) body.cost_currency = benefitForm.cost_currency
      body.taxable = benefitForm.taxable
      if (editingBenefit) {
        body.active = editingBenefit.active
        await api.put(`/compensation/benefits/${editingBenefit.id}`, body)
      } else {
        await api.post('/compensation/benefits', body)
      }
      setShowBenefitModal(false)
      fetchBenefits()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar beneficio')
    } finally {
      setSavingBenefit(false)
    }
  }

  const openAssignBenefit = () => {
    setAssignBenefitForm({ ...emptyAssignBenefitForm })
    setShowAssignBenefit(true)
  }

  const handleAssignBenefit = async () => {
    setSavingAssignBenefit(true)
    try {
      const body: Record<string, any> = {
        benefit_id: assignBenefitForm.benefit_id,
        effective_from: assignBenefitForm.effective_from,
        currency: assignBenefitForm.currency,
      }
      if (assignBenefitForm.employee_cost) body.employee_cost = parseFloat(assignBenefitForm.employee_cost)
      if (assignBenefitForm.company_cost) body.company_cost = parseFloat(assignBenefitForm.company_cost)
      await api.post(`/compensation/employees/${assignBenefitForm.employee_id}/benefits`, body)
      setShowAssignBenefit(false)
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al asignar beneficio')
    } finally {
      setSavingAssignBenefit(false)
    }
  }

  // --- Presupuestos ---
  const openCreateBudget = () => {
    setBudgetForm({ ...emptyBudgetForm })
    setShowBudgetModal(true)
  }

  const handleSaveBudget = async () => {
    setSavingBudget(true)
    try {
      const body: Record<string, any> = {
        year: parseInt(budgetForm.year, 10),
        budget_amount: parseFloat(budgetForm.budget_amount),
        currency: budgetForm.currency,
      }
      if (budgetForm.department_id) body.department_id = budgetForm.department_id
      await api.post('/compensation/budgets', body)
      setShowBudgetModal(false)
      fetchBudgets()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar presupuesto')
    } finally {
      setSavingBudget(false)
    }
  }

  const gradeOptions: SelectOption[] = grades.map(g => ({ value: g.id, label: `${g.code} - ${g.name}` }))
  const structureOptions: SelectOption[] = structures.map(s => ({ value: s.id, label: s.name }))
  const planOptions: SelectOption[] = bonusPlans.map(p => ({ value: p.id, label: p.name }))

  const statCards = [
    { label: 'Costo salarial total', value: fmtMoney(stats.total_salary_cost), sub: stats.currency },
    { label: 'Compensación promedio', value: fmtMoney(stats.average_compensation), sub: stats.currency },
    { label: 'Costo de beneficios', value: fmtMoney(stats.benefit_cost), sub: stats.currency },
    { label: 'Total bonos', value: fmtMoney(stats.total_bonuses), sub: stats.currency },
    { label: 'Presupuesto total', value: fmtMoney(stats.budget_total), sub: stats.currency },
    { label: 'Presupuesto usado', value: fmtMoney(stats.budget_used), sub: stats.currency },
    { label: 'Propuestas pendientes', value: String(stats.pending_proposals ?? 0) },
    { label: 'Empleados fuera de banda', value: String(stats.employees_out_of_band ?? 0) },
    { label: 'Revisiones activas', value: String(stats.active_reviews ?? 0) },
  ]

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-slate-900">Compensaciones</h1>
      </div>

      {error && <div className="mb-4 p-3 bg-red-50 text-red-700 text-sm rounded-lg">{error}</div>}

      <Tabs defaultValue="dashboard">
        <TabsList>
          <TabsTrigger value="dashboard">Dashboard</TabsTrigger>
          <TabsTrigger value="structures">Estructuras</TabsTrigger>
          <TabsTrigger value="grades">Grados</TabsTrigger>
          <TabsTrigger value="bands">Bandas</TabsTrigger>
          <TabsTrigger value="components">Componentes</TabsTrigger>
          <TabsTrigger value="adjustments">Ajustes</TabsTrigger>
          <TabsTrigger value="bonuses">Bonos</TabsTrigger>
          <TabsTrigger value="benefits">Beneficios</TabsTrigger>
          <TabsTrigger value="budgets">Presupuestos</TabsTrigger>
        </TabsList>

        {/* ---------------- Dashboard ---------------- */}
        <TabsContent value="dashboard">
          <div className="mb-4 flex items-center justify-between">
            <span className="text-sm text-slate-500">Resumen de compensaciones {stats.currency ? `(${stats.currency})` : ''}</span>
            <Button size="sm" variant="outline" onClick={fetchStats}>Actualizar</Button>
          </div>
          {statsLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
          : <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-3 gap-4">
              {statCards.map(c => (
                <Card key={c.label}>
                  <CardContent className="p-4">
                    <div className="text-sm text-slate-500">{c.label}</div>
                    <div className="mt-1 text-xl font-bold text-slate-900">{c.value}</div>
                    {c.sub && <div className="text-xs text-slate-400">{c.sub}</div>}
                  </CardContent>
                </Card>
              ))}
            </div>}
        </TabsContent>

        {/* ---------------- Estructuras ---------------- */}
        <TabsContent value="structures">
          <div className="flex items-center justify-between mb-4">
            <span className="text-sm text-slate-500">{structures.length} estructuras</span>
            <Button size="sm" onClick={openCreateStruct}><Plus size={16} className="mr-1" /> Nueva</Button>
          </div>

          <Card>
            <CardContent className="p-0">
              {structLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
              : structures.length === 0 && !error ? <div className="p-6 text-center text-slate-500">No hay estructuras salariales</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Descripción</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Moneda</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Vigencia</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                        <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                      </tr>
                    </thead>
                    <tbody>
                      {structures.map(s => (
                        <tr key={s.id} className="border-b border-slate-100 hover:bg-slate-50">
                          <td className="px-4 py-3 font-medium text-slate-900">{s.name}</td>
                          <td className="px-4 py-3 text-slate-600">{s.description || '-'}</td>
                          <td className="px-4 py-3 text-slate-600">{s.currency || '-'}</td>
                          <td className="px-4 py-3 text-slate-600">
                            {fmtDate(s.effective_from)}{s.effective_to ? ` → ${fmtDate(s.effective_to)}` : ''}
                          </td>
                          <td className="px-4 py-3">{statusBadge(s.status)}</td>
                          <td className="px-4 py-3 text-right">
                            <Button variant="ghost" size="sm" title="Editar" onClick={() => openEditStruct(s)}>
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

        {/* ---------------- Grados ---------------- */}
        <TabsContent value="grades">
          <div className="flex items-center gap-3 mb-4">
            <div className="w-64">
              <Label className="mb-1 block text-xs text-slate-500">Estructura</Label>
              <Select options={structureOptions} placeholder="Seleccionar..." value={structureId} onChange={e => setStructureId(e.target.value)} />
            </div>
            <div className="ml-auto">
              <Button size="sm" onClick={openCreateGrade}><Plus size={16} className="mr-1" /> Nuevo</Button>
            </div>
          </div>

          <Card>
            <CardContent className="p-0">
              {gradeLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
              : !structureId ? <div className="p-6 text-center text-slate-500">Seleccioná una estructura salarial</div>
              : grades.length === 0 ? <div className="p-6 text-center text-slate-500">No hay grados en esta estructura</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Código</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Orden</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                        <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                      </tr>
                    </thead>
                    <tbody>
                      {grades.map(g => (
                        <tr key={g.id} className="border-b border-slate-100 hover:bg-slate-50">
                          <td className="px-4 py-3 text-slate-500">{g.code}</td>
                          <td className="px-4 py-3 font-medium text-slate-900">{g.name}</td>
                          <td className="px-4 py-3 text-slate-600">{g.sort_order ?? '-'}</td>
                          <td className="px-4 py-3">{statusBadge(g.status)}</td>
                          <td className="px-4 py-3 text-right">
                            <Button variant="ghost" size="sm" title="Editar" onClick={() => openEditGrade(g)}>
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

        {/* ---------------- Bandas ---------------- */}
        <TabsContent value="bands">
          <div className="flex items-center gap-3 mb-4">
            <div className="w-64">
              <Label className="mb-1 block text-xs text-slate-500">Estructura</Label>
              <Select options={structureOptions} placeholder="Seleccionar..." value={structureId} onChange={e => setStructureId(e.target.value)} />
            </div>
            <div className="ml-auto">
              <Button size="sm" onClick={openCreateBand}><Plus size={16} className="mr-1" /> Nueva</Button>
            </div>
          </div>

          <Card>
            <CardContent className="p-0">
              {bandLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
              : !structureId ? <div className="p-6 text-center text-slate-500">Seleccioná una estructura salarial</div>
              : bands.length === 0 ? <div className="p-6 text-center text-slate-500">No hay bandas en esta estructura</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Código</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Grado</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Mínimo</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Medio</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Máximo</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Moneda</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                        <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                      </tr>
                    </thead>
                    <tbody>
                      {bands.map(b => (
                        <tr key={b.id} className="border-b border-slate-100 hover:bg-slate-50">
                          <td className="px-4 py-3 text-slate-500">{b.code}</td>
                          <td className="px-4 py-3 font-medium text-slate-900">{b.name}</td>
                          <td className="px-4 py-3 text-slate-600">{b.grade_id ? gradeOptions.find(o => o.value === b.grade_id)?.label?.split(' - ')[1] ?? b.grade_id.slice(0, 8) : '-'}</td>
                          <td className="px-4 py-3 text-slate-600">{fmtMoney(b.minimum_amount)}</td>
                          <td className="px-4 py-3 text-slate-600">{fmtMoney(b.midpoint_amount)}</td>
                          <td className="px-4 py-3 text-slate-600">{fmtMoney(b.maximum_amount)}</td>
                          <td className="px-4 py-3 text-slate-600">{b.currency || '-'}</td>
                          <td className="px-4 py-3">{statusBadge(b.status)}</td>
                          <td className="px-4 py-3 text-right">
                            <Button variant="ghost" size="sm" title="Editar" onClick={() => openEditBand(b)}>
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

        {/* ---------------- Componentes ---------------- */}
        <TabsContent value="components">
          <div className="flex items-center justify-between mb-4">
            <span className="text-sm text-slate-500">{components.length} componentes</span>
            <div className="flex gap-2">
              <Button size="sm" variant="outline" onClick={openAssignComp}>Asignar a empleado</Button>
              <Button size="sm" onClick={openCreateComp}><Plus size={16} className="mr-1" /> Nuevo</Button>
            </div>
          </div>

          <Card>
            <CardContent className="p-0">
              {compLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
              : components.length === 0 ? <div className="p-6 text-center text-slate-500">No hay componentes</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Código</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Tipo</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Gravable</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Recurrente</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                      </tr>
                    </thead>
                    <tbody>
                      {components.map(c => (
                        <tr key={c.id} className="border-b border-slate-100 hover:bg-slate-50">
                          <td className="px-4 py-3 text-slate-500">{c.code}</td>
                          <td className="px-4 py-3 font-medium text-slate-900">{c.name}</td>
                          <td className="px-4 py-3 text-slate-600">{c.component_type || '-'}</td>
                          <td className="px-4 py-3 text-slate-600">{c.taxable ? 'Sí' : 'No'}</td>
                          <td className="px-4 py-3 text-slate-600">{c.recurring ? 'Sí' : 'No'}</td>
                          <td className="px-4 py-3">{statusBadge(c.active ? 'active' : 'inactive')}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>}
            </CardContent>
          </Card>
        </TabsContent>

        {/* ---------------- Ajustes ---------------- */}
        <TabsContent value="adjustments">
          <div className="flex items-center justify-between mb-4">
            <span className="text-sm text-slate-500">{adjustments.length} ajustes</span>
            <Button size="sm" onClick={openCreateAdj}><Plus size={16} className="mr-1" /> Nuevo</Button>
          </div>

          <Card>
            <CardContent className="p-0">
              {adjLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
              : adjustments.length === 0 ? <div className="p-6 text-center text-slate-500">No hay ajustes</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Empleado</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Tipo</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Monto</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Moneda</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Motivo</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Vigencia</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                        <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                      </tr>
                    </thead>
                    <tbody>
                      {adjustments.map(a => (
                        <tr key={a.id} className="border-b border-slate-100 hover:bg-slate-50">
                          <td className="px-4 py-3 text-slate-600">{employeeName(a.employee_id)}</td>
                          <td className="px-4 py-3 text-slate-600">{adjustmentTypeLabel[a.adjustment_type] ?? a.adjustment_type}</td>
                          <td className="px-4 py-3 font-medium text-slate-900">{fmtMoney(a.value)}</td>
                          <td className="px-4 py-3 text-slate-600">{a.currency || '-'}</td>
                          <td className="px-4 py-3 text-slate-600">{a.reason || '-'}</td>
                          <td className="px-4 py-3 text-slate-600">{fmtDate(a.effective_from)}</td>
                          <td className="px-4 py-3">{statusBadge(a.status)}</td>
                          <td className="px-4 py-3 text-right whitespace-nowrap">
                            {adjustmentActions(a)}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>}
            </CardContent>
          </Card>
        </TabsContent>

        {/* ---------------- Bonos ---------------- */}
        <TabsContent value="bonuses">
          <div className="flex items-center justify-between mb-4">
            <span className="text-sm text-slate-500">{bonuses.length} bonos</span>
            <Button size="sm" onClick={openCreateBonus}><Plus size={16} className="mr-1" /> Nuevo bono</Button>
          </div>

          <Card>
            <CardContent className="p-0">
              {bonusLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
              : bonuses.length === 0 ? <div className="p-6 text-center text-slate-500">No hay bonos</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Empleado</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Tipo</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Monto</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Moneda</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Periodo</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                        <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                      </tr>
                    </thead>
                    <tbody>
                      {bonuses.map(b => (
                        <tr key={b.id} className="border-b border-slate-100 hover:bg-slate-50">
                          <td className="px-4 py-3 text-slate-600">{employeeName(b.employee_id)}</td>
                          <td className="px-4 py-3 font-medium text-slate-900">{b.name}</td>
                          <td className="px-4 py-3 text-slate-600">{b.bonus_type || '-'}</td>
                          <td className="px-4 py-3 font-medium text-slate-900">{fmtMoney(b.amount)}</td>
                          <td className="px-4 py-3 text-slate-600">{b.currency || '-'}</td>
                          <td className="px-4 py-3 text-slate-600">{b.period || '-'}</td>
                          <td className="px-4 py-3">{statusBadge(b.status)}</td>
                          <td className="px-4 py-3 text-right whitespace-nowrap">
                            {bonusActions(b)}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>}
            </CardContent>
          </Card>

          <div className="flex items-center justify-between mt-8 mb-4">
            <span className="text-sm text-slate-500">{bonusPlans.length} planes de bonos</span>
            <Button size="sm" variant="outline" onClick={openCreatePlan}><Plus size={16} className="mr-1" /> Nuevo plan</Button>
          </div>

          <Card>
            <CardContent className="p-0">
              {bonusPlans.length === 0 ? <div className="p-6 text-center text-slate-500">No hay planes de bonos</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Periodo</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Target %</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Máximo %</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                      </tr>
                    </thead>
                    <tbody>
                      {bonusPlans.map(p => (
                        <tr key={p.id} className="border-b border-slate-100 hover:bg-slate-50">
                          <td className="px-4 py-3 font-medium text-slate-900">{p.name}</td>
                          <td className="px-4 py-3 text-slate-600">{p.period || '-'}</td>
                          <td className="px-4 py-3 text-slate-600">{p.target_percentage != null ? `${fmtMoney(p.target_percentage)}%` : '-'}</td>
                          <td className="px-4 py-3 text-slate-600">{p.maximum_percentage != null ? `${fmtMoney(p.maximum_percentage)}%` : '-'}</td>
                          <td className="px-4 py-3">{statusBadge(p.status)}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>}
            </CardContent>
          </Card>
        </TabsContent>

        {/* ---------------- Beneficios ---------------- */}
        <TabsContent value="benefits">
          <div className="flex items-center justify-between mb-4">
            <span className="text-sm text-slate-500">{benefits.length} beneficios</span>
            <div className="flex gap-2">
              <Button size="sm" variant="outline" onClick={openAssignBenefit}>Asignar a empleado</Button>
              <Button size="sm" onClick={openCreateBenefit}><Plus size={16} className="mr-1" /> Nuevo</Button>
            </div>
          </div>

          <Card>
            <CardContent className="p-0">
              {benefitLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
              : benefits.length === 0 ? <div className="p-6 text-center text-slate-500">No hay beneficios</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Código</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Tipo</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Proveedor</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Costo</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Frecuencia</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                        <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                      </tr>
                    </thead>
                    <tbody>
                      {benefits.map(b => (
                        <tr key={b.id} className="border-b border-slate-100 hover:bg-slate-50">
                          <td className="px-4 py-3 text-slate-500">{b.code}</td>
                          <td className="px-4 py-3 font-medium text-slate-900">{b.name}</td>
                          <td className="px-4 py-3 text-slate-600">{b.benefit_type || '-'}</td>
                          <td className="px-4 py-3 text-slate-600">{b.provider || '-'}</td>
                          <td className="px-4 py-3 text-slate-600">{b.cost_amount != null ? `${fmtMoney(b.cost_amount)} ${b.cost_currency || ''}`.trim() : '-'}</td>
                          <td className="px-4 py-3 text-slate-600">{b.frequency || '-'}</td>
                          <td className="px-4 py-3">{statusBadge(b.active ? 'active' : 'inactive')}</td>
                          <td className="px-4 py-3 text-right">
                            <Button variant="ghost" size="sm" title="Editar" onClick={() => openEditBenefit(b)}>
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

        {/* ---------------- Presupuestos ---------------- */}
        <TabsContent value="budgets">
          <div className="flex items-center justify-between mb-4">
            <span className="text-sm text-slate-500">{budgets.length} presupuestos</span>
            <Button size="sm" onClick={openCreateBudget}><Plus size={16} className="mr-1" /> Nuevo</Button>
          </div>

          <Card>
            <CardContent className="p-0">
              {budgetLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
              : budgets.length === 0 ? <div className="p-6 text-center text-slate-500">No hay presupuestos</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Año</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Departamento</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Presupuesto</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Comprometido</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Gastado</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Moneda</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                      </tr>
                    </thead>
                    <tbody>
                      {budgets.map(b => (
                        <tr key={b.id} className="border-b border-slate-100 hover:bg-slate-50">
                          <td className="px-4 py-3 text-slate-600">{b.year}</td>
                          <td className="px-4 py-3 text-slate-600">{departmentName(b.department_id)}</td>
                          <td className="px-4 py-3 font-medium text-slate-900">{fmtMoney(b.budget_amount)}</td>
                          <td className="px-4 py-3 text-slate-600">{fmtMoney(b.committed_amount)}</td>
                          <td className="px-4 py-3 text-slate-600">{fmtMoney(b.spent_amount)}</td>
                          <td className="px-4 py-3 text-slate-600">{b.currency || '-'}</td>
                          <td className="px-4 py-3">{statusBadge(b.status)}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      {/* --- Dialogo nueva/editar estructura --- */}
      <Dialog open={showStructModal} onOpenChange={setShowStructModal}>
        <DialogContent className="max-w-xl">
          <DialogHeader>
            <DialogTitle>{editingStruct ? 'Editar Estructura Salarial' : 'Nueva Estructura Salarial'}</DialogTitle>
            <DialogDescription>
              {editingStruct ? 'Modificá los datos de la estructura' : 'Completá los datos para crear una nueva estructura'}
            </DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2 col-span-2">
              <Label htmlFor="struct-name">Nombre *</Label>
              <Input id="struct-name" value={structForm.name} onChange={e => setStructForm({ ...structForm, name: e.target.value })} required />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="struct-desc">Descripción</Label>
              <Input id="struct-desc" value={structForm.description} onChange={e => setStructForm({ ...structForm, description: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="struct-currency">Moneda</Label>
              <Select id="struct-currency" options={currencyOptions} value={structForm.currency} onChange={e => setStructForm({ ...structForm, currency: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="struct-from">Vigencia desde *</Label>
              <Input id="struct-from" type="date" value={structForm.effective_from} onChange={e => setStructForm({ ...structForm, effective_from: e.target.value })} required />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="struct-to">Vigencia hasta</Label>
              <Input id="struct-to" type="date" value={structForm.effective_to} onChange={e => setStructForm({ ...structForm, effective_to: e.target.value })} />
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowStructModal(false)}>Cancelar</Button>
            <Button onClick={handleSaveStruct} disabled={savingStruct || !structForm.name || !structForm.effective_from}>
              {savingStruct ? 'Guardando...' : editingStruct ? 'Guardar Cambios' : 'Crear Estructura'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* --- Dialogo nuevo/editar grado --- */}
      <Dialog open={showGradeModal} onOpenChange={setShowGradeModal}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>{editingGrade ? 'Editar Grado' : 'Nuevo Grado'}</DialogTitle>
            <DialogDescription>
              {editingGrade ? 'Modificá los datos del grado' : 'Completá los datos para crear un nuevo grado'}
            </DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="grade-code">Código *</Label>
              <Input id="grade-code" value={gradeForm.code} onChange={e => setGradeForm({ ...gradeForm, code: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="grade-order">Orden</Label>
              <Input id="grade-order" type="number" value={gradeForm.sort_order} onChange={e => setGradeForm({ ...gradeForm, sort_order: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="grade-name">Nombre *</Label>
              <Input id="grade-name" value={gradeForm.name} onChange={e => setGradeForm({ ...gradeForm, name: e.target.value })} required />
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowGradeModal(false)}>Cancelar</Button>
            <Button onClick={handleSaveGrade} disabled={savingGrade || !gradeForm.code || !gradeForm.name}>
              {savingGrade ? 'Guardando...' : editingGrade ? 'Guardar Cambios' : 'Crear Grado'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* --- Dialogo nueva/editar banda --- */}
      <Dialog open={showBandModal} onOpenChange={setShowBandModal}>
        <DialogContent className="max-w-xl">
          <DialogHeader>
            <DialogTitle>{editingBand ? 'Editar Banda Salarial' : 'Nueva Banda Salarial'}</DialogTitle>
            <DialogDescription>
              {editingBand ? 'Modificá los datos de la banda' : 'Completá los datos para crear una nueva banda'}
            </DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="band-code">Código *</Label>
              <Input id="band-code" value={bandForm.code} onChange={e => setBandForm({ ...bandForm, code: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="band-grade">Grado</Label>
              <Select id="band-grade" options={gradeOptions} placeholder="Sin grado" value={bandForm.grade_id} onChange={e => setBandForm({ ...bandForm, grade_id: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="band-name">Nombre *</Label>
              <Input id="band-name" value={bandForm.name} onChange={e => setBandForm({ ...bandForm, name: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="band-min">Mínimo *</Label>
              <Input id="band-min" type="number" step="0.01" value={bandForm.minimum_amount} onChange={e => setBandForm({ ...bandForm, minimum_amount: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="band-mid">Medio *</Label>
              <Input id="band-mid" type="number" step="0.01" value={bandForm.midpoint_amount} onChange={e => setBandForm({ ...bandForm, midpoint_amount: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="band-max">Máximo *</Label>
              <Input id="band-max" type="number" step="0.01" value={bandForm.maximum_amount} onChange={e => setBandForm({ ...bandForm, maximum_amount: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="band-currency">Moneda</Label>
              <Select id="band-currency" options={currencyOptions} value={bandForm.currency} onChange={e => setBandForm({ ...bandForm, currency: e.target.value })} />
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowBandModal(false)}>Cancelar</Button>
            <Button onClick={handleSaveBand} disabled={savingBand || !bandForm.code || !bandForm.name || !bandForm.minimum_amount || !bandForm.midpoint_amount || !bandForm.maximum_amount}>
              {savingBand ? 'Guardando...' : editingBand ? 'Guardar Cambios' : 'Crear Banda'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* --- Dialogo nuevo componente --- */}
      <Dialog open={showCompModal} onOpenChange={setShowCompModal}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>Nuevo Componente</DialogTitle>
            <DialogDescription>Completá los datos para crear un nuevo componente de compensación</DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="comp-code">Código *</Label>
              <Input id="comp-code" value={compForm.code} onChange={e => setCompForm({ ...compForm, code: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="comp-type">Tipo *</Label>
              <Select id="comp-type" options={componentTypeOptions} value={compForm.component_type} onChange={e => setCompForm({ ...compForm, component_type: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="comp-name">Nombre *</Label>
              <Input id="comp-name" value={compForm.name} onChange={e => setCompForm({ ...compForm, name: e.target.value })} required />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="comp-desc">Descripción</Label>
              <Input id="comp-desc" value={compForm.description} onChange={e => setCompForm({ ...compForm, description: e.target.value })} />
            </div>
            <label className="flex items-center gap-2 text-sm text-slate-700">
              <input type="checkbox" checked={compForm.taxable} onChange={e => setCompForm({ ...compForm, taxable: e.target.checked })} />
              Gravable
            </label>
            <label className="flex items-center gap-2 text-sm text-slate-700">
              <input type="checkbox" checked={compForm.recurring} onChange={e => setCompForm({ ...compForm, recurring: e.target.checked })} />
              Recurrente
            </label>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowCompModal(false)}>Cancelar</Button>
            <Button onClick={handleSaveComp} disabled={savingComp || !compForm.code || !compForm.name}>
              {savingComp ? 'Guardando...' : 'Crear Componente'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* --- Dialogo asignar componente --- */}
      <Dialog open={showAssignComp} onOpenChange={setShowAssignComp}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>Asignar Componente</DialogTitle>
            <DialogDescription>Asigná un componente de compensación a un empleado</DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2 col-span-2">
              <Label htmlFor="assigncomp-emp">Empleado *</Label>
              <Select id="assigncomp-emp" options={employees} placeholder="Seleccionar..." value={assignCompForm.employee_id} onChange={e => setAssignCompForm({ ...assignCompForm, employee_id: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="assigncomp-comp">Componente *</Label>
              <Select id="assigncomp-comp" options={components.map(c => ({ value: c.id, label: `${c.code} - ${c.name}` }))} placeholder="Seleccionar..." value={assignCompForm.component_id} onChange={e => setAssignCompForm({ ...assignCompForm, component_id: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="assigncomp-amount">Monto *</Label>
              <Input id="assigncomp-amount" type="number" step="0.01" value={assignCompForm.amount} onChange={e => setAssignCompForm({ ...assignCompForm, amount: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="assigncomp-currency">Moneda</Label>
              <Select id="assigncomp-currency" options={currencyOptions} value={assignCompForm.currency} onChange={e => setAssignCompForm({ ...assignCompForm, currency: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="assigncomp-from">Vigencia desde *</Label>
              <Input id="assigncomp-from" type="date" value={assignCompForm.effective_from} onChange={e => setAssignCompForm({ ...assignCompForm, effective_from: e.target.value })} required />
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowAssignComp(false)}>Cancelar</Button>
            <Button onClick={handleAssignComp} disabled={savingAssignComp || !assignCompForm.employee_id || !assignCompForm.component_id || !assignCompForm.amount || !assignCompForm.effective_from}>
              {savingAssignComp ? 'Guardando...' : 'Asignar'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* --- Dialogo nuevo ajuste --- */}
      <Dialog open={showAdjModal} onOpenChange={setShowAdjModal}>
        <DialogContent className="max-w-xl">
          <DialogHeader>
            <DialogTitle>Nuevo Ajuste</DialogTitle>
            <DialogDescription>Completá los datos para crear un nuevo ajuste de compensación</DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2 col-span-2">
              <Label htmlFor="adj-emp">Empleado *</Label>
              <Select id="adj-emp" options={employees} placeholder="Seleccionar..." value={adjForm.employee_id} onChange={e => setAdjForm({ ...adjForm, employee_id: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="adj-type">Tipo de ajuste *</Label>
              <Select id="adj-type" options={adjustmentTypeOptions} value={adjForm.adjustment_type} onChange={e => setAdjForm({ ...adjForm, adjustment_type: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="adj-value">Valor *</Label>
              <Input id="adj-value" type="number" step="0.01" value={adjForm.value} onChange={e => setAdjForm({ ...adjForm, value: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="adj-currency">Moneda</Label>
              <Select id="adj-currency" options={currencyOptions} value={adjForm.currency} onChange={e => setAdjForm({ ...adjForm, currency: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="adj-from">Vigencia desde *</Label>
              <Input id="adj-from" type="date" value={adjForm.effective_from} onChange={e => setAdjForm({ ...adjForm, effective_from: e.target.value })} required />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="adj-reason">Motivo *</Label>
              <Input id="adj-reason" value={adjForm.reason} onChange={e => setAdjForm({ ...adjForm, reason: e.target.value })} required />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="adj-notes">Notas</Label>
              <Input id="adj-notes" value={adjForm.notes} onChange={e => setAdjForm({ ...adjForm, notes: e.target.value })} />
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowAdjModal(false)}>Cancelar</Button>
            <Button onClick={handleSaveAdj} disabled={savingAdj || !adjForm.employee_id || !adjForm.value || !adjForm.reason || !adjForm.effective_from}>
              {savingAdj ? 'Guardando...' : 'Crear Ajuste'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* --- Dialogo nuevo bono --- */}
      <Dialog open={showBonusModal} onOpenChange={setShowBonusModal}>
        <DialogContent className="max-w-xl">
          <DialogHeader>
            <DialogTitle>Nuevo Bono</DialogTitle>
            <DialogDescription>Completá los datos para crear un nuevo bono</DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2 col-span-2">
              <Label htmlFor="bonus-emp">Empleado *</Label>
              <Select id="bonus-emp" options={employees} placeholder="Seleccionar..." value={bonusForm.employee_id} onChange={e => setBonusForm({ ...bonusForm, employee_id: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="bonus-name">Nombre *</Label>
              <Input id="bonus-name" value={bonusForm.name} onChange={e => setBonusForm({ ...bonusForm, name: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="bonus-type">Tipo</Label>
              <Select id="bonus-type" options={bonusTypeOptions} value={bonusForm.bonus_type} onChange={e => setBonusForm({ ...bonusForm, bonus_type: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="bonus-plan">Plan de bonos</Label>
              <Select id="bonus-plan" options={planOptions} placeholder="Sin plan" value={bonusForm.bonus_plan_id} onChange={e => setBonusForm({ ...bonusForm, bonus_plan_id: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="bonus-amount">Monto *</Label>
              <Input id="bonus-amount" type="number" step="0.01" value={bonusForm.amount} onChange={e => setBonusForm({ ...bonusForm, amount: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="bonus-currency">Moneda</Label>
              <Select id="bonus-currency" options={currencyOptions} value={bonusForm.currency} onChange={e => setBonusForm({ ...bonusForm, currency: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="bonus-period">Periodo</Label>
              <Select id="bonus-period" options={periodOptions} placeholder="Opcional" value={bonusForm.period} onChange={e => setBonusForm({ ...bonusForm, period: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="bonus-reason">Motivo</Label>
              <Input id="bonus-reason" value={bonusForm.reason} onChange={e => setBonusForm({ ...bonusForm, reason: e.target.value })} />
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowBonusModal(false)}>Cancelar</Button>
            <Button onClick={handleSaveBonus} disabled={savingBonus || !bonusForm.employee_id || !bonusForm.name || !bonusForm.amount}>
              {savingBonus ? 'Guardando...' : 'Crear Bono'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* --- Dialogo nuevo plan de bonos --- */}
      <Dialog open={showPlanModal} onOpenChange={setShowPlanModal}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>Nuevo Plan de Bonos</DialogTitle>
            <DialogDescription>Completá los datos para crear un nuevo plan de bonos</DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2 col-span-2">
              <Label htmlFor="plan-name">Nombre *</Label>
              <Input id="plan-name" value={planForm.name} onChange={e => setPlanForm({ ...planForm, name: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="plan-period">Periodo</Label>
              <Select id="plan-period" options={periodOptions} value={planForm.period} onChange={e => setPlanForm({ ...planForm, period: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="plan-target">Target %</Label>
              <Input id="plan-target" type="number" step="0.01" value={planForm.target_percentage} onChange={e => setPlanForm({ ...planForm, target_percentage: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="plan-max">Máximo %</Label>
              <Input id="plan-max" type="number" step="0.01" value={planForm.maximum_percentage} onChange={e => setPlanForm({ ...planForm, maximum_percentage: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="plan-rules">Reglas de elegibilidad</Label>
              <Input id="plan-rules" value={planForm.eligibility_rules} onChange={e => setPlanForm({ ...planForm, eligibility_rules: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="plan-desc">Descripción</Label>
              <Input id="plan-desc" value={planForm.description} onChange={e => setPlanForm({ ...planForm, description: e.target.value })} />
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowPlanModal(false)}>Cancelar</Button>
            <Button onClick={handleSavePlan} disabled={savingPlan || !planForm.name}>
              {savingPlan ? 'Guardando...' : 'Crear Plan'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* --- Dialogo nuevo/editar beneficio --- */}
      <Dialog open={showBenefitModal} onOpenChange={setShowBenefitModal}>
        <DialogContent className="max-w-xl">
          <DialogHeader>
            <DialogTitle>{editingBenefit ? 'Editar Beneficio' : 'Nuevo Beneficio'}</DialogTitle>
            <DialogDescription>
              {editingBenefit ? 'Modificá los datos del beneficio' : 'Completá los datos para crear un nuevo beneficio'}
            </DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="benefit-code">Código *</Label>
              <Input id="benefit-code" value={benefitForm.code} onChange={e => setBenefitForm({ ...benefitForm, code: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="benefit-type">Tipo *</Label>
              <Select id="benefit-type" options={benefitTypeOptions} value={benefitForm.benefit_type} onChange={e => setBenefitForm({ ...benefitForm, benefit_type: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="benefit-name">Nombre *</Label>
              <Input id="benefit-name" value={benefitForm.name} onChange={e => setBenefitForm({ ...benefitForm, name: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="benefit-provider">Proveedor</Label>
              <Input id="benefit-provider" value={benefitForm.provider} onChange={e => setBenefitForm({ ...benefitForm, provider: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="benefit-freq">Frecuencia</Label>
              <Select id="benefit-freq" options={frequencyOptions} value={benefitForm.frequency} onChange={e => setBenefitForm({ ...benefitForm, frequency: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="benefit-cost">Costo</Label>
              <Input id="benefit-cost" type="number" step="0.01" value={benefitForm.cost_amount} onChange={e => setBenefitForm({ ...benefitForm, cost_amount: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="benefit-cur">Moneda</Label>
              <Select id="benefit-cur" options={currencyOptions} value={benefitForm.cost_currency} onChange={e => setBenefitForm({ ...benefitForm, cost_currency: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="benefit-desc">Descripción</Label>
              <Input id="benefit-desc" value={benefitForm.description} onChange={e => setBenefitForm({ ...benefitForm, description: e.target.value })} />
            </div>
            <label className="flex items-center gap-2 text-sm text-slate-700 col-span-2">
              <input type="checkbox" checked={benefitForm.taxable} onChange={e => setBenefitForm({ ...benefitForm, taxable: e.target.checked })} />
              Gravable
            </label>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowBenefitModal(false)}>Cancelar</Button>
            <Button onClick={handleSaveBenefit} disabled={savingBenefit || !benefitForm.code || !benefitForm.name}>
              {savingBenefit ? 'Guardando...' : editingBenefit ? 'Guardar Cambios' : 'Crear Beneficio'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* --- Dialogo asignar beneficio --- */}
      <Dialog open={showAssignBenefit} onOpenChange={setShowAssignBenefit}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>Asignar Beneficio</DialogTitle>
            <DialogDescription>Asigná un beneficio a un empleado</DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2 col-span-2">
              <Label htmlFor="assignbenefit-emp">Empleado *</Label>
              <Select id="assignbenefit-emp" options={employees} placeholder="Seleccionar..." value={assignBenefitForm.employee_id} onChange={e => setAssignBenefitForm({ ...assignBenefitForm, employee_id: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="assignbenefit-benefit">Beneficio *</Label>
              <Select id="assignbenefit-benefit" options={benefits.map(b => ({ value: b.id, label: `${b.code} - ${b.name}` }))} placeholder="Seleccionar..." value={assignBenefitForm.benefit_id} onChange={e => setAssignBenefitForm({ ...assignBenefitForm, benefit_id: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="assignbenefit-from">Vigencia desde *</Label>
              <Input id="assignbenefit-from" type="date" value={assignBenefitForm.effective_from} onChange={e => setAssignBenefitForm({ ...assignBenefitForm, effective_from: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="assignbenefit-empcost">Costo empleado</Label>
              <Input id="assignbenefit-empcost" type="number" step="0.01" value={assignBenefitForm.employee_cost} onChange={e => setAssignBenefitForm({ ...assignBenefitForm, employee_cost: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="assignbenefit-compcost">Costo empresa</Label>
              <Input id="assignbenefit-compcost" type="number" step="0.01" value={assignBenefitForm.company_cost} onChange={e => setAssignBenefitForm({ ...assignBenefitForm, company_cost: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="assignbenefit-currency">Moneda</Label>
              <Select id="assignbenefit-currency" options={currencyOptions} value={assignBenefitForm.currency} onChange={e => setAssignBenefitForm({ ...assignBenefitForm, currency: e.target.value })} />
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowAssignBenefit(false)}>Cancelar</Button>
            <Button onClick={handleAssignBenefit} disabled={savingAssignBenefit || !assignBenefitForm.employee_id || !assignBenefitForm.benefit_id || !assignBenefitForm.effective_from}>
              {savingAssignBenefit ? 'Guardando...' : 'Asignar'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* --- Dialogo nuevo presupuesto --- */}
      <Dialog open={showBudgetModal} onOpenChange={setShowBudgetModal}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>Nuevo Presupuesto</DialogTitle>
            <DialogDescription>Completá los datos para crear un nuevo presupuesto de compensación</DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="budget-year">Año *</Label>
              <Input id="budget-year" type="number" value={budgetForm.year} onChange={e => setBudgetForm({ ...budgetForm, year: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="budget-dept">Departamento</Label>
              <Select id="budget-dept" options={departments} placeholder="Sin departamento" value={budgetForm.department_id} onChange={e => setBudgetForm({ ...budgetForm, department_id: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="budget-amount">Monto *</Label>
              <Input id="budget-amount" type="number" step="0.01" value={budgetForm.budget_amount} onChange={e => setBudgetForm({ ...budgetForm, budget_amount: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="budget-currency">Moneda</Label>
              <Select id="budget-currency" options={currencyOptions} value={budgetForm.currency} onChange={e => setBudgetForm({ ...budgetForm, currency: e.target.value })} />
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowBudgetModal(false)}>Cancelar</Button>
            <Button onClick={handleSaveBudget} disabled={savingBudget || !budgetForm.year || !budgetForm.budget_amount}>
              {savingBudget ? 'Guardando...' : 'Crear Presupuesto'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
