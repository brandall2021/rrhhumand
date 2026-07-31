import { useEffect, useState } from 'react'
import { Plus, Pencil } from 'lucide-react'
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

interface Structure {
  id: string
  name: string
  description?: string
  currency: string
  effective_from: string
  effective_to?: string
  status: string
}

const emptyForm = {
  name: '',
  description: '',
  currency: 'USD',
  effective_from: '',
  effective_to: '',
}

export default function CompensationPage() {
  const [items, setItems] = useState<Structure[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showModal, setShowModal] = useState(false)
  const [editing, setEditing] = useState<Structure | null>(null)
  const [form, setForm] = useState({ ...emptyForm })
  const [saving, setSaving] = useState(false)

  const fetchData = async () => {
    setLoading(true)
    try {
      const res = await api.get('/compensation/structures')
      setItems(Array.isArray(res.data) ? res.data : res.data.data ?? [])
      setError('')
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar estructuras')
      setItems([])
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchData() }, [])

  const openCreate = () => {
    setEditing(null)
    setForm({ ...emptyForm })
    setShowModal(true)
  }

  const openEdit = (s: Structure) => {
    setEditing(s)
    setForm({
      name: s.name,
      description: s.description ?? '',
      currency: s.currency || 'USD',
      effective_from: s.effective_from ? s.effective_from.slice(0, 10) : '',
      effective_to: s.effective_to ? s.effective_to.slice(0, 10) : '',
    })
    setShowModal(true)
  }

  const handleSave = async () => {
    setSaving(true)
    try {
      const body: Record<string, any> = {}
      if (form.name) body.name = form.name
      if (form.description) body.description = form.description
      if (form.currency) body.currency = form.currency
      if (form.effective_from) body.effective_from = form.effective_from
      if (form.effective_to) body.effective_to = form.effective_to
      if (editing) {
        body.status = editing.status
        await api.put(`/compensation/structures/${editing.id}`, body)
      } else {
        await api.post('/compensation/structures', body)
      }
      setShowModal(false)
      fetchData()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar estructura')
    } finally {
      setSaving(false)
    }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-slate-900">Compensaciones</h1>
        <Button size="sm" onClick={openCreate}><Plus size={16} className="mr-1" /> Nueva</Button>
      </div>

      {error && <div className="mb-4 p-3 bg-red-50 text-red-700 text-sm rounded-lg">{error}</div>}

      <Card>
        <CardContent className="p-0">
          {loading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
          : items.length === 0 && !error ? <div className="p-6 text-center text-slate-500">No hay estructuras salariales</div>
          : <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-slate-200 bg-slate-50">
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Descripción</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Moneda</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Vigencia</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                    <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                  </tr>
                </thead>
                <tbody>
                  {items.map(s => (
                    <tr key={s.id} className="border-b border-slate-100 hover:bg-slate-50">
                      <td className="px-4 py-3 font-medium text-slate-900">{s.name}</td>
                      <td className="px-4 py-3 text-slate-600">{s.description || '-'}</td>
                      <td className="px-4 py-3 text-slate-600">{s.currency || '-'}</td>
                      <td className="px-4 py-3 text-slate-600">
                        {s.effective_from ? s.effective_from.slice(0, 10) : '-'}
                        {s.effective_to ? ` → ${s.effective_to.slice(0, 10)}` : ''}
                      </td>
                      <td className="px-4 py-3">
                        <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${s.status === 'ACTIVE' ? 'bg-emerald-50 text-emerald-700' : 'bg-slate-100 text-slate-600'}`}>
                          {s.status === 'ACTIVE' ? 'Activo' : s.status === 'DRAFT' ? 'Borrador' : s.status}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-right">
                        <Button variant="ghost" size="sm" onClick={() => openEdit(s)}><Pencil size={14} /></Button>
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
            <DialogTitle>{editing ? 'Editar Estructura Salarial' : 'Nueva Estructura Salarial'}</DialogTitle>
            <DialogDescription>
              {editing ? 'Modificá los datos de la estructura' : 'Completá los datos para crear una nueva estructura'}
            </DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2 col-span-2">
              <Label htmlFor="comp-name">Nombre *</Label>
              <Input id="comp-name" value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} required />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="comp-desc">Descripción</Label>
              <textarea
                id="comp-desc"
                value={form.description}
                onChange={e => setForm({ ...form, description: e.target.value })}
                rows={3}
                className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-brand-500"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="comp-currency">Moneda</Label>
              <Select
                id="comp-currency"
                options={[{ value: 'USD', label: 'USD' }, { value: 'ARS', label: 'ARS' }, { value: 'EUR', label: 'EUR' }]}
                value={form.currency}
                onChange={e => setForm({ ...form, currency: e.target.value })}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="comp-from">Vigencia desde *</Label>
              <Input id="comp-from" type="date" value={form.effective_from} onChange={e => setForm({ ...form, effective_from: e.target.value })} required />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="comp-to">Vigencia hasta</Label>
              <Input id="comp-to" type="date" value={form.effective_to} onChange={e => setForm({ ...form, effective_to: e.target.value })} />
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button>
            <Button onClick={handleSave} disabled={saving || !form.name || !form.effective_from}>
              {saving ? 'Guardando...' : editing ? 'Guardar Cambios' : 'Crear Estructura'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
