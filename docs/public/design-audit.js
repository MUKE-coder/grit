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
    sectionPaddings: 3,
    shadows: 2,
    sections: 12,
    headings: 12,
    accentPerViewport: 3,
    lineLength: 72,
    contrastFailures: 0,
  };

  /*
   * The radius scale, in px, as --radius-* declares it in globals.css and as
   * the design system's component CSS draws it. 0 is always fine.
   *
   * This used to be a count ("no more than four distinct radii"), which was
   * the right check while there was no published scale: counting distinct
   * values was a proxy for having one. There is a scale now, so the question
   * is no longer how many corners a page has but whether any of them is a
   * value nobody chose. A page may legitimately use every step.
   *
   * 7 is here because .g-tab sits inside .g-tabs (10) with a 3px inset, and
   * 10 - 3 is what keeps the two curves concentric. 2 is the line-tab
   * indicator, a 2px bar rounded to its own height. Both are derived from a
   * scale value rather than chosen beside it.
   */
  const ALLOWED_RADII = [0, 2, 4, 6, 7, 8, 10, 12, 16, 20, 9999];

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
  // Parse any colour the browser can produce, by making the browser do it.
  //
  // A regex over the numbers worked while every computed colour came back as
  // rgb(0-255). It stopped the moment the tokens moved to color-mix() and
  // relative colour syntax: those compute to `color(srgb 0.63 0.65 0.69)`,
  // whose 0-1 floats read as near-black, and the audit reported 276 contrast
  // failures at exactly 1:1 on text anybody could see was fine.
  //
  // A 1x1 canvas normalises every syntax there is, including ones that do not
  // exist yet.
  const _probe = document.createElement('canvas');
  _probe.width = _probe.height = 1;
  const _ctx = _probe.getContext('2d', { willReadFrequently: true });
  const _rgbCache = new Map();
  const rgb = (c) => {
    if (!c || c === 'transparent' || c === 'none') return null;
    if (_rgbCache.has(c)) return _rgbCache.get(c);
    let out = null;
    try {
      _ctx.clearRect(0, 0, 1, 1);
      _ctx.fillStyle = '#000';
      _ctx.fillStyle = c;
      // An unparseable value leaves fillStyle at the previous one, so a
      // colour that silently failed is caught rather than measured as black.
      _ctx.fillRect(0, 0, 1, 1);
      const d = _ctx.getImageData(0, 0, 1, 1).data;
      out = { r: d[0], g: d[1], b: d[2], a: d[3] / 255 };
    } catch {
      out = null;
    }
    _rgbCache.set(c, out);
    return out;
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
  // A compound value is several corners of one radius, so each component is
  // checked on its own rather than the whole string being waved through.
  const onScale = (part) => {
    const n = Math.round(parseFloat(part));
    if (Number.isNaN(n)) return true; // a percentage or a keyword, not a px value
    return ALLOWED_RADII.includes(n === 0 ? 0 : n > 500 ? 9999 : n);
  };
  const offScale = radiusValues.filter(([v]) => !v.split(' ').every(onScale));
  if (offScale.length) {
    fail('off-scale border radii', offScale.length, 0,
      'the scale is 4 / 6 / 8 / 10 / 12 / 16 / 20 / full. Off it here: ' +
      offScale.map(([v, n]) => `${v} (x${n})`).join(', '));
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
    /*
     * A fill, not a wash.
     *
     * `bg-accent/10` computes to rgba(169, 155, 255, 0.1): the same three
     * channels as the fill, at a tenth of the strength. Comparing r, g and
     * b alone counted every soft accent panel on the page as a call to
     * action, which on a docs page is five of them and on no page is true.
     * A 10% wash is a surface; an opaque fill is a button.
     */
    const near = (c) => c && c.a >= 0.85 &&
      Math.abs(c.r - accentRGB.r) < 24 &&
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

  /* --- 7. line length ----------------------------------------------------
   *
   * Measured, not assumed, and definitely not in `ch`.
   *
   * The `ch` unit is the advance of the digit "0". Geist's digits are
   * tabular and wide: `ch` is 0.663em where a real prose character averages
   * 0.461em, so `max-width: 72ch` holds 104 characters in this typeface.
   * Anything reported in ch here would be 44% optimistic, and a page could
   * pass while running a third longer than the rule allows.
   *
   * So each paragraph's own font renders a sentence off-screen, and the
   * width of that divided by its length is what one character costs. It
   * costs one layout per distinct font, cached by the font shorthand. */
  const charWidthCache = new Map();
  const SAMPLE =
    'The quick brown fox jumps over the lazy dog, and then it does so again, ' +
    'because typography is measured in characters and not in zeroes.';
  const charWidth = (cs) => {
    const key = `${cs.fontWeight} ${cs.fontSize} ${cs.fontFamily}`;
    if (charWidthCache.has(key)) return charWidthCache.get(key);
    const span = document.createElement('span');
    span.style.cssText = 'position:absolute;visibility:hidden;white-space:pre;top:-9999px';
    span.style.fontFamily = cs.fontFamily;
    span.style.fontSize = cs.fontSize;
    span.style.fontWeight = cs.fontWeight;
    span.style.letterSpacing = cs.letterSpacing;
    span.textContent = SAMPLE;
    document.body.appendChild(span);
    const w = span.getBoundingClientRect().width / SAMPLE.length;
    span.remove();
    charWidthCache.set(key, w);
    return w;
  };

  const longLines = [];
  for (const p of [...document.querySelectorAll('p')].filter(visible)) {
    const cs = getComputedStyle(p);
    const w = charWidth(cs);
    if (!w) continue;
    const chars = p.getBoundingClientRect().width / w;
    if (chars > BUDGETS.lineLength && p.textContent.trim().length > 120) {
      longLines.push(Math.round(chars));
    }
  }
  if (longLines.length) {
    fail('paragraphs over the measure', longLines.length, 0,
      `widest ~${Math.max(...longLines)} characters, against ${BUDGETS.lineLength}. ` +
      'Set max-width to var(--measure), never to a ch value.');
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
      /*
       * WCAG 2.5.8 exempts a target "in a sentence or block of text", and
       * what that means in practice is a text link whose height is its line
       * box. A breadcrumb, a meta row and a list of text actions are all
       * made of those, and flagging them taught people to ignore the check.
       *
       * So the exemption is: visible text, no padding, no background and no
       * border. An icon button has no text and fails it; a 20px-tall link
       * with a background is a control and fails it too.
       */
      const s = getComputedStyle(e);
      // "A text link" means its height is its line box: no vertical padding,
      // no fill, no border. Horizontal padding is allowed, because a link in
      // a row of links often carries some and it does not make the target
      // any shorter.
      const lineHeight = parseFloat(s.lineHeight) || parseFloat(s.fontSize) * 1.2;
      const bare =
        (e.innerText || '').trim().length > 0 &&
        s.paddingTop === '0px' && s.paddingBottom === '0px' &&
        s.borderTopWidth === '0px' && s.borderBottomWidth === '0px' &&
        (rgb(s.backgroundColor)?.a ?? 0) < 0.05 &&
        r.height <= lineHeight + 4;
      if (bare) return false;
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
      offScaleRadii: { found: offScale.length, limit: 0, distinctInUse: radiusValues.length },
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
