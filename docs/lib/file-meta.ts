// Shared file-type metadata for the <Files> and <FileTree> components.
// Edit the COLORS map to retint file types, or FOLDER emojis to taste.

const COLORS: Record<string, string> = {
  go: 'text-cyan-400',
  mod: 'text-cyan-400',
  sum: 'text-cyan-400',
  ts: 'text-blue-400',
  tsx: 'text-sky-400',
  js: 'text-yellow-400',
  jsx: 'text-yellow-400',
  mjs: 'text-yellow-400',
  cjs: 'text-yellow-400',
  json: 'text-amber-400',
  md: 'text-slate-300',
  mdx: 'text-slate-300',
  css: 'text-sky-400',
  scss: 'text-pink-400',
  html: 'text-orange-400',
  env: 'text-emerald-400',
  yml: 'text-violet-400',
  yaml: 'text-violet-400',
  toml: 'text-orange-300',
  sql: 'text-orange-300',
  sh: 'text-emerald-400',
  bash: 'text-emerald-400',
  prisma: 'text-teal-400',
  lock: 'text-muted-foreground/60',
  png: 'text-fuchsia-400',
  jpg: 'text-fuchsia-400',
  jpeg: 'text-fuchsia-400',
  svg: 'text-fuchsia-400',
  ico: 'text-fuchsia-400',
}

/** Tailwind text-colour class for a file, based on its extension. */
export function fileColor(name: string): string {
  // dotfiles like .env, .gitignore → use the trailing token as the "ext"
  const clean = name.replace(/\/+$/, '')
  const ext = clean.includes('.')
    ? clean.split('.').pop()!.toLowerCase()
    : clean.replace(/^\./, '').toLowerCase()
  return COLORS[ext] ?? 'text-muted-foreground/70'
}

/** Emoji for a folder, open or closed. Swap these to change the look site-wide. */
export function folderEmoji(open: boolean): string {
  return open ? '📂' : '📁'
}
