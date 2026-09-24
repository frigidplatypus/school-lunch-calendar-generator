const MONTHS = ['January','February','March','April','May','June','July','August','September','October','November','December'];
const DAYS_SHORT = ['Sun','Mon','Tue','Wed','Thu','Fri','Sat'];
const DAYS_FULL = ['Sunday','Monday','Tuesday','Wednesday','Thursday','Friday','Saturday'];
const DAY_NAMES = ['Monday','Tuesday','Wednesday','Thursday','Friday'];

const schoolInput = document.getElementById('school-input');
const footerMessageInput = document.getElementById('footer-message-input');
const logoInput = document.getElementById('logo-input');
const monthSelect = document.getElementById('month-select');
const yearSelect = document.getElementById('year-select');
const presetSelect = document.getElementById('preset-select');
const colorPicker = document.getElementById('color-picker');
const fontSelect = document.getElementById('font-select');
const bgStyle = document.getElementById('bg-style');
const borderStyle = document.getElementById('border-style');
const inputGrid = document.getElementById('input-grid');
const generateBtn = document.getElementById('generate-btn');
const calendarOutput = document.getElementById('calendar-output');
const downloadSection = document.getElementById('download-section');
const downloadBtn = document.getElementById('download-btn');
const exportBtn = document.getElementById('export-btn');
const importHolidaysFile = document.getElementById('import-holidays-file');
const importMenuItemsFile = document.getElementById('import-menu-items-file');
const importCalendarFile = document.getElementById('import-calendar-file');
const manageHolidaysBtn = document.getElementById('manage-holidays-btn');
const manageMenusBtn = document.getElementById('manage-menus-btn');
const logoutBtn = document.getElementById('logout-btn');
const changePwBtn = document.getElementById('change-pw-btn');

let logoDataUrl = null;
let holidays = [];
let menuItems = { breakfast: [], lunch: [] };

const FOOTER_STORAGE_KEY = 'cal_footer_message';

footerMessageInput.value = localStorage.getItem(FOOTER_STORAGE_KEY) || '';
footerMessageInput.addEventListener('input', function() {
  localStorage.setItem(FOOTER_STORAGE_KEY, footerMessageInput.value);
});

function applyThemeToHeading(color) {
  const h1 = document.querySelector('h1');
  if (h1) h1.style.color = color;
}

function initPresets() {
  var cards = document.getElementById('preset-cards');
  if (!cards) return;
  Object.keys(DESIGN_PRESETS).forEach(function(id) {
    var p = DESIGN_PRESETS[id];
    var opt = document.createElement('option');
    opt.value = id;
    opt.textContent = p.name;
    presetSelect.appendChild(opt);

    var card = document.createElement('div');
    card.className = 'preset-card';
    card.dataset.preset = id;
    card.innerHTML =
      '<div class="preset-card-name">' + p.name + '</div>' +
      '<div class="preset-card-desc">' + p.description + '</div>' +
      '<div class="preset-card-swatches">' +
        '<div class="preset-card-swatch" style="background:' + p.palette.primary + ';"></div>' +
        '<div class="preset-card-swatch" style="background:' + p.palette.secondary + ';"></div>' +
        '<div class="preset-card-swatch" style="background:' + p.palette.accent + ';"></div>' +
      '</div>';
    card.addEventListener('click', function() {
      presetSelect.value = id;
      applyPreset(id);
    });
    cards.appendChild(card);
  });
  presetSelect.value = 'classic-school';
  presetSelect.addEventListener('change', function() {
    applyPreset(presetSelect.value);
  });
}

function applyPreset(id) {
  var p = getPreset(id);
  colorPicker.value = p.palette.primary;
  applyThemeToHeading(p.palette.primary);
  document.querySelectorAll('.preset-card').forEach(function(c) {
    c.classList.toggle('active', c.dataset.preset === id);
  });
  scheduleStylePreview();
}

function getActivePreset() {
  return getPreset(presetSelect.value);
}

function showToast(message, type) {
  let container = document.getElementById('toast-container');
  if (!container) {
    container = document.createElement('div');
    container.id = 'toast-container';
    document.body.appendChild(container);
  }
  const toast = document.createElement('div');
  toast.className = 'toast' + (type === 'error' ? ' toast-error' : '');
  toast.textContent = message;
  container.appendChild(toast);
  requestAnimationFrame(() => toast.classList.add('show'));
  setTimeout(() => {
    toast.classList.remove('show');
    setTimeout(() => toast.remove(), 300);
  }, 3500);
}

function markDirty() {
  generateBtn.classList.add('btn-dirty');
}

function clearDirty() {
  generateBtn.classList.remove('btn-dirty');
}

document.addEventListener('input', function(e) {
  if (e.target.closest('#input-grid') || e.target === schoolInput) markDirty();
});

document.addEventListener('change', function(e) {
  if (e.target.classList && e.target.classList.contains('day-mode-toggle')) markDirty();
});

document.addEventListener('keydown', function(e) {
  if (e.key === 'Escape') {
    document.querySelectorAll('.modal.active').forEach(m => m.classList.remove('active'));
  }
});

function getMenuItemNames(type) {
  return menuItems[type].map(item => item.item_name);
}

