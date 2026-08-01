import { useEffect, useState } from 'react'
import { LogIn, LogOut, Coffee, Clock, Check, X } from 'lucide-react'
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

interface AttendanceRecord {
  id: string
  employee_id: string
  employee_name?: string
  work_date: string
  scheduled_start?: string
  scheduled_end?: string
  actual_start?: string
  actual_end?: string
  scheduled_minutes: number
  worked_minutes: number
  late_minutes: number
  early_leave_minutes: number
  overtime_minutes: number
  break_minutes: number
  status: string
  notes?: string
}

interface AttendanceCorrection {
  id: string
  employee_id: string
  employee_name?: string
  attendance_id?: string
  correction_type: string
  requested_value: string
  original_value?: string
  reason: string
  status: string
  created_at: string
}

interface AttendanceDashboard {
  total_employees: number
  present: number
  absent: number
  late: number
  early_leave: number
  on_vacation: number
  on_leave: number
  holiday: number
  remote: number
  average_clock_in: string
  average_clock_out: string
  average_hours: number
}

const statusColors: Record<string, string> = {
  PRESENT: 'bg-emerald-50 text-emerald-700',
  LATE: 'bg-amber-50 text-amber-700',
  EARLY_LEAVE: 'bg-orange-50 text-orange-700',
  ABSENT: 'bg-red-50 text-red-700',
  INCOMPLETE: 'bg-slate-100 text-slate-600',
  VACATION: 'bg-blue-50 text-blue-700',
  LEAVE: 'bg-indigo-50 text-indigo-700',
  HOLIDAY: 'bg-purple-50 text-purple-700',
  REMOTE: 'bg-teal-50 text-teal-700',
}

const statusLabels: Record<string, string> = {
  PRESENT: 'Presente', LATE: 'Tardanza', EARLY_LEAVE: 'Salida temprana',
  ABSENT: 'Ausente', INCOMPLETE: 'Incompleto', VACATION: 'Vacaciones',
  LEAVE: 'Licencia', HOLIDAY: 'Feriado', REMOTE: 'Remoto',
}

const fmtTime = (d?: string) => (d ? new Date(d).toLocaleTimeString('es-AR', { hour: '2-digit', minute: '2-digit' }) : '-')
const fmtDate = (d: string) => new Date(d).toLocaleDateString('es-AR')

