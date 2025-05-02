// Global variables
let sockets = {};
let reconnectIntervals = {};
let messageHistory = {};
let connectionCounter = 0;

// DOM ready event
document.addEventListener('DOMContentLoaded', function() {
    // Theme handling
    initTheme();
    
    // Initialize token
    initToken();
    
    // Add event listeners
    document.getElementById('save-token').addEventListener('click', saveToken);
    document.getElementById('add-connection-btn').addEventListener('click', addConnection);
    
    // Load saved connections
    loadSavedConnections();
});

/**
 * Initialize theme preference
 */
function initTheme() {
    const themeToggle = document.getElementById('theme-toggle');
    
    // Check for saved theme preference
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme === 'dark') {
        document.documentElement.setAttribute('data-theme', 'dark');
        themeToggle.checked = true;
    }
    
    // Theme toggle event listener
    themeToggle.addEventListener('change', function() {
        if (this.checked) {
            document.documentElement.setAttribute('data-theme', 'dark');
            localStorage.setItem('theme', 'dark');
        } else {
            document.documentElement.removeAttribute('data-theme');
            localStorage.setItem('theme', 'light');
        }
    });
}

/**
 * Initialize authentication token
 */
function initToken() {
    const savedToken = localStorage.getItem("ws_auth_token");
    if (savedToken) {
        document.getElementById("token").value = savedToken;
    }
}

/**
 * Save authentication token
 */
function saveToken() {
    const token = document.getElementById("token").value;
    if (token) {
        localStorage.setItem("ws_auth_token", token);
        alert("Token saved successfully!");
    } else {
        alert("Please enter a token first");
    }
}

/**
 * Add a new connection
 */
function addConnection() {
    connectionCounter++;
    const connectionId = connectionCounter;

    // Initialize message history for this connection
    messageHistory[connectionId] = [];

    // Create connection element
    const connectionElement = createConnectionElement(connectionId);

    // Add to DOM
    document.getElementById('connections-container').appendChild(connectionElement);

    // Add event listeners
    setupConnectionEventListeners(connectionId);

    // Save connection IDs to localStorage
    saveConnectionsMetadata();

    return connectionId;
}

/**
 * Create connection HTML element
 */
function createConnectionElement(connectionId) {
    const connectionElement = document.createElement('div');
    connectionElement.className = 'conn-block';
    connectionElement.id = `connection-${connectionId}`;
    
    // Default URL
    const defaultUrl = `ws://localhost:8090/websocket/conn${connectionId}`;
    
    // Load saved URL if exists
    const savedUrl = localStorage.getItem(`wsurl${connectionId}`) || defaultUrl;
    
    connectionElement.innerHTML = `
        <button class="remove-connection" data-id="${connectionId}">×</button>
        <div class="conn-header">
            <span id="status-indicator${connectionId}" class="status-indicator status-disconnected" title="Disconnected"></span>
            <h3>🔌 Connection ${connectionId}</h3>
        </div>
        <div class="control-row">
            <input id="url${connectionId}" value="${savedUrl}" size="60" placeholder="WebSocket URL">
            <button id="connection-btn${connectionId}" class="connection-btn connect" data-id="${connectionId}">Connect</button>
            <div class="auto-reconnect">
                <label for="auto-reconnect${connectionId}">Auto-reconnect:</label>
                <input type="checkbox" id="auto-reconnect${connectionId}">
            </div>
        </div>
        <textarea id="log${connectionId}" readonly placeholder="Connection logs will appear here..."></textarea>
        <div class="control-row">
            <input id="msg${connectionId}" size="50" placeholder="Type message here">
            <div class="message-type">
                <label for="msg-type${connectionId}">Format:</label>
                <select id="msg-type${connectionId}">
                    <option value="text">Text</option>
                    <option value="json">JSON</option>
                </select>
            </div>
            <button id="send-btn${connectionId}" data-id="${connectionId}">Send</button>
            <button id="clear-btn${connectionId}" class="clear-btn" data-id="${connectionId}">Clear Log</button>
        </div>
        <div class="history-controls">
            <select id="message-history${connectionId}" class="history-dropdown">
                <option value="">-- Recent Messages --</option>
            </select>
            <button id="load-history-btn${connectionId}" data-id="${connectionId}">Load</button>
            <button id="save-history-btn${connectionId}" data-id="${connectionId}">Save Current</button>
        </div>
    `;
    
    return connectionElement;
}

/**
 * Setup event listeners for a connection
 */
