import { useEffect, useState } from 'react'
import { Plus, Radar } from 'lucide-react'
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

interface OvertimeRecord {
  id: string
  employee_id: string
  employee_name?: string
  work_date: string
  planned_minutes: number
  actual_minutes: number
  overtime_minutes: number
  approved_minutes: number
  compensated_minutes: number
  paid_minutes: number
  overtime_type: string
  status: string
  is_weekend: boolean
  is_holiday: boolean
  is_night: boolean
  reason?: string
}

interface OvertimeRequest {
  id: string
  employee_id: string
  employee_name?: string
  work_date: string
  requested_minutes: number
  approved_minutes: number
  reason: string
  status: string
  rejection_reason?: string
}

interface Compensation {
  id: string
  employee_id: string
  employee_name?: string
  work_date: string
  minutes: number
  reason: string
  status: string
  rejection_reason?: string
}

interface Policy {
  id: string
  name: string
  description?: string
  max_daily_minutes: number
  max_weekly_minutes: number
  max_monthly_minutes: number
  requires_approval: boolean
  allows_compensation: boolean
  allows_payment: boolean
  minimum_overtime_minutes: number
  rounding_minutes: number
  overtime_expiration_days: number
  weekend_multiplier: number
  holiday_multiplier: number
  night_multiplier: number
  status: string
}

interface Dashboard {
  total_detected: number
  total_pending: number
  total_approved: number
  total_rejected: number
  total_compensated: number
  total_paid: number
  total_minutes: number
  balance_minutes: number
}

function useEmployees() {
  const [employees, setEmployees] = useState<{ value: string; label: string }[]>([])
  useEffect(() => {
    api.get('/employees', { params: { limit: '200' } })
      .then(res => setEmployees((res.data.data ?? []).map((e: any) => ({ value: e.id, label: `${e.first_name} ${e.last_name}` }))))
      .catch(() => {})
  }, [])
  return employees
}

