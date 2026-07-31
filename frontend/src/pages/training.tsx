import { useEffect, useState } from 'react'
import { Plus, Pencil, Rocket } from 'lucide-react'
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

interface Course {
  id: string
  code: string
  name: string
  short_description?: string
  category_id?: string
  difficulty?: string
  modality?: string
  status: string
  duration_minutes?: number
  mandatory?: boolean
  passing_score?: number
  certificate_enabled?: boolean
}

interface SelectOption {
  value: string
  label: string
}

const emptyForm = {
  code: '',
  name: '',
  category_id: '',
  short_description: '',
  description: '',
  objectives: '',
  difficulty: 'BEGINNER',
  duration_minutes: '60',
  modality: 'ONLINE',
  mandatory: false,
  passing_score: '',
  certificate_enabled: false,
  min_attendance_percentage: '80',
}

const difficultyOptions: SelectOption[] = [
  { value: 'BEGINNER', label: 'Principiante' },
  { value: 'INTERMEDIATE', label: 'Intermedio' },
  { value: 'ADVANCED', label: 'Avanzado' },
  { value: 'EXPERT', label: 'Experto' },
]

const modalityOptions: SelectOption[] = [
  { value: 'ONLINE', label: 'Online' },
  { value: 'PRESENTIAL', label: 'Presencial' },
  { value: 'HYBRID', label: 'Híbrido' },
  { value: 'ASYNC', label: 'Asincrónico' },
]

