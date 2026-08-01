import { useEffect, useState } from 'react'
import { Plus, RefreshCw } from 'lucide-react'
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

interface CalendarEntry {
  id: string
  employee_id: string
  employee_name?: string
  work_date: string
  shift_id?: string
  shift_name?: string
  planned_start?: string
  planned_end?: string
  planned_break_minutes: number
  status: string
  source: string
}

interface ScheduleException {
  id: string
  employee_id?: string
  employee_name?: string
  exception_date: string
  exception_type: string
  start_time?: string
  end_time?: string
  reason?: string
  status: string
}

const exceptionTypes: Record<string, string> = {
  HOLIDAY: 'Feriado',
  SHIFT_CHANGE: 'Cambio de turno',
  TRAINING: 'Capacitación',
  MEETING: 'Reunión',
  ABSENCE: 'Ausencia',
  OTHER: 'Otro',
}

const calendarStatus: Record<string, string> = {
  PLANNED: 'Planificado', CONFIRMED: 'Confirmado', CANCELLED: 'Cancelado', COMPLETED: 'Completado',
}

export default function SchedulingPage() {
  const [activeTab, setActiveTab] = useState('calendar')

  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-900 mb-6">Turnos y Horarios</h1>

      <div className="flex gap-1 mb-6 border-b border-slate-200">
        {[{ key: 'calendar', label: 'Calendario' }, { key: 'exceptions', label: 'Excepciones' }, { key: 'swap', label: 'Intercambios' }].map((tab) => (
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

      {activeTab === 'calendar' && <CalendarTab />}
      {activeTab === 'exceptions' && <ExceptionsTab />}
      {activeTab === 'swap' && <SwapTab />}
    </div>
  )
}

function CalendarTab() {
  const [entries, setEntries] = useState<CalendarEntry[]>([])
  const [loading, setLoading] = useState(true)
  const [employees, setEmployees] = useState<{ value: string; label: string }[]>([])
  const [employeeId, setEmployeeId] = useState('')
  const [dateFrom, setDateFrom] = useState(new Date(Date.now() - 7 * 86400000).toISOString().slice(0, 10))
  const [dateTo, setDateTo] = useState(new Date(Date.now() + 23 * 86400000).toISOString().slice(0, 10))
  const [showGenerate, setShowGenerate] = useState(false)
  const [genForm, setGenForm] = useState({ employee_id: '', from: '', to: '' })
  const [generating, setGenerating] = useState(false)

  const fetch = async () => {
    setLoading(true)
    try {
      const params: Record<string, any> = { per_page: '100', date_from: dateFrom, date_to: dateTo }
      if (employeeId) params.employee_id = employeeId
      const res = await api.get('/scheduling/calendar', { params })
      setEntries(res.data.data ?? [])
    } catch { setEntries([]) }
    finally { setLoading(false) }
  }

  useEffect(() => {
    api.get('/employees', { params: { limit: '200' } })
      .then(res => setEmployees((res.data.data ?? []).map((e: any) => ({ value: e.id, label: `${e.first_name} ${e.last_name}` }))))
      .catch(() => {})
  }, [])

  useEffect(() => { fetch() }, [employeeId, dateFrom, dateTo])

  const openGenerate = async () => {
    if (employees.length === 0) {
      try {
        const res = await api.get('/employees', { params: { limit: '200' } })
        setEmployees((res.data.data ?? []).map((e: any) => ({ value: e.id, label: `${e.first_name} ${e.last_name}` })))
      } catch {}
    }
    setGenForm({ employee_id: '', from: dateFrom, to: dateTo })
    setShowGenerate(true)
  }

  const handleGenerate = async () => {
    setGenerating(true)
    try {
      const res = await api.post('/scheduling/calendar/generate', genForm)
      if (Array.isArray(res.data?.data)) {
        setEntries((prev) => [...res.data.data, ...prev])
      }
      setShowGenerate(false)
      fetch()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al generar calendario')
    } finally { setGenerating(false) }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4 flex-wrap gap-2">
        <div className="flex items-end gap-3 flex-wrap">
          <div>
            <Label>Empleado</Label>
            <Select options={employees} placeholder="Todos" value={employeeId} onChange={e => setEmployeeId(e.target.value)} />
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
        <Button size="sm" onClick={openGenerate}><Plus size={14} /> Generar Calendario</Button>
      </div>

      <Card>
        <CardContent className="p-0">
          {loading ? (
            <div className="p-6 text-center text-slate-500">Cargando...</div>
          ) : entries.length === 0 ? (
            <div className="p-6 text-center text-slate-500">No hay turnos en el rango. Usá "Generar Calendario" para crear turnos.</div>
          ) : (
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-slate-200 bg-slate-50">
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Fecha</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Empleado</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Turno</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Inicio</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Fin</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                </tr>
              </thead>
              <tbody>
                {entries.map((e) => (
                  <tr key={e.id} className="border-b border-slate-100 hover:bg-slate-50">
                    <td className="px-4 py-3 text-slate-900">{new Date(e.work_date).toLocaleDateString('es-AR')}</td>
                    <td className="px-4 py-3 font-medium text-slate-900">{e.employee_name || e.employee_id}</td>
                    <td className="px-4 py-3 text-slate-600">{e.shift_name || '-'}</td>
                    <td className="px-4 py-3 text-slate-600">{e.planned_start ? new Date(e.planned_start).toLocaleTimeString('es-AR', { hour: '2-digit', minute: '2-digit' }) : '-'}</td>
                    <td className="px-4 py-3 text-slate-600">{e.planned_end ? new Date(e.planned_end).toLocaleTimeString('es-AR', { hour: '2-digit', minute: '2-digit' }) : '-'}</td>
                    <td className="px-4 py-3">
                      <span className="inline-flex px-2 py-0.5 rounded-full text-xs font-medium bg-slate-100 text-slate-600">
                        {calendarStatus[e.status] || e.status}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </CardContent>
      </Card>

      <Dialog open={showGenerate} onOpenChange={setShowGenerate}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Generar Calendario</DialogTitle>
            <DialogDescription>Genera turnos para un empleado en un rango de fechas</DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label>Empleado *</Label>
              <Select
                options={employees}
                placeholder="Seleccionar..."
                value={genForm.employee_id}
                onChange={e => setGenForm({ ...genForm, employee_id: e.target.value })}
              />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Desde *</Label>
                <Input type="date" value={genForm.from} onChange={e => setGenForm({ ...genForm, from: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label>Hasta *</Label>
                <Input type="date" value={genForm.to} onChange={e => setGenForm({ ...genForm, to: e.target.value })} />
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowGenerate(false)}>Cancelar</Button>
            <Button onClick={handleGenerate} disabled={generating || !genForm.employee_id || !genForm.from || !genForm.to}>
              {generating ? 'Generando...' : 'Generar'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function ExceptionsTab() {
  const [exceptions, setExceptions] = useState<ScheduleException[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [employees, setEmployees] = useState<{ value: string; label: string }[]>([])
  const [form, setForm] = useState({ employee_id: '', exception_date: '', exception_type: 'HOLIDAY', start_time: '', end_time: '', reason: '' })
  const [saving, setSaving] = useState(false)

  const fetch = async () => {
    setLoading(true)
    try {
      const res = await api.get('/scheduling/exceptions')
      setExceptions(res.data.data ?? [])
    } catch { setExceptions([]) }
    finally { setLoading(false) }
  }

  useEffect(() => {
    api.get('/employees', { params: { limit: '200' } })
      .then(res => setEmployees((res.data.data ?? []).map((e: any) => ({ value: e.id, label: `${e.first_name} ${e.last_name}` }))))
      .catch(() => {})
    fetch()
  }, [])

  const handleCreate = async () => {
    setSaving(true)
    try {
      const payload: Record<string, any> = {
        exception_date: form.exception_date,
        exception_type: form.exception_type,
      }
      if (form.employee_id) payload.employee_id = form.employee_id
      if (form.start_time) payload.start_time = form.start_time
      if (form.end_time) payload.end_time = form.end_time
      if (form.reason) payload.reason = form.reason
      await api.post('/scheduling/exceptions', payload)
      setShowModal(false)
      setForm({ employee_id: '', exception_date: '', exception_type: 'HOLIDAY', start_time: '', end_time: '', reason: '' })
      fetch()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al crear excepción')
    } finally { setSaving(false) }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-slate-500">{exceptions.length} excepciones</p>
        <Button size="sm" onClick={() => setShowModal(true)}><Plus size={14} /> Nueva Excepción</Button>
      </div>

      <Card>
        <CardContent className="p-0">
          {loading ? (
            <div className="p-6 text-center text-slate-500">Cargando...</div>
          ) : exceptions.length === 0 ? (
            <div className="p-6 text-center text-slate-500">No hay excepciones</div>
          ) : (
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-slate-200 bg-slate-50">
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Fecha</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Empleado</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Tipo</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Horario</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Motivo</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                </tr>
              </thead>
              <tbody>
                {exceptions.map((x) => (
                  <tr key={x.id} className="border-b border-slate-100 hover:bg-slate-50">
                    <td className="px-4 py-3 text-slate-900">{new Date(x.exception_date).toLocaleDateString('es-AR')}</td>
                    <td className="px-4 py-3 font-medium text-slate-900">{x.employee_name || 'Todos'}</td>
                    <td className="px-4 py-3 text-slate-600">{exceptionTypes[x.exception_type] || x.exception_type}</td>
                    <td className="px-4 py-3 text-slate-600">{[x.start_time, x.end_time].filter(Boolean).join(' - ') || '-'}</td>
                    <td className="px-4 py-3 text-slate-600">{x.reason || '-'}</td>
                    <td className="px-4 py-3">
                      <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${x.status === 'APPROVED' ? 'bg-emerald-50 text-emerald-700' : 'bg-amber-50 text-amber-700'}`}>
                        {x.status}
                      </span>
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
            <DialogTitle>Nueva Excepción</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label>Empleado</Label>
              <Select options={employees} placeholder="Toda la empresa" value={form.employee_id} onChange={e => setForm({ ...form, employee_id: e.target.value })} />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Fecha *</Label>
                <Input type="date" value={form.exception_date} onChange={e => setForm({ ...form, exception_date: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label>Tipo *</Label>
                <Select
                  options={Object.entries(exceptionTypes).map(([value, label]) => ({ value, label }))}
                  value={form.exception_type}
                  onChange={e => setForm({ ...form, exception_type: e.target.value })}
                />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Desde</Label>
                <Input type="time" value={form.start_time} onChange={e => setForm({ ...form, start_time: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label>Hasta</Label>
                <Input type="time" value={form.end_time} onChange={e => setForm({ ...form, end_time: e.target.value })} />
              </div>
            </div>
            <div className="space-y-2">
              <Label>Motivo</Label>
              <Input value={form.reason} onChange={e => setForm({ ...form, reason: e.target.value })} />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button>
            <Button onClick={handleCreate} disabled={saving || !form.exception_date}>
              {saving ? 'Guardando...' : 'Crear'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function SwapTab() {
  const [employees, setEmployees] = useState<{ value: string; label: string }[]>([])
  const [form, setForm] = useState({ target_employee_id: '', requester_date: '', target_date: '', reason: '' })
  const [saving, setSaving] = useState(false)
  const [message, setMessage] = useState('')

  useEffect(() => {
    api.get('/employees', { params: { limit: '200' } })
      .then(res => setEmployees((res.data.data ?? []).map((e: any) => ({ value: e.id, label: `${e.first_name} ${e.last_name}` }))))
      .catch(() => {})
  }, [])

  const handleCreate = async () => {
    setSaving(true)
    try {
      const payload: Record<string, any> = {
        target_employee_id: form.target_employee_id,
        requester_date: form.requester_date,
        target_date: form.target_date,
      }
      if (form.reason) payload.reason = form.reason
      await api.post('/scheduling/swap', payload)
      setMessage('Intercambio solicitado correctamente.')
      setForm({ target_employee_id: '', requester_date: '', target_date: '', reason: '' })
    } catch (err: any) {
      setMessage('')
      alert(err?.response?.data?.error || 'Error al solicitar intercambio')
    } finally { setSaving(false) }
  }

  return (
    <Card className="max-w-xl">
      <CardContent className="p-6 space-y-4">
        <div>
          <h3 className="font-semibold text-slate-900">Solicitar Intercambio de Turno</h3>
          <p className="text-sm text-slate-500">Pedí a un compañero intercambiar tu turno por el suyo</p>
        </div>
        <div className="space-y-2">
          <Label>Compañero *</Label>
          <Select
            options={employees}
            placeholder="Seleccionar..."
            value={form.target_employee_id}
            onChange={e => setForm({ ...form, target_employee_id: e.target.value })}
          />
        </div>
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label>Tu turno *</Label>
            <Input type="date" value={form.requester_date} onChange={e => setForm({ ...form, requester_date: e.target.value })} />
          </div>
          <div className="space-y-2">
            <Label>Turno del compañero *</Label>
            <Input type="date" value={form.target_date} onChange={e => setForm({ ...form, target_date: e.target.value })} />
          </div>
        </div>
        <div className="space-y-2">
          <Label>Motivo</Label>
          <Input value={form.reason} onChange={e => setForm({ ...form, reason: e.target.value })} />
        </div>
        {message && <p className="text-sm text-emerald-600 font-medium">{message}</p>}
        <Button onClick={handleCreate} disabled={saving || !form.target_employee_id || !form.requester_date || !form.target_date}>
          <RefreshCw size={14} /> {saving ? 'Enviando...' : 'Solicitar Intercambio'}
        </Button>
      </CardContent>
    </Card>
  )
}
