import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { ArrowLeft, Pencil, Plus, Trash2 } from 'lucide-react'
import api from '@/lib/api'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
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
  created_at: string
  updated_at: string
}

interface Contact {
  id: string
  contact_type: string
  contact_value: string
  is_primary: boolean
}

interface Address {
  id: string
  address_type: string
  street?: string
  street_number?: string
  apartment?: string
  city?: string
  state?: string
  country: string
  postal_code?: string
  is_primary: boolean
}

interface EmergencyContact {
  id: string
  name: string
  relationship?: string
  phone: string
  alt_phone?: string
  is_primary: boolean
}

interface HistoryEntry {
  id: string
  event_type: string
  old_value?: string
  new_value?: string
  description?: string
  performed_by?: string
  created_at: string
}

export default function EmployeeDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [emp, setEmp] = useState<Employee | null>(null)
  const [loading, setLoading] = useState(true)

  const [contacts, setContacts] = useState<Contact[]>([])
  const [addresses, setAddresses] = useState<Address[]>([])
  const [emergencyContacts, setEmergencyContacts] = useState<EmergencyContact[]>([])
  const [history, setHistory] = useState<HistoryEntry[]>([])
  const [activeTab, setActiveTab] = useState('info')

  const [showContactModal, setShowContactModal] = useState(false)
  const [contactForm, setContactForm] = useState({ contact_type: '', contact_value: '', is_primary: false })
  const [editingContact, setEditingContact] = useState<Contact | null>(null)

  const [showAddressModal, setShowAddressModal] = useState(false)
  const [addressForm, setAddressForm] = useState({ address_type: 'home', street: '', street_number: '', apartment: '', city: '', state: '', country: 'Argentina', postal_code: '', is_primary: false })
  const [editingAddress, setEditingAddress] = useState<Address | null>(null)

  const [showEmergencyModal, setShowEmergencyModal] = useState(false)
  const [emergencyForm, setEmergencyForm] = useState({ name: '', relationship: '', phone: '', alt_phone: '', is_primary: false })
  const [editingEmergency, setEditingEmergency] = useState<EmergencyContact | null>(null)

  useEffect(() => {
    if (!id) return
    setLoading(true)
    Promise.all([
      api.get(`/employees/${id}`),
      api.get(`/employees/${id}/contacts`),
      api.get(`/employees/${id}/addresses`),
      api.get(`/employees/${id}/emergency-contacts`),
      api.get(`/employees/${id}/history`),
    ])
      .then(([eRes, cRes, aRes, ecRes, hRes]) => {
        setEmp(eRes.data.data ?? eRes.data)
        setContacts(cRes.data.data ?? [])
        setAddresses(aRes.data.data ?? [])
        setEmergencyContacts(ecRes.data.data ?? [])
        setHistory(hRes.data.data ?? [])
      })
      .catch(() => navigate('/employees'))
      .finally(() => setLoading(false))
  }, [id])

  if (loading || !emp) {
    return <div className="p-6 text-center text-slate-500">Cargando...</div>
  }

  const tabs = [
    { key: 'info', label: 'Información' },
    { key: 'contacts', label: 'Contactos' },
    { key: 'addresses', label: 'Direcciones' },
    { key: 'emergency', label: 'Emergencia' },
    { key: 'history', label: 'Historial' },
  ]

  return (
    <div>
      <div className="flex items-center gap-3 mb-6">
        <Button variant="ghost" size="sm" onClick={() => navigate('/employees')}>
          <ArrowLeft size={16} />
        </Button>
        <h1 className="text-2xl font-bold text-slate-900">
          {emp.first_name} {emp.last_name}
        </h1>
        <span className={`ml-2 inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${emp.status === 'active' ? 'bg-emerald-50 text-emerald-700' : 'bg-slate-100 text-slate-600'}`}>
          {emp.status === 'active' ? 'Activo' : 'Inactivo'}
        </span>
      </div>

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

      {activeTab === 'info' && <InfoTab emp={emp} />}
      {activeTab === 'contacts' && (
        <ContactsTab
          contacts={contacts}
          onAdd={() => {
            setEditingContact(null)
            setContactForm({ contact_type: '', contact_value: '', is_primary: false })
            setShowContactModal(true)
          }}
          onEdit={(c) => {
            setEditingContact(c)
            setContactForm({ contact_type: c.contact_type, contact_value: c.contact_value, is_primary: c.is_primary })
            setShowContactModal(true)
          }}
          onDelete={async (c) => {
            const updated = contacts.filter((x) => x.id !== c.id)
            await api.put(`/employees/${id}/contacts`, updated)
            setContacts(updated)
          }}
        />
      )}
      {activeTab === 'addresses' && (
        <AddressesTab
          addresses={addresses}
          onAdd={() => {
            setEditingAddress(null)
            setAddressForm({ address_type: 'home', street: '', street_number: '', apartment: '', city: '', state: '', country: 'Argentina', postal_code: '', is_primary: false })
            setShowAddressModal(true)
          }}
          onEdit={(a) => {
            setEditingAddress(a)
            setAddressForm({
              address_type: a.address_type,
              street: a.street ?? '',
              street_number: a.street_number ?? '',
              apartment: a.apartment ?? '',
              city: a.city ?? '',
              state: a.state ?? '',
              country: a.country,
              postal_code: a.postal_code ?? '',
              is_primary: a.is_primary,
            })
            setShowAddressModal(true)
          }}
          onDelete={async (a) => {
            const updated = addresses.filter((x) => x.id !== a.id)
            await api.put(`/employees/${id}/addresses`, updated)
            setAddresses(updated)
          }}
        />
      )}
      {activeTab === 'emergency' && (
        <EmergencyTab
          items={emergencyContacts}
          onAdd={() => {
            setEditingEmergency(null)
            setEmergencyForm({ name: '', relationship: '', phone: '', alt_phone: '', is_primary: false })
            setShowEmergencyModal(true)
          }}
          onEdit={(e) => {
            setEditingEmergency(e)
            setEmergencyForm({ name: e.name, relationship: e.relationship ?? '', phone: e.phone, alt_phone: e.alt_phone ?? '', is_primary: e.is_primary })
            setShowEmergencyModal(true)
          }}
          onDelete={async (e) => {
            const updated = emergencyContacts.filter((x) => x.id !== e.id)
            await api.put(`/employees/${id}/emergency-contacts`, updated)
            setEmergencyContacts(updated)
          }}
        />
      )}
      {activeTab === 'history' && <HistoryTab items={history} />}

      <ContactModal
        open={showContactModal}
        onOpenChange={setShowContactModal}
        form={contactForm}
        setForm={setContactForm}
        onSave={async () => {
          const updated = editingContact
            ? contacts.map((c) => (c.id === editingContact.id ? { ...c, ...contactForm } : c))
            : [...contacts, { id: crypto.randomUUID(), ...contactForm }]
          await api.put(`/employees/${id}/contacts`, updated)
          setContacts(updated)
          setShowContactModal(false)
        }}
      />

      <AddressModal
        open={showAddressModal}
        onOpenChange={setShowAddressModal}
        form={addressForm}
        setForm={setAddressForm}
        onSave={async () => {
          const updated = editingAddress
            ? addresses.map((a) => (a.id === editingAddress.id ? { ...a, ...addressForm } : a))
            : [...addresses, { id: crypto.randomUUID(), ...addressForm }]
          await api.put(`/employees/${id}/addresses`, updated)
          setAddresses(updated)
          setShowAddressModal(false)
        }}
      />

      <EmergencyModal
        open={showEmergencyModal}
        onOpenChange={setShowEmergencyModal}
        form={emergencyForm}
        setForm={setEmergencyForm}
        onSave={async () => {
          const updated = editingEmergency
            ? emergencyContacts.map((e) => (e.id === editingEmergency.id ? { ...e, ...emergencyForm } : e))
            : [...emergencyContacts, { id: crypto.randomUUID(), ...emergencyForm }]
          await api.put(`/employees/${id}/emergency-contacts`, updated)
          setEmergencyContacts(updated)
          setShowEmergencyModal(false)
        }}
      />
    </div>
  )
}

