export interface User {
  id: number
  email: string
  first_name: string
  last_name: string
  phone: string
  role_id: number
  role?: { id: number; name: string }
  account_type: 'individual' | 'company'
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface Company {
  id: number
  user_id: number
  company_name: string
  legal_name: string
  bin: string
  contact_person: string
  created_at: string
  updated_at: string
}

export interface AuthResponse {
  access_token: string
  refresh_token: string
  expires_at: string
  user: User
}

export interface Subscription {
  id: number
  user_id: number
  user?: User
  status: string
  start_date: string | null
  end_date: string | null
  amount: number
  admin_comment: string
  requested_at: string
  approved_at: string | null
  created_at: string
  // Company fields (present when account_type === 'company')
  company_name?: string
  legal_name?: string
  bin?: string
}

export interface Document {
  id: number
  user_id: number
  user?: User
  original_name: string
  status: string
  error_message: string
  file_size: number
  queued_at: string | null
  processed_at: string | null
  created_at: string
}

export interface AIData {
  document_type: string
  fields: Record<string, string | number>
  items: Record<string, string | number>[]
}

export interface AdminStats {
  total_companies: number
  active_subscriptions: number
  pending_subscriptions: number
  total_documents: number
  completed_documents: number
}

export interface RegisterRequest {
  email: string
  password: string
  first_name: string
  last_name: string
  phone: string
  account_type: 'individual' | 'company'
  company_name?: string
  legal_name?: string
  bin?: string
  contact_person?: string
}

export interface LoginRequest {
  email: string
  password: string
}