document.querySelectorAll('.preset-color').forEach(el => {
  el.addEventListener('click', () => {
    colorPicker.value = el.dataset.color;
    applyThemeToHeading(el.dataset.color);
    scheduleStylePreview();
  });
});

colorPicker.addEventListener('input', () => applyThemeToHeading(colorPicker.value));

logoInput.addEventListener('change', async function(e) {
  const file = e.target.files[0];
  if (!file) return;
  const formData = new FormData();
  formData.append('logo', file);
  try {
    const res = await fetch('/api/logo', { method: 'POST', body: formData });
    if (res.ok) {
      logoDataUrl = URL.createObjectURL(file);
      renderCalendar(true);
    } else {
      showToast('Logo upload failed (must be png/jpg/gif/svg/webp, max 5MB)', 'error');
    }
  } catch(err) { showToast('Error uploading logo', 'error'); }
  logoInput.value = '';
});

document.querySelectorAll('.modal-close').forEach(btn => {
  btn.addEventListener('click', function() {
    document.getElementById(this.dataset.modal).classList.remove('active');
  });
});

window.addEventListener('click', function(e) {
  if (e.target.classList.contains('modal')) e.target.classList.remove('active');
});

logoutBtn.addEventListener('click', async function() {
  await fetch('/logout', { method: 'POST' });
  window.location.href = '/login';
});

changePwBtn.addEventListener('click', function() {
  document.getElementById('password-modal').classList.add('active');
});

document.getElementById('change-pw-submit').addEventListener('click', async function() {
  const newPw = document.getElementById('new-password').value;
  const confirmPw = document.getElementById('confirm-password').value;
  const pwError = document.getElementById('pw-error');
  pwError.style.display = 'none';

  if (newPw.length < 4) {
    pwError.textContent = 'Password must be at least 4 characters';
    pwError.style.display = 'block';
    return;
  }
  if (newPw !== confirmPw) {
    pwError.textContent = 'Passwords do not match';
    pwError.style.display = 'block';
    return;
  }

  const res = await fetch('/api/change-password', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ new_password: newPw })
  });

  if (res.ok) {
    showToast('Password changed successfully');
    document.getElementById('password-modal').classList.remove('active');
    document.getElementById('new-password').value = '';
    document.getElementById('confirm-password').value = '';
  } else {
    const data = await res.json();
    pwError.textContent = data.error || 'Failed to change password';
    pwError.style.display = 'block';
  }
});

manageHolidaysBtn.addEventListener('click', () => {
  renderHolidaysModal();
  document.getElementById('holidays-modal').classList.add('active');
});

manageMenusBtn.addEventListener('click', () => {
  renderMenusModal();
  document.getElementById('menus-modal').classList.add('active');
});

function renderHolidaysModal() {
  const list = document.getElementById('holidays-list');
  list.innerHTML = '';
  const byMonth = {};
  holidays.forEach(h => {
    if (!byMonth[h.month]) byMonth[h.month] = [];
    byMonth[h.month].push(h);
  });
  Object.keys(byMonth).sort((a,b) => a-b).forEach(month => {
    const section = document.createElement('div');
    section.className = 'month-section';
    section.innerHTML = '<h4>' + MONTHS[month] + '</h4>';
    byMonth[month].forEach(h => {
      const item = document.createElement('div');
      item.className = 'list-item';
      item.innerHTML = '<span>' + h.date + ' - ' + escapeHtml(h.name) + '</span><button onclick="deleteHoliday(' + h.id + ')">Delete</button>';
      section.appendChild(item);
    });
    list.appendChild(section);
  });
}

window.deleteHoliday = async function(id) {
  await fetch('/api/holidays/' + id, { method: 'DELETE' });
  holidays = holidays.filter(h => h.id !== id);
  renderHolidaysModal();
  renderInputGrid();
};

document.getElementById('add-holiday-btn').addEventListener('click', async function() {
  const month = parseInt(document.getElementById('holiday-month').value);
  const date = parseInt(document.getElementById('holiday-date').value);
  const name = document.getElementById('holiday-name').value.trim();
  if (!date || !name) { showToast('Please fill in all fields', 'error'); return; }
  const res = await fetch('/api/holidays', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ month: month, date: date, name: name })
  });
  if (res.ok) {
    const h = await res.json();
    holidays.push(h);
    holidays.sort((a,b) => a.month - b.month || a.date - b.date);
    renderHolidaysModal();
    renderInputGrid();
    document.getElementById('holiday-date').value = '';
    document.getElementById('holiday-name').value = '';
  }
});

function renderMenusModal() {
  const list = document.getElementById('menus-list');
  list.innerHTML = '';
  ['breakfast', 'lunch'].forEach(type => {
    const section = document.createElement('div');
    section.className = 'month-section';
    section.innerHTML = '<h4>' + (type === 'breakfast' ? 'Breakfast' : 'Lunch') + ' Items</h4>';
    menuItems[type].forEach(item => {
      const itemDiv = document.createElement('div');
      itemDiv.className = 'list-item';
      itemDiv.innerHTML = '<span>' + escapeHtml(item.item_name) + '</span><button onclick="deleteMenuItem(' + item.id + ')">Delete</button>';
      section.appendChild(itemDiv);
    });
    if (menuItems[type].length === 0) {
      section.innerHTML += '<p style="color:#999;font-size:0.85rem;padding:8px;">No items yet</p>';
    }
    list.appendChild(section);
  });
}

