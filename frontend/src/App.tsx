import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from '@/contexts/auth-context'
import LoginPage from '@/pages/login'
import DashboardPage from '@/pages/dashboard'
import EmployeesPage from '@/pages/employees'
import DepartmentsPage from '@/pages/departments'
import PositionsPage from '@/pages/positions'
import RecruitmentPage from '@/pages/recruitment'
import TrainingPage from '@/pages/training'
import CompensationPage from '@/pages/compensation'
import ExpensesPage from '@/pages/expenses'
import DocumentsPage from '@/pages/documents'
import LeavesPage from '@/pages/leaves'
import EmployeeDetailPage from '@/pages/employee-detail'
import ProfilePage from '@/pages/profile'
import OrgChartPage from '@/pages/org-chart'
import FeedPage from '@/pages/feed'
import SurveysPage from '@/pages/surveys'
import SettingsPage from '@/pages/settings'
import AttendancePage from '@/pages/attendance'
import SchedulingPage from '@/pages/scheduling'
import OvertimePage from '@/pages/overtime'
import PayrollPage from '@/pages/payroll'
import BenefitsPage from '@/pages/benefits'
import PerformancePage from '@/pages/performance'
import OnboardingPage from '@/pages/onboarding'
import AppLayout from '@/layouts/app-layout'

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth()
  if (loading) return <div className="min-h-screen flex items-center justify-center text-slate-500">Cargando...</div>
  if (!user) return <Navigate to="/login" replace />
  return <>{children}</>
}

function PublicRoute({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth()
  if (loading) return <div className="min-h-screen flex items-center justify-center text-slate-500">Cargando...</div>
  if (user) return <Navigate to="/" replace />
  return <>{children}</>
}

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route
            path="/login"
            element={
              <PublicRoute>
                <LoginPage />
              </PublicRoute>
            }
          />
          <Route
            element={
              <ProtectedRoute>
                <AppLayout />
              </ProtectedRoute>
            }
          >
            <Route index element={<DashboardPage />} />
            <Route path="employees" element={<EmployeesPage />} />
            <Route path="employees/:id" element={<EmployeeDetailPage />} />
            <Route path="leaves" element={<LeavesPage />} />
            <Route path="departments" element={<DepartmentsPage />} />
            <Route path="positions" element={<PositionsPage />} />
            <Route path="recruitment" element={<RecruitmentPage />} />
            <Route path="training" element={<TrainingPage />} />
            <Route path="compensation" element={<CompensationPage />} />
            <Route path="expenses" element={<ExpensesPage />} />
            <Route path="documents" element={<DocumentsPage />} />
            <Route path="profile" element={<ProfilePage />} />
            <Route path="organization" element={<OrgChartPage />} />
            <Route path="feed" element={<FeedPage />} />
            <Route path="surveys" element={<SurveysPage />} />
            <Route path="settings" element={<SettingsPage />} />
            <Route path="attendance" element={<AttendancePage />} />
            <Route path="scheduling" element={<SchedulingPage />} />
            <Route path="overtime" element={<OvertimePage />} />
            <Route path="payroll" element={<PayrollPage />} />
            <Route path="benefits" element={<BenefitsPage />} />
            <Route path="performance" element={<PerformancePage />} />
            <Route path="onboarding" element={<OnboardingPage />} />
          </Route>
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  )
}

export default App
