package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// TextGenerator generates text for typing practice
type TextGenerator struct {
	apiKey string
	model  string
}

// NewTextGenerator creates a new text generator
func NewTextGenerator(apiKey string) *TextGenerator {
	return &TextGenerator{
		apiKey: apiKey,
		model:  "gemini-1.5-flash",
	}
}

// GeminiRequest represents a request to the Gemini API
type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

// GeminiContent represents the content of a Gemini request
type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart represents a part of a Gemini content
type GeminiPart struct {
	Text string `json:"text"`
}

// GeminiResponse represents a response from the Gemini API
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// GenerateText generates text based on a prompt
func (g *TextGenerator) GenerateText(prompt string) (string, error) {
	// If API key is empty, fall back to the template-based approach
	if g.apiKey == "" {
		fmt.Println("No API key provided, using fallback text generation")
		if strings.Contains(prompt, "AT LEAST") || strings.Contains(prompt, "Practice:") {
			return g.generatePracticeText(prompt)
		}
		return g.generateRegularText(prompt)
	}

	// Use Gemini API for text generation
	fmt.Println("Using Gemini API for text generation")
	return g.generateWithGemini(prompt)
}

// generateWithGemini generates text using the Gemini API
func (g *TextGenerator) generateWithGemini(prompt string) (string, error) {
	// Enhance the prompt with typing-specific instructions
	enhancedPrompt := enhancePromptForTyping(prompt)

	fmt.Printf("Calling Gemini API with enhanced prompt: %s\n", enhancedPrompt)
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", g.model, g.apiKey)
	fmt.Printf("API URL: %s\n", url[:60]+"...") // Only show part of URL to avoid exposing full API key

	// Create the request body
	requestBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{
						Text: enhancedPrompt,
					},
				},
			},
		},
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}
	fmt.Printf("Request JSON: %s\n", string(jsonData))

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	fmt.Println("Sending request to Gemini API...")
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	fmt.Printf("Response status: %d\n", resp.StatusCode)
	fmt.Printf("Response body: %s\n", string(body))

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var geminiResponse GeminiResponse
	if err := json.Unmarshal(body, &geminiResponse); err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	// Extract the generated text
	if len(geminiResponse.Candidates) > 0 && len(geminiResponse.Candidates[0].Content.Parts) > 0 {
		generatedText := geminiResponse.Candidates[0].Content.Parts[0].Text
		fmt.Printf("Successfully generated text (%d characters)\n", len(generatedText))
		return generatedText, nil
	}

	return "", fmt.Errorf("no text generated in response")
}

// enhancePromptForTyping adds typing-specific instructions to the prompt
func enhancePromptForTyping(originalPrompt string) string {
	// If it's already a practice prompt with specific words, don't modify it too much
	if strings.Contains(originalPrompt, "AT LEAST") {
		return originalPrompt + "\n\nAdditional instructions: Make the text flow naturally while incorporating the required words. Use simple sentence structures that are easy to type."
	}

	// For regular prompts, add more comprehensive typing instructions
	baseInstructions := `
Generate a typing practice text with the following characteristics:
1. Keep it between 30-50 words unless a different length is specified
2. Use a mix of common and less common words to practice different finger movements
3. Include some punctuation for practice (commas, periods, question marks)
4. Avoid very long words or extremely technical terms unless specifically requested
5. Create coherent, meaningful content that's engaging to type
6. Include a balanced mix of letters that exercise both hands evenly
7. Incorporate some capital letters naturally within the text
8. Use simple sentence structures that flow well for typing practice

Based on this request: `

	return baseInstructions + originalPrompt
}

// generatePracticeText creates a practice text with repeated words (fallback method)
func (g *TextGenerator) generatePracticeText(prompt string) (string, error) {
	// Extract words from the prompt
	words := extractWordsFromPrompt(prompt)
	if len(words) == 0 {
		return "Please provide words to practice.", nil
	}

	// Create sentences that use these words multiple times
	sentences := []string{}
	templates := []string{
		"I need to practice typing the word %s correctly.",
		"The word %s is challenging for me to type accurately.",
		"When I type %s, I should focus on each letter carefully.",
		"Typing %s requires attention to detail and precision.",
		"I will improve my accuracy when typing %s with practice.",
		"The more I practice typing %s, the better I will become.",
		"Each time I type %s, I should check for errors.",
		"Careful typing of %s will help me build muscle memory.",
		"I should slow down when typing %s to avoid mistakes.",
		"Repetition of typing %s will help me master it.",
	}

	// Generate 3 sentences for each word
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, word := range words {
		// Use each template at least once for this word
		for i := 0; i < 3; i++ {
			templateIndex := r.Intn(len(templates))
			sentence := fmt.Sprintf(templates[templateIndex], word)
			sentences = append(sentences, sentence)
		}
	}

	// Mix the sentences
	r.Shuffle(len(sentences), func(i, j int) {
		sentences[i], sentences[j] = sentences[j], sentences[i]
	})

	// Add some connecting phrases between groups of sentences
	connectors := []string{
		"Let's continue practicing. ",
		"Moving on to more practice. ",
		"Now for some more typing practice. ",
		"Let's focus on these words again. ",
		"Continuing with our practice session. ",
	}

	// Build the final text with connectors between groups of sentences
	var result strings.Builder
	for i := 0; i < len(sentences); i++ {
		if i > 0 && i%3 == 0 {
			result.WriteString(connectors[r.Intn(len(connectors))])
		}
		result.WriteString(sentences[i] + " ")
	}

	return result.String(), nil
}