function InfoTab({ emp }: { emp: Employee }) {
  const fields = [
    { label: 'N° Empleado', value: emp.employee_number },
    { label: 'Nombre', value: `${emp.first_name} ${emp.last_name}` },
    { label: 'DNI', value: emp.dni },
    { label: 'Email', value: emp.email },
    { label: 'Teléfono', value: emp.phone },
    { label: 'Fecha Nacimiento', value: emp.birth_date },
    { label: 'Sucursal', value: emp.branch_name },
    { label: 'Departamento', value: emp.department_name },
    { label: 'Posición', value: emp.position_name },
    { label: 'Manager', value: emp.manager_name },
    { label: 'F. Ingreso', value: emp.hire_date },
    { label: 'F. Egreso', value: emp.termination_date },
  ]

  return (
    <Card>
      <CardContent className="p-6">
        <div className="grid grid-cols-2 md:grid-cols-3 gap-6">
          {fields.map((f) => (
            <div key={f.label}>
              <p className="text-xs font-medium text-slate-500 uppercase tracking-wider mb-1">{f.label}</p>
              <p className="text-sm text-slate-900">{f.value || '-'}</p>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}

function ContactsTab({ contacts, onAdd, onEdit, onDelete }: {
  contacts: Contact[]
  onAdd: () => void
  onEdit: (c: Contact) => void
  onDelete: (c: Contact) => void
}) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle className="text-base">Contactos</CardTitle>
        <Button size="sm" onClick={onAdd}><Plus size={14} /> Agregar</Button>
      </CardHeader>
      <CardContent className="p-0">
        {contacts.length === 0 ? (
          <div className="p-6 text-center text-slate-500">Sin contactos registrados</div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-slate-200 bg-slate-50">
                <th className="text-left px-4 py-3 font-medium text-slate-600">Tipo</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Valor</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Principal</th>
                <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
              </tr>
            </thead>
            <tbody>
              {contacts.map((c) => (
                <tr key={c.id} className="border-b border-slate-100 hover:bg-slate-50">
                  <td className="px-4 py-3 text-slate-900 capitalize">{c.contact_type}</td>
                  <td className="px-4 py-3 text-slate-600">{c.contact_value}</td>
                  <td className="px-4 py-3">{c.is_primary ? <span className="text-emerald-600 font-medium">Sí</span> : 'No'}</td>
                  <td className="px-4 py-3 text-right">
                    <Button variant="ghost" size="sm" onClick={() => onEdit(c)}><Pencil size={14} /></Button>
                    <Button variant="ghost" size="sm" className="text-red-500" onClick={() => onDelete(c)}><Trash2 size={14} /></Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </CardContent>
    </Card>
  )
}

function AddressesTab({ addresses, onAdd, onEdit, onDelete }: {
  addresses: Address[]
  onAdd: () => void
  onEdit: (a: Address) => void
  onDelete: (a: Address) => void
}) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle className="text-base">Direcciones</CardTitle>
        <Button size="sm" onClick={onAdd}><Plus size={14} /> Agregar</Button>
      </CardHeader>
      <CardContent className="p-0">
        {addresses.length === 0 ? (
          <div className="p-6 text-center text-slate-500">Sin direcciones registradas</div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-slate-200 bg-slate-50">
                <th className="text-left px-4 py-3 font-medium text-slate-600">Tipo</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Dirección</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Ciudad</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">País</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Principal</th>
                <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
              </tr>
            </thead>
            <tbody>
              {addresses.map((a) => (
                <tr key={a.id} className="border-b border-slate-100 hover:bg-slate-50">
                  <td className="px-4 py-3 text-slate-900 capitalize">{a.address_type}</td>
                  <td className="px-4 py-3 text-slate-600">{[a.street, a.street_number, a.apartment].filter(Boolean).join(' ')}</td>
                  <td className="px-4 py-3 text-slate-600">{a.city}</td>
                  <td className="px-4 py-3 text-slate-600">{a.country}</td>
                  <td className="px-4 py-3">{a.is_primary ? <span className="text-emerald-600 font-medium">Sí</span> : 'No'}</td>
                  <td className="px-4 py-3 text-right">
                    <Button variant="ghost" size="sm" onClick={() => onEdit(a)}><Pencil size={14} /></Button>
                    <Button variant="ghost" size="sm" className="text-red-500" onClick={() => onDelete(a)}><Trash2 size={14} /></Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </CardContent>
    </Card>
  )
}

function EmergencyTab({ items, onAdd, onEdit, onDelete }: {
  items: EmergencyContact[]
  onAdd: () => void
  onEdit: (e: EmergencyContact) => void
  onDelete: (e: EmergencyContact) => void
}) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle className="text-base">Contactos de Emergencia</CardTitle>
        <Button size="sm" onClick={onAdd}><Plus size={14} /> Agregar</Button>
      </CardHeader>
      <CardContent className="p-0">
        {items.length === 0 ? (
          <div className="p-6 text-center text-slate-500">Sin contactos de emergencia</div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-slate-200 bg-slate-50">
                <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Parentesco</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Teléfono</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Tel. Alternativo</th>
                <th className="text-left px-4 py-3 font-medium text-slate-600">Principal</th>
                <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
              </tr>
            </thead>
            <tbody>
              {items.map((e) => (
                <tr key={e.id} className="border-b border-slate-100 hover:bg-slate-50">
                  <td className="px-4 py-3 font-medium text-slate-900">{e.name}</td>
                  <td className="px-4 py-3 text-slate-600">{e.relationship || '-'}</td>
                  <td className="px-4 py-3 text-slate-600">{e.phone}</td>
                  <td className="px-4 py-3 text-slate-600">{e.alt_phone || '-'}</td>
                  <td className="px-4 py-3">{e.is_primary ? <span className="text-emerald-600 font-medium">Sí</span> : 'No'}</td>
                  <td className="px-4 py-3 text-right">
                    <Button variant="ghost" size="sm" onClick={() => onEdit(e)}><Pencil size={14} /></Button>
                    <Button variant="ghost" size="sm" className="text-red-500" onClick={() => onDelete(e)}><Trash2 size={14} /></Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </CardContent>
    </Card>
  )
}

function HistoryTab({ items }: { items: HistoryEntry[] }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Historial</CardTitle>
      </CardHeader>
      <CardContent className="p-0">
        {items.length === 0 ? (
          <div className="p-6 text-center text-slate-500">Sin historial registrado</div>
        ) : (
          <div className="divide-y divide-slate-100">
            {items.map((h) => (
              <div key={h.id} className="px-4 py-3">
                <div className="flex items-start justify-between">
                  <div>
                    <span className="inline-flex px-2 py-0.5 rounded text-xs font-medium bg-slate-100 text-slate-700 mb-1">{h.event_type}</span>
                    <p className="text-sm text-slate-700">{h.description}</p>
                    {(h.old_value || h.new_value) && (
                      <p className="text-xs text-slate-500 mt-0.5">
                        {h.old_value && <span className="line-through text-red-500 mr-2">{h.old_value}</span>}
                        {h.new_value && <span className="text-emerald-600">{h.new_value}</span>}
                      </p>
                    )}
                  </div>
                  <span className="text-xs text-slate-400 whitespace-nowrap ml-4">
                    {new Date(h.created_at).toLocaleDateString('es-AR')}
                  </span>
                </div>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  )
}

function ContactModal({ open, onOpenChange, form, setForm, onSave }: {
  open: boolean
  onOpenChange: (v: boolean) => void
  form: { contact_type: string; contact_value: string; is_primary: boolean }
  setForm: (f: any) => void
  onSave: () => Promise<void>
}) {
  const [saving, setSaving] = useState(false)
  const handleSave = async () => {
    setSaving(true)
    try { await onSave() } finally { setSaving(false) }
  }
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader><DialogTitle>Contacto</DialogTitle></DialogHeader>
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Tipo</Label>
            <select
              className="flex h-10 w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm"
              value={form.contact_type}
              onChange={e => setForm({ ...form, contact_type: e.target.value })}
            >
              <option value="">Seleccionar...</option>
              <option value="email">Email</option>
              <option value="phone">Teléfono</option>
              <option value="mobile">Celular</option>
              <option value="other">Otro</option>
            </select>
          </div>
          <div className="space-y-2">
            <Label>Valor</Label>
            <Input value={form.contact_value} onChange={e => setForm({ ...form, contact_value: e.target.value })} />
          </div>
          <label className="flex items-center gap-2 text-sm">
            <input type="checkbox" checked={form.is_primary} onChange={e => setForm({ ...form, is_primary: e.target.checked })} className="rounded" />
            Contacto principal
          </label>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>Cancelar</Button>
          <Button onClick={handleSave} disabled={saving || !form.contact_type || !form.contact_value}>
            {saving ? 'Guardando...' : 'Guardar'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

function AddressModal({ open, onOpenChange, form, setForm, onSave }: {
  open: boolean
  onOpenChange: (v: boolean) => void
  form: { address_type: string; street: string; street_number: string; apartment: string; city: string; state: string; country: string; postal_code: string; is_primary: boolean }
  setForm: (f: any) => void
  onSave: () => Promise<void>
}) {
  const [saving, setSaving] = useState(false)
  const handleSave = async () => {
    setSaving(true)
    try { await onSave() } finally { setSaving(false) }
  }
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader><DialogTitle>Dirección</DialogTitle></DialogHeader>
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label>Tipo</Label>
            <select className="flex h-10 w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm" value={form.address_type} onChange={e => setForm({ ...form, address_type: e.target.value })}>
              <option value="home">Casa</option>
              <option value="work">Trabajo</option>
              <option value="other">Otro</option>
            </select>
          </div>
          <div className="space-y-2">
            <Label>País</Label>
            <Input value={form.country} onChange={e => setForm({ ...form, country: e.target.value })} />
          </div>
          <div className="space-y-2">
            <Label>Calle</Label>
            <Input value={form.street} onChange={e => setForm({ ...form, street: e.target.value })} />
          </div>
          <div className="space-y-2">
            <Label>Número</Label>
            <Input value={form.street_number} onChange={e => setForm({ ...form, street_number: e.target.value })} />
          </div>
          <div className="space-y-2">
            <Label>Piso/Dto</Label>
            <Input value={form.apartment} onChange={e => setForm({ ...form, apartment: e.target.value })} />
          </div>
          <div className="space-y-2">
            <Label>Código Postal</Label>
            <Input value={form.postal_code} onChange={e => setForm({ ...form, postal_code: e.target.value })} />
          </div>
          <div className="space-y-2">
            <Label>Ciudad</Label>
            <Input value={form.city} onChange={e => setForm({ ...form, city: e.target.value })} />
          </div>
          <div className="space-y-2">
            <Label>Provincia</Label>
            <Input value={form.state} onChange={e => setForm({ ...form, state: e.target.value })} />
          </div>
          <div className="col-span-2">
            <label className="flex items-center gap-2 text-sm">
              <input type="checkbox" checked={form.is_primary} onChange={e => setForm({ ...form, is_primary: e.target.checked })} className="rounded" />
              Dirección principal
            </label>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>Cancelar</Button>
          <Button onClick={handleSave} disabled={saving}>{saving ? 'Guardando...' : 'Guardar'}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

function EmergencyModal({ open, onOpenChange, form, setForm, onSave }: {
  open: boolean
  onOpenChange: (v: boolean) => void
  form: { name: string; relationship: string; phone: string; alt_phone: string; is_primary: boolean }
  setForm: (f: any) => void
  onSave: () => Promise<void>
}) {
  const [saving, setSaving] = useState(false)
  const handleSave = async () => {
    setSaving(true)
    try { await onSave() } finally { setSaving(false) }
  }
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader><DialogTitle>Contacto de Emergencia</DialogTitle></DialogHeader>
        <div className="grid grid-cols-2 gap-4">
          <div className="col-span-2 space-y-2">
            <Label>Nombre completo *</Label>
            <Input value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} />
          </div>
          <div className="space-y-2">
            <Label>Parentesco</Label>
            <select className="flex h-10 w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm" value={form.relationship} onChange={e => setForm({ ...form, relationship: e.target.value })}>
              <option value="">Seleccionar...</option>
              <option value="spouse">Cónyuge</option>
              <option value="parent">Padre/Madre</option>
              <option value="sibling">Hermano/a</option>
              <option value="child">Hijo/a</option>
              <option value="friend">Amigo/a</option>
              <option value="other">Otro</option>
            </select>
          </div>
          <div className="space-y-2">
            <Label>Teléfono *</Label>
            <Input value={form.phone} onChange={e => setForm({ ...form, phone: e.target.value })} />
          </div>
          <div className="space-y-2">
            <Label>Tel. Alternativo</Label>
            <Input value={form.alt_phone} onChange={e => setForm({ ...form, alt_phone: e.target.value })} />
          </div>
          <div className="col-span-2">
            <label className="flex items-center gap-2 text-sm">
              <input type="checkbox" checked={form.is_primary} onChange={e => setForm({ ...form, is_primary: e.target.checked })} className="rounded" />
              Contacto principal
            </label>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>Cancelar</Button>
          <Button onClick={handleSave} disabled={saving || !form.name || !form.phone}>{saving ? 'Guardando...' : 'Guardar'}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
