"use client";

import { use } from "react";
import { ResourceDetailPage } from "@/components/resource/resource-detail-page";
import { invoiceResource } from "@/resources/invoices/invoices";

export default function InvoicesDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  return <ResourceDetailPage resource={invoiceResource} id={id} />;
}
