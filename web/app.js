const S = {
    token: localStorage.getItem('tok'),
    email: localStorage.getItem('em'),
    page: 'list',
    pid: null,
    off: 0,
    lim: 10,
    uid: localStorage.getItem('uid'),
};

async function api(path, o = {}) {
    const h = { 'Content-Type': 'application/json' };
    if (S.token) h['Authorization'] = 'Bearer ' + S.token;
    const r = await fetch('/api' + path, { ...o, headers: h });
    const t = await r.text();
    const d = t ? JSON.parse(t) : null;
    if (!r.ok) throw { status: r.status, msg: d?.error || 'Error' };
    return d;
}

function parseJWT(tok) {
    try {
        return JSON.parse(atob(tok.split('.')[1]));
    } catch { return null; }
}

function setAuth(tok, em, uid) {
    S.token = tok; S.email = em; S.uid = uid;
    localStorage.setItem('tok', tok);
    localStorage.setItem('em', em);
    localStorage.setItem('uid', uid);
}

function clearAuth() {
    S.token = null; S.email = null; S.uid = null;
    localStorage.removeItem('tok');
    localStorage.removeItem('em');
    localStorage.removeItem('uid');
}

let tt;
function toast(m) {
    const e = document.getElementById('toast');
    e.textContent = m; e.className = 'show';
    clearTimeout(tt); tt = setTimeout(() => e.className = '', 2500);
}

function nav(page, pid) {
    S.page = page; S.pid = pid || null;
    if (page === 'list') S.off = 0;
    render();
}

