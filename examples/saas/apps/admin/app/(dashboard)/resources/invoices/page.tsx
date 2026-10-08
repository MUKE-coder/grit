"use client";

import { ResourcePage } from "@/components/resource/resource-page";
import { invoiceResource } from "@/resources/invoices/invoices";

export default function InvoicesPage() {
  return <ResourcePage resource={invoiceResource} />;
}
