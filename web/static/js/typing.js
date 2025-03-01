document.addEventListener('DOMContentLoaded', function() {
    // Initialize typing functionality
    initTyping();
    
    // Also listen for HTMX content swaps (when new text is generated)
    document.body.addEventListener('htmx:afterSwap', function(event) {
        console.log("HTMX content swapped, reinitializing typing");
        initTyping();
    });
});

function initTyping() {
    const textDisplay = document.getElementById('text-display');
    const textContainer = document.getElementById('typing-text');
    
    if (!textDisplay || !textContainer) {
        console.error("Missing required elements");
        return;
    }
    
    const textId = textContainer.dataset.textId;
    const originalText = textContainer.dataset.content;
    
    if (!textId || !originalText) {
        console.error("Missing text ID or content");
        return;
    }
    
    console.log("Initializing typing with text ID:", textId);
    
    // Remove any existing cursor before creating a new one
    const existingCursor = document.getElementById('typing-cursor');
    if (existingCursor) {
        existingCursor.remove();
    }
    
    // Create a persistent cursor element
    let cursor = document.createElement('div');
    cursor.id = 'typing-cursor';
    document.body.appendChild(cursor);
    
    // Variables to track typing state
    let typedText = '';
    let startTime = null;
    let isSessionActive = false;
    let errorCount = 0;
    let errorDetails = [];
    let wordsWithErrors = new Set();
    let timerInterval = null;
    let metricsUpdateInterval = null;
    
    // Create a timer element if it doesn't exist
    let timerElement = document.getElementById('typing-timer');
    if (!timerElement) {
        timerElement = document.createElement('div');
        timerElement.id = 'typing-timer';
        timerElement.className = 'text-xl font-bold text-yellow-400 text-center mb-4';
        timerElement.textContent = '00:00.00';
        textContainer.parentNode.insertBefore(timerElement, textContainer);
    }
    
    // Initialize the display
    initializeDisplay();
    
    // Make sure the text display is focusable
    textDisplay.setAttribute('tabindex', '0');
    textDisplay.focus();
    
    // Add a global click handler to ensure focus returns to the text display
    document.addEventListener('click', function() {
        textDisplay.focus();
    });
    
    // Add a variable to track typing activity
    let typingTimer = null;
    const typingDelay = 100; // 100ms delay before considering typing stopped
    
    // Handle keydown events
    textDisplay.addEventListener('keydown', function(e) {
        console.log("Key pressed:", e.key);
        
        // Prevent default behavior for all keys
        e.preventDefault();
        
        // Show solid cursor during typing (no blink)
        cursor.classList.add('typing');
        
        // Clear any existing typing timer
        if (typingTimer) {
            clearTimeout(typingTimer);
        }
        
        // Set a timer to remove the typing class after delay
        typingTimer = setTimeout(function() {
            cursor.classList.remove('typing');
        }, typingDelay);
        
        // If session is not active, start it on the first key press
        if (!isSessionActive && e.key.length === 1) {
            console.log("Starting session");
            startTime = new Date();
            isSessionActive = true;
            
            // Start the timer and metrics updates
            startTimer();
            startMetricsUpdates();
        }
        
        // Handle Backspace
        if (e.key === 'Backspace') {
            if (typedText.length > 0) {
                typedText = typedText.slice(0, -1);
                updateDisplay(typedText);
            }
            return;
        }
        
        // Handle regular typing
        if (e.key.length === 1) {
            // Check if this character is an error
            if (typedText.length < originalText.length && e.key !== originalText[typedText.length]) {
                errorCount++;
                document.getElementById('errors').textContent = errorCount;
            }
            
            typedText += e.key;
            updateDisplay(typedText);
            
            // Check if we've completed the text
            if (typedText.length >= originalText.length) {
                isSessionActive = false;
                stopTimer();
                submitResult();
            }
        }
    });
    
    // Function to initialize the display
    function initializeDisplay() {
        let displayHTML = '';
        
        for (let i = 0; i < originalText.length; i++) {
            displayHTML += `<span class="text-gray-300">${originalText[i]}</span>`;
        }
        
        textDisplay.innerHTML = displayHTML;
        
        // Position cursor at the beginning
        updateCursorPosition(0);
    }
    
    // Function to update the display based on typed text
    function updateDisplay(currentInput) {
        let displayHTML = '';
        
        for (let i = 0; i < originalText.length; i++) {
            if (i < currentInput.length) {
                if (currentInput[i] === originalText[i]) {
                    // Correct character
                    displayHTML += `<span class="text-white">${originalText[i]}</span>`;
                } else {
                    // Incorrect character
                    displayHTML += `<span class="text-red-500 bg-red-900">${originalText[i]}</span>`;
                }
            } else {
                // Not yet typed
                displayHTML += `<span class="text-gray-300">${originalText[i]}</span>`;
            }
        }
        
        textDisplay.innerHTML = displayHTML;
        
        // Update cursor position
        updateCursorPosition(currentInput.length);
        
        // Update metrics
        updateMetrics();
    }
    
    // Function to update cursor position
    function updateCursorPosition(position) {
        // Find the position where the cursor should be
        if (position < originalText.length) {
            const spans = textDisplay.querySelectorAll('span');
            if (spans.length > position) {
                const currentSpan = spans[position];
                const rect = currentSpan.getBoundingClientRect();
                
                // Get the position relative to the viewport
                const viewportX = rect.left;
                const viewportY = rect.top;
                
                // Get the scroll position
                const scrollX = window.pageXOffset || document.documentElement.scrollLeft;
                const scrollY = window.pageYOffset || document.documentElement.scrollTop;
                
                // Calculate absolute position (viewport + scroll)
                cursor.style.left = (viewportX + scrollX) + 'px';
                
                // Adjust the vertical position to align with text baseline
                cursor.style.top = (viewportY + scrollY) + 'px';
                
                // Make sure the cursor is visible
                cursor.style.display = 'block';
            }
        } else {
            // Position at the end
            const lastSpan = textDisplay.querySelector('span:last-child');
            if (lastSpan) {
                const rect = lastSpan.getBoundingClientRect();
                
                // Get the position relative to the viewport
                const viewportX = rect.right;
                const viewportY = rect.top;
                
                // Get the scroll position
                const scrollX = window.pageXOffset || document.documentElement.scrollLeft;
                const scrollY = window.pageYOffset || document.documentElement.scrollTop;
                
                // Calculate absolute position (viewport + scroll)
                cursor.style.left = (viewportX + scrollX) + 'px';
                cursor.style.top = (viewportY + scrollY) + 'px';
                
                // Make sure the cursor is visible
                cursor.style.display = 'block';
            }
        }
    }
    
    // Function to update metrics
    function updateMetrics() {
        if (!isSessionActive || !startTime) return;
        
        // Calculate WPM
        const elapsedTime = (new Date() - startTime) / 1000 / 60; // in minutes
        let wpm = 0;
        if (elapsedTime > 0) {
            wpm = (typedText.length / 5) / elapsedTime;
        }
        
        // Calculate accuracy
        let correctChars = 0;
        for (let i = 0; i < typedText.length; i++) {
            if (i < originalText.length && typedText[i] === originalText[i]) {
                correctChars++;
            }
        }
        
        let accuracy = 100;
        if (typedText.length > 0) {
            accuracy = (correctChars / typedText.length) * 100;
        }
        
        // Update the UI
        document.getElementById('wpm').textContent = wpm.toFixed(1);
        document.getElementById('accuracy').textContent = accuracy.toFixed(1) + '%';
        document.getElementById('errors').textContent = errorCount;
    }
    
    // Function to start the timer
    function startTimer() {
        if (timerInterval) clearInterval(timerInterval);
        
        timerInterval = setInterval(function() {
            const elapsedTime = new Date() - startTime;
            const minutes = Math.floor(elapsedTime / 60000);
            const seconds = Math.floor((elapsedTime % 60000) / 1000);
            const milliseconds = Math.floor((elapsedTime % 1000) / 10);
            
            timerElement.textContent = 
                (minutes < 10 ? '0' : '') + minutes + ':' + 
                (seconds < 10 ? '0' : '') + seconds + '.' + 
                (milliseconds < 10 ? '0' : '') + milliseconds;
        }, 10);
    }
    
    // Function to stop the timer
    function stopTimer() {
        if (timerInterval) {
            clearInterval(timerInterval);
            timerInterval = null;
        }
        
        if (metricsUpdateInterval) {
            clearInterval(metricsUpdateInterval);
            metricsUpdateInterval = null;
        }
    }
    
    // Function to start metrics updates
    function startMetricsUpdates() {
        if (metricsUpdateInterval) clearInterval(metricsUpdateInterval);
        
        metricsUpdateInterval = setInterval(function() {
            if (isSessionActive) {
                updateMetrics();
            }
        }, 500);
    }
    
    // Function to submit the result
    function submitResult() {
        // Collect error words
        const errorWords = Array.from(wordsWithErrors);
        
        // Create the result object
        const result = {
            text_id: parseInt(textId),
            wpm: parseFloat(document.getElementById('wpm').textContent),
            accuracy: parseFloat(document.getElementById('accuracy').textContent),
            errors: errorCount,
            error_details: errorDetails,
            error_words: errorWords
        };
        
        // Submit the result
        fetch('/submit-result', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(result)
        })
        .then(response => response.json())
        .then(data => {
            console.log("Result submitted:", data);
        })
        .catch(error => {
            console.error("Error submitting result:", error);
        });
    }
} 