function esc(s) {
    return String(s)
        .replace(/&/g,'&amp;').replace(/</g,'&lt;')
        .replace(/>/g,'&gt;').replace(/"/g,'&quot;');
}
function date(s) {
    return new Date(s).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
}

function renderNav() {
    document.getElementById('nav').innerHTML = S.token
        ? `<span class="user-pill">${esc(S.email)}</span>
       <button class="btn" onclick="nav('create')">+ Write</button>
       <button class="btn ghost" onclick="logout()">Logout</button>`
        : `<button class="btn ghost" onclick="nav('login')">Login</button>
       <button class="btn primary" onclick="nav('register')">Register</button>`;
}

function render() {
    renderNav();
    const el = document.getElementById('page');
    ({
        list:     () => renderList(el),
        post:     () => renderPost(el),
        login:    () => renderLogin(el),
        register: () => renderReg(el),
        create:   () => renderEditor(el, null),
        edit:     () => loadEditor(el),
    }[S.page] || (() => {}))();
}

async function renderList(el) {
    el.innerHTML = '<div class="spin"></div>';
    try {
        const posts = await api(`/posts?limit=${S.lim}&offset=${S.off}`) || [];
        const hasPrev = S.off > 0, hasNext = posts.length === S.lim;
        const pg = Math.floor(S.off / S.lim) + 1;
        el.innerHTML = `
      <div class="page-head fade">
        <h1 class="page-h1">Posts</h1>
      </div>
      <div class="post-list fade">
        ${posts.length === 0
                ? '<div class="empty">No posts yet.</div>'
                : posts.map(p => `
            <div class="post-row" onclick="nav('post','${p.id}')">
              <div class="post-row-meta">${date(p.created_at)} · ${esc(p.author_email.split('@')[0])}</div>
              <div class="post-row-title">${esc(p.title)}</div>
              <div class="post-row-excerpt">${esc(p.body)}</div>
            </div>`).join('')}
      </div>
      ${hasPrev || hasNext ? `
        <div class="pager">
          <button class="btn" onclick="prevPg()" ${!hasPrev ? 'disabled' : ''}>← Prev</button>
          <span>Page ${pg}</span>
          <button class="btn" onclick="nextPg()" ${!hasNext ? 'disabled' : ''}>Next →</button>
        </div>` : ''}`;
    } catch(e) {
        el.innerHTML = `<div class="alert error">${e.msg}</div>`;
    }
}

function nextPg() { S.off += S.lim; renderList(document.getElementById('page')); }
function prevPg() { S.off = Math.max(0, S.off - S.lim); renderList(document.getElementById('page')); }

async function renderPost(el) {
    el.innerHTML = '<div class="spin"></div>';
    try {
        const p = await api(`/posts/${S.pid}`);
        el.innerHTML = `
      <div class="fade">
        <button class="back-btn" onclick="nav('list')">← All posts</button>
        <div class="post-meta"><span>${date(p.created_at)} · ${esc(p.author_email.split('@')[0])}</span></div>
        <h1 class="post-title">${esc(p.title)}</h1>
        <div class="post-body">${esc(p.body)}</div>
        ${S.token && S.uid === p.author_id ? `
          <div class="post-footer">
            <button class="btn" onclick="nav('edit','${p.id}')">Edit</button>
            <button class="btn danger" onclick="confirmDel('${p.id}','${esc(p.title)}')">Delete</button>
          </div>` : ''}
      </div>`;
    } catch(e) {
        el.innerHTML = `<div class="alert error">${e.msg}</div>`;
    }
}

function confirmDel(id, title) {
    const m = document.createElement('div');
    m.className = 'modal-bg';
    m.innerHTML = `
    <div class="modal">
      <h3>Delete post?</h3>
      <p>"${esc(title)}" will be removed permanently.</p>
      <div class="modal-btns">
        <button class="btn ghost" onclick="this.closest('.modal-bg').remove()">Cancel</button>
        <button class="btn danger" id="del-ok">Delete</button>
      </div>
    </div>`;
    document.body.appendChild(m);
    document.getElementById('del-ok').onclick = async () => {
        document.getElementById('del-ok').disabled = true;
        try {
            await api(`/posts/${id}`, { method: 'DELETE' });
            m.remove(); toast('Deleted.'); nav('list');
        } catch(e) {
            m.remove(); toast(e.msg);
        }
    };
}

function renderEditor(el, post) {
    const edit = !!post;
    el.innerHTML = `
    <div class="fade">
      <div class="editor-head">
        <span class="editor-label">${edit ? 'Edit' : 'New post'}</span>
        <div class="editor-actions">
          <button class="btn ghost" onclick="nav(${edit ? `'post','${post?.id}'` : "'list'"})">Cancel</button>
          <button class="btn primary" id="save-btn">Publish</button>
        </div>
      </div>
      <div id="ed-err"></div>
      <div class="field">
        <label>Title</label>
        <input class="title-input" id="ed-title" type="text" placeholder="Post title…" value="${edit ? esc(post.title) : ''}"/>
      </div>
      <div class="field">
        <label>Body</label>
        <textarea id="ed-body" placeholder="Write something…">${edit ? esc(post.body) : ''}</textarea>
      </div>
    </div>`;

    document.getElementById('save-btn').onclick = async () => {
        const title = document.getElementById('ed-title').value.trim();
        const body  = document.getElementById('ed-body').value.trim();
        const err   = document.getElementById('ed-err');
        const btn   = document.getElementById('save-btn');
        if (!title || !body) {
            err.innerHTML = '<div class="alert error">Title and body are required.</div>';
            return;
        }
        btn.disabled = true; btn.textContent = '…'; err.innerHTML = '';
        try {
            if (edit) {
                await api(`/posts/${post.id}`, { method: 'PUT', body: JSON.stringify({ title, body }) });
                toast('Updated.'); nav('post', post.id);
            } else {
                const p = await api('/posts', { method: 'POST', body: JSON.stringify({ title, body }) });
                toast('Published.'); nav('post', p.id);
            }
        } catch(e) {
            err.innerHTML = `<div class="alert error">${e.msg}</div>`;
            btn.disabled = false; btn.textContent = 'Publish';
        }
    };
}

async function loadEditor(el) {
    el.innerHTML = '<div class="spin"></div>';
    try {
        renderEditor(el, await api(`/posts/${S.pid}`));
    } catch(e) {
        el.innerHTML = `<div class="alert error">${e.msg}</div>`;
    }
}

function renderLogin(el) {
    el.innerHTML = `
    <div class="form-wrap fade">
      <h1 class="form-title">Welcome back.</h1>
      <p class="form-hint">No account? <a onclick="nav('register')">Register</a></p>
      <div id="l-err"></div>
      <div class="field"><label>Email</label><input id="l-em" type="email" autocomplete="email" placeholder="you@example.com"/></div>
      <div class="field"><label>Password</label><input id="l-pw" type="password" autocomplete="current-password" placeholder="••••••••"/></div>
      <button class="btn primary btn-full" id="l-btn">Login</button>
    </div>`;

    const go = async () => {
        const email    = document.getElementById('l-em').value.trim();
        const password = document.getElementById('l-pw').value;
        const err = document.getElementById('l-err');
        const btn = document.getElementById('l-btn');
        btn.disabled = true; btn.textContent = '…'; err.innerHTML = '';
        try {
            const d = await api('/auth/login', { method: 'POST', body: JSON.stringify({ email, password }) });
            const claims = parseJWT(d.token);
            setAuth(d.token, email, claims?.sub); toast('Welcome back.'); nav('list');
        } catch(e) {
            err.innerHTML = `<div class="alert error">${e.msg}</div>`;
            btn.disabled = false; btn.textContent = 'Login';
        }
    };
    document.getElementById('l-btn').onclick = go;
    document.getElementById('l-pw').addEventListener('keydown', e => e.key === 'Enter' && go());
}

function renderReg(el) {
    el.innerHTML = `
    <div class="form-wrap fade">
      <h1 class="form-title">Create account.</h1>
      <p class="form-hint">Have an account? <a onclick="nav('login')">Login</a></p>
      <div id="r-err"></div>
      <div class="field"><label>Email</label><input id="r-em" type="email" autocomplete="email" placeholder="you@example.com"/></div>
      <div class="field"><label>Password <span style="font-size:10px;color:var(--muted)">(min 8 chars)</span></label><input id="r-pw" type="password" autocomplete="new-password" placeholder="••••••••"/></div>
      <button class="btn primary btn-full" id="r-btn">Create account</button>
    </div>`;

    const go = async () => {
        const email    = document.getElementById('r-em').value.trim();
        const password = document.getElementById('r-pw').value;
        const err = document.getElementById('r-err');
        const btn = document.getElementById('r-btn');
        btn.disabled = true; btn.textContent = '…'; err.innerHTML = '';
        try {
            await api('/auth/register', { method: 'POST', body: JSON.stringify({ email, password }) });
            toast('Account created. Please login.'); nav('login');
        } catch(e) {
            err.innerHTML = `<div class="alert error">${e.msg}</div>`;
            btn.disabled = false; btn.textContent = 'Create account';
        }
    };
    document.getElementById('r-btn').onclick = go;
    document.getElementById('r-pw').addEventListener('keydown', e => e.key === 'Enter' && go());
}


function logout() { clearAuth(); toast('Logged out.'); nav('list'); }

render();
