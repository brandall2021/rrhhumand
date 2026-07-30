import { useEffect, useState } from 'react'
import api from '@/lib/api'
import { Card, CardContent } from '@/components/ui/card'
import { Users, GraduationCap, Receipt, Handshake, DollarSign, Clock } from 'lucide-react'

interface Stats {
  employee_count: number
  active_courses: number
  pending_expenses: number
  open_recruitments: number
  payroll_runs: number
  pending_approvals: number
}

const defaultStats: Stats = {
  employee_count: 0,
  active_courses: 0,
  pending_expenses: 0,
  open_recruitments: 0,
  payroll_runs: 0,
  pending_approvals: 0,
}

const cards = [
  { key: 'employee_count', label: 'Empleados', icon: Users, color: 'text-brand-600', bg: 'bg-brand-50' },
  { key: 'active_courses', label: 'Cursos activos', icon: GraduationCap, color: 'text-emerald-600', bg: 'bg-emerald-50' },
  { key: 'pending_expenses', label: 'Gastos pendientes', icon: Receipt, color: 'text-amber-600', bg: 'bg-amber-50' },
  { key: 'open_recruitments', label: 'Reclutamientos abiertos', icon: Handshake, color: 'text-violet-600', bg: 'bg-violet-50' },
  { key: 'payroll_runs', label: 'Nóminas del mes', icon: DollarSign, color: 'text-cyan-600', bg: 'bg-cyan-50' },
  { key: 'pending_approvals', label: 'Aprobaciones pendientes', icon: Clock, color: 'text-rose-600', bg: 'bg-rose-50' },
]

export default function DashboardPage() {
  const [stats, setStats] = useState<Stats>(defaultStats)

  useEffect(() => {
    api.get('/compensation/dashboard')
      .then((res) => setStats(res.data.data))
      .catch(() => {})
  }, [])

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-slate-900">Dashboard</h1>
        <p className="text-slate-500 mt-1">Resumen general de la organización</p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {cards.map((card) => {
          const Icon = card.icon
          const value = stats[card.key as keyof Stats]
          return (
            <Card key={card.key}>
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-slate-500">{card.label}</p>
                    <p className="text-3xl font-bold text-slate-900 mt-1">{value}</p>
                  </div>
                  <div className={`p-3 rounded-xl ${card.bg}`}>
                    <Icon className={`h-6 w-6 ${card.color}`} />
                  </div>
                </div>
              </CardContent>
            </Card>
          )
        })}
      </div>
    </div>
  )
}
