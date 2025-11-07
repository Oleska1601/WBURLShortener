const API_BASE = 'http://localhost:8082';

// ==================== СОЗДАНИЕ ССЫЛКИ ====================
async function createShortUrl() {
    const url = document.getElementById('urlInput').value.trim();
    const custom = document.getElementById('customInput').value.trim();
    const result = document.getElementById('createResult');

    if (!url) {
        showResult(result, '❌ Введите URL', 'error');
        return;
    }

    try {
        const payload = { url: url };
        if (custom) payload.short_url = custom;

        const response = await fetch(`${API_BASE}/api/shorten`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        const data = await response.json();

        if (response.status === 201) {
            const html = `
                <div class="success">
                    <strong>✅ Ссылка успешно создана!</strong>
                    <div class="url">${data.short_url}</div>
                    <div class="action-buttons">
                        <button onclick="copyToClipboard('${data.short_url}')">📋 Копировать</button>
                        <button onclick="testRedirect('${data.short_url}')">🔗 Тест перехода</button>
                        <button onclick="getAnalytics('${data.short_url}')">📊 Аналитика</button>
                    </div>
                    <div><small>ID: ${data.id}</small></div>
                </div>
            `;
            showResult(result, html, 'success');
            document.getElementById('urlInput').value = '';
            document.getElementById('customInput').value = '';
        } else {
            let errorMessage = '❌ Неизвестная ошибка';
            if (response.status === 400) errorMessage = '❌ Невозможно создать ссылку';
            else if (response.status === 409) errorMessage = '❌ Ссылка уже существует';
            else if (data.error) errorMessage = `❌ ${data.error}`;
            showResult(result, errorMessage, 'error');
        }
    } catch (error) {
        showResult(result, `❌ Ошибка сети: ${error.message}`, 'error');
    }
}

// ==================== АНАЛИТИКА ====================
async function getAnalytics(shortUrlParam = null) {
    const short_url = shortUrlParam || document.getElementById('analyticsInput').value.trim();
    const result = document.getElementById('analyticsResult');

    if (!short_url) {
        showResult(result, '❌ Введите короткую ссылку', 'error');
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/api/analytics/${short_url}`);
        const data = await response.json();

        if (response.ok) {
            const html = `
                <div class="success">
                    <strong>📊 Аналитика для: ${data.short_url}</strong>
                    <div class="analytics-data">
Всего переходов: ${data.total_count || 0}

📅 ПЕРЕХОДЫ ПО ДНЯМ:
${formatObjectData(data.day_count)}

📊 ПЕРЕХОДЫ ПО МЕСЯЦАМ:
${formatObjectData(data.month_count)}

🖥️ USER-AGENT:
${formatObjectData(data.user_agent_count)}
                    </div>
                </div>
            `;
            showResult(result, html, 'success');
            if (!shortUrlParam) document.getElementById('analyticsInput').value = '';
        } else {
            let errorMessage = '❌ Ошибка получения аналитики';
            if (response.status === 400) errorMessage = '❌ Неверный формат';
            else if (data.error) errorMessage = `❌ ${data.error}`;
            showResult(result, errorMessage, 'error');
        }
    } catch (error) {
        showResult(result, `❌ Ошибка сети: ${error.message}`, 'error');
    }
}

// ==================== ТЕСТ РЕДИРЕКТА ====================
async function testRedirect(shortUrlParam = null) {
    const short_url = shortUrlParam || document.getElementById('testInput').value.trim();
    const result = document.getElementById('testResult');

    if (!short_url) {
        showResult(result, '❌ Введите короткую ссылку', 'error');
        return;
    }

    try {
        // 🔴 ПРОБЛЕМА БЫЛА ЗДЕСЬ: используем обычный fetch без redirect: 'manual'
        const response = await fetch(`${API_BASE}/api/s/${short_url}`);

        let message, type;

        switch (response.status) {
            case 200:
            case 301:
            case 302:
                message = '✅ Редирект работает! Открываю оригинальный URL...';
                type = 'success';
                showResult(result, message, type);
                
                // ✅ ОТКРЫВАЕМ ОРИГИНАЛЬНЫЙ URL ИЗ ЗАГОЛОВКОВ
                const originalUrl = response.url; // Браузер уже сделал редирект
                window.open(originalUrl, '_blank');
                break;
                
            case 400:
                message = '❌ Неверный формат короткой ссылки';
                type = 'error';
                showResult(result, message, type);
                break;
                
            case 404:
                message = '❌ Ссылка не найдена';
                type = 'error';
                showResult(result, message, type);
                break;
                
            case 500:
                message = '❌ Ошибка сервера';
                type = 'error';
                showResult(result, message, type);
                break;
                
            default:
                message = `⚠️ Статус: ${response.status}`;
                type = 'warning';
                showResult(result, message, type);
        }

        if (!shortUrlParam) document.getElementById('testInput').value = '';

    } catch (error) {
        showResult(result, `❌ Ошибка сети: ${error.message}`, 'error');
    }
}

// ==================== ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ====================
function showResult(element, content, type) {
    element.innerHTML = content;
    element.className = `result ${type}`;
}

function formatObjectData(obj) {
    if (!obj || Object.keys(obj).length === 0) return '  Нет данных';
    return Object.entries(obj).map(([key, value]) => `  ${key}: ${value}`).join('\n');
}

function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(() => {
        alert('✅ Ссылка скопирована!');
    }).catch(err => {
        alert('❌ Ошибка копирования: ' + err);
    });
}

// ==================== ОБРАБОТЧИКИ ENTER ====================
document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('urlInput').addEventListener('keypress', function(e) {
        if (e.key === 'Enter') createShortUrl();
    });
    
    document.getElementById('analyticsInput').addEventListener('keypress', function(e) {
        if (e.key === 'Enter') getAnalytics();
    });
    
    document.getElementById('testInput').addEventListener('keypress', function(e) {
        if (e.key === 'Enter') testRedirect();
    });
});