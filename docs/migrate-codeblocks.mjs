import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

function findPageFiles(dir) {
  const results = [];
  const items = fs.readdirSync(dir, { withFileTypes: true });
  for (const item of items) {
    const fp = path.join(dir, item.name);
    if (item.isDirectory()) results.push(...findPageFiles(fp));
    else if (item.name === 'page.tsx') results.push(fp);
  }
  return results;
}

function inferLang(label) {
  if (!label) return 'go';
  const l = label.trim().toLowerCase();
  if (l.endsWith('.go')) return 'go';
  if (l.endsWith('.tsx')) return 'tsx';
  if (l.endsWith('.ts') && !l.endsWith('.tsx')) return 'typescript';
  if (l.endsWith('.yaml') || l.endsWith('.yml') || l.includes('docker-compose')) return 'yaml';
  if (l.endsWith('.json') || l === 'package.json' || l === 'tsconfig.json') return 'json';
  if (l.endsWith('.css')) return 'css';
  if (l.endsWith('.env') || l === '.env' || l === '.env.example') return 'bash';
  if (l.startsWith('dockerfile') || l === 'dockerfile') return 'dockerfile';
  if (l === 'terminal') return 'bash';
  if (l.endsWith('.html') || l.endsWith('.htm')) return 'markup';
  if (l.includes('response') || l === 'request body') return 'json';
  if (l === 'output' || l === 'structure' || l === 'generated files' || l === 'file structure') return 'bash';
  return 'go';
}

const docsDir = path.join(__dirname, 'app/docs');
const allFiles = findPageFiles(docsDir);

let totalFiles = 0, totalBlocks = 0;

for (const fp of allFiles) {
  let src = fs.readFileSync(fp, 'utf8');
  if (!src.includes('<pre className="p-5')) continue;

  // Match code blocks: optional prefix classes, outer div > header div with optional dots > optional span > /header > pre > code > /pre > /outer
  // Handles: ml-10 prefix, leading-relaxed on pre, inline flex gap-1.5 dots
  const rx = /<div className="[^"]*rounded-xl border border-border\/40 bg-card\/80 overflow-hidden[^"]*">\s*<div className="flex items-center[^>]*>\s*(?:<div className="flex[^"]*gap-1\.5[^"]*">[\s\S]*?<\/div>\s*)?(?:<span className="[^"]*text-\[11px\] font-mono text-muted-foreground\/40">([^<]*)<\/span>\s*)?<\/div>\s*<pre className="p-5 text-sm font-mono text-foreground\/80 overflow-x-auto[^"]*">\s*([\s\S]*?)\s*<\/pre>\s*<\/div>/g;

  let count = 0;
  const result = src.replace(rx, (m, rawLabel, rawCode) => {
    count++;
    const label = (rawLabel || '').trim();
    const code = (rawCode || '').trim();
    const lang = inferLang(label);

    const attrs = [];
    if (lang !== 'go') attrs.push(`language="${lang}"`);
    if (label && label !== 'terminal') attrs.push(`filename="${label}"`);
    attrs.push(`code=${code}`);

    return `<CodeBlock ${attrs.join(' ')} />`;
  });

  if (count === 0) continue;

  let final = result;

  // Add import if missing
  if (!final.includes("code-block'") && !final.includes('code-block"')) {
    const lines = final.split('\n');
    let insertAt = 0;
    for (let i = 0; i < Math.min(lines.length, 30); i++) {
      const trimmed = lines[i].trimStart();
      if (trimmed.startsWith('import ')) {
        if (lines[i].includes('{') && !lines[i].includes('}')) {
          for (let j = i + 1; j < lines.length; j++) {
            if (lines[j].includes('}')) { insertAt = j + 1; break; }
          }
        } else {
          insertAt = i + 1;
        }
      }
    }
    lines.splice(insertAt, 0, "import { CodeBlock } from '@/components/code-block'");
    final = lines.join('\n');
  }

  fs.writeFileSync(fp, final);
  totalFiles++;
  totalBlocks += count;
  const rel = path.relative(process.cwd(), fp);
  console.log(`${rel}: ${count} blocks`);
}

console.log(`\nDone: ${totalBlocks} blocks in ${totalFiles} files`);
