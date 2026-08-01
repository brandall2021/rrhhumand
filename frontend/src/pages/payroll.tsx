import { useEffect, useState } from 'react'
import { Plus } from 'lucide-react'
import api from '@/lib/api'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select } from '@/components/ui/select'
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter,
} from '@/components/ui/dialog'

const money = (n: any) => (n == null ? '-' : `$${Number(n).toLocaleString('es-AR', { minimumFractionDigits: 2 })}`)
const fmt = (d?: string) => (d ? new Date(d).toLocaleDateString('es-AR') : '-')
const months = ['Enero', 'Febrero', 'Marzo', 'Abril', 'Mayo', 'Junio', 'Julio', 'Agosto', 'Septiembre', 'Octubre', 'Noviembre', 'Diciembre']
const statusPill = (s: string) => (
  <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${s === 'CLOSED' || s === 'APPROVED' ? 'bg-emerald-50 text-emerald-700' : s === 'DRAFT' ? 'bg-slate-100 text-slate-600' : s === 'CALCULATED' || s === 'VALIDATED' ? 'bg-sky-50 text-sky-700' : 'bg-amber-50 text-amber-700'}`}>{s}</span>
)
const raw = (r: any) => (Array.isArray(r.data) ? r.data : Array.isArray(r) ? r : r.data?.data ?? [])

export default function PayrollPage() {
  const [activeTab, setActiveTab] = useState('periods')
  const [stats, setStats] = useState<any>(null)
  useEffect(() => { api.get('/payroll/dashboard').then(r => setStats(r.data)).catch(() => {}) }, [])
  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-900 mb-6">Nómina</h1>
      {stats && (
        <div className="grid grid-cols-2 md:grid-cols-4 gap-3 mb-6">
          {[['Períodos activos', stats.active_periods], ['Liquidaciones pendientes', stats.pending_runs], ['Errores', stats.total_errors], ['Errores bloqueantes', stats.blocking_errors]].map(([l, v]) => (
            <Card key={l}><CardContent className="p-4"><p className="text-xs text-slate-500 mb-1">{l}</p><p className="text-lg font-bold text-slate-900">{v}</p></CardContent></Card>
          ))}
        </div>
      )}
      <div className="flex gap-1 mb-6 border-b border-slate-200">
        {[['periods', 'Períodos'], ['concepts', 'Conceptos'], ['rules', 'Reglas'], ['novelties', 'Novedades'], ['advances', 'Adelantos']].map(([k, l]) => (
          <button key={k} onClick={() => setActiveTab(k)} className={`px-4 py-2.5 text-sm font-medium border-b-2 transition-colors ${activeTab === k ? 'border-brand-600 text-brand-700' : 'border-transparent text-slate-500 hover:text-slate-700'}`}>{l}</button>
        ))}
      </div>
      {activeTab === 'periods' && <PeriodsTab />}
      {activeTab === 'concepts' && <ConceptsTab />}
      {activeTab === 'rules' && <RulesTab />}
      {activeTab === 'novelties' && <NoveltiesTab />}
      {activeTab === 'advances' && <AdvancesTab />}
    </div>
  )
}

function PeriodsTab() {
  const [periods, setPeriods] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [selected, setSelected] = useState<any>(null)
  const [runs, setRuns] = useState<any[]>([])
  const [runLoading, setRunLoading] = useState(false)
  const [showCreate, setShowCreate] = useState(false)
  const [showRun, setShowRun] = useState(false)
  const [saving, setSaving] = useState(false)
  const [runType, setRunType] = useState('NORMAL')
  const [showDetail, setShowDetail] = useState(false)
  const [detail, setDetail] = useState<any>(null)
  const [emps, setEmps] = useState<any[]>([])
  const [summary, setSummary] = useState<any>(null)
  const [runErrors, setRunErrors] = useState<any[]>([])
  const [showAddEmp, setShowAddEmp] = useState(false)
  const [empOptions, setEmpOptions] = useState<{ value: string; label: string }[]>([])
  const [empSel, setEmpSel] = useState('')
  const [form, setForm] = useState({ name: '', year: String(new Date().getFullYear()), month: String(new Date().getMonth() + 1), period_type: 'MENSUAL', start_date: '', end_date: '', payment_date: '' })

  const fetch = async () => {
    setLoading(true)
    try { const r = await api.get('/payroll/periods', { params: { limit: '100' } }); setPeriods(raw(r)) } catch { setPeriods([]) } finally { setLoading(false) }
  }
  useEffect(() => { fetch() }, [])

  const selectPeriod = async (p: any) => {
    setSelected(p)
    setRunLoading(true)
    try { const r = await api.get('/payroll/runs', { params: { period_id: p.id, limit: '100' } }); setRuns(raw(r)) } catch { setRuns([]) } finally { setRunLoading(false) }
  }

  const create = async () => {
    setSaving(true)
    try {
      await api.post('/payroll/periods', { year: parseInt(form.year), month: parseInt(form.month), period_type: form.period_type, name: form.name, start_date: new Date(form.start_date).toISOString(), end_date: new Date(form.end_date).toISOString(), ...(form.payment_date ? { payment_date: new Date(form.payment_date).toISOString() } : {}) })
      setShowCreate(false); setForm({ name: '', year: form.year, month: form.month, period_type: 'MENSUAL', start_date: '', end_date: '', payment_date: '' }); fetch()
    } catch (e: any) { alert(e?.response?.data?.error || 'Error al crear período') } finally { setSaving(false) }
  }

  const createRun = async () => {
    if (!selected) return
    setSaving(true)
    try { await api.post(`/payroll/periods/${selected.id}/runs`, { run_type: runType }); setShowRun(false); selectPeriod(selected) } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }

  const runAction = async (run: any, action: string) => {
    try { await api.post(`/payroll/runs/${run.id}/${action}`); if (selected) selectPeriod(selected) } catch (e: any) { alert(e?.response?.data?.error || 'Error') }
  }

  const openDetail = async (run: any) => {
    setDetail(run); setShowDetail(true)
    try {
      const [e, s, er] = await Promise.all([api.get(`/payroll/runs/${run.id}/employees`), api.get(`/payroll/runs/${run.id}/summary`), api.get(`/payroll/runs/${run.id}/errors`)])
      setEmps(raw(e)); setSummary(s.data ?? null); setRunErrors(raw(er))
    } catch { setEmps([]); setSummary(null); setRunErrors([]) }
  }

  const openAddEmp = async () => {
    setEmpSel(''); setShowAddEmp(true)
    try { const r = await api.get('/employees', { params: { limit: '200' } }); setEmpOptions((r.data.data ?? []).map((x: any) => ({ value: x.id, label: `${x.first_name} ${x.last_name}` }))) } catch { setEmpOptions([]) }
  }

  const addEmp = async () => {
    if (!detail || !empSel) return
    try { await api.post(`/payroll/runs/${detail.id}/employees`, { employee_id: empSel }); setShowAddEmp(false); openDetail(detail) } catch (e: any) { alert(e?.response?.data?.error || 'Error') }
  }

  const closePeriod = async () => {
    if (!selected || !window.confirm(`¿Cerrar el período "${selected.name}"?`)) return
    try { await api.post(`/payroll/periods/${selected.id}/close`); fetch(); selectPeriod(selected) } catch (e: any) { alert(e?.response?.data?.error || 'Error') }
  }

  const TabButton = ({ label, onClick }: any) => <button onClick={onClick} className="text-left w-full px-4 py-3 border-b border-slate-100 hover:bg-slate-50 cursor-pointer">
    <p className="text-sm font-medium text-slate-900">{label.name}</p>
    <p className="text-xs text-slate-500">{months[label.month - 1]} {label.year} · {fmt(label.start_date)} → {fmt(label.end_date)}</p>
  </button>

  return (
    <div className="grid gap-6 lg:grid-cols-[1fr_1.4fr]">
      <div>
        <div className="flex items-center justify-between mb-4">
          <p className="text-sm text-slate-500">{periods.length} períodos</p>
          <Button size="sm" onClick={() => setShowCreate(true)}><Plus size={14} /> Nuevo Período</Button>
        </div>
        <Card>
          <CardContent className="p-0">
            {loading ? <div className="p-6 text-center text-slate-500">Cargando...</div> : periods.length === 0 ? <div className="p-6 text-center text-slate-500">No hay períodos</div> : (
              <div>{periods.map((p: any) => (
                <div key={p.id} onClick={() => selectPeriod(p)} className={`cursor-pointer ${selected?.id === p.id ? 'bg-brand-50/60' : ''}`}>
                  <TabButton label={p} />
                </div>
              ))}</div>
            )}
          </CardContent>
        </Card>
      </div>
      <div>
        {!selected ? <Card><CardContent className="p-8 text-center text-slate-500">Seleccioná un período para ver sus liquidaciones</CardContent></Card> : (
          <div>
            <div className="flex items-center justify-between mb-4">
              <div>
                <h2 className="font-bold text-slate-900">{selected.name}</h2>
                <p className="text-sm text-slate-500">{months[selected.month - 1]} {selected.year} · {fmt(selected.start_date)} → {fmt(selected.end_date)}</p>
              </div>
              <div className="flex gap-2">
                {selected.status !== 'CLOSED' && <Button size="sm" variant="outline" onClick={closePeriod}>Cerrar Período</Button>}
                <Button size="sm" onClick={() => { setRunType('NORMAL'); setShowRun(true) }}><Plus size={14} /> Nueva Liquidación</Button>
              </div>
            </div>
            <Card>
              <CardContent className="p-0">
                {runLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div> : runs.length === 0 ? <div className="p-6 text-center text-slate-500">Sin liquidaciones</div> : (
                  <table className="w-full text-sm">
                    <thead><tr className="border-b border-slate-200 bg-slate-50">
                      <th className="text-left px-4 py-3 font-medium text-slate-600">Nº</th>
                      <th className="text-left px-4 py-3 font-medium text-slate-600">Tipo</th>
                      <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                      <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                    </tr></thead>
                    <tbody>{runs.map((r: any) => (
                      <tr key={r.id} className="border-b border-slate-100 hover:bg-slate-50">
                        <td className="px-4 py-3 font-semibold text-slate-900">#{r.run_number}</td>
                        <td className="px-4 py-3 text-slate-600">{r.run_type}</td>
                        <td className="px-4 py-3">{statusPill(r.status)}</td>
                        <td className="px-4 py-3">
                          <div className="flex justify-end gap-1">
                            <Button size="sm" variant="outline" onClick={() => openDetail(r)}>Detalle</Button>
                            {r.status === 'DRAFT' && <><Button size="sm" variant="outline" onClick={() => runAction(r, 'calculate')}>Calcular</Button><Button size="sm" variant="ghost" onClick={() => runAction(r, 'close')}>Cerrar</Button></>}
                            {r.status === 'CALCULATED' && <><Button size="sm" variant="outline" onClick={() => runAction(r, 'validate')}>Validar</Button><Button size="sm" variant="ghost" onClick={() => runAction(r, 'close')}>Cerrar</Button></>}
                            {r.status === 'VALIDATED' && <><Button size="sm" variant="outline" onClick={() => runAction(r, 'approve')}>Aprobar</Button><Button size="sm" variant="ghost" onClick={() => runAction(r, 'close')}>Cerrar</Button></>}
                            {r.status === 'APPROVED' && <Button size="sm" variant="outline" onClick={() => runAction(r, 'close')}>Cerrar</Button>}
                          </div>
                        </td>
                      </tr>
                    ))}</tbody>
                  </table>
                )}
              </CardContent>
            </Card>
          </div>
        )}
      </div>

      <Dialog open={showCreate} onOpenChange={setShowCreate}>
        <DialogContent>
          <DialogHeader><DialogTitle>Nuevo Período</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Año *</Label><Input type="number" value={form.year} onChange={e => setForm({ ...form, year: e.target.value })} /></div>
              <div className="space-y-2"><Label>Mes *</Label><Select options={months.map((m, i) => ({ value: String(i + 1), label: m }))} value={form.month} onChange={e => setForm({ ...form, month: e.target.value })} /></div>
            </div>
            <div className="space-y-2"><Label>Tipo *</Label><Select options={['MENSUAL', 'QUINCENAL', 'SEMANAL'].map(t => ({ value: t, label: t }))} value={form.period_type} onChange={e => setForm({ ...form, period_type: e.target.value })} /></div>
            <div className="space-y-2"><Label>Nombre *</Label><Input value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} placeholder="Ej: Mensual Febrero 2026" /></div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Inicio *</Label><Input type="date" value={form.start_date} onChange={e => setForm({ ...form, start_date: e.target.value })} /></div>
              <div className="space-y-2"><Label>Fin *</Label><Input type="date" value={form.end_date} onChange={e => setForm({ ...form, end_date: e.target.value })} /></div>
            </div>
            <div className="space-y-2"><Label>Fecha de pago</Label><Input type="date" value={form.payment_date} onChange={e => setForm({ ...form, payment_date: e.target.value })} /></div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowCreate(false)}>Cancelar</Button>
            <Button onClick={create} disabled={saving || !form.name || !form.start_date || !form.end_date}>{saving ? 'Creando...' : 'Crear'}</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={showRun} onOpenChange={setShowRun}>
        <DialogContent>
          <DialogHeader><DialogTitle>Nueva Liquidación</DialogTitle><DialogDescription>Período: {selected?.name}</DialogDescription></DialogHeader>
          <div className="space-y-2"><Label>Tipo *</Label><Select options={['NORMAL', 'SUPLEMENTARIA', 'VACACIONES', 'AGUINALDO'].map(t => ({ value: t, label: t }))} value={runType} onChange={e => setRunType(e.target.value)} /></div>
          <DialogFooter><Button variant="outline" onClick={() => setShowRun(false)}>Cancelar</Button><Button onClick={createRun} disabled={saving}>{saving ? 'Creando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={showDetail} onOpenChange={setShowDetail}>
        <DialogContent className="max-w-3xl">
          <DialogHeader>
            <DialogTitle>Liquidación #{detail?.run_number} · {detail?.run_type}</DialogTitle>
            {summary && <div className="grid grid-cols-2 md:grid-cols-4 gap-2 text-sm mt-2">
              <div className="bg-slate-50 rounded-lg p-2"><p className="text-xs text-slate-500">Empleados</p><p className="font-semibold">{summary.total_employees}</p></div>
              <div className="bg-slate-50 rounded-lg p-2"><p className="text-xs text-slate-500">Total neto</p><p className="font-semibold">{money(summary.total_net)}</p></div>
              <div className="bg-slate-50 rounded-lg p-2"><p className="text-xs text-slate-500">Deducciones</p><p className="font-semibold">{money(summary.total_deductions)}</p></div>
              <div className="bg-slate-50 rounded-lg p-2"><p className="text-xs text-slate-500">Costo empleador</p><p className="font-semibold">{money(summary.total_employer_cost)}</p></div>
            </div>}
          </DialogHeader>
          {runErrors.length > 0 && (
            <div className="bg-rose-50 border border-rose-200 rounded-lg p-3 mb-3">
              <p className="text-sm font-medium text-rose-700 mb-1">{runErrors.length} error(es) de cálculo</p>
              {runErrors.slice(0, 5).map((e: any) => <p key={e.id} className="text-xs text-rose-600">[{e.severity} · {e.code}] {e.message}</p>)}
            </div>
          )}
          <div className="flex justify-end mb-2"><Button size="sm" variant="outline" onClick={openAddEmp}><Plus size={14} /> Agregar Empleado</Button></div>
          <div className="max-h-96 overflow-auto">
            <table className="w-full text-sm">
              <thead className="sticky top-0 bg-white"><tr className="border-b border-slate-200 bg-slate-50">
                <th className="text-left px-4 py-3 font-medium text-slate-600">Empleado</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                <th className="text-right px-4 py-3 font-medium text-slate-600">Bruto</th>
                <th className="text-right px-4 py-3 font-medium text-slate-600">Neto</th>
              </tr></thead>
              <tbody>{emps.map((e: any) => (
                <tr key={e.id} className="border-b border-slate-100">
                  <td className="px-4 py-3 font-medium text-slate-900">{e.employee_id}</td>
                  <td className="px-4 py-3"><span className={`text-xs font-medium ${e.error_message ? 'text-rose-600' : 'text-slate-600'}`}>{e.status}{e.error_message ? ' · error' : ''}</span></td>
                  <td className="px-4 py-3 text-right text-slate-600">{money(e.gross_remunerative)}</td>
                  <td className="px-4 py-3 text-right font-semibold text-slate-900">{money(e.net_amount)}</td>
                </tr>
              ))}</tbody>
            </table>
          </div>
        </DialogContent>
      </Dialog>

      <Dialog open={showAddEmp} onOpenChange={setShowAddEmp}>
        <DialogContent>
          <DialogHeader><DialogTitle>Agregar Empleado</DialogTitle></DialogHeader>
          <div className="space-y-2"><Label>Empleado *</Label><Select options={empOptions} placeholder="Seleccionar..." value={empSel} onChange={e => setEmpSel(e.target.value)} /></div>
          <DialogFooter><Button variant="outline" onClick={() => setShowAddEmp(false)}>Cancelar</Button><Button onClick={addEmp} disabled={!empSel}>Agregar</Button></DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function ConceptsTab() {
  const [concepts, setConcepts] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [filter, setFilter] = useState({ concept_type: '', taxability: '' })
  const [showModal, setShowModal] = useState(false)
  const [saving, setSaving] = useState(false)
  const empty = { code: '', name: '', concept_type: 'REMUNERATIVO', taxability: 'GRAVADO', calculation_type: 'FIJO', description: '', sort_order: '0' }
  const [form, setForm] = useState(empty)

  const fetch = async () => {
    setLoading(true)
    try { const params: Record<string, any> = {}; if (filter.concept_type) params.concept_type = filter.concept_type; if (filter.taxability) params.taxability = filter.taxability; const r = await api.get('/payroll/concepts', { params }); setConcepts(raw(r)) } catch { setConcepts([]) } finally { setLoading(false) }
  }
  useEffect(() => { fetch() }, [filter])

  const create = async () => {
    setSaving(true)
    try {
      await api.post('/payroll/concepts', { code: form.code, name: form.name, concept_type: form.concept_type, taxability: form.taxability, calculation_type: form.calculation_type, ...(form.description ? { description: form.description } : {}), sort_order: parseInt(form.sort_order) || 0 })
      setShowModal(false); setForm(empty); fetch()
    } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4 flex-wrap gap-2">
        <div className="flex items-end gap-3">
          <div><Label>Tipo</Label><Select options={[{ value: '', label: 'Todos' }, ...['REMUNERATIVO', 'NO_REMUNERATIVO', 'DEDUCCION', 'APORTE', 'CONTRIBUCION'].map(t => ({ value: t, label: t }))]} value={filter.concept_type} onChange={e => setFilter({ ...filter, concept_type: e.target.value })} /></div>
          <div><Label>Tratamiento fiscal</Label><Select options={[{ value: '', label: 'Todos' }, ...['GRAVADO', 'NO_GRAVADO', 'EXENTO'].map(t => ({ value: t, label: t }))]} value={filter.taxability} onChange={e => setFilter({ ...filter, taxability: e.target.value })} /></div>
        </div>
        <Button size="sm" onClick={() => setShowModal(true)}><Plus size={14} /> Nuevo Concepto</Button>
      </div>
      <Card>
        <CardContent className="p-0">
          {loading ? <div className="p-6 text-center text-slate-500">Cargando...</div> : concepts.length === 0 ? <div className="p-6 text-center text-slate-500">No hay conceptos</div> : (
            <table className="w-full text-sm">
              <thead><tr className="border-b border-slate-200 bg-slate-50">
                <th className="text-left px-4 py-3 font-medium text-slate-600">Código</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Tipo</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Tratamiento</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Cálculo</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
              </tr></thead>
              <tbody>{concepts.map((c: any) => (
                <tr key={c.id} className="border-b border-slate-100 hover:bg-slate-50">
                  <td className="px-4 py-3 font-mono text-slate-600">{c.code}</td>
                  <td className="px-4 py-3 font-medium text-slate-900">{c.name}</td>
                  <td className="px-4 py-3 text-slate-600">{c.concept_type}</td>
                  <td className="px-4 py-3 text-slate-600">{c.taxability}</td>
                  <td className="px-4 py-3 text-slate-600">{c.calculation_type}</td>
                  <td className="px-4 py-3">{c.active ? <span className="text-xs font-medium text-emerald-700 bg-emerald-50 px-2 py-0.5 rounded-full">ACTIVO</span> : <span className="text-xs font-medium text-slate-500 bg-slate-100 px-2 py-0.5 rounded-full">INACTIVO</span>}</td>
                </tr>
              ))}</tbody>
            </table>
          )}
        </CardContent>
      </Card>
      <Dialog open={showModal} onOpenChange={setShowModal}>
        <DialogContent>
          <DialogHeader><DialogTitle>Nuevo Concepto</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Código *</Label><Input value={form.code} onChange={e => setForm({ ...form, code: e.target.value })} /></div>
              <div className="space-y-2"><Label>Nombre *</Label><Input value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} /></div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Tipo *</Label><Select options={['REMUNERATIVO', 'NO_REMUNERATIVO', 'DEDUCCION', 'APORTE', 'CONTRIBUCION', 'HABER'].map(t => ({ value: t, label: t }))} value={form.concept_type} onChange={e => setForm({ ...form, concept_type: e.target.value })} /></div>
              <div className="space-y-2"><Label>Tratamiento *</Label><Select options={['GRAVADO', 'NO_GRAVADO', 'EXENTO'].map(t => ({ value: t, label: t }))} value={form.taxability} onChange={e => setForm({ ...form, taxability: e.target.value })} /></div>
            </div>
            <div className="space-y-2"><Label>Método de cálculo *</Label><Select options={['FIJO', 'PORCENTAJE', 'HORAS', 'DIAS', 'FORMULA'].map(t => ({ value: t, label: t }))} value={form.calculation_type} onChange={e => setForm({ ...form, calculation_type: e.target.value })} /></div>
            <div className="space-y-2"><Label>Descripción</Label><Input value={form.description} onChange={e => setForm({ ...form, description: e.target.value })} /></div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button><Button onClick={create} disabled={saving || !form.code || !form.name}>{saving ? 'Guardando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function RulesTab() {
  const [rules, setRules] = useState<any[]>([])
  const [concepts, setConcepts] = useState<{ value: string; label: string }[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [saving, setSaving] = useState(false)
  const [form, setForm] = useState({ concept_id: '', rule_type: 'CALCULO', formula: '', effective_from: '', priority: '0', parameters: '' })

  const fetch = async () => {
    setLoading(true)
    try { const r = await api.get('/payroll/rules'); setRules(raw(r)) } catch { setRules([]) } finally { setLoading(false) }
  }
  useEffect(() => {
    fetch()
    api.get('/payroll/concepts').then(r => setConcepts((raw(r)).map((c: any) => ({ value: c.id, label: `${c.code} · ${c.name}` })))).catch(() => {})
  }, [])

  const create = async () => {
    setSaving(true)
    try {
      const payload: Record<string, any> = { concept_id: form.concept_id, rule_type: form.rule_type, effective_from: form.effective_from, priority: parseInt(form.priority) || 0 }
      if (form.formula) payload.formula = form.formula
      if (form.parameters) { try { payload.parameters = JSON.parse(form.parameters) } catch { payload.parameters = {} } }
      await api.post('/payroll/rules', payload)
      setShowModal(false); setForm({ concept_id: '', rule_type: 'CALCULO', formula: '', effective_from: '', priority: '0', parameters: '' }); fetch()
    } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-slate-500">{rules.length} reglas</p>
        <Button size="sm" onClick={() => setShowModal(true)}><Plus size={14} /> Nueva Regla</Button>
      </div>
      <Card>
        <CardContent className="p-0">
          {loading ? <div className="p-6 text-center text-slate-500">Cargando...</div> : rules.length === 0 ? <div className="p-6 text-center text-slate-500">No hay reglas</div> : (
            <table className="w-full text-sm">
              <thead><tr className="border-b border-slate-200 bg-slate-50">
                <th className="text-left px-4 py-3 font-medium text-slate-600">Concepto</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Tipo</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Fórmula</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Prioridad</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Vigente desde</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
              </tr></thead>
              <tbody>{rules.map((r: any) => (
                <tr key={r.id} className="border-b border-slate-100 hover:bg-slate-50">
                  <td className="px-4 py-3 font-mono text-slate-600">{r.concept_id}</td>
                  <td className="px-4 py-3 text-slate-600">{r.rule_type}</td>
                  <td className="px-4 py-3 font-mono text-xs text-slate-600 max-w-[200px] truncate">{r.formula || '-'}</td>
                  <td className="px-4 py-3 text-slate-600">{r.priority}</td>
                  <td className="px-4 py-3 text-slate-600">{fmt(r.effective_from)}</td>
                  <td className="px-4 py-3">{r.active ? <span className="text-xs font-medium text-emerald-700 bg-emerald-50 px-2 py-0.5 rounded-full">ACTIVA</span> : <span className="text-xs font-medium text-slate-500 bg-slate-100 px-2 py-0.5 rounded-full">INACTIVA</span>}</td>
                </tr>
              ))}</tbody>
            </table>
          )}
        </CardContent>
      </Card>
      <Dialog open={showModal} onOpenChange={setShowModal}>
        <DialogContent>
          <DialogHeader><DialogTitle>Nueva Regla</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2"><Label>Concepto *</Label><Select options={concepts} placeholder="Seleccionar..." value={form.concept_id} onChange={e => setForm({ ...form, concept_id: e.target.value })} /></div>
            <div className="space-y-2"><Label>Tipo *</Label><Select options={['CALCULO', 'LIMITE', 'APORTE', 'REDONDEO', 'GARANTIA'].map(t => ({ value: t, label: t }))} value={form.rule_type} onChange={e => setForm({ ...form, rule_type: e.target.value })} /></div>
            <div className="space-y-2"><Label>Fórmula</Label><Input value={form.formula} onChange={e => setForm({ ...form, formula: e.target.value })} placeholder="Ej: base * rate" /></div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Vigente desde *</Label><Input type="date" value={form.effective_from} onChange={e => setForm({ ...form, effective_from: e.target.value })} /></div>
              <div className="space-y-2"><Label>Prioridad</Label><Input type="number" value={form.priority} onChange={e => setForm({ ...form, priority: e.target.value })} /></div>
            </div>
            <div className="space-y-2"><Label>Parámetros (JSON)</Label><Input value={form.parameters} onChange={e => setForm({ ...form, parameters: e.target.value })} placeholder='{"rate": 0.17}' /></div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button><Button onClick={create} disabled={saving || !form.concept_id || !form.effective_from}>{saving ? 'Guardando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function NoveltiesTab() {
  const [novelties, setNovelties] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [employees, setEmployees] = useState<{ value: string; label: string }[]>([])
  const [showModal, setShowModal] = useState(false)
  const [saving, setSaving] = useState(false)
  const empty = { employee_id: '', novelty_type: 'INGRESO', quantity: '', amount: '', description: '' }
  const [form, setForm] = useState(empty)

  const fetch = async () => {
    setLoading(true)
    try { const r = await api.get('/payroll/novelties', { params: { limit: '100' } }); setNovelties(raw(r)) } catch { setNovelties([]) } finally { setLoading(false) }
  }
  useEffect(() => {
    fetch()
    api.get('/employees', { params: { limit: '200' } }).then(r => setEmployees((r.data.data ?? []).map((e: any) => ({ value: e.id, label: `${e.first_name} ${e.last_name}` })))).catch(() => {})
  }, [])

  const create = async () => {
    setSaving(true)
    try {
      const payload: Record<string, any> = { employee_id: form.employee_id, novelty_type: form.novelty_type, source: 'MANUAL' }
      if (form.quantity) payload.quantity = parseFloat(form.quantity)
      if (form.amount) payload.amount = parseFloat(form.amount)
      if (form.description) payload.description = form.description
      await api.post('/payroll/novelties', payload)
      setShowModal(false); setForm(empty); fetch()
    } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }

  const approve = async (id: string) => { try { await api.post(`/payroll/novelties/${id}/approve`); fetch() } catch (e: any) { alert(e?.response?.data?.error || 'Error') } }
  const remove = async (id: string) => { if (!window.confirm('¿Eliminar la novedad?')) return; try { await api.delete(`/payroll/novelties/${id}`); fetch() } catch (e: any) { alert(e?.response?.data?.error || 'Error') } }

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-slate-500">{novelties.length} novedades</p>
        <Button size="sm" onClick={() => setShowModal(true)}><Plus size={14} /> Nueva Novedad</Button>
      </div>
      <Card>
        <CardContent className="p-0">
          {loading ? <div className="p-6 text-center text-slate-500">Cargando...</div> : novelties.length === 0 ? <div className="p-6 text-center text-slate-500">No hay novedades</div> : (
            <table className="w-full text-sm">
              <thead><tr className="border-b border-slate-200 bg-slate-50">
                <th className="text-left px-4 py-3 font-medium text-slate-600">Empleado</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Tipo</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Cantidad</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Importe</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Origen</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
              </tr></thead>
              <tbody>{novelties.map((n: any) => (
                <tr key={n.id} className="border-b border-slate-100 hover:bg-slate-50">
                  <td className="px-4 py-3 font-medium text-slate-900">{n.employee_id}</td>
                  <td className="px-4 py-3 text-slate-600">{n.novelty_type}</td>
                  <td className="px-4 py-3 text-slate-600">{n.quantity ?? '-'}</td>
                  <td className="px-4 py-3 text-slate-600">{money(n.amount)}</td>
                  <td className="px-4 py-3 text-slate-600">{n.source}</td>
                  <td className="px-4 py-3">{statusPill(n.status)}</td>
                  <td className="px-4 py-3">
                    <div className="flex justify-end gap-1">
                      {n.status === 'PENDING' && <Button size="sm" variant="outline" onClick={() => approve(n.id)}>Aprobar</Button>}
                      <Button size="sm" variant="ghost" onClick={() => remove(n.id)}>Eliminar</Button>
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
          <DialogHeader><DialogTitle>Nueva Novedad</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2"><Label>Empleado *</Label><Select options={employees} placeholder="Seleccionar..." value={form.employee_id} onChange={e => setForm({ ...form, employee_id: e.target.value })} /></div>
            <div className="space-y-2"><Label>Tipo *</Label><Select options={['INGRESO', 'EGRESO', 'AUSENCIA', 'LICENCIA', 'ADICIONAL', 'DESCUENTO', 'AJUSTE', 'OTRO'].map(t => ({ value: t, label: t }))} value={form.novelty_type} onChange={e => setForm({ ...form, novelty_type: e.target.value })} /></div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Cantidad</Label><Input type="number" value={form.quantity} onChange={e => setForm({ ...form, quantity: e.target.value })} /></div>
              <div className="space-y-2"><Label>Importe</Label><Input type="number" value={form.amount} onChange={e => setForm({ ...form, amount: e.target.value })} /></div>
            </div>
            <div className="space-y-2"><Label>Descripción</Label><Input value={form.description} onChange={e => setForm({ ...form, description: e.target.value })} /></div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button><Button onClick={create} disabled={saving || !form.employee_id}>{saving ? 'Guardando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function AdvancesTab() {
  const [employees, setEmployees] = useState<{ value: string; label: string }[]>([])
  const [selected, setSelected] = useState('')
  const [advances, setAdvances] = useState<any[]>([])
  const [showModal, setShowModal] = useState(false)
  const [saving, setSaving] = useState(false)
  const [form, setForm] = useState({ amount: '', request_date: new Date().toISOString().slice(0, 10), installments: '1', reason: '' })

  const fetch = async () => {
    if (!selected) { setAdvances([]); return }
    try { const r = await api.get(`/payroll/employees/${selected}/advances`); setAdvances(raw(r)) } catch { setAdvances([]) }
  }
  useEffect(() => {
    api.get('/employees', { params: { limit: '200' } }).then(r => setEmployees((r.data.data ?? []).map((e: any) => ({ value: e.id, label: `${e.first_name} ${e.last_name}` })))).catch(() => {})
  }, [])
  useEffect(() => { fetch() }, [selected])

  const create = async () => {
    setSaving(true)
    try {
      await api.post('/payroll/advances', { employee_id: selected, amount: parseFloat(form.amount), request_date: form.request_date, installments: parseInt(form.installments) || 1, ...(form.reason ? { reason: form.reason } : {}) })
      setShowModal(false); setForm({ amount: '', request_date: new Date().toISOString().slice(0, 10), installments: '1', reason: '' }); fetch()
    } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4 flex-wrap gap-2">
        <div className="w-72"><Label>Empleado</Label><Select options={employees} placeholder="Seleccionar..." value={selected} onChange={e => setSelected(e.target.value)} /></div>
        <Button size="sm" onClick={() => setShowModal(true)} disabled={!selected}><Plus size={14} /> Nuevo Adelanto</Button>
      </div>
      <Card>
        <CardContent className="p-0">
          {!selected ? <div className="p-6 text-center text-slate-500">Seleccioná un empleado para ver sus adelantos</div> : advances.length === 0 ? <div className="p-6 text-center text-slate-500">Sin adelantos</div> : (
            <table className="w-full text-sm">
              <thead><tr className="border-b border-slate-200 bg-slate-50">
                <th className="text-left px-4 py-3 font-medium text-slate-600">Fecha</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Importe</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Cuotas</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Restante</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Motivo</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
              </tr></thead>
              <tbody>{advances.map((a: any) => (
                <tr key={a.id} className="border-b border-slate-100 hover:bg-slate-50">
                  <td className="px-4 py-3 text-slate-900">{fmt(a.request_date)}</td>
                  <td className="px-4 py-3 font-medium text-slate-900">{money(a.amount)}</td>
                  <td className="px-4 py-3 text-slate-600">{a.installments}</td>
                  <td className="px-4 py-3 text-slate-600">{money(a.remaining_amount)}</td>
                  <td className="px-4 py-3 text-slate-600">{a.reason || '-'}</td>
                  <td className="px-4 py-3">{statusPill(a.status)}</td>
                </tr>
              ))}</tbody>
            </table>
          )}
        </CardContent>
      </Card>
      <Dialog open={showModal} onOpenChange={setShowModal}>
        <DialogContent>
          <DialogHeader><DialogTitle>Nuevo Adelanto</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2"><Label>Importe *</Label><Input type="number" value={form.amount} onChange={e => setForm({ ...form, amount: e.target.value })} /></div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Fecha</Label><Input type="date" value={form.request_date} onChange={e => setForm({ ...form, request_date: e.target.value })} /></div>
              <div className="space-y-2"><Label>Cuotas</Label><Input type="number" value={form.installments} onChange={e => setForm({ ...form, installments: e.target.value })} /></div>
            </div>
            <div className="space-y-2"><Label>Motivo</Label><Input value={form.reason} onChange={e => setForm({ ...form, reason: e.target.value })} /></div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button><Button onClick={create} disabled={saving || !form.amount}>{saving ? 'Guardando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
