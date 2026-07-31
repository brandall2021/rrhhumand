import { useEffect, useState } from 'react'
import { Check, X, Plus, Ban } from 'lucide-react'
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
import { useAuth } from '@/contexts/auth-context'

interface LeaveType {
  id: string
  name: string
  code: string
  category: string
  requires_approval: boolean
  requires_document: boolean
  affects_balance: boolean
  is_paid: boolean
  is_active: boolean
  description?: string
}

interface LeaveBalance {
  leave_type_id: string
  leave_type_name: string
  year: number
  allocated_days: number
  carried_over_days: number
  adjustment_days: number
  used_days: number
  reserved_days: number
  available_days: number
}

interface LeaveRequest {
  id: string
  employee_id: string
  employee_name?: string
  leave_type_id: string
  leave_type_name?: string
  start_date: string
  end_date: string
  requested_days: number
  reason?: string
  status: string
  created_at: string
  approvals?: LeaveApproval[]
}

interface LeaveApproval {
  id: string
  approver_name?: string
  level: number
  status: string
  comments?: string
  decided_at?: string
}

interface Holiday {
  id: string
  date: string
  name: string
  is_recurring: boolean
}

const statusColors: Record<string, string> = {
  PENDING: 'bg-amber-50 text-amber-700',
  APPROVED: 'bg-emerald-50 text-emerald-700',
  REJECTED: 'bg-red-50 text-red-700',
  CANCELLED: 'bg-slate-100 text-slate-600',
}

const statusLabels: Record<string, string> = {
  PENDING: 'Pendiente',
  APPROVED: 'Aprobado',
  REJECTED: 'Rechazado',
  CANCELLED: 'Cancelado',
}

