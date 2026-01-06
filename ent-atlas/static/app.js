// DOM Elements
const createSection = document.getElementById('create-section');
const viewSection = document.getElementById('view-section');
const recentSection = document.getElementById('recent-section');
const pasteForm = document.getElementById('paste-form');
const recentList = document.getElementById('recent-list');

// Check if viewing a specific paste
const path = window.location.pathname;
const slugMatch = path.match(/^\/p\/([a-zA-Z0-9_-]+)$/);

if (slugMatch) {
    loadPaste(slugMatch[1]);
} else {
    loadRecentPastes();
}

// Form submission
pasteForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const formData = {
        title: document.getElementById('title').value,
        content: document.getElementById('content').value,
        language: document.getElementById('language').value,
        is_public: document.getElementById('is_public').checked,
        expires_in_minutes: parseInt(document.getElementById('expires').value) || null
    };

    if (!formData.content.trim()) {
        alert('Content is required');
        return;
    }

    try {
        const response = await fetch('/pastes', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Failed to create paste');
        }

        const paste = await response.json();
        window.location.href = `/p/${paste.slug}`;
    } catch (error) {
        alert('Error: ' + error.message);
    }
});

// Load a specific paste
async function loadPaste(slug) {
    createSection.classList.add('hidden');
    recentSection.classList.add('hidden');
    viewSection.classList.remove('hidden');

    try {
        const response = await fetch(`/pastes/${slug}`);
        
        if (!response.ok) {
            if (response.status === 410) {
                throw new Error('This paste has expired');
            }
            if (response.status === 404) {
                throw new Error('Paste not found');
            }
            throw new Error('Failed to load paste');
        }

        const paste = await response.json();
        displayPaste(paste);
    } catch (error) {
        viewSection.innerHTML = `
            <div class="error">
                <h2>Error</h2>
                <p>${error.message}</p>
                <a href="/" class="btn" style="margin-top: 15px;">Create New Paste</a>
            </div>
        `;
    }
}

// Display paste content
function displayPaste(paste) {
    document.getElementById('paste-title').textContent = paste.title || 'Untitled';
    document.getElementById('paste-content').textContent = paste.content;
    document.getElementById('paste-language').textContent = paste.language;
    document.getElementById('paste-date').textContent = formatDate(paste.created_at);

    // Copy button
    document.getElementById('copy-btn').addEventListener('click', () => {
        navigator.clipboard.writeText(paste.content).then(() => {
            const btn = document.getElementById('copy-btn');
            btn.textContent = 'Copied!';
            setTimeout(() => btn.textContent = 'Copy', 2000);
        });
    });

    // Raw button
    document.getElementById('raw-btn').addEventListener('click', () => {
        const blob = new Blob([paste.content], { type: 'text/plain' });
        const url = URL.createObjectURL(blob);
        window.open(url, '_blank');
    });

    // Update page title
    document.title = `${paste.title || 'Untitled'} - Pastebin`;
}

// Load recent pastes
async function loadRecentPastes() {
    try {
        const response = await fetch('/pastes?limit=10');
        
        if (!response.ok) {
            throw new Error('Failed to load pastes');
        }

        const data = await response.json();
        displayRecentPastes(data.pastes || []);
    } catch (error) {
        recentList.innerHTML = `<li class="error">Failed to load recent pastes</li>`;
    }
}

// Display recent pastes list
function displayRecentPastes(pastes) {
    if (pastes.length === 0) {
        recentList.innerHTML = '<li class="loading">No pastes yet</li>';
        return;
    }

    recentList.innerHTML = pastes.map(paste => `
        <li>
            <a href="/p/${paste.slug}">
                ${escapeHtml(paste.title || 'Untitled')}
            </a>
            <span class="meta">
                <span class="badge">${paste.language}</span>
                ${formatDate(paste.created_at)}
            </span>
        </li>
    `).join('');
}

// Format date
function formatDate(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diff = now - date;

    if (diff < 60000) return 'just now';
    if (diff < 3600000) return `${Math.floor(diff / 60000)}m ago`;
    if (diff < 86400000) return `${Math.floor(diff / 3600000)}h ago`;
    if (diff < 604800000) return `${Math.floor(diff / 86400000)}d ago`;

    return date.toLocaleDateString();
}

// Escape HTML to prevent XSS
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}
