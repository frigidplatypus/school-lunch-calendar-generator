var DESIGN_PRESETS = {
  'classic-school': {
    name: 'Classic School',
    description: 'Clean and professional',
    palette: { primary: '#2980b9', secondary: '#1a5276', accent: '#3498db', bg: '#eaf2f8', cellBg: '#ffffff' },
    headerSvg: null,
    footerSvg: null,
    badgeStyle: 'rounded',
    cellBorderRadius: '4px',
  },
  'back-to-school': {
    name: 'Back to School',
    description: 'Bright and energetic for the new year',
    palette: { primary: '#e74c3c', secondary: '#c0392b', accent: '#f39c12', bg: '#fdf2e9', cellBg: '#fffdf7' },
    headerSvg: 'backpack',
    footerSvg: 'pencil-line',
    badgeStyle: 'pill',
    cellBorderRadius: '6px',
  },
  'autumn-harvest': {
    name: 'Autumn Harvest',
    description: 'Warm fall colors and leaves',
    palette: { primary: '#d35400', secondary: '#a04000', accent: '#f39c12', bg: '#fdf2e9', cellBg: '#fffaf5' },
    headerSvg: 'leaves',
    footerSvg: 'acorn-line',
    badgeStyle: 'leaf',
    cellBorderRadius: '6px',
  },
  'winter-snow': {
    name: 'Winter Snow',
    description: 'Cool blues and snowflakes',
    palette: { primary: '#2c3e50', secondary: '#1a252f', accent: '#5dade2', bg: '#eaf2f8', cellBg: '#f8fbfe' },
    headerSvg: 'snowflakes',
    footerSvg: 'snowflake-line',
    badgeStyle: 'rounded',
    cellBorderRadius: '4px',
  },
  'spring-garden': {
    name: 'Spring Garden',
    description: 'Fresh greens and flowers',
    palette: { primary: '#27ae60', secondary: '#1e8449', accent: '#f39c12', bg: '#eafaf1', cellBg: '#f9fef9' },
    headerSvg: 'flowers',
    footerSvg: 'leaf-line',
    badgeStyle: 'pill',
    cellBorderRadius: '8px',
  },
  'summer-sunshine': {
    name: 'Summer Sunshine',
    description: 'Bright yellows and warm tones',
    palette: { primary: '#f39c12', secondary: '#d68910', accent: '#e74c3c', bg: '#fef9e7', cellBg: '#fffdf5' },
    headerSvg: 'sun',
    footerSvg: 'wave-line',
    badgeStyle: 'pill',
    cellBorderRadius: '8px',
  },
  'testing-week': {
    name: 'Testing Week',
    description: 'Focused and calm for assessments',
    palette: { primary: '#8e44ad', secondary: '#6c3483', accent: '#3498db', bg: '#f4ecf7', cellBg: '#faf8fc' },
    headerSvg: 'pencil',
    footerSvg: 'check-line',
    badgeStyle: 'rounded',
    cellBorderRadius: '4px',
  },
  'spirit-week': {
    name: 'Spirit Week',
    description: 'Fun and colorful celebration',
    palette: { primary: '#e74c3c', secondary: '#c0392b', accent: '#f1c40f', bg: '#fdedec', cellBg: '#fffdfc' },
    headerSvg: 'stars',
    footerSvg: 'confetti-line',
    badgeStyle: 'pill',
    cellBorderRadius: '10px',
  },
  'graduation': {
    name: 'Graduation',
    description: 'Celebrating achievement',
    palette: { primary: '#1a5276', secondary: '#154360', accent: '#f39c12', bg: '#eaf2f8', cellBg: '#f8fbfe' },
    headerSvg: 'gradcap',
    footerSvg: 'star-line',
    badgeStyle: 'rounded',
    cellBorderRadius: '4px',
  },
  'field-day': {
    name: 'Field Day',
    description: 'Outdoor fun and activities',
    palette: { primary: '#27ae60', secondary: '#1e8449', accent: '#e74c3c', bg: '#eafaf1', cellBg: '#f9fef9' },
    headerSvg: 'sun',
    footerSvg: 'wave-line',
    badgeStyle: 'pill',
    cellBorderRadius: '10px',
  },
  'teacher-appreciation': {
    name: 'Teacher Appreciation',
    description: 'Honoring educators',
    palette: { primary: '#e74c3c', secondary: '#c0392b', accent: '#f39c12', bg: '#fdedec', cellBg: '#fffdfc' },
    headerSvg: 'apple',
    footerSvg: 'heart-line',
    badgeStyle: 'pill',
    cellBorderRadius: '8px',
  },
  'halloween': {
    name: 'Halloween',
    description: 'Spooky and fun',
    palette: { primary: '#e67e22', secondary: '#d35400', accent: '#8e44ad', bg: '#fdf2e9', cellBg: '#fffaf5' },
    headerSvg: 'pumpkin',
    footerSvg: 'bat-line',
    badgeStyle: 'rounded',
    cellBorderRadius: '6px',
  },
};

