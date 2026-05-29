// === User & Auth ===
export interface User {
  id: number
  name: string
  email: string
}

export interface TokenPair {
  access_token: string
  refresh_token: string
  expires_at: number
}

export type WorkspaceAccess = 'loans' | 'spares' | 'both'

export interface BusinessInfo {
  id: number
  name: string
  role: string
  branch_id: number | null
  /** Controls which sidebar workspaces (loans / spares) this user sees. */
  workspace_access: WorkspaceAccess
}

export interface RoleInfo {
  business_id: number
  business_name: string
  role: string
  branch_id: number | null
  branch_name?: string
  workspace_access: WorkspaceAccess
}

// === Business & Branch ===
export interface Business {
  id: number
  name: string
  owner_id: number
  created_at: string
}

export interface Branch {
  id: number
  business_id: number
  name: string
  address: string
  is_default: boolean
  created_at: string
}

// === Products ===
export interface Category {
  id: number
  business_id: number
  name: string
}

export interface Product {
  id: number
  business_id: number
  category_id: number
  category?: Category
  title: string
  description: string
  barcode: string
  selling_price: number
  cost_price: number
  image_url?: string
  low_stock_threshold: number
  created_at: string
}

// === Stock ===
export interface Stock {
  id: number
  product_id: number
  branch_id: number
  quantity: number
  product?: Product
  branch?: Branch
}

export interface StockMovement {
  id: number
  product_id: number
  branch_id: number
  movement_type: 'stock_in' | 'sale' | 'transfer_out' | 'transfer_in' | 'adjustment'
  quantity: number
  reference_id: number | null
  note: string
  created_by: number
  created_at: string
  product?: Product
  branch?: Branch
}

// === Sales ===
export interface SaleItem {
  id: number
  sale_id: number
  product_id: number
  quantity: number
  unit_price: number
  unit_cost: number
  line_total: number
  product?: Product
}

export interface Sale {
  id: number
  branch_id: number
  cashier_id: number
  payment_method: 'cash' | 'mobile_money'
  customer_phone: string
  subtotal: number
  discount_amount: number
  total: number
  items?: SaleItem[]
  branch?: Branch
  cashier?: User
  created_at: string
}

// === Roles ===
export type Role = 'admin' | 'manager' | 'cashier' | 'stock_clerk' | 'loan_officer' | 'accountant'

// === Motorcycles ===
export type MotorcycleStatus = 'available' | 'reserved' | 'sold' | 'on_loan' | 'repossessed'

export interface Motorcycle {
  id: number
  business_id: number
  branch_id: number
  branch?: { id: number; name: string }
  name: string
  number_plate: string
  chassis_no: string
  engine_no: string
  color: string
  year_of_make: number
  cost_price: number
  selling_price: number
  status: MotorcycleStatus
  notes: string
  image_url?: string
  created_at: string
  updated_at: string
}

// === Loan Products ===
export type InterestMethod = 'flat' | 'reducing_balance'
export type RepaymentCycle = 'weekly' | 'biweekly' | 'monthly'

export interface LoanProduct {
  id: number
  business_id: number
  name: string
  description: string
  min_amount: number
  max_amount: number
  min_duration: number
  max_duration: number
  interest_method: InterestMethod
  interest_rate: number
  repayment_cycle: RepaymentCycle
  requires_collateral: boolean
  grace_period_days: number
  is_active: boolean
  created_at: string
}

// === Borrowers ===
export type RiskLevel = 'low' | 'medium' | 'high' | 'critical'

export interface Borrower {
  id: number
  business_id: number
  branch_id: number
  branch?: { id: number; name: string }
  first_name: string
  last_name: string
  full_name: string
  phone: string
  alt_phone: string
  email: string
  date_of_birth?: string
  gender: string
  national_id: string
  address: string
  employment_status: string
  occupation: string
  employer: string
  monthly_income: number
  next_of_kin_name: string
  next_of_kin_phone: string
  next_of_kin_relation: string
  credit_score: number
  risk_level: RiskLevel
  photo_url?: string
  id_document_url?: string
  created_at: string
}