window.deleteMenuItem = async function(id) {
  await fetch('/api/menu-items/' + id, { method: 'DELETE' });
  menuItems.breakfast = menuItems.breakfast.filter(item => item.id !== id);
  menuItems.lunch = menuItems.lunch.filter(item => item.id !== id);
  renderMenusModal();
};

document.getElementById('add-breakfast-btn').addEventListener('click', async function() {
  const input = document.getElementById('breakfast-item-input');
  const item = input.value.trim();
  if (!item) return;
  const res = await fetch('/api/menu-items', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ meal_type: 'breakfast', item_name: item })
  });
  if (res.ok) {
    const newItem = await res.json();
    menuItems.breakfast.push(newItem);
    renderMenusModal();
  }
  input.value = '';
});

document.getElementById('add-lunch-btn').addEventListener('click', async function() {
  const input = document.getElementById('lunch-item-input');
  const item = input.value.trim();
  if (!item) return;
  const res = await fetch('/api/menu-items', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ meal_type: 'lunch', item_name: item })
  });
  if (res.ok) {
    const newItem = await res.json();
    menuItems.lunch.push(newItem);
    renderMenusModal();
  }
  input.value = '';
});

function initSelectors() {
  MONTHS.forEach((m, i) => {
    const opt = document.createElement('option');
    opt.value = i;
    opt.textContent = m;
    monthSelect.appendChild(opt);
  });
  const holidayMonthSelect = document.getElementById('holiday-month');
  MONTHS.forEach((m, i) => {
    const opt = document.createElement('option');
    opt.value = i;
    opt.textContent = m;
    holidayMonthSelect.appendChild(opt);
  });
  const now = new Date();
  const currentYear = now.getFullYear();
  for (let y = currentYear - 1; y <= currentYear + 3; y++) {
    const opt = document.createElement('option');
    opt.value = y;
    opt.textContent = y;
    if (y === currentYear) opt.selected = true;
    yearSelect.appendChild(opt);
  }
  monthSelect.value = now.getMonth();
}

function isHoliday(date, month) {
  return holidays.find(h => h.month === month && h.date === date);
}

function getWeekdaysInMonth(year, month) {
  const firstDay = new Date(year, month, 1);
  const lastDay = new Date(year, month + 1, 0);
  const weekStart = new Date(firstDay);
  weekStart.setDate(weekStart.getDate() - ((weekStart.getDay() + 6) % 7));
  const weeks = [];
  while (weekStart <= lastDay || weeks.length === 0) {
    const week = [];
    for (let weekday = 0; weekday < 5; weekday++) {
      const day = new Date(weekStart);
      day.setDate(day.getDate() + weekday);
      week.push(day.getMonth() === month ? {
        date: day.getDate(),
        dayOfWeek: day.getDay()
      } : null);
    }
    if (week.some(day => day !== null)) weeks.push(week);
    weekStart.setDate(weekStart.getDate() + 7);
  }
  return weeks;
}

function getFullMonthWeeks(year, month) {
  const firstDay = new Date(year, month, 1);
  const lastDay = new Date(year, month + 1, 0);
  const startOffset = firstDay.getDay();
  let day = new Date(firstDay);
  day.setDate(day.getDate() - startOffset);
  const weeks = [];
  while (day <= lastDay || weeks.length === 0) {
    const week = [];
    for (let d = 0; d < 7; d++) {
      if (day.getMonth() === month) {
        week.push({ date: day.getDate(), dayIndex: day.getDay(), inMonth: true });
      } else {
        week.push(null);
      }
      day.setDate(day.getDate() + 1);
    }
    weeks.push(week);
    if (day > lastDay) break;
  }
  return weeks;
}

async function loadLogo() {
  try {
    const res = await fetch('/api/logo');
    if (res.ok) {
      const blob = await res.blob();
      logoDataUrl = URL.createObjectURL(blob);
    }
  } catch(e) {}
}

async function loadHolidays() {
  const res = await fetch('/api/holidays');
  if (res.status === 401) { window.location.href = '/login'; return; }
  holidays = await res.json();
}

async function loadMenuItems() {
  const res = await fetch('/api/menu-items');
  if (res.status === 401) { window.location.href = '/login'; return; }
  menuItems = await res.json();
}

async function loadCalendarData(year, month) {
  const res = await fetch('/api/calendar/' + year + '/' + month);
  if (res.ok) return await res.json();
  return null;
}

async function saveCalendarData(year, month, data) {
  await fetch('/api/calendar/' + year + '/' + month, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  });
}

