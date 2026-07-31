import { useEffect, useState } from 'react'
import { Plus, Pencil, Users, FileText, BarChart3, Send, CheckCircle2, Rocket, XCircle } from 'lucide-react'
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

type Tab = 'offers' | 'postings' | 'candidates'

interface Offer {
  id: string
  position_title: string
  offer_type?: string
  department_id?: string
  employment_type?: string
  work_mode?: string
  salary_amount?: number
  salary_currency?: string
  salary_period?: string
  status: string
  created_at: string
}

interface Posting {
  id: string
  title: string
  position_id?: string
  description?: string
  employment_type?: string
  work_mode?: string
  location?: string
  salary_min?: number
  salary_max?: number
  currency?: string
  closing_at?: string
  is_public: boolean
  status: string
}

interface Candidate {
  id: string
  first_name: string
  last_name: string
  email: string
  phone?: string
  current_company?: string
  current_position?: string
  source?: string
  status: string
  created_at: string
}

interface SelectOption {
  value: string
  label: string
}

const emptyOfferForm = {
  position_title: '',
  offer_type: '',
  department_id: '',
  employment_type: '',
  work_mode: '',
  salary_amount: '',
  salary_currency: 'ARS',
  salary_period: '',
  start_date: '',
  response_deadline: '',
  benefits_summary: '',
  conditions: '',
  notes: '',
}

const emptyPostingForm = {
  title: '',
  position_id: '',
  description: '',
  employment_type: '',
  work_mode: '',
  location: '',
  salary_min: '',
  salary_max: '',
  currency: 'ARS',
  closing_at: '',
  is_public: false,
}

const emptyCandidateForm = {
  first_name: '',
  last_name: '',
  email: '',
  phone: '',
  document_type: '',
  document_number: '',
  location: '',
  current_company: '',
  current_position: '',
  salary_expectation_min: '',
  salary_expectation_max: '',
  salary_currency: 'ARS',
  source: '',
  notes: '',
}