// generateRegularText returns a predefined text based on the prompt keywords (fallback method)
func (g *TextGenerator) generateRegularText(prompt string) (string, error) {
	prompt = strings.ToLower(prompt)

	// Check for keywords in the prompt and return appropriate text
	if strings.Contains(prompt, "python") || strings.Contains(prompt, "code") || strings.Contains(prompt, "programming") {
		return codingText, nil
	} else if strings.Contains(prompt, "poem") || strings.Contains(prompt, "poetry") {
		return poemText, nil
	} else if strings.Contains(prompt, "science") || strings.Contains(prompt, "tech") {
		return scienceText, nil
	} else {
		return generalText, nil
	}
}

// extractWordsFromPrompt extracts practice words from a prompt
func extractWordsFromPrompt(prompt string) []string {
	fmt.Printf("Extracting words from prompt: %s\n", prompt)

	// For the practice text from HandleGeneratePractice, extract words from the AT LEAST part
	if strings.Contains(prompt, "AT LEAST") {
		parts := strings.Split(prompt, "AT LEAST")
		if len(parts) >= 2 {
			// Extract the words after "AT LEAST" and before the period
			wordsPart := strings.Split(parts[1], ":")
			if len(wordsPart) >= 2 {
				wordsList := strings.Split(wordsPart[1], ".")
				if len(wordsList) >= 1 {
					cleanWords := strings.TrimSpace(wordsList[0])
					wordList := strings.Split(cleanWords, ",")

					// Clean up each word
					var result []string
					for _, word := range wordList {
						word = strings.TrimSpace(word)
						if word != "" {
							fmt.Printf("Found practice word: %s\n", word)
							result = append(result, word)
						}
					}
					return result
				}
			}
		}
	}

	// Check if this is a practice prompt with the format "Practice: word1, word2"
	if strings.Contains(prompt, "Practice:") {
		// Extract the words part
		parts := strings.SplitN(prompt, "Practice:", 2)
		if len(parts) < 2 {
			return []string{}
		}

		// Split by commas and clean up
		wordPart := parts[1]
		wordList := strings.Split(wordPart, ",")

		// Clean up each word
		var result []string
		for _, word := range wordList {
			word = strings.TrimSpace(word)
			if word != "" {
				fmt.Printf("Found practice word: %s\n", word)
				result = append(result, word)
			}
		}

		return result
	}

	return []string{}
}

// Predefined texts for different categories (used as fallback)
var generalText = `The ability to type quickly and accurately is an essential skill in today's digital world. Regular practice can significantly improve your typing speed and reduce errors. Focus on maintaining proper finger positioning on the home row keys and try to look at the screen instead of your keyboard. With consistent practice, typing will become second nature, allowing you to focus more on content creation rather than the mechanical process of typing.`

var codingText = `def calculate_fibonacci(n):
    """
    Calculate the Fibonacci sequence up to the nth term.
    The Fibonacci sequence starts with 0 and 1, and each subsequent number
    is the sum of the two preceding ones.
    
    Args:
        n: The number of terms to calculate
        
    Returns:
        A list containing the Fibonacci sequence
    """
    if n <= 0:
        return []
    elif n == 1:
        return [0]
    
    fibonacci = [0, 1]
    for i in range(2, n):
        fibonacci.append(fibonacci[i-1] + fibonacci[i-2])
    
    return fibonacci`

var poemText = `The Programmer's Lament

I write the code, line by line,
Debugging errors, spending time,
On functions, loops, and arrays too,
Fixing bugs both old and new.

The cursor blinks, awaiting thought,
Solutions that cannot be bought,
With logic clear and purpose true,
I build the systems, through and through.

In silence deep, I contemplate,
Algorithms that calculate,
The patterns hidden in the noise,
Data structures are my toys.

When darkness falls, I'm still awake,
One more feature I must make,
The perfect code remains my quest,
Until at last I take my rest.`

var scienceText = `Quantum computing represents a significant leap in computational power by harnessing the principles of quantum mechanics. Unlike classical computers that use bits (0s and 1s), quantum computers use quantum bits or "qubits" that can exist in multiple states simultaneously due to superposition. This property allows quantum computers to process vast amounts of information in parallel, potentially solving complex problems that are currently intractable. Researchers are exploring applications in cryptography, drug discovery, optimization problems, and material science. Despite significant progress, quantum computers still face challenges related to qubit stability, error correction, and scaling up to practical sizes.`
