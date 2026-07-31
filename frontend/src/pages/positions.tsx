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

interface Position {
  id: string
  name: string
  code?: string
  description?: string
  department_id?: string
  department_name?: string
  level?: number
  min_salary?: number
  max_salary?: number
  active: boolean
}

interface SelectOption {
  value: string
  label: string
}

const emptyForm = {
  name: '',
  code: '',
  description: '',
  department_id: '',
  level: '1',
  min_salary: '',
  max_salary: '',
}

export default function PositionsPage() {
  const [items, setItems] = useState<Position[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [editing, setEditing] = useState<Position | null>(null)
  const [form, setForm] = useState({ ...emptyForm })
  const [saving, setSaving] = useState(false)
  const [departments, setDepartments] = useState<SelectOption[]>([])

  const fetchData = async () => {
    setLoading(true)
    try {
      const res = await api.get('/positions', { params: { limit: '100' } })
      setItems(res.data.data ?? [])
    } catch { setItems([]) }
    finally { setLoading(false) }
  }

  const fetchDepartments = async () => {
    try {
      const res = await api.get('/departments', { params: { limit: '100' } })
      setDepartments((res.data.data ?? []).map((d: any) => ({ value: d.id, label: d.name })))
    } catch {}
  }

  useEffect(() => { fetchData() }, [])

  const openCreate = () => {
    setEditing(null)
    setForm({ ...emptyForm })
    fetchDepartments()
    setShowModal(true)
  }

  const openEdit = (p: Position) => {
    setEditing(p)
    setForm({
      name: p.name,
      code: p.code ?? '',
      description: p.description ?? '',
      department_id: p.department_id ?? '',
      level: String(p.level ?? 1),
      min_salary: p.min_salary != null ? String(p.min_salary) : '',
      max_salary: p.max_salary != null ? String(p.max_salary) : '',
    })
    fetchDepartments()
    setShowModal(true)
  }

  const handleSave = async () => {
    setSaving(true)
    try {
      const body: Record<string, any> = {
        level: parseInt(form.level || '1', 10),
      }
      if (form.name) body.name = form.name
      if (form.code) body.code = form.code
      if (form.description) body.description = form.description
      if (form.department_id) body.department_id = form.department_id
      if (form.min_salary) body.min_salary = parseFloat(form.min_salary)
      if (form.max_salary) body.max_salary = parseFloat(form.max_salary)
      if (editing) {
        body.active = editing.active
        await api.put(`/positions/${editing.id}`, body)
      } else {
        await api.post('/positions', body)
      }
      setShowModal(false)
      fetchData()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar posición')
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (p: Position) => {
    if (!confirm(`¿Eliminar la posición "${p.name}"?`)) return
    try {
      await api.delete(`/positions/${p.id}`)
      fetchData()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al eliminar posición')
    }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-slate-900">Posiciones</h1>
        <Button size="sm" onClick={openCreate}><Plus size={16} className="mr-1" /> Nueva</Button>
      </div>
      <Card>
        <CardContent className="p-0">
          {loading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
          : items.length === 0 ? <div className="p-6 text-center text-slate-500">No hay posiciones registradas</div>
          : <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-slate-200 bg-slate-50">
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Código</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Departamento</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Nivel</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Rango Salarial</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                    <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                  </tr>
                </thead>
                <tbody>
                  {items.map(p => (
                    <tr key={p.id} className="border-b border-slate-100 hover:bg-slate-50">
                      <td className="px-4 py-3 text-slate-500">{p.code || '-'}</td>
                      <td className="px-4 py-3 font-medium text-slate-900">{p.name}</td>
                      <td className="px-4 py-3 text-slate-600">{p.department_name || '-'}</td>
                      <td className="px-4 py-3 text-slate-600">{p.level ?? 1}</td>
                      <td className="px-4 py-3 text-slate-600">
                        {p.min_salary != null && p.max_salary != null ? `$${p.min_salary.toLocaleString()} - $${p.max_salary.toLocaleString()}` : '-'}
                      </td>
                      <td className="px-4 py-3">
                        <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${p.active ? 'bg-emerald-50 text-emerald-700' : 'bg-slate-100 text-slate-600'}`}>
                          {p.active ? 'Activo' : 'Inactivo'}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-right">
                        <Button variant="ghost" size="sm" onClick={() => openEdit(p)}><Pencil size={14} /></Button>
                        <Button variant="ghost" size="sm" className="text-red-500" onClick={() => handleDelete(p)}><Trash2 size={14} /></Button>
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
            <DialogTitle>{editing ? 'Editar Posición' : 'Nueva Posición'}</DialogTitle>
            <DialogDescription>
              {editing ? 'Modificá los datos de la posición' : 'Completá los datos para registrar una nueva posición'}
            </DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="pos-name">Nombre *</Label>
              <Input id="pos-name" value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="pos-code">Código</Label>
              <Input id="pos-code" value={form.code} onChange={e => setForm({ ...form, code: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="pos-desc">Descripción</Label>
              <Input id="pos-desc" value={form.description} onChange={e => setForm({ ...form, description: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="pos-dept">Departamento</Label>
              <Select id="pos-dept" options={departments} placeholder="Seleccionar..." value={form.department_id} onChange={e => setForm({ ...form, department_id: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="pos-level">Nivel</Label>
              <Input id="pos-level" type="number" min={1} value={form.level} onChange={e => setForm({ ...form, level: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="pos-min">Salario Mínimo</Label>
              <Input id="pos-min" type="number" value={form.min_salary} onChange={e => setForm({ ...form, min_salary: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="pos-max">Salario Máximo</Label>
              <Input id="pos-max" type="number" value={form.max_salary} onChange={e => setForm({ ...form, max_salary: e.target.value })} />
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button>
            <Button onClick={handleSave} disabled={saving || !form.name}>
              {saving ? 'Guardando...' : editing ? 'Guardar Cambios' : 'Crear Posición'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
