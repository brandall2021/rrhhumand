import { useEffect, useRef, useState } from 'react'
import {
  Plus,
  Pencil,
  Trash2,
  Download,
  Upload,
  FileText,
  Archive,
  ArchiveRestore,
  History,
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
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'

interface Document {
  id: string
  category_id?: string
  category_name?: string
  employee_id?: string
  department_id?: string
  title: string
  description?: string
  original_filename: string
  mime_type: string
  file_size: number
  status: string
  is_public: boolean
  expires_at?: string
  created_at: string
  current_version?: number
}

interface DocumentVersion {
  id: string
  document_id: string
  version: number
  original_filename: string
  mime_type: string
  file_size: number
  uploaded_by: string
  created_at: string
}

interface Category {
  id: string
  name: string
  description?: string
  is_active: boolean
  created_at: string
}

interface SelectOption {
  value: string
  label: string
}

const emptyDocForm = {
  title: '',
  description: '',
  category_id: '',
  employee_id: '',
  is_public: false,
}

const emptyCatForm = {
  name: '',
  description: '',
}

const formatSize = (bytes: number) => {
  if (!bytes) return '-'
  const mb = bytes / (1024 * 1024)
  return mb >= 1 ? `${mb.toFixed(1)} MB` : `${(bytes / 1024).toFixed(0)} KB`
}

export default function DocumentsPage() {
  const [items, setItems] = useState<Document[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [search, setSearch] = useState('')
  const [stats, setStats] = useState({ total: 0, expiring: 0 })

  const [showDocModal, setShowDocModal] = useState(false)
  const [editingDoc, setEditingDoc] = useState<Document | null>(null)
  const [docForm, setDocForm] = useState({ ...emptyDocForm })
  const [file, setFile] = useState<File | null>(null)
  const [uploading, setUploading] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const [categories, setCategories] = useState<SelectOption[]>([])
  const [employees, setEmployees] = useState<SelectOption[]>([])

  const [showVersions, setShowVersions] = useState(false)
  const [versionDoc, setVersionDoc] = useState<Document | null>(null)
  const [versions, setVersions] = useState<DocumentVersion[]>([])
  const [versionFile, setVersionFile] = useState<File | null>(null)
  const [uploadingVersion, setUploadingVersion] = useState(false)
  const versionInputRef = useRef<HTMLInputElement>(null)

  const [catItems, setCatItems] = useState<Category[]>([])
  const [catLoading, setCatLoading] = useState(true)
  const [showCatModal, setShowCatModal] = useState(false)
  const [editingCat, setEditingCat] = useState<Category | null>(null)
  const [catForm, setCatForm] = useState({ ...emptyCatForm })
  const [savingCat, setSavingCat] = useState(false)

  const computeExpiring = (docs: Document[]) => {
    const limit = Date.now() + 30 * 24 * 60 * 60 * 1000
    return docs.filter(d => d.expires_at && new Date(d.expires_at).getTime() <= limit).length
  }

  const fetchData = async () => {
    setLoading(true)
    try {
      const params: Record<string, string> = { limit: '100' }
      if (search) params.q = search
      const res = await api.get('/documents', { params })
      const list = res.data.data ?? []
      setItems(list)
      setError('')
      try {
        const [sRes, eRes] = await Promise.all([
          api.get('/documents/stats'),
          api.get('/documents/expiring', { params: { days: '30' } }),
        ])
        setStats({
          total: sRes.data?.data?.total_documents ?? list.length,
          expiring: Array.isArray(eRes.data?.data) ? eRes.data.data.length : computeExpiring(list),
        })
      } catch {
        setStats({ total: list.length, expiring: computeExpiring(list) })
      }
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar documentos')
      setItems([])
      setStats({ total: 0, expiring: 0 })
    } finally {
      setLoading(false)
    }
  }

  const fetchSelects = async () => {
    try {
      const [cRes, eRes] = await Promise.all([
        api.get('/document-categories'),
        api.get('/employees', { params: { limit: '200' } }),
      ])
      setCategories((cRes.data.data ?? []).map((c: any) => ({ value: c.id, label: c.name })))
      setEmployees((eRes.data.data ?? []).map((e: any) => ({ value: e.id, label: `${e.first_name} ${e.last_name}` })))
    } catch {}
  }

  const fetchCategories = async () => {
    setCatLoading(true)
    try {
      const res = await api.get('/document-categories')
      setCatItems(res.data.data ?? [])
    } catch {
      setCatItems([])
    } finally {
      setCatLoading(false)
    }
  }

  useEffect(() => { fetchData() }, [search])
  useEffect(() => { fetchSelects(); fetchCategories() }, [])

  const employeeName = (id?: string) => {
    if (!id) return '-'
    return employees.find(e => e.value === id)?.label ?? '-'
  }

  const openCreateDoc = () => {
    setEditingDoc(null)
    setDocForm({ ...emptyDocForm })
    setFile(null)
    if (fileInputRef.current) fileInputRef.current.value = ''
    fetchSelects()
    setShowDocModal(true)
  }

  const openEditDoc = (d: Document) => {
    setEditingDoc(d)
    setDocForm({
      title: d.title,
      description: d.description ?? '',
      category_id: d.category_id ?? '',
      employee_id: d.employee_id ?? '',
      is_public: d.is_public,
    })
    fetchSelects()
    setShowDocModal(true)
  }

  const handleSaveDoc = async () => {
    setUploading(true)
    try {
      if (editingDoc) {
        const body: Record<string, any> = { title: docForm.title, is_public: docForm.is_public }
        if (docForm.description) body.description = docForm.description
        if (docForm.category_id) body.category_id = docForm.category_id
        await api.put(`/documents/${editingDoc.id}`, body)
      } else {
        if (!file) return
        const fd = new FormData()
        fd.append('file', file)
        if (docForm.title) fd.append('title', docForm.title)
        if (docForm.description) fd.append('description', docForm.description)
        if (docForm.category_id) fd.append('category_id', docForm.category_id)
        if (docForm.employee_id) fd.append('employee_id', docForm.employee_id)
        fd.append('is_public', docForm.is_public ? 'true' : 'false')
        await api.post('/documents', fd)
      }
      setShowDocModal(false)
      fetchData()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar documento')
    } finally {
      setUploading(false)
    }
  }

  const handleDownload = async (d: Document) => {
    try {
      const res = await api.get(`/documents/${d.id}/download`, { responseType: 'blob' })
      const disposition = res.headers?.['content-disposition'] as string | undefined
      let filename = d.original_filename || d.title
      if (disposition) {
        const m = /filename="?([^";]+)"?/.exec(disposition)
        if (m?.[1]) filename = m[1]
      }
      const url = URL.createObjectURL(res.data)
      const a = document.createElement('a')
      a.href = url
      a.download = filename
      document.body.appendChild(a)
      a.click()
      a.remove()
      URL.revokeObjectURL(url)
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al descargar documento')
    }
  }

  const handleArchive = async (d: Document) => {
    try {
      await api.post(d.status === 'ARCHIVED' ? `/documents/${d.id}/restore` : `/documents/${d.id}/archive`)
      fetchData()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al archivar documento')
    }
  }

  const handleDelete = async (d: Document) => {
    if (!confirm(`¿Eliminar "${d.title}"?`)) return
    try {
      await api.delete(`/documents/${d.id}`)
      fetchData()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al eliminar documento')
    }
  }

  const openVersions = async (d: Document) => {
    setVersionDoc(d)
    setVersions([])
    setVersionFile(null)
    if (versionInputRef.current) versionInputRef.current.value = ''
    setShowVersions(true)
    try {
      const res = await api.get(`/documents/${d.id}/versions`)
      setVersions(res.data.data ?? [])
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al cargar versiones')
    }
  }

  const handleUploadVersion = async () => {
    if (!versionDoc || !versionFile) return
    setUploadingVersion(true)
    try {
      const fd = new FormData()
      fd.append('file', versionFile)
      await api.post(`/documents/${versionDoc.id}/versions`, fd)
      if (versionInputRef.current) versionInputRef.current.value = ''
      setVersionFile(null)
      const res = await api.get(`/documents/${versionDoc.id}/versions`)
      setVersions(res.data.data ?? [])
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al subir versión')
    } finally {
      setUploadingVersion(false)
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
        body.is_active = editingCat.is_active
        await api.put(`/document-categories/${editingCat.id}`, body)
      } else {
        await api.post('/document-categories', body)
      }
      setShowCatModal(false)
      fetchCategories()
      fetchSelects()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar categoría')
    } finally {
      setSavingCat(false)
    }
  }

  const handleDeleteCat = async (c: Category) => {
    if (!confirm(`¿Eliminar la categoría "${c.name}"?`)) return
    try {
      await api.delete(`/document-categories/${c.id}`)
      fetchCategories()
      fetchSelects()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al eliminar categoría')
    }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-slate-900">Documentos</h1>
      </div>

      <Tabs defaultValue="documents">
        <TabsList>
          <TabsTrigger value="documents">Documentos</TabsTrigger>
          <TabsTrigger value="categories">Categorías</TabsTrigger>
        </TabsList>

        <TabsContent value="documents">
          <div className="grid grid-cols-2 gap-4 mb-4">
            <Card>
              <CardContent className="p-4">
                <div className="text-sm text-slate-500">Total documentos</div>
                <div className="text-2xl font-bold text-slate-900">{stats.total}</div>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4">
                <div className="text-sm text-slate-500">Por vencer (30 días)</div>
                <div className="text-2xl font-bold text-slate-900">{stats.expiring}</div>
              </CardContent>
            </Card>
          </div>

          <div className="flex items-center gap-3 mb-4">
            <Input
              placeholder="Buscar por nombre..."
              value={search}
              onChange={e => setSearch(e.target.value)}
              className="max-w-xs"
            />
            <div className="ml-auto">
              <Button size="sm" onClick={openCreateDoc}><Plus size={16} className="mr-1" /> Nuevo</Button>
            </div>
          </div>

          {error && <div className="mb-4 p-3 bg-red-50 text-red-700 text-sm rounded-lg">{error}</div>}

          <Card>
            <CardContent className="p-0">
              {loading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
              : items.length === 0 && !error ? <div className="p-6 text-center text-slate-500">No hay documentos subidos</div>
              : <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50">
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Categoría</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Empleado</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Tamaño</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                        <th className="text-left px-4 py-3 font-medium text-slate-600">Subido</th>
                        <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                      </tr>
                    </thead>
                    <tbody>
                      {items.map(d => (
                        <tr key={d.id} className="border-b border-slate-100 hover:bg-slate-50">
                          <td className="px-4 py-3">
                            <div className="flex items-center gap-2">
                              <FileText size={16} className="text-slate-400 shrink-0" />
                              <span className="font-medium text-slate-900">{d.title || d.original_filename}</span>
                            </div>
                          </td>
                          <td className="px-4 py-3 text-slate-600">{d.category_name || '-'}</td>
                          <td className="px-4 py-3 text-slate-600">{employeeName(d.employee_id)}</td>
                          <td className="px-4 py-3 text-slate-600">{formatSize(d.file_size)}</td>
                          <td className="px-4 py-3">
                            <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${d.status === 'ARCHIVED' ? 'bg-slate-100 text-slate-600' : 'bg-emerald-50 text-emerald-700'}`}>
                              {d.status === 'ARCHIVED' ? 'Archivado' : 'Activo'}
                            </span>
                          </td>
                          <td className="px-4 py-3 text-slate-600">{d.created_at?.slice(0, 10)}</td>
                          <td className="px-4 py-3 text-right whitespace-nowrap">
                            <Button variant="ghost" size="sm" onClick={() => handleDownload(d)} title="Descargar"><Download size={14} /></Button>
                            <Button variant="ghost" size="sm" onClick={() => openVersions(d)} title="Versiones"><History size={14} /></Button>
                            <Button variant="ghost" size="sm" onClick={() => openEditDoc(d)} title="Editar"><Pencil size={14} /></Button>
                            <Button variant="ghost" size="sm" onClick={() => handleArchive(d)} title={d.status === 'ARCHIVED' ? 'Restaurar' : 'Archivar'}>
                              {d.status === 'ARCHIVED' ? <ArchiveRestore size={14} /> : <Archive size={14} />}
                            </Button>
                            <Button variant="ghost" size="sm" className="text-red-500" onClick={() => handleDelete(d)} title="Eliminar"><Trash2 size={14} /></Button>
                          </td>
                        </tr>
                      ))}
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
                            <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${c.is_active ? 'bg-emerald-50 text-emerald-700' : 'bg-slate-100 text-slate-600'}`}>
                              {c.is_active ? 'Activa' : 'Inactiva'}
                            </span>
                          </td>
                          <td className="px-4 py-3 text-right whitespace-nowrap">
                            <Button variant="ghost" size="sm" onClick={() => openEditCat(c)} title="Editar"><Pencil size={14} /></Button>
                            <Button variant="ghost" size="sm" className="text-red-500" onClick={() => handleDeleteCat(c)} title="Eliminar"><Trash2 size={14} /></Button>
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

      <Dialog open={showDocModal} onOpenChange={setShowDocModal}>
        <DialogContent className="max-w-xl">
          <DialogHeader>
            <DialogTitle>{editingDoc ? 'Editar Documento' : 'Nuevo Documento'}</DialogTitle>
            <DialogDescription>
              {editingDoc ? 'Modificá los datos del documento' : 'Completá los datos para subir un nuevo documento'}
            </DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            {!editingDoc && (
              <div className="space-y-2 col-span-2">
                <Label htmlFor="doc-file">Archivo *</Label>
                <input
                  ref={fileInputRef}
                  id="doc-file"
                  type="file"
                  onChange={e => {
                    const f = e.target.files?.[0] ?? null
                    setFile(f)
                    if (f && !docForm.title) setDocForm({ ...docForm, title: f.name })
                  }}
                  className="block w-full text-sm text-slate-600 file:mr-3 file:rounded-lg file:border-0 file:bg-brand-50 file:px-3 file:py-1.5 file:text-sm file:font-medium file:text-brand-700 hover:file:bg-brand-100"
                />
              </div>
            )}
            <div className="space-y-2 col-span-2">
              <Label htmlFor="doc-title">{editingDoc ? 'Título *' : 'Título'}</Label>
              <Input id="doc-title" value={docForm.title} onChange={e => setDocForm({ ...docForm, title: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="doc-desc">Descripción</Label>
              <Input id="doc-desc" value={docForm.description} onChange={e => setDocForm({ ...docForm, description: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="doc-cat">Categoría</Label>
              <Select id="doc-cat" options={categories} placeholder="Seleccionar..." value={docForm.category_id} onChange={e => setDocForm({ ...docForm, category_id: e.target.value })} />
            </div>
            {!editingDoc && (
              <div className="space-y-2">
                <Label htmlFor="doc-emp">Empleado</Label>
                <Select id="doc-emp" options={employees} placeholder="Seleccionar..." value={docForm.employee_id} onChange={e => setDocForm({ ...docForm, employee_id: e.target.value })} />
              </div>
            )}
            <div className="space-y-2 col-span-2">
              <label className="flex items-center gap-2 text-sm text-slate-700">
                <input type="checkbox" checked={docForm.is_public} onChange={e => setDocForm({ ...docForm, is_public: e.target.checked })} className="h-4 w-4 rounded border-slate-300 text-brand-600 focus:ring-brand-500" />
                Documento público
              </label>
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowDocModal(false)}>Cancelar</Button>
            <Button
              onClick={handleSaveDoc}
              disabled={uploading || (editingDoc ? !docForm.title : !file)}
            >
              {uploading ? 'Guardando...' : editingDoc ? 'Guardar Cambios' : 'Subir'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={showVersions} onOpenChange={setShowVersions}>
        <DialogContent className="max-w-xl">
          <DialogHeader>
            <DialogTitle>Versiones de {versionDoc?.title}</DialogTitle>
            <DialogDescription>Historial de versiones del documento</DialogDescription>
          </DialogHeader>

          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="version-file">Nueva versión *</Label>
              <div className="flex gap-3">
                <input
                  ref={versionInputRef}
                  id="version-file"
                  type="file"
                  onChange={e => setVersionFile(e.target.files?.[0] ?? null)}
                  className="block flex-1 text-sm text-slate-600 file:mr-3 file:rounded-lg file:border-0 file:bg-brand-50 file:px-3 file:py-1.5 file:text-sm file:font-medium file:text-brand-700 hover:file:bg-brand-100"
                />
                <Button size="sm" onClick={handleUploadVersion} disabled={uploadingVersion || !versionFile}>
                  <Upload size={14} className="mr-1" /> {uploadingVersion ? 'Subiendo...' : 'Subir'}
                </Button>
              </div>
            </div>

            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-slate-200 bg-slate-50">
                    <th className="text-left px-3 py-2 font-medium text-slate-600">Versión</th>
                    <th className="text-left px-3 py-2 font-medium text-slate-600">Archivo</th>
                    <th className="text-left px-3 py-2 font-medium text-slate-600">Tamaño</th>
                    <th className="text-left px-3 py-2 font-medium text-slate-600">Fecha</th>
                  </tr>
                </thead>
                <tbody>
                  {versions.length === 0 ? (
                    <tr><td colSpan={4} className="px-3 py-4 text-center text-slate-500">Sin versiones</td></tr>
                  ) : versions.map(v => (
                    <tr key={v.id} className="border-b border-slate-100">
                      <td className="px-3 py-2 text-slate-600">v{v.version}</td>
                      <td className="px-3 py-2 font-medium text-slate-900">{v.original_filename}</td>
                      <td className="px-3 py-2 text-slate-600">{formatSize(v.file_size)}</td>
                      <td className="px-3 py-2 text-slate-600">{v.created_at?.slice(0, 10)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowVersions(false)}>Cerrar</Button>
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