function setupConnectionEventListeners(connectionId) {
    // Connect/disconnect button
    document.getElementById(`connection-btn${connectionId}`).addEventListener('click', () => toggleConnection(connectionId));
    
    // Send message button
    document.getElementById(`send-btn${connectionId}`).addEventListener('click', () => sendMsg(connectionId));
    
    // Clear log button
    document.getElementById(`clear-btn${connectionId}`).addEventListener('click', () => clearLog(connectionId));
    
    // Load history button
    document.getElementById(`load-history-btn${connectionId}`).addEventListener('click', () => loadFromHistory(connectionId));
    
    // Save to history button
    document.getElementById(`save-history-btn${connectionId}`).addEventListener('click', () => saveToHistory(connectionId));
    
    // Remove connection button
    const removeBtn = document.querySelector(`.remove-connection[data-id="${connectionId}"]`);
    if (removeBtn) {
        removeBtn.addEventListener('click', () => removeConnection(connectionId));
    }
    
    // Enter key to send message
    document.getElementById(`msg${connectionId}`).addEventListener('keypress', (event) => {
        if (event.key === 'Enter') {
            sendMsg(connectionId);
        }
    });
    
    // Load saved log
    loadHistory(connectionId);
    
    // Load message history
    loadMessageHistory(connectionId);
}

/**
 * Toggle WebSocket connection (connect/disconnect)
 */
function toggleConnection(connectionId) {
    const connectionBtn = document.getElementById(`connection-btn${connectionId}`);
    
    if (connectionBtn.innerText === "Connect") {
        connect(connectionId);
    } else {
        disconnect(connectionId);
    }
}

/**
 * Connect to WebSocket
 */
function connect(connectionId) {
    // Clear any existing reconnect interval
    if (reconnectIntervals[connectionId]) {
        clearInterval(reconnectIntervals[connectionId]);
        reconnectIntervals[connectionId] = null;
    }
    
    // Close existing connection if any
    if (sockets[connectionId]) {
        sockets[connectionId].close();
    }
    
    const url = document.getElementById(`url${connectionId}`).value;
    const token = document.getElementById("token").value || localStorage.getItem("ws_auth_token");
    
    try {
        // Save URL to localStorage
        localStorage.setItem(`wsurl${connectionId}`, url);
        
        // Create WebSocket connection
        let wsUrl = new URL(url);
        
        // Add token as query parameter if provided
        if (token) {
            wsUrl.searchParams.set("token", token);
        }
        
        const socket = new WebSocket(wsUrl);
        
        socket.onopen = () => {
            log(connectionId, "🟢 Connected successfully");
            updateStatusIndicator(connectionId, true);
        };
        
        socket.onmessage = (e) => {
            try {
                // Try to parse as JSON for prettier display
                const data = JSON.parse(e.data);
                log(connectionId, "📩 " + JSON.stringify(data, null, 2));
            } catch {
                // If not JSON, display as regular text
                log(connectionId, "📩 " + e.data);
            }
        };
        
        socket.onclose = (e) => {
            updateStatusIndicator(connectionId, false);
            log(connectionId, `🔴 Connection closed (code: ${e.code}, reason: ${e.reason || "none"})`);
            
            // Check if auto-reconnect is enabled
            const autoReconnect = document.getElementById(`auto-reconnect${connectionId}`).checked;
            if (autoReconnect && !reconnectIntervals[connectionId]) {
                log(connectionId, "⏱️ Attempting to reconnect in 5 seconds...");
                reconnectIntervals[connectionId] = setInterval(() => {
                    log(connectionId, "🔄 Reconnecting...");
                    connect(connectionId);
                }, 5000);
            }
        };
        
        socket.onerror = (e) => {
            updateStatusIndicator(connectionId, false);
            log(connectionId, "❌ Connection error occurred");
        };
        
        sockets[connectionId] = socket;
        
    } catch (error) {
        log(connectionId, `❌ Error: ${error.message}`);
    }
}

/**
 * Disconnect from WebSocket
 */
function disconnect(connectionId) {
    // Clear any reconnect interval
    if (reconnectIntervals[connectionId]) {
        clearInterval(reconnectIntervals[connectionId]);
        reconnectIntervals[connectionId] = null;
    }
    
    const socket = sockets[connectionId];
    if (socket) {
        socket.close(1000, "User initiated disconnect");
        updateStatusIndicator(connectionId, false);
    } else {
        log(connectionId, "⚠️ No active connection to disconnect");
    }
}

/**
 * Update connection status indicator
 */
function updateStatusIndicator(connectionId, connected) {
    const indicator = document.getElementById(`status-indicator${connectionId}`);
    const connectionBtn = document.getElementById(`connection-btn${connectionId}`);
    
    if (connected) {
        indicator.className = "status-indicator status-connected";
        indicator.title = "Connected";
        connectionBtn.innerText = "Disconnect";
        connectionBtn.className = "connection-btn disconnect";
    } else {
        indicator.className = "status-indicator status-disconnected";
        indicator.title = "Disconnected";
        connectionBtn.innerText = "Connect";
        connectionBtn.className = "connection-btn connect";
    }
}

