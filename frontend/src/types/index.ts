export interface User {
  id: string
  email: string
  name: string
  role: string
  company_id: string
}

export interface AuthResponse {
  success: boolean
  data: {
    token: string
    refresh_token: string
    user: User
  }
}

export interface ApiResponse<T = unknown> {
  success: boolean
  data: T
  error?: string
}

export interface DashboardStats {
  employee_count: number
  active_courses: number
  pending_expenses: number
  open_recruitments: number
  payroll_runs: number
  pending_approvals: number
}
