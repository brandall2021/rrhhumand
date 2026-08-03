import { useEffect, useState } from 'react'
import { Plus, Pencil, Rocket, Eye, Layers } from 'lucide-react'
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
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'

interface Course {
  id: string
  code: string
  name: string
  category_id?: string
  short_description?: string
  description?: string
  objectives?: string
  difficulty?: string
  duration_minutes?: number
  modality?: string
  status: string
  mandatory?: boolean
  passing_score?: number
  certificate_enabled?: boolean
  min_attendance_percentage?: number
}

interface Category {
  id: string
  name: string
  description?: string
  active: boolean
  created_at: string
}

interface CourseVersion {
  id: string
  course_id: string
  version: string
  description?: string
  status: string
  is_published: boolean
  created_by?: string
  published_at?: string
  created_at: string
}

interface CourseContent {
  id: string
  course_version_id: string
  title: string
  description?: string
  content_type: string
  storage_provider?: string
  storage_key?: string
  external_url?: string
  duration_seconds: number
  sort_order: number
  required: boolean
  published: boolean
  created_at: string
}

interface TrainingStats {
  total_courses: number
  total_enrollments: number
  completed_enrollments: number
  in_progress_enrollments: number
  active_offerings: number
  average_rating: number
}

interface SelectOption {
  value: string
  label: string
}

