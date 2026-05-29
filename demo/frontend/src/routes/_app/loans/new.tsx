import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { useState, useMemo } from 'react'
import { ArrowLeft } from 'lucide-react'
import toast from 'react-hot-toast'
import { useCreateLoan } from '@/hooks/useLoans'
import { useBorrowers } from '@/hooks/useBorrowers'
import { useMotorcycles } from '@/hooks/useMotorcycles'
import { useLoanProducts } from '@/hooks/useLoanProducts'
import { useBranches } from '@/hooks/useBusiness'
import { SearchableSelect } from '@/components/ui/SearchableSelect'
import { formatUGX } from '@/lib/utils'

export const Route = createFileRoute('/_app/loans/new')({
  component: NewLoanPage,
  validateSearch: (s: Record<string, unknown>) => ({
    borrower_id: typeof s.borrower_id === 'string' ? s.borrower_id : undefined,
    motorcycle_id: typeof s.motorcycle_id === 'string' ? s.motorcycle_id : undefined,
  }),
})

function NewLoanPage() {
  const navigate = useNavigate()
  const search = Route.useSearch()
  const create = useCreateLoan()

  const { data: branches } = useBranches()
  const { data: borrowersData } = useBorrowers()
  const { data: motosData } = useMotorcycles({ status: 'available' })
  const { data: products } = useLoanProducts(true)

  const [form, setForm] = useState({
    branch_id: '',
    borrower_id: search.borrower_id || '',
    motorcycle_id: search.motorcycle_id || '',
    loan_product_id: '',
    principal_amount: '',
    initial_deposit: '0',
    duration: '',
    first_payment_date: '',
  })

  const selectedProduct = useMemo(
    () => products?.find((p) => p.id === Number(form.loan_product_id)),
    [products, form.loan_product_id],
  )

  // Estimated installment for live preview using flat math (rough — server is authoritative).
  const estimatedInstallment = useMemo(() => {
    const p = Number(form.principal_amount) - Number(form.initial_deposit || 0)
    const d = Number(form.duration)
    if (!selectedProduct || !p || !d) return null
    const cyclesPerYear = selectedProduct.repayment_cycle === 'weekly' ? 52 : selectedProduct.repayment_cycle === 'biweekly' ? 26 : 12
    const totalInterest = p * (selectedProduct.interest_rate / 100) * (d / cyclesPerYear)
    return (p + totalInterest) / d
  }, [form.principal_amount, form.initial_deposit, form.duration, selectedProduct])

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const data = await create.mutateAsync({
        branch_id: Number(form.branch_id),
        borrower_id: Number(form.borrower_id),
        motorcycle_id: form.motorcycle_id ? Number(form.motorcycle_id) : undefined,
        loan_product_id: Number(form.loan_product_id),
        principal_amount: Number(form.principal_amount),
        initial_deposit: Number(form.initial_deposit) || 0,
        duration: Number(form.duration),
        first_payment_date: form.first_payment_date,
      })
      toast.success(`Loan ${data.loan_number} created`)
      navigate({ to: `/loans/${data.id}` })
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed')
    }
  }

  return (
    <div className="p-6 lg:p-8 max-w-3xl">
      <Link to="/loans" className="inline-flex items-center gap-2 text-sm text-text-muted hover:text-text-primary mb-4">
        <ArrowLeft size={14} /> Back to loans
      </Link>

      <h1 className="text-2xl font-bold text-text-primary mb-6">New Loan</h1>

      <form onSubmit={submit} className="space-y-5 bg-surface border border-border rounded-2xl p-6 shadow-sm">
        <SearchableSelect label="Branch" required placeholder="Select branch" value={form.branch_id} onChange={(v) => setForm({ ...form, branch_id: v })}
          options={branches?.map((b) => ({ value: String(b.id), label: b.name })) || []} />

        <SearchableSelect label="Borrower" required placeholder="Search borrowers" value={form.borrower_id} onChange={(v) => setForm({ ...form, borrower_id: v })}
          options={borrowersData?.data.map((b) => ({ value: String(b.id), label: `${b.full_name} · ${b.phone}` })) || []} />

        <SearchableSelect label="Motorcycle (optional)" placeholder="Search available motorcycles" value={form.motorcycle_id} onChange={(v) => setForm({ ...form, motorcycle_id: v })}
          options={motosData?.data.map((m) => ({ value: String(m.id), label: `${m.name} · ${m.number_plate} · ${formatUGX(m.selling_price)}` })) || []} />

        <SearchableSelect label="Loan Product" required placeholder="Select loan product" value={form.loan_product_id} onChange={(v) => setForm({ ...form, loan_product_id: v })}
          options={products?.map((p) => ({ value: String(p.id), label: `${p.name} · ${p.interest_rate}% · ${p.interest_method.replace('_', ' ')}` })) || []} />

        {selectedProduct && (
          <div className="text-xs text-text-muted bg-background rounded-lg p-3">
            Range: {formatUGX(selectedProduct.min_amount)} – {formatUGX(selectedProduct.max_amount)} ·
            Duration: {selectedProduct.min_duration}–{selectedProduct.max_duration} {selectedProduct.repayment_cycle === 'monthly' ? 'months' : selectedProduct.repayment_cycle === 'biweekly' ? 'biweeks' : 'weeks'} ·
            Interest: {selectedProduct.interest_rate}% / yr ({selectedProduct.interest_method.replace('_', ' ')})
          </div>
        )}

        <div className="grid grid-cols-2 gap-3">
          <Field label="Principal Amount (UGX)" type="number" required value={form.principal_amount} onChange={(v) => setForm({ ...form, principal_amount: v })} />
          <Field label="Initial Deposit (UGX)" type="number" value={form.initial_deposit} onChange={(v) => setForm({ ...form, initial_deposit: v })} />
        </div>

        <div className="grid grid-cols-2 gap-3">
          <Field label={`Duration (${selectedProduct?.repayment_cycle === 'monthly' ? 'months' : selectedProduct?.repayment_cycle === 'biweekly' ? 'biweeks' : 'weeks'})`} type="number" required value={form.duration} onChange={(v) => setForm({ ...form, duration: v })} />
          <Field label="First Payment Date" type="date" required value={form.first_payment_date} onChange={(v) => setForm({ ...form, first_payment_date: v })} />
        </div>

        {estimatedInstallment !== null && (
          <div className="bg-grit-blue-50 border border-grit-blue/30 rounded-lg p-4">
            <p className="text-xs text-text-muted">Estimated installment (per {selectedProduct?.repayment_cycle})</p>
            <p className="text-2xl font-bold text-grit-blue">{formatUGX(estimatedInstallment)}</p>
            <p className="text-xs text-text-muted mt-1">Server will compute the exact schedule on submit.</p>
          </div>
        )}

        <button type="submit" disabled={create.isPending} className="w-full py-3 rounded-xl bg-grit-blue text-white font-semibold hover:bg-grit-blue-dark disabled:opacity-50">
          {create.isPending ? 'Creating...' : 'Create Loan'}
        </button>
      </form>
    </div>
  )
}

function Field({ label, value, onChange, required, type = 'text' }: { label: string; value: string; onChange: (v: string) => void; required?: boolean; type?: string }) {
  return (
    <div>
      <label className="block text-sm font-semibold text-text-primary mb-1.5">{label}{required && <span className="text-error"> *</span>}</label>
      <input type={type} required={required} value={value} onChange={(e) => onChange(e.target.value)}
        className="w-full px-4 py-2.5 rounded-xl border border-border bg-white text-sm focus:outline-none focus:ring-2 focus:ring-grit-blue/20 focus:border-grit-blue" />
    </div>
  )
}