export default function LeavesPage() {
  const { user } = useAuth()
  const [activeTab, setActiveTab] = useState('requests')

  const tabs = [
    { key: 'requests', label: 'Solicitudes' },
    { key: 'balance', label: 'Saldo' },
    { key: 'types', label: 'Tipos' },
    { key: 'holidays', label: 'Feriados' },
  ]

  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-900 mb-6">Vacaciones / Licencias</h1>

      <div className="flex gap-1 mb-6 border-b border-slate-200">
        {tabs.map((tab) => (
          <button
            key={tab.key}
            onClick={() => setActiveTab(tab.key)}
            className={`px-4 py-2.5 text-sm font-medium border-b-2 transition-colors ${
              activeTab === tab.key
                ? 'border-brand-600 text-brand-700'
                : 'border-transparent text-slate-500 hover:text-slate-700'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {activeTab === 'requests' && <RequestsTab userId={user?.id} />}
      {activeTab === 'balance' && <BalanceTab userId={user?.id} />}
      {activeTab === 'types' && <LeaveTypesTab />}
      {activeTab === 'holidays' && <HolidaysTab />}
    </div>
  )
}

function RequestsTab({ userId }: { userId?: string }) {
  const [requests, setRequests] = useState<LeaveRequest[]>([])
  const [loading, setLoading] = useState(true)
  const [showCreate, setShowCreate] = useState(false)
  const [leaveTypes, setLeaveTypes] = useState<LeaveType[]>([])
  const [employees, setEmployees] = useState<{ value: string; label: string }[]>([])
  const [form, setForm] = useState({
    employee_id: '',
    leave_type_id: '',
    start_date: '',
    end_date: '',
    reason: '',
  })
  const [saving, setSaving] = useState(false)
  const [actionLoading, setActionLoading] = useState<string | null>(null)

  const fetch = async () => {
    setLoading(true)
    try {
      const res = await api.get('/leave/requests', { params: { limit: '100' } })
      setRequests(res.data.data ?? [])
    } catch { setRequests([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetch() }, [])

  const openCreate = async () => {
    try {
      const [ltRes, empRes] = await Promise.all([
        api.get('/leave/types', { params: { limit: '100' } }),
        api.get('/employees', { params: { limit: '200' } }),
      ])
      setLeaveTypes(ltRes.data.data ?? [])
      setEmployees((empRes.data.data ?? []).map((e: any) => ({ value: e.id, label: `${e.first_name} ${e.last_name}` })))
    } catch {}
    setForm({ employee_id: userId ?? '', leave_type_id: '', start_date: '', end_date: '', reason: '' })
    setShowCreate(true)
  }

  const handleCreate = async () => {
    setSaving(true)
    try {
      await api.post('/leave/requests', form)
      setShowCreate(false)
      fetch()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al crear solicitud')
    } finally { setSaving(false) }
  }

  const handleAction = async (id: string, action: 'approve' | 'reject' | 'cancel') => {
    setActionLoading(id)
    try {
      await api.post(`/leave/requests/${id}/${action}`)
      fetch()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al procesar solicitud')
    } finally { setActionLoading(null) }
  }

  const totalDays = form.start_date && form.end_date
    ? Math.max(1, Math.ceil((new Date(form.end_date).getTime() - new Date(form.start_date).getTime()) / (1000 * 60 * 60 * 24)) + 1)
    : 0

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-slate-500">{requests.length} solicitudes</p>
        <Button size="sm" onClick={openCreate}><Plus size={14} /> Nueva Solicitud</Button>
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
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Empleado</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Tipo</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Desde</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Hasta</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Días</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                  <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                </tr>
              </thead>
              <tbody>
                {requests.map((r) => (
                  <tr key={r.id} className="border-b border-slate-100 hover:bg-slate-50">
                    <td className="px-4 py-3 font-medium text-slate-900">{r.employee_name || r.employee_id}</td>
                    <td className="px-4 py-3 text-slate-600">{r.leave_type_name || '-'}</td>
                    <td className="px-4 py-3 text-slate-600">{new Date(r.start_date).toLocaleDateString('es-AR')}</td>
                    <td className="px-4 py-3 text-slate-600">{new Date(r.end_date).toLocaleDateString('es-AR')}</td>
                    <td className="px-4 py-3 text-slate-600">{r.requested_days}</td>
                    <td className="px-4 py-3">
                      <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${statusColors[r.status] || 'bg-slate-100 text-slate-600'}`}>
                        {statusLabels[r.status] || r.status}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-right whitespace-nowrap">
                      {r.status === 'PENDING' && (
                        <>
                          <Button
                            variant="ghost" size="sm"
                            className="text-emerald-600"
                            onClick={() => handleAction(r.id, 'approve')}
                            disabled={actionLoading === r.id}
                          >
                            <Check size={14} />
                          </Button>
                          <Button
                            variant="ghost" size="sm"
                            className="text-red-500"
                            onClick={() => handleAction(r.id, 'reject')}
                            disabled={actionLoading === r.id}
                          >
                            <X size={14} />
                          </Button>
                          <Button
                            variant="ghost" size="sm"
                            className="text-slate-400"
                            onClick={() => handleAction(r.id, 'cancel')}
                            disabled={actionLoading === r.id}
                          >
                            <Ban size={14} />
                          </Button>
                        </>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </CardContent>
      </Card>

      <Dialog open={showCreate} onOpenChange={setShowCreate}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Nueva Solicitud de Licencia</DialogTitle>
            <DialogDescription>Completá los datos para solicitar una licencia o vacación</DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label>Empleado *</Label>
              <Select
                options={employees}
                placeholder="Seleccionar..."
                value={form.employee_id}
                onChange={e => setForm({ ...form, employee_id: e.target.value })}
              />
            </div>
            <div className="space-y-2">
              <Label>Tipo de Licencia *</Label>
              <Select
                options={leaveTypes.filter(t => t.is_active).map(t => ({ value: t.id, label: t.name }))}
                placeholder="Seleccionar..."
                value={form.leave_type_id}
                onChange={e => setForm({ ...form, leave_type_id: e.target.value })}
              />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Fecha Inicio *</Label>
                <Input type="date" value={form.start_date} onChange={e => setForm({ ...form, start_date: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label>Fecha Fin *</Label>
                <Input type="date" value={form.end_date} onChange={e => setForm({ ...form, end_date: e.target.value })} />
              </div>
            </div>
            {totalDays > 0 && (
              <p className="text-sm text-slate-600">Días solicitados: <strong>{totalDays}</strong></p>
            )}
            <div className="space-y-2">
              <Label>Motivo</Label>
              <textarea
                className="flex min-h-[80px] w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm"
                value={form.reason}
                onChange={e => setForm({ ...form, reason: e.target.value })}
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowCreate(false)}>Cancelar</Button>
            <Button onClick={handleCreate} disabled={saving || !form.employee_id || !form.leave_type_id || !form.start_date || !form.end_date}>
              {saving ? 'Guardando...' : 'Solicitar'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function BalanceTab({ userId }: { userId?: string }) {
  const [balances, setBalances] = useState<LeaveBalance[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedEmployee, setSelectedEmployee] = useState(userId || '')
  const [employees, setEmployees] = useState<{ value: string; label: string }[]>([])

  const fetch = async () => {
    if (!selectedEmployee) { setBalances([]); setLoading(false); return }
    setLoading(true)
    try {
      const res = await api.get('/leave/balance', { params: { employee_id: selectedEmployee } })
      setBalances(res.data.data?.balances ?? res.data.data ?? [])
    } catch { setBalances([]) }
    finally { setLoading(false) }
  }

  useEffect(() => {
    api.get('/employees', { params: { limit: '200' } })
      .then(res => setEmployees((res.data.data ?? []).map((e: any) => ({ value: e.id, label: `${e.first_name} ${e.last_name}` }))))
      .catch(() => {})
  }, [])

  useEffect(() => { fetch() }, [selectedEmployee])

  return (
    <div>
      <div className="mb-4 max-w-xs">
        <Label>Empleado</Label>
        <Select
          options={employees}
          placeholder="Seleccionar..."
          value={selectedEmployee}
          onChange={e => setSelectedEmployee(e.target.value)}
        />
      </div>

      {loading ? (
        <div className="text-center text-slate-500 py-6">Cargando...</div>
      ) : balances.length === 0 ? (
        <div className="text-center text-slate-500 py-6">Sin datos de saldo</div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {balances.map((b) => {
            const pct = b.allocated_days > 0 ? Math.round((b.used_days / b.allocated_days) * 100) : 0
            return (
              <Card key={b.leave_type_id}>
                <CardContent className="p-5">
                  <h3 className="font-semibold text-slate-900 mb-3">{b.leave_type_name}</h3>
                  <div className="space-y-2 text-sm">
                    <div className="flex justify-between">
                      <span className="text-slate-500">Asignados</span>
                      <span className="font-medium">{b.allocated_days}</span>
                    </div>
                    {b.carried_over_days > 0 && (
                      <div className="flex justify-between">
                        <span className="text-slate-500">Del año anterior</span>
                        <span className="font-medium text-amber-600">+{b.carried_over_days}</span>
                      </div>
                    )}
                    {b.adjustment_days !== 0 && (
                      <div className="flex justify-between">
                        <span className="text-slate-500">Ajustes</span>
                        <span className="font-medium text-blue-600">{b.adjustment_days > 0 ? '+' : ''}{b.adjustment_days}</span>
                      </div>
                    )}
                    <div className="flex justify-between">
                      <span className="text-slate-500">Usados</span>
                      <span className="font-medium text-red-500">-{b.used_days}</span>
                    </div>
                    {b.reserved_days > 0 && (
                      <div className="flex justify-between">
                        <span className="text-slate-500">Reservados</span>
                        <span className="font-medium text-amber-500">-{b.reserved_days}</span>
                      </div>
                    )}
                    <div className="border-t pt-2 mt-2">
                      <div className="flex justify-between text-base">
                        <span className="font-semibold text-slate-700">Disponibles</span>
                        <span className="font-bold text-brand-700">{b.available_days}</span>
                      </div>
                    </div>
                    <div className="mt-3">
                      <div className="h-2 bg-slate-100 rounded-full overflow-hidden">
                        <div
                          className="h-full bg-brand-500 rounded-full transition-all"
                          style={{ width: `${Math.min(pct, 100)}%` }}
                        />
                      </div>
                      <p className="text-xs text-slate-400 mt-1">{pct}% utilizado</p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            )
          })}
        </div>
      )}
    </div>
  )
}

function LeaveTypesTab() {
  const [types, setTypes] = useState<LeaveType[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [editing, setEditing] = useState<LeaveType | null>(null)
  const [form, setForm] = useState({ name: '', code: '', category: 'vacation', description: '', requires_approval: true, requires_document: false, affects_balance: true, is_paid: true })
  const [saving, setSaving] = useState(false)

  const fetch = async () => {
    setLoading(true)
    try {
      const res = await api.get('/leave/types', { params: { limit: '100' } })
      setTypes(res.data.data ?? [])
    } catch { setTypes([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetch() }, [])

  const openCreate = () => {
    setEditing(null)
    setForm({ name: '', code: '', category: 'vacation', description: '', requires_approval: true, requires_document: false, affects_balance: true, is_paid: true })
    setShowModal(true)
  }

  const openEdit = (t: LeaveType) => {
    setEditing(t)
    setForm({
      name: t.name,
      code: t.code,
      category: t.category,
      description: t.description ?? '',
      requires_approval: t.requires_approval,
      requires_document: t.requires_document,
      affects_balance: t.affects_balance,
      is_paid: t.is_paid,
    })
    setShowModal(true)
  }

  const handleSave = async () => {
    setSaving(true)
    try {
      if (editing) {
        await api.put(`/leave/types/${editing.id}`, form)
      } else {
        await api.post('/leave/types', form)
      }
      setShowModal(false)
      fetch()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar tipo')
    } finally { setSaving(false) }
  }

  const categories: Record<string, string> = { vacation: 'Vacación', sick: 'Enfermedad', personal: 'Personal', family: 'Familiar', other: 'Otro' }

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-slate-500">{types.length} tipos</p>
        <Button size="sm" onClick={openCreate}><Plus size={14} /> Nuevo Tipo</Button>
      </div>

      <Card>
        <CardContent className="p-0">
          {loading ? (
            <div className="p-6 text-center text-slate-500">Cargando...</div>
          ) : types.length === 0 ? (
            <div className="p-6 text-center text-slate-500">No hay tipos de licencia</div>
          ) : (
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-slate-200 bg-slate-50">
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Código</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Categoría</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Activo</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Pagado</th>
                  <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                </tr>
              </thead>
              <tbody>
                {types.map((t) => (
                  <tr key={t.id} className="border-b border-slate-100 hover:bg-slate-50">
                    <td className="px-4 py-3 font-medium text-slate-900">{t.name}</td>
                    <td className="px-4 py-3 text-slate-600">{t.code}</td>
                    <td className="px-4 py-3 text-slate-600">{categories[t.category] || t.category}</td>
                    <td className="px-4 py-3">{t.is_active ? <span className="text-emerald-600 font-medium">Sí</span> : 'No'}</td>
                    <td className="px-4 py-3">{t.is_paid ? <span className="text-emerald-600 font-medium">Sí</span> : 'No'}</td>
                    <td className="px-4 py-3 text-right">
                      <Button variant="ghost" size="sm" onClick={() => openEdit(t)}>Editar</Button>
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
            <DialogTitle>{editing ? 'Editar Tipo' : 'Nuevo Tipo de Licencia'}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Nombre *</Label>
                <Input value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label>Código *</Label>
                <Input value={form.code} onChange={e => setForm({ ...form, code: e.target.value })} />
              </div>
            </div>
            <div className="space-y-2">
              <Label>Categoría</Label>
              <select
                className="flex h-10 w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm"
                value={form.category}
                onChange={e => setForm({ ...form, category: e.target.value })}
              >
                <option value="vacation">Vacación</option>
                <option value="sick">Enfermedad</option>
                <option value="personal">Personal</option>
                <option value="family">Familiar</option>
                <option value="other">Otro</option>
              </select>
            </div>
            <div className="space-y-2">
              <Label>Descripción</Label>
              <textarea
                className="flex min-h-[60px] w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm"
                value={form.description}
                onChange={e => setForm({ ...form, description: e.target.value })}
              />
            </div>
            <div className="grid grid-cols-2 gap-4 text-sm">
              <label className="flex items-center gap-2">
                <input type="checkbox" checked={form.requires_approval} onChange={e => setForm({ ...form, requires_approval: e.target.checked })} className="rounded" />
                Requiere aprobación
              </label>
              <label className="flex items-center gap-2">
                <input type="checkbox" checked={form.requires_document} onChange={e => setForm({ ...form, requires_document: e.target.checked })} className="rounded" />
                Requiere documento
              </label>
              <label className="flex items-center gap-2">
                <input type="checkbox" checked={form.affects_balance} onChange={e => setForm({ ...form, affects_balance: e.target.checked })} className="rounded" />
                Afecta saldo
              </label>
              <label className="flex items-center gap-2">
                <input type="checkbox" checked={form.is_paid} onChange={e => setForm({ ...form, is_paid: e.target.checked })} className="rounded" />
                Pagado
              </label>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button>
            <Button onClick={handleSave} disabled={saving || !form.name || !form.code}>
              {saving ? 'Guardando...' : editing ? 'Guardar Cambios' : 'Crear'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function HolidaysTab() {
  const [holidays, setHolidays] = useState<Holiday[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [form, setForm] = useState({ date: '', name: '', is_recurring: false })
  const [saving, setSaving] = useState(false)

  const fetch = async () => {
    setLoading(true)
    try {
      const res = await api.get('/leave/holidays')
      setHolidays(res.data.data ?? [])
    } catch { setHolidays([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetch() }, [])

  const handleCreate = async () => {
    setSaving(true)
    try {
      await api.post('/leave/holidays', form)
      setShowModal(false)
      setForm({ date: '', name: '', is_recurring: false })
      fetch()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al crear feriado')
    } finally { setSaving(false) }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('¿Eliminar este feriado?')) return
    try {
      await api.delete(`/leave/holidays/${id}`)
      fetch()
    } catch { alert('Error al eliminar') }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-slate-500">{holidays.length} feriados</p>
        <Button size="sm" onClick={() => setShowModal(true)}><Plus size={14} /> Agregar Feriado</Button>
      </div>

      <Card>
        <CardContent className="p-0">
          {loading ? (
            <div className="p-6 text-center text-slate-500">Cargando...</div>
          ) : holidays.length === 0 ? (
            <div className="p-6 text-center text-slate-500">No hay feriados registrados</div>
          ) : (
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-slate-200 bg-slate-50">
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Fecha</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Recurrente</th>
                  <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                </tr>
              </thead>
              <tbody>
                {holidays.map((h) => (
                  <tr key={h.id} className="border-b border-slate-100 hover:bg-slate-50">
                    <td className="px-4 py-3 text-slate-900">{new Date(h.date).toLocaleDateString('es-AR')}</td>
                    <td className="px-4 py-3 font-medium text-slate-900">{h.name}</td>
                    <td className="px-4 py-3">{h.is_recurring ? <span className="text-emerald-600 font-medium">Sí</span> : 'No'}</td>
                    <td className="px-4 py-3 text-right">
                      <Button variant="ghost" size="sm" className="text-red-500" onClick={() => handleDelete(h.id)}>Eliminar</Button>
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
          <DialogHeader><DialogTitle>Nuevo Feriado</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label>Fecha *</Label>
              <Input type="date" value={form.date} onChange={e => setForm({ ...form, date: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label>Nombre *</Label>
              <Input value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} />
            </div>
            <label className="flex items-center gap-2 text-sm">
              <input type="checkbox" checked={form.is_recurring} onChange={e => setForm({ ...form, is_recurring: e.target.checked })} className="rounded" />
              Recurrente (se repite cada año)
            </label>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button>
            <Button onClick={handleCreate} disabled={saving || !form.date || !form.name}>
              {saving ? 'Guardando...' : 'Crear'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
