import { useEffect, useState } from 'react'
import api from '@/lib/api'
import { Card, CardContent } from '@/components/ui/card'
import { Users, Banknote, PiggyBank, Gift, Wallet, BarChart3 } from 'lucide-react'

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

const cards = [
  { key: 'total_salary_cost', label: 'Costo Salarial Total', icon: Users, color: 'text-brand-600', bg: 'bg-brand-50' },
  { key: 'average_compensation', label: 'Compensación Promedio', icon: Banknote, color: 'text-emerald-600', bg: 'bg-emerald-50' },
  { key: 'benefit_cost', label: 'Costo en Beneficios', icon: Gift, color: 'text-amber-600', bg: 'bg-amber-50' },
  { key: 'total_bonuses', label: 'Bonos Totales', icon: PiggyBank, color: 'text-violet-600', bg: 'bg-violet-50' },
  { key: 'budget_total', label: 'Presupuesto Total', icon: Wallet, color: 'text-cyan-600', bg: 'bg-cyan-50' },
  { key: 'budget_used', label: 'Presupuesto Utilizado', icon: BarChart3, color: 'text-rose-600', bg: 'bg-rose-50' },
]

export default function DashboardPage() {
  const [stats, setStats] = useState<Stats>(defaultStats)

  useEffect(() => {
    api.get('/compensation/dashboard')
      .then((res) => setStats(res.data))
      .catch(() => {})
  }, [])

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-slate-900">Dashboard</h1>
        <p className="text-slate-500 mt-1">Resumen general de compensaciones</p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {cards.map((card) => {
          const Icon = card.icon
          const value = (stats as unknown as Record<string, string | number>)[card.key]
          return (
            <Card key={card.key}>
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-slate-500">{card.label}</p>
                    <p className="text-3xl font-bold text-slate-900 mt-1">
                      {typeof value === 'number' ? value : `$${value}`}
                    </p>
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
