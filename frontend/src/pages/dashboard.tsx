import { useEffect, useState } from 'react'
import api from '@/lib/api'
import { Card, CardContent } from '@/components/ui/card'
import { Users, Banknote, PiggyBank, Gift, Wallet, BarChart3, Building2, Briefcase, FileText, Receipt } from 'lucide-react'

interface Stats {
  total_salary_cost: string
  average_compensation: string
  benefit_cost: string
  total_bonuses: string
  budget_total: string
  budget_used: string
  pending_proposals: number
  employees_out_of_band: number
  active_reviews: number
  currency: string
}

const defaultStats: Stats = {
  total_salary_cost: '0',
  average_compensation: '0',
  benefit_cost: '0',
  total_bonuses: '0',
  budget_total: '0',
  budget_used: '0',
  pending_proposals: 0,
  employees_out_of_band: 0,
  active_reviews: 0,
  currency: 'USD',
}

const compensationCards = [
  { key: 'total_salary_cost', label: 'Costo Salarial Total', icon: Users, color: 'text-brand-600', bg: 'bg-brand-50' },
  { key: 'average_compensation', label: 'Compensación Promedio', icon: Banknote, color: 'text-emerald-600', bg: 'bg-emerald-50' },
  { key: 'benefit_cost', label: 'Costo en Beneficios', icon: Gift, color: 'text-amber-600', bg: 'bg-amber-50' },
  { key: 'total_bonuses', label: 'Bonos Totales', icon: PiggyBank, color: 'text-violet-600', bg: 'bg-violet-50' },
  { key: 'budget_total', label: 'Presupuesto Total', icon: Wallet, color: 'text-cyan-600', bg: 'bg-cyan-50' },
  { key: 'budget_used', label: 'Presupuesto Utilizado', icon: BarChart3, color: 'text-rose-600', bg: 'bg-rose-50' },
]

const generalCards = [
  { key: 'employees', label: 'Empleados', icon: Users, color: 'text-brand-600', bg: 'bg-brand-50' },
  { key: 'departments', label: 'Departamentos', icon: Building2, color: 'text-emerald-600', bg: 'bg-emerald-50' },
  { key: 'positions', label: 'Posiciones', icon: Briefcase, color: 'text-amber-600', bg: 'bg-amber-50' },
  { key: 'documents', label: 'Documentos', icon: FileText, color: 'text-violet-600', bg: 'bg-violet-50' },
  { key: 'expenses', label: 'Gastos', icon: Receipt, color: 'text-cyan-600', bg: 'bg-cyan-50' },
]

export default function DashboardPage() {
  const [stats, setStats] = useState<Stats>(defaultStats)
  const [general, setGeneral] = useState<Record<string, number>>({})

  useEffect(() => {
    api.get('/compensation/dashboard')
      .then((res) => setStats(res.data))
      .catch(() => {})

    Promise.allSettled([
      api.get('/employees', { params: { limit: '1' } }),
      api.get('/departments', { params: { limit: '1' } }),
      api.get('/positions', { params: { limit: '1' } }),
      api.get('/documents', { params: { limit: '1' } }),
      api.get('/expenses', { params: { limit: '100' } }),
    ]).then(([emp, dept, pos, docs, exp]) => {
      const total = (r: PromiseSettledResult<any>) => (r.status === 'fulfilled' ? r.value.data?.meta?.total ?? r.value.data?.data?.length ?? 0 : 0)
      setGeneral({
        employees: total(emp),
        departments: total(dept),
        positions: total(pos),
        documents: total(docs),
        expenses: exp.status === 'fulfilled' ? (Array.isArray(exp.value.data?.data) ? exp.value.data.data.length : 0) : 0,
      })
    })
  }, [])

  const renderCards = (list: typeof compensationCards, data: Record<string, unknown>) =>
    list.map((card) => {
      const Icon = card.icon
      const value = data[card.key]
      return (
        <Card key={card.key}>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-slate-500">{card.label}</p>
                <p className="text-3xl font-bold text-slate-900 mt-1">
                  {typeof value === 'number' ? value.toLocaleString() : `$${Number(value || 0).toLocaleString()}`}
                </p>
              </div>
              <div className={`p-3 rounded-xl ${card.bg}`}>
                <Icon className={`h-6 w-6 ${card.color}`} />
              </div>
            </div>
          </CardContent>
        </Card>
      )
    })

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-slate-900">Dashboard</h1>
        <p className="text-slate-500 mt-1">Resumen general de la organización</p>
      </div>

      <div className="mb-8">
        <h2 className="text-sm font-medium text-slate-500 mb-3">General</h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4">
          {renderCards(generalCards, general as unknown as Record<string, unknown>)}
        </div>
      </div>

      <div>
        <h2 className="text-sm font-medium text-slate-500 mb-3">Compensaciones</h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {renderCards(compensationCards, stats as unknown as Record<string, unknown>)}
        </div>
      </div>
    </div>
  )
}
