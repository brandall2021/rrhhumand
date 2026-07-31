import { useEffect, useState } from 'react'
import { Plus, Pencil, Trash2, Send } from 'lucide-react'
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

interface Expense {
  id: string
  description: string
  original_amount?: string | number
  original_currency?: string
  category_id?: string
  expense_date?: string
  merchant_name?: string
  status: string
  created_at: string
  employee_id?: string
}

interface SelectOption {
  value: string
  label: string
}

const emptyForm = {
  category_id: '',
  expense_date: '',
  description: '',
  original_amount: '',
  original_currency: 'ARS',
  merchant_name: '',
  payment_method_id: '',
  observation: '',
}

const statusStyles: Record<string, { label: string; cls: string }> = {
  draft: { label: 'Borrador', cls: 'bg-slate-100 text-slate-600' },
  submitted: { label: 'Enviado', cls: 'bg-blue-50 text-blue-700' },
  pending_approval: { label: 'Pendiente', cls: 'bg-amber-50 text-amber-700' },
  approved: { label: 'Aprobado', cls: 'bg-emerald-50 text-emerald-700' },
  rejected: { label: 'Rechazado', cls: 'bg-red-50 text-red-700' },
  observed: { label: 'Observado', cls: 'bg-violet-50 text-violet-700' },
  paid: { label: 'Pagado', cls: 'bg-teal-50 text-teal-700' },
  canceled: { label: 'Cancelado', cls: 'bg-slate-100 text-slate-500' },
}

