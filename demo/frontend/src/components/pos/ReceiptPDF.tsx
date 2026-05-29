import { Document, Page, Text, View, StyleSheet, Font } from '@react-pdf/renderer'

/**
 * Onest webfont — registered once at module load. We host nothing; we pull
 * directly from Google Fonts CDN. If the CDN is unreachable, react-pdf
 * falls back to its built-in Helvetica which is also fine.
 */
Font.register({
  family: 'Onest',
  fonts: [
    { src: 'https://fonts.gstatic.com/s/onest/v7/gNMKW3F-SZkJdHse31CzagS_K20.woff2', fontWeight: 400 },
    { src: 'https://fonts.gstatic.com/s/onest/v7/gNMKW3F-SZkJdHse31CzagS_Dmo.woff2', fontWeight: 600 },
    { src: 'https://fonts.gstatic.com/s/onest/v7/gNMKW3F-SZkJdHse31CzagS_OGo.woff2', fontWeight: 700 },
  ],
})

const C = {
  ink:         '#0F172A',
  inkSoft:     '#475569',
  muted:       '#94A3B8',
  line:        '#E2E8F0',
  lineSoft:    '#F1F5F9',
  paper:       '#FFFFFF',
  panel:       '#F8FAFC',
  brand:       '#1E7EF5',
  brandLight:  '#E8F1FE',
  brandDark:   '#1565C0',
  good:        '#059669',
  goodLight:   '#D1FAE5',
  warn:        '#D97706',
  bad:         '#DC2626',
  badLight:    '#FEE2E2',
}

