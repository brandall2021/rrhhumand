import { useEffect, useState } from 'react'
import { Plus, Pencil, Trash2 } from 'lucide-react'
import api from '@/lib/api'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'

interface Company {
  id: string
  name: string
  slug: string
  logo_url?: string
  plan: string
  active: boolean
  created_at: string
}

interface Branch {
  id: string
  company_id: string
  name: string
  code?: string
  address?: string
  city?: string
  state?: string
  country: string
  phone?: string
  email?: string
  timezone: string
  active: boolean
}

interface Role {
  id: string
  name: string
  description?: string
  created_at: string
}

export default function SettingsPage() {
  const [activeTab, setActiveTab] = useState('companies')

  const tabs = [
    { key: 'companies', label: 'Empresas' },
    { key: 'branches', label: 'Sucursales' },
    { key: 'roles', label: 'Roles' },
  ]

  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-900 mb-6">Configuración</h1>

      <div className="flex gap-1 mb-6 border-b border-slate-200">
        {tabs.map((tab) => (
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

      {activeTab === 'companies' && <CompaniesTab />}
      {activeTab === 'branches' && <BranchesTab />}
      {activeTab === 'roles' && <RolesTab />}
    </div>
  )
}

function CompaniesTab() {
  const [companies, setCompanies] = useState<Company[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [editing, setEditing] = useState<Company | null>(null)
  const [form, setForm] = useState({ name: '', slug: '', plan: 'free' })
  const [saving, setSaving] = useState(false)

  const fetch = async () => {
    setLoading(true)
    try {
      const res = await api.get('/companies')
      setCompanies(res.data.data ?? [])
    } catch { setCompanies([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetch() }, [])

  const openCreate = () => {
    setEditing(null)
    setForm({ name: '', slug: '', plan: 'free' })
    setShowModal(true)
  }

  const openEdit = (c: Company) => {
    setEditing(c)
    setForm({ name: c.name, slug: c.slug, plan: c.plan })
    setShowModal(true)
  }

  const handleSave = async () => {
    setSaving(true)
    try {
      if (editing) {
        await api.put(`/companies/${editing.id}`, {
          name: form.name,
          slug: form.slug,
          plan: form.plan,
        })
      } else {
        await api.post('/companies', { name: form.name, slug: form.slug })
      }
      setShowModal(false)
      fetch()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar la empresa')
    } finally { setSaving(false) }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-slate-500">{companies.length} empresas</p>
        <Button size="sm" onClick={openCreate}><Plus size={14} /> Nueva Empresa</Button>
      </div>

      <Card>
        <CardContent className="p-0">
          {loading ? (
            <div className="p-6 text-center text-slate-500">Cargando...</div>
          ) : companies.length === 0 ? (
            <div className="p-6 text-center text-slate-500">No hay empresas registradas</div>
          ) : (
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-slate-200 bg-slate-50">
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Slug</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Plan</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                  <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                </tr>
              </thead>
              <tbody>
                {companies.map((c) => (
                  <tr key={c.id} className="border-b border-slate-100 hover:bg-slate-50">
                    <td className="px-4 py-3 font-medium text-slate-900">{c.name}</td>
                    <td className="px-4 py-3 text-slate-600">{c.slug}</td>
                    <td className="px-4 py-3 text-slate-600">{c.plan}</td>
                    <td className="px-4 py-3">{c.active ? <span className="text-emerald-600 font-medium">Activa</span> : 'Inactiva'}</td>
                    <td className="px-4 py-3 text-right">
                      <Button variant="ghost" size="sm" onClick={() => openEdit(c)}><Pencil size={14} /></Button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </CardContent>
      </Card>

      <Dialog open={showModal} onOpenChange={setShowModal}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{editing ? 'Editar Empresa' : 'Nueva Empresa'}</DialogTitle>
            <DialogDescription>El slug se usa para identificar el dominio de la empresa</DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label>Nombre *</Label>
              <Input value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label>Slug *</Label>
              <Input value={form.slug} onChange={e => setForm({ ...form, slug: e.target.value })} placeholder="mi-empresa" />
            </div>
            {editing && (
              <div className="space-y-2">
                <Label>Plan</Label>
                <select
                  className="flex h-10 w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm"
                  value={form.plan}
                  onChange={e => setForm({ ...form, plan: e.target.value })}
                >
                  <option value="free">Free</option>
                  <option value="pro">Pro</option>
                  <option value="enterprise">Enterprise</option>
                </select>
              </div>
            )}
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button>
            <Button onClick={handleSave} disabled={saving || !form.name || !form.slug}>
              {saving ? 'Guardando...' : editing ? 'Guardar Cambios' : 'Crear'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function BranchesTab() {
  const [branches, setBranches] = useState<Branch[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [editing, setEditing] = useState<Branch | null>(null)
  const [form, setForm] = useState({
    name: '', code: '', address: '', city: '', state: '', country: 'Argentina', phone: '', email: '', timezone: 'America/Argentina/Buenos_Aires',
  })
  const [saving, setSaving] = useState(false)

  const fetch = async () => {
    setLoading(true)
    try {
      const res = await api.get('/branches', { params: { limit: '100' } })
      setBranches(res.data.data ?? [])
    } catch { setBranches([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetch() }, [])

  const openCreate = () => {
    setEditing(null)
    setForm({ name: '', code: '', address: '', city: '', state: '', country: 'Argentina', phone: '', email: '', timezone: 'America/Argentina/Buenos_Aires' })
    setShowModal(true)
  }

  const openEdit = (b: Branch) => {
    setEditing(b)
    setForm({
      name: b.name, code: b.code ?? '', address: b.address ?? '', city: b.city ?? '',
      state: b.state ?? '', country: b.country, phone: b.phone ?? '', email: b.email ?? '', timezone: b.timezone,
    })
    setShowModal(true)
  }

  const handleSave = async () => {
    setSaving(true)
    try {
      const payload: Record<string, any> = {
        name: form.name, code: form.code || null, address: form.address || null,
        city: form.city || null, state: form.state || null, country: form.country || null,
        phone: form.phone || null, email: form.email || null, timezone: form.timezone || null,
      }
      if (editing) {
        await api.put(`/branches/${editing.id}`, payload)
      } else {
        await api.post('/branches', payload)
      }
      setShowModal(false)
      fetch()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar la sucursal')
    } finally { setSaving(false) }
  }

  const handleDelete = async (b: Branch) => {
    if (!confirm(`¿Eliminar la sucursal "${b.name}"?`)) return
    try {
      await api.delete(`/branches/${b.id}`)
      fetch()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al eliminar')
    }
  }

  const handleToggle = async (b: Branch) => {
    try {
      await api.put(`/branches/${b.id}`, { active: !b.active })
      fetch()
    } catch { alert('Error al actualizar la sucursal') }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-slate-500">{branches.length} sucursales</p>
        <Button size="sm" onClick={openCreate}><Plus size={14} /> Nueva Sucursal</Button>
      </div>

      <Card>
        <CardContent className="p-0">
          {loading ? (
            <div className="p-6 text-center text-slate-500">Cargando...</div>
          ) : branches.length === 0 ? (
            <div className="p-6 text-center text-slate-500">No hay sucursales registradas</div>
          ) : (
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-slate-200 bg-slate-50">
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Código</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Ciudad</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                  <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                </tr>
              </thead>
              <tbody>
                {branches.map((b) => (
                  <tr key={b.id} className="border-b border-slate-100 hover:bg-slate-50">
                    <td className="px-4 py-3 font-medium text-slate-900">{b.name}</td>
                    <td className="px-4 py-3 text-slate-600">{b.code || '-'}</td>
                    <td className="px-4 py-3 text-slate-600">{[b.city, b.state].filter(Boolean).join(', ') || '-'}</td>
                    <td className="px-4 py-3">
                      <button onClick={() => handleToggle(b)} className={b.active ? 'text-emerald-600 font-medium hover:underline' : 'text-slate-400 font-medium hover:underline'}>
                        {b.active ? 'Activa' : 'Inactiva'}
                      </button>
                    </td>
                    <td className="px-4 py-3 text-right whitespace-nowrap">
                      <Button variant="ghost" size="sm" onClick={() => openEdit(b)}><Pencil size={14} /></Button>
                      <Button variant="ghost" size="sm" className="text-red-500" onClick={() => handleDelete(b)}><Trash2 size={14} /></Button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </CardContent>
      </Card>

      <Dialog open={showModal} onOpenChange={setShowModal}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{editing ? 'Editar Sucursal' : 'Nueva Sucursal'}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Nombre *</Label>
                <Input value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label>Código</Label>
                <Input value={form.code} onChange={e => setForm({ ...form, code: e.target.value })} />
              </div>
            </div>
            <div className="space-y-2">
              <Label>Dirección</Label>
              <Input value={form.address} onChange={e => setForm({ ...form, address: e.target.value })} />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Ciudad</Label>
                <Input value={form.city} onChange={e => setForm({ ...form, city: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label>Provincia / Estado</Label>
                <Input value={form.state} onChange={e => setForm({ ...form, state: e.target.value })} />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>País</Label>
                <Input value={form.country} onChange={e => setForm({ ...form, country: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label>Teléfono</Label>
                <Input value={form.phone} onChange={e => setForm({ ...form, phone: e.target.value })} />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Email</Label>
                <Input value={form.email} onChange={e => setForm({ ...form, email: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label>Zona horaria</Label>
                <Input value={form.timezone} onChange={e => setForm({ ...form, timezone: e.target.value })} />
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowModal(false)}>Cancelar</Button>
            <Button onClick={handleSave} disabled={saving || !form.name}>
              {saving ? 'Guardando...' : editing ? 'Guardar Cambios' : 'Crear'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function RolesTab() {
  const [roles, setRoles] = useState<Role[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    setLoading(true)
    api.get('/roles')
      .then((res) => setRoles(res.data.data ?? []))
      .catch(() => setRoles([]))
      .finally(() => setLoading(false))
  }, [])

  return (
    <div>
      <p className="text-sm text-slate-500 mb-4">Roles disponibles en el sistema (solo lectura)</p>

      <Card>
        <CardContent className="p-0">
          {loading ? (
            <div className="p-6 text-center text-slate-500">Cargando...</div>
          ) : roles.length === 0 ? (
            <div className="p-6 text-center text-slate-500">No hay roles registrados</div>
          ) : (
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-slate-200 bg-slate-50">
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                  <th className="text-left px-4 py-3 font-medium text-slate-600">Descripción</th>
                </tr>
              </thead>
              <tbody>
                {roles.map((r) => (
                  <tr key={r.id} className="border-b border-slate-100 hover:bg-slate-50">
                    <td className="px-4 py-3 font-medium text-slate-900">{r.name}</td>
                    <td className="px-4 py-3 text-slate-600">{r.description || '-'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
