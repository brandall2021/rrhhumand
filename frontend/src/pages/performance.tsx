import { useEffect, useState } from 'react'
import { Plus, Calculator } from 'lucide-react'
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
  <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${s === 'COMPLETED' || s === 'APPROVED' || s === 'ACTIVE' ? 'bg-emerald-50 text-emerald-700' : s === 'PENDING' || s === 'IN_PROGRESS' ? 'bg-amber-50 text-amber-700' : s === 'CANCELLED' ? 'bg-slate-100 text-slate-500' : 'bg-sky-50 text-sky-700'}`}>{s}</span>
)

export default function PerformancePage() {
  const [activeTab, setActiveTab] = useState('cycles')
  const [dash, setDash] = useState<any>(null)
  useEffect(() => { api.get('/performance/dashboard').then(r => setDash(r.data.data)).catch(() => {}) }, [])
  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-900 mb-6">Desempeño</h1>
      {dash && (
        <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-8 gap-3 mb-6">
          {[
            ['Ciclos totales', dash.total_cycles], ['Ciclos activos', dash.active_cycles],
            ['Evaluaciones', dash.total_evaluations], ['Pendientes', dash.pending_evaluations],
            ['Score promedio', dash.average_score?.toFixed?.(2) ?? dash.average_score],
            ['Objetivos', dash.total_objectives], ['Feedback', dash.total_feedback],
            ['Planes', (dash.total_improvement_plans ?? 0) + (dash.total_development_plans ?? 0)],
          ].map(([l, v]) => (
            <Card key={l as string}><CardContent className="p-4"><p className="text-xs text-slate-500 mb-1">{l}</p><p className="text-lg font-bold text-slate-900">{v}</p></CardContent></Card>
          ))}
        </div>
      )}
      <div className="flex gap-1 mb-6 border-b border-slate-200">
        {[['cycles', 'Ciclos'], ['objectives', 'Objetivos'], ['evaluations', 'Evaluaciones'], ['feedback', 'Feedback'], ['plans', 'Planes']].map(([k, l]) => (
          <button key={k} onClick={() => setActiveTab(k)} className={`px-4 py-2.5 text-sm font-medium border-b-2 transition-colors ${activeTab === k ? 'border-brand-600 text-brand-700' : 'border-transparent text-slate-500 hover:text-slate-700'}`}>{l}</button>
        ))}
      </div>
      {activeTab === 'cycles' && <CyclesTab />}
      {activeTab === 'objectives' && <ObjectivesTab />}
      {activeTab === 'evaluations' && <EvaluationsTab />}
      {activeTab === 'feedback' && <FeedbackTab />}
      {activeTab === 'plans' && <PlansTab />}
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

function CyclesTab() {
  const [cycles, setCycles] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [saving, setSaving] = useState(false)
  const [form, setForm] = useState({ name: '', cycle_type: 'ANUAL', description: '', start_date: '', end_date: '', evaluation_start_date: '', evaluation_end_date: '' })

  const fetch = async () => {
    setLoading(true)
    try { const r = await api.get('/performance/cycles'); setCycles(r.data.data ?? []) } catch { setCycles([]) } finally { setLoading(false) }
  }
  useEffect(() => { fetch() }, [])

  const create = async () => {
    setSaving(true)
    try {
      const payload: Record<string, any> = { name: form.name, cycle_type: form.cycle_type }
      if (form.description) payload.description = form.description
      for (const k of ['start_date', 'end_date', 'evaluation_start_date', 'evaluation_end_date']) if (form[k as keyof typeof form]) payload[k] = new Date(form[k as keyof typeof form]).toISOString()
      await api.post('/performance/cycles', payload)
      setShowModal(false); setForm({ name: '', cycle_type: 'ANUAL', description: '', start_date: '', end_date: '', evaluation_start_date: '', evaluation_end_date: '' }); fetch()
    } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }

  const changeStatus = async (c: any, status: string) => {
    try { await api.post(`/performance/cycles/${c.id}/status`, { status }); fetch() } catch (e: any) { alert(e?.response?.data?.error || 'Error') }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-slate-500">{cycles.length} ciclos</p>
        <Button size="sm" onClick={() => setShowModal(true)}><Plus size={14} /> Nuevo Ciclo</Button>
      </div>
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {loading ? <div className="p-6 text-center text-slate-500 col-span-full">Cargando...</div> : cycles.length === 0 ? <div className="p-6 text-center text-slate-500 col-span-full">No hay ciclos</div> : cycles.map((c: any) => (
          <Card key={c.id}>
            <CardContent className="p-5">
              <div className="flex items-start justify-between mb-2">
                <div>
                  <h3 className="font-semibold text-slate-900">{c.name}</h3>
                  <p className="text-xs text-slate-500">{c.cycle_type}</p>
                </div>
                {pill(c.status)}
              </div>
              <p className="text-sm text-slate-600 mb-3">{c.description || 'Sin descripción'}</p>
              <div className="grid grid-cols-2 gap-2 text-xs text-slate-500 mb-3">
                <div>Inicio: <span className="font-medium text-slate-700">{fmt(c.start_date)}</span></div>
                <div>Fin: <span className="font-medium text-slate-700">{fmt(c.end_date)}</span></div>
                {c.objective_weight != null && <div>Objetivos: {c.objective_weight}%</div>}
                {c.competency_weight != null && <div>Competencias: {c.competency_weight}%</div>}
              </div>
              <div className="flex gap-2">
                {c.status === 'DRAFT' && <Button size="sm" variant="outline" onClick={() => changeStatus(c, 'OPEN')}>Abrir</Button>}
                {['OPEN', 'IN_PROGRESS'].includes(c.status) && <Button size="sm" variant="outline" onClick={() => changeStatus(c, 'REVIEW')}>Revisión</Button>}
                {c.status === 'REVIEW' && <Button size="sm" variant="outline" onClick={() => changeStatus(c, 'CLOSED')}>Cerrar</Button>}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
      <Dialog open={showModal} onOpenChange={setShowModal}>
        <DialogContent>
          <DialogHeader><DialogTitle>Nuevo Ciclo de Desempeño</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Nombre *</Label><Input value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} placeholder="Ej: Evaluación 2026" /></div>
              <div className="space-y-2"><Label>Tipo</Label><Select options={['ANUAL', 'SEMESTRAL', 'TRIMESTRAL', 'PROYECTO'].map(t => ({ value: t, label: t }))} value={form.cycle_type} onChange={e => setForm({ ...form, cycle_type: e.target.value })} /></div>
            </div>
            <div className="space-y-2"><Label>Descripción</Label><Input value={form.description} onChange={e => setForm({ ...form, description: e.target.value })} /></div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Inicio</Label><Input type="date" value={form.start_date} onChange={e => setForm({ ...form, start_date: e.target.value })} /></div>
              <div className="space-y-2"><Label>Fin</Label><Input type="date" value={form.end_date} onChange={e => setForm({ ...form, end_date: e.target.value })} /></div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Evaluación desde</Label><Input type="date" value={form.evaluation_start_date} onChange={e => setForm({ ...form, evaluation_start_date: e.target.value })} /></div>
              <div className="space-y-2"><Label>Evaluación hasta</Label><Input type="date" value={form.evaluation_end_date} onChange={e => setForm({ ...form, evaluation_end_date: e.target.value })} /></div>
            </div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button><Button onClick={create} disabled={saving || !form.name}>{saving ? 'Creando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function ObjectivesTab() {
  const [objectives, setObjectives] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const employees = useEmployees()
  const [filterEmp, setFilterEmp] = useState('')
  const [showModal, setShowModal] = useState(false)
  const [saving, setSaving] = useState(false)
  const [form, setForm] = useState({ cycle_id: '', employee_id: '', title: '', objective_type: 'OKR', weight: '', target_value: '', unit: '' })

  const fetch = async () => {
    setLoading(true)
    try { const params: Record<string, any> = {}; if (filterEmp) params.employee_id = filterEmp; const r = await api.get('/performance/objectives', { params }); setObjectives(r.data.data ?? []) } catch { setObjectives([]) } finally { setLoading(false) }
  }
  useEffect(() => { fetch() }, [filterEmp])

  const create = async () => {
    setSaving(true)
    try {
      const payload: Record<string, any> = { title: form.title, objective_type: form.objective_type, ...(form.cycle_id ? { cycle_id: form.cycle_id } : {}), ...(form.employee_id ? { employee_id: form.employee_id } : {}) }
      if (form.weight) payload.weight = parseFloat(form.weight)
      if (form.target_value) payload.target_value = parseFloat(form.target_value)
      if (form.unit) payload.unit = form.unit
      await api.post('/performance/objectives', payload)
      setShowModal(false); setForm({ cycle_id: '', employee_id: '', title: '', objective_type: 'OKR', weight: '', target_value: '', unit: '' }); fetch()
    } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }

  const updateProgress = async (o: any) => {
    const val = window.prompt('Progreso actual:', String(o.current_value ?? 0))
    if (val === null) return
    try { await api.post(`/performance/objectives/${o.id}/progress`, { current_value: parseFloat(val) || 0 }); fetch() } catch (e: any) { alert(e?.response?.data?.error || 'Error') }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4 flex-wrap gap-2">
        <div className="w-72"><Label>Empleado</Label><Select options={[{ value: '', label: 'Todos' }, ...employees]} value={filterEmp} onChange={e => setFilterEmp(e.target.value)} /></div>
        <Button size="sm" onClick={() => setShowModal(true)}><Plus size={14} /> Nuevo Objetivo</Button>
      </div>
      <Card>
        <CardContent className="p-0">
          {loading ? <div className="p-6 text-center text-slate-500">Cargando...</div> : objectives.length === 0 ? <div className="p-6 text-center text-slate-500">No hay objetivos</div> : (
            <table className="w-full text-sm">
              <thead><tr className="border-b border-slate-200 bg-slate-50">
                <th className="text-left px-4 py-3 font-medium text-slate-600">Objetivo</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Tipo</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Progreso</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Meta</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
              </tr></thead>
              <tbody>{objectives.map((o: any) => {
                const pct = o.target_value ? Math.min(100, Math.round(((o.current_value ?? 0) / o.target_value) * 100)) : 0
                return (
                  <tr key={o.id} className="border-b border-slate-100 hover:bg-slate-50">
                    <td className="px-4 py-3">
                      <p className="font-medium text-slate-900">{o.title}</p>
                      <p className="text-xs text-slate-500">{o.employee_id}</p>
                    </td>
                    <td className="px-4 py-3 text-slate-600">{o.objective_type}</td>
                    <td className="px-4 py-3 w-40">
                      <div className="flex items-center gap-2">
                        <div className="flex-1 h-1.5 bg-slate-100 rounded-full overflow-hidden"><div className="h-full bg-brand-600 rounded-full" style={{ width: `${pct}%` }} /></div>
                        <span className="text-xs text-slate-500">{pct}%</span>
                      </div>
                    </td>
                    <td className="px-4 py-3 text-slate-600">{o.current_value ?? 0} / {o.target_value ?? '-'} {o.unit || ''}</td>
                    <td className="px-4 py-3">{pill(o.status)}</td>
                    <td className="px-4 py-3"><div className="flex justify-end"><Button size="sm" variant="outline" onClick={() => updateProgress(o)}>Actualizar</Button></div></td>
                  </tr>
                )
              })}</tbody>
            </table>
          )}
        </CardContent>
      </Card>
      <Dialog open={showModal} onOpenChange={setShowModal}>
        <DialogContent>
          <DialogHeader><DialogTitle>Nuevo Objetivo</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2"><Label>Título *</Label><Input value={form.title} onChange={e => setForm({ ...form, title: e.target.value })} /></div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Tipo</Label><Select options={['OKR', 'KPI', 'TAREA', 'DESARROLLO'].map(t => ({ value: t, label: t }))} value={form.objective_type} onChange={e => setForm({ ...form, objective_type: e.target.value })} /></div>
              <div className="space-y-2"><Label>Empleado</Label><Select options={employees} placeholder="Opcional" value={form.employee_id} onChange={e => setForm({ ...form, employee_id: e.target.value })} /></div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Peso (%)</Label><Input type="number" value={form.weight} onChange={e => setForm({ ...form, weight: e.target.value })} /></div>
              <div className="space-y-2"><Label>Meta</Label><Input type="number" value={form.target_value} onChange={e => setForm({ ...form, target_value: e.target.value })} /></div>
            </div>
            <div className="space-y-2"><Label>Unidad</Label><Input value={form.unit} onChange={e => setForm({ ...form, unit: e.target.value })} placeholder="Ej: %, USD, unidades" /></div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button><Button onClick={create} disabled={saving || !form.title}>{saving ? 'Guardando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function EvaluationsTab() {
  const [evals, setEvals] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const employees = useEmployees()
  const [showModal, setShowModal] = useState(false)
  const [saving, setSaving] = useState(false)
  const [form, setForm] = useState({ cycle_id: '', employee_id: '', evaluator_id: '', evaluation_type: 'AUTO' })
  const [calcForm, setCalcForm] = useState({ cycle_id: '', employee_id: '' })
  const [showCalc, setShowCalc] = useState(false)

  const fetch = async () => {
    setLoading(true)
    try { const r = await api.get('/performance/evaluations'); setEvals(r.data.data ?? []) } catch { setEvals([]) } finally { setLoading(false) }
  }
  useEffect(() => { fetch() }, [])

  const create = async () => {
    setSaving(true)
    try {
      const payload: Record<string, any> = { employee_id: form.employee_id, evaluator_id: form.evaluator_id, evaluation_type: form.evaluation_type }
      if (form.cycle_id) payload.cycle_id = form.cycle_id
      await api.post('/performance/evaluations', payload)
      setShowModal(false); setForm({ cycle_id: '', employee_id: '', evaluator_id: '', evaluation_type: 'AUTO' }); fetch()
    } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }

  const act = async (id: string, action: string) => {
    try { await api.post(`/performance/evaluations/${id}/${action}`); fetch() } catch (e: any) { alert(e?.response?.data?.error || 'Error') }
  }

  const calculate = async () => {
    setSaving(true)
    try { await api.post('/performance/results/calculate', calcForm); setShowCalc(false); setCalcForm({ cycle_id: '', employee_id: '' }) } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-slate-500">{evals.length} evaluaciones</p>
        <div className="flex gap-2">
          <Button size="sm" variant="outline" onClick={() => setShowCalc(true)}><Calculator size={14} /> Calcular Resultado</Button>
          <Button size="sm" onClick={() => setShowModal(true)}><Plus size={14} /> Nueva Evaluación</Button>
        </div>
      </div>
      <Card>
        <CardContent className="p-0">
          {loading ? <div className="p-6 text-center text-slate-500">Cargando...</div> : evals.length === 0 ? <div className="p-6 text-center text-slate-500">No hay evaluaciones</div> : (
            <table className="w-full text-sm">
              <thead><tr className="border-b border-slate-200 bg-slate-50">
                <th className="text-left px-4 py-3 font-medium text-slate-600">Empleado</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Evaluador</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Tipo</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
              </tr></thead>
              <tbody>{evals.map((e: any) => (
                <tr key={e.id} className="border-b border-slate-100 hover:bg-slate-50">
                  <td className="px-4 py-3 font-medium text-slate-900">{e.employee_id}</td>
                  <td className="px-4 py-3 text-slate-600">{e.evaluator_id}</td>
                  <td className="px-4 py-3 text-slate-600">{e.evaluation_type}</td>
                  <td className="px-4 py-3">{pill(e.status)}</td>
                  <td className="px-4 py-3">
                    <div className="flex justify-end gap-1">
                      {['DRAFT', 'PENDING'].includes(e.status) && <Button size="sm" variant="outline" onClick={() => act(e.id, 'submit')}>Enviar</Button>}
                      {e.status === 'SUBMITTED' && <Button size="sm" variant="outline" onClick={() => act(e.id, 'approve')}>Aprobar</Button>}
                      {['SUBMITTED', 'APPROVED'].includes(e.status) && <Button size="sm" variant="ghost" onClick={() => act(e.id, 'lock')}>Bloquear</Button>}
                      {['LOCKED', 'APPROVED'].includes(e.status) && <Button size="sm" variant="ghost" onClick={() => act(e.id, 'reopen')}>Reabrir</Button>}
                    </div>
                  </td>
                </tr>
              ))}</tbody>
            </table>
          )}
        </CardContent>
      </Card>
      <Dialog open={showModal} onOpenChange={setShowModal}>
        <DialogContent>
          <DialogHeader><DialogTitle>Nueva Evaluación</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2"><Label>Empleado *</Label><Select options={employees} placeholder="Seleccionar..." value={form.employee_id} onChange={e => setForm({ ...form, employee_id: e.target.value })} /></div>
            <div className="space-y-2"><Label>Evaluador *</Label><Select options={employees} placeholder="Seleccionar..." value={form.evaluator_id} onChange={e => setForm({ ...form, evaluator_id: e.target.value })} /></div>
            <div className="space-y-2"><Label>Tipo *</Label><Select options={['AUTO', 'SELF', 'PEER', 'MANAGER', '360'].map(t => ({ value: t, label: t }))} value={form.evaluation_type} onChange={e => setForm({ ...form, evaluation_type: e.target.value })} /></div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button><Button onClick={create} disabled={saving || !form.employee_id || !form.evaluator_id}>{saving ? 'Guardando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>
      <Dialog open={showCalc} onOpenChange={setShowCalc}>
        <DialogContent>
          <DialogHeader><DialogTitle>Calcular Resultado</DialogTitle><DialogDescription>Genera el puntaje de desempeño de un empleado en un ciclo</DialogDescription></DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2"><Label>Empleado *</Label><Select options={employees} placeholder="Seleccionar..." value={calcForm.employee_id} onChange={e => setCalcForm({ ...calcForm, employee_id: e.target.value })} /></div>
            <div className="space-y-2"><Label>Ciclo</Label><Input value={calcForm.cycle_id} onChange={e => setCalcForm({ ...calcForm, cycle_id: e.target.value })} placeholder="ID del ciclo (opcional)" /></div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowCalc(false)}>Cancelar</Button><Button onClick={calculate} disabled={saving || !calcForm.employee_id}>{saving ? 'Calculando...' : 'Calcular'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function FeedbackTab() {
  const [feedback, setFeedback] = useState<any[]>([])
  const [recognitions, setRecognitions] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const employees = useEmployees()
  const [showFb, setShowFb] = useState(false)
  const [showRec, setShowRec] = useState(false)
  const [saving, setSaving] = useState(false)
  const [fbForm, setFbForm] = useState({ employee_id: '', feedback_type: 'PRAISE', visibility: 'PRIVATE', content: '' })
  const [recForm, setRecForm] = useState({ employee_id: '', recognition_type: 'EXCELLENT', message: '' })

  const fetch = async () => {
    setLoading(true)
    try {
      const [f, r] = await Promise.all([api.get('/performance/feedback'), api.get('/performance/recognitions')])
      setFeedback(f.data.data ?? []); setRecognitions(r.data.data ?? [])
    } catch { setFeedback([]); setRecognitions([]) } finally { setLoading(false) }
  }
  useEffect(() => { fetch() }, [])

  const createFb = async () => {
    setSaving(true)
    try { await api.post('/performance/feedback', fbForm); setShowFb(false); setFbForm({ employee_id: '', feedback_type: 'PRAISE', visibility: 'PRIVATE', content: '' }); fetch() } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }
  const createRec = async () => {
    setSaving(true)
    try { await api.post('/performance/recognitions', recForm); setShowRec(false); setRecForm({ employee_id: '', recognition_type: 'EXCELLENT', message: '' }); fetch() } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-slate-500">{feedback.length} feedback · {recognitions.length} reconocimientos</p>
        <div className="flex gap-2">
          <Button size="sm" variant="outline" onClick={() => setShowRec(true)}><Plus size={14} /> Reconocimiento</Button>
          <Button size="sm" onClick={() => setShowFb(true)}><Plus size={14} /> Nuevo Feedback</Button>
        </div>
      </div>
      <div className="grid gap-6 lg:grid-cols-2">
        <div>
          <h2 className="font-semibold text-slate-900 mb-3">Feedback</h2>
          {loading ? <p className="text-sm text-slate-500">Cargando...</p> : feedback.length === 0 ? <p className="text-sm text-slate-500">Sin feedback</p> : (
            <div className="space-y-3">{feedback.map((f: any) => (
              <Card key={f.id}><CardContent className="p-4">
                <div className="flex items-center justify-between mb-2">
                  <div className="flex items-center gap-2">
                    <span className="text-xs font-medium text-brand-700 bg-brand-50 px-2 py-0.5 rounded-full">{f.feedback_type}</span>
                    <span className="text-xs text-slate-500">{f.is_anonymous ? 'Anónimo' : f.author_id}</span>
                  </div>
                  <span className="text-xs text-slate-500">{f.visibility}</span>
                </div>
                <p className="text-sm text-slate-700">{f.content}</p>
                <p className="text-xs text-slate-400 mt-2">→ {f.employee_id}</p>
              </CardContent></Card>
            ))}</div>
          )}
        </div>
        <div>
          <h2 className="font-semibold text-slate-900 mb-3">Reconocimientos</h2>
          {loading ? <p className="text-sm text-slate-500">Cargando...</p> : recognitions.length === 0 ? <p className="text-sm text-slate-500">Sin reconocimientos</p> : (
            <div className="space-y-3">{recognitions.map((r: any) => (
              <Card key={r.id}><CardContent className="p-4">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-xs font-medium text-violet-700 bg-violet-50 px-2 py-0.5 rounded-full">{r.recognition_type}</span>
                  <span className="text-xs text-slate-400">de {r.author_id}</span>
                </div>
                <p className="text-sm text-slate-700">{r.message}</p>
                <p className="text-xs text-slate-400 mt-2">→ {r.employee_id}</p>
              </CardContent></Card>
            ))}</div>
          )}
        </div>
      </div>
      <Dialog open={showFb} onOpenChange={setShowFb}>
        <DialogContent>
          <DialogHeader><DialogTitle>Nuevo Feedback</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2"><Label>Empleado *</Label><Select options={employees} placeholder="Seleccionar..." value={fbForm.employee_id} onChange={e => setFbForm({ ...fbForm, employee_id: e.target.value })} /></div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Tipo *</Label><Select options={['PRAISE', 'SUGGESTION', 'WARNING', 'OTHER'].map(t => ({ value: t, label: t }))} value={fbForm.feedback_type} onChange={e => setFbForm({ ...fbForm, feedback_type: e.target.value })} /></div>
              <div className="space-y-2"><Label>Visibilidad</Label><Select options={['PRIVATE', 'PUBLIC', 'TEAM'].map(t => ({ value: t, label: t }))} value={fbForm.visibility} onChange={e => setFbForm({ ...fbForm, visibility: e.target.value })} /></div>
            </div>
            <div className="space-y-2"><Label>Contenido *</Label><Input value={fbForm.content} onChange={e => setFbForm({ ...fbForm, content: e.target.value })} /></div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowFb(false)}>Cancelar</Button><Button onClick={createFb} disabled={saving || !fbForm.employee_id || !fbForm.content}>{saving ? 'Guardando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>
      <Dialog open={showRec} onOpenChange={setShowRec}>
        <DialogContent>
          <DialogHeader><DialogTitle>Nuevo Reconocimiento</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2"><Label>Empleado *</Label><Select options={employees} placeholder="Seleccionar..." value={recForm.employee_id} onChange={e => setRecForm({ ...recForm, employee_id: e.target.value })} /></div>
            <div className="space-y-2"><Label>Tipo *</Label><Select options={['EXCELLENT', 'GOOD', 'COLLABORATION', 'INNOVATION', 'LEADERSHIP'].map(t => ({ value: t, label: t }))} value={recForm.recognition_type} onChange={e => setRecForm({ ...recForm, recognition_type: e.target.value })} /></div>
            <div className="space-y-2"><Label>Mensaje *</Label><Input value={recForm.message} onChange={e => setRecForm({ ...recForm, message: e.target.value })} /></div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowRec(false)}>Cancelar</Button><Button onClick={createRec} disabled={saving || !recForm.employee_id || !recForm.message}>{saving ? 'Guardando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function PlansTab() {
  const [improvement, setImprovement] = useState<any[]>([])
  const [development, setDevelopment] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const employees = useEmployees()
  const [showImp, setShowImp] = useState(false)
  const [showDev, setShowDev] = useState(false)
  const [saving, setSaving] = useState(false)
  const [impForm, setImpForm] = useState({ employee_id: '', reason: '', start_date: '', end_date: '', success_criteria: '' })
  const [devForm, setDevForm] = useState({ employee_id: '', title: '', description: '', career_goal: '' })

  const fetch = async () => {
    setLoading(true)
    try {
      const [i, d] = await Promise.all([api.get('/performance/improvement-plans'), api.get('/performance/development-plans')])
      setImprovement(i.data.data ?? []); setDevelopment(d.data.data ?? [])
    } catch { setImprovement([]); setDevelopment([]) } finally { setLoading(false) }
  }
  useEffect(() => { fetch() }, [])

  const createImp = async () => {
    setSaving(true)
    try {
      const payload: Record<string, any> = { employee_id: impForm.employee_id, reason: impForm.reason, start_date: new Date(impForm.start_date).toISOString(), end_date: new Date(impForm.end_date).toISOString() }
      if (impForm.success_criteria) payload.success_criteria = impForm.success_criteria
      await api.post('/performance/improvement-plans', payload)
      setShowImp(false); setImpForm({ employee_id: '', reason: '', start_date: '', end_date: '', success_criteria: '' }); fetch()
    } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }
  const createDev = async () => {
    setSaving(true)
    try {
      const payload: Record<string, any> = { employee_id: devForm.employee_id, title: devForm.title }
      if (devForm.description) payload.description = devForm.description
      if (devForm.career_goal) payload.career_goal = devForm.career_goal
      await api.post('/performance/development-plans', payload)
      setShowDev(false); setDevForm({ employee_id: '', title: '', description: '', career_goal: '' }); fetch()
    } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }

  return (
    <div className="grid gap-6 lg:grid-cols-2">
      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="font-semibold text-slate-900">Planes de Mejora</h2>
          <Button size="sm" onClick={() => setShowImp(true)}><Plus size={14} /> Nuevo</Button>
        </div>
        <div className="space-y-3">
          {loading ? <p className="text-sm text-slate-500">Cargando...</p> : improvement.length === 0 ? <p className="text-sm text-slate-500">Sin planes</p> : improvement.map((p: any) => (
            <Card key={p.id}><CardContent className="p-4">
              <div className="flex items-center justify-between mb-2">
                <p className="font-medium text-slate-900">{p.employee_id}</p>
                {pill(p.status)}
              </div>
              <p className="text-sm text-slate-600">{p.reason}</p>
              <p className="text-xs text-slate-500 mt-1">{fmt(p.start_date)} → {fmt(p.end_date)}</p>
            </CardContent></Card>
          ))}
        </div>
      </div>
      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="font-semibold text-slate-900">Planes de Desarrollo</h2>
          <Button size="sm" onClick={() => setShowDev(true)}><Plus size={14} /> Nuevo</Button>
        </div>
        <div className="space-y-3">
          {loading ? <p className="text-sm text-slate-500">Cargando...</p> : development.length === 0 ? <p className="text-sm text-slate-500">Sin planes</p> : development.map((p: any) => (
            <Card key={p.id}><CardContent className="p-4">
              <div className="flex items-center justify-between mb-2">
                <p className="font-medium text-slate-900">{p.title}</p>
                {pill(p.status)}
              </div>
              <p className="text-sm text-slate-600">{p.description || '-'}</p>
              {p.career_goal && <p className="text-xs text-brand-700 mt-1">Meta: {p.career_goal}</p>}
              <p className="text-xs text-slate-500 mt-1">Empleado: {p.employee_id}</p>
            </CardContent></Card>
          ))}
        </div>
      </div>
      <Dialog open={showImp} onOpenChange={setShowImp}>
        <DialogContent>
          <DialogHeader><DialogTitle>Nuevo Plan de Mejora</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2"><Label>Empleado *</Label><Select options={employees} placeholder="Seleccionar..." value={impForm.employee_id} onChange={e => setImpForm({ ...impForm, employee_id: e.target.value })} /></div>
            <div className="space-y-2"><Label>Motivo *</Label><Input value={impForm.reason} onChange={e => setImpForm({ ...impForm, reason: e.target.value })} /></div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Inicio *</Label><Input type="date" value={impForm.start_date} onChange={e => setImpForm({ ...impForm, start_date: e.target.value })} /></div>
              <div className="space-y-2"><Label>Fin *</Label><Input type="date" value={impForm.end_date} onChange={e => setImpForm({ ...impForm, end_date: e.target.value })} /></div>
            </div>
            <div className="space-y-2"><Label>Criterios de éxito</Label><Input value={impForm.success_criteria} onChange={e => setImpForm({ ...impForm, success_criteria: e.target.value })} /></div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowImp(false)}>Cancelar</Button><Button onClick={createImp} disabled={saving || !impForm.employee_id || !impForm.reason || !impForm.start_date || !impForm.end_date}>{saving ? 'Guardando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>
      <Dialog open={showDev} onOpenChange={setShowDev}>
        <DialogContent>
          <DialogHeader><DialogTitle>Nuevo Plan de Desarrollo</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2"><Label>Empleado *</Label><Select options={employees} placeholder="Seleccionar..." value={devForm.employee_id} onChange={e => setDevForm({ ...devForm, employee_id: e.target.value })} /></div>
            <div className="space-y-2"><Label>Título *</Label><Input value={devForm.title} onChange={e => setDevForm({ ...devForm, title: e.target.value })} /></div>
            <div className="space-y-2"><Label>Descripción</Label><Input value={devForm.description} onChange={e => setDevForm({ ...devForm, description: e.target.value })} /></div>
            <div className="space-y-2"><Label>Meta de carrera</Label><Input value={devForm.career_goal} onChange={e => setDevForm({ ...devForm, career_goal: e.target.value })} /></div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowDev(false)}>Cancelar</Button><Button onClick={createDev} disabled={saving || !devForm.employee_id || !devForm.title}>{saving ? 'Guardando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