const s = StyleSheet.create({
  page: {
    paddingHorizontal: 28,
    paddingVertical: 26,
    fontFamily: 'Onest',
    fontSize: 10,
    color: C.ink,
    backgroundColor: C.paper,
  },

  /* Brand bar */
  brandBar: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingBottom: 14,
    borderBottomWidth: 2,
    borderBottomColor: C.ink,
  },
  brandLeft: { flexDirection: 'row', alignItems: 'center', gap: 10 },
  logo: {
    width: 34, height: 34, borderRadius: 8,
    backgroundColor: C.brand,
    alignItems: 'center', justifyContent: 'center',
  },
  logoText: { color: '#FFFFFF', fontSize: 14, fontWeight: 700, letterSpacing: 0.4 },
  brandName: { fontSize: 15, fontWeight: 700, color: C.ink, letterSpacing: 0.2 },
  brandSub: { fontSize: 8, fontWeight: 600, color: C.muted, marginTop: 1, textTransform: 'uppercase', letterSpacing: 0.8 },
  brandRight: { alignItems: 'flex-end' },
  receiptLabel: { fontSize: 8, fontWeight: 700, color: C.muted, textTransform: 'uppercase', letterSpacing: 1.2 },
  receiptNo: { fontSize: 14, fontWeight: 700, color: C.ink, marginTop: 2 },
  receiptDate: { fontSize: 9, color: C.inkSoft, marginTop: 2 },

  /* Meta block */
  metaCard: {
    flexDirection: 'row',
    backgroundColor: C.panel,
    borderRadius: 6,
    padding: 12,
    marginTop: 16,
    gap: 16,
  },
  metaCol: { flex: 1, gap: 6 },
  metaLabel: { fontSize: 7.5, fontWeight: 600, color: C.muted, textTransform: 'uppercase', letterSpacing: 0.6 },
  metaValue: { fontSize: 10.5, color: C.ink, fontWeight: 600, marginTop: 1 },

  /* Items table */
  sectionTitle: {
    fontSize: 8,
    fontWeight: 700,
    color: C.muted,
    textTransform: 'uppercase',
    letterSpacing: 1,
    marginTop: 18,
    marginBottom: 8,
  },
  table: {
    borderTopWidth: 1,
    borderTopColor: C.ink,
  },
  thead: {
    flexDirection: 'row',
    paddingVertical: 7,
    borderBottomWidth: 0.5,
    borderBottomColor: C.line,
  },
  th: {
    fontSize: 8,
    fontWeight: 700,
    color: C.inkSoft,
    textTransform: 'uppercase',
    letterSpacing: 0.6,
  },
  trow: {
    flexDirection: 'row',
    paddingVertical: 8,
    borderBottomWidth: 0.5,
    borderBottomColor: C.lineSoft,
    alignItems: 'flex-start',
  },
  cItem:  { flex: 5, paddingRight: 6 },
  cQty:   { flex: 1, textAlign: 'center' },
  cPrice: { flex: 2, textAlign: 'right' },
  cTotal: { flex: 2, textAlign: 'right' },

  itemName: { fontSize: 10, fontWeight: 600, color: C.ink },
  itemSku:  { fontSize: 8, color: C.muted, marginTop: 1 },
  cellQty:  { fontSize: 10, color: C.inkSoft, textAlign: 'center' },
  cellMoney:{ fontSize: 10, color: C.ink, textAlign: 'right' },
  cellTotal:{ fontSize: 10, color: C.ink, textAlign: 'right', fontWeight: 600 },

  /* Totals card */
  totalsWrap: { flexDirection: 'row', justifyContent: 'flex-end', marginTop: 14 },
  totalsCard: {
    width: '60%',
    borderRadius: 6,
    overflow: 'hidden',
    borderWidth: 1,
    borderColor: C.line,
  },
  totalsRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    paddingHorizontal: 12,
    paddingVertical: 7,
    borderBottomWidth: 0.5,
    borderBottomColor: C.line,
  },
  totalsLabel: { fontSize: 10, color: C.inkSoft },
  totalsValue: { fontSize: 10, color: C.ink, fontWeight: 600 },
  grandRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    paddingHorizontal: 12,
    paddingVertical: 11,
    backgroundColor: C.ink,
  },
  grandLabel: { fontSize: 11, fontWeight: 700, color: '#FFFFFF', letterSpacing: 0.6 },
  grandValue: { fontSize: 14, fontWeight: 700, color: '#FFFFFF' },

  /* Payment chip */
  payChip: {
    alignSelf: 'flex-start',
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
    paddingHorizontal: 8,
    paddingVertical: 3,
    borderRadius: 999,
    fontSize: 9,
    fontWeight: 600,
  },
  payChipCash:  { backgroundColor: C.goodLight,  color: C.good },
  payChipMomo:  { backgroundColor: C.brandLight, color: C.brandDark },
  payChipCredit:{ backgroundColor: C.badLight,   color: C.bad  },

  /* Footer */
  footer: {
    marginTop: 28,
    paddingTop: 14,
    borderTopWidth: 0.5,
    borderTopColor: C.line,
    alignItems: 'center',
  },
  thanks: { fontSize: 11, fontWeight: 700, color: C.ink },
  footerLine: { fontSize: 8.5, color: C.muted, marginTop: 4, textAlign: 'center' },
  footerStrong: { fontSize: 8.5, color: C.inkSoft, marginTop: 6 },
})

function fmtUGX(n: number) {
  return `UGX ${(n || 0).toLocaleString('en-US')}`
}

function methodLabel(m?: string) {
  if (!m) return '—'
  if (m === 'cash') return 'Cash'
  if (m === 'mobile_money') return 'Mobile Money'
  if (m === 'credit') return 'Credit'
  return m.replace('_', ' ').replace(/\b\w/g, (c) => c.toUpperCase())
}

interface ReceiptPDFProps {
  sale: any
  businessName?: string
  tagline?: string
}

