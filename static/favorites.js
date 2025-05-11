// API endpoint
const API_URL = 'http://localhost:8080/api';

// DOM Elements
const favoritesGrid = document.getElementById('favoritesGrid');

// Helper to convert ["U+1F1E6", "U+1F1EB"] to the actual emoji
function unicodeArrayToEmoji(unicodeArr) {
    // Debug log
    console.log('Flag unicode array:', unicodeArr);
    return String.fromCodePoint(...unicodeArr.map(u => parseInt(u.replace('U+', ''), 16)));
}

// Load favorites
async function loadFavorites() {
    try {
        const response = await fetch(`${API_URL}/favorites`);
        if (!response.ok) throw new Error('Failed to load favorites');
        const favorites = await response.json();
        renderFavorites(favorites);
    } catch (error) {
        console.error('Error loading favorites:', error);
        favoritesGrid.innerHTML = `
            <div class="empty-state">
                <p>You haven't added any emojis to your favorites yet.</p>
                <a href="/catalog.html" class="cta-button">Browse Emojis</a>
            </div>
        `;
    }
}

// Remove from favorites
async function removeFavorite(id) {
    try {
        const response = await fetch(`${API_URL}/favorites/${id}`, { method: 'DELETE' });
        if (!response.ok) throw new Error('Failed to remove favorite');
        loadFavorites(); // Reload the list
    } catch (error) {
        console.error('Error removing favorite:', error);
    }
}

// Render favorites
function renderFavorites(favorites) {
    if (favorites.length === 0) {
        favoritesGrid.innerHTML = `
            <div class="empty-state">
                <p>You haven't added any emojis to your favorites yet.</p>
                <a href="/catalog.html" class="cta-button">Browse Emojis</a>
            </div>
        `;
        return;
    }

    favoritesGrid.innerHTML = favorites.map(emoji => {
        const emojiChar = emoji.category === 'flags'
            ? unicodeArrayToEmoji(emoji.unicode)
            : emoji.htmlCode[0];
        return `
            <div class="emoji-card">
                <div class="emoji-char">${emojiChar}</div>
                <div class="emoji-name">${emoji.name}</div>
                <div class="emoji-category">${emoji.category}</div>
                <button 
                    class="favorite-button active"
                    onclick="removeFavorite('${emoji.id}')"
                >
                    ❤️
                </button>
            </div>
        `;
    }).join('');
}

// Initial load
loadFavorites(); 