var SVG_DECORATIONS = {
  backpack: '<svg viewBox="0 0 40 40" width="36" height="36" fill="none" xmlns="http://www.w3.org/2000/svg"><rect x="8" y="12" width="24" height="22" rx="4" stroke="currentColor" stroke-width="2"/><path d="M14 12V8a6 6 0 0 1 12 0v4" stroke="currentColor" stroke-width="2"/><rect x="14" y="20" width="12" height="6" rx="2" stroke="currentColor" stroke-width="1.5"/></svg>',

  pencil: '<svg viewBox="0 0 40 40" width="36" height="36" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M28 8l4 4-20 20H8v-4L28 8z" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/><path d="M24 12l4 4" stroke="currentColor" stroke-width="2"/></svg>',

  'pencil-line': '<svg viewBox="0 0 120 12" width="120" height="12" fill="none" xmlns="http://www.w3.org/2000/svg"><line x1="0" y1="6" x2="100" y2="6" stroke="currentColor" stroke-width="1" stroke-dasharray="4 3"/><path d="M104 2l8 4-8 4" stroke="currentColor" stroke-width="1.5" fill="none"/></svg>',

  leaves: '<svg viewBox="0 0 50 40" width="46" height="36" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M10 30C10 15 25 5 40 8C35 18 20 22 10 30Z" stroke="currentColor" stroke-width="1.5" fill="currentColor" opacity="0.15"/><path d="M10 30C18 25 28 18 40 8" stroke="currentColor" stroke-width="1.5"/><path d="M5 28C5 18 15 10 28 10C25 18 15 24 5 28Z" stroke="currentColor" stroke-width="1.5" fill="currentColor" opacity="0.1"/></svg>',

  'acorn-line': '<svg viewBox="0 0 120 12" width="120" height="12" fill="none" xmlns="http://www.w3.org/2000/svg"><line x1="0" y1="6" x2="45" y2="6" stroke="currentColor" stroke-width="1" stroke-dasharray="4 3"/><ellipse cx="60" cy="7" rx="4" ry="5" stroke="currentColor" stroke-width="1.2"/><path d="M56 3h8" stroke="currentColor" stroke-width="1.2"/><line x1="75" y1="6" x2="120" y2="6" stroke="currentColor" stroke-width="1" stroke-dasharray="4 3"/></svg>',

  snowflakes: '<svg viewBox="0 0 50 40" width="46" height="36" fill="none" xmlns="http://www.w3.org/2000/svg"><g stroke="currentColor" stroke-width="1.5" stroke-linecap="round"><line x1="15" y1="8" x2="15" y2="28"/><line x1="5" y1="18" x2="25" y2="18"/><line x1="8" y1="11" x2="22" y2="25"/><line x1="22" y1="11" x2="8" y2="25"/><circle cx="15" cy="18" r="2" fill="currentColor" opacity="0.2"/><line x1="35" y1="10" x2="35" y2="26"/><line x1="27" y1="18" x2="43" y2="18"/><line x1="30" y1="13" x2="40" y2="23"/><line x1="40" y1="13" x2="30" y2="23"/></g></svg>',

  'snowflake-line': '<svg viewBox="0 0 120 12" width="120" height="12" fill="none" xmlns="http://www.w3.org/2000/svg"><line x1="0" y1="6" x2="50" y2="6" stroke="currentColor" stroke-width="1" stroke-dasharray="4 3"/><g transform="translate(60,6)" stroke="currentColor" stroke-width="1.2" stroke-linecap="round"><line x1="0" y1="-4" x2="0" y2="4"/><line x1="-3.5" y1="-2" x2="3.5" y2="2"/><line x1="3.5" y1="-2" x2="-3.5" y2="2"/></g><line x1="70" y1="6" x2="120" y2="6" stroke="currentColor" stroke-width="1" stroke-dasharray="4 3"/></svg>',

  flowers: '<svg viewBox="0 0 50 40" width="46" height="36" fill="none" xmlns="http://www.w3.org/2000/svg"><circle cx="20" cy="16" r="3" fill="currentColor" opacity="0.3"/><g stroke="currentColor" stroke-width="1.2"><ellipse cx="20" cy="10" rx="3" ry="4" fill="currentColor" opacity="0.15"/><ellipse cx="20" cy="22" rx="3" ry="4" fill="currentColor" opacity="0.15"/><ellipse cx="14" cy="16" rx="4" ry="3" fill="currentColor" opacity="0.15"/><ellipse cx="26" cy="16" rx="4" ry="3" fill="currentColor" opacity="0.15"/></g><line x1="20" y1="24" x2="20" y2="36" stroke="currentColor" stroke-width="1.5"/><path d="M20 30c-4-2-6 0-8 2" stroke="currentColor" stroke-width="1.2"/><circle cx="38" cy="24" r="2" fill="currentColor" opacity="0.25"/><g stroke="currentColor" stroke-width="1"><ellipse cx="38" cy="20" rx="2" ry="3" fill="currentColor" opacity="0.1"/><ellipse cx="38" cy="28" rx="2" ry="3" fill="currentColor" opacity="0.1"/><ellipse cx="34" cy="24" rx="3" ry="2" fill="currentColor" opacity="0.1"/><ellipse cx="42" cy="24" rx="3" ry="2" fill="currentColor" opacity="0.1"/></g></svg>',

  'leaf-line': '<svg viewBox="0 0 120 12" width="120" height="12" fill="none" xmlns="http://www.w3.org/2000/svg"><line x1="0" y1="6" x2="48" y2="6" stroke="currentColor" stroke-width="1" stroke-dasharray="4 3"/><path d="M55 10C55 4 62 1 70 3C66 7 60 9 55 10Z" stroke="currentColor" stroke-width="1.2" fill="currentColor" opacity="0.15"/><line x1="80" y1="6" x2="120" y2="6" stroke="currentColor" stroke-width="1" stroke-dasharray="4 3"/></svg>',

  sun: '<svg viewBox="0 0 50 40" width="46" height="36" fill="none" xmlns="http://www.w3.org/2000/svg"><circle cx="25" cy="20" r="8" stroke="currentColor" stroke-width="1.5" fill="currentColor" opacity="0.15"/><g stroke="currentColor" stroke-width="1.5" stroke-linecap="round"><line x1="25" y1="6" x2="25" y2="10"/><line x1="25" y1="30" x2="25" y2="34"/><line x1="11" y1="20" x2="15" y2="20"/><line x1="35" y1="20" x2="39" y2="20"/><line x1="15" y1="10" x2="18" y2="13"/><line x1="32" y1="27" x2="35" y2="30"/><line x1="35" y1="10" x2="32" y2="13"/><line x1="18" y1="27" x2="15" y2="30"/></g></svg>',

  'wave-line': '<svg viewBox="0 0 120 12" width="120" height="12" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M0 6 Q10 0 20 6 Q30 12 40 6 Q50 0 60 6 Q70 12 80 6 Q90 0 100 6 Q110 12 120 6" stroke="currentColor" stroke-width="1.2" fill="none"/></svg>',

  stars: '<svg viewBox="0 0 50 40" width="46" height="36" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M15 8l2 5h5l-4 3 1.5 5L15 18l-4.5 3L12 16l-4-3h5z" stroke="currentColor" stroke-width="1.2" fill="currentColor" opacity="0.2"/><path d="M35 14l1.5 3.5h3.5l-3 2.2 1 3.5L35 21l-3 2.2 1-3.5-3-2.2h3.5z" stroke="currentColor" stroke-width="1" fill="currentColor" opacity="0.15"/><path d="M25 26l1 2.5h2.5l-2 1.5.7 2.5L25 31l-2.2 1.5.7-2.5-2-1.5h2.5z" stroke="currentColor" stroke-width="1" fill="currentColor" opacity="0.1"/></svg>',

  'confetti-line': '<svg viewBox="0 0 120 12" width="120" height="12" fill="none" xmlns="http://www.w3.org/2000/svg"><circle cx="15" cy="4" r="2" fill="currentColor" opacity="0.3"/><rect x="35" y="2" width="3" height="3" rx="0.5" transform="rotate(20 36 3.5)" fill="currentColor" opacity="0.25"/><circle cx="60" cy="8" r="1.5" fill="currentColor" opacity="0.2"/><rect x="80" y="3" width="2.5" height="2.5" rx="0.5" transform="rotate(-15 81 4.2)" fill="currentColor" opacity="0.3"/><circle cx="105" cy="6" r="2" fill="currentColor" opacity="0.2"/></svg>',

  gradcap: '<svg viewBox="0 0 50 40" width="46" height="36" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M5 18L25 8l20 10-20 10L5 18z" stroke="currentColor" stroke-width="1.5" fill="currentColor" opacity="0.1"/><path d="M15 22v8c0 3 4.5 6 10 6s10-3 10-6v-8" stroke="currentColor" stroke-width="1.5"/><line x1="40" y1="18" x2="40" y2="32" stroke="currentColor" stroke-width="1.5"/><circle cx="40" cy="33" r="1.5" fill="currentColor" opacity="0.3"/></svg>',

  'star-line': '<svg viewBox="0 0 120 12" width="120" height="12" fill="none" xmlns="http://www.w3.org/2000/svg"><line x1="0" y1="6" x2="48" y2="6" stroke="currentColor" stroke-width="1" stroke-dasharray="4 3"/><path d="M60 2l1.8 4h4.2l-3.4 2.5 1.3 4.1L60 10.2l-3.9 2.4 1.3-4.1L54 6h4.2z" stroke="currentColor" stroke-width="1" fill="currentColor" opacity="0.2"/><line x1="72" y1="6" x2="120" y2="6" stroke="currentColor" stroke-width="1" stroke-dasharray="4 3"/></svg>',

  pumpkin: '<svg viewBox="0 0 50 40" width="46" height="36" fill="none" xmlns="http://www.w3.org/2000/svg"><ellipse cx="25" cy="24" rx="14" ry="12" stroke="currentColor" stroke-width="1.5" fill="currentColor" opacity="0.12"/><path d="M18 24c0-8 3-14 7-14s7 6 7 14" stroke="currentColor" stroke-width="1"/><path d="M25 12c0-4 2-6 4-6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/><path d="M22 14c-1-2 0-5 2-6" stroke="currentColor" stroke-width="1" stroke-linecap="round"/></svg>',

  'bat-line': '<svg viewBox="0 0 120 12" width="120" height="12" fill="none" xmlns="http://www.w3.org/2000/svg"><line x1="0" y1="6" x2="45" y2="6" stroke="currentColor" stroke-width="1" stroke-dasharray="4 3"/><path d="M55 3c3 0 5 2 7 3-2 1-4 1-5 0-1 1-3 1-5 0 2-1 4-3 3-3z" stroke="currentColor" stroke-width="1" fill="currentColor" opacity="0.25"/><line x1="75" y1="6" x2="120" y2="6" stroke="currentColor" stroke-width="1" stroke-dasharray="4 3"/></svg>',

  apple: '<svg viewBox="0 0 50 40" width="46" height="36" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M25 14c-2-4-1-8 1-10" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/><path d="M25 14c6-2 14 2 16 12-2 10-10 14-16 10-6 4-14 0-16-10 2-10 10-14 16-12z" stroke="currentColor" stroke-width="1.5" fill="currentColor" opacity="0.12"/><path d="M25 14c-3 0-6 4-6 10" stroke="currentColor" stroke-width="1" opacity="0.4"/></svg>',

  'heart-line': '<svg viewBox="0 0 120 12" width="120" height="12" fill="none" xmlns="http://www.w3.org/2000/svg"><line x1="0" y1="6" x2="48" y2="6" stroke="currentColor" stroke-width="1" stroke-dasharray="4 3"/><path d="M60 9l-5-4.5c-1.5-1.3-3.5-.3-3.5 1.5 0 2.5 3.5 5 8.5 8 5-3 8.5-5.5 8.5-8 0-1.8-2-2.8-3.5-1.5L60 9z" stroke="currentColor" stroke-width="1" fill="currentColor" opacity="0.2"/><line x1="72" y1="6" x2="120" y2="6" stroke="currentColor" stroke-width="1" stroke-dasharray="4 3"/></svg>',
};

function getPreset(id) {
  return DESIGN_PRESETS[id] || DESIGN_PRESETS['classic-school'];
}

function presetSvgHtml(key, color) {
  if (!key || !SVG_DECORATIONS[key]) return '';
  return '<span class="preset-decoration" style="color:' + color + ';">' + SVG_DECORATIONS[key] + '</span>';
}

function presetBadgeStyle(preset, color) {
  var radius = preset.badgeStyle === 'pill' ? '10px' : '3px';
  return 'display:inline-block;background:' + color + ';color:#fff;font-weight:700;font-size:8pt;line-height:1;border-radius:' + radius + ';padding:2px 5px;margin-right:3px;vertical-align:baseline;';
}