const emptyCourseForm = {
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

const emptyCatForm = {
  name: '',
  description: '',
}

const emptyStats: TrainingStats = {
  total_courses: 0,
  total_enrollments: 0,
  completed_enrollments: 0,
  in_progress_enrollments: 0,
  active_offerings: 0,
  average_rating: 0,
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

const statusLabel = (status: string) => {
  const s = (status || '').toLowerCase()
  return s === 'published' ? 'Publicado' : s === 'draft' ? 'Borrador' : s
}

export default function TrainingPage() {
  const [items, setItems] = useState<Course[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [search, setSearch] = useState('')
  const [stats, setStats] = useState<TrainingStats>({ ...emptyStats })

  const [showModal, setShowModal] = useState(false)
  const [editing, setEditing] = useState<Course | null>(null)
  const [courseForm, setCourseForm] = useState({ ...emptyCourseForm })
  const [saving, setSaving] = useState(false)

  const [categories, setCategories] = useState<SelectOption[]>([])
  const [categoryMap, setCategoryMap] = useState<Record<string, string>>({})

  const [showDetail, setShowDetail] = useState(false)
  const [detailCourse, setDetailCourse] = useState<Course | null>(null)
  const [detailVersions, setDetailVersions] = useState<CourseVersion[]>([])
  const [detailContents, setDetailContents] = useState<Record<string, CourseContent[]>>({})
  const [versionForm, setVersionForm] = useState({ version: '', description: '' })
  const [creatingVersion, setCreatingVersion] = useState(false)

  const [catItems, setCatItems] = useState<Category[]>([])
  const [catLoading, setCatLoading] = useState(true)
  const [showCatModal, setShowCatModal] = useState(false)
  const [editingCat, setEditingCat] = useState<Category | null>(null)
  const [catForm, setCatForm] = useState({ ...emptyCatForm })
  const [savingCat, setSavingCat] = useState(false)

  const fetchData = async () => {
    setLoading(true)
    try {
      const params: Record<string, string> = { limit: '100' }
      if (search) params.search = search
      const res = await api.get('/training/courses', { params })
      setItems(res.data.data ?? [])
      setError('')
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar cursos')
      setItems([])
    } finally {
      setLoading(false)
    }
  }

  const fetchStats = async () => {
    try {
      const res = await api.get('/training/dashboard/stats')
      setStats({ ...emptyStats, ...(res.data?.data ?? res.data ?? {}) })
    } catch {
      setStats({ ...emptyStats })
    }
  }

  const fetchCategories = async () => {
    try {
      const res = await api.get('/training/categories')
      const list = res.data?.data ?? res.data ?? []
      setCategories(list.map((c: any) => ({ value: c.id, label: c.name })))
      const map: Record<string, string> = {}
      list.forEach((c: any) => { map[c.id] = c.name })
      setCategoryMap(map)
      return list
    } catch {
      return []
    }
  }

  const fetchCatItems = async () => {
    setCatLoading(true)
    try {
      const res = await api.get('/training/categories')
      setCatItems(res.data?.data ?? res.data ?? [])
    } catch {
      setCatItems([])
    } finally {
      setCatLoading(false)
    }
  }

  useEffect(() => { fetchData() }, [search])
  useEffect(() => {
    fetchStats()
    fetchCategories()
    fetchCatItems()
  }, [])

  const refreshDetail = async (courseId: string) => {
    try {
      const [dRes, vRes] = await Promise.all([
        api.get(`/training/courses/${courseId}/details`),
        api.get(`/training/courses/${courseId}/versions`),
      ])
      const details = dRes.data?.data ?? dRes.data ?? {}
      const versions: CourseVersion[] = vRes.data?.data ?? vRes.data ?? details.versions ?? []
      setDetailCourse(prev => prev ? { ...prev, ...(details.course ?? {}) } : prev)
      setDetailVersions(versions)
      const map: Record<string, CourseContent[]> = {}
      if (versions.length) {
        const lists = await Promise.all(
          versions.map(v =>
            api.get(`/training/versions/${v.id}/contents`)
              .then(r => ({ vid: v.id, list: (r.data?.data ?? r.data ?? []) as CourseContent[] }))
              .catch(() => ({ vid: v.id, list: [] as CourseContent[] })),
          ),
        )
        lists.forEach(({ vid, list }) => { map[vid] = list })
      }
      if (!versions.length && Array.isArray(details.contents)) {
        details.contents.forEach((x: CourseContent) => {
          map[x.course_version_id] = [...(map[x.course_version_id] ?? []), x]
        })
      }
      setDetailContents(map)
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al cargar el detalle del curso')
    }
  }

  const openDetail = (c: Course) => {
    setDetailCourse(c)
    setDetailVersions([])
    setDetailContents({})
    setVersionForm({ version: '', description: '' })
    setShowDetail(true)
    refreshDetail(c.id)
  }

  const handleCreateVersion = async () => {
    if (!detailCourse || !versionForm.version.trim()) return
    setCreatingVersion(true)
    try {
      const body: Record<string, any> = { version: versionForm.version.trim() }
      if (versionForm.description) body.description = versionForm.description
      await api.post(`/training/courses/${detailCourse.id}/versions`, body)
      setVersionForm({ version: '', description: '' })
      await refreshDetail(detailCourse.id)
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al crear la versión')
    } finally {
      setCreatingVersion(false)
    }
  }

  const openCreate = () => {
    setEditing(null)
    setCourseForm({ ...emptyCourseForm })
    fetchCategories()
    setShowModal(true)
  }

  const openEdit = async (c: Course) => {
    setEditing(c)
    setCourseForm({
      code: c.code,
      name: c.name,
      category_id: c.category_id ?? '',
      short_description: c.short_description ?? '',
      description: c.description ?? '',
      objectives: c.objectives ?? '',
      difficulty: c.difficulty || 'BEGINNER',
      duration_minutes: c.duration_minutes != null ? String(c.duration_minutes) : '60',
      modality: c.modality || 'ONLINE',
      mandatory: c.mandatory ?? false,
      passing_score: c.passing_score != null ? String(c.passing_score) : '',
      certificate_enabled: c.certificate_enabled ?? false,
      min_attendance_percentage: c.min_attendance_percentage != null ? String(c.min_attendance_percentage) : '80',
    })
    fetchCategories()
    setShowModal(true)
    try {
      const res = await api.get(`/training/courses/${c.id}`)
      const full = res.data?.data ?? res.data ?? c
      setCourseForm({
        code: full.code ?? c.code,
        name: full.name ?? c.name,
        category_id: full.category_id ?? c.category_id ?? '',
        short_description: full.short_description ?? '',
        description: full.description ?? '',
        objectives: full.objectives ?? '',
        difficulty: full.difficulty || 'BEGINNER',
        duration_minutes: full.duration_minutes != null ? String(full.duration_minutes) : '60',
        modality: full.modality || 'ONLINE',
        mandatory: full.mandatory ?? false,
        passing_score: full.passing_score != null ? String(full.passing_score) : '',
        certificate_enabled: full.certificate_enabled ?? false,
        min_attendance_percentage: full.min_attendance_percentage != null ? String(full.min_attendance_percentage) : '80',
      })
    } catch {}
  }

  const handleSaveCourse = async () => {
    setSaving(true)
    try {
      const body: Record<string, any> = {
        code: courseForm.code,
        name: courseForm.name,
        difficulty: courseForm.difficulty,
        modality: courseForm.modality,
        duration_minutes: parseInt(courseForm.duration_minutes || '0', 10),
        mandatory: courseForm.mandatory,
        certificate_enabled: courseForm.certificate_enabled,
      }
      if (courseForm.category_id) body.category_id = courseForm.category_id
      if (courseForm.short_description) body.short_description = courseForm.short_description
      if (courseForm.description) body.description = courseForm.description
      if (courseForm.objectives) body.objectives = courseForm.objectives
      if (courseForm.passing_score) body.passing_score = parseFloat(courseForm.passing_score)
      if (courseForm.min_attendance_percentage) body.min_attendance_percentage = parseFloat(courseForm.min_attendance_percentage)
      if (editing) {
        await api.put(`/training/courses/${editing.id}`, body)
      } else {
        await api.post('/training/courses', body)
      }
      setShowModal(false)
      fetchData()
      fetchStats()
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
      fetchStats()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al publicar curso')
    }
  }

  const openCreateCat = () => {
    setEditingCat(null)
    setCatForm({ ...emptyCatForm })
    setShowCatModal(true)
  }

  const openEditCat = (c: Category) => {
    setEditingCat(c)
    setCatForm({ name: c.name, description: c.description ?? '' })
    setShowCatModal(true)
  }

  const handleSaveCat = async () => {
    setSavingCat(true)
    try {
      const body: Record<string, any> = { name: catForm.name }
      if (catForm.description) body.description = catForm.description
      if (editingCat) {
        body.active = editingCat.active
        await api.put(`/training/categories/${editingCat.id}`, body)
      } else {
        await api.post('/training/categories', body)
      }
      setShowCatModal(false)
      fetchCatItems()
      fetchCategories()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar categoría')
    } finally {
      setSavingCat(false)
    }
  }

  const statCards: { label: string; value: number | string }[] = [
    { label: 'Cursos totales', value: stats.total_courses },
    { label: 'Inscripciones', value: stats.total_enrollments },
    { label: 'Completadas', value: stats.completed_enrollments },
    { label: 'En curso', value: stats.in_progress_enrollments },
    { label: 'Ofertas activas', value: stats.active_offerings },
    { label: 'Rating promedio', value: stats.average_rating ? stats.average_rating.toFixed(1) : '-' },
  ]

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-slate-900">Capacitación</h1>
      </div>

      {error && <div className="mb-4 p-3 bg-red-50 text-red-700 text-sm rounded-lg">{error}</div>}

      <Tabs defaultValue="courses">
        <TabsList>
          <TabsTrigger value="courses">Cursos</TabsTrigger>
          <TabsTrigger value="categories">Categorías</TabsTrigger>
        </TabsList>

        <TabsContent value="courses">
          <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-4 mb-4">
            {statCards.map(s => (
              <Card key={s.label}>
                <CardContent className="p-4">
                  <div className="text-xs text-slate-500">{s.label}</div>
                  <div className="text-xl font-bold text-slate-900 mt-1">{s.value}</div>
                </CardContent>
              </Card>
            ))}
          </div>

          <div className="flex items-center gap-3 mb-4">
            <Input
              placeholder="Buscar por código o nombre..."
              value={search}
              onChange={e => setSearch(e.target.value)}
              className="max-w-xs"
            />
            <div className="ml-auto">
              <Button size="sm" onClick={openCreate}><Plus size={16} className="mr-1" /> Nuevo curso</Button>
            </div>
          </div>

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
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Categoría</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Modalidad</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Duración</th>
                        <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                      </tr>
                    </thead>
                    <tbody>
                      {items.map(c => {
                        const st = (c.status || '').toLowerCase()
                        return (
                          <tr key={c.id} className="border-b border-slate-100 hover:bg-slate-50">
                            <td className="px-4 py-3 text-slate-500">{c.code}</td>
                            <td className="px-4 py-3 font-medium text-slate-900">{c.name}</td>
                            <td className="px-4 py-3 text-slate-600">{c.category_id ? (categoryMap[c.category_id] || '-') : '-'}</td>
                            <td className="px-4 py-3 text-slate-600 capitalize">{(c.modality || '').toLowerCase()}</td>
                            <td className="px-4 py-3">
                              <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${st === 'published' ? 'bg-emerald-50 text-emerald-700' : 'bg-slate-100 text-slate-600'}`}>
                                {statusLabel(st)}
                              </span>
                            </td>
                            <td className="px-4 py-3 text-slate-600">{c.duration_minutes ? `${Math.round(c.duration_minutes / 60)}h` : '-'}</td>
                            <td className="px-4 py-3 text-right whitespace-nowrap">
                              <Button variant="ghost" size="sm" onClick={() => openDetail(c)} title="Detalle"><Eye size={14} /></Button>
                              {st !== 'published' && (
                                <Button variant="ghost" size="sm" className="text-emerald-600" onClick={() => handlePublish(c)} title="Publicar"><Rocket size={14} /></Button>
                              )}
                              <Button variant="ghost" size="sm" onClick={() => openEdit(c)} title="Editar"><Pencil size={14} /></Button>
                            </td>
                          </tr>
                        )
                      })}
                    </tbody>
                  </table>
                </div>}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="categories">
          <div className="flex items-center justify-between mb-4">
            <span className="text-sm text-slate-500">{catItems.length} categorías</span>
            <Button size="sm" onClick={openCreateCat}><Plus size={16} className="mr-1" /> Nueva</Button>
          </div>

          <Card>
            <CardContent className="p-0">
              {catLoading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
              : catItems.length === 0 ? <div className="p-6 text-center text-slate-500">No hay categorías</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Descripción</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                        <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                      </tr>
                    </thead>
                    <tbody>
                      {catItems.map(c => (
                        <tr key={c.id} className="border-b border-slate-100 hover:bg-slate-50">
                          <td className="px-4 py-3 font-medium text-slate-900">{c.name}</td>
                          <td className="px-4 py-3 text-slate-600">{c.description || '-'}</td>
                          <td className="px-4 py-3">
                            <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${c.active ? 'bg-emerald-50 text-emerald-700' : 'bg-slate-100 text-slate-600'}`}>
                              {c.active ? 'Activa' : 'Inactiva'}
                            </span>
                          </td>
                          <td className="px-4 py-3 text-right whitespace-nowrap">
                            <Button variant="ghost" size="sm" onClick={() => openEditCat(c)} title="Editar"><Pencil size={14} /></Button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      <Dialog open={showModal} onOpenChange={setShowModal}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>{editing ? 'Editar Curso' : 'Nuevo Curso'}</DialogTitle>
            <DialogDescription>
              {editing ? 'Modificá los datos del curso' : 'Completá los datos para registrar un nuevo curso'}
            </DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="course-code">Código *</Label>
              <Input id="course-code" value={courseForm.code} disabled={!!editing} onChange={e => setCourseForm({ ...courseForm, code: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="course-name">Nombre *</Label>
              <Input id="course-name" value={courseForm.name} onChange={e => setCourseForm({ ...courseForm, name: e.target.value })} required />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="course-short">Descripción corta</Label>
              <Input id="course-short" value={courseForm.short_description} onChange={e => setCourseForm({ ...courseForm, short_description: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="course-desc">Descripción</Label>
              <textarea
                id="course-desc"
                value={courseForm.description}
                onChange={e => setCourseForm({ ...courseForm, description: e.target.value })}
                rows={3}
                className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-brand-500"
              />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="course-obj">Objetivos</Label>
              <textarea
                id="course-obj"
                value={courseForm.objectives}
                onChange={e => setCourseForm({ ...courseForm, objectives: e.target.value })}
                rows={3}
                className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-brand-500"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="course-cat">Categoría</Label>
              <Select id="course-cat" options={categories} placeholder="Seleccionar..." value={courseForm.category_id} onChange={e => setCourseForm({ ...courseForm, category_id: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="course-difficulty">Dificultad</Label>
              <Select id="course-difficulty" options={difficultyOptions} value={courseForm.difficulty} onChange={e => setCourseForm({ ...courseForm, difficulty: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="course-modality">Modalidad</Label>
              <Select id="course-modality" options={modalityOptions} value={courseForm.modality} onChange={e => setCourseForm({ ...courseForm, modality: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="course-duration">Duración (minutos)</Label>
              <Input id="course-duration" type="number" min={0} value={courseForm.duration_minutes} onChange={e => setCourseForm({ ...courseForm, duration_minutes: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="course-passing">Nota de aprobación</Label>
              <Input id="course-passing" type="number" step="0.01" value={courseForm.passing_score} onChange={e => setCourseForm({ ...courseForm, passing_score: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="course-attendance">Asistencia mínima (%)</Label>
              <Input id="course-attendance" type="number" step="0.01" min={0} max={100} value={courseForm.min_attendance_percentage} onChange={e => setCourseForm({ ...courseForm, min_attendance_percentage: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <label className="flex items-center gap-2 text-sm text-slate-700">
                <input type="checkbox" checked={courseForm.mandatory} onChange={e => setCourseForm({ ...courseForm, mandatory: e.target.checked })} className="h-4 w-4 rounded border-slate-300 text-brand-600 focus:ring-brand-500" />
                Curso obligatorio
              </label>
            </div>
            <div className="space-y-2 col-span-2">
              <label className="flex items-center gap-2 text-sm text-slate-700">
                <input type="checkbox" checked={courseForm.certificate_enabled} onChange={e => setCourseForm({ ...courseForm, certificate_enabled: e.target.checked })} className="h-4 w-4 rounded border-slate-300 text-brand-600 focus:ring-brand-500" />
                Emitir certificado
              </label>
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button>
            <Button onClick={handleSaveCourse} disabled={saving || !courseForm.code || !courseForm.name}>
              {saving ? 'Guardando...' : editing ? 'Guardar Cambios' : 'Crear Curso'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={showDetail} onOpenChange={setShowDetail}>
        <DialogContent className="max-w-3xl">
          <DialogHeader>
            <DialogTitle>Detalle del curso</DialogTitle>
            <DialogDescription>{detailCourse ? `${detailCourse.code} · ${detailCourse.name}` : ''}</DialogDescription>
          </DialogHeader>

          {detailCourse && (
            <div className="space-y-6">
              <div className="grid grid-cols-2 sm:grid-cols-3 gap-4 text-sm">
                <div>
                  <div className="text-xs text-slate-500">Categoría</div>
                  <div className="text-slate-900">{detailCourse.category_id ? (categoryMap[detailCourse.category_id] || '-') : '-'}</div>
                </div>
                <div>
                  <div className="text-xs text-slate-500">Modalidad</div>
                  <div className="text-slate-900 capitalize">{(detailCourse.modality || '').toLowerCase()}</div>
                </div>
                <div>
                  <div className="text-xs text-slate-500">Dificultad</div>
                  <div className="text-slate-900 capitalize">{(detailCourse.difficulty || '').toLowerCase()}</div>
                </div>
                <div>
                  <div className="text-xs text-slate-500">Estado</div>
                  <div className="text-slate-900">{statusLabel(detailCourse.status)}</div>
                </div>
                <div>
                  <div className="text-xs text-slate-500">Duración</div>
                  <div className="text-slate-900">{detailCourse.duration_minutes ? `${Math.round(detailCourse.duration_minutes / 60)}h` : '-'}</div>
                </div>
                <div>
                  <div className="text-xs text-slate-500">Nota de aprobación</div>
                  <div className="text-slate-900">{detailCourse.passing_score != null ? detailCourse.passing_score : '-'}</div>
                </div>
                <div>
                  <div className="text-xs text-slate-500">Obligatorio</div>
                  <div className="text-slate-900">{detailCourse.mandatory ? 'Sí' : 'No'}</div>
                </div>
                <div>
                  <div className="text-xs text-slate-500">Certificado</div>
                  <div className="text-slate-900">{detailCourse.certificate_enabled ? 'Sí' : 'No'}</div>
                </div>
                <div>
                  <div className="text-xs text-slate-500">Asistencia mínima</div>
                  <div className="text-slate-900">{detailCourse.min_attendance_percentage != null ? `${detailCourse.min_attendance_percentage}%` : '-'}</div>
                </div>
              </div>

              {detailCourse.short_description && (
                <div>
                  <div className="text-xs text-slate-500 mb-1">Descripción corta</div>
                  <div className="text-sm text-slate-700">{detailCourse.short_description}</div>
                </div>
              )}
              {detailCourse.description && (
                <div>
                  <div className="text-xs text-slate-500 mb-1">Descripción</div>
                  <div className="text-sm text-slate-700">{detailCourse.description}</div>
                </div>
              )}
              {detailCourse.objectives && (
                <div>
                  <div className="text-xs text-slate-500 mb-1">Objetivos</div>
                  <div className="text-sm text-slate-700 whitespace-pre-wrap">{detailCourse.objectives}</div>
                </div>
              )}

              <div>
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-2 text-sm font-medium text-slate-900">
                    <Layers size={16} className="text-slate-400" />
                    Versiones ({detailVersions.length})
                  </div>
                  <Button size="sm" onClick={handleCreateVersion} disabled={creatingVersion || !versionForm.version.trim()}>
                    <Plus size={14} className="mr-1" /> {creatingVersion ? 'Creando...' : 'Crear versión'}
                  </Button>
                </div>

                <div className="flex gap-3 mb-3">
                  <Input
                    placeholder="Versión (ej. 1.0)"
                    value={versionForm.version}
                    onChange={e => setVersionForm({ ...versionForm, version: e.target.value })}
                    className="max-w-[140px]"
                  />
                  <Input
                    placeholder="Descripción de la versión"
                    value={versionForm.description}
                    onChange={e => setVersionForm({ ...versionForm, description: e.target.value })}
                  />
                </div>

                <div className="space-y-4">
                  {detailVersions.length === 0 ? (
                    <div className="p-4 text-center text-sm text-slate-500 border border-slate-100 rounded-lg">Sin versiones. Creá la primera versión para poder publicar el curso.</div>
                  ) : detailVersions.map(v => (
                    <div key={v.id} className="border border-slate-100 rounded-lg">
                      <div className="flex items-center justify-between px-4 py-2 bg-slate-50 rounded-t-lg">
                        <div className="text-sm font-medium text-slate-900">v{v.version}</div>
                        <div className="flex items-center gap-3">
                          {v.description && <span className="text-xs text-slate-500">{v.description}</span>}
                          <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${v.is_published ? 'bg-emerald-50 text-emerald-700' : 'bg-slate-100 text-slate-600'}`}>
                            {v.is_published ? 'Publicada' : 'No publicada'}
                          </span>
                          <span className="text-xs text-slate-400">{v.created_at?.slice(0, 10)}</span>
                        </div>
                      </div>
                      {(detailContents[v.id] ?? []).length === 0 ? (
                        <div className="px-4 py-3 text-sm text-slate-500">Sin contenidos</div>
                      ) : (
                        <table className="w-full text-sm">
                          <thead>
                            <tr className="border-b border-slate-100">
                              <th className="text-left px-4 py-2 font-medium text-slate-600">Contenido</th>
                              <th className="text-left px-4 py-2 font-medium text-slate-600">Tipo</th>
                              <th className="text-left px-4 py-2 font-medium text-slate-600">Duración</th>
                              <th className="text-left px-4 py-2 font-medium text-slate-600">Estado</th>
                            </tr>
                          </thead>
                          <tbody>
                            {(detailContents[v.id] ?? []).map(ct => (
                              <tr key={ct.id} className="border-b border-slate-50">
                                <td className="px-4 py-2 text-slate-900">{ct.title}</td>
                                <td className="px-4 py-2 text-slate-600 capitalize">{(ct.content_type || '').toLowerCase()}</td>
                                <td className="px-4 py-2 text-slate-600">{ct.duration_seconds ? `${Math.round(ct.duration_seconds / 60)}m` : '-'}</td>
                                <td className="px-4 py-2">
                                  <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${ct.published ? 'bg-emerald-50 text-emerald-700' : 'bg-slate-100 text-slate-600'}`}>
                                    {ct.published ? 'Publicado' : 'Borrador'}
                                  </span>
                                </td>
                              </tr>
                            ))}
                          </tbody>
                        </table>
                      )}
                    </div>
                  ))}
                </div>
              </div>
            </div>
          )}

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowDetail(false)}>Cerrar</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={showCatModal} onOpenChange={setShowCatModal}>
        <DialogContent className="max-w-xl">
          <DialogHeader>
            <DialogTitle>{editingCat ? 'Editar Categoría' : 'Nueva Categoría'}</DialogTitle>
            <DialogDescription>
              {editingCat ? 'Modificá los datos de la categoría' : 'Completá los datos para registrar una nueva categoría'}
            </DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2 col-span-2">
              <Label htmlFor="cat-name">Nombre *</Label>
              <Input id="cat-name" value={catForm.name} onChange={e => setCatForm({ ...catForm, name: e.target.value })} required />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="cat-desc">Descripción</Label>
              <Input id="cat-desc" value={catForm.description} onChange={e => setCatForm({ ...catForm, description: e.target.value })} />
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowCatModal(false)}>Cancelar</Button>
            <Button onClick={handleSaveCat} disabled={savingCat || !catForm.name}>
              {savingCat ? 'Guardando...' : editingCat ? 'Guardar Cambios' : 'Crear Categoría'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