export default function TrainingPage() {
  const [items, setItems] = useState<Course[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showModal, setShowModal] = useState(false)
  const [editing, setEditing] = useState<Course | null>(null)
  const [form, setForm] = useState({ ...emptyForm })
  const [saving, setSaving] = useState(false)
  const [categories, setCategories] = useState<SelectOption[]>([])

  const fetchData = async () => {
    setLoading(true)
    try {
      const res = await api.get('/training/courses', { params: { limit: '100' } })
      setItems(res.data.data ?? [])
      setError('')
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar cursos')
      setItems([])
    } finally {
      setLoading(false)
    }
  }

  const fetchCategories = async () => {
    try {
      const res = await api.get('/training/categories')
      setCategories((res.data.data ?? res.data ?? []).map((c: any) => ({ value: c.id, label: c.name })))
    } catch {}
  }

  useEffect(() => { fetchData() }, [])

  const openCreate = () => {
    setEditing(null)
    setForm({ ...emptyForm })
    fetchCategories()
    setShowModal(true)
  }

  const openEdit = (c: Course) => {
    setEditing(c)
    setForm({
      code: c.code,
      name: c.name,
      category_id: c.category_id ?? '',
      short_description: c.short_description ?? '',
      description: '',
      objectives: '',
      difficulty: c.difficulty || 'BEGINNER',
      duration_minutes: c.duration_minutes != null ? String(c.duration_minutes) : '60',
      modality: c.modality || 'ONLINE',
      mandatory: c.mandatory ?? false,
      passing_score: c.passing_score != null ? String(c.passing_score) : '',
      certificate_enabled: c.certificate_enabled ?? false,
      min_attendance_percentage: '80',
    })
    fetchCategories()
    setShowModal(true)
  }

  const handleSave = async () => {
    setSaving(true)
    try {
      const body: Record<string, any> = {
        code: form.code,
        name: form.name,
        difficulty: form.difficulty,
        modality: form.modality,
        duration_minutes: parseInt(form.duration_minutes || '0', 10),
        mandatory: form.mandatory,
        certificate_enabled: form.certificate_enabled,
      }
      if (form.category_id) body.category_id = form.category_id
      if (form.short_description) body.short_description = form.short_description
      if (form.description) body.description = form.description
      if (form.objectives) body.objectives = form.objectives
      if (form.passing_score) body.passing_score = parseFloat(form.passing_score)
      if (form.min_attendance_percentage) body.min_attendance_percentage = parseFloat(form.min_attendance_percentage)
      if (editing) {
        await api.put(`/training/courses/${editing.id}`, body)
      } else {
        await api.post('/training/courses', body)
      }
      setShowModal(false)
      fetchData()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar curso')
    } finally {
      setSaving(false)
    }
  }

  const handlePublish = async (c: Course) => {
    try {
      await api.post(`/training/courses/${c.id}/publish`)
      fetchData()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al publicar curso')
    }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-slate-900">Capacitación</h1>
        <Button size="sm" onClick={openCreate}><Plus size={16} className="mr-1" /> Nuevo curso</Button>
      </div>

      {error && <div className="mb-4 p-3 bg-red-50 text-red-700 text-sm rounded-lg">{error}</div>}

      <Card>
        <CardContent className="p-0">
          {loading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
          : items.length === 0 && !error ? <div className="p-6 text-center text-slate-500">No hay cursos registrados</div>
          : <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-slate-200 bg-slate-50">
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Código</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Modalidad</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Duración</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                    <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                  </tr>
                </thead>
                <tbody>
                  {items.map(c => (
                    <tr key={c.id} className="border-b border-slate-100 hover:bg-slate-50">
                      <td className="px-4 py-3 text-slate-500">{c.code}</td>
                      <td className="px-4 py-3 font-medium text-slate-900">{c.name}</td>
                      <td className="px-4 py-3 text-slate-600 capitalize">{(c.modality || '').toLowerCase()}</td>
                      <td className="px-4 py-3 text-slate-600">{c.duration_minutes ? `${Math.round(c.duration_minutes / 60)}h` : '-'}</td>
                      <td className="px-4 py-3">
                        <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${c.status === 'PUBLISHED' ? 'bg-emerald-50 text-emerald-700' : 'bg-slate-100 text-slate-600'}`}>
                          {c.status === 'PUBLISHED' ? 'Publicado' : c.status === 'DRAFT' ? 'Borrador' : c.status}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-right whitespace-nowrap">
                        {c.status !== 'PUBLISHED' && (
                          <Button variant="ghost" size="sm" className="text-emerald-600" onClick={() => handlePublish(c)} title="Publicar">
                            <Rocket size={14} />
                          </Button>
                        )}
                        <Button variant="ghost" size="sm" onClick={() => openEdit(c)}><Pencil size={14} /></Button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>}
        </CardContent>
      </Card>

      <Dialog open={showModal} onOpenChange={setShowModal}>
        <DialogContent className="max-w-xl">
          <DialogHeader>
            <DialogTitle>{editing ? 'Editar Curso' : 'Nuevo Curso'}</DialogTitle>
            <DialogDescription>
              {editing ? 'Modificá los datos del curso' : 'Completá los datos para registrar un nuevo curso'}
            </DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="course-code">Código *</Label>
              <Input id="course-code" value={form.code} onChange={e => setForm({ ...form, code: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="course-name">Nombre *</Label>
              <Input id="course-name" value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} required />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="course-short">Descripción corta</Label>
              <Input id="course-short" value={form.short_description} onChange={e => setForm({ ...form, short_description: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="course-desc">Descripción</Label>
              <textarea
                id="course-desc"
                value={form.description}
                onChange={e => setForm({ ...form, description: e.target.value })}
                rows={3}
                className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-brand-500"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="course-cat">Categoría</Label>
              <Select id="course-cat" options={categories} placeholder="Seleccionar..." value={form.category_id} onChange={e => setForm({ ...form, category_id: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="course-difficulty">Dificultad</Label>
              <Select id="course-difficulty" options={difficultyOptions} value={form.difficulty} onChange={e => setForm({ ...form, difficulty: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="course-modality">Modalidad</Label>
              <Select id="course-modality" options={modalityOptions} value={form.modality} onChange={e => setForm({ ...form, modality: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="course-duration">Duración (minutos)</Label>
              <Input id="course-duration" type="number" min={0} value={form.duration_minutes} onChange={e => setForm({ ...form, duration_minutes: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="course-passing">Nota de aprobación</Label>
              <Input id="course-passing" type="number" step="0.01" value={form.passing_score} onChange={e => setForm({ ...form, passing_score: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <label className="flex items-center gap-2 text-sm text-slate-700">
                <input type="checkbox" checked={form.mandatory} onChange={e => setForm({ ...form, mandatory: e.target.checked })} className="h-4 w-4 rounded border-slate-300 text-brand-600 focus:ring-brand-500" />
                Curso obligatorio
              </label>
            </div>
            <div className="space-y-2 col-span-2">
              <label className="flex items-center gap-2 text-sm text-slate-700">
                <input type="checkbox" checked={form.certificate_enabled} onChange={e => setForm({ ...form, certificate_enabled: e.target.checked })} className="h-4 w-4 rounded border-slate-300 text-brand-600 focus:ring-brand-500" />
                Emitir certificado
              </label>
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button>
            <Button onClick={handleSave} disabled={saving || !form.code || !form.name}>
              {saving ? 'Guardando...' : editing ? 'Guardar Cambios' : 'Crear Curso'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
