"use client";

import { ResourcePage } from "@/components/resource/resource-page";
import { optionResource } from "@/resources/options/options";

export default function OptionsPage() {
  return <ResourcePage resource={optionResource} />;
}
