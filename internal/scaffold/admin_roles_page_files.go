package scaffold

import "strings"

// adminRolesPage emits the Roles & permissions screen.
//
// One component, used by the Next.js admin and the Vite/TanStack admin.
// Building it once is deliberate: the two admins have drifted apart repeatedly
// (see #69), and a permission editor that disagreed between them would be a
// security-shaped bug, not a cosmetic one.
//
// UI contract:
//   - tri-state checkboxes at module/group/feature level. `indeterminate` is a
//     DOM property, not an attribute, so it has to be set through a ref.
//   - a CRUD matrix per feature, rendering an em dash for actions the feature
//     doesn't declare, rather than a checkbox that silently does nothing.
//   - live "N / total granted" counter and a text filter.
//   - copy-permissions-from, to start a new role from an existing one.
//   - built-in roles: name locked, delete hidden. The server enforces both — the
//     UI lock is a convenience, not the security boundary.
//
// Selection is seeded from the server's `expanded` list and collapsed back to
// wildcards on save, so a role granted a whole resource keeps inheriting actions
// added to it later.
func adminRolesPage() string {
	src := `"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import {
	usePermissionCatalog,
	useRoles,
	useCreateRole,
	useUpdateRole,
	useDeleteRole,
	collapseGrants,
	featureKeys,
	groupKeys,
	moduleKeys,
	type PermAction,
	type PermFeature,
	type PermGroup,
	type PermModule,
	type Role,
} from "@/hooks/use-roles";
import { PageHeader } from "@/components/chrome/PageHeader";
import { ShieldCheck, Lock, Users, Plus, ArrowLeft, Save, Trash2, Search, Copy } from "@/lib/icons";

const ACTIONS: PermAction[] = ["create", "view", "edit", "delete"];

type CheckState = "on" | "off" | "partial";

function checkState(granted: number, total: number): CheckState {
	if (total === 0 || granted === 0) return "off";
	return granted === total ? "on" : "partial";
}

function TriCheckbox({
	state,
	onChange,
	disabled,
	label,
}: {
	state: CheckState;
	onChange: (next: boolean) => void;
	disabled?: boolean;
	label: string;
}) {
	const ref = useRef<HTMLInputElement>(null);
	// indeterminate can only be set imperatively.
	useEffect(() => {
		if (ref.current) ref.current.indeterminate = state === "partial";
	}, [state]);

	return (
		<input
			ref={ref}
			type="checkbox"
			aria-label={label}
			className="h-4 w-4 shrink-0 cursor-pointer rounded border-border accent-accent disabled:cursor-not-allowed disabled:opacity-50"
			checked={state === "on"}
			disabled={disabled}
			onChange={(e) => onChange(e.target.checked)}
		/>
	);
}

function RolesList({
	roles,
	onEdit,
	onNew,
}: {
	roles: Role[];
	onEdit: (r: Role) => void;
	onNew: () => void;
}) {
	return (
		<div>
			<div className="mb-6 flex flex-wrap items-center justify-between gap-3">
				<p className="text-sm text-text-secondary">
					Define what each role can do. A role granted a whole resource keeps any
					actions added to it later.
				</p>
				<button
					onClick={onNew}
					className="inline-flex shrink-0 items-center gap-2 rounded-lg bg-accent px-4 py-2 text-sm font-semibold text-white transition-opacity hover:opacity-90"
				>
					<Plus className="h-4 w-4" />
					New role
				</button>
			</div>

			<div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
				{roles.map((role) => (
					<button
						key={role.id}
						onClick={() => onEdit(role)}
						className="rounded-xl border border-border bg-bg-elevated p-5 text-left transition-colors hover:border-accent/40"
					>
						<div className="mb-3 flex items-start justify-between gap-3">
							<div className="flex items-center gap-2">
								{role.is_system ? (
									<Lock className="h-4 w-4 shrink-0 text-text-muted" />
								) : (
									<ShieldCheck className="h-4 w-4 shrink-0 text-accent" />
								)}
								<span className="font-semibold text-text-primary">{role.name}</span>
							</div>
							{role.is_system ? (
								<span className="shrink-0 rounded border border-border px-1.5 py-0.5 font-mono text-[10px] uppercase tracking-wider text-text-muted">
									Built-in
								</span>
							) : null}
						</div>
						<p className="mb-4 line-clamp-2 text-sm text-text-secondary">
							{role.description || "No description."}
						</p>
						<div className="flex items-center justify-between text-xs text-text-muted">
							<span className="inline-flex items-center gap-1.5">
								<Users className="h-3.5 w-3.5" />
								{role.user_count} {role.user_count === 1 ? "user" : "users"}
							</span>
							<span className="font-mono">
								{role.grants.indexOf("*") >= 0
									? "all permissions"
									: role.expanded.length + " granted"}
							</span>
						</div>
					</button>
				))}
			</div>
		</div>
	);
}

function FeatureRow({
	feature,
	selected,
	toggle,
	disabled,
}: {
	feature: PermFeature;
	selected: Set<string>;
	toggle: (keys: string[], on: boolean) => void;
	disabled: boolean;
}) {
	const keys = featureKeys(feature);
	const granted = keys.filter((k) => selected.has(k)).length;

	return (
		<tr className="border-t border-border/60">
			<td className="py-2 pl-10 pr-4">
				<label className="flex cursor-pointer items-center gap-2.5">
					<TriCheckbox
						state={checkState(granted, keys.length)}
						onChange={(on) => toggle(keys, on)}
						disabled={disabled}
						label={feature.name}
					/>
					<span className="text-sm text-text-primary">{feature.name}</span>
				</label>
			</td>
			{ACTIONS.map((a) => {
				// Only actions the feature declares are meaningful. A view-only
				// report shows a dash rather than a checkbox that does nothing.
				if (feature.actions.indexOf(a) < 0) {
					return (
						<td key={a} className="px-3 py-2 text-center text-text-muted">
							&mdash;
						</td>
					);
				}
				const key = feature.key + "." + a;
				return (
					<td key={a} className="px-3 py-2 text-center">
						<input
							type="checkbox"
							aria-label={feature.name + " " + a}
							className="h-4 w-4 cursor-pointer rounded border-border accent-accent disabled:cursor-not-allowed disabled:opacity-50"
							checked={selected.has(key)}
							disabled={disabled}
							onChange={(e) => toggle([key], e.target.checked)}
						/>
					</td>
				);
			})}
		</tr>
	);
}

function ModuleSection({
	module,
	selected,
	toggle,
	disabled,
	filter,
}: {
	module: PermModule;
	selected: Set<string>;
	toggle: (keys: string[], on: boolean) => void;
	disabled: boolean;
	filter: string;
}) {
	const [open, setOpen] = useState(true);
	const q = filter.trim().toLowerCase();

	const groups = useMemo(() => {
		if (!q) return module.groups;
		const hit = (s: string) => s.toLowerCase().indexOf(q) >= 0;
		if (hit(module.name)) return module.groups;
		return module.groups
			.map((g) => ({
				...g,
				features: hit(g.name) ? g.features : g.features.filter((f) => hit(f.name)),
			}))
			.filter((g) => g.features.length > 0);
	}, [module, q]);

	if (groups.length === 0) return null;

	const keys = moduleKeys(module);
	const granted = keys.filter((k) => selected.has(k)).length;

	return (
		<div className="border-b border-border last:border-b-0">
			<div className="flex items-center gap-3 px-4 py-3">
				<TriCheckbox
					state={checkState(granted, keys.length)}
					onChange={(on) => toggle(keys, on)}
					disabled={disabled}
					label={module.name}
				/>
				<button
					onClick={() => setOpen(!open)}
					className="flex flex-1 items-center justify-between text-left"
				>
					<span className="font-semibold text-text-primary">{module.name}</span>
					<span className="font-mono text-xs text-text-muted">
						{granted} / {keys.length}
					</span>
				</button>
			</div>

			{open ? (
				<div className="pb-2">
					{groups.map((g) => {
						const gKeys = groupKeys(g as PermGroup);
						const gGranted = gKeys.filter((k) => selected.has(k)).length;
						return (
							<div key={g.key}>
								<div className="flex items-center gap-2.5 px-4 py-2 pl-8">
									<TriCheckbox
										state={checkState(gGranted, gKeys.length)}
										onChange={(on) => toggle(gKeys, on)}
										disabled={disabled}
										label={g.name}
									/>
									<span className="text-sm font-medium text-text-secondary">{g.name}</span>
								</div>
								<table className="w-full">
									<thead>
										<tr className="text-[10px] uppercase tracking-wider text-text-muted">
											<th className="pb-1 pl-10 pr-4 text-left font-medium">Feature</th>
											{ACTIONS.map((a) => (
												<th key={a} className="px-3 pb-1 text-center font-medium">
													{a}
												</th>
											))}
										</tr>
									</thead>
									<tbody>
										{g.features.map((f) => (
											<FeatureRow
												key={f.key}
												feature={f}
												selected={selected}
												toggle={toggle}
												disabled={disabled}
											/>
										))}
									</tbody>
								</table>
							</div>
						);
					})}
				</div>
			) : null}
		</div>
	);
}

function RoleEditor({
	role,
	roles,
	modules,
	totalKeys,
	onBack,
}: {
	role: Role | null;
	roles: Role[];
	modules: PermModule[];
	totalKeys: number;
	onBack: () => void;
}) {
	const create = useCreateRole();
	const update = useUpdateRole();
	const remove = useDeleteRole();

	const [name, setName] = useState(role ? role.name : "");
	const [description, setDescription] = useState(role ? role.description : "");
	const [filter, setFilter] = useState("");
	const [error, setError] = useState("");
	// Seeded from the server-expanded list — the client never resolves wildcards.
	const [selected, setSelected] = useState<Set<string>>(
		new Set(role ? role.expanded : [])
	);

	const isSuper = role ? role.grants.indexOf("*") >= 0 : false;
	const saving = create.isPending || update.isPending || remove.isPending;

	function toggle(keys: string[], on: boolean) {
		setSelected((prev) => {
			const next = new Set(prev);
			keys.forEach((k) => (on ? next.add(k) : next.delete(k)));
			return next;
		});
	}

	async function save() {
		setError("");
		if (!name.trim()) {
			setError("Name is required.");
			return;
		}
		// Collapse to wildcards so the role inherits future actions.
		const grants = isSuper ? ["*"] : collapseGrants(selected, modules);
		try {
			if (role) {
				await update.mutateAsync({ id: role.id, name, description, grants });
			} else {
				await create.mutateAsync({ name, description, grants });
			}
			onBack();
		} catch (e: any) {
			setError(
				(e && e.response && e.response.data && e.response.data.error
					? e.response.data.error.message
					: null) || "Could not save the role."
			);
		}
	}

	async function destroy() {
		if (!role) return;
		const ok = confirm(
			"Delete the role " + role.name + "? Users holding it will lose its permissions."
		);
		if (!ok) return;
		setError("");
		try {
			await remove.mutateAsync(role.id);
			onBack();
		} catch (e: any) {
			setError(
				(e && e.response && e.response.data && e.response.data.error
					? e.response.data.error.message
					: null) || "Could not delete the role."
			);
		}
	}

	return (
		<div>
			<div className="mb-6 flex items-center justify-between gap-4">
				<button
					onClick={onBack}
					className="inline-flex items-center gap-2 text-sm text-text-secondary hover:text-text-primary"
				>
					<ArrowLeft className="h-4 w-4" />
					All roles
				</button>
				<div className="flex items-center gap-2">
					{role && !role.is_system ? (
						<button
							onClick={destroy}
							disabled={saving}
							className="inline-flex items-center gap-2 rounded-lg border border-border px-3 py-2 text-sm text-danger transition-colors hover:bg-danger/10 disabled:opacity-50"
						>
							<Trash2 className="h-4 w-4" />
							Delete
						</button>
					) : null}
					<button
						onClick={save}
						disabled={saving}
						className="inline-flex items-center gap-2 rounded-lg bg-accent px-4 py-2 text-sm font-semibold text-white transition-opacity hover:opacity-90 disabled:opacity-50"
					>
						<Save className="h-4 w-4" />
						{saving ? "Saving..." : "Save role"}
					</button>
				</div>
			</div>

			{error ? (
				<p className="mb-4 rounded-lg border border-danger/30 bg-danger/10 px-4 py-3 text-sm text-danger">
					{error}
				</p>
			) : null}

			<div className="mb-4 grid gap-4 md:grid-cols-2">
				<div>
					<label className="mb-1.5 block text-xs font-medium uppercase tracking-wider text-text-muted">
						Name
					</label>
					<input
						value={name}
						onChange={(e) => setName(e.target.value)}
						disabled={role ? role.is_system : false}
						placeholder="e.g. Support"
						className="w-full rounded-lg border border-border bg-bg-primary px-3 py-2 text-sm text-text-primary disabled:opacity-60"
					/>
				</div>
				<div>
					<label className="mb-1.5 block text-xs font-medium uppercase tracking-wider text-text-muted">
						Description
					</label>
					<input
						value={description}
						onChange={(e) => setDescription(e.target.value)}
						placeholder="What is this role for?"
						className="w-full rounded-lg border border-border bg-bg-primary px-3 py-2 text-sm text-text-primary"
					/>
				</div>
			</div>

			{role && role.is_system ? (
				<p className="mb-4 rounded-lg border border-warning/30 bg-warning/10 px-4 py-3 text-sm text-text-secondary">
					<strong className="text-text-primary">Built-in role.</strong> Its permissions are
					editable, but the name is fixed and it cannot be deleted &mdash; routes and the
					upgrade path resolve this role by name.
				</p>
			) : null}

			{roles.length > 0 ? (
				<div className="mb-4 rounded-xl border border-border bg-bg-elevated p-4">
					<p className="mb-2 text-xs font-medium uppercase tracking-wider text-text-muted">
						Copy permissions from
					</p>
					<div className="flex flex-wrap gap-2">
						{roles
							.filter((r) => !role || r.id !== role.id)
							.map((r) => (
								<button
									key={r.id}
									onClick={() => setSelected(new Set(r.expanded))}
									className="inline-flex items-center gap-1.5 rounded-lg border border-border px-3 py-1.5 text-xs text-text-secondary transition-colors hover:border-accent/40 hover:text-text-primary"
								>
									<Copy className="h-3 w-3" />
									{r.name}
								</button>
							))}
					</div>
				</div>
			) : null}

			<div className="rounded-xl border border-border bg-bg-elevated">
				<div className="flex flex-wrap items-center justify-between gap-3 border-b border-border px-4 py-3">
					<span className="font-semibold text-text-primary">Permissions</span>
					<div className="flex items-center gap-3">
						<span className="font-mono text-xs text-text-muted">
							{isSuper ? "all permissions" : selected.size + " / " + totalKeys + " granted"}
						</span>
						<div className="relative">
							<Search className="pointer-events-none absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-text-muted" />
							<input
								value={filter}
								onChange={(e) => setFilter(e.target.value)}
								placeholder="Filter..."
								className="w-44 rounded-lg border border-border bg-bg-primary py-1.5 pl-8 pr-3 text-xs text-text-primary"
							/>
						</div>
					</div>
				</div>

				{isSuper ? (
					<p className="px-4 py-6 text-sm text-text-secondary">
						This role holds the <code className="font-mono text-accent">*</code> grant
						&mdash; every permission, including any added in future. Remove that grant to
						pick individual permissions.
					</p>
				) : (
					modules.map((m) => (
						<ModuleSection
							key={m.key}
							module={m}
							selected={selected}
							toggle={toggle}
							disabled={saving}
							filter={filter}
						/>
					))
				)}
			</div>
		</div>
	);
}

export default function RolesPage() {
	const { data: catalog, isLoading: catalogLoading } = usePermissionCatalog();
	const { data: roles, isLoading: rolesLoading } = useRoles();
	// undefined = list view, null = creating, Role = editing.
	const [editing, setEditing] = useState<Role | null | undefined>(undefined);

	const loading = catalogLoading || rolesLoading;

	// Re-read the edited role from the query cache so counts and expanded
	// grants stay fresh after a save.
	const current = useMemo(() => {
		if (editing === undefined || editing === null) return editing;
		const found = roles ? roles.find((r) => r.id === editing.id) : undefined;
		return found || editing;
	}, [editing, roles]);

	return (
		<div className="space-y-6">
			<PageHeader
				title="Roles & permissions"
				subtitle="Control exactly what each role can see and do."
			/>

			{loading ? (
				<p className="text-sm text-text-secondary">Loading...</p>
			) : editing === undefined ? (
				<RolesList
					roles={roles || []}
					onEdit={(r) => setEditing(r)}
					onNew={() => setEditing(null)}
				/>
			) : (
				<RoleEditor
					role={current || null}
					roles={roles || []}
					modules={catalog ? catalog.modules : []}
					totalKeys={catalog ? catalog.keys.length : 0}
					onBack={() => setEditing(undefined)}
				/>
			)}
		</div>
	);
}
`
	return strings.ReplaceAll(src, "~", "`")
}
