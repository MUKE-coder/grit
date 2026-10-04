"use client";

import { ResourcePage } from "@/components/resource/resource-page";
import { loanResource } from "@/resources/loans/loans";

export default function LoansPage() {
  return <ResourcePage resource={loanResource} />;
}