async function renderInputGrid() {
  const year = parseInt(yearSelect.value);
  const month = parseInt(monthSelect.value);
  const weeks = getWeekdaysInMonth(year, month);

  const saved = await loadCalendarData(year, month);
  const storedMap = {};
  if (saved) {
    if (saved.school_name) schoolInput.value = saved.school_name;
    try {
      const parsed = JSON.parse(saved.data_json);
      if (parsed.menus) {
        parsed.menus.forEach(m => { storedMap[m.date + '-' + m.day] = m; });
      }
    } catch(e) {}
  }

  inputGrid.innerHTML = '';

  weeks.forEach((week, wi) => {
    const validDays = week.filter(d => d !== null);
    if (validDays.length === 0) return;
    const group = document.createElement('div');
    group.className = 'week-group';
    const startDate = validDays[0].date;
    const endDate = validDays[validDays.length - 1].date;
    const h3 = document.createElement('h3');
    h3.textContent = 'Week ' + (wi+1) + ': ' + MONTHS[month] + ' ' + startDate + ' \u2013 ' + endDate;
    group.appendChild(h3);

    const grid = document.createElement('div');
    grid.className = 'week-grid';

    week.forEach((dayInfo, di) => {
      const cell = document.createElement('div');
      cell.className = 'day-cell';
      if (dayInfo === null) {
        cell.style.opacity = '0.3';
        cell.innerHTML = '<div class="day-label">' + DAY_NAMES[di] + '</div><p style="font-size:0.8rem;color:#999;">\u2014</p>';
      } else {
        const dayName = DAYS_FULL[dayInfo.dayOfWeek];
        const dow = dayInfo.dayOfWeek;
        const holiday = isHoliday(dayInfo.date, month);
        const stored = storedMap[dayInfo.date + '-' + dow];
        const isMessage = stored?.type === 'message' || (holiday && !stored);
        const messageText = stored?.message || (holiday?.name) || '';
        const breakfast = stored?.breakfast || '';
        const lunch = stored?.lunch || '';
        if (holiday) cell.classList.add('holiday-cell');

        cell.innerHTML =
          '<div class="day-label">' + dayName + ', ' + MONTHS[month] + ' ' + dayInfo.date +
          (holiday ? ' <span style="color:#e67e22;font-size:0.75rem;">(' + escapeHtml(holiday.name) + ')</span>' : '') + '</div>' +
          '<div class="toggle-row">' +
            '<input type="checkbox" id="toggle-' + dayInfo.date + '-' + dow + '" data-date="' + dayInfo.date + '" data-day="' + dow + '" class="day-mode-toggle"' + (isMessage ? ' checked' : '') + '>' +
            '<label for="toggle-' + dayInfo.date + '-' + dow + '">No meal (set message)</label>' +
          '</div>' +
          '<div class="menu-section" style="' + (isMessage ? 'display:none;' : '') + '">' +
            '<label>B: (Breakfast)</label>' +
            '<textarea data-day="' + dow + '" data-date="' + dayInfo.date + '" data-meal="breakfast" placeholder="Enter breakfast menu..." class="menu-textarea">' + escapeHtml(breakfast) + '</textarea>' +
            '<label>L: (Lunch)</label>' +
            '<textarea data-day="' + dow + '" data-date="' + dayInfo.date + '" data-meal="lunch" placeholder="Enter lunch menu..." class="menu-textarea">' + escapeHtml(lunch) + '</textarea>' +
          '</div>' +
          '<div class="message-section" style="' + (isMessage ? '' : 'display:none;') + '">' +
            '<label>Message</label>' +
            '<input type="text" class="message-input" data-date="' + dayInfo.date + '" data-day="' + dow + '" placeholder="e.g. Labor Day, No School" value="' + escapeHtml(messageText) + '">' +
          '</div>';
      }
      grid.appendChild(cell);
    });
    group.appendChild(grid);
    inputGrid.appendChild(group);
  });

  inputGrid.querySelectorAll('.day-mode-toggle').forEach(toggle => {
    toggle.addEventListener('change', function() {
      const cell = this.closest('.day-cell');
      const menuSection = cell.querySelector('.menu-section');
      const messageSection = cell.querySelector('.message-section');
      menuSection.style.display = this.checked ? 'none' : 'block';
      messageSection.style.display = this.checked ? 'block' : 'none';
    });
  });

  inputGrid.querySelectorAll('.menu-textarea').forEach(textarea => { setupAutocomplete(textarea); });
}

function setupAutocomplete(textarea) {
  let autocompleteDiv = null;
  const mealType = textarea.dataset.meal;
  textarea.addEventListener('input', function() {
    const value = this.value.toLowerCase();
    if (autocompleteDiv) { autocompleteDiv.remove(); autocompleteDiv = null; }
    if (value.length < 2) return;
    const items = getMenuItemNames(mealType);
    const matches = items.filter(item => item.toLowerCase().includes(value));
    if (matches.length === 0) return;
    autocompleteDiv = document.createElement('div');
    autocompleteDiv.className = 'autocomplete-list';
    autocompleteDiv.style.left = '10px';
    autocompleteDiv.style.top = (textarea.offsetTop + textarea.offsetHeight) + 'px';
    matches.slice(0, 5).forEach(match => {
      const item = document.createElement('div');
      item.className = 'autocomplete-item';
      item.textContent = match;
      item.addEventListener('mousedown', function(e) {
        e.preventDefault();
        textarea.value = match;
        autocompleteDiv.remove();
        autocompleteDiv = null;
      });
      autocompleteDiv.appendChild(item);
    });
    textarea.closest('.day-cell').appendChild(autocompleteDiv);
  });
  textarea.addEventListener('blur', function() {
    setTimeout(() => { if (autocompleteDiv) { autocompleteDiv.remove(); autocompleteDiv = null; } }, 200);
  });
}

