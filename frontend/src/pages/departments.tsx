import { useEffect, useState } from 'react'
import { Pencil, Trash2, Plus } from 'lucide-react'
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

interface Department {
  id: string
  name: string
  code?: string
  description?: string
  branch_id?: string
  active: boolean
  created_at: string
}

interface SelectOption {
  value: string
  label: string
}

const emptyForm = {
  name: '',
  code: '',
  description: '',
  branch_id: '',
}

export default function DepartmentsPage() {
  const [items, setItems] = useState<Department[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [editing, setEditing] = useState<Department | null>(null)
  const [form, setForm] = useState({ ...emptyForm })
  const [saving, setSaving] = useState(false)
  const [branches, setBranches] = useState<SelectOption[]>([])

  const fetchData = async () => {
    setLoading(true)
    try {
      const res = await api.get('/departments', { params: { limit: '100' } })
      setItems(res.data.data ?? [])
    } catch { setItems([]) }
    finally { setLoading(false) }
  }

  const fetchBranches = async () => {
    try {
      const res = await api.get('/branches', { params: { limit: '100' } })
      setBranches((res.data.data ?? []).map((b: any) => ({ value: b.id, label: b.name })))
    } catch {}
  }

  useEffect(() => { fetchData() }, [])

  const openCreate = () => {
    setEditing(null)
    setForm({ ...emptyForm })
    fetchBranches()
    setShowModal(true)
  }

  const openEdit = (d: Department) => {
    setEditing(d)
    setForm({
      name: d.name,
      code: d.code ?? '',
      description: d.description ?? '',
      branch_id: d.branch_id ?? '',
    })
    fetchBranches()
    setShowModal(true)
  }

  const handleSave = async () => {
    setSaving(true)
    try {
      const body: Record<string, any> = {}
      if (form.name) body.name = form.name
      if (form.code) body.code = form.code
      if (form.description) body.description = form.description
      if (form.branch_id) body.branch_id = form.branch_id
      if (editing) {
        body.active = editing.active
        await api.put(`/departments/${editing.id}`, body)
      } else {
        await api.post('/departments', body)
      }
      setShowModal(false)
      fetchData()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar departamento')
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (d: Department) => {
    if (!confirm(`¿Eliminar el departamento "${d.name}"?`)) return
    try {
      await api.delete(`/departments/${d.id}`)
      fetchData()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al eliminar departamento')
    }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-slate-900">Departamentos</h1>
        <Button size="sm" onClick={openCreate}><Plus size={16} className="mr-1" /> Nuevo</Button>
      </div>
      <Card>
        <CardContent className="p-0">
          {loading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
          : items.length === 0 ? <div className="p-6 text-center text-slate-500">No hay departamentos</div>
          : <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-slate-200 bg-slate-50">
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Código</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Descripción</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                    <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                  </tr>
                </thead>
                <tbody>
                  {items.map(d => (
                    <tr key={d.id} className="border-b border-slate-100 hover:bg-slate-50">
                      <td className="px-4 py-3 text-slate-500">{d.code || '-'}</td>
                      <td className="px-4 py-3 font-medium text-slate-900">{d.name}</td>
                      <td className="px-4 py-3 text-slate-600">{d.description || '-'}</td>
                      <td className="px-4 py-3">
                        <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${d.active ? 'bg-emerald-50 text-emerald-700' : 'bg-slate-100 text-slate-600'}`}>
                          {d.active ? 'Activo' : 'Inactivo'}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-right">
                        <Button variant="ghost" size="sm" onClick={() => openEdit(d)}><Pencil size={14} /></Button>
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
            <DialogTitle>{editing ? 'Editar Departamento' : 'Nuevo Departamento'}</DialogTitle>
            <DialogDescription>
              {editing ? 'Modificá los datos del departamento' : 'Completá los datos para registrar un nuevo departamento'}
            </DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="dept-name">Nombre *</Label>
              <Input id="dept-name" value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="dept-code">Código</Label>
              <Input id="dept-code" value={form.code} onChange={e => setForm({ ...form, code: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="dept-desc">Descripción</Label>
              <Input id="dept-desc" value={form.description} onChange={e => setForm({ ...form, description: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="dept-branch">Sucursal</Label>
              <Select id="dept-branch" options={branches} placeholder="Seleccionar..." value={form.branch_id} onChange={e => setForm({ ...form, branch_id: e.target.value })} />
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button>
            <Button onClick={handleSave} disabled={saving || !form.name}>
              {saving ? 'Guardando...' : editing ? 'Guardar Cambios' : 'Crear Departamento'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
