import { useEffect, useState } from 'react'
import { Pencil } from 'lucide-react'
import api from '@/lib/api'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'

interface Profile {
  id: string
  employee_number: string
  first_name: string
  last_name: string
  dni?: string
  email?: string
  phone?: string
  birth_date?: string
  photo_url?: string
  hire_date: string
  tenure: string
  status: string
  branch_name?: string
  department_name?: string
  position_name?: string
  position_level?: number
  manager_name?: string
  phone_contact?: string
  personal_email?: string
  created_at: string
}

export default function ProfilePage() {
  const [profile, setProfile] = useState<Profile | null>(null)
  const [loading, setLoading] = useState(true)
  const [showEdit, setShowEdit] = useState(false)
  const [form, setForm] = useState({ phone: '', photo_url: '' })
  const [saving, setSaving] = useState(false)

  const fetch = async () => {
    setLoading(true)
    try {
      const res = await api.get('/me/profile')
      setProfile(res.data.data ?? null)
    } catch { setProfile(null) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetch() }, [])

  const openEdit = () => {
    setForm({
      phone: profile?.phone ?? '',
      photo_url: profile?.photo_url ?? '',
    })
    setShowEdit(true)
  }

  const handleSave = async () => {
    setSaving(true)
    try {
      await api.put('/me/profile', {
        phone: form.phone || null,
        photo_url: form.photo_url || null,
      })
      setShowEdit(false)
      fetch()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al guardar el perfil')
    } finally { setSaving(false) }
  }

  if (loading) return <div className="text-center text-slate-500 py-12">Cargando...</div>
  if (!profile) return <div className="text-center text-slate-500 py-12">No se encontró un perfil de empleado para este usuario</div>

  const initials = `${profile.first_name} ${profile.last_name}`.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)

  const fields: { label: string; value?: string }[] = [
    { label: 'Legajo', value: profile.employee_number },
    { label: 'DNI', value: profile.dni },
    { label: 'Email corporativo', value: profile.email },
    { label: 'Email personal', value: profile.personal_email },
    { label: 'Teléfono', value: profile.phone },
    { label: 'Contacto de emergencia', value: profile.phone_contact },
    { label: 'Fecha de nacimiento', value: profile.birth_date ? new Date(profile.birth_date).toLocaleDateString('es-AR') : undefined },
    { label: 'Sucursal', value: profile.branch_name },
    { label: 'Departamento', value: profile.department_name },
    { label: 'Posición', value: profile.position_name },
    { label: 'Nivel', value: profile.position_level != null ? String(profile.position_level) : undefined },
    { label: 'Jefe directo', value: profile.manager_name },
    { label: 'Fechas de alta', value: `${new Date(profile.hire_date).toLocaleDateString('es-AR')} · ${profile.tenure}` },
  ]

  return (
    <div className="max-w-3xl">
      <h1 className="text-2xl font-bold text-slate-900 mb-6">Mi Perfil</h1>

      <Card className="mb-6">
        <CardContent className="p-6">
          <div className="flex items-start justify-between">
            <div className="flex items-center gap-4">
              <Avatar className="h-20 w-20">
                {profile.photo_url && <AvatarImage src={profile.photo_url} alt="" />}
                <AvatarFallback className="bg-brand-600 text-white text-xl">{initials}</AvatarFallback>
              </Avatar>
              <div>
                <h2 className="text-xl font-semibold text-slate-900">
                  {profile.first_name} {profile.last_name}
                </h2>
                <p className="text-sm text-slate-500">
                  {[profile.position_name, profile.department_name].filter(Boolean).join(' · ') || 'Sin posición asignada'}
                </p>
                <span className="inline-flex mt-2 px-2 py-0.5 rounded-full text-xs font-medium bg-emerald-50 text-emerald-700">
                  Activo
                </span>
              </div>
            </div>
            <Button variant="outline" size="sm" onClick={openEdit}>
              <Pencil size={14} /> Editar
            </Button>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="p-6">
          <h3 className="font-semibold text-slate-900 mb-4">Datos de contacto y organización</h3>
          <dl className="grid grid-cols-1 md:grid-cols-2 gap-x-8 gap-y-4 text-sm">
            {fields.map((f) => (
              <div key={f.label} className="flex justify-between gap-4 border-b border-slate-100 pb-2">
                <dt className="text-slate-500">{f.label}</dt>
                <dd className="font-medium text-slate-900 text-right">{f.value || '-'}</dd>
              </div>
            ))}
          </dl>
        </CardContent>
      </Card>

      <Dialog open={showEdit} onOpenChange={setShowEdit}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Editar Mi Perfil</DialogTitle>
            <DialogDescription>Solo podés actualizar teléfono y foto de perfil</DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label>Teléfono</Label>
              <Input value={form.phone} onChange={e => setForm({ ...form, phone: e.target.value })} placeholder="+54 11 1234 5678" />
            </div>
            <div className="space-y-2">
              <Label>URL de foto</Label>
              <Input value={form.photo_url} onChange={e => setForm({ ...form, photo_url: e.target.value })} placeholder="https://..." />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowEdit(false)}>Cancelar</Button>
            <Button onClick={handleSave} disabled={saving}>
              {saving ? 'Guardando...' : 'Guardar Cambios'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
