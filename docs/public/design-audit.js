/**
 * design-audit.js — measures a rendered page against DESIGN.md §9.
 *
 * Self-contained on purpose: it is one expression, with no imports and no
 * build step, so it can be dropped straight into a browser evaluate call by
 * whatever tool is driving the page. Paste the whole file as the function
 * body, or read it and pass it through.
 *
 * It reports numbers, not opinions. "This page has 40 type styles" is a fact
 * somebody can act on; "the typography feels inconsistent" is not.
 *
 * Returns { budgets, failures, details }. `failures` is the list to fix.
 */
(() => {
  // --- DESIGN.md §9 -------------------------------------------------------
  const BUDGETS = {
    typeStyles: 12,
    radii: 4,
    sectionPaddings: 3,
    shadows: 2,
    sections: 12,
    headings: 12,
    accentPerViewport: 3,
    lineLength: 72,
    contrastFailures: 0,
  };

  // The four radii DESIGN.md §4 allows, in px. 0 is always fine.
  const ALLOWED_RADII = [0, 6, 10, 16, 9999];

  // The three section paddings DESIGN.md §3 allows, in px.
  const ALLOWED_SECTION_PADDING = [0, 64, 96, 128];

  const visible = (el) => {
    const r = el.getBoundingClientRect();
    if (r.width === 0 && r.height === 0) return false;
    const s = getComputedStyle(el);
    if (s.display === 'none' || s.visibility === 'hidden' || s.opacity === '0') return false;
    // Screen-reader-only text keeps a 1x1 box and clips it away, so a
    // rect test calls it visible. It is not, and flagging its contrast is
    // noise that teaches people to ignore the audit.
    if (r.width <= 1 && r.height <= 1) return false;
    if (s.clipPath === 'inset(50%)' || s.clip === 'rect(0px, 0px, 0px, 0px)') return false;
    if (el.closest('[aria-hidden="true"]')) return false;
    return true;
  };

  const els = [...document.querySelectorAll('body *')].filter(visible);

  // --- colour helpers -----------------------------------------------------
  const rgb = (c) => {
    const m = (c || '').match(/-?\d+(\.\d+)?/g);
    if (!m || m.length < 3) return null;
    return { r: +m[0], g: +m[1], b: +m[2], a: m.length > 3 ? +m[3] : 1 };
  };
  const lum = (c) => {
    const f = (v) => {
      v /= 255;
      return v <= 0.03928 ? v / 12.92 : Math.pow((v + 0.055) / 1.055, 2.4);
    };
    return 0.2126 * f(c.r) + 0.7152 * f(c.g) + 0.0722 * f(c.b);
  };
  const ratio = (a, b) => {
    const l1 = lum(a), l2 = lum(b);
    return (Math.max(l1, l2) + 0.05) / (Math.min(l1, l2) + 0.05);
  };
  // The nearest ancestor that actually paints a background, which is what the
  // text is really sitting on. Walking up matters: most text sits on a
  // transparent parent.
  const backdrop = (el) => {
    let node = el;
    while (node && node !== document.documentElement) {
      const c = rgb(getComputedStyle(node).backgroundColor);
      if (c && c.a > 0.95) return c;
      node = node.parentElement;
    }
    return rgb(getComputedStyle(document.body).backgroundColor) || { r: 255, g: 255, b: 255, a: 1 };
  };

  const tally = (list) => {
    const m = new Map();
    for (const k of list) if (k != null) m.set(k, (m.get(k) || 0) + 1);
    return [...m.entries()].sort((a, b) => b[1] - a[1]);
  };

  const failures = [];
  const fail = (what, found, limit, hint) =>
    failures.push({ what, found, limit, hint });

  // --- 1. type styles -----------------------------------------------------
  // Only elements whose own text is the leaf text, so a wrapper does not
  // count its children's style as a second one.
  const textLeaves = els.filter((e) =>
    [...e.childNodes].some((n) => n.nodeType === 3 && n.textContent.trim())
  );
  const typeStyles = tally(
    textLeaves.map((e) => {
      const s = getComputedStyle(e);
      return `${Math.round(parseFloat(s.fontSize))}px/${s.fontWeight}`;
    })
  );
  if (typeStyles.length > BUDGETS.typeStyles) {
    fail('type styles', typeStyles.length, BUDGETS.typeStyles,
      'the rarest are the ones to delete: ' +
      typeStyles.slice(-6).map(([k, n]) => `${k} (x${n})`).join(', '));
  }

  // --- 2. radii -----------------------------------------------------------
  // A compound radius ("0px 10px 10px 0px") is one value applied to some
  // corners, not a new one. Normalising to the distinct non-zero components
  // means a card with a rounded top does not read as three radii.
  const normaliseRadius = (r) => {
    const parts = [...new Set(r.split('/')[0].trim().split(/\s+/))]
      .filter((p) => p && p !== '0px');
    return parts.length ? parts.sort().join(' ') : null;
  };
  const radiusValues = tally(
    els.map((e) => {
      const r = getComputedStyle(e).borderRadius;
      return r && r !== '0px' ? normaliseRadius(r) : null;
    })
  );
  const offScale = radiusValues.filter(([v]) => {
    const n = Math.round(parseFloat(v));
    return !v.includes(' ') && !ALLOWED_RADII.includes(n === 0 ? 0 : n > 500 ? 9999 : n);
  });
  if (radiusValues.length > BUDGETS.radii) {
    fail('border radii', radiusValues.length, BUDGETS.radii,
      'allowed: 6, 10, 16, full. Off-scale here: ' +
      (offScale.length ? offScale.map(([v, n]) => `${v} (x${n})`).join(', ') : 'none, just too many'));
  }

  // --- 3. section rhythm --------------------------------------------------
  const sections = [...document.querySelectorAll('section')].filter(visible);
  const paddings = tally(
    sections.map((s) => {
      const c = getComputedStyle(s);
      return `${Math.round(parseFloat(c.paddingTop))}/${Math.round(parseFloat(c.paddingBottom))}`;
    })
  );
  if (paddings.length > BUDGETS.sectionPaddings) {
    fail('section paddings', paddings.length, BUDGETS.sectionPaddings,
      'allowed: 64, 96, 128. Found: ' + paddings.map(([v, n]) => `${v} (x${n})`).join(', '));
  }
  if (sections.length > BUDGETS.sections) {
    fail('sections', sections.length, BUDGETS.sections,
      'a reader gives a page three scrolls; sections past twelve bury the ones above them');
  }

  // --- 4. shadows ---------------------------------------------------------
  const shadows = tally(
    els.map((e) => {
      const s = getComputedStyle(e).boxShadow;
      return s && s !== 'none' ? s : null;
    })
  );
  if (shadows.length > BUDGETS.shadows) {
    fail('shadows', shadows.length, BUDGETS.shadows,
      'structure comes from 1px lines; a shadow is for something that floats');
  }

  // --- 5. headings --------------------------------------------------------
  const h1s = [...document.querySelectorAll('h1')].filter(visible);
  const headings = [...document.querySelectorAll('h1,h2')].filter(visible);
  if (headings.length > BUDGETS.headings) {
    fail('h1 + h2', headings.length, BUDGETS.headings,
      'every heading competes with every other one');
  }
  if (h1s.length !== 1) {
    fail('h1 count', h1s.length, 1, 'exactly one h1 per page');
  }
  // Levels must not skip: an h4 under an h2 with no h3 between.
  const levels = [...document.querySelectorAll('h1,h2,h3,h4,h5,h6')]
    .filter(visible).map((h) => +h.tagName[1]);
  const skips = [];
  for (let i = 1; i < levels.length; i++) {
    if (levels[i] - levels[i - 1] > 1) skips.push(`h${levels[i - 1]} -> h${levels[i]}`);
  }
  if (skips.length) {
    fail('heading level skips', skips.length, 0, skips.slice(0, 4).join(', '));
  }

  // --- 6. accent rationing ------------------------------------------------
  const accentRaw = getComputedStyle(document.documentElement)
    .getPropertyValue('--accent').trim();
  let accentRGB = null;
  if (accentRaw) {
    const probe = document.createElement('div');
    probe.style.color = accentRaw.startsWith('#') || accentRaw.startsWith('rgb')
      ? accentRaw : `hsl(${accentRaw})`;
    document.body.appendChild(probe);
    accentRGB = rgb(getComputedStyle(probe).color);
    probe.remove();
  }
  let accentInView = 0;
  if (accentRGB) {
    const near = (c) => c && Math.abs(c.r - accentRGB.r) < 24 &&
      Math.abs(c.g - accentRGB.g) < 24 && Math.abs(c.b - accentRGB.b) < 24;
    for (const e of els) {
      const r = e.getBoundingClientRect();
      if (r.top > window.innerHeight || r.bottom < 0) continue;
      const s = getComputedStyle(e);
      // Fills only. Accent-tinted text and icons are a set; an accent fill
      // is a call to action, and two of those compete.
      if (near(rgb(s.backgroundColor))) accentInView++;
    }
    if (accentInView > BUDGETS.accentPerViewport) {
      fail('accent-filled elements in the first viewport', accentInView,
        BUDGETS.accentPerViewport,
        'a fill is a call to action, and two of those compete: one primary per view');
    }
  }

  // --- 7. line length -----------------------------------------------------
  const longLines = [];
  for (const p of [...document.querySelectorAll('p')].filter(visible)) {
    const s = getComputedStyle(p);
    const ch = parseFloat(s.fontSize) * 0.5; // approx advance of one character
    const chars = p.getBoundingClientRect().width / ch;
    if (chars > BUDGETS.lineLength && p.textContent.trim().length > 120) {
      longLines.push(Math.round(chars));
    }
  }
  if (longLines.length) {
    fail('paragraphs wider than 72ch', longLines.length, 0,
      `widest ~${Math.max(...longLines)}ch. Put prose in max-w-3xl`);
  }

  // --- 8. contrast --------------------------------------------------------
  const contrastFails = [];
  for (const e of textLeaves) {
    const s = getComputedStyle(e);
    const fg = rgb(s.color);
    if (!fg || fg.a < 0.95) continue;
    const size = parseFloat(s.fontSize);
    const bold = +s.fontWeight >= 700;
    const large = size >= 24 || (size >= 18.66 && bold);
    const need = large ? 3 : 4.5;
    const bg = backdrop(e);
    const got = ratio(fg, bg);
    if (got < need) {
      contrastFails.push({
        text: e.textContent.trim().slice(0, 40),
        where: (e.tagName.toLowerCase() +
          (e.className ? '.' + String(e.className).trim().split(/\s+/).slice(0, 2).join('.') : '')),
        ratio: +got.toFixed(2), need, size: Math.round(size),
        fg: s.color, bg: `rgb(${bg.r}, ${bg.g}, ${bg.b})`,
      });
    }
  }
  if (contrastFails.length) {
    fail('contrast failures', contrastFails.length, 0,
      contrastFails.slice(0, 3)
        .map((c) => `"${c.text}" ${c.ratio}:1 needs ${c.need}`).join(' | '));
  }

  // --- 9. accessible names on icon-only controls --------------------------
  const unnamed = [];
  for (const el of [...document.querySelectorAll('a,button,[role="button"]')].filter(visible)) {
    const text = el.textContent.trim();
    if (text) continue;
    if (el.getAttribute('aria-label') || el.getAttribute('aria-labelledby') ||
        el.getAttribute('title') || el.querySelector('.sr-only')) continue;
    unnamed.push(el.tagName.toLowerCase() + (el.className ? '.' + String(el.className).split(' ')[0] : ''));
  }
  if (unnamed.length) {
    fail('icon-only controls with no accessible name', unnamed.length, 0,
      unnamed.slice(0, 4).join(', ') + ' — a screen reader reads these as the bare URL');
  }

  // --- 10. target size ----------------------------------------------------
  const small = [...document.querySelectorAll('a,button,[role="button"],input,select')]
    .filter(visible)
    .filter((e) => {
      const r = e.getBoundingClientRect();
      // Inline links inside a paragraph are exempt, as WCAG 2.5.8 allows.
      if (e.tagName === 'A' && e.closest('p')) return false;
      return r.width < 24 || r.height < 24;
    });
  if (small.length) {
    fail('targets under 24x24', small.length, 0,
      'WCAG 2.5.8. Pad the control rather than the icon');
  }

  return {
    url: location.pathname,
    viewport: `${window.innerWidth}x${window.innerHeight}`,
    budgets: {
      typeStyles: { found: typeStyles.length, limit: BUDGETS.typeStyles },
      radii: { found: radiusValues.length, limit: BUDGETS.radii },
      sectionPaddings: { found: paddings.length, limit: BUDGETS.sectionPaddings },
      shadows: { found: shadows.length, limit: BUDGETS.shadows },
      sections: { found: sections.length, limit: BUDGETS.sections },
      headings: { found: headings.length, limit: BUDGETS.headings },
      accentFillsInFirstViewport: { found: accentInView, limit: BUDGETS.accentPerViewport },
      contrastFailures: { found: contrastFails.length, limit: 0 },
    },
    failures,
    details: {
      typeStyles: typeStyles.map(([k, n]) => `${k} x${n}`),
      radii: radiusValues.map(([k, n]) => `${k} x${n}`),
      sectionPaddings: paddings.map(([k, n]) => `${k} x${n}`),
      contrast: contrastFails.slice(0, 10),
      unnamedControls: unnamed.slice(0, 10),
    },
    verdict: failures.length === 0
      ? 'within budget'
      : `${failures.length} budget(s) over`,
  };
})();