/**
 * Send a message through WebSocket
 */
function sendMsg(connectionId) {
    const msgElem = document.getElementById(`msg${connectionId}`);
    const msgType = document.getElementById(`msg-type${connectionId}`).value;
    let msg = msgElem.value;
    const socket = sockets[connectionId];
    
    if (!msg.trim()) {
        log(connectionId, "⚠️ Please enter a message to send");
        return;
    }
    
    if (socket && socket.readyState === WebSocket.OPEN) {
        // Format as JSON if selected
        if (msgType === "json") {
            try {
                // Check if it's already valid JSON
                JSON.parse(msg);
            } catch {
                // If not, try to convert it to a JSON object
                try {
                    msg = JSON.stringify({ message: msg });
                } catch (e) {
                    log(connectionId, `⚠️ Failed to format as JSON: ${e.message}`);
                    return;
                }
            }
        }
        
        socket.send(msg);
        log(connectionId, "📤 " + msg);
        
        // Save to message history
        saveMessageToHistory(connectionId, msg);

//        msgElem.value = "";
    } else {
        log(connectionId, "⚠️ Socket not connected. Message not sent.");
    }
}

/**
 * Add timestamp to log messages
 */
function getTimestamp() {
    const now = new Date();
    return now.toLocaleTimeString();
}

/**
 * Log message to connection's log area
 */
function log(connectionId, msg) {
    const logElem = document.getElementById(`log${connectionId}`);
    if (!logElem) return;
    
    const timestamp = getTimestamp();
    logElem.value += `[${timestamp}] ${msg}\n`;
    logElem.scrollTop = logElem.scrollHeight;
    saveHistory(connectionId, logElem.value);
}

/**
 * Clear log for a connection
 */
function clearLog(connectionId) {
    document.getElementById(`log${connectionId}`).value = "";
    saveHistory(connectionId, "");
}

/**
 * Save connection log to localStorage
 */
function saveHistory(connectionId, content) {
    localStorage.setItem(`wslog${connectionId}`, content);
}

/**
 * Load connection log from localStorage
 */
function loadHistory(connectionId) {
    const saved = localStorage.getItem(`wslog${connectionId}`);
    if (saved) {
        const logElem = document.getElementById(`log${connectionId}`);
        if (logElem) {
            logElem.value = saved;
        }
    }
}

/**
 * Save message to history
 */
function saveMessageToHistory(connectionId, message) {
    // Ensure the history array exists
    if (!messageHistory[connectionId]) {
        messageHistory[connectionId] = [];
    }
    
    // Add to our message history array
    if (!messageHistory[connectionId].includes(message)) {
        messageHistory[connectionId].unshift(message);
        // Keep only last 10 messages
        if (messageHistory[connectionId].length > 10) {
            messageHistory[connectionId].pop();
        }
        
        // Save to localStorage
        localStorage.setItem(`ws_msg_history${connectionId}`, JSON.stringify(messageHistory[connectionId]));
        
        // Update dropdown
        updateHistoryDropdown(connectionId);
    }
}

/**
 * Update history dropdown for a connection
 */
function updateHistoryDropdown(connectionId) {
    const dropdown = document.getElementById(`message-history${connectionId}`);
    if (!dropdown) return;
    
    // Clear current options except first one
    while (dropdown.options.length > 1) {
        dropdown.remove(1);
    }
    
    // Add history items
    if (messageHistory[connectionId]) {
        messageHistory[connectionId].forEach(msg => {
            const displayText = msg.length > 30 ? msg.substring(0, 27) + '...' : msg;
            const option = new Option(displayText, msg);
            dropdown.add(option);
        });
    }
}

/**
 * Load message from history dropdown
 */
function loadFromHistory(connectionId) {
    const dropdown = document.getElementById(`message-history${connectionId}`);
    const selectedValue = dropdown.value;
    
    if (selectedValue) {
        document.getElementById(`msg${connectionId}`).value = selectedValue;
    }
}

/**
 * Save current message to history
 */
function saveToHistory(connectionId) {
    const msg = document.getElementById(`msg${connectionId}`).value;
    if (msg.trim()) {
        saveMessageToHistory(connectionId, msg);
    }
}

/**
 * Load message history from localStorage
 */
function loadMessageHistory(connectionId) {
    const saved = localStorage.getItem(`ws_msg_history${connectionId}`);
    if (saved) {
        try {
            messageHistory[connectionId] = JSON.parse(saved);
            updateHistoryDropdown(connectionId);
        } catch (e) {
            console.error(`Failed to parse message history for connection ${connectionId}:`, e);
        }
    }
}