const statusBadge = (status: string) => {
  const map: Record<string, { label: string; cls: string }> = {
    DRAFT: { label: 'Borrador', cls: 'bg-slate-100 text-slate-600' },
    PENDING_APPROVAL: { label: 'Pendiente aprobación', cls: 'bg-amber-50 text-amber-700' },
    APPROVED: { label: 'Aprobado', cls: 'bg-emerald-50 text-emerald-700' },
    SENT: { label: 'Enviado', cls: 'bg-blue-50 text-blue-700' },
    ACCEPTED: { label: 'Aceptado', cls: 'bg-teal-50 text-teal-700' },
    REJECTED: { label: 'Rechazado', cls: 'bg-red-50 text-red-700' },
    WITHDRAWN: { label: 'Retirado', cls: 'bg-slate-100 text-slate-500' },
    PUBLISHED: { label: 'Publicado', cls: 'bg-emerald-50 text-emerald-700' },
    CLOSED: { label: 'Cerrado', cls: 'bg-slate-100 text-slate-600' },
    CANCELLED: { label: 'Cancelado', cls: 'bg-red-50 text-red-700' },
    ACTIVE: { label: 'Activo', cls: 'bg-emerald-50 text-emerald-700' },
    INACTIVE: { label: 'Inactivo', cls: 'bg-slate-100 text-slate-600' },
    BLACKLISTED: { label: 'Bloqueado', cls: 'bg-red-50 text-red-700' },
    HIRED: { label: 'Contratado', cls: 'bg-teal-50 text-teal-700' },
  }
  const s = map[status] ?? { label: status, cls: 'bg-slate-100 text-slate-600' }
  return <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${s.cls}`}>{s.label}</span>
}

const tabs = [
  { key: 'offers', label: 'Ofertas', icon: FileText },
  { key: 'postings', label: 'Publicaciones', icon: BarChart3 },
  { key: 'candidates', label: 'Candidatos', icon: Users },
] as const

export default function RecruitmentPage() {
  const [tab, setTab] = useState<Tab>('offers')
  const [items, setItems] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showModal, setShowModal] = useState(false)
  const [editing, setEditing] = useState<any | null>(null)
  const [saving, setSaving] = useState(false)
  const [departments, setDepartments] = useState<SelectOption[]>([])
  const [recPositions, setRecPositions] = useState<SelectOption[]>([])

  const [offerForm, setOfferForm] = useState({ ...emptyOfferForm })
  const [postingForm, setPostingForm] = useState({ ...emptyPostingForm })
  const [candidateForm, setCandidateForm] = useState({ ...emptyCandidateForm })

  const fetchData = async () => {
    setLoading(true)
    try {
      const endpoint = tab === 'offers' ? '/recruitment/offers' : tab === 'postings' ? '/recruitment/postings' : '/recruitment/candidates'
      const res = await api.get(endpoint, { params: { limit: '100' } })
      const data = res.data.data ?? []
      if (tab === 'offers') setItems(data as Offer[])
      else if (tab === 'postings') setItems(data as Posting[])
      else setItems(data as Candidate[])
      setError('')
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar registros')
      setItems([])
    } finally {
      setLoading(false)
    }
  }

  const fetchSelects = async () => {
    try {
      const [dRes, pRes] = await Promise.all([
        api.get('/departments', { params: { limit: '100' } }),
        api.get('/recruitment/positions', { params: { limit: '100' } }),
      ])
      setDepartments((dRes.data.data ?? []).map((d: any) => ({ value: d.id, label: d.name })))
      setRecPositions((pRes.data.data ?? []).map((p: any) => ({ value: p.id, label: p.title })))
    } catch {}
  }

  useEffect(() => { fetchData() }, [tab])

  const openCreate = () => {
    setEditing(null)
    setOfferForm({ ...emptyOfferForm })
    setPostingForm({ ...emptyPostingForm })
    setCandidateForm({ ...emptyCandidateForm })
    fetchSelects()
    setShowModal(true)
  }

  const openEdit = (item: any) => {
    setEditing(item)
    fetchSelects()
    if (tab === 'offers') {
      setOfferForm({
        position_title: item.position_title ?? '',
        offer_type: item.offer_type ?? '',
        department_id: item.department_id ?? '',
        employment_type: item.employment_type ?? '',
        work_mode: item.work_mode ?? '',
        salary_amount: item.salary_amount != null ? String(item.salary_amount) : '',
        salary_currency: item.salary_currency ?? 'ARS',
        salary_period: item.salary_period ?? '',
        start_date: item.start_date ? item.start_date.slice(0, 10) : '',
        response_deadline: item.response_deadline ? item.response_deadline.slice(0, 10) : '',
        benefits_summary: item.benefits_summary ?? '',
        conditions: item.conditions ?? '',
        notes: item.notes ?? '',
      })
    } else if (tab === 'postings') {
      setPostingForm({
        title: item.title ?? '',
        position_id: item.position_id ?? '',
        description: item.description ?? '',
        employment_type: item.employment_type ?? '',
        work_mode: item.work_mode ?? '',
        location: item.location ?? '',
        salary_min: item.salary_min != null ? String(item.salary_min) : '',
        salary_max: item.salary_max != null ? String(item.salary_max) : '',
        currency: item.currency ?? 'ARS',
        closing_at: item.closing_at ? item.closing_at.slice(0, 10) : '',
        is_public: item.is_public ?? false,
      })
    } else {
      setCandidateForm({
        first_name: item.first_name ?? '',
        last_name: item.last_name ?? '',
        email: item.email ?? '',
        phone: item.phone ?? '',
        document_type: item.document_type ?? '',
        document_number: item.document_number ?? '',
        location: item.location ?? '',
        current_company: item.current_company ?? '',
        current_position: item.current_position ?? '',
        salary_expectation_min: item.salary_expectation_min != null ? String(item.salary_expectation_min) : '',
        salary_expectation_max: item.salary_expectation_max != null ? String(item.salary_expectation_max) : '',
        salary_currency: item.salary_currency ?? 'ARS',
        source: item.source ?? '',
        notes: item.notes ?? '',
      })
    }
    setShowModal(true)
  }

  const handleSave = async () => {
    setSaving(true)
    try {
      if (tab === 'offers') {
        const f = offerForm
        const body: Record<string, any> = {}
        if (f.position_title) body.position_title = f.position_title
        if (f.offer_type) body.offer_type = f.offer_type
        if (f.department_id) body.department_id = f.department_id
        if (f.employment_type) body.employment_type = f.employment_type
        if (f.work_mode) body.work_mode = f.work_mode
        if (f.salary_amount) body.salary_amount = parseFloat(f.salary_amount)
        if (f.salary_currency) body.salary_currency = f.salary_currency
        if (f.salary_period) body.salary_period = f.salary_period
        if (f.start_date) body.start_date = f.start_date
        if (f.response_deadline) body.response_deadline = f.response_deadline
        if (f.benefits_summary) body.benefits_summary = f.benefits_summary
        if (f.conditions) body.conditions = f.conditions
        if (f.notes) body.notes = f.notes
        if (editing) await api.put(`/recruitment/offers/${editing.id}`, body)
        else await api.post('/recruitment/offers', body)
      } else if (tab === 'postings') {
        const f = postingForm
        const body: Record<string, any> = { is_public: f.is_public }
        if (f.title) body.title = f.title
        if (f.position_id) body.position_id = f.position_id
        if (f.description) body.description = f.description
        if (f.employment_type) body.employment_type = f.employment_type
        if (f.work_mode) body.work_mode = f.work_mode
        if (f.location) body.location = f.location
        if (f.salary_min) body.salary_min = parseFloat(f.salary_min)
        if (f.salary_max) body.salary_max = parseFloat(f.salary_max)
        if (f.currency) body.currency = f.currency
        if (f.closing_at) body.closing_at = f.closing_at
        if (editing) await api.put(`/recruitment/postings/${editing.id}`, body)
        else await api.post('/recruitment/postings', body)
      } else {
        const f = candidateForm
        const body: Record<string, any> = {}
        if (f.first_name) body.first_name = f.first_name
        if (f.last_name) body.last_name = f.last_name
        if (f.email) body.email = f.email
        if (f.phone) body.phone = f.phone
        if (f.document_type) body.document_type = f.document_type
        if (f.document_number) body.document_number = f.document_number
        if (f.location) body.location = f.location
        if (f.current_company) body.current_company = f.current_company
        if (f.current_position) body.current_position = f.current_position
        if (f.salary_expectation_min) body.salary_expectation_min = parseFloat(f.salary_expectation_min)
        if (f.salary_expectation_max) body.salary_expectation_max = parseFloat(f.salary_expectation_max)
        if (f.salary_currency) body.salary_currency = f.salary_currency
        if (f.source) body.source = f.source
        if (f.notes) body.notes = f.notes
        if (editing) await api.put(`/recruitment/candidates/${editing.id}`, body)
        else await api.post('/recruitment/candidates', body)
      }
      setShowModal(false)
      fetchData()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar registro')
    } finally {
      setSaving(false)
    }
  }

  const runAction = async (url: string, confirmMsg?: string) => {
    if (confirmMsg && !confirm(confirmMsg)) return
    try {
      await api.post(url)
      fetchData()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al ejecutar acción')
    }
  }

  const canSubmit = (item: any) => item.status === 'DRAFT'
  const canApprove = (item: any) => item.status === 'PENDING_APPROVAL'
  const canSend = (item: any) => item.status === 'APPROVED'

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-slate-900">Reclutamiento</h1>
        <Button size="sm" onClick={openCreate}><Plus size={16} className="mr-1" />
          {tab === 'offers' ? 'Nueva oferta' : tab === 'postings' ? 'Nueva publicación' : 'Nuevo candidato'}
        </Button>
      </div>

      <div className="flex gap-2 mb-4">
        {tabs.map(t => (
          <button key={t.key} onClick={() => setTab(t.key)}
            className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${tab === t.key ? 'bg-brand-50 text-brand-700' : 'text-slate-600 hover:bg-slate-100'}`}>
            <t.icon size={16} /> {t.label}
          </button>
        ))}
      </div>

      {error && <div className="mb-4 p-3 bg-red-50 text-red-700 text-sm rounded-lg">{error}</div>}

      <Card>
        <CardContent className="p-0">
          {loading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
          : items.length === 0 && !error ? <div className="p-6 text-center text-slate-500">No hay registros</div>
          : <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-slate-200 bg-slate-50">
                    {tab === 'offers' && (
                      <>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Posición</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Salario</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                      </>
                    )}
                    {tab === 'postings' && (
                      <>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Título</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Ubicación</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                      </>
                    )}
                    {tab === 'candidates' && (
                      <>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Email</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Posición actual</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                      </>
                    )}
                    <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                  </tr>
                </thead>
                <tbody>
                  {items.map(o => (
                    <tr key={o.id} className="border-b border-slate-100 hover:bg-slate-50">
                      {tab === 'offers' && (
                        <>
                          <td className="px-4 py-3 font-medium text-slate-900">{o.position_title}</td>
                          <td className="px-4 py-3 text-slate-600">
                            {o.salary_amount != null ? `${o.salary_amount.toLocaleString()} ${o.salary_currency || ''}` : '-'}
                          </td>
                          <td className="px-4 py-3">{statusBadge(o.status)}</td>
                        </>
                      )}
                      {tab === 'postings' && (
                        <>
                          <td className="px-4 py-3 font-medium text-slate-900">{o.title}</td>
                          <td className="px-4 py-3 text-slate-600">{o.location || '-'}</td>
                          <td className="px-4 py-3">{statusBadge(o.status)}</td>
                        </>
                      )}
                      {tab === 'candidates' && (
                        <>
                          <td className="px-4 py-3 font-medium text-slate-900">{o.first_name} {o.last_name}</td>
                          <td className="px-4 py-3 text-slate-600">{o.email}</td>
                          <td className="px-4 py-3 text-slate-600">{o.current_position || '-'}</td>
                          <td className="px-4 py-3">{statusBadge(o.status)}</td>
                        </>
                      )}
                      <td className="px-4 py-3 text-right whitespace-nowrap">
                        {tab === 'offers' && (
                          <>
                            {canSubmit(o) && (
                              <Button variant="ghost" size="sm" className="text-amber-600" onClick={() => runAction(`/recruitment/offers/${o.id}/submit`)} title="Enviar a aprobación"><Send size={14} /></Button>
                            )}
                            {canApprove(o) && (
                              <Button variant="ghost" size="sm" className="text-emerald-600" onClick={() => runAction(`/recruitment/offers/${o.id}/approve`)} title="Aprobar"><CheckCircle2 size={14} /></Button>
                            )}
                            {canSend(o) && (
                              <Button variant="ghost" size="sm" className="text-blue-600" onClick={() => runAction(`/recruitment/offers/${o.id}/send`)} title="Enviar oferta"><Rocket size={14} /></Button>
                            )}
                            {canApprove(o) && (
                              <Button variant="ghost" size="sm" className="text-red-500" onClick={() => runAction(`/recruitment/offers/${o.id}/reject`)} title="Rechazar"><XCircle size={14} /></Button>
                            )}
                          </>
                        )}
                        {tab === 'postings' && (
                          <>
                            {o.status === 'DRAFT' && (
                              <Button variant="ghost" size="sm" className="text-emerald-600" onClick={() => runAction(`/recruitment/postings/${o.id}/publish`)} title="Publicar"><Rocket size={14} /></Button>
                            )}
                            {o.status === 'PUBLISHED' && (
                              <Button variant="ghost" size="sm" className="text-slate-500" onClick={() => runAction(`/recruitment/postings/${o.id}/close`, `¿Cerrar la publicación "${o.title}"?`)} title="Cerrar"><XCircle size={14} /></Button>
                            )}
                          </>
                        )}
                        <Button variant="ghost" size="sm" onClick={() => openEdit(o)}><Pencil size={14} /></Button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>}
        </CardContent>
      </Card>

      <Dialog open={showModal} onOpenChange={setShowModal}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>
              {tab === 'offers' ? (editing ? 'Editar Oferta' : 'Nueva Oferta')
                : tab === 'postings' ? (editing ? 'Editar Publicación' : 'Nueva Publicación')
                : (editing ? 'Editar Candidato' : 'Nuevo Candidato')}
            </DialogTitle>
            <DialogDescription>Completá los datos del registro</DialogDescription>
          </DialogHeader>

          {tab === 'offers' && (
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2 col-span-2">
                <Label htmlFor="offer-title">Posición *</Label>
                <Input id="offer-title" value={offerForm.position_title} onChange={e => setOfferForm({ ...offerForm, position_title: e.target.value })} required />
              </div>
              <div className="space-y-2">
                <Label htmlFor="offer-type">Tipo de oferta</Label>
                <Input id="offer-type" value={offerForm.offer_type} onChange={e => setOfferForm({ ...offerForm, offer_type: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="offer-dept">Departamento</Label>
                <Select id="offer-dept" options={departments} placeholder="Seleccionar..." value={offerForm.department_id} onChange={e => setOfferForm({ ...offerForm, department_id: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="offer-emp">Tipo de contratación</Label>
                <Select
                  id="offer-emp"
                  options={[{ value: 'FULL_TIME', label: 'Tiempo completo' }, { value: 'PART_TIME', label: 'Medio tiempo' }, { value: 'CONTRACTOR', label: 'Contratista' }, { value: 'INTERNSHIP', label: 'Pasantía' }]}
                  placeholder="Seleccionar..."
                  value={offerForm.employment_type}
                  onChange={e => setOfferForm({ ...offerForm, employment_type: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="offer-mode">Modalidad</Label>
                <Select
                  id="offer-mode"
                  options={[{ value: 'REMOTE', label: 'Remoto' }, { value: 'ONSITE', label: 'Presencial' }, { value: 'HYBRID', label: 'Híbrido' }]}
                  placeholder="Seleccionar..."
                  value={offerForm.work_mode}
                  onChange={e => setOfferForm({ ...offerForm, work_mode: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="offer-amount">Salario</Label>
                <Input id="offer-amount" type="number" step="0.01" value={offerForm.salary_amount} onChange={e => setOfferForm({ ...offerForm, salary_amount: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="offer-currency">Moneda</Label>
                <Select
                  id="offer-currency"
                  options={[{ value: 'ARS', label: 'ARS' }, { value: 'USD', label: 'USD' }, { value: 'EUR', label: 'EUR' }]}
                  value={offerForm.salary_currency}
                  onChange={e => setOfferForm({ ...offerForm, salary_currency: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="offer-period">Periodicidad</Label>
                <Select
                  id="offer-period"
                  options={[{ value: 'MONTHLY', label: 'Mensual' }, { value: 'ANNUAL', label: 'Anual' }, { value: 'HOURLY', label: 'Por hora' }]}
                  placeholder="Seleccionar..."
                  value={offerForm.salary_period}
                  onChange={e => setOfferForm({ ...offerForm, salary_period: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="offer-start">Fecha de inicio</Label>
                <Input id="offer-start" type="date" value={offerForm.start_date} onChange={e => setOfferForm({ ...offerForm, start_date: e.target.value })} />
              </div>
              <div className="space-y-2 col-span-2">
                <Label htmlFor="offer-deadline">Fecha límite de respuesta</Label>
                <Input id="offer-deadline" type="date" value={offerForm.response_deadline} onChange={e => setOfferForm({ ...offerForm, response_deadline: e.target.value })} />
              </div>
              <div className="space-y-2 col-span-2">
                <Label htmlFor="offer-benefits">Beneficios</Label>
                <textarea
                  id="offer-benefits"
                  value={offerForm.benefits_summary}
                  onChange={e => setOfferForm({ ...offerForm, benefits_summary: e.target.value })}
                  rows={2}
                  className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-brand-500"
                />
              </div>
              <div className="space-y-2 col-span-2">
                <Label htmlFor="offer-conditions">Condiciones</Label>
                <textarea
                  id="offer-conditions"
                  value={offerForm.conditions}
                  onChange={e => setOfferForm({ ...offerForm, conditions: e.target.value })}
                  rows={2}
                  className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-brand-500"
                />
              </div>
              <div className="space-y-2 col-span-2">
                <Label htmlFor="offer-notes">Notas</Label>
                <textarea
                  id="offer-notes"
                  value={offerForm.notes}
                  onChange={e => setOfferForm({ ...offerForm, notes: e.target.value })}
                  rows={2}
                  className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-brand-500"
                />
              </div>
            </div>
          )}

          {tab === 'postings' && (
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2 col-span-2">
                <Label htmlFor="post-title">Título *</Label>
                <Input id="post-title" value={postingForm.title} onChange={e => setPostingForm({ ...postingForm, title: e.target.value })} required />
              </div>
              <div className="space-y-2 col-span-2">
                <Label htmlFor="post-pos">Posición (reclutamiento)</Label>
                <Select id="post-pos" options={recPositions} placeholder="Seleccionar..." value={postingForm.position_id} onChange={e => setPostingForm({ ...postingForm, position_id: e.target.value })} />
              </div>
              <div className="space-y-2 col-span-2">
                <Label htmlFor="post-desc">Descripción *</Label>
                <textarea
                  id="post-desc"
                  value={postingForm.description}
                  onChange={e => setPostingForm({ ...postingForm, description: e.target.value })}
                  rows={3}
                  className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-brand-500"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="post-emp">Tipo de contratación</Label>
                <Select
                  id="post-emp"
                  options={[{ value: 'FULL_TIME', label: 'Tiempo completo' }, { value: 'PART_TIME', label: 'Medio tiempo' }, { value: 'CONTRACTOR', label: 'Contratista' }, { value: 'INTERNSHIP', label: 'Pasantía' }]}
                  placeholder="Seleccionar..."
                  value={postingForm.employment_type}
                  onChange={e => setPostingForm({ ...postingForm, employment_type: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="post-mode">Modalidad</Label>
                <Select
                  id="post-mode"
                  options={[{ value: 'REMOTE', label: 'Remoto' }, { value: 'ONSITE', label: 'Presencial' }, { value: 'HYBRID', label: 'Híbrido' }]}
                  placeholder="Seleccionar..."
                  value={postingForm.work_mode}
                  onChange={e => setPostingForm({ ...postingForm, work_mode: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="post-loc">Ubicación</Label>
                <Input id="post-loc" value={postingForm.location} onChange={e => setPostingForm({ ...postingForm, location: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="post-close">Cierre de publicación</Label>
                <Input id="post-close" type="date" value={postingForm.closing_at} onChange={e => setPostingForm({ ...postingForm, closing_at: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="post-min">Salario mínimo</Label>
                <Input id="post-min" type="number" step="0.01" value={postingForm.salary_min} onChange={e => setPostingForm({ ...postingForm, salary_min: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="post-max">Salario máximo</Label>
                <Input id="post-max" type="number" step="0.01" value={postingForm.salary_max} onChange={e => setPostingForm({ ...postingForm, salary_max: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="post-currency">Moneda</Label>
                <Select
                  id="post-currency"
                  options={[{ value: 'ARS', label: 'ARS' }, { value: 'USD', label: 'USD' }, { value: 'EUR', label: 'EUR' }]}
                  value={postingForm.currency}
                  onChange={e => setPostingForm({ ...postingForm, currency: e.target.value })}
                />
              </div>
              <div className="space-y-2 col-span-2">
                <label className="flex items-center gap-2 text-sm text-slate-700">
                  <input type="checkbox" checked={postingForm.is_public} onChange={e => setPostingForm({ ...postingForm, is_public: e.target.checked })} className="h-4 w-4 rounded border-slate-300 text-brand-600 focus:ring-brand-500" />
                  Publicación pública
                </label>
              </div>
            </div>
          )}

          {tab === 'candidates' && (
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="cand-first">Nombre *</Label>
                <Input id="cand-first" value={candidateForm.first_name} onChange={e => setCandidateForm({ ...candidateForm, first_name: e.target.value })} required />
              </div>
              <div className="space-y-2">
                <Label htmlFor="cand-last">Apellido *</Label>
                <Input id="cand-last" value={candidateForm.last_name} onChange={e => setCandidateForm({ ...candidateForm, last_name: e.target.value })} required />
              </div>
              <div className="space-y-2 col-span-2">
                <Label htmlFor="cand-email">Email *</Label>
                <Input id="cand-email" type="email" value={candidateForm.email} onChange={e => setCandidateForm({ ...candidateForm, email: e.target.value })} required />
              </div>
              <div className="space-y-2">
                <Label htmlFor="cand-phone">Teléfono</Label>
                <Input id="cand-phone" value={candidateForm.phone} onChange={e => setCandidateForm({ ...candidateForm, phone: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="cand-doctype">Tipo de documento</Label>
                <Input id="cand-doctype" value={candidateForm.document_type} onChange={e => setCandidateForm({ ...candidateForm, document_type: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="cand-docnum">N° de documento</Label>
                <Input id="cand-docnum" value={candidateForm.document_number} onChange={e => setCandidateForm({ ...candidateForm, document_number: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="cand-loc">Ubicación</Label>
                <Input id="cand-loc" value={candidateForm.location} onChange={e => setCandidateForm({ ...candidateForm, location: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="cand-company">Empresa actual</Label>
                <Input id="cand-company" value={candidateForm.current_company} onChange={e => setCandidateForm({ ...candidateForm, current_company: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="cand-position">Posición actual</Label>
                <Input id="cand-position" value={candidateForm.current_position} onChange={e => setCandidateForm({ ...candidateForm, current_position: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="cand-src">Fuente</Label>
                <Input id="cand-src" value={candidateForm.source} onChange={e => setCandidateForm({ ...candidateForm, source: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="cand-salmin">Expectativa salarial mínima</Label>
                <Input id="cand-salmin" type="number" step="0.01" value={candidateForm.salary_expectation_min} onChange={e => setCandidateForm({ ...candidateForm, salary_expectation_min: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="cand-salmax">Expectativa salarial máxima</Label>
                <Input id="cand-salmax" type="number" step="0.01" value={candidateForm.salary_expectation_max} onChange={e => setCandidateForm({ ...candidateForm, salary_expectation_max: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="cand-salcur">Moneda</Label>
                <Select
                  id="cand-salcur"
                  options={[{ value: 'ARS', label: 'ARS' }, { value: 'USD', label: 'USD' }, { value: 'EUR', label: 'EUR' }]}
                  value={candidateForm.salary_currency}
                  onChange={e => setCandidateForm({ ...candidateForm, salary_currency: e.target.value })}
                />
              </div>
              <div className="space-y-2 col-span-2">
                <Label htmlFor="cand-notes">Notas</Label>
                <textarea
                  id="cand-notes"
                  value={candidateForm.notes}
                  onChange={e => setCandidateForm({ ...candidateForm, notes: e.target.value })}
                  rows={2}
                  className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-brand-500"
                />
              </div>
            </div>
          )}

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button>
            <Button
              onClick={handleSave}
              disabled={
                saving ||
                (tab === 'offers' && !offerForm.position_title) ||
                (tab === 'postings' && !postingForm.title) ||
                (tab === 'candidates' && (!candidateForm.first_name || !candidateForm.last_name || !candidateForm.email))
              }
            >
              {saving ? 'Guardando...' : editing ? 'Guardar Cambios' : 'Crear'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
