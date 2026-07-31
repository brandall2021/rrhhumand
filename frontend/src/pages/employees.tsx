import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Pencil, Trash2, Plus, Eye } from 'lucide-react'
import api from '@/lib/api'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
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

interface Employee {
  id: string
  employee_number: string
  first_name: string
  last_name: string
  dni?: string
  email?: string
  phone?: string
  birth_date?: string
  photo_url?: string
  branch_id?: string
  department_id?: string
  position_id?: string
  manager_id?: string
  hire_date: string
  termination_date?: string
  status: string
  branch_name?: string
  department_name?: string
  position_name?: string
  manager_name?: string
}

interface SelectOption {
  value: string
  label: string
}

const emptyForm = {
  employee_number: '',
  first_name: '',
  last_name: '',
  dni: '',
  email: '',
  phone: '',
  birth_date: '',
  branch_id: '',
  department_id: '',
  position_id: '',
  manager_id: '',
  hire_date: '',
}

export default function EmployeesPage() {
  const navigate = useNavigate()
  const [items, setItems] = useState<Employee[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [showModal, setShowModal] = useState(false)
  const [editing, setEditing] = useState<Employee | null>(null)
  const [form, setForm] = useState({ ...emptyForm })
  const [saving, setSaving] = useState(false)

  const [branches, setBranches] = useState<SelectOption[]>([])
  const [departments, setDepartments] = useState<SelectOption[]>([])
  const [positions, setPositions] = useState<SelectOption[]>([])
  const [managers, setManagers] = useState<SelectOption[]>([])

  const fetchData = async () => {
    setLoading(true)
    try {
      const params: Record<string, string> = { limit: '100' }
      if (search) params.search = search
      const res = await api.get('/employees', { params })
      setItems(res.data.data ?? [])
    } catch { setItems([]) }
    finally { setLoading(false) }
  }

  const fetchSelects = async () => {
    try {
      const [bRes, dRes, pRes, mRes] = await Promise.all([
        api.get('/branches', { params: { limit: '100' } }),
        api.get('/departments', { params: { limit: '100' } }),
        api.get('/positions', { params: { limit: '100' } }),
        api.get('/employees', { params: { limit: '200' } }),
      ])
      setBranches((bRes.data.data ?? []).map((b: any) => ({ value: b.id, label: b.name })))
      setDepartments((dRes.data.data ?? []).map((d: any) => ({ value: d.id, label: d.name })))
      setPositions((pRes.data.data ?? []).map((p: any) => ({ value: p.id, label: `${p.name} (${p.code})` })))
      setManagers((mRes.data.data ?? []).map((m: any) => ({ value: m.id, label: `${m.first_name} ${m.last_name}` })))
    } catch {}
  }

  useEffect(() => { fetchData() }, [search])

  const openCreate = () => {
    setEditing(null)
    setForm({ ...emptyForm })
    fetchSelects()
    setShowModal(true)
  }

  const openEdit = (emp: Employee) => {
    setEditing(emp)
    setForm({
      employee_number: emp.employee_number,
      first_name: emp.first_name,
      last_name: emp.last_name,
      dni: emp.dni ?? '',
      email: emp.email ?? '',
      phone: emp.phone ?? '',
      birth_date: emp.birth_date ?? '',
      branch_id: emp.branch_id ?? '',
      department_id: emp.department_id ?? '',
      position_id: emp.position_id ?? '',
      manager_id: emp.manager_id ?? '',
      hire_date: emp.hire_date ?? '',
    })
    fetchSelects()
    setShowModal(true)
  }

  const handleSave = async () => {
    setSaving(true)
    try {
      const body: Record<string, any> = {
        employee_number: form.employee_number,
        first_name: form.first_name,
        last_name: form.last_name,
        hire_date: form.hire_date,
      }
      for (const k of ['dni','email','phone','birth_date','branch_id','department_id','position_id','manager_id']) {
        if (form[k as keyof typeof form]) body[k] = form[k as keyof typeof form]
      }
      if (editing) {
        await api.put(`/employees/${editing.id}`, body)
      } else {
        await api.post('/employees', body)
      }
      setShowModal(false)
      fetchData()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar empleado')
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (emp: Employee) => {
    if (!confirm(`¿Eliminar a ${emp.first_name} ${emp.last_name}?`)) return
    try {
      await api.delete(`/employees/${emp.id}`)
      fetchData()
    } catch { alert('Error al eliminar empleado') }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-slate-900">Empleados</h1>
        <Button size="sm" onClick={openCreate}><Plus size={16} className="mr-1" /> Nuevo</Button>
      </div>

      <Card>
        <CardHeader className="pb-3">
          <div className="flex items-center gap-3">
            <Input placeholder="Buscar por nombre o email..." value={search} onChange={e => setSearch(e.target.value)} className="max-w-xs" />
          </div>
        </CardHeader>
        <CardContent className="p-0">
          {loading ? (
            <div className="p-6 text-center text-slate-500">Cargando...</div>
          ) : items.length === 0 ? (
            <div className="p-6 text-center text-slate-500">No hay empleados registrados</div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-slate-200 bg-slate-50">
                    <th className="text-left px-4 py-3 font-medium text-slate-600">N°</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Email</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Departamento</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Posición</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                    <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                  </tr>
                </thead>
                <tbody>
                  {items.map((emp) => (
                    <tr key={emp.id} className="border-b border-slate-100 hover:bg-slate-50">
                      <td className="px-4 py-3 text-slate-500">{emp.employee_number}</td>
                      <td className="px-4 py-3 font-medium text-slate-900">{emp.first_name} {emp.last_name}</td>
                      <td className="px-4 py-3 text-slate-600">{emp.email}</td>
                      <td className="px-4 py-3 text-slate-600">{emp.department_name || '-'}</td>
                      <td className="px-4 py-3 text-slate-600">{emp.position_name || '-'}</td>
                      <td className="px-4 py-3">
                        <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${emp.status === 'active' ? 'bg-emerald-50 text-emerald-700' : 'bg-slate-100 text-slate-600'}`}>
                          {emp.status === 'active' ? 'Activo' : 'Inactivo'}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-right whitespace-nowrap">
                        <Button variant="ghost" size="sm" onClick={() => navigate(`/employees/${emp.id}`)}>
                          <Eye size={14} />
                        </Button>
                        <Button variant="ghost" size="sm" onClick={() => openEdit(emp)}>
                          <Pencil size={14} />
                        </Button>
                        <Button variant="ghost" size="sm" className="text-red-500" onClick={() => handleDelete(emp)}>
                          <Trash2 size={14} />
                        </Button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </CardContent>
      </Card>

      <Dialog open={showModal} onOpenChange={setShowModal}>
        <DialogContent className="max-w-xl">
          <DialogHeader>
            <DialogTitle>{editing ? 'Editar Empleado' : 'Nuevo Empleado'}</DialogTitle>
            <DialogDescription>
              {editing ? 'Modificá los datos del empleado' : 'Completá los datos para registrar un nuevo empleado'}
            </DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="emp-number">N° Empleado *</Label>
              <Input id="emp-number" value={form.employee_number} onChange={e => setForm({ ...form, employee_number: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="hire-date">Fecha Ingreso *</Label>
              <Input id="hire-date" type="date" value={form.hire_date} onChange={e => setForm({ ...form, hire_date: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="first-name">Nombre *</Label>
              <Input id="first-name" value={form.first_name} onChange={e => setForm({ ...form, first_name: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="last-name">Apellido *</Label>
              <Input id="last-name" value={form.last_name} onChange={e => setForm({ ...form, last_name: e.target.value })} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="dni">DNI</Label>
              <Input id="dni" value={form.dni} onChange={e => setForm({ ...form, dni: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input id="email" type="email" value={form.email} onChange={e => setForm({ ...form, email: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="phone">Teléfono</Label>
              <Input id="phone" value={form.phone} onChange={e => setForm({ ...form, phone: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="birth-date">Fecha Nacimiento</Label>
              <Input id="birth-date" type="date" value={form.birth_date} onChange={e => setForm({ ...form, birth_date: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="branch">Sucursal</Label>
              <Select id="branch" options={branches} placeholder="Seleccionar..." value={form.branch_id} onChange={e => setForm({ ...form, branch_id: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="department">Departamento</Label>
              <Select id="department" options={departments} placeholder="Seleccionar..." value={form.department_id} onChange={e => setForm({ ...form, department_id: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="position">Posición</Label>
              <Select id="position" options={positions} placeholder="Seleccionar..." value={form.position_id} onChange={e => setForm({ ...form, position_id: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="manager">Manager</Label>
              <Select id="manager" options={managers} placeholder="Seleccionar..." value={form.manager_id} onChange={e => setForm({ ...form, manager_id: e.target.value })} />
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button>
            <Button onClick={handleSave} disabled={saving || !form.employee_number || !form.first_name || !form.last_name || !form.hire_date}>
              {saving ? 'Guardando...' : editing ? 'Guardar Cambios' : 'Crear Empleado'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
