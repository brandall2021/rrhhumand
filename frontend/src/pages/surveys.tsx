import { useEffect, useState } from 'react'
import {
  Plus,
  Trash2,
  Send,
  FileDown,
  ClipboardList,
  CheckCircle2,
  XCircle,
} from 'lucide-react'
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

interface Survey {
  id: string
  title: string
  description?: string
  type: string
  status: string
  anonymous: boolean
  multiple_responses: boolean
  starts_at?: string
  ends_at?: string
  created_by_name?: string
  created_at: string
  questions?: SurveyQuestion[]
  targets?: SurveyTarget[]
  response_count?: number
  target_count?: number
  participation_rate?: number
}

interface SurveyQuestion {
  id: string
  question: string
  type: string
  position: number
  required: boolean
  options?: SurveyOption[]
}

interface SurveyOption {
  id: string
  question_id: string
  option_text: string
  position: number
}

interface SurveyTarget {
  id: string
  target_type: string
  target_id?: string
}

interface SurveyStats {
  total_targeted: number
  total_responded: number
  participation_rate: number
  questions: QuestionStats[]
}

interface QuestionStats {
  question_id: string
  question: string
  type: string
  total_answers: number
  average?: number
  min?: number
  max?: number
  distribution?: OptionDistribution[]
  yes_count?: number
  no_count?: number
  yes_percentage?: number
  sample_texts?: string[]
}

interface OptionDistribution {
  option_id: string
  option_text: string
  count: number
  percentage: number
}

const surveyTypes: Record<string, string> = {
  GENERAL: 'General',
  CLIMATE: 'Clima laboral',
  SATISFACTION: 'Satisfacción',
  FEEDBACK: 'Feedback',
  PULSE: 'Pulso',
  TRAINING: 'Capacitación',
  INTERNAL: 'Interna',
}

const questionTypes: Record<string, string> = {
  TEXT: 'Texto libre',
  NUMBER: 'Número',
  RATING: 'Calificación 1-5',
  SINGLE_CHOICE: 'Opción única',
  MULTIPLE_CHOICE: 'Opción múltiple',
  YES_NO: 'Sí / No',
}

const targetTypes: Record<string, string> = {
  ALL: 'Toda la empresa',
  DEPARTMENT: 'Departamento',
  BRANCH: 'Sucursal',
  POSITION: 'Posición',
  EMPLOYEE: 'Empleado',
}

const statusLabels: Record<string, string> = {
  DRAFT: 'Borrador',
  PUBLISHED: 'Publicada',
  CLOSED: 'Cerrada',
  ARCHIVED: 'Archivada',
}

const statusColors: Record<string, string> = {
  DRAFT: 'bg-slate-100 text-slate-600',
  PUBLISHED: 'bg-emerald-50 text-emerald-700',
  CLOSED: 'bg-amber-50 text-amber-700',
  ARCHIVED: 'bg-slate-100 text-slate-500',
}

