"use client";

// address-card-01 — Saved address card (Shopify My Account-inspired)
export function AddressCard({
  name = "Sarah Chen",
  line1 = "123 Market Street",
  line2 = "Apt 4B",
  city = "San Francisco",
  state = "CA",
  zip = "94105",
  country = "United States",
  phone = "+1 (415) 555-0142",
  isDefault = false,
  onEdit,
  onDelete,
  onSetDefault,
  onSelect,
  selected = false,
}: {
  name?: string; line1?: string; line2?: string; city?: string;
  state?: string; zip?: string; country?: string; phone?: string;
  isDefault?: boolean; selected?: boolean;
  onEdit?: () => void; onDelete?: () => void;
  onSetDefault?: () => void; onSelect?: () => void;
}) {
  return (
    <div onClick={onSelect} className={"p-4 rounded-xl border-2 transition-all cursor-pointer " + (selected ? "border-accent bg-accent/5" : "border-border hover:border-accent/30")}>
      <div className="flex items-start justify-between mb-3">
        <div className="flex items-center gap-2">
          <span className="text-sm font-semibold text-foreground">{name}</span>
          {isDefault && <span className="text-xs bg-accent/10 text-accent border border-accent/20 px-1.5 py-0.5 rounded-full">Default</span>}
        </div>
        {selected && <svg className="w-5 h-5 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>}
      </div>
      <address className="not-italic text-sm text-text-secondary flex flex-col gap-0.5">
        <span>{line1}{line2 && ", " + line2}</span>
        <span>{city}, {state} {zip}</span>
        <span>{country}</span>
        {phone && <span className="mt-1 text-text-muted">{phone}</span>}
      </address>
      <div className="flex items-center gap-3 mt-4 pt-3 border-t border-border">
        <button onClick={(e) => { e.stopPropagation(); onEdit?.(); }} className="text-xs text-text-secondary hover:text-foreground transition-colors">Edit</button>
        {!isDefault && <button onClick={(e) => { e.stopPropagation(); onSetDefault?.(); }} className="text-xs text-text-secondary hover:text-foreground transition-colors">Set as default</button>}
        <button onClick={(e) => { e.stopPropagation(); onDelete?.(); }} className="text-xs text-danger hover:underline ml-auto">Delete</button>
      </div>
    </div>
  );
}

export default AddressCard;