// === Loans ===
export type LoanStatus = 'pending' | 'approved' | 'active' | 'completed' | 'defaulted' | 'rejected' | 'cancelled'

export interface Loan {
  id: number
  business_id: number
  branch_id: number
  borrower_id: number
  motorcycle_id: number | null
  loan_product_id: number
  loan_number: string
  principal_amount: number
  initial_deposit: number
  disbursed_amount: number
  duration: number
  interest_method: InterestMethod
  interest_rate: number
  repayment_cycle: RepaymentCycle
  total_interest: number
  total_repayments: number
  installment_amount: number
  total_amount: number
  balance_remaining: number
  status: LoanStatus
  disbursement_date: string | null
  first_payment_date: string | null
  next_payment_date: string | null
  maturity_date: string | null
  completed_at: string | null
  created_at: string
  borrower?: Borrower
  motorcycle?: Motorcycle
  loan_product?: LoanProduct
  branch?: { id: number; name: string }
}

export interface RepaymentSchedule {
  id: number
  loan_id: number
  installment_number: number
  due_date: string
  principal_amount: number
  interest_amount: number
  total_amount: number
  balance_before: number
  balance_after: number
  is_paid: boolean
  paid_amount: number
  paid_at: string | null
  is_overdue: boolean
  days_past_due: number
}

// === Repayments ===
export type RepaymentStatus = 'pending' | 'approved' | 'rejected' | 'failed'
export type RepaymentMethod = 'cash' | 'mobile_money' | 'bank_transfer' | 'cheque'

export interface Repayment {
  id: number
  loan_id: number
  schedule_id: number | null
  amount: number
  collection_date: string
  payment_method: RepaymentMethod
  transaction_ref: string
  receipt: string
  status: RepaymentStatus
  notes: string
  dgateway_reference: string
  dgateway_provider: string
  dgateway_phone?: string
  dgateway_fee: number
  dgateway_net_amount: number
  dgateway_fail_reason: string
  dgateway_confirmed_at: string | null
  collected_by: number
  verified_by: number | null
  verified_at: string | null
  created_at: string
  loan?: Loan
  collector?: User
}

// === Cash Sales ===
export interface CashSale {
  id: number
  business_id: number
  branch_id: number
  motorcycle_id: number
  sale_number: string
  sold_by: number
  customer_name: string
  customer_phone: string
  customer_email: string
  customer_nin: string
  list_price: number
  discount_amount: number
  total: number
  payment_method: string
  transaction_ref: string
  notes: string
  created_at: string
  motorcycle?: Motorcycle
  branch?: { id: number; name: string }
  seller?: User
}

// === Daily Boda ===
export type DailyBodaMotoStatus = 'available' | 'occupied' | 'returned' | 'in_service'
export type DailyBodaPaymentStatus = 'paid' | 'partial' | 'pending'

export interface DailyBodaDriver {
  id: number
  business_id: number
  branch_id: number
  full_name: string
  phone: string
  national_id: string
  address: string
  daily_rate: number
  is_active: boolean
  branch?: { id: number; name: string }
}

export interface DailyBodaMotorcycle {
  id: number
  business_id: number
  branch_id: number
  name: string
  number_plate: string
  status: DailyBodaMotoStatus
  assigned_driver_id: number | null
  branch?: { id: number; name: string }
}

export interface DailyBodaPayment {
  id: number
  driver_id: number
  motorcycle_id: number
  branch_id: number
  amount: number
  daily_rate: number
  balance: number
  payment_date: string
  payment_method: string
  status: DailyBodaPaymentStatus
  verified_at: string | null
  notes: string
  driver?: DailyBodaDriver
  motorcycle?: DailyBodaMotorcycle
}