function getMenusFromInputs() {
  const data = {};
  document.querySelectorAll('#input-grid .day-mode-toggle').forEach(toggle => {
    const key = toggle.dataset.date + '-' + toggle.dataset.day;
    if (!data[key]) data[key] = { date: toggle.dataset.date, day: parseInt(toggle.dataset.day), type: 'menu', breakfast: '', lunch: '', message: '' };
    data[key].type = toggle.checked ? 'message' : 'menu';
  });
  document.querySelectorAll('#input-grid textarea').forEach(ta => {
    const key = ta.dataset.date + '-' + ta.dataset.day;
    if (!data[key]) data[key] = { date: ta.dataset.date, day: parseInt(ta.dataset.day), type: 'menu', breakfast: '', lunch: '', message: '' };
    data[key][ta.dataset.meal] = ta.value.trim();
  });
  document.querySelectorAll('#input-grid .message-input').forEach(input => {
    const key = input.dataset.date + '-' + input.dataset.day;
    if (!data[key]) data[key] = { date: input.dataset.date, day: parseInt(input.dataset.day), type: 'menu', breakfast: '', lunch: '', message: '' };
    data[key].message = input.value.trim();
  });
  return Object.values(data);
}

function escapeHtml(text) {
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
}

function lightenColor(hex, factor) {
  const r = parseInt(hex.slice(1,3), 16);
  const g = parseInt(hex.slice(3,5), 16);
  const b = parseInt(hex.slice(5,7), 16);
  const lr = Math.round(r + (255 - r) * factor);
  const lg = Math.round(g + (255 - g) * factor);
  const lb = Math.round(b + (255 - b) * factor);
  return '#' + lr.toString(16).padStart(2,'0') + lg.toString(16).padStart(2,'0') + lb.toString(16).padStart(2,'0');
}

function getBackgroundStyle(style, color, lightColor) {
  switch(style) {
    case 'gradient-diagonal': return 'background: linear-gradient(135deg, ' + lightColor + ' 0%, ' + lightenColor(color, 0.7) + ' 100%);';
    case 'gradient-vertical': return 'background: linear-gradient(to bottom, ' + lightColor + ' 0%, ' + lightenColor(color, 0.75) + ' 100%);';
    case 'gradient-radial': return 'background: radial-gradient(circle, ' + lightColor + ' 0%, ' + lightenColor(color, 0.7) + ' 100%);';
    case 'pattern-dots': return 'background: ' + lightColor + '; background-image: radial-gradient(circle, ' + lightenColor(color, 0.6) + ' 1px, transparent 1px); background-size: 10px 10px;';
    case 'pattern-lines': return 'background: repeating-linear-gradient(0deg, ' + lightColor + ', ' + lightColor + ' 3px, ' + lightenColor(color, 0.75) + ' 3px, ' + lightenColor(color, 0.75) + ' 4px);';
    case 'pattern-grid': return 'background: ' + lightColor + '; background-image: linear-gradient(' + lightenColor(color, 0.75) + ' 1px, transparent 1px), linear-gradient(90deg, ' + lightenColor(color, 0.75) + ' 1px, transparent 1px); background-size: 12px 12px;';
    default: return 'background: ' + lightColor + ';';
  }
}

function getBorderStyle(style, color) {
  switch(style) {
    case 'rounded': return 'border: 2px solid ' + lightenColor(color, 0.4) + '; border-radius: 12px;';
    case 'dashed': return 'border: 2px dashed ' + lightenColor(color, 0.4) + '; border-radius: 6px;';
    case 'double': return 'border: 4px double ' + lightenColor(color, 0.4) + '; border-radius: 6px;';
    case 'shadow': return 'border: 1px solid ' + lightenColor(color, 0.5) + '; border-radius: 6px; box-shadow: 2px 2px 6px rgba(0,0,0,0.15);';
    default: return 'border: 2px solid ' + lightenColor(color, 0.4) + '; border-radius: 6px;';
  }
}

function getBorderStyleInline(style, color) {
  switch(style) {
    case 'rounded': return 'border:2px solid ' + color + ';border-radius:8px;';
    case 'dashed': return 'border:2px dashed ' + color + ';';
    case 'double': return 'border:3px double ' + color + ';';
    case 'shadow': return 'border:1px solid ' + color + ';box-shadow:1px 1px 4px rgba(0,0,0,0.12);';
    default: return 'border:1px solid ' + color + ';';
  }
}

function mealBadge(letter, color) {
  return '<span style="display:inline-block;background:' + color + ';color:#fff;font-weight:700;font-size:8pt;line-height:1;border-radius:3px;padding:2px 4px;margin-right:3px;vertical-align:baseline;">' + letter + '</span>';
}