export function ReceiptPDF({ sale, businessName = 'Grit Motors', tagline = 'Spares · Motorcycles · Loans' }: ReceiptPDFProps) {
  const date = new Date(sale.created_at)
  const items: any[] = sale.items || []
  const method = sale.payment_method as string | undefined
  const chipStyle =
    method === 'cash' ? s.payChipCash :
    method === 'mobile_money' ? s.payChipMomo :
    s.payChipCredit

  return (
    <Document
      title={`Receipt #${sale.id}`}
      author={businessName}
      subject={`Sales receipt #${sale.id}`}
      creator={businessName}
      producer={businessName}
    >
      {/* A5 portrait — 148 × 210 mm. Looks crisp printed on standard paper
          and reads well as a download. We render `auto` for height so long
          item lists flow onto a second page. */}
      <Page size="A5" style={s.page}>
        {/* Brand bar */}
        <View style={s.brandBar}>
          <View style={s.brandLeft}>
            <View style={s.logo}>
              <Text style={s.logoText}>KM</Text>
            </View>
            <View>
              <Text style={s.brandName}>{businessName}</Text>
              <Text style={s.brandSub}>{tagline}</Text>
            </View>
          </View>
          <View style={s.brandRight}>
            <Text style={s.receiptLabel}>Receipt</Text>
            <Text style={s.receiptNo}>#{String(sale.id).padStart(6, '0')}</Text>
            <Text style={s.receiptDate}>
              {date.toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })}
              {' · '}
              {date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
            </Text>
          </View>
        </View>

        {/* Meta */}
        <View style={s.metaCard}>
          <View style={s.metaCol}>
            <View>
              <Text style={s.metaLabel}>Branch</Text>
              <Text style={s.metaValue}>{sale.branch_name || '—'}</Text>
            </View>
            <View>
              <Text style={s.metaLabel}>Cashier</Text>
              <Text style={s.metaValue}>{sale.cashier_name || '—'}</Text>
            </View>
          </View>
          <View style={s.metaCol}>
            <View>
              <Text style={s.metaLabel}>Customer</Text>
              <Text style={s.metaValue}>{sale.customer_phone || 'Walk-in'}</Text>
            </View>
            <View>
              <Text style={s.metaLabel}>Payment</Text>
              <View style={[s.payChip, chipStyle]}>
                <Text>{methodLabel(method)}</Text>
              </View>
            </View>
          </View>
        </View>

        {/* Items */}
        <Text style={s.sectionTitle}>Items ({items.length})</Text>
        <View style={s.table}>
          <View style={s.thead}>
            <Text style={[s.th, s.cItem]}>Item</Text>
            <Text style={[s.th, s.cQty]}>Qty</Text>
            <Text style={[s.th, s.cPrice]}>Unit</Text>
            <Text style={[s.th, s.cTotal]}>Amount</Text>
          </View>
          {items.map((it: any, idx: number) => (
            <View key={it.id ?? idx} style={s.trow} wrap={false}>
              <View style={s.cItem}>
                <Text style={s.itemName}>{it.product_name || it.product?.title || `Item #${it.product_id}`}</Text>
                {it.product?.sku && <Text style={s.itemSku}>SKU {it.product.sku}</Text>}
              </View>
              <Text style={[s.cellQty, s.cQty]}>{it.quantity}</Text>
              <Text style={[s.cellMoney, s.cPrice]}>{fmtUGX(it.unit_price)}</Text>
              <Text style={[s.cellTotal, s.cTotal]}>{fmtUGX(it.line_total)}</Text>
            </View>
          ))}
        </View>

        {/* Totals */}
        <View style={s.totalsWrap}>
          <View style={s.totalsCard}>
            <View style={s.totalsRow}>
              <Text style={s.totalsLabel}>Subtotal</Text>
              <Text style={s.totalsValue}>{fmtUGX(sale.subtotal)}</Text>
            </View>
            {sale.discount_amount > 0 && (
              <View style={s.totalsRow}>
                <Text style={[s.totalsLabel, { color: C.bad }]}>Discount</Text>
                <Text style={[s.totalsValue, { color: C.bad }]}>−{fmtUGX(sale.discount_amount)}</Text>
              </View>
            )}
            <View style={s.grandRow}>
              <Text style={s.grandLabel}>TOTAL</Text>
              <Text style={s.grandValue}>{fmtUGX(sale.total)}</Text>
            </View>
          </View>
        </View>

        {/* Footer */}
        <View style={s.footer}>
          <Text style={s.thanks}>Thank you for your business!</Text>
          <Text style={s.footerLine}>Goods sold are not returnable without this receipt.</Text>
          <Text style={s.footerStrong}>{businessName} · {sale.branch_name || ''}</Text>
        </View>
      </Page>
    </Document>
  )
}
