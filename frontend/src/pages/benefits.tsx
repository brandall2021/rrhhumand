import { useEffect, useState } from 'react'
import { Plus } from 'lucide-react'
import api from '@/lib/api'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select } from '@/components/ui/select'
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter,
} from '@/components/ui/dialog'

const money = (n: any) => (n == null ? '-' : `$${Number(n).toLocaleString('es-AR', { minimumFractionDigits: 2 })}`)
const fmt = (d?: string) => (d ? new Date(d).toLocaleDateString('es-AR') : '-')
const pill = (s: string) => (
  <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${s === 'ACTIVE' || s === 'ACTIVO' || s === 'APPROVED' ? 'bg-emerald-50 text-emerald-700' : s === 'PENDING' ? 'bg-amber-50 text-amber-700' : s === 'REJECTED' ? 'bg-rose-50 text-rose-700' : s === 'CANCELLED' ? 'bg-slate-100 text-slate-500' : 'bg-sky-50 text-sky-700'}`}>{s}</span>
)
const data = (r: any) => (Array.isArray(r.data?.data) ? r.data.data : Array.isArray(r.data) ? r.data : r.data?.data ?? [])

export default function BenefitsPage() {
  const [activeTab, setActiveTab] = useState('benefits')
  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-900 mb-6">Beneficios</h1>
      <div className="flex gap-1 mb-6 border-b border-slate-200">
        {[['benefits', 'Catálogo'], ['assignments', 'Asignaciones'], ['requests', 'Solicitudes'], ['catalog', 'Categorías / Tipos']].map(([k, l]) => (
          <button key={k} onClick={() => setActiveTab(k)} className={`px-4 py-2.5 text-sm font-medium border-b-2 transition-colors ${activeTab === k ? 'border-brand-600 text-brand-700' : 'border-transparent text-slate-500 hover:text-slate-700'}`}>{l}</button>
        ))}
      </div>
      {activeTab === 'benefits' && <BenefitsTab />}
      {activeTab === 'assignments' && <AssignmentsTab />}
      {activeTab === 'requests' && <RequestsTab />}
      {activeTab === 'catalog' && <CatalogTab />}
    </div>
  )
}

function BenefitsTab() {
  const [benefits, setBenefits] = useState<any[]>([])
  const [types, setTypes] = useState<{ value: string; label: string }[]>([])
  const [providers, setProviders] = useState<{ value: string; label: string }[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [saving, setSaving] = useState(false)
  const [filter, setFilter] = useState('')
  const [form, setForm] = useState({ type_id: '', provider_id: '', code: '', name: '', short_description: '' })

  const fetch = async () => {
    setLoading(true)
    try { const r = await api.get('/benefits'); setBenefits(data(r)) } catch { setBenefits([]) } finally { setLoading(false) }
  }
  useEffect(() => {
    fetch()
    api.get('/benefits/types').then(r => setTypes(data(r).map((t: any) => ({ value: t.id, label: t.name })))).catch(() => {})
    api.get('/benefits/providers').then(r => setProviders(data(r).map((p: any) => ({ value: p.id, label: p.name })))).catch(() => {})
  }, [])

  const create = async () => {
    setSaving(true)
    try {
      const payload: Record<string, any> = { type_id: form.type_id, code: form.code, name: form.name }
      if (form.provider_id) payload.provider_id = form.provider_id
      if (form.short_description) payload.short_description = form.short_description
      await api.post('/benefits', payload)
      setShowModal(false); setForm({ type_id: '', provider_id: '', code: '', name: '', short_description: '' }); fetch()
    } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }

  const visible = benefits.filter((b: any) => !filter || (b.name || '').toLowerCase().includes(filter.toLowerCase()) || (b.code || '').toLowerCase().includes(filter.toLowerCase()))

  return (
    <div>
      <div className="flex items-center justify-between mb-4 flex-wrap gap-2">
        <Input className="max-w-xs" placeholder="Buscar..." value={filter} onChange={e => setFilter(e.target.value)} />
        <Button size="sm" onClick={() => setShowModal(true)}><Plus size={14} /> Nuevo Beneficio</Button>
      </div>
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {loading ? <div className="p-6 text-center text-slate-500 col-span-full">Cargando...</div> : visible.length === 0 ? <div className="p-6 text-center text-slate-500 col-span-full">No hay beneficios</div> : visible.map((b: any) => (
          <Card key={b.id}>
            <CardContent className="p-5">
              <div className="flex items-start justify-between mb-2">
                <div>
                  <h3 className="font-semibold text-slate-900">{b.name}</h3>
                  <p className="text-xs text-slate-500 font-mono">{b.code}</p>
                </div>
                {pill(b.status)}
              </div>
              <p className="text-sm text-slate-600 mb-3 line-clamp-2">{b.short_description || b.description || 'Sin descripción'}</p>
              <div className="flex flex-wrap gap-2 text-xs">
                {b.auto_enroll && <span className="bg-violet-50 text-violet-700 px-2 py-0.5 rounded-full">Alta automática</span>}
                {b.requires_evidence && <span className="bg-amber-50 text-amber-700 px-2 py-0.5 rounded-full">Requiere comprobante</span>}
                <span className="bg-slate-100 text-slate-600 px-2 py-0.5 rounded-full">{b.visibility}</span>
                <span className="bg-slate-100 text-slate-600 px-2 py-0.5 rounded-full">{b.current_beneficiaries} inscritos</span>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
      <Dialog open={showModal} onOpenChange={setShowModal}>
        <DialogContent>
          <DialogHeader><DialogTitle>Nuevo Beneficio</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2"><Label>Tipo *</Label><Select options={types} placeholder="Seleccionar..." value={form.type_id} onChange={e => setForm({ ...form, type_id: e.target.value })} /></div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Código *</Label><Input value={form.code} onChange={e => setForm({ ...form, code: e.target.value })} /></div>
              <div className="space-y-2"><Label>Proveedor</Label><Select options={providers} placeholder="Opcional" value={form.provider_id} onChange={e => setForm({ ...form, provider_id: e.target.value })} /></div>
            </div>
            <div className="space-y-2"><Label>Nombre *</Label><Input value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} /></div>
            <div className="space-y-2"><Label>Descripción corta</Label><Input value={form.short_description} onChange={e => setForm({ ...form, short_description: e.target.value })} /></div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button><Button onClick={create} disabled={saving || !form.type_id || !form.code || !form.name}>{saving ? 'Guardando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function AssignmentsTab() {
  const [assignments, setAssignments] = useState<any[]>([])
  const [employees, setEmployees] = useState<{ value: string; label: string }[]>([])
  const [benefits, setBenefits] = useState<{ value: string; label: string }[]>([])
  const [loading, setLoading] = useState(true)
  const [filterEmp, setFilterEmp] = useState('')
  const [showModal, setShowModal] = useState(false)
  const [saving, setSaving] = useState(false)
  const [form, setForm] = useState({ employee_id: '', benefit_id: '' })

  const fetch = async () => {
    setLoading(true)
    try { const params: Record<string, any> = {}; if (filterEmp) params.employee_id = filterEmp; const r = await api.get('/benefits/assignments', { params }); setAssignments(data(r)) } catch { setAssignments([]) } finally { setLoading(false) }
  }
  useEffect(() => {
    fetch()
    api.get('/employees', { params: { limit: '200' } }).then(r => setEmployees((r.data.data ?? []).map((e: any) => ({ value: e.id, label: `${e.first_name} ${e.last_name}` })))).catch(() => {})
    api.get('/benefits').then(r => setBenefits(data(r).map((b: any) => ({ value: b.id, label: b.name })))).catch(() => {})
  }, [])
  useEffect(() => { fetch() }, [filterEmp])

  const enroll = async () => {
    setSaving(true)
    try {
      await api.post('/benefits/assignments', { employee_id: form.employee_id, benefit_id: form.benefit_id, source: 'MANUAL' })
      setShowModal(false); setForm({ employee_id: '', benefit_id: '' }); fetch()
    } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }

  const cancel = async (id: string) => {
    const reason = window.prompt('Motivo de la baja:')
    if (reason === null) return
    try { await api.post(`/benefits/assignments/${id}/cancel`, { reason }); fetch() } catch (e: any) { alert(e?.response?.data?.error || 'Error') }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4 flex-wrap gap-2">
        <div className="w-72"><Label>Empleado</Label><Select options={[{ value: '', label: 'Todos' }, ...employees]} value={filterEmp} onChange={e => setFilterEmp(e.target.value)} /></div>
        <Button size="sm" onClick={() => setShowModal(true)}><Plus size={14} /> Asignar Beneficio</Button>
      </div>
      <Card>
        <CardContent className="p-0">
          {loading ? <div className="p-6 text-center text-slate-500">Cargando...</div> : assignments.length === 0 ? <div className="p-6 text-center text-slate-500">Sin asignaciones</div> : (
            <table className="w-full text-sm">
              <thead><tr className="border-b border-slate-200 bg-slate-50">
                <th className="text-left px-4 py-3 font-medium text-slate-600">Beneficio</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Empleado</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Alta</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Costo</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
              </tr></thead>
              <tbody>{assignments.map((a: any) => (
                <tr key={a.id} className="border-b border-slate-100 hover:bg-slate-50">
                  <td className="px-4 py-3 font-medium text-slate-900">{a.benefit_id}</td>
                  <td className="px-4 py-3 text-slate-600">{a.employee_id}</td>
                  <td className="px-4 py-3 text-slate-600">{fmt(a.enrollment_date)}</td>
                  <td className="px-4 py-3 text-slate-600">{money(a.employee_cost)}</td>
                  <td className="px-4 py-3">{pill(a.status)}</td>
                  <td className="px-4 py-3">
                    <div className="flex justify-end gap-1">
                      {a.status === 'ACTIVE' && <Button size="sm" variant="ghost" onClick={() => cancel(a.id)}>Dar de baja</Button>}
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
          <DialogHeader><DialogTitle>Asignar Beneficio</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2"><Label>Empleado *</Label><Select options={employees} placeholder="Seleccionar..." value={form.employee_id} onChange={e => setForm({ ...form, employee_id: e.target.value })} /></div>
            <div className="space-y-2"><Label>Beneficio *</Label><Select options={benefits} placeholder="Seleccionar..." value={form.benefit_id} onChange={e => setForm({ ...form, benefit_id: e.target.value })} /></div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button><Button onClick={enroll} disabled={saving || !form.employee_id || !form.benefit_id}>{saving ? 'Guardando...' : 'Asignar'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function RequestsTab() {
  const [requests, setRequests] = useState<any[]>([])
  const [employees, setEmployees] = useState<{ value: string; label: string }[]>([])
  const [benefits, setBenefits] = useState<{ value: string; label: string }[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [saving, setSaving] = useState(false)
  const [form, setForm] = useState({ employee_id: '', benefit_id: '', request_type: 'ENROLLMENT' })

  const fetch = async () => {
    setLoading(true)
    try { const r = await api.get('/benefits/requests'); setRequests(data(r)) } catch { setRequests([]) } finally { setLoading(false) }
  }
  useEffect(() => {
    fetch()
    api.get('/employees', { params: { limit: '200' } }).then(r => setEmployees((r.data.data ?? []).map((e: any) => ({ value: e.id, label: `${e.first_name} ${e.last_name}` })))).catch(() => {})
    api.get('/benefits').then(r => setBenefits(data(r).map((b: any) => ({ value: b.id, label: b.name })))).catch(() => {})
  }, [])

  const create = async () => {
    setSaving(true)
    try {
      await api.post('/benefits/requests', { employee_id: form.employee_id, benefit_id: form.benefit_id, request_type: form.request_type })
      setShowModal(false); setForm({ employee_id: '', benefit_id: '', request_type: 'ENROLLMENT' }); fetch()
    } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }

  const act = async (id: string, actName: string) => {
    const body: Record<string, any> = {}
    if (actName === 'reject') { const reason = window.prompt('Motivo del rechazo:'); if (!reason) return; body.reason = reason }
    if (actName === 'approve') { const comment = window.prompt('Comentario (opcional):'); if (comment !== null) body.comment = comment }
    try { await api.post(`/benefits/requests/${id}/${actName}`, body); fetch() } catch (e: any) { alert(e?.response?.data?.error || 'Error') }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-slate-500">{requests.length} solicitudes</p>
        <Button size="sm" onClick={() => setShowModal(true)}><Plus size={14} /> Nueva Solicitud</Button>
      </div>
      <Card>
        <CardContent className="p-0">
          {loading ? <div className="p-6 text-center text-slate-500">Cargando...</div> : requests.length === 0 ? <div className="p-6 text-center text-slate-500">No hay solicitudes</div> : (
            <table className="w-full text-sm">
              <thead><tr className="border-b border-slate-200 bg-slate-50">
                <th className="text-left px-4 py-3 font-medium text-slate-600">Empleado</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Beneficio</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Tipo</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Fecha</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
              </tr></thead>
              <tbody>{requests.map((r: any) => (
                <tr key={r.id} className="border-b border-slate-100 hover:bg-slate-50">
                  <td className="px-4 py-3 font-medium text-slate-900">{r.employee_id}</td>
                  <td className="px-4 py-3 text-slate-600">{r.benefit_id}</td>
                  <td className="px-4 py-3 text-slate-600">{r.request_type}</td>
                  <td className="px-4 py-3 text-slate-600">{fmt(r.created_at)}</td>
                  <td className="px-4 py-3">{pill(r.status)}</td>
                  <td className="px-4 py-3">
                    <div className="flex justify-end gap-1">
                      {r.status === 'DRAFT' && <Button size="sm" variant="outline" onClick={() => act(r.id, 'submit')}>Enviar</Button>}
                      {r.status === 'SUBMITTED' && <><Button size="sm" variant="outline" onClick={() => act(r.id, 'approve')}>Aprobar</Button><Button size="sm" variant="ghost" onClick={() => act(r.id, 'reject')}>Rechazar</Button></>}
                      {['DRAFT', 'SUBMITTED'].includes(r.status) && <Button size="sm" variant="ghost" onClick={() => act(r.id, 'cancel')}>Cancelar</Button>}
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
          <DialogHeader><DialogTitle>Nueva Solicitud</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2"><Label>Empleado *</Label><Select options={employees} placeholder="Seleccionar..." value={form.employee_id} onChange={e => setForm({ ...form, employee_id: e.target.value })} /></div>
            <div className="space-y-2"><Label>Beneficio *</Label><Select options={benefits} placeholder="Seleccionar..." value={form.benefit_id} onChange={e => setForm({ ...form, benefit_id: e.target.value })} /></div>
            <div className="space-y-2"><Label>Tipo *</Label><Select options={['ENROLLMENT', 'CHANGE_PLAN', 'CANCELLATION', 'CLAIM'].map(t => ({ value: t, label: t }))} value={form.request_type} onChange={e => setForm({ ...form, request_type: e.target.value })} /></div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button><Button onClick={create} disabled={saving || !form.employee_id || !form.benefit_id}>{saving ? 'Guardando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function CatalogTab() {
  const [categories, setCategories] = useState<any[]>([])
  const [types, setTypes] = useState<any[]>([])
  const [showCat, setShowCat] = useState(false)
  const [showType, setShowType] = useState(false)
  const [saving, setSaving] = useState(false)
  const [catForm, setCatForm] = useState({ name: '', description: '' })
  const [typeForm, setTypeForm] = useState({ name: '', code: '', nature: 'BENEFICIO', tax_treatment: 'NO_GRAVADO' })

  const fetchAll = async () => {
    api.get('/benefits/categories').then(r => setCategories(data(r))).catch(() => {})
    api.get('/benefits/types').then(r => setTypes(data(r))).catch(() => {})
  }
  useEffect(() => { fetchAll() }, [])

  const createCat = async () => {
    setSaving(true)
    try { await api.post('/benefits/categories', { ...catForm, sort_order: 0 }); setShowCat(false); setCatForm({ name: '', description: '' }); fetchAll() } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }
  const createType = async () => {
    setSaving(true)
    try { await api.post('/benefits/types', { ...typeForm, requires_approval: false, is_reimbursable: false, is_flexible: false, has_wallet: false, sort_order: 0 }); setShowType(false); setTypeForm({ name: '', code: '', nature: 'BENEFICIO', tax_treatment: 'NO_GRAVADO' }); fetchAll() } catch (e: any) { alert(e?.response?.data?.error || 'Error') } finally { setSaving(false) }
  }

  return (
    <div className="grid gap-6 lg:grid-cols-2">
      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="font-semibold text-slate-900">Categorías</h2>
          <Button size="sm" onClick={() => setShowCat(true)}><Plus size={14} /> Nueva</Button>
        </div>
        <Card>
          <CardContent className="p-0">
            {categories.length === 0 ? <div className="p-6 text-center text-slate-500">Sin categorías</div> : (
              <table className="w-full text-sm">
                <thead><tr className="border-b border-slate-200 bg-slate-50">
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Descripción</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                </tr></thead>
                <tbody>{categories.map((c: any) => (
                  <tr key={c.id} className="border-b border-slate-100">
                    <td className="px-4 py-3 font-medium text-slate-900">{c.name}</td>
                    <td className="px-4 py-3 text-slate-600">{c.description || '-'}</td>
                    <td className="px-4 py-3">{pill(c.is_active ? 'ACTIVE' : 'INACTIVE')}</td>
                  </tr>
                ))}</tbody>
              </table>
            )}
          </CardContent>
        </Card>
      </div>
      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="font-semibold text-slate-900">Tipos de Beneficio</h2>
          <Button size="sm" onClick={() => setShowType(true)}><Plus size={14} /> Nuevo</Button>
        </div>
        <Card>
          <CardContent className="p-0">
            {types.length === 0 ? <div className="p-6 text-center text-slate-500">Sin tipos</div> : (
              <table className="w-full text-sm">
                <thead><tr className="border-b border-slate-200 bg-slate-50">
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Código</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Naturaleza</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                </tr></thead>
                <tbody>{types.map((t: any) => (
                  <tr key={t.id} className="border-b border-slate-100">
                    <td className="px-4 py-3 font-mono text-slate-600">{t.code}</td>
                    <td className="px-4 py-3 font-medium text-slate-900">{t.name}</td>
                    <td className="px-4 py-3 text-slate-600">{t.nature}</td>
                    <td className="px-4 py-3">{pill(t.is_active ? 'ACTIVE' : 'INACTIVE')}</td>
                  </tr>
                ))}</tbody>
              </table>
            )}
          </CardContent>
        </Card>
      </div>

      <Dialog open={showCat} onOpenChange={setShowCat}>
        <DialogContent>
          <DialogHeader><DialogTitle>Nueva Categoría</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2"><Label>Nombre *</Label><Input value={catForm.name} onChange={e => setCatForm({ ...catForm, name: e.target.value })} /></div>
            <div className="space-y-2"><Label>Descripción</Label><Input value={catForm.description} onChange={e => setCatForm({ ...catForm, description: e.target.value })} /></div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowCat(false)}>Cancelar</Button><Button onClick={createCat} disabled={saving || !catForm.name}>{saving ? 'Guardando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={showType} onOpenChange={setShowType}>
        <DialogContent>
          <DialogHeader><DialogTitle>Nuevo Tipo de Beneficio</DialogTitle></DialogHeader>
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"><Label>Nombre *</Label><Input value={typeForm.name} onChange={e => setTypeForm({ ...typeForm, name: e.target.value })} /></div>
              <div className="space-y-2"><Label>Código *</Label><Input value={typeForm.code} onChange={e => setTypeForm({ ...typeForm, code: e.target.value })} /></div>
            </div>
            <div className="space-y-2"><Label>Naturaleza *</Label><Select options={['BENEFICIO', 'SERVICIO', 'REEMBOLSO', 'FLEXIBLE', 'RECONOCIMIENTO'].map(t => ({ value: t, label: t }))} value={typeForm.nature} onChange={e => setTypeForm({ ...typeForm, nature: e.target.value })} /></div>
            <div className="space-y-2"><Label>Tratamiento fiscal</Label><Select options={['GRAVADO', 'NO_GRAVADO', 'EXENTO'].map(t => ({ value: t, label: t }))} value={typeForm.tax_treatment} onChange={e => setTypeForm({ ...typeForm, tax_treatment: e.target.value })} /></div>
          </div>
          <DialogFooter><Button variant="outline" onClick={() => setShowType(false)}>Cancelar</Button><Button onClick={createType} disabled={saving || !typeForm.name || !typeForm.code}>{saving ? 'Guardando...' : 'Crear'}</Button></DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