function buildCalendarHTML() {
  const year = parseInt(yearSelect.value);
  const month = parseInt(monthSelect.value);
  const school = schoolInput.value.trim();
  const footerMessage = footerMessageInput.value.trim();
  const preset = getActivePreset();
  const themeColor = colorPicker.value;
  const font = fontSelect.value;
  const bgStyleValue = bgStyle.value;
  const lightColor = lightenColor(themeColor, 0.85);
  const borderColor = lightenColor(themeColor, 0.4);
  const weeks = getFullMonthWeeks(year, month);
  const menus = getMenusFromInputs();

  const menuMap = {};
  menus.forEach(m => { menuMap[m.date + '-' + m.day] = m; });

  let html = '<div id="pdf-page" style="width:1056px;max-width:100%;padding:16px;font-family:' + font + ';background:#fff;color:#222;display:flex;flex-direction:column;">';

  html += '<div class="cal-header" style="display:flex;align-items:center;justify-content:center;gap:14px;margin-bottom:10px;flex-shrink:0;">';
  if (preset.headerSvg) {
    html += presetSvgHtml(preset.headerSvg, themeColor);
  }
  if (logoDataUrl) {
    html += '<img src="' + logoDataUrl + '" style="max-height:48px;max-width:96px;object-fit:contain;flex-shrink:0;">';
  }
  html += '<div style="text-align:center;">';
  if (school) html += '<div style="font-size:13pt;font-weight:700;color:' + themeColor + ';line-height:1.2;">' + escapeHtml(school) + '</div>';
  html += '<div style="font-size:20pt;font-weight:700;color:' + themeColor + ';line-height:1.2;">' + MONTHS[month] + ' ' + year + '</div>';
  html += '</div>';
  if (preset.headerSvg) {
    html += presetSvgHtml(preset.headerSvg, themeColor);
  }
  html += '</div>';

  html += '<div class="cal-grid" style="display:grid;grid-template-columns:repeat(7,1fr);gap:4px;flex:1;">';

  DAYS_SHORT.forEach((d, i) => {
    const isWeekend = i === 0 || i === 6;
    const bg = isWeekend ? lightenColor(themeColor, 0.35) : themeColor;
    html += '<div style="text-align:center;font-weight:700;font-size:10pt;color:#fff;background:' + bg + ';padding:5px 0;border-radius:3px;">' + d + '</div>';
  });

  weeks.forEach(week => {
    week.forEach((dayInfo, di) => {
      const isWeekend = di === 0 || di === 6;

      if (dayInfo === null || !dayInfo.inMonth) {
        let s = 'border:1px solid #e0e0e0;border-radius:' + preset.cellBorderRadius + ';padding:4px;opacity:0.3;background:#f9f9f9;';
        html += '<div class="cal-cell" style="' + s + '"></div>';
      } else if (isWeekend) {
        let s = 'border:1px solid ' + borderColor + ';border-radius:' + preset.cellBorderRadius + ';padding:4px;background:#f5f5f5;';
        html += '<div class="cal-cell" style="' + s + '"><div style="font-weight:600;font-size:11px;color:#888;">' + dayInfo.date + '</div></div>';
      } else {
        const key = dayInfo.date + '-' + di;
        const entry = menuMap[key] || { type: 'menu', breakfast: '', lunch: '', message: '' };
        let s;
        if (entry.type === 'message') {
          s = 'border:1px solid ' + themeColor + ';border-radius:' + preset.cellBorderRadius + ';padding:5px;background:#fef9e7;';
        } else {
          s = 'border:1px solid ' + borderColor + ';border-radius:' + preset.cellBorderRadius + ';padding:5px;' + getBackgroundStyle(bgStyleValue, themeColor, lightColor);
        }
        let content = '<div style="font-weight:700;font-size:12px;color:' + themeColor + ';margin-bottom:2px;">' + DAYS_SHORT[di] + ' ' + dayInfo.date + '</div>';
        if (entry.type === 'message') {
          if (entry.message) content += '<div style="font-size:10px;line-height:1.3;text-align:center;font-weight:600;color:' + themeColor + ';">' + escapeHtml(entry.message).replace(/\n/g, '<br>') + '</div>';
        } else {
          if (entry.breakfast) content += '<div style="font-size:10px;line-height:1.3;margin-bottom:2px;display:-webkit-box;-webkit-box-orient:vertical;-webkit-line-clamp:2;overflow:hidden;overflow-wrap:anywhere;">' + presetBadgeStyle(preset, themeColor) + escapeHtml(entry.breakfast).replace(/\n/g, '<br>') + '</div>';
          if (entry.lunch) content += '<div style="font-size:10px;line-height:1.3;display:-webkit-box;-webkit-box-orient:vertical;-webkit-line-clamp:2;overflow:hidden;overflow-wrap:anywhere;">' + presetBadgeStyle(preset, lightenColor(themeColor, 0.3)) + escapeHtml(entry.lunch).replace(/\n/g, '<br>') + '</div>';
        }
        html += '<div class="cal-cell" style="' + s + '">' + content + '</div>';
      }
    });
  });

  html += '</div>';

  if (footerMessage || preset.footerSvg) {
    html += '<div class="cal-footer" style="margin-top:8px;padding-top:6px;border-top:1px solid ' + borderColor + ';text-align:center;font-size:9pt;color:#666;flex-shrink:0;display:flex;align-items:center;justify-content:center;gap:12px;">';
    if (preset.footerSvg) {
      html += presetSvgHtml(preset.footerSvg, themeColor);
    }
    if (footerMessage) {
      html += '<span>' + escapeHtml(footerMessage) + '</span>';
    }
    if (preset.footerSvg) {
      html += presetSvgHtml(preset.footerSvg, themeColor);
    }
    html += '</div>';
  }

  html += '</div>';

  return html;
}

