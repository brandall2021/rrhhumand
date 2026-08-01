import { useEffect, useState } from 'react'
import { Plus, Play, CheckCircle, XCircle, AlertTriangle } from 'lucide-react'
import api from '@/lib/api'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select } from '@/components/ui/select'
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter,
} from '@/components/ui/dialog'

const fmt = (d?: string) => (d ? new Date(d).toLocaleDateString('es-AR') : '-')
const pill = (s: string) => (
  <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${s === 'COMPLETED' || s === 'APPROVED' || s === 'APPROVAL' || s === 'REVOKED' || s === 'RETURNED' ? 'bg-emerald-50 text-emerald-700' : s === 'PENDING' || s === 'IN_PROGRESS' || s === 'PENDING_APPROVAL' || s === 'PENDING_RETURN' ? 'bg-amber-50 text-amber-700' : s === 'CANCELLED' || s === 'BLOCKED' || s === 'DAMAGED' || s === 'LOST' ? 'bg-rose-50 text-rose-700' : s === 'DRAFT' ? 'bg-slate-100 text-slate-600' : 'bg-sky-50 text-sky-700'}`}>{s}</span>
)

export default function OnboardingPage() {
  const [mode, setMode] = useState<'onboarding' | 'offboarding'>('onboarding')
  const [dash, setDash] = useState<any>(null)

  useEffect(() => {
    if (mode === 'onboarding') {
      api.get('/onboarding/dashboard').then(r => setDash(r.data.data)).catch(() => setDash(null))
    } else {
      api.get('/offboarding/dashboard').then(r => setDash(r.data.data)).catch(() => setDash(null))
    }
  }, [mode])

  const dashCards = mode === 'onboarding'
    ? [['Activos', dash?.active_onboardings], ['Pendientes', dash?.pending_onboardings], ['Completados', dash?.completed_onboardings], ['Atrasados', dash?.overdue_onboardings], ['Progreso prom.', dash?.average_progress?.toFixed?.(1) ?? dash?.average_progress], ['Docs pend.', dash?.pending_documents], ['Assets pend.', dash?.pending_assets], ['Accesos pend.', dash?.pending_access]]
    : [['Activos', dash?.active_offboardings], ['Pendientes', dash?.pending_offboardings], ['Completados', dash?.completed_offboardings], ['Atrasados', dash?.overdue_offboardings], ['Assets no devueltos', dash?.assets_not_returned], ['Accesos sin revocar', dash?.access_not_revoked]]

  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-900 mb-6">Onboarding y Offboarding</h1>
      {dash && (
        <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-3 mb-6">
          {dashCards.map(([l, v]) => (
            <Card key={l as string}><CardContent className="p-4"><p className="text-xs text-slate-500 mb-1">{l}</p><p className="text-lg font-bold text-slate-900">{v ?? 0}</p></CardContent></Card>
          ))}
        </div>
      )}
      <div className="flex gap-1 mb-6 border-b border-slate-200">
        {([['onboarding', 'Onboarding'], ['offboarding', 'Offboarding']] as const).map(([k, l]) => (
          <button key={k} onClick={() => setMode(k)} className={`px-4 py-2.5 text-sm font-medium border-b-2 transition-colors ${mode === k ? 'border-brand-600 text-brand-700' : 'border-transparent text-slate-500 hover:text-slate-700'}`}>{l}</button>
        ))}
      </div>
      {mode === 'onboarding' ? <OnboardingTab /> : <OffboardingTab />}
    </div>
  )
}

function useEmployees() {
  const [employees, setEmployees] = useState<{ value: string; label: string }[]>([])
  useEffect(() => {
    api.get('/employees', { params: { limit: '200' } }).then(r => setEmployees((r.data.data ?? []).map((e: any) => ({ value: e.id, label: `${e.first_name} ${e.last_name}` })))).catch(() => {})
  }, [])
  return employees
}

function OnboardingTab() {
  const [list, setList] = useState<any[]>([])
  const [detail, setDetail] = useState<any>(null)
  const [tasks, setTasks] = useState<any[]>([])
  const [docs, setDocs] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [statusFilter, setStatusFilter] = useState('')
  const [search, setSearch] = useState('')
  const [showModal, setShowModal] = useState(false)
  const [saving, setSaving] = useState(false)
  const [showTasks, setShowTasks] = useState(false)
  const [showDocs, setShowDocs] = useState(false)
  const [taskForm, setTaskForm] = useState({ title: '', task_type: 'OTHER', assigned_role: 'EMPLOYEE', due_date: '' })
  const [docForm, setDocForm] = useState({ document_type: '', name: '' })
  const employees = useEmployees()
  const [form, setForm] = useState({ employee_id: '', start_date: '', employee_type: 'ADMINISTRATIVO', work_mode: 'HIBRIDO', template_id: '' })

  const fetchList = async () => {
    setLoading(true)
    try {
      const params: Record<string, any> = {}
      if (statusFilter) params.status = statusFilter
      if (search) params.search = search
      const r = await api.get('/onboarding', { params })
      setList(r.data.data ?? [])
    } catch { setList([]) } finally { setLoading(false) }
  }
  useEffect(() => { fetchList() }, [statusFilter])

  const openDetail = async (id: string) => {
    try {
      const [d, t, dc] = await Promise.all([
        api.get(`/onboarding/${id}`),
        api.get(`/onboarding/${id}/tasks`),
        api.get(`/onboarding/${id}/documents`),
      ])
      setDetail(d.data.data); setTasks(t.data.data ?? []); setDocs(dc.data.data ?? [])
    } catch { alert('No se pudo cargar el detalle') }
  }

  const create = async () => {
    setSaving(true)
    try {
      const payload: Record<string, any> = { employee_id: form.employee_id, start_date: new Date(form.start_date).toISOString() }
      if (form.employee_type) payload.employee_type = form.employee_type
      if (form.work_mode) payload.work_mode = form.work_mode
      if (form.template_id) payload.template_id = form.template_id
      await api.post('/onboarding', payload)
      setShowModal(false); setForm({ employee_id: '', start_date: '', employee_type: 'ADMINISTRATIVO', work_mode: 'HIBRIDO', template_id: '' }); fetchList()
    } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }

  const act = async (id: string, action: string, extra?: Record<string, any>) => {
    try {
      if (action === 'cancel') {
        const reason = window.prompt('Motivo de cancelación:')
        if (reason === null) return
        await api.post(`/onboarding/${id}/cancel`, { reason })
      } else {
        await api.post(`/onboarding/${id}/${action}`, extra ?? {})
      }
      openDetail(id); fetchList()
    } catch (e: any) { alert(e?.response?.data?.error || 'Error') }
  }

  const addTask = async () => {
    setSaving(true)
    try {
      const payload: Record<string, any> = { title: taskForm.title, task_type: taskForm.task_type, assigned_role: taskForm.assigned_role }
      if (taskForm.due_date) payload.due_date = new Date(taskForm.due_date).toISOString()
      await api.post(`/onboarding/${detail.id}/tasks`, payload)
      setShowTasks(false); setTaskForm({ title: '', task_type: 'OTHER', assigned_role: 'EMPLOYEE', due_date: '' }); openDetail(detail.id)
    } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }

  const completeTask = async (taskId: string) => {
    try { await api.post(`/onboarding/tasks/${taskId}/complete`); openDetail(detail.id) } catch (e: any) { alert(e?.response?.data?.error || 'Error') }
  }

  const addDoc = async () => {
    setSaving(true)
    try {
      await api.post(`/onboarding/${detail.id}/documents`, docForm)
      setShowDocs(false); setDocForm({ document_type: '', name: '' }); openDetail(detail.id)
    } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }

  const docAction = async (docId: string, action: 'approve' | 'reject') => {
    try { await api.post(`/onboarding/documents/${docId}/${action}`); openDetail(detail.id) } catch (e: any) { alert(e?.response?.data?.error || 'Error') }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4 flex-wrap gap-2">
        <div className="flex gap-2">
          <div className="w-48"><Select options={[{ value: '', label: 'Todos los estados' }, ...['DRAFT', 'PENDING', 'IN_PROGRESS', 'BLOCKED', 'COMPLETED', 'CANCELLED'].map(s => ({ value: s, label: s }))]} value={statusFilter} onChange={e => setStatusFilter(e.target.value)} /></div>
          <div className="w-56"><Input placeholder="Buscar..." value={search} onChange={e => setSearch(e.target.value)} onKeyDown={e => e.key === 'Enter' && fetchList()} /></div>
        </div>
        <Button size="sm" onClick={() => setShowModal(true)}><Plus size={14} /> Nuevo Onboarding</Button>
      </div>
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {loading ? <div className="p-6 text-center text-slate-500 col-span-full">Cargando...</div> : list.length === 0 ? <div className="p-6 text-center text-slate-500 col-span-full">No hay procesos</div> : list.map((p: any) => (
          <Card key={p.id}>
            <CardContent className="p-5">
              <div className="flex items-start justify-between mb-2">
                <div>
                  <p className="font-semibold text-slate-900">{p.employee_id}</p>
                  <p className="text-xs text-slate-500">{fmt(p.start_date)} · {p.employee_type || '-'} · {p.work_mode || '-'}</p>
                </div>
                {pill(p.status)}
              </div>
              <div className="flex items-center gap-2 mb-3">
                <div className="flex-1 h-1.5 bg-slate-100 rounded-full overflow-hidden"><div className="h-full rounded-full" style={{ width: `${Math.round(p.progress ?? 0)}%`, backgroundColor: (p.progress ?? 0) >= 100 ? '#10b981' : (p.progress ?? 0) > 0 ? '#f59e0b' : '#6366f1' }} /></div>
                <span className="text-xs text-slate-500">{Math.round(p.progress ?? 0)}%</span>
              </div>
              <div className="flex gap-2">
                <Button size="sm" variant="outline" onClick={() => openDetail(p.id)}>Ver detalle</Button>
                {p.status === 'DRAFT' && <Button size="sm" variant="outline" onClick={() => act(p.id, 'start')}><Play size={14} /> Iniciar</Button>}
                {p.status === 'IN_PROGRESS' && <Button size="sm" variant="outline" onClick={() => act(p.id, 'complete')}><CheckCircle size={14} /> Completar</Button>}
                {['DRAFT', 'PENDING', 'IN_PROGRESS', 'BLOCKED'].includes(p.status) && <Button size="sm" variant="ghost" onClick={() => act(p.id, 'cancel')}><XCircle size={14} /></Button>}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      <Dialog open={showModal} onOpenChange={setShowModal}>
        <DialogContent>
          <DialogHeader><DialogTitle>Nuevo Onboarding</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2"><Label>Empleado *</Label><Select options={employees} placeholder="Seleccionar..." value={form.employee_id} onChange={e => setForm({ ...form, employee_id: e.target.value })} /></div>
            <div className="space-y-2"><Label>Fecha de inicio *</Label><Input type="date" value={form.start_date} onChange={e => setForm({ ...form, start_date: e.target.value })} /></div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Tipo de empleado</Label><Select options={['ADMINISTRATIVO', 'DOCENTE', 'TECNICO', 'DESARROLLADOR', 'MANAGER', 'DIRECTOR', 'PASANTE', 'CONTRATISTA'].map(t => ({ value: t, label: t }))} value={form.employee_type} onChange={e => setForm({ ...form, employee_type: e.target.value })} /></div>
              <div className="space-y-2"><Label>Modalidad</Label><Select options={['REMOTO', 'PRESENCIAL', 'HIBRIDO'].map(t => ({ value: t, label: t }))} value={form.work_mode} onChange={e => setForm({ ...form, work_mode: e.target.value })} /></div>
            </div>
            <div className="space-y-2"><Label>Plantilla (opcional)</Label><Input value={form.template_id} onChange={e => setForm({ ...form, template_id: e.target.value })} placeholder="ID de plantilla" /></div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button><Button onClick={create} disabled={saving || !form.employee_id || !form.start_date}>{saving ? 'Creando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={!!detail} onOpenChange={(o) => { if (!o) setDetail(null) }}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">Onboarding · {detail?.employee_id} {detail && pill(detail.status)}</DialogTitle>
            <DialogDescription>Inicio {fmt(detail?.start_date)} · Progreso {Math.round(detail?.progress ?? 0)}% · Política {detail?.completion_policy}</DialogDescription>
          </DialogHeader>
          {detail && (
            <div className="space-y-4">
              <div className="h-2 bg-slate-100 rounded-full overflow-hidden"><div className="h-full rounded-full" style={{ width: `${Math.round(detail.progress ?? 0)}%`, backgroundColor: '#6366f1' }} /></div>
              <div className="flex flex-wrap gap-2">
                <Button size="sm" variant="outline" onClick={() => act(detail.id, 'start')}><Play size={14} /> Iniciar</Button>
                <Button size="sm" variant="outline" onClick={() => act(detail.id, 'block')}><AlertTriangle size={14} /> Bloquear</Button>
                <Button size="sm" variant="outline" onClick={() => act(detail.id, 'complete')}><CheckCircle size={14} /> Completar</Button>
                <Button size="sm" variant="ghost" onClick={() => act(detail.id, 'cancel')}><XCircle size={14} /> Cancelar</Button>
                <Button size="sm" variant="outline" onClick={() => setShowTasks(true)}><Plus size={14} /> Tarea</Button>
                <Button size="sm" variant="outline" onClick={() => setShowDocs(true)}><Plus size={14} /> Documento</Button>
              </div>
              <div className="grid gap-6 lg:grid-cols-2">
                <div>
                  <h3 className="text-sm font-semibold text-slate-900 mb-2">Tareas ({tasks.length})</h3>
                  {tasks.length === 0 ? <p className="text-xs text-slate-500">Sin tareas</p> : (
                    <div className="space-y-2">{tasks.map((t: any) => (
                      <div key={t.id} className="flex items-center justify-between bg-slate-50 rounded-lg p-3">
                        <div>
                          <p className="text-sm font-medium text-slate-800">{t.title}</p>
                          <p className="text-xs text-slate-500">{t.task_type} · {t.assigned_role || t.assigned_to || '-'} {t.due_date ? `· ${fmt(t.due_date)}` : ''}</p>
                        </div>
                        <div className="flex items-center gap-2">
                          {pill(t.status)}
                          {t.status === 'PENDING' && <Button size="sm" variant="ghost" onClick={() => completeTask(t.id)}><CheckCircle size={14} /></Button>}
                        </div>
                      </div>
                    ))}</div>
                  )}
                </div>
                <div>
                  <h3 className="text-sm font-semibold text-slate-900 mb-2">Documentos ({docs.length})</h3>
                  {docs.length === 0 ? <p className="text-xs text-slate-500">Sin documentos</p> : (
                    <div className="space-y-2">{docs.map((d: any) => (
                      <div key={d.id} className="flex items-center justify-between bg-slate-50 rounded-lg p-3">
                        <div>
                          <p className="text-sm font-medium text-slate-800">{d.name}</p>
                          <p className="text-xs text-slate-500">{d.document_type} {d.required ? '· Requerido' : ''}</p>
                        </div>
                        <div className="flex items-center gap-1">
                          {pill(d.status)}
                          {d.status === 'PENDING' && <>
                            <Button size="sm" variant="ghost" onClick={() => docAction(d.id, 'approve')}><CheckCircle size={14} /></Button>
                            <Button size="sm" variant="ghost" onClick={() => docAction(d.id, 'reject')}><XCircle size={14} /></Button>
                          </>}
                        </div>
                      </div>
                    ))}</div>
                  )}
                </div>
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>

      <Dialog open={showTasks} onOpenChange={setShowTasks}>
        <DialogContent>
          <DialogHeader><DialogTitle>Nueva Tarea</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2"><Label>Título *</Label><Input value={taskForm.title} onChange={e => setTaskForm({ ...taskForm, title: e.target.value })} /></div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Tipo</Label><Select options={['DOCUMENT', 'APPROVAL', 'TRAINING', 'ACCOUNT', 'ASSET', 'MEETING', 'SIGNATURE', 'CHECKLIST', 'INFORMATION', 'SYSTEM', 'OTHER'].map(t => ({ value: t, label: t }))} value={taskForm.task_type} onChange={e => setTaskForm({ ...taskForm, task_type: e.target.value })} /></div>
              <div className="space-y-2"><Label>Responsable</Label><Select options={['EMPLOYEE', 'MANAGER', 'HR', 'IT', 'FINANCE', 'SECURITY', 'TRAINING', 'LEGAL', 'EXTERNAL'].map(t => ({ value: t, label: t }))} value={taskForm.assigned_role} onChange={e => setTaskForm({ ...taskForm, assigned_role: e.target.value })} /></div>
            </div>
            <div className="space-y-2"><Label>Vence</Label><Input type="date" value={taskForm.due_date} onChange={e => setTaskForm({ ...taskForm, due_date: e.target.value })} /></div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowTasks(false)}>Cancelar</Button><Button onClick={addTask} disabled={saving || !taskForm.title}>{saving ? 'Guardando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={showDocs} onOpenChange={setShowDocs}>
        <DialogContent>
          <DialogHeader><DialogTitle>Nuevo Documento</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2"><Label>Tipo de documento *</Label><Input value={docForm.document_type} onChange={e => setDocForm({ ...docForm, document_type: e.target.value })} placeholder="Ej: DNI, CONTRATO" /></div>
            <div className="space-y-2"><Label>Nombre *</Label><Input value={docForm.name} onChange={e => setDocForm({ ...docForm, name: e.target.value })} /></div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowDocs(false)}>Cancelar</Button><Button onClick={addDoc} disabled={saving || !docForm.document_type || !docForm.name}>{saving ? 'Guardando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function OffboardingTab() {
  const [list, setList] = useState<any[]>([])
  const [detail, setDetail] = useState<any>(null)
  const [tasks, setTasks] = useState<any[]>([])
  const [assets, setAssets] = useState<any[]>([])
  const [access, setAccess] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [statusFilter, setStatusFilter] = useState('')
  const [showModal, setShowModal] = useState(false)
  const [saving, setSaving] = useState(false)
  const [showTask, setShowTask] = useState(false)
  const [showAsset, setShowAsset] = useState(false)
  const [showAccess, setShowAccess] = useState(false)
  const [taskForm, setTaskForm] = useState({ title: '', task_type: 'OTHER', assigned_role: 'EMPLOYEE' })
  const [assetForm, setAssetForm] = useState({ asset_type: '', description: '', serial_number: '' })
  const [accessForm, setAccessForm] = useState({ system_name: '', access_type: '' })
  const employees = useEmployees()
  const [form, setForm] = useState({ employee_id: '', termination_type: 'RESIGNATION', notice_date: '', last_working_date: '', template_id: '' })

  const fetchList = async () => {
    setLoading(true)
    try {
      const params: Record<string, any> = {}
      if (statusFilter) params.status = statusFilter
      const r = await api.get('/offboarding', { params })
      setList(r.data.data ?? [])
    } catch { setList([]) } finally { setLoading(false) }
  }
  useEffect(() => { fetchList() }, [statusFilter])

  const openDetail = async (id: string) => {
    try {
      const [d, t, a, ac] = await Promise.all([
        api.get(`/offboarding/${id}`),
        api.get(`/offboarding/${id}/tasks`),
        api.get(`/offboarding/${id}/assets`),
        api.get(`/offboarding/${id}/access`),
      ])
      setDetail(d.data.data); setTasks(t.data.data ?? []); setAssets(a.data.data ?? []); setAccess(ac.data.data ?? [])
    } catch { alert('No se pudo cargar el detalle') }
  }

  const create = async () => {
    setSaving(true)
    try {
      const payload: Record<string, any> = {
        employee_id: form.employee_id,
        termination_type: form.termination_type,
        notice_date: new Date(form.notice_date).toISOString(),
        last_working_date: new Date(form.last_working_date).toISOString(),
      }
      if (form.template_id) payload.template_id = form.template_id
      await api.post('/offboarding', payload)
      setShowModal(false); setForm({ employee_id: '', termination_type: 'RESIGNATION', notice_date: '', last_working_date: '', template_id: '' }); fetchList()
    } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }

  const act = async (id: string, action: string) => {
    try { await api.post(`/offboarding/${id}/${action}`); openDetail(id); fetchList() } catch (e: any) { alert(e?.response?.data?.error || 'Error') }
  }

  const completeTask = async (taskId: string) => {
    try { await api.post(`/offboarding/tasks/${taskId}/complete`); openDetail(detail.id) } catch (e: any) { alert(e?.response?.data?.error || 'Error') }
  }

  const addTask = async () => {
    setSaving(true)
    try {
      await api.post(`/offboarding/${detail.id}/tasks`, taskForm)
      setShowTask(false); setTaskForm({ title: '', task_type: 'OTHER', assigned_role: 'EMPLOYEE' }); openDetail(detail.id)
    } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }

  const addAsset = async () => {
    setSaving(true)
    try {
      const payload: Record<string, any> = { asset_type: assetForm.asset_type }
      if (assetForm.description) payload.description = assetForm.description
      if (assetForm.serial_number) payload.serial_number = assetForm.serial_number
      await api.post(`/offboarding/${detail.id}/assets`, payload)
      setShowAsset(false); setAssetForm({ asset_type: '', description: '', serial_number: '' }); openDetail(detail.id)
    } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }

  const assetAction = async (assetId: string, action: 'return' | 'report-damaged' | 'report-lost') => {
    try {
      if (action === 'return') {
        const condition = window.prompt('Condición de devolución:', 'BUENO')
        if (condition === null) return
        await api.post(`/offboarding/assets/${assetId}/${action}`, { condition_on_return: condition, status: 'RETURNED' })
      } else {
        await api.post(`/offboarding/assets/${assetId}/${action}`, { status: action === 'report-damaged' ? 'DAMAGED' : 'LOST' })
      }
      openDetail(detail.id)
    } catch (e: any) { alert(e?.response?.data?.error || 'Error') }
  }

  const addAccess = async () => {
    setSaving(true)
    try {
      await api.post(`/offboarding/${detail.id}/access`, accessForm)
      setShowAccess(false); setAccessForm({ system_name: '', access_type: '' }); openDetail(detail.id)
    } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }

  const accessAction = async (accessId: string, action: 'revoke' | 'retry') => {
    try { await api.post(`/offboarding/access/${accessId}/${action}`); openDetail(detail.id) } catch (e: any) { alert(e?.response?.data?.error || 'Error') }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4 flex-wrap gap-2">
        <div className="w-48"><Select options={[{ value: '', label: 'Todos los estados' }, ...['DRAFT', 'PENDING_APPROVAL', 'APPROVED', 'IN_PROGRESS', 'BLOCKED', 'COMPLETED', 'CANCELLED'].map(s => ({ value: s, label: s }))]} value={statusFilter} onChange={e => setStatusFilter(e.target.value)} /></div>
        <Button size="sm" onClick={() => setShowModal(true)}><Plus size={14} /> Nuevo Offboarding</Button>
      </div>
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {loading ? <div className="p-6 text-center text-slate-500 col-span-full">Cargando...</div> : list.length === 0 ? <div className="p-6 text-center text-slate-500 col-span-full">No hay procesos</div> : list.map((p: any) => (
          <Card key={p.id}>
            <CardContent className="p-5">
              <div className="flex items-start justify-between mb-2">
                <div>
                  <p className="font-semibold text-slate-900">{p.employee_id}</p>
                  <p className="text-xs text-slate-500">{p.termination_type} · Último día {fmt(p.last_working_date)}</p>
                </div>
                {pill(p.status)}
              </div>
              <div className="flex items-center gap-2 mb-3">
                <div className="flex-1 h-1.5 bg-slate-100 rounded-full overflow-hidden"><div className="h-full rounded-full" style={{ width: `${Math.round(p.progress ?? 0)}%`, backgroundColor: '#6366f1' }} /></div>
                <span className="text-xs text-slate-500">{Math.round(p.progress ?? 0)}%</span>
              </div>
              <div className="flex gap-2 flex-wrap">
                <Button size="sm" variant="outline" onClick={() => openDetail(p.id)}>Ver detalle</Button>
                {p.status === 'DRAFT' && <Button size="sm" variant="outline" onClick={() => act(p.id, 'approve')}>Aprobar</Button>}
                {['APPROVED', 'DRAFT'].includes(p.status) && <Button size="sm" variant="outline" onClick={() => act(p.id, 'start')}><Play size={14} /> Iniciar</Button>}
                {p.status === 'IN_PROGRESS' && <Button size="sm" variant="outline" onClick={() => act(p.id, 'complete')}><CheckCircle size={14} /> Completar</Button>}
                {['DRAFT', 'PENDING_APPROVAL', 'APPROVED', 'IN_PROGRESS', 'BLOCKED'].includes(p.status) && <Button size="sm" variant="ghost" onClick={() => act(p.id, 'cancel')}><XCircle size={14} /></Button>}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      <Dialog open={showModal} onOpenChange={setShowModal}>
        <DialogContent>
          <DialogHeader><DialogTitle>Nuevo Offboarding</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2"><Label>Empleado *</Label><Select options={employees} placeholder="Seleccionar..." value={form.employee_id} onChange={e => setForm({ ...form, employee_id: e.target.value })} /></div>
            <div className="space-y-2"><Label>Tipo de baja *</Label><Select options={['RESIGNATION', 'RETIREMENT', 'TERMINATION', 'END_OF_CONTRACT', 'MUTUAL_AGREEMENT', 'LAYOFF', 'TRANSFER', 'OTHER'].map(t => ({ value: t, label: t }))} value={form.termination_type} onChange={e => setForm({ ...form, termination_type: e.target.value })} /></div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Notificación *</Label><Input type="date" value={form.notice_date} onChange={e => setForm({ ...form, notice_date: e.target.value })} /></div>
              <div className="space-y-2"><Label>Último día *</Label><Input type="date" value={form.last_working_date} onChange={e => setForm({ ...form, last_working_date: e.target.value })} /></div>
            </div>
            <div className="space-y-2"><Label>Plantilla (opcional)</Label><Input value={form.template_id} onChange={e => setForm({ ...form, template_id: e.target.value })} placeholder="ID de plantilla" /></div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button><Button onClick={create} disabled={saving || !form.employee_id || !form.termination_type || !form.notice_date || !form.last_working_date}>{saving ? 'Creando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={!!detail} onOpenChange={(o) => { if (!o) setDetail(null) }}>
        <DialogContent className="max-w-3xl">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">Offboarding · {detail?.employee_id} {detail && pill(detail.status)}</DialogTitle>
            <DialogDescription>{detail?.termination_type} · Notificado {fmt(detail?.notice_date)} · Último día {fmt(detail?.last_working_date)}</DialogDescription>
          </DialogHeader>
          {detail && (
            <div className="space-y-4">
              <div className="h-2 bg-slate-100 rounded-full overflow-hidden"><div className="h-full rounded-full" style={{ width: `${Math.round(detail.progress ?? 0)}%`, backgroundColor: '#6366f1' }} /></div>
              <div className="flex flex-wrap gap-2">
                <Button size="sm" variant="outline" onClick={() => act(detail.id, 'approve')}>Aprobar</Button>
                <Button size="sm" variant="outline" onClick={() => act(detail.id, 'start')}><Play size={14} /> Iniciar</Button>
                <Button size="sm" variant="outline" onClick={() => act(detail.id, 'complete')}><CheckCircle size={14} /> Completar</Button>
                <Button size="sm" variant="ghost" onClick={() => act(detail.id, 'cancel')}><XCircle size={14} /> Cancelar</Button>
                <Button size="sm" variant="outline" onClick={() => setShowTask(true)}><Plus size={14} /> Tarea</Button>
                <Button size="sm" variant="outline" onClick={() => setShowAsset(true)}><Plus size={14} /> Asset</Button>
                <Button size="sm" variant="outline" onClick={() => setShowAccess(true)}><Plus size={14} /> Acceso</Button>
              </div>
              <div className="grid gap-6 lg:grid-cols-3">
                <div>
                  <h3 className="text-sm font-semibold text-slate-900 mb-2">Tareas ({tasks.length})</h3>
                  {tasks.length === 0 ? <p className="text-xs text-slate-500">Sin tareas</p> : (
                    <div className="space-y-2">{tasks.map((t: any) => (
                      <div key={t.id} className="flex items-center justify-between bg-slate-50 rounded-lg p-3">
                        <div>
                          <p className="text-sm font-medium text-slate-800">{t.title}</p>
                          <p className="text-xs text-slate-500">{t.task_type} · {t.assigned_role || t.assigned_to || '-'}</p>
                        </div>
                        <div className="flex items-center gap-2">
                          {pill(t.status)}
                          {t.status === 'PENDING' && <Button size="sm" variant="ghost" onClick={() => completeTask(t.id)}><CheckCircle size={14} /></Button>}
                        </div>
                      </div>
                    ))}</div>
                  )}
                </div>
                <div>
                  <h3 className="text-sm font-semibold text-slate-900 mb-2">Assets ({assets.length})</h3>
                  {assets.length === 0 ? <p className="text-xs text-slate-500">Sin assets</p> : (
                    <div className="space-y-2">{assets.map((a: any) => (
                      <div key={a.id} className="flex items-center justify-between bg-slate-50 rounded-lg p-3">
                        <div>
                          <p className="text-sm font-medium text-slate-800">{a.asset_type}</p>
                          <p className="text-xs text-slate-500">{a.description || a.serial_number || '-'}</p>
                        </div>
                        <div className="flex items-center gap-1">
                          {pill(a.status)}
                          {['PENDING_RETURN', 'PENDING'].includes(a.status) && <Button size="sm" variant="ghost" onClick={() => assetAction(a.id, 'return')}><CheckCircle size={14} /></Button>}
                        </div>
                      </div>
                    ))}</div>
                  )}
                </div>
                <div>
                  <h3 className="text-sm font-semibold text-slate-900 mb-2">Accesos ({access.length})</h3>
                  {access.length === 0 ? <p className="text-xs text-slate-500">Sin accesos</p> : (
                    <div className="space-y-2">{access.map((a: any) => (
                      <div key={a.id} className="flex items-center justify-between bg-slate-50 rounded-lg p-3">
                        <div>
                          <p className="text-sm font-medium text-slate-800">{a.system_name}</p>
                          <p className="text-xs text-slate-500">{a.access_type || '-'}</p>
                        </div>
                        <div className="flex items-center gap-1">
                          {pill(a.status)}
                          {['PENDING', 'FAILED'].includes(a.status) && <Button size="sm" variant="ghost" onClick={() => accessAction(a.id, a.status === 'FAILED' ? 'retry' : 'revoke')}><CheckCircle size={14} /></Button>}
                        </div>
                      </div>
                    ))}</div>
                  )}
                </div>
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>

      <Dialog open={showTask} onOpenChange={setShowTask}>
        <DialogContent>
          <DialogHeader><DialogTitle>Nueva Tarea</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2"><Label>Título *</Label><Input value={taskForm.title} onChange={e => setTaskForm({ ...taskForm, title: e.target.value })} /></div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Tipo</Label><Select options={['DOCUMENT', 'APPROVAL', 'TRAINING', 'ACCOUNT', 'ASSET', 'MEETING', 'SIGNATURE', 'CHECKLIST', 'INFORMATION', 'SYSTEM', 'OTHER'].map(t => ({ value: t, label: t }))} value={taskForm.task_type} onChange={e => setTaskForm({ ...taskForm, task_type: e.target.value })} /></div>
              <div className="space-y-2"><Label>Responsable</Label><Select options={['EMPLOYEE', 'MANAGER', 'HR', 'IT', 'FINANCE', 'SECURITY', 'TRAINING', 'LEGAL', 'EXTERNAL'].map(t => ({ value: t, label: t }))} value={taskForm.assigned_role} onChange={e => setTaskForm({ ...taskForm, assigned_role: e.target.value })} /></div>
            </div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowTask(false)}>Cancelar</Button><Button onClick={addTask} disabled={saving || !taskForm.title}>{saving ? 'Guardando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={showAsset} onOpenChange={setShowAsset}>
        <DialogContent>
          <DialogHeader><DialogTitle>Nuevo Asset</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2"><Label>Tipo *</Label><Input value={assetForm.asset_type} onChange={e => setAssetForm({ ...assetForm, asset_type: e.target.value })} placeholder="Ej: NOTEBOOK, CELULAR" /></div>
            <div className="space-y-2"><Label>Descripción</Label><Input value={assetForm.description} onChange={e => setAssetForm({ ...assetForm, description: e.target.value })} /></div>
            <div className="space-y-2"><Label>N° de serie</Label><Input value={assetForm.serial_number} onChange={e => setAssetForm({ ...assetForm, serial_number: e.target.value })} /></div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowAsset(false)}>Cancelar</Button><Button onClick={addAsset} disabled={saving || !assetForm.asset_type}>{saving ? 'Guardando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={showAccess} onOpenChange={setShowAccess}>
        <DialogContent>
          <DialogHeader><DialogTitle>Nueva Revocación</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2"><Label>Sistema *</Label><Input value={accessForm.system_name} onChange={e => setAccessForm({ ...accessForm, system_name: e.target.value })} placeholder="Ej: GMAIL, SLACK, ERP" /></div>
            <div className="space-y-2"><Label>Tipo de acceso</Label><Input value={accessForm.access_type} onChange={e => setAccessForm({ ...accessForm, access_type: e.target.value })} placeholder="Ej: EMAIL, VPN" /></div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowAccess(false)}>Cancelar</Button><Button onClick={addAccess} disabled={saving || !accessForm.system_name}>{saving ? 'Guardando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