export default function AttendancePage() {
  const [activeTab, setActiveTab] = useState('records')

  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-900 mb-6">Asistencia</h1>

      <div className="flex gap-1 mb-6 border-b border-slate-200">
        {[{ key: 'records', label: 'Registros' }, { key: 'corrections', label: 'Correcciones' }].map((tab) => (
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

      <ClockInPanel />
      {activeTab === 'records' && <RecordsTab />}
      {activeTab === 'corrections' && <CorrectionsTab />}
    </div>
  )
}

function ClockInPanel() {
  const [myRecords, setMyRecords] = useState<AttendanceRecord[]>([])
  const [dashboard, setDashboard] = useState<AttendanceDashboard | null>(null)
  const [busy, setBusy] = useState<string | null>(null)

  const fetch = async () => {
    try {
      const today = new Date().toISOString().slice(0, 10)
      const res = await api.get('/attendance/me', { params: { from: today, to: today } })
      setMyRecords(res.data.data ?? [])
    } catch { setMyRecords([]) }
    try {
      const dRes = await api.get('/attendance/dashboard')
      setDashboard(dRes.data.data ?? null)
    } catch { setDashboard(null) }
  }

  useEffect(() => { fetch() }, [])

  const act = async (action: string, label: string) => {
    setBusy(action)
    try {
      await api.post(`/attendance/${action}`, { source: 'WEB' })
      fetch()
    } catch (err: any) {
      alert(err?.response?.data?.error || err?.response?.data?.message || `Error al ${label}`)
    } finally { setBusy(null) }
  }

  const today = myRecords[0] ?? null
  const isClockedIn = !!today?.actual_start && !today?.actual_end

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
      <Card>
        <CardContent className="p-5">
          <h3 className="text-sm font-semibold text-slate-500 mb-3">Mi jornada de hoy</h3>
          <div className="flex flex-wrap items-center gap-2">
            <Button size="sm" onClick={() => act('clock-in', 'fichar entrada')} disabled={!!busy || isClockedIn || !!today?.actual_end}>
              <LogIn size={14} /> Entrada
            </Button>
            <Button size="sm" variant="outline" onClick={() => act('break/start', 'iniciar pausa')} disabled={!!busy || !isClockedIn}>
              <Coffee size={14} /> Pausa
            </Button>
            <Button size="sm" variant="outline" onClick={() => act('break/end', 'finalizar pausa')} disabled={!!busy || !isClockedIn}>
              <Coffee size={14} /> Fin pausa
            </Button>
            <Button size="sm" onClick={() => act('clock-out', 'fichar salida')} disabled={!!busy || !isClockedIn}>
              <LogOut size={14} /> Salida
            </Button>
          </div>
          <div className="mt-3 text-sm space-y-1">
            <div className="flex justify-between">
              <span className="text-slate-500">Entrada</span>
              <span className="font-medium text-slate-900">{fmtTime(today?.actual_start)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-slate-500">Salida</span>
              <span className="font-medium text-slate-900">{fmtTime(today?.actual_end)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-slate-500">Estado</span>
              <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${statusColors[today?.status ?? ''] || 'bg-slate-100 text-slate-600'}`}>
                {statusLabels[today?.status ?? ''] || 'Sin fichar'}
              </span>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card className="md:col-span-2">
        <CardContent className="p-5">
          <h3 className="text-sm font-semibold text-slate-500 mb-3">Resumen del día</h3>
          {dashboard ? (
            <>
              <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-4">
                {[
                  { label: 'Total', value: dashboard.total_employees },
                  { label: 'Presentes', value: dashboard.present },
                  { label: 'Ausentes', value: dashboard.absent },
                  { label: 'Tardanzas', value: dashboard.late },
                ].map((s) => (
                  <div key={s.label} className="rounded-lg bg-slate-50 p-3 text-center">
                    <p className="text-xl font-bold text-slate-900">{s.value}</p>
                    <p className="text-xs text-slate-500">{s.label}</p>
                  </div>
                ))}
              </div>
              <div className="flex flex-wrap gap-4 text-sm text-slate-600">
                <span className="flex items-center gap-1"><Clock size={14} /> Prom. entrada: <strong>{dashboard.average_clock_in || '-'}</strong></span>
                <span className="flex items-center gap-1"><Clock size={14} /> Prom. salida: <strong>{dashboard.average_clock_out || '-'}</strong></span>
                <span className="flex items-center gap-1"><Clock size={14} /> Prom. horas: <strong>{dashboard.average_hours}</strong></span>
              </div>
            </>
          ) : (
            <p className="text-sm text-slate-500">Sin datos</p>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

function RecordsTab() {
  const [records, setRecords] = useState<AttendanceRecord[]>([])
  const [loading, setLoading] = useState(true)
  const [employees, setEmployees] = useState<{ value: string; label: string }[]>([])
  const [employeeId, setEmployeeId] = useState('')
  const [status, setStatus] = useState('')
  const [dateFrom, setDateFrom] = useState(new Date(Date.now() - 30 * 86400000).toISOString().slice(0, 10))
  const [dateTo, setDateTo] = useState(new Date().toISOString().slice(0, 10))

  const fetch = async () => {
    setLoading(true)
    try {
      const params: Record<string, any> = { limit: '100', date_from: dateFrom, date_to: dateTo }
      if (employeeId) params.employee_id = employeeId
      if (status) params.status = status
      const res = await api.get('/attendance', { params })
      setRecords(res.data.data ?? [])
    } catch { setRecords([]) }
    finally { setLoading(false) }
  }

  useEffect(() => {
    api.get('/employees', { params: { limit: '200' } })
      .then(res => setEmployees((res.data.data ?? []).map((e: any) => ({ value: e.id, label: `${e.first_name} ${e.last_name}` }))))
      .catch(() => {})
  }, [])

  useEffect(() => { fetch() }, [employeeId, status, dateFrom, dateTo])

  return (
    <div>
      <div className="grid grid-cols-1 sm:grid-cols-4 gap-3 mb-4">
        <div>
          <Label>Empleado</Label>
          <Select options={employees} placeholder="Todos" value={employeeId} onChange={e => setEmployeeId(e.target.value)} />
        </div>
        <div>
          <Label>Estado</Label>
          <Select
            options={Object.entries(statusLabels).map(([value, label]) => ({ value, label }))}
            placeholder="Todos"
            value={status}
            onChange={e => setStatus(e.target.value)}
          />
        </div>
        <div>
          <Label>Desde</Label>
          <Input type="date" value={dateFrom} onChange={e => setDateFrom(e.target.value)} />
        </div>
        <div>
          <Label>Hasta</Label>
          <Input type="date" value={dateTo} onChange={e => setDateTo(e.target.value)} />
        </div>
      </div>

      <Card>
        <CardContent className="p-0">
          {loading ? (
            <div className="p-6 text-center text-slate-500">Cargando...</div>
          ) : records.length === 0 ? (
            <div className="p-6 text-center text-slate-500">No hay registros</div>
          ) : (
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-slate-200 bg-slate-50">
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Fecha</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Empleado</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Entrada</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Salida</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Trabajado</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Tardanza</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                </tr>
              </thead>
              <tbody>
                {records.map((r) => (
                  <tr key={r.id} className="border-b border-slate-100 hover:bg-slate-50">
                    <td className="px-4 py-3 text-slate-900">{fmtDate(r.work_date)}</td>
                    <td className="px-4 py-3 font-medium text-slate-900">{r.employee_name || '-'}</td>
                    <td className="px-4 py-3 text-slate-600">{fmtTime(r.actual_start)}</td>
                    <td className="px-4 py-3 text-slate-600">{fmtTime(r.actual_end)}</td>
                    <td className="px-4 py-3 text-slate-600">{Math.round(r.worked_minutes / 60)}h {r.worked_minutes % 60}m</td>
                    <td className="px-4 py-3 text-slate-600">{r.late_minutes > 0 ? `${r.late_minutes}m` : '-'}</td>
                    <td className="px-4 py-3">
                      <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${statusColors[r.status] || 'bg-slate-100 text-slate-600'}`}>
                        {statusLabels[r.status] || r.status}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

function CorrectionsTab() {
  const [corrections, setCorrections] = useState<AttendanceCorrection[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [form, setForm] = useState({ correction_type: 'CLOCK_IN', requested_value: '', reason: '' })
  const [saving, setSaving] = useState(false)
  const [actionLoading, setActionLoading] = useState<string | null>(null)

  const fetch = async () => {
    setLoading(true)
    try {
      const res = await api.get('/attendance/corrections', { params: { limit: '100' } })
      setCorrections(res.data.data ?? [])
    } catch { setCorrections([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetch() }, [])

  const handleCreate = async () => {
    setSaving(true)
    try {
      await api.post('/attendance/corrections', {
        correction_type: form.correction_type,
        requested_value: new Date(form.requested_value).toISOString(),
        reason: form.reason,
      })
      setShowModal(false)
      setForm({ correction_type: 'CLOCK_IN', requested_value: '', reason: '' })
      fetch()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al crear corrección')
    } finally { setSaving(false) }
  }

  const handleAction = async (id: string, action: 'approve' | 'reject') => {
    setActionLoading(id)
    try {
      await api.post(`/attendance/corrections/${id}/${action}`)
      fetch()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al procesar corrección')
    } finally { setActionLoading(null) }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-slate-500">{corrections.length} correcciones</p>
        <Button size="sm" onClick={() => setShowModal(true)}>Nueva Corrección</Button>
      </div>

      <Card>
        <CardContent className="p-0">
          {loading ? (
            <div className="p-6 text-center text-slate-500">Cargando...</div>
          ) : corrections.length === 0 ? (
            <div className="p-6 text-center text-slate-500">No hay correcciones</div>
          ) : (
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-slate-200 bg-slate-50">
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Empleado</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Tipo</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Valor solicitado</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Motivo</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                  <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                </tr>
              </thead>
              <tbody>
                {corrections.map((c) => (
                  <tr key={c.id} className="border-b border-slate-100 hover:bg-slate-50">
                    <td className="px-4 py-3 font-medium text-slate-900">{c.employee_name || c.employee_id}</td>
                    <td className="px-4 py-3 text-slate-600">{c.correction_type === 'CLOCK_IN' ? 'Entrada' : 'Salida'}</td>
                    <td className="px-4 py-3 text-slate-600">{new Date(c.requested_value).toLocaleString('es-AR')}</td>
                    <td className="px-4 py-3 text-slate-600">{c.reason}</td>
                    <td className="px-4 py-3">
                      <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${c.status === 'PENDING' ? 'bg-amber-50 text-amber-700' : c.status === 'APPROVED' ? 'bg-emerald-50 text-emerald-700' : 'bg-red-50 text-red-700'}`}>
                        {c.status}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-right whitespace-nowrap">
                      {c.status === 'PENDING' && (
                        <>
                          <Button variant="ghost" size="sm" className="text-emerald-600" disabled={actionLoading === c.id} onClick={() => handleAction(c.id, 'approve')}>
                            <Check size={14} />
                          </Button>
                          <Button variant="ghost" size="sm" className="text-red-500" disabled={actionLoading === c.id} onClick={() => handleAction(c.id, 'reject')}>
                            <X size={14} />
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

      <Dialog open={showModal} onOpenChange={setShowModal}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Nueva Corrección</DialogTitle>
            <DialogDescription>Solicitá ajustar tu horario de entrada o salida</DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label>Tipo *</Label>
              <Select
                options={[{ value: 'CLOCK_IN', label: 'Entrada' }, { value: 'CLOCK_OUT', label: 'Salida' }]}
                value={form.correction_type}
                onChange={e => setForm({ ...form, correction_type: e.target.value })}
              />
            </div>
            <div className="space-y-2">
              <Label>Nuevo horario *</Label>
              <Input
                type="datetime-local"
                value={form.requested_value}
                onChange={e => setForm({ ...form, requested_value: e.target.value })}
              />
            </div>
            <div className="space-y-2">
              <Label>Motivo *</Label>
              <textarea
                className="flex min-h-[60px] w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm"
                value={form.reason}
                onChange={e => setForm({ ...form, reason: e.target.value })}
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button>
            <Button onClick={handleCreate} disabled={saving || !form.requested_value || !form.reason}>
              {saving ? 'Guardando...' : 'Solicitar'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