export default function SurveysPage() {
  const [activeTab, setActiveTab] = useState('available')

  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-900 mb-6">Encuestas</h1>

      <div className="flex gap-1 mb-6 border-b border-slate-200">
        {[{ key: 'available', label: 'Disponibles' }, { key: 'admin', label: 'Administración' }].map((tab) => (
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

      {activeTab === 'available' ? <AvailableTab /> : <AdminTab />}
    </div>
  )
}

function AvailableTab() {
  const [surveys, setSurveys] = useState<Survey[]>([])
  const [loading, setLoading] = useState(true)
  const [responding, setResponding] = useState<Survey | null>(null)
  const [answers, setAnswers] = useState<Record<string, any>>({})
  const [submitting, setSubmitting] = useState(false)

  const fetch = async () => {
    setLoading(true)
    try {
      const res = await api.get('/me/surveys')
      setSurveys(res.data.data ?? [])
    } catch { setSurveys([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetch() }, [])

  const openRespond = async (survey: Survey) => {
    try {
      const res = await api.get(`/surveys/${survey.id}`)
      setResponding(res.data.data ?? survey)
      const initial: Record<string, any> = {}
      for (const q of (res.data.data?.questions ?? [])) {
        initial[q.id] = q.type === 'MULTIPLE_CHOICE' ? [] : ''
      }
      setAnswers(initial)
    } catch { setResponding(survey) }
  }

  const submit = async () => {
    if (!responding) return
    const payload: any[] = []
    for (const q of responding.questions ?? []) {
      const v = answers[q.id]
      if (q.type === 'TEXT') {
        if (v && v.trim()) payload.push({ question_id: q.id, text: v })
      } else if (q.type === 'NUMBER' || q.type === 'RATING') {
        if (v !== '' && v != null && !isNaN(Number(v))) payload.push({ question_id: q.id, number: Number(v) })
      } else if (q.type === 'MULTIPLE_CHOICE') {
        const selected = (v ?? []) as string[]
        if (selected.length > 0) payload.push({ question_id: q.id, option_ids: selected })
      } else {
        if (v) payload.push({ question_id: q.id, option_id: v })
      }
    }
    setSubmitting(true)
    try {
      await api.post(`/surveys/${responding.id}/respond`, { answers: payload })
      setResponding(null)
      fetch()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al enviar la encuesta')
    } finally { setSubmitting(false) }
  }

  return (
    <div>
      <p className="text-sm text-slate-500 mb-4">
        Encuestas disponibles para vos
      </p>

      {loading ? (
        <div className="text-center text-slate-500 py-12">Cargando...</div>
      ) : surveys.length === 0 ? (
        <Card>
          <CardContent className="p-10 text-center text-slate-500">No tenés encuestas disponibles</CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {surveys.map((s) => (
            <Card key={s.id}>
              <CardContent className="p-5">
                <div className="flex items-start justify-between gap-3">
                  <div>
                    <h3 className="font-semibold text-slate-900">{s.title}</h3>
                    <p className="text-xs text-slate-400 mt-0.5">{surveyTypes[s.type] || s.type}</p>
                  </div>
                  {s.anonymous && (
                    <span className="px-2 py-0.5 rounded-full bg-slate-100 text-slate-600 text-xs font-medium shrink-0">Anónima</span>
                  )}
                </div>
                {s.description && <p className="text-sm text-slate-600 mt-2">{s.description}</p>}
                <div className="flex items-center justify-between mt-4">
                  <p className="text-xs text-slate-400">
                    {s.ends_at ? `Vence: ${new Date(s.ends_at).toLocaleDateString('es-AR')}` : 'Sin fecha límite'}
                  </p>
                  <Button size="sm" onClick={() => openRespond(s)}>
                    <ClipboardList size={14} /> Responder
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      <Dialog open={!!responding} onOpenChange={(o) => !o && setResponding(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{responding?.title}</DialogTitle>
            <DialogDescription>{responding?.description || 'Respondé la encuesta'}</DialogDescription>
          </DialogHeader>
          <div className="space-y-5 max-h-[60vh] overflow-y-auto pr-1">
            {(responding?.questions ?? []).map((q) => (
              <div key={q.id} className="space-y-2">
                <Label>
                  {q.position}. {q.question} {q.required && <span className="text-red-500">*</span>}
                </Label>

                {q.type === 'TEXT' && (
                  <textarea
                    className="flex min-h-[70px] w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm"
                    value={answers[q.id] ?? ''}
                    onChange={e => setAnswers({ ...answers, [q.id]: e.target.value })}
                  />
                )}

                {q.type === 'NUMBER' && (
                  <Input type="number" value={answers[q.id] ?? ''} onChange={e => setAnswers({ ...answers, [q.id]: e.target.value })} />
                )}

                {q.type === 'RATING' && (
                  <div className="flex gap-1">
                    {[1, 2, 3, 4, 5].map((n) => (
                      <button
                        key={n}
                        onClick={() => setAnswers({ ...answers, [q.id]: String(n) })}
                        className={`h-9 w-9 rounded-lg border text-sm font-medium transition-colors ${
                          Number(answers[q.id]) === n
                            ? 'bg-brand-600 text-white border-brand-600'
                            : 'text-slate-600 border-slate-200 hover:bg-slate-50'
                        }`}
                      >
                        {n}
                      </button>
                    ))}
                  </div>
                )}

                {(q.type === 'SINGLE_CHOICE' || q.type === 'YES_NO') && (
                  <div className="space-y-1.5">
                    {(q.options ?? []).map((o) => (
                      <label key={o.id} className={`flex items-center gap-2 px-3 py-2 rounded-lg border text-sm cursor-pointer transition-colors ${answers[q.id] === o.id ? 'border-brand-600 bg-brand-50 text-brand-700' : 'border-slate-200 hover:bg-slate-50'}`}>
                        <input
                          type="radio"
                          name={q.id}
                          checked={answers[q.id] === o.id}
                          onChange={() => setAnswers({ ...answers, [q.id]: o.id })}
                          className="accent-brand-600"
                        />
                        {o.option_text}
                      </label>
                    ))}
                  </div>
                )}

                {q.type === 'MULTIPLE_CHOICE' && (
                  <div className="space-y-1.5">
                    {(q.options ?? []).map((o) => {
                      const selected = (answers[q.id] ?? []) as string[]
                      const isChecked = selected.includes(o.id)
                      return (
                        <label key={o.id} className={`flex items-center gap-2 px-3 py-2 rounded-lg border text-sm cursor-pointer transition-colors ${isChecked ? 'border-brand-600 bg-brand-50 text-brand-700' : 'border-slate-200 hover:bg-slate-50'}`}>
                          <input
                            type="checkbox"
                            checked={isChecked}
                            onChange={() =>
                              setAnswers({
                                ...answers,
                                [q.id]: isChecked ? selected.filter((x) => x !== o.id) : [...selected, o.id],
                              })
                            }
                            className="accent-brand-600"
                          />
                          {o.option_text}
                        </label>
                      )
                    })}
                  </div>
                )}
              </div>
            ))}
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setResponding(null)}>Cancelar</Button>
            <Button onClick={submit} disabled={submitting}>
              <Send size={14} /> {submitting ? 'Enviando...' : 'Enviar Respuesta'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function AdminTab() {
  const [surveys, setSurveys] = useState<Survey[]>([])
  const [loading, setLoading] = useState(true)
  const [showCreate, setShowCreate] = useState(false)
  const [saving, setSaving] = useState(false)
  const [actionLoading, setActionLoading] = useState<string | null>(null)
  const [editor, setEditor] = useState<Survey | null>(null)
  const [results, setResults] = useState<{ survey: Survey; stats: SurveyStats } | null>(null)

  const fetch = async () => {
    setLoading(true)
    try {
      const res = await api.get('/surveys', { params: { limit: '100' } })
      setSurveys(res.data.data ?? [])
    } catch { setSurveys([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetch() }, [])

  const [form, setForm] = useState({
    title: '',
    description: '',
    type: 'GENERAL',
    anonymous: false,
    multiple_responses: false,
    starts_at: '',
    ends_at: '',
  })

  const openCreate = () => {
    setForm({ title: '', description: '', type: 'GENERAL', anonymous: false, multiple_responses: false, starts_at: '', ends_at: '' })
    setShowCreate(true)
  }

  const handleCreate = async () => {
    setSaving(true)
    try {
      const payload: Record<string, any> = { title: form.title, type: form.type }
      if (form.description) payload.description = form.description
      if (form.anonymous) payload.anonymous = true
      if (form.multiple_responses) payload.multiple_responses = true
      if (form.starts_at) payload.starts_at = new Date(form.starts_at).toISOString()
      if (form.ends_at) payload.ends_at = new Date(form.ends_at).toISOString()
      await api.post('/surveys', payload)
      setShowCreate(false)
      fetch()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al crear la encuesta')
    } finally { setSaving(false) }
  }

  const handleAction = async (s: Survey, action: 'publish' | 'close' | 'archive' | 'delete') => {
    const label = { publish: 'publicar', close: 'cerrar', archive: 'archivar', delete: 'eliminar' }[action]
    if (action === 'delete' && !confirm(`¿Eliminar la encuesta "${s.title}"?`)) return
    setActionLoading(s.id)
    try {
      await api.post(`/surveys/${s.id}/${action}`)
      fetch()
    } catch (err: any) {
      alert(err?.response?.data?.error || `Error al ${label}`)
    } finally { setActionLoading(null) }
  }

  const exportCsv = async (s: Survey) => {
    try {
      const res = await api.get(`/surveys/${s.id}/export`, { responseType: 'blob' })
      const url = URL.createObjectURL(res.data)
      const a = document.createElement('a')
      a.href = url
      a.download = `encuesta_${s.id}.csv`
      a.click()
      URL.revokeObjectURL(url)
    } catch { alert('Error al exportar') }
  }

  const openResults = async (s: Survey) => {
    try {
      const res = await api.get(`/surveys/${s.id}/results`)
      setResults({ survey: s, stats: res.data.data })
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al cargar resultados')
    }
  }

  const openEditor = async (s: Survey) => {
    try {
      const res = await api.get(`/surveys/${s.id}`)
      setEditor(res.data.data ?? s)
    } catch { setEditor(s) }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-slate-500">{surveys.length} encuestas</p>
        <Button size="sm" onClick={openCreate}><Plus size={14} /> Nueva Encuesta</Button>
      </div>

      <Card>
        <CardContent className="p-0">
          {loading ? (
            <div className="p-6 text-center text-slate-500">Cargando...</div>
          ) : surveys.length === 0 ? (
            <div className="p-6 text-center text-slate-500">No hay encuestas creadas</div>
          ) : (
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-slate-200 bg-slate-50">
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Título</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Tipo</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Respuestas</th>
                  <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                </tr>
              </thead>
              <tbody>
                {surveys.map((s) => (
                  <tr key={s.id} className="border-b border-slate-100 hover:bg-slate-50">
                    <td className="px-4 py-3">
                      <p className="font-medium text-slate-900">{s.title}</p>
                      <p className="text-xs text-slate-400">{s.anonymous ? 'Anónima' : 'Nominativa'}</p>
                    </td>
                    <td className="px-4 py-3 text-slate-600">{surveyTypes[s.type] || s.type}</td>
                    <td className="px-4 py-3">
                      <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${statusColors[s.status] || 'bg-slate-100 text-slate-600'}`}>
                        {statusLabels[s.status] || s.status}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-slate-600">
                      {s.response_count ?? 0}
                      {s.participation_rate != null && (
                        <span className="text-xs text-slate-400"> ({Math.round(s.participation_rate)}%)</span>
                      )}
                    </td>
                    <td className="px-4 py-3 text-right whitespace-nowrap">
                      {s.status === 'DRAFT' && (
                        <>
                          <Button variant="ghost" size="sm" onClick={() => openEditor(s)}>Editar</Button>
                          <Button variant="ghost" size="sm" className="text-emerald-600" disabled={actionLoading === s.id} onClick={() => handleAction(s, 'publish')}>Publicar</Button>
                          <Button variant="ghost" size="sm" className="text-red-500" disabled={actionLoading === s.id} onClick={() => handleAction(s, 'delete')}>Eliminar</Button>
                        </>
                      )}
                      {s.status === 'PUBLISHED' && (
                        <>
                          <Button variant="ghost" size="sm" onClick={() => openResults(s)}>Resultados</Button>
                          <Button variant="ghost" size="sm" className="text-amber-600" disabled={actionLoading === s.id} onClick={() => handleAction(s, 'close')}>Cerrar</Button>
                        </>
                      )}
                      {(s.status === 'CLOSED') && (
                        <>
                          <Button variant="ghost" size="sm" onClick={() => openResults(s)}>Resultados</Button>
                          <Button variant="ghost" size="sm" onClick={() => exportCsv(s)}><FileDown size={14} /> CSV</Button>
                          <Button variant="ghost" size="sm" className="text-slate-500" disabled={actionLoading === s.id} onClick={() => handleAction(s, 'archive')}>Archivar</Button>
                        </>
                      )}
                      {s.status === 'ARCHIVED' && (
                        <>
                          <Button variant="ghost" size="sm" onClick={() => openResults(s)}>Resultados</Button>
                          <Button variant="ghost" size="sm" onClick={() => exportCsv(s)}><FileDown size={14} /> CSV</Button>
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
            <DialogTitle>Nueva Encuesta</DialogTitle>
            <DialogDescription>Definí los datos básicos; las preguntas se agregan después</DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label>Título *</Label>
              <Input value={form.title} onChange={e => setForm({ ...form, title: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label>Descripción</Label>
              <textarea
                className="flex min-h-[60px] w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm"
                value={form.description}
                onChange={e => setForm({ ...form, description: e.target.value })}
              />
            </div>
            <div className="space-y-2">
              <Label>Tipo *</Label>
              <Select
                options={Object.entries(surveyTypes).map(([value, label]) => ({ value, label }))}
                value={form.type}
                onChange={e => setForm({ ...form, type: e.target.value })}
              />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Inicio</Label>
                <Input type="datetime-local" value={form.starts_at} onChange={e => setForm({ ...form, starts_at: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label>Fin</Label>
                <Input type="datetime-local" value={form.ends_at} onChange={e => setForm({ ...form, ends_at: e.target.value })} />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4 text-sm">
              <label className="flex items-center gap-2">
                <input type="checkbox" checked={form.anonymous} onChange={e => setForm({ ...form, anonymous: e.target.checked })} className="rounded" />
                Anónima
              </label>
              <label className="flex items-center gap-2">
                <input type="checkbox" checked={form.multiple_responses} onChange={e => setForm({ ...form, multiple_responses: e.target.checked })} className="rounded" />
                Permitir múltiples respuestas
              </label>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowCreate(false)}>Cancelar</Button>
            <Button onClick={handleCreate} disabled={saving || !form.title}>
              {saving ? 'Creando...' : 'Crear'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {editor && <SurveyEditor survey={editor} onClose={() => { setEditor(null); fetch() }} />}
      {results && <ResultsDialog data={results} onClose={() => setResults(null)} />}
    </div>
  )
}

function SurveyEditor({ survey, onClose }: { survey: Survey; onClose: () => void }) {
  const [data, setData] = useState<Survey>(survey)
  const [questionForm, setQuestionForm] = useState<{ question: string; type: string; required: boolean }>({ question: '', type: 'TEXT', required: false })
  const [showQuestion, setShowQuestion] = useState(false)
  const [optionText, setOptionText] = useState<Record<string, string>>({})
  const [saving, setSaving] = useState(false)

  const [targetType, setTargetType] = useState('ALL')
  const [targetId, setTargetId] = useState('')
  const [targetsList, setTargetsList] = useState<SurveyTarget[]>(survey.targets ?? [])
  const [entityOptions, setEntityOptions] = useState<{ value: string; label: string }[]>([])

  const addQuestion = async () => {
    if (!questionForm.question.trim()) return
    setSaving(true)
    try {
      await api.post(`/surveys/${data.id}/questions`, {
        question: questionForm.question,
        type: questionForm.type,
        required: questionForm.required,
      })
      setShowQuestion(false)
      setQuestionForm({ question: '', type: 'TEXT', required: false })
      const res = await api.get(`/surveys/${data.id}`)
      setData(res.data.data)
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al agregar pregunta')
    } finally { setSaving(false) }
  }

  const deleteQuestion = async (q: SurveyQuestion) => {
    if (!confirm(`¿Eliminar la pregunta "${q.question}"?`)) return
    try {
      await api.delete(`/surveys/questions/${q.id}`)
      const res = await api.get(`/surveys/${data.id}`)
      setData(res.data.data)
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al eliminar pregunta')
    }
  }

  const addOption = async (q: SurveyQuestion) => {
    const text = optionText[q.id]?.trim()
    if (!text) return
    try {
      await api.post(`/surveys/questions/${q.id}/options`, { option_text: text })
      setOptionText((prev) => ({ ...prev, [q.id]: '' }))
      const res = await api.get(`/surveys/${data.id}`)
      setData(res.data.data)
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al agregar opción')
    }
  }

  const deleteOption = async (optionId: string) => {
    try {
      await api.delete(`/surveys/options/${optionId}`)
      const res = await api.get(`/surveys/${data.id}`)
      setData(res.data.data)
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al eliminar opción')
    }
  }

  const loadEntities = async (type: string) => {
    setTargetType(type)
    setTargetId('')
    if (type === 'ALL') { setEntityOptions([]); return }
    const path = type === 'DEPARTMENT' ? '/departments' : type === 'BRANCH' ? '/branches' : type === 'POSITION' ? '/positions' : '/employees'
    try {
      const res = await api.get(path, { params: { limit: '200' } })
      const items = res.data.data ?? []
      setEntityOptions(
        items.map((i: any) => ({
          value: i.id,
          label: type === 'EMPLOYEE' ? `${i.first_name} ${i.last_name}` : (i.name || `${i.first_name} ${i.last_name}`),
        })),
      )
    } catch { setEntityOptions([]) }
  }

  const addTarget = () => {
    const t: SurveyTarget = { id: `${Date.now()}`, target_type: targetType, target_id: targetType === 'ALL' ? undefined : (targetId || undefined) }
    setTargetsList((prev) => [...prev, t])
    setTargetId('')
  }

  const saveTargets = async () => {
    setSaving(true)
    try {
      await api.put(`/surveys/${data.id}/targets`, {
        targets: targetsList.map((t) => ({ target_type: t.target_type, target_id: t.target_id ?? null })),
      })
      alert('Destinatarios guardados')
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar destinatarios')
    } finally { setSaving(false) }
  }

  const choiceTypes = ['SINGLE_CHOICE', 'MULTIPLE_CHOICE', 'YES_NO']

  return (
    <Dialog open onOpenChange={(o) => !o && onClose()}>
      <DialogContent className="max-h-[90vh] flex flex-col">
        <DialogHeader>
          <DialogTitle>{data.title}</DialogTitle>
          <DialogDescription>
            {surveyTypes[data.type] || data.type} · {statusLabels[data.status] || data.status} · {data.anonymous ? 'Anónima' : 'Nominativa'}
          </DialogDescription>
        </DialogHeader>

        <div className="flex-1 overflow-y-auto space-y-6 pr-1">
          <div>
            <h4 className="font-semibold text-sm text-slate-900 mb-3">Preguntas ({data.questions?.length ?? 0})</h4>
            {(data.questions ?? []).length === 0 && (
              <p className="text-sm text-slate-500 mb-3">Todavía no hay preguntas. Agregá al menos una antes de publicar.</p>
            )}
            <div className="space-y-3">
              {(data.questions ?? []).map((q) => (
                <div key={q.id} className="rounded-lg border border-slate-200 p-3">
                  <div className="flex items-start justify-between gap-2">
                    <div>
                      <p className="text-sm font-medium text-slate-900">
                        {q.position}. {q.question} {q.required && <span className="text-red-500">*</span>}
                      </p>
                      <p className="text-xs text-slate-400 mt-0.5">{questionTypes[q.type] || q.type}</p>
                    </div>
                    <Button variant="ghost" size="sm" className="text-red-400" onClick={() => deleteQuestion(q)}>
                      <Trash2 size={14} />
                    </Button>
                  </div>

                  {choiceTypes.includes(q.type) && (
                    <div className="mt-2 space-y-1.5">
                      {(q.options ?? []).map((o) => (
                        <div key={o.id} className="flex items-center justify-between gap-2 px-3 py-1.5 rounded bg-slate-50 text-sm">
                          <span className="text-slate-700">{o.option_text}</span>
                          <button className="text-red-400 hover:text-red-600" onClick={() => deleteOption(o.id)}>
                            <XCircle size={14} />
                          </button>
                        </div>
                      ))}
                      <div className="flex gap-2">
                        <Input
                          placeholder="Nueva opción..."
                          value={optionText[q.id] ?? ''}
                          onChange={e => setOptionText((prev) => ({ ...prev, [q.id]: e.target.value }))}
                          onKeyDown={e => { if (e.key === 'Enter') addOption(q) }}
                        />
                        <Button variant="outline" size="sm" onClick={() => addOption(q)}>
                          <Plus size={14} />
                        </Button>
                      </div>
                    </div>
                  )}
                </div>
              ))}
            </div>
            <Button variant="outline" size="sm" className="mt-3" onClick={() => setShowQuestion(!showQuestion)}>
              <Plus size={14} /> Agregar Pregunta
            </Button>
          </div>

          {showQuestion && (
            <div className="rounded-lg border border-brand-200 bg-brand-50 p-4 space-y-3">
              <div className="space-y-2">
                <Label>Pregunta *</Label>
                <Input value={questionForm.question} onChange={e => setQuestionForm({ ...questionForm, question: e.target.value })} />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div className="space-y-2">
                  <Label>Tipo</Label>
                  <Select
                    options={Object.entries(questionTypes).map(([value, label]) => ({ value, label }))}
                    value={questionForm.type}
                    onChange={e => setQuestionForm({ ...questionForm, type: e.target.value })}
                  />
                </div>
                <label className="flex items-end gap-2 text-sm pb-2">
                  <input type="checkbox" checked={questionForm.required} onChange={e => setQuestionForm({ ...questionForm, required: e.target.checked })} className="rounded" />
                  Obligatoria
                </label>
              </div>
              <Button size="sm" onClick={addQuestion} disabled={saving || !questionForm.question.trim()}>
                <CheckCircle2 size={14} /> Agregar
              </Button>
            </div>
          )}

          <div>
            <h4 className="font-semibold text-sm text-slate-900 mb-3">Destinatarios</h4>
            <div className="flex gap-2 mb-2">
              <Select
                options={Object.entries(targetTypes).map(([value, label]) => ({ value, label }))}
                value={targetType}
                onChange={e => loadEntities(e.target.value)}
              />
              {targetType !== 'ALL' && (
                <Select
                  options={entityOptions}
                  placeholder="Seleccionar..."
                  value={targetId}
                  onChange={e => setTargetId(e.target.value)}
                />
              )}
              <Button variant="outline" size="sm" onClick={addTarget} disabled={targetType !== 'ALL' && !targetId}>
                <Plus size={14} />
              </Button>
            </div>
            {targetsList.length > 0 && (
              <div className="space-y-1.5 mb-2">
                {targetsList.map((t) => (
                  <div key={t.id} className="flex items-center justify-between px-3 py-1.5 rounded bg-slate-50 text-sm">
                    <span className="text-slate-700">
                      {targetTypes[t.target_type] || t.target_type}{t.target_id ? `: ${t.target_id}` : ''}
                    </span>
                    <button className="text-red-400 hover:text-red-600" onClick={() => setTargetsList((prev) => prev.filter((x) => x.id !== t.id))}>
                      <Trash2 size={14} />
                    </button>
                  </div>
                ))}
              </div>
            )}
            <Button size="sm" variant="outline" onClick={saveTargets} disabled={saving}>
              Guardar Destinatarios
            </Button>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={onClose}>Cerrar</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

function ResultsDialog({ data, onClose }: { data: { survey: Survey; stats: SurveyStats }; onClose: () => void }) {
  const { survey, stats } = data
  const maxCount = (dist: OptionDistribution[]) => Math.max(1, ...dist.map((d) => d.count))

  return (
    <Dialog open onOpenChange={(o) => !o && onClose()}>
      <DialogContent className="max-h-[90vh] flex flex-col">
        <DialogHeader>
          <DialogTitle>{survey.title} — Resultados</DialogTitle>
          <DialogDescription>
            {stats.total_responded} de {stats.total_targeted} destinatarios · {Math.round(stats.participation_rate)}% de participación
          </DialogDescription>
        </DialogHeader>

        <div className="flex-1 overflow-y-auto space-y-5 pr-1">
          {stats.questions.length === 0 ? (
            <p className="text-sm text-slate-500">Todavía no hay respuestas</p>
          ) : (
            stats.questions.map((q) => (
              <div key={q.question_id} className="rounded-lg border border-slate-200 p-4">
                <p className="text-sm font-medium text-slate-900">{q.question}</p>
                <p className="text-xs text-slate-400 mb-3">{questionTypes[q.type] || q.type} · {q.total_answers} respuestas</p>

                {(q.type === 'NUMBER' || q.type === 'RATING') && q.average != null && (
                  <div className="grid grid-cols-3 gap-3 text-center">
                    <div className="rounded-lg bg-slate-50 p-3">
                      <p className="text-xs text-slate-500">Promedio</p>
                      <p className="text-lg font-bold text-brand-700">{q.average.toFixed(1)}</p>
                    </div>
                    <div className="rounded-lg bg-slate-50 p-3">
                      <p className="text-xs text-slate-500">Mínimo</p>
                      <p className="text-lg font-bold text-slate-700">{q.min}</p>
                    </div>
                    <div className="rounded-lg bg-slate-50 p-3">
                      <p className="text-xs text-slate-500">Máximo</p>
                      <p className="text-lg font-bold text-slate-700">{q.max}</p>
                    </div>
                  </div>
                )}

                {q.type === 'YES_NO' && q.yes_count != null && (
                  <div className="flex gap-6 text-sm">
                    <span className="text-emerald-600 font-medium">Sí: {q.yes_count} ({Math.round(q.yes_percentage ?? 0)}%)</span>
                    <span className="text-red-500 font-medium">No: {q.no_count}</span>
                  </div>
                )}

                {(q.distribution ?? []).length > 0 && (
                  <div className="space-y-2">
                    {q.distribution!.map((d) => (
                      <div key={d.option_id}>
                        <div className="flex justify-between text-xs mb-1">
                          <span className="text-slate-700">{d.option_text}</span>
                          <span className="text-slate-500">{d.count} ({Math.round(d.percentage)}%)</span>
                        </div>
                        <div className="h-2 bg-slate-100 rounded-full overflow-hidden">
                          <div
                            className="h-full bg-brand-500 rounded-full"
                            style={{ width: `${(d.count / maxCount(q.distribution!)) * 100}%` }}
                          />
                        </div>
                      </div>
                    ))}
                  </div>
                )}

                {(q.sample_texts ?? []).length > 0 && (
                  <div className="mt-3 space-y-1">
                    {q.sample_texts!.slice(0, 5).map((t, i) => (
                      <p key={i} className="text-sm text-slate-600 rounded bg-slate-50 px-3 py-2">“{t}”</p>
                    ))}
                  </div>
                )}
              </div>
            ))
          )}
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={onClose}>Cerrar</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