function renderCalendar(force) {
  if (!force && calendarOutput.style.display !== 'block') return;
  calendarOutput.innerHTML = buildCalendarHTML();
  calendarOutput.style.display = 'block';
  downloadSection.style.display = 'block';
}

async function generateCalendar() {
  const year = parseInt(yearSelect.value);
  const month = parseInt(monthSelect.value);
  const school = schoolInput.value.trim();
  const menus = getMenusFromInputs();

  await saveCalendarData(year, month, { school_name: school, menus: menus });
  renderCalendar(true);
  clearDirty();
}

downloadBtn.onclick = function() {
  const page = document.getElementById('pdf-page');
  if (!page) { showToast('Generate a calendar first', 'error'); return; }
  const font = fontSelect.value;
  const grid = page.querySelector('.cal-grid');
  const totalCells = grid.children.length;
  const weekRows = Math.floor((totalCells - 7) / 7);

  const PAGE_W = 10.6;
  const PAGE_H = 8.1;
  const HEADER_H = 0.7;
  const DAY_HDR_H = 0.3;
  const GAP = 0.04;
  const hasFooter = !!page.querySelector('.cal-footer');
  const FOOTER_H = hasFooter ? 0.3 : 0;
  const gridH = PAGE_H - HEADER_H - GAP - DAY_HDR_H - (hasFooter ? FOOTER_H + GAP : 0);
  const rowH = (gridH - (weekRows - 1) * GAP) / weekRows;
  const W = (PAGE_W - 6 * GAP) / (2 + 5 * 1.33);
  const weekendW = W;
  const weekdayW = W * 1.33;

  const colTemplate = [weekendW, weekdayW, weekdayW, weekdayW, weekdayW, weekdayW, weekendW].map(w => w.toFixed(3) + 'in').join(' ');
  const rowTemplate = Array(weekRows).fill(rowH.toFixed(3) + 'in').join(' ');

  const printWindow = window.open('', '_blank');
  printWindow.document.write(
    '<!DOCTYPE html><html><head><title>Calendar</title>' +
    '<link rel="stylesheet" href="' + window.location.origin + '/fonts/fonts.css">' +
    '<style>' +
    '@page { size: letter landscape; margin: 0.2in; }' +
    '* { margin: 0; padding: 0; box-sizing: border-box; }' +
    'html, body { width: 100%; height: 100%; overflow: hidden; }' +
    'body { font-family: ' + font + '; background: #fff; color: #222; -webkit-print-color-adjust: exact; print-color-adjust: exact; }' +
    '#pdf-page { width: ' + PAGE_W + 'in; height: ' + PAGE_H + 'in; display: flex; flex-direction: column; padding: 0 !important; }' +
    '.cal-header { height: ' + HEADER_H + 'in; margin: 0 0 ' + GAP + 'in 0 !important; flex-shrink: 0; }' +
    '.cal-header img { max-height: 0.5in !important; }' +
    '.cal-grid { display: grid !important; grid-template-columns: ' + colTemplate + '; grid-template-rows: ' + DAY_HDR_H + 'in ' + rowTemplate + '; gap: ' + GAP + 'in; flex: 1 1 auto; min-height: 0; }' +
    '.cal-cell { overflow: hidden; }' +
    '.cal-footer { height: ' + FOOTER_H + 'in; margin: ' + GAP + 'in 0 0 0 !important; padding: 0 !important; display: flex !important; align-items: center; justify-content: center; overflow: hidden; }' +
    '</style>' +
    '</head><body>' +
    page.innerHTML +
    '</body></html>'
  );
  printWindow.document.close();
};

exportBtn.addEventListener('click', async function() {
  const res = await fetch('/api/export');
  const blob = await res.blob();
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = 'calendar-export.json';
  a.click();
  URL.revokeObjectURL(url);
});

function parseCSV(text) {
  const rows = [];
  let row = [], value = '', quoted = false;
  for (let i = 0; i < text.length; i++) {
    const char = text[i];
    if (char === '"' && quoted && text[i + 1] === '"') { value += '"'; i++; }
    else if (char === '"') quoted = !quoted;
    else if (char === ',' && !quoted) { row.push(value.trim()); value = ''; }
    else if ((char === '\n' || char === '\r') && !quoted) {
      if (char === '\r' && text[i + 1] === '\n') i++;
      row.push(value.trim());
      if (row.some(cell => cell !== '')) rows.push(row);
      row = []; value = '';
    } else value += char;
  }
  if (value || row.length) {
    row.push(value.trim());
    if (row.some(cell => cell !== '')) rows.push(row);
  }
  if (rows.length < 2) return [];
  const headers = rows.shift().map(h => h.replace(/^\uFEFF/, '').toLowerCase().trim());
  return rows.map(values => Object.fromEntries(headers.map((header, i) => [header, values[i] || ''])));
}

function csvValue(row, ...names) {
  for (const name of names) if (row[name] !== undefined) return row[name];
  return '';
}

function parseMonth(value, csv) {
  const text = String(value).trim();
  const named = MONTHS.findIndex(month => month.toLowerCase() === text.toLowerCase());
  if (named >= 0) return named;
  const number = Number(text);
  if (!Number.isInteger(number)) throw new Error('month must be a month name or number');
  const month = csv ? (number === 0 ? 0 : number - 1) : number;
  if (month < 0 || month > 11) throw new Error('month must be January-December or 0-11 in JSON');
  return month;
}