const statusBadge = (status: string) => {
  const s = statusStyles[status] ?? { label: status, cls: 'bg-slate-100 text-slate-600' }
  return <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${s.cls}`}>{s.label}</span>
}

export default function ExpensesPage() {
  const [items, setItems] = useState<Expense[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showModal, setShowModal] = useState(false)
  const [editing, setEditing] = useState<Expense | null>(null)
  const [form, setForm] = useState({ ...emptyForm })
  const [saving, setSaving] = useState(false)
  const [categories, setCategories] = useState<SelectOption[]>([])
  const [paymentMethods, setPaymentMethods] = useState<SelectOption[]>([])

  const fetchData = async () => {
    setLoading(true)
    try {
      const res = await api.get('/expenses', { params: { limit: '100' } })
      const data = res.data.data ?? res.data ?? []
      setItems(Array.isArray(data) ? data : [])
      setError('')
    } catch (e: any) {
      setError(e.response?.data?.error || 'Error al cargar gastos')
      setItems([])
    } finally {
      setLoading(false)
    }
  }

  const fetchSelects = async () => {
    try {
      const [cRes, pRes] = await Promise.all([
        api.get('/expenses/categories'),
        api.get('/expenses/payment-methods'),
      ])
      setCategories((cRes.data.data ?? cRes.data ?? []).map((c: any) => ({ value: c.id, label: `${c.code ? c.code + ' - ' : ''}${c.name}` })))
      setPaymentMethods((pRes.data.data ?? pRes.data ?? []).map((p: any) => ({ value: p.id, label: p.name })))
    } catch {}
  }

  useEffect(() => { fetchData() }, [])

  const openCreate = () => {
    setEditing(null)
    setForm({ ...emptyForm })
    fetchSelects()
    setShowModal(true)
  }

  const openEdit = (e: Expense) => {
    setEditing(e)
    setForm({
      category_id: e.category_id ?? '',
      expense_date: e.expense_date ? e.expense_date.slice(0, 10) : '',
      description: e.description,
      original_amount: e.original_amount != null ? String(e.original_amount) : '',
      original_currency: e.original_currency ?? 'ARS',
      merchant_name: e.merchant_name ?? '',
      payment_method_id: '',
      observation: '',
    })
    fetchSelects()
    setShowModal(true)
  }

  const handleSave = async () => {
    setSaving(true)
    try {
      const body: Record<string, any> = {
        category_id: form.category_id,
        expense_date: form.expense_date,
        description: form.description,
        original_amount: form.original_amount,
        original_currency: form.original_currency,
      }
      if (form.merchant_name) body.merchant_name = form.merchant_name
      if (form.payment_method_id) body.payment_method_id = form.payment_method_id
      if (form.observation) body.observation = form.observation
      if (editing) {
        await api.put(`/expenses/${editing.id}`, body)
      } else {
        await api.post('/expenses', body)
      }
      setShowModal(false)
      fetchData()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar gasto')
    } finally {
      setSaving(false)
    }
  }

  const handleSubmit = async (e: Expense) => {
    try {
      await api.post(`/expenses/${e.id}/submit`)
      fetchData()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al enviar gasto')
    }
  }

  const handleDelete = async (e: Expense) => {
    if (!confirm('¿Eliminar este gasto?')) return
    try {
      await api.delete(`/expenses/${e.id}`)
      fetchData()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al eliminar gasto')
    }
  }

  const categoryName = (id?: string) => {
    if (!id) return '-'
    return categories.find(c => c.value === id)?.label ?? id.slice(0, 8)
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-slate-900">Gastos</h1>
        <Button size="sm" onClick={openCreate}><Plus size={16} className="mr-1" /> Nuevo gasto</Button>
      </div>

      {error && <div className="mb-4 p-3 bg-red-50 text-red-700 text-sm rounded-lg">{error}</div>}

      <Card>
        <CardContent className="p-0">
          {loading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
          : items.length === 0 && !error ? <div className="p-6 text-center text-slate-500">No hay gastos registrados</div>
          : <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-slate-200 bg-slate-50">
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Descripción</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Categoría</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Fecha</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Monto</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                    <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                  </tr>
                </thead>
                <tbody>
                  {items.map(e => (
                    <tr key={e.id} className="border-b border-slate-100 hover:bg-slate-50">
                      <td className="px-4 py-3 font-medium text-slate-900">{e.description}</td>
                      <td className="px-4 py-3 text-slate-600">{categoryName(e.category_id)}</td>
                      <td className="px-4 py-3 text-slate-600">{e.expense_date ? e.expense_date.slice(0, 10) : '-'}</td>
                      <td className="px-4 py-3 text-slate-900 font-medium">
                        {(e.original_amount != null && e.original_amount !== '' ? parseFloat(String(e.original_amount)) : 0).toLocaleString()} {e.original_currency || ''}
                      </td>
                      <td className="px-4 py-3">{statusBadge(e.status)}</td>
                      <td className="px-4 py-3 text-right whitespace-nowrap">
                        {e.status === 'draft' && (
                          <Button variant="ghost" size="sm" className="text-blue-600" onClick={() => handleSubmit(e)} title="Enviar a aprobación">
                            <Send size={14} />
                          </Button>
                        )}
                        <Button variant="ghost" size="sm" onClick={() => openEdit(e)}><Pencil size={14} /></Button>
                        <Button variant="ghost" size="sm" className="text-red-500" onClick={() => handleDelete(e)}><Trash2 size={14} /></Button>
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
            <DialogTitle>{editing ? 'Editar Gasto' : 'Nuevo Gasto'}</DialogTitle>
            <DialogDescription>
              {editing ? 'Modificá los datos del gasto' : 'Completá los datos para registrar un nuevo gasto'}
            </DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2 col-span-2">
              <Label htmlFor="exp-desc">Descripción *</Label>
              <Input id="exp-desc" value={form.description} onChange={e => setForm({ ...form, description: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="exp-cat">Categoría *</Label>
              <Select id="exp-cat" options={categories} placeholder="Seleccionar..." value={form.category_id} onChange={e => setForm({ ...form, category_id: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="exp-date">Fecha *</Label>
              <Input id="exp-date" type="date" value={form.expense_date} onChange={e => setForm({ ...form, expense_date: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="exp-amount">Monto *</Label>
              <Input id="exp-amount" type="number" step="0.01" value={form.original_amount} onChange={e => setForm({ ...form, original_amount: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="exp-currency">Moneda *</Label>
              <Select
                id="exp-currency"
                options={[{ value: 'ARS', label: 'ARS' }, { value: 'USD', label: 'USD' }, { value: 'EUR', label: 'EUR' }]}
                value={form.original_currency}
                onChange={e => setForm({ ...form, original_currency: e.target.value })}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="exp-merchant">Comercio</Label>
              <Input id="exp-merchant" value={form.merchant_name} onChange={e => setForm({ ...form, merchant_name: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="exp-payment">Método de Pago</Label>
              <Select id="exp-payment" options={paymentMethods} placeholder="Seleccionar..." value={form.payment_method_id} onChange={e => setForm({ ...form, payment_method_id: e.target.value })} />
            </div>
            <div className="space-y-2 col-span-2">
              <Label htmlFor="exp-obs">Observación</Label>
              <Input id="exp-obs" value={form.observation} onChange={e => setForm({ ...form, observation: e.target.value })} />
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button>
            <Button onClick={handleSave} disabled={saving || !form.description || !form.category_id || !form.expense_date || !form.original_amount}>
              {saving ? 'Guardando...' : editing ? 'Guardar Cambios' : 'Crear Gasto'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