/**
 * Remove a connection
 */
function removeConnection(connectionId) {
    // Disconnect if connected
    if (sockets[connectionId]) {
        disconnect(connectionId);
    }

    // Remove from DOM
    const connectionElement = document.getElementById(`connection-${connectionId}`);
    if (connectionElement) {
        connectionElement.remove();
    }

    // Clean up resources
    delete sockets[connectionId];
    delete reconnectIntervals[connectionId];
    delete messageHistory[connectionId];

    // Update storage
    saveConnectionsMetadata();
}

/**
 * Save connections metadata to localStorage
 */
function saveConnectionsMetadata() {
    const connectionElements = document.querySelectorAll('.conn-block');
    const connectionIds = Array.from(connectionElements).map(el =>
        parseInt(el.id.replace('connection-', ''))
    );

    localStorage.setItem('ws_connections', JSON.stringify(connectionIds));
    localStorage.setItem('ws_connection_counter', connectionCounter);
}

/**
 * Load saved connections
 */
function loadSavedConnections() {
    // Get connection counter
    const savedCounter = localStorage.getItem('ws_connection_counter');
    if (savedCounter) {
        connectionCounter = parseInt(savedCounter);
    }

    // Get saved connection IDs
    const savedConnections = localStorage.getItem('ws_connections');
    if (savedConnections) {
        try {
            const connectionIds = JSON.parse(savedConnections);

            // Create each saved connection
            connectionIds.forEach(id => {
                // Create connection element with this ID
                const connectionId = parseInt(id);

                // Create connection element
                const connectionElement = createConnectionElement(connectionId);

                // Add to DOM
                document.getElementById('connections-container').appendChild(connectionElement);

                // Add event listeners
                setupConnectionEventListeners(connectionId);

                // Initialize message history array if needed
                if (!messageHistory[connectionId]) {
                    messageHistory[connectionId] = [];
                }

                // Check if auto-reconnect was enabled and reconnect if necessary
                const autoReconnect = localStorage.getItem(`ws_auto_reconnect${connectionId}`);
                if (autoReconnect === 'true') {
                    document.getElementById(`auto-reconnect${connectionId}`).checked = true;

                    // Check if it was previously connected
                    const wasConnected = localStorage.getItem(`ws_was_connected${connectionId}`);
                    if (wasConnected === 'true') {
                        // Attempt to reconnect with a small delay to ensure DOM is ready
                        setTimeout(() => {
                            connect(connectionId);
                        }, 500);
                    }
                }
            });

        } catch (e) {
            console.error("Failed to load saved connections:", e);
        }
    } else {
        // If no saved connections, create a default one
        addConnection();
    }
}

/**
 * Save connection state before page unload
 */
window.addEventListener('beforeunload', function() {
    // Save connection states for potential auto-reconnect
    Object.keys(sockets).forEach(connectionId => {
        const socket = sockets[connectionId];
        const isConnected = socket && socket.readyState === WebSocket.OPEN;
        const autoReconnect = document.getElementById(`auto-reconnect${connectionId}`).checked;
        
        localStorage.setItem(`ws_was_connected${connectionId}`, isConnected);
        localStorage.setItem(`ws_auto_reconnect${connectionId}`, autoReconnect);
    });
});

/**
 * Check WebSocket URL format
 * @param {string} url - WebSocket URL to validate
 * @returns {boolean} - True if valid WebSocket URL
 */
function isValidWebSocketUrl(url) {
    try {
        const parsedUrl = new URL(url);
        return parsedUrl.protocol === 'ws:' || parsedUrl.protocol === 'wss:';
    } catch (e) {
        return false;
    }
}

/**
 * Export all connection logs
 */
function exportLogs() {
    const exportData = {};
    
    // Get all connection IDs
    const connectionElements = document.querySelectorAll('.conn-block');
    const connectionIds = Array.from(connectionElements).map(el => 
        parseInt(el.id.replace('connection-', ''))
    );
    
    // Collect all logs
    connectionIds.forEach(id => {
        const logElem = document.getElementById(`log${id}`);
        const url = document.getElementById(`url${id}`).value;
        
        exportData[`Connection ${id}`] = {
            url: url,
            log: logElem.value
        };
    });
    
    // Create download link
    const dataStr = "data:text/json;charset=utf-8," + encodeURIComponent(JSON.stringify(exportData, null, 2));
    const downloadAnchorNode = document.createElement('a');
    downloadAnchorNode.setAttribute("href", dataStr);
    downloadAnchorNode.setAttribute("download", "websocket_logs_" + new Date().toISOString().slice(0, 10) + ".json");
    document.body.appendChild(downloadAnchorNode);
    downloadAnchorNode.click();
    downloadAnchorNode.remove();
}