function parseImportData(text, fileName, kind) {
  const isCSV = fileName.toLowerCase().endsWith('.csv');
  const source = isCSV ? parseCSV(text) : JSON.parse(text);
  if (kind === 'holidays') {
    const rows = isCSV ? source : (Array.isArray(source) ? source : source.holidays);
    if (!Array.isArray(rows)) throw new Error('JSON must contain holidays array');
    return rows.map(row => ({
      month: parseMonth(row.month, isCSV),
      date: Number(csvValue(row, 'date')),
      name: String(csvValue(row, 'name', 'holiday')).trim()
    }));
  }
  if (kind === 'menu-items') {
    let rows;
    if (isCSV) rows = source;
    else if (Array.isArray(source)) rows = source;
    else if (Array.isArray(source.items)) rows = source.items;
    else {
      const groups = source.menu_items || source;
      rows = ['breakfast', 'lunch'].flatMap(type => (groups[type] || []).map(item => typeof item === 'string' ? { meal_type: type, item_name: item } : { ...item, meal_type: type }));
    }
    return rows.map(row => ({
      meal_type: String(csvValue(row, 'meal_type', 'type')).toLowerCase().trim(),
      item_name: String(csvValue(row, 'item_name', 'name')).trim()
    }));
  }
  const rows = isCSV ? source : (Array.isArray(source) ? source : source.calendar_data);
  if (!Array.isArray(rows)) throw new Error('JSON must contain calendar_data array');
  return rows.map(row => {
    const dataJSON = csvValue(row, 'data_json');
    let data = dataJSON;
    if (!data) {
      const menus = typeof row.menus === 'string' ? JSON.parse(row.menus) : (row.menus || []);
      data = JSON.stringify({ menus: menus });
    }
    else JSON.parse(data);
    return {
      year: Number(csvValue(row, 'year')),
      month: parseMonth(row.month, isCSV),
      school_name: String(csvValue(row, 'school_name', 'school')).trim(),
      data_json: data
    };
  });
}

async function importFile(file, kind, input) {
  if (!file) return;
  try {
    const records = parseImportData(await file.text(), file.name, kind);
    const endpoint = kind === 'holidays' ? '/api/import/holidays' : kind === 'menu-items' ? '/api/import/menu-items' : '/api/import/calendar';
    const key = kind === 'holidays' ? 'holidays' : kind === 'menu-items' ? 'items' : 'calendar_data';
    const res = await fetch(endpoint, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ [key]: records })
    });
    if (!res.ok) {
      const error = await res.json().catch(() => ({}));
      throw new Error(error.error || 'server rejected import');
    }
    showToast(kind + ' imported successfully');
    if (kind === 'holidays' || kind === 'menu-items') await Promise.all([loadHolidays(), loadMenuItems()]);
    await renderInputGrid();
    renderCalendar(true);
  } catch (error) {
    showToast('Import failed: ' + error.message, 'error');
  }
  input.value = '';
}

importHolidaysFile.addEventListener('change', e => importFile(e.target.files[0], 'holidays', importHolidaysFile));
importMenuItemsFile.addEventListener('change', e => importFile(e.target.files[0], 'menu-items', importMenuItemsFile));
importCalendarFile.addEventListener('change', e => importFile(e.target.files[0], 'calendar', importCalendarFile));

function shiftMonth(delta) {
  let m = parseInt(monthSelect.value) + delta;
  let y = parseInt(yearSelect.value);
  if (m < 0) { m = 11; y--; } else if (m > 11) { m = 0; y++; }
  const years = Array.from(yearSelect.options).map(o => parseInt(o.value));
  const minYear = Math.min(...years);
  const maxYear = Math.max(...years);
  if (y < minYear || y > maxYear) return;
  monthSelect.value = m;
  yearSelect.value = y;
  renderInputGrid().then(() => renderCalendar(true));
}

let stylePreviewTimer = null;
function scheduleStylePreview() {
  clearTimeout(stylePreviewTimer);
  stylePreviewTimer = setTimeout(() => renderCalendar(false), 150);
}

(async function init() {
  try {
    await Promise.all([loadHolidays(), loadMenuItems(), loadLogo()]);
  } catch(e) {
    window.location.href = '/login';
    return;
  }
  initSelectors();
  initPresets();
  await renderInputGrid();

  monthSelect.addEventListener('change', () => renderInputGrid().then(() => renderCalendar(true)));
  yearSelect.addEventListener('change', () => renderInputGrid().then(() => renderCalendar(true)));
  document.getElementById('prev-month-btn').addEventListener('click', () => shiftMonth(-1));
  document.getElementById('next-month-btn').addEventListener('click', () => shiftMonth(1));
  generateBtn.addEventListener('click', generateCalendar);

  [colorPicker, fontSelect, bgStyle, borderStyle].forEach(el => {
    el.addEventListener('input', scheduleStylePreview);
    el.addEventListener('change', scheduleStylePreview);
  });
  schoolInput.addEventListener('input', scheduleStylePreview);
  footerMessageInput.addEventListener('input', scheduleStylePreview);
})();