const statusPill = (status: string) => {
  const map: Record<string, string> = {
    PENDING: 'bg-amber-50 text-amber-700', APPROVED: 'bg-emerald-50 text-emerald-700',
    REJECTED: 'bg-rose-50 text-rose-700', CANCELLED: 'bg-slate-100 text-slate-500',
    PAID: 'bg-sky-50 text-sky-700', COMPENSATED: 'bg-violet-50 text-violet-700', DETECTED: 'bg-blue-50 text-blue-700',
  }
  return (
    <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${map[status] || 'bg-slate-100 text-slate-600'}`}>
      {status}
    </span>
  )
}

const minutesLabel = (m: number) => `${Math.floor(m / 60)}h ${m % 60}m`

export default function OvertimePage() {
  const [activeTab, setActiveTab] = useState('records')

  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-900 mb-6">Horas Extra</h1>
      <DashboardCards />
      <div className="flex gap-1 mb-6 mt-6 border-b border-slate-200">
        {[{ key: 'records', label: 'Registros' }, { key: 'requests', label: 'Solicitudes' }, { key: 'compensations', label: 'Compensaciones' }, { key: 'policies', label: 'Políticas' }, { key: 'balance', label: 'Balance' }].map((tab) => (
          <button
            key={tab.key}
            onClick={() => setActiveTab(tab.key)}
            className={`px-4 py-2.5 text-sm font-medium border-b-2 transition-colors ${
              activeTab === tab.key ? 'border-brand-600 text-brand-700' : 'border-transparent text-slate-500 hover:text-slate-700'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>
      {activeTab === 'records' && <RecordsTab />}
      {activeTab === 'requests' && <RequestsTab />}
      {activeTab === 'compensations' && <CompensationsTab />}
      {activeTab === 'policies' && <PoliciesTab />}
      {activeTab === 'balance' && <BalanceTab />}
    </div>
  )
}

function DashboardCards() {
  const [dash, setDash] = useState<Dashboard | null>(null)
  useEffect(() => {
    api.get('/overtime/dashboard').then(res => setDash(res.data.data)).catch(() => {})
  }, [])
  if (!dash) return null
  const cards = [
    { label: 'Detectadas', value: dash.total_detected },
    { label: 'Pendientes', value: dash.total_pending },
    { label: 'Aprobadas', value: dash.total_approved },
    { label: 'Compensadas', value: dash.total_compensated },
    { label: 'Pagadas', value: dash.total_paid },
    { label: 'Total minutos', value: minutesLabel(dash.total_minutes) },
    { label: 'Balance', value: minutesLabel(dash.balance_minutes) },
  ]
  return (
    <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-7 gap-3">
      {cards.map(c => (
        <Card key={c.label}>
          <CardContent className="p-4">
            <p className="text-xs text-slate-500 mb-1">{c.label}</p>
            <p className="text-lg font-bold text-slate-900">{c.value}</p>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}

function RecordsTab() {
  const [records, setRecords] = useState<OvertimeRecord[]>([])
  const [loading, setLoading] = useState(true)
  const [filters, setFilters] = useState({ status: '', overtime_type: '', date_from: '', date_to: '' })
  const [showDetect, setShowDetect] = useState(false)
  const [detectForm, setDetectForm] = useState({ date_from: '', date_to: '' })
  const [detecting, setDetecting] = useState(false)

  const fetch = async () => {
    setLoading(true)
    try {
      const params: Record<string, any> = {}
      for (const [k, v] of Object.entries(filters)) if (v) params[k] = v
      const res = await api.get('/overtime', { params })
      setRecords(res.data.data ?? [])
    } catch { setRecords([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetch() }, [filters])

  const approve = async (id: string) => {
    const minutes = window.prompt('Minutos a aprobar:', String(records.find(r => r.id === id)?.overtime_minutes ?? 0))
    if (minutes === null) return
    try {
      await api.post(`/overtime/${id}/approve`, { approved_minutes: parseInt(minutes) || 0 })
      fetch()
    } catch (err: any) { alert(err?.response?.data?.error || 'Error') }
  }

  const reject = async (id: string) => {
    const reason = window.prompt('Motivo del rechazo:')
    if (!reason) return
    try {
      await api.post(`/overtime/${id}/reject`, { reason })
      fetch()
    } catch (err: any) { alert(err?.response?.data?.error || 'Error') }
  }

  const handleDetect = async () => {
    setDetecting(true)
    try {
      await api.post('/overtime/detect', detectForm)
      setShowDetect(false)
      fetch()
    } catch (err: any) { alert(err?.response?.data?.error || 'Error al detectar') }
    finally { setDetecting(false) }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4 flex-wrap gap-2">
        <div className="flex items-end gap-3 flex-wrap">
          <div>
            <Label>Estado</Label>
            <Select
              options={[{ value: '', label: 'Todos' }, ...['PENDING', 'APPROVED', 'REJECTED', 'PAID', 'COMPENSATED'].map(s => ({ value: s, label: s }))]}
              value={filters.status}
              onChange={e => setFilters({ ...filters, status: e.target.value })}
            />
          </div>
          <div>
            <Label>Tipo</Label>
            <Select
              options={[{ value: '', label: 'Todos' }, ...['REGULAR', 'WEEKEND', 'HOLIDAY', 'NIGHT'].map(s => ({ value: s, label: s }))]}
              value={filters.overtime_type}
              onChange={e => setFilters({ ...filters, overtime_type: e.target.value })}
            />
          </div>
          <div>
            <Label>Desde</Label>
            <Input type="date" value={filters.date_from} onChange={e => setFilters({ ...filters, date_from: e.target.value })} />
          </div>
          <div>
            <Label>Hasta</Label>
            <Input type="date" value={filters.date_to} onChange={e => setFilters({ ...filters, date_to: e.target.value })} />
          </div>
        </div>
        <Button size="sm" onClick={() => { setDetectForm({ date_from: filters.date_from || '', date_to: filters.date_to || '' }); setShowDetect(true) }}>
          <Radar size={14} /> Detectar
        </Button>
      </div>

      <Card>
        <CardContent className="p-0">
          {loading ? (
            <div className="p-6 text-center text-slate-500">Cargando...</div>
          ) : records.length === 0 ? (
            <div className="p-6 text-center text-slate-500">No hay registros de horas extra</div>
          ) : (
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-slate-200 bg-slate-50">
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Fecha</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Empleado</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Tipo</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Horas extra</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                  <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                </tr>
              </thead>
              <tbody>
                {records.map(r => (
                  <tr key={r.id} className="border-b border-slate-100 hover:bg-slate-50">
                    <td className="px-4 py-3 text-slate-900">{new Date(r.work_date).toLocaleDateString('es-AR')}</td>
                    <td className="px-4 py-3 font-medium text-slate-900">{r.employee_name || r.employee_id}</td>
                    <td className="px-4 py-3 text-slate-600">{r.overtime_type}{r.is_weekend ? ' · finde' : ''}{r.is_holiday ? ' · feriado' : ''}{r.is_night ? ' · noche' : ''}</td>
                    <td className="px-4 py-3 text-slate-900 font-medium">{minutesLabel(r.overtime_minutes)}</td>
                    <td className="px-4 py-3">{statusPill(r.status)}</td>
                    <td className="px-4 py-3">
                      <div className="flex justify-end gap-1">
                        {['PENDING', 'DETECTED'].includes(r.status) && (
                          <>
                            <Button size="sm" variant="outline" onClick={() => approve(r.id)}>Aprobar</Button>
                            <Button size="sm" variant="ghost" onClick={() => reject(r.id)}>Rechazar</Button>
                          </>
                        )}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </CardContent>
      </Card>

      <Dialog open={showDetect} onOpenChange={setShowDetect}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Detectar Horas Extra</DialogTitle>
            <DialogDescription>Analiza el registro de asistencia en el rango y crea horas extra detectadas</DialogDescription>
          </DialogHeader>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label>Desde *</Label>
              <Input type="date" value={detectForm.date_from} onChange={e => setDetectForm({ ...detectForm, date_from: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label>Hasta *</Label>
              <Input type="date" value={detectForm.date_to} onChange={e => setDetectForm({ ...detectForm, date_to: e.target.value })} />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowDetect(false)}>Cancelar</Button>
            <Button onClick={handleDetect} disabled={detecting || !detectForm.date_from || !detectForm.date_to}>
              {detecting ? 'Detectando...' : 'Detectar'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function RequestsTab() {
  const [requests, setRequests] = useState<OvertimeRequest[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [saving, setSaving] = useState(false)
  const [form, setForm] = useState({ work_date: '', requested_minutes: '', reason: '' })

  const fetch = async () => {
    setLoading(true)
    try {
      const res = await api.get('/overtime/requests')
      setRequests(res.data.data ?? [])
    } catch { setRequests([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetch() }, [])

  const create = async () => {
    setSaving(true)
    try {
      await api.post('/overtime/requests', {
        work_date: form.work_date,
        requested_minutes: parseInt(form.requested_minutes) || 0,
        reason: form.reason,
      })
      setShowModal(false)
      setForm({ work_date: '', requested_minutes: '', reason: '' })
      fetch()
    } catch (err: any) { alert(err?.response?.data?.error || 'Error') }
    finally { setSaving(false) }
  }

  const approve = async (id: string) => {
    const minutes = window.prompt('Minutos a aprobar:', String(requests.find(r => r.id === id)?.requested_minutes ?? 0))
    if (minutes === null) return
    try {
      await api.post(`/overtime/requests/${id}/approve`, { approved_minutes: parseInt(minutes) || 0 })
      fetch()
    } catch (err: any) { alert(err?.response?.data?.error || 'Error') }
  }

  const reject = async (id: string) => {
    const reason = window.prompt('Motivo del rechazo:')
    if (!reason) return
    try {
      await api.post(`/overtime/requests/${id}/reject`, { reason })
      fetch()
    } catch (err: any) { alert(err?.response?.data?.error || 'Error') }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-slate-500">{requests.length} solicitudes</p>
        <Button size="sm" onClick={() => setShowModal(true)}><Plus size={14} /> Nueva Solicitud</Button>
      </div>
      <Card>
        <CardContent className="p-0">
          {loading ? (
            <div className="p-6 text-center text-slate-500">Cargando...</div>
          ) : requests.length === 0 ? (
            <div className="p-6 text-center text-slate-500">No hay solicitudes</div>
          ) : (
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-slate-200 bg-slate-50">
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Fecha</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Empleado</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Solicitadas</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Aprobadas</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Motivo</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                  <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                </tr>
              </thead>
              <tbody>
                {requests.map(r => (
                  <tr key={r.id} className="border-b border-slate-100 hover:bg-slate-50">
                    <td className="px-4 py-3 text-slate-900">{new Date(r.work_date).toLocaleDateString('es-AR')}</td>
                    <td className="px-4 py-3 font-medium text-slate-900">{r.employee_name || r.employee_id}</td>
                    <td className="px-4 py-3 text-slate-600">{minutesLabel(r.requested_minutes)}</td>
                    <td className="px-4 py-3 text-slate-600">{minutesLabel(r.approved_minutes)}</td>
                    <td className="px-4 py-3 text-slate-600 max-w-[200px] truncate" title={r.reason}>{r.reason}</td>
                    <td className="px-4 py-3">{statusPill(r.status)}</td>
                    <td className="px-4 py-3">
                      <div className="flex justify-end gap-1">
                        {r.status === 'PENDING' && (
                          <>
                            <Button size="sm" variant="outline" onClick={() => approve(r.id)}>Aprobar</Button>
                            <Button size="sm" variant="ghost" onClick={() => reject(r.id)}>Rechazar</Button>
                          </>
                        )}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </CardContent>
      </Card>

      <Dialog open={showModal} onOpenChange={setShowModal}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Nueva Solicitud de Horas Extra</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label>Fecha *</Label>
              <Input type="date" value={form.work_date} onChange={e => setForm({ ...form, work_date: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label>Minutos *</Label>
              <Input type="number" value={form.requested_minutes} onChange={e => setForm({ ...form, requested_minutes: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label>Motivo *</Label>
              <Input value={form.reason} onChange={e => setForm({ ...form, reason: e.target.value })} />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button>
            <Button onClick={create} disabled={saving || !form.work_date || !form.requested_minutes || !form.reason}>
              {saving ? 'Guardando...' : 'Crear'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function CompensationsTab() {
  const [comps, setComps] = useState<Compensation[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [saving, setSaving] = useState(false)
  const [form, setForm] = useState({ work_date: '', minutes: '', reason: '' })

  const fetch = async () => {
    setLoading(true)
    try {
      const res = await api.get('/compensations')
      setComps(res.data.data ?? [])
    } catch { setComps([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetch() }, [])

  const create = async () => {
    setSaving(true)
    try {
      await api.post('/compensations', { work_date: form.work_date, minutes: parseInt(form.minutes) || 0, reason: form.reason })
      setShowModal(false)
      setForm({ work_date: '', minutes: '', reason: '' })
      fetch()
    } catch (err: any) { alert(err?.response?.data?.error || 'Error') }
    finally { setSaving(false) }
  }

  const action = async (id: string, act: string) => {
    let body: Record<string, any> = {}
    if (act === 'reject') {
      const reason = window.prompt('Motivo del rechazo:')
      if (!reason) return
      body = { reason }
    }
    try {
      await api.post(`/compensations/${id}/${act}`, body)
      fetch()
    } catch (err: any) { alert(err?.response?.data?.error || 'Error') }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-slate-500">{comps.length} compensaciones</p>
        <Button size="sm" onClick={() => setShowModal(true)}><Plus size={14} /> Nueva Compensación</Button>
      </div>
      <Card>
        <CardContent className="p-0">
          {loading ? (
            <div className="p-6 text-center text-slate-500">Cargando...</div>
          ) : comps.length === 0 ? (
            <div className="p-6 text-center text-slate-500">No hay compensaciones</div>
          ) : (
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-slate-200 bg-slate-50">
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Fecha</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Empleado</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Minutos</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Motivo</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                  <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                </tr>
              </thead>
              <tbody>
                {comps.map(c => (
                  <tr key={c.id} className="border-b border-slate-100 hover:bg-slate-50">
                    <td className="px-4 py-3 text-slate-900">{new Date(c.work_date).toLocaleDateString('es-AR')}</td>
                    <td className="px-4 py-3 font-medium text-slate-900">{c.employee_name || c.employee_id}</td>
                    <td className="px-4 py-3 text-slate-900 font-medium">{minutesLabel(c.minutes)}</td>
                    <td className="px-4 py-3 text-slate-600 max-w-[200px] truncate" title={c.reason}>{c.reason}</td>
                    <td className="px-4 py-3">{statusPill(c.status)}</td>
                    <td className="px-4 py-3">
                      <div className="flex justify-end gap-1">
                        {c.status === 'PENDING' && (
                          <>
                            <Button size="sm" variant="outline" onClick={() => action(c.id, 'approve')}>Aprobar</Button>
                            <Button size="sm" variant="ghost" onClick={() => action(c.id, 'reject')}>Rechazar</Button>
                            <Button size="sm" variant="ghost" onClick={() => action(c.id, 'cancel')}>Cancelar</Button>
                          </>
                        )}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </CardContent>
      </Card>

      <Dialog open={showModal} onOpenChange={setShowModal}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Nueva Compensación</DialogTitle>
            <DialogDescription>Banco de horas: convertí horas extra en tiempo libre</DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label>Fecha *</Label>
              <Input type="date" value={form.work_date} onChange={e => setForm({ ...form, work_date: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label>Minutos *</Label>
              <Input type="number" value={form.minutes} onChange={e => setForm({ ...form, minutes: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label>Motivo *</Label>
              <Input value={form.reason} onChange={e => setForm({ ...form, reason: e.target.value })} />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button>
            <Button onClick={create} disabled={saving || !form.work_date || !form.minutes || !form.reason}>
              {saving ? 'Guardando...' : 'Crear'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function PoliciesTab() {
  const [policies, setPolicies] = useState<Policy[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [saving, setSaving] = useState(false)
  const emptyForm = { name: '', description: '', max_daily_minutes: '', max_weekly_minutes: '', max_monthly_minutes: '', minimum_overtime_minutes: '', rounding_minutes: '' }
  const [form, setForm] = useState(emptyForm)

  const fetch = async () => {
    setLoading(true)
    try {
      const res = await api.get('/overtime-policies')
      setPolicies(res.data.data ?? [])
    } catch { setPolicies([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetch() }, [])

  const create = async () => {
    setSaving(true)
    try {
      const payload: Record<string, any> = { name: form.name }
      if (form.description) payload.description = form.description
      for (const k of ['max_daily_minutes', 'max_weekly_minutes', 'max_monthly_minutes', 'minimum_overtime_minutes', 'rounding_minutes']) {
        if (form[k as keyof typeof form]) payload[k] = parseInt(form[k as keyof typeof form])
      }
      await api.post('/overtime-policies', payload)
      setShowModal(false)
      setForm(emptyForm)
      fetch()
    } catch (err: any) { alert(err?.response?.data?.error || 'Error') }
    finally { setSaving(false) }
  }

  const toggleStatus = async (p: Policy) => {
    try {
      await api.put(`/overtime-policies/${p.id}`, { status: p.status === 'ACTIVE' ? 'INACTIVE' : 'ACTIVE' })
      fetch()
    } catch (err: any) { alert(err?.response?.data?.error || 'Error') }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-slate-500">{policies.length} políticas</p>
        <Button size="sm" onClick={() => setShowModal(true)}><Plus size={14} /> Nueva Política</Button>
      </div>
      <div className="grid gap-4 md:grid-cols-2">
        {loading ? (
          <div className="p-6 text-center text-slate-500">Cargando...</div>
        ) : policies.length === 0 ? (
          <div className="p-6 text-center text-slate-500">No hay políticas</div>
        ) : policies.map(p => (
          <Card key={p.id}>
            <CardContent className="p-5">
              <div className="flex items-start justify-between mb-3">
                <div>
                  <h3 className="font-semibold text-slate-900">{p.name}</h3>
                  <p className="text-xs text-slate-500">{p.description || 'Sin descripción'}</p>
                </div>
                <button onClick={() => toggleStatus(p)} className={`text-xs font-medium px-2 py-0.5 rounded-full ${p.status === 'ACTIVE' ? 'bg-emerald-50 text-emerald-700' : 'bg-slate-100 text-slate-500'}`}>
                  {p.status}
                </button>
              </div>
              <div className="grid grid-cols-2 gap-2 text-sm">
                <div><span className="text-slate-500">Máx diario:</span> <span className="font-medium text-slate-900">{minutesLabel(p.max_daily_minutes)}</span></div>
                <div><span className="text-slate-500">Máx semanal:</span> <span className="font-medium text-slate-900">{minutesLabel(p.max_weekly_minutes)}</span></div>
                <div><span className="text-slate-500">Máx mensual:</span> <span className="font-medium text-slate-900">{minutesLabel(p.max_monthly_minutes)}</span></div>
                <div><span className="text-slate-500">Redondeo:</span> <span className="font-medium text-slate-900">{p.rounding_minutes} min</span></div>
                <div><span className="text-slate-500">Fin de semana:</span> <span className="font-medium text-slate-900">x{p.weekend_multiplier}</span></div>
                <div><span className="text-slate-500">Feriado:</span> <span className="font-medium text-slate-900">x{p.holiday_multiplier}</span></div>
                <div><span className="text-slate-500">Aprobación:</span> <span className="font-medium text-slate-900">{p.requires_approval ? 'Sí' : 'No'}</span></div>
                <div><span className="text-slate-500">Expiración:</span> <span className="font-medium text-slate-900">{p.overtime_expiration_days} días</span></div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      <Dialog open={showModal} onOpenChange={setShowModal}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Nueva Política de Horas Extra</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label>Nombre *</Label>
              <Input value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label>Descripción</Label>
              <Input value={form.description} onChange={e => setForm({ ...form, description: e.target.value })} />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Máx diario (min)</Label>
                <Input type="number" value={form.max_daily_minutes} onChange={e => setForm({ ...form, max_daily_minutes: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label>Máx semanal (min)</Label>
                <Input type="number" value={form.max_weekly_minutes} onChange={e => setForm({ ...form, max_weekly_minutes: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label>Máx mensual (min)</Label>
                <Input type="number" value={form.max_monthly_minutes} onChange={e => setForm({ ...form, max_monthly_minutes: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label>Mínimo (min)</Label>
                <Input type="number" value={form.minimum_overtime_minutes} onChange={e => setForm({ ...form, minimum_overtime_minutes: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label>Redondeo (min)</Label>
                <Input type="number" value={form.rounding_minutes} onChange={e => setForm({ ...form, rounding_minutes: e.target.value })} />
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button>
            <Button onClick={create} disabled={saving || !form.name}>{saving ? 'Guardando...' : 'Crear'}</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function BalanceTab() {
  const employees = useEmployees()
  const [employeeId, setEmployeeId] = useState('')
  const [balance, setBalance] = useState<{ balance_minutes: number } | null>(null)
  const [txs, setTxs] = useState<any[]>([])
  const [loading, setLoading] = useState(false)
  const [adjustForm, setAdjustForm] = useState({ minutes: '', reason: '' })

  const fetchBalance = async () => {
    if (!employeeId) { setBalance(null); setTxs([]); return }
    setLoading(true)
    try {
      const [bal, tx] = await Promise.all([
        api.get(`/employees/${employeeId}/time-balance`),
        api.get(`/employees/${employeeId}/time-balance/transactions`),
      ])
      setBalance(bal.data.data ?? null)
      setTxs(tx.data.data ?? [])
    } catch { setBalance(null); setTxs([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchBalance() }, [employeeId])

  const adjust = async () => {
    if (!employeeId) return
    try {
      await api.post(`/employees/${employeeId}/time-balance/adjust`, {
        minutes: parseInt(adjustForm.minutes) || 0,
        reason: adjustForm.reason,
      })
      setAdjustForm({ minutes: '', reason: '' })
      fetchBalance()
    } catch (err: any) { alert(err?.response?.data?.error || 'Error') }
  }

  return (
    <div className="grid gap-4 md:grid-cols-3">
      <Card className="md:col-span-2">
        <CardContent className="p-5">
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-semibold text-slate-900">Saldo de Horas</h3>
            <div className="w-64">
              <Select options={employees} placeholder="Seleccionar empleado..." value={employeeId} onChange={e => setEmployeeId(e.target.value)} />
            </div>
          </div>
          {loading ? (
            <p className="text-sm text-slate-500 py-4">Cargando...</p>
          ) : balance ? (
            <p className="text-3xl font-bold text-slate-900 mb-6">{minutesLabel(balance.balance_minutes)}</p>
          ) : (
            <p className="text-sm text-slate-500 py-4">Seleccioná un empleado para ver su saldo</p>
          )}
          <h4 className="text-sm font-medium text-slate-700 mb-2">Movimientos</h4>
          {txs.length === 0 ? (
            <p className="text-sm text-slate-500">Sin movimientos</p>
          ) : (
            <div className="space-y-2">
              {txs.map(t => (
                <div key={t.id} className="flex items-center justify-between text-sm border-b border-slate-100 pb-2">
                  <div>
                    <span className="font-medium text-slate-900">{t.transaction_type}</span>
                    <p className="text-xs text-slate-500">{new Date(t.created_at).toLocaleString('es-AR')}</p>
                  </div>
                  <span className={`font-semibold ${t.minutes >= 0 ? 'text-emerald-600' : 'text-rose-600'}`}>{t.minutes >= 0 ? '+' : ''}{minutesLabel(t.minutes)}</span>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
      <Card>
        <CardContent className="p-5 space-y-4">
          <h3 className="font-semibold text-slate-900">Ajustar Saldo</h3>
          <div className="space-y-2">
            <Label>Minutos *</Label>
            <Input type="number" value={adjustForm.minutes} onChange={e => setAdjustForm({ ...adjustForm, minutes: e.target.value })} placeholder="Puede ser negativo" />
          </div>
          <div className="space-y-2">
            <Label>Motivo *</Label>
            <Input value={adjustForm.reason} onChange={e => setAdjustForm({ ...adjustForm, reason: e.target.value })} />
          </div>
          <Button onClick={adjust} disabled={!employeeId || !adjustForm.minutes || !adjustForm.reason}>Ajustar</Button>
        </CardContent>
      </Card>
    </div>
  )
}
