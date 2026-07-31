import { useEffect, useState, useRef } from 'react'
import { Upload, Download, Trash2, FileText, Archive, ArchiveRestore } from 'lucide-react'
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

interface Document {
  id: string
  title: string
  original_filename?: string
  category_name?: string
  category_id?: string
  file_size: number
  mime_type?: string
  created_at: string
  employee_id?: string
  status: string
}

interface SelectOption {
  value: string
  label: string
}

export default function DocumentsPage() {
  const [items, setItems] = useState<Document[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showModal, setShowModal] = useState(false)
  const [uploading, setUploading] = useState(false)
  const [file, setFile] = useState<File | null>(null)
  const [form, setForm] = useState({ title: '', description: '', category_id: '', employee_id: '', is_public: false })
  const [categories, setCategories] = useState<SelectOption[]>([])
  const [employees, setEmployees] = useState<SelectOption[]>([])
  const fileInputRef = useRef<HTMLInputElement>(null)

  const formatSize = (bytes: number) => {
    if (!bytes) return '-'
    const mb = bytes / (1024 * 1024)
    return mb >= 1 ? `${mb.toFixed(1)} MB` : `${(bytes / 1024).toFixed(0)} KB`
  }

  const fetchData = async () => {
    setLoading(true)
    try {
      const res = await api.get('/documents', { params: { limit: '100' } })
      setItems(res.data.data ?? [])
      setError('')
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar documentos')
      setItems([])
    } finally {
      setLoading(false)
    }
  }

  const fetchSelects = async () => {
    try {
      const [cRes, eRes] = await Promise.all([
        api.get('/document-categories', { params: { limit: '100' } }),
        api.get('/employees', { params: { limit: '200' } }),
      ])
      setCategories((cRes.data.data ?? []).map((c: any) => ({ value: c.id, label: c.name })))
      setEmployees((eRes.data.data ?? []).map((e: any) => ({ value: e.id, label: `${e.first_name} ${e.last_name}` })))
    } catch {}
  }

  useEffect(() => { fetchData() }, [])

  const openUpload = () => {
    setForm({ title: '', description: '', category_id: '', employee_id: '', is_public: false })
    setFile(null)
    if (fileInputRef.current) fileInputRef.current.value = ''
    fetchSelects()
    setShowModal(true)
  }

  const handleUpload = async () => {
    if (!file) return
    setUploading(true)
    try {
      const fd = new FormData()
      fd.append('file', file)
      if (form.title) fd.append('title', form.title)
      if (form.description) fd.append('description', form.description)
      if (form.category_id) fd.append('category_id', form.category_id)
      if (form.employee_id) fd.append('employee_id', form.employee_id)
      fd.append('is_public', form.is_public ? 'true' : 'false')
      await api.post('/documents', fd)
      setShowModal(false)
      fetchData()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al subir documento')
    } finally {
      setUploading(false)
    }
  }

  const handleDownload = async (d: Document) => {
    try {
      const res = await api.get(`/documents/${d.id}/download`, { responseType: 'blob' })
      const url = URL.createObjectURL(res.data)
      const a = document.createElement('a')
      a.href = url
      a.download = d.original_filename || d.title
      document.body.appendChild(a)
      a.click()
      a.remove()
      URL.revokeObjectURL(url)
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al descargar documento')
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

  const handleArchive = async (d: Document) => {
    try {
      await api.post(d.status === 'ARCHIVED' ? `/documents/${d.id}/restore` : `/documents/${d.id}/archive`)
      fetchData()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al archivar documento')
    }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-slate-900">Documentos</h1>
        <Button size="sm" onClick={openUpload}><Upload size={16} className="mr-1" /> Subir</Button>
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
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Tamaño</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Subido</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                    <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                  </tr>
                </thead>
                <tbody>
                  {items.map(d => (
                    <tr key={d.id} className="border-b border-slate-100 hover:bg-slate-50">
                      <td className="px-4 py-3">
                        <div className="flex items-center gap-2">
                          <FileText size={16} className="text-slate-400" />
                          <span className="font-medium text-slate-900">{d.title || d.original_filename}</span>
                        </div>
                      </td>
                      <td className="px-4 py-3 text-slate-600">{d.category_name || '-'}</td>
                      <td className="px-4 py-3 text-slate-600">{formatSize(d.file_size)}</td>
                      <td className="px-4 py-3 text-slate-600">{d.created_at?.slice(0, 10)}</td>
                      <td className="px-4 py-3">
                        <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${d.status === 'ARCHIVED' ? 'bg-slate-100 text-slate-600' : 'bg-emerald-50 text-emerald-700'}`}>
                          {d.status === 'ARCHIVED' ? 'Archivado' : 'Activo'}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-right whitespace-nowrap">
                        <Button variant="ghost" size="sm" onClick={() => handleDownload(d)} title="Descargar"><Download size={14} /></Button>
                        <Button variant="ghost" size="sm" onClick={() => handleArchive(d)} title={d.status === 'ARCHIVED' ? 'Restaurar' : 'Archivar'}>
                          {d.status === 'ARCHIVED' ? <ArchiveRestore size={14} /> : <Archive size={14} />}
                        </Button>
                        <Button variant="ghost" size="sm" className="text-red-500" onClick={() => handleDelete(d)}><Trash2 size={14} /></Button>
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
            <DialogTitle>Subir Documento</DialogTitle>
            <DialogDescription>Seleccioná un archivo y los datos del documento</DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2 col-span-2">
              <Label htmlFor="doc-file">Archivo *</Label>
              <input
                ref={fileInputRef}
                id="doc-file"
                type="file"
                onChange={e => {
                  const f = e.target.files?.[0] ?? null
                  setFile(f)
                  if (f && !form.title) setForm({ ...form, title: f.name })
                }}
                className="block w-full text-sm text-slate-600 file:mr-3 file:rounded-lg file:border-0 file:bg-brand-50 file:px-3 file:py-1.5 file:text-sm file:font-medium file:text-brand-700 hover:file:bg-brand-100"
              />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="doc-title">Título</Label>
              <Input id="doc-title" value={form.title} onChange={e => setForm({ ...form, title: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="doc-desc">Descripción</Label>
              <Input id="doc-desc" value={form.description} onChange={e => setForm({ ...form, description: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="doc-cat">Categoría</Label>
              <Select id="doc-cat" options={categories} placeholder="Seleccionar..." value={form.category_id} onChange={e => setForm({ ...form, category_id: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="doc-emp">Empleado</Label>
              <Select id="doc-emp" options={employees} placeholder="Seleccionar..." value={form.employee_id} onChange={e => setForm({ ...form, employee_id: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <label className="flex items-center gap-2 text-sm text-slate-700">
                <input type="checkbox" checked={form.is_public} onChange={e => setForm({ ...form, is_public: e.target.checked })} className="h-4 w-4 rounded border-slate-300 text-brand-600 focus:ring-brand-500" />
                Documento público
              </label>
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button>
            <Button onClick={handleUpload} disabled={uploading || !file}>
              {uploading ? 'Subiendo...' : 'Subir'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
