import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { useState, useMemo } from 'react'
import { ArrowLeft } from 'lucide-react'
import toast from 'react-hot-toast'
import { useCreateCashSale } from '@/hooks/useCashSales'
import { useMotorcycles } from '@/hooks/useMotorcycles'
import { useBranches } from '@/hooks/useBusiness'
import { SearchableSelect } from '@/components/ui/SearchableSelect'
import { formatUGX } from '@/lib/utils'

export const Route = createFileRoute('/_app/cash-sales/new')({
  component: NewCashSalePage,
  validateSearch: (s: Record<string, unknown>) => ({
    motorcycle_id: typeof s.motorcycle_id === 'string' ? s.motorcycle_id : undefined,
  }),
})

function NewCashSalePage() {
  const navigate = useNavigate()
  const search = Route.useSearch()
  const create = useCreateCashSale()

  const { data: branches } = useBranches()
  const { data: motosData } = useMotorcycles({ status: 'available' })

  const [form, setForm] = useState({
    branch_id: '',
    motorcycle_id: search.motorcycle_id || '',
    customer_name: '', customer_phone: '', customer_email: '', customer_nin: '', customer_address: '',
    list_price: '', discount_amount: '0', total: '',
    payment_method: 'cash', transaction_ref: '', notes: '',
  })

  const selectedMoto = useMemo(
    () => motosData?.data.find((m) => m.id === Number(form.motorcycle_id)),
    [motosData, form.motorcycle_id],
  )

  const onPickMoto = (id: string) => {
    const moto = motosData?.data.find((m) => m.id === Number(id))
    setForm({
      ...form,
      motorcycle_id: id,
      list_price: moto ? String(moto.selling_price) : form.list_price,
      total: moto ? String(moto.selling_price) : form.total,
      branch_id: moto ? String(moto.branch_id) : form.branch_id,
    })
  }

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const data = await create.mutateAsync({
        branch_id: Number(form.branch_id),
        motorcycle_id: Number(form.motorcycle_id),
        customer_name: form.customer_name,
        customer_phone: form.customer_phone,
        customer_email: form.customer_email,
        customer_nin: form.customer_nin,
        customer_address: form.customer_address,
        list_price: Number(form.list_price) || 0,
        discount_amount: Number(form.discount_amount) || 0,
        total: Number(form.total),
        payment_method: form.payment_method,
        transaction_ref: form.transaction_ref,
        notes: form.notes,
      })
      toast.success(`Cash sale ${data.sale_number} recorded`)
      navigate({ to: '/cash-sales' })
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed')
    }
  }

  return (
    <div className="p-6 lg:p-8 max-w-3xl">
      <Link to="/cash-sales" className="inline-flex items-center gap-2 text-sm text-text-muted hover:text-text-primary mb-4">
        <ArrowLeft size={14} /> Back
      </Link>
      <h1 className="text-2xl font-bold text-text-primary mb-6">New Cash Sale</h1>

      <form onSubmit={submit} className="space-y-5 bg-surface border border-border rounded-2xl p-6 shadow-sm">
        <SearchableSelect label="Motorcycle" required placeholder="Search available motorcycles" value={form.motorcycle_id} onChange={onPickMoto}
          options={motosData?.data.map((m) => ({ value: String(m.id), label: `${m.name} · ${m.number_plate} · ${formatUGX(m.selling_price)}` })) || []} />

        {selectedMoto && (
          <div className="bg-grit-blue-50 border border-grit-blue/30 rounded-lg p-3 text-sm">
            Selected: <span className="font-semibold">{selectedMoto.name} {selectedMoto.number_plate}</span> · List {formatUGX(selectedMoto.selling_price)}
          </div>
        )}

        <SearchableSelect label="Branch" required placeholder="Select branch" value={form.branch_id} onChange={(v) => setForm({ ...form, branch_id: v })}
          options={branches?.map((b) => ({ value: String(b.id), label: b.name })) || []} />

        <div className="border-t border-border pt-4">
          <p className="text-xs font-semibold text-text-muted uppercase tracking-wider mb-3">Customer</p>
          <div className="grid grid-cols-2 gap-3">
            <Field label="Name" required value={form.customer_name} onChange={(v) => setForm({ ...form, customer_name: v })} />
            <Field label="Phone" value={form.customer_phone} onChange={(v) => setForm({ ...form, customer_phone: v })} />
          </div>
          <div className="grid grid-cols-2 gap-3 mt-3">
            <Field label="Email" type="email" value={form.customer_email} onChange={(v) => setForm({ ...form, customer_email: v })} />
            <Field label="NIN" value={form.customer_nin} onChange={(v) => setForm({ ...form, customer_nin: v })} />
          </div>
          <div className="mt-3"><Field label="Address" value={form.customer_address} onChange={(v) => setForm({ ...form, customer_address: v })} /></div>
        </div>

        <div className="border-t border-border pt-4">
          <p className="text-xs font-semibold text-text-muted uppercase tracking-wider mb-3">Payment</p>
          <div className="grid grid-cols-3 gap-3">
            <Field label="List Price" type="number" value={form.list_price} onChange={(v) => setForm({ ...form, list_price: v })} />
            <Field label="Discount" type="number" value={form.discount_amount} onChange={(v) => setForm({ ...form, discount_amount: v })} />
            <Field label="Total" type="number" required value={form.total} onChange={(v) => setForm({ ...form, total: v })} />
          </div>
          <div className="grid grid-cols-2 gap-3 mt-3">
            <div>
              <label className="block text-sm font-semibold text-text-primary mb-1.5">Method <span className="text-error">*</span></label>
              <select required value={form.payment_method} onChange={(e) => setForm({ ...form, payment_method: e.target.value })}
                className="w-full px-4 py-2.5 rounded-xl border border-border bg-white text-sm focus:outline-none focus:ring-2 focus:ring-grit-blue/20">
                <option value="cash">Cash</option>
                <option value="mobile_money">Mobile Money</option>
                <option value="bank_transfer">Bank Transfer</option>
              </select>
            </div>
            <Field label="Transaction Ref" value={form.transaction_ref} onChange={(v) => setForm({ ...form, transaction_ref: v })} />
          </div>
        </div>

        <button type="submit" disabled={create.isPending} className="w-full py-3 rounded-xl bg-grit-blue text-white font-semibold hover:bg-grit-blue-dark disabled:opacity-50">
          {create.isPending ? 'Recording...' : 'Record Cash Sale'}
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
