import { useState } from 'react'
import { Outlet, Link, useLocation } from 'react-router-dom'
import { useAuth } from '@/contexts/auth-context'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Separator } from '@/components/ui/separator'
import { cn } from '@/lib/utils'
import {
  LayoutDashboard,
  Users,
  Building2,
  Briefcase,
  GraduationCap,
  DollarSign,
  Receipt,
  Handshake,
  FileText,
  CalendarCheck,
  Clock,
  CalendarDays,
  Timer,
  Banknote,
  Gift,
  Star,
  Rocket,
  Network,
  MessageSquare,
  ClipboardList,
  UserCircle,
  Settings,
  ChevronLeft,
  ChevronRight,
  LogOut,
  Menu,
  X,
} from 'lucide-react'

const navItems = [
  { href: '/', label: 'Dashboard', icon: LayoutDashboard },
  { href: '/feed', label: 'Feed', icon: MessageSquare },
  { href: '/employees', label: 'Empleados', icon: Users },
  { href: '/organization', label: 'Organigrama', icon: Network },
  { href: '/departments', label: 'Departamentos', icon: Building2 },
  { href: '/positions', label: 'Posiciones', icon: Briefcase },
  { href: '/recruitment', label: 'Reclutamiento', icon: Handshake },
  { href: '/training', label: 'Capacitación', icon: GraduationCap },
  { href: '/compensation', label: 'Compensaciones', icon: DollarSign },
  { href: '/expenses', label: 'Gastos', icon: Receipt },
  { href: '/documents', label: 'Documentos', icon: FileText },
  { href: '/leaves', label: 'Licencias', icon: CalendarCheck },
  { href: '/attendance', label: 'Asistencia', icon: Clock },
  { href: '/scheduling', label: 'Turnos', icon: CalendarDays },
  { href: '/overtime', label: 'Horas Extras', icon: Timer },
  { href: '/payroll', label: 'Nómina', icon: Banknote },
  { href: '/benefits', label: 'Beneficios', icon: Gift },
  { href: '/performance', label: 'Desempeño', icon: Star },
  { href: '/onboarding', label: 'Onboarding', icon: Rocket },
  { href: '/surveys', label: 'Encuestas', icon: ClipboardList },
  { href: '/profile', label: 'Mi Perfil', icon: UserCircle },
  { href: '/settings', label: 'Configuración', icon: Settings },
]

export default function AppLayout() {
  const [collapsed, setCollapsed] = useState(false)
  const [mobileOpen, setMobileOpen] = useState(false)
  const { user, logout } = useAuth()
  const location = useLocation()

  const initials = user?.name
    ? user.name.split(' ').map((n) => n[0]).join('').toUpperCase().slice(0, 2)
    : 'U'

  return (
    <div className="min-h-screen bg-surface flex">
      {/* Mobile backdrop */}
      {mobileOpen && (
        <div
          className="fixed inset-0 bg-black/50 z-40 lg:hidden"
          onClick={() => setMobileOpen(false)}
        />
      )}

      {/* Sidebar */}
      <aside
        className={cn(
          'fixed lg:static inset-y-0 left-0 z-50 flex flex-col bg-sidebar text-white transition-all duration-300',
          collapsed ? 'w-16' : 'w-60',
          mobileOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0',
        )}
      >
        <div className="flex items-center h-16 px-4 gap-3 border-b border-sidebar-hover">
          {!collapsed && (
            <span className="text-lg font-bold tracking-tight">RRHHumand</span>
          )}
          <button
            onClick={() => { setCollapsed(!collapsed); setMobileOpen(false) }}
            className="ml-auto p-1.5 rounded-lg hover:bg-sidebar-hover transition-colors hidden lg:block"
          >
            {collapsed ? <ChevronRight size={18} /> : <ChevronLeft size={18} />}
          </button>
          <button
            onClick={() => setMobileOpen(false)}
            className="ml-auto p-1.5 rounded-lg hover:bg-sidebar-hover transition-colors lg:hidden"
          >
            <X size={18} />
          </button>
        </div>

        <nav className="flex-1 overflow-y-auto p-2 space-y-1">
          {navItems.map((item) => {
            const Icon = item.icon
            const active = location.pathname === item.href ||
              (item.href !== '/' && location.pathname.startsWith(item.href))
            return (
              <Link
                key={item.href}
                to={item.href}
                onClick={() => setMobileOpen(false)}
                className={cn(
                  'flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors',
                  active
                    ? 'bg-sidebar-active text-white'
                    : 'text-slate-400 hover:bg-sidebar-hover hover:text-white',
                )}
              >
                <Icon size={20} className="shrink-0" />
                {!collapsed && <span>{item.label}</span>}
              </Link>
            )
          })}
        </nav>

        <Separator className="bg-sidebar-hover" />
        <div className="p-3">
          <Link to="/profile" onClick={() => setMobileOpen(false)} className={cn('flex items-center gap-3 rounded-lg hover:bg-sidebar-hover transition-colors', collapsed && 'justify-center')}>
            <Avatar className="h-8 w-8">
              <AvatarFallback className="bg-brand-600 text-white text-xs">{initials}</AvatarFallback>
            </Avatar>
            {!collapsed && (
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-white truncate">{user?.name}</p>
                <p className="text-xs text-slate-400 truncate">{user?.email}</p>
              </div>
            )}
          </Link>
        </div>
      </aside>

      {/* Main content */}
      <div className="flex-1 flex flex-col min-w-0">
        <header className="h-16 border-b border-slate-200 bg-white flex items-center px-4 lg:px-6 gap-4">
          <button
            onClick={() => setMobileOpen(true)}
            className="p-2 rounded-lg hover:bg-slate-100 transition-colors lg:hidden"
          >
            <Menu size={20} />
          </button>

          <div className="flex-1" />

          <Button variant="ghost" size="sm" onClick={logout} className="text-slate-600">
            <LogOut size={16} />
            <span className="hidden sm:inline ml-1">Salir</span>
          </Button>
        </header>

        <main className="flex-1 p-4 lg:p-6 overflow-auto">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
