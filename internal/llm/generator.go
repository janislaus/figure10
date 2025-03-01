package llm

import (
	"errors"
	"strings"
)

// TextGenerator generates text using a local LLM
type TextGenerator struct {
	ModelPath string
}

// NewTextGenerator creates a new text generator
func NewTextGenerator(modelPath string) *TextGenerator {
	return &TextGenerator{
		ModelPath: modelPath,
	}
}

// GenerateText generates text based on a prompt
// This is a placeholder implementation that would be replaced with actual LLM integration
func (g *TextGenerator) GenerateText(prompt string) (string, error) {
	// For now, we'll return a mock response
	// In a real implementation, you would call your local LLM here

	// Example of how you might call a local LLM like llama.cpp
	// cmd := exec.Command("./llama", "-m", g.ModelPath, "--prompt", prompt, "--temp", "0.7", "--tokens", "300")
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	//     return "", err
	// }
	// return string(output), nil

	// Mock implementation for now
	if prompt == "" {
		return "", errors.New("prompt cannot be empty")
	}

	// Return different mock texts based on the prompt
	if strings.Contains(strings.ToLower(prompt), "python") {
		return `def fibonacci(n):
    """Return the nth Fibonacci number."""
    if n <= 0:
        return 0
    elif n == 1:
        return 1
    else:
        return fibonacci(n-1) + fibonacci(n-2)

# Calculate the first 10 Fibonacci numbers
for i in range(10):
    print(f"Fibonacci({i}) = {fibonacci(i)}")`, nil
	} else if strings.Contains(strings.ToLower(prompt), "poem") {
		return `Fingers dancing across the keys,
A symphony of clicks and taps,
Each keystroke a note in the melody,
Of thoughts transformed to digital maps.

Ten fingers working in harmony,
A bridge between mind and machine,
Practice makes perfect, they say to me,
As I type faster than I've ever been.`, nil
	} else if strings.Contains(strings.ToLower(prompt), "go") || strings.Contains(strings.ToLower(prompt), "golang") {
		return `package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello, Figure10!")
	
	// Create a ticker that ticks every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	// Count down from 10
	for i := 10; i > 0; i-- {
		<-ticker.C
		fmt.Printf("%d...\n", i)
	}
	
	fmt.Println("Start typing!")
}`, nil
	} else {
		return `The quick brown fox jumps over the lazy dog. This pangram contains every letter of the English alphabet at least once. Typing practice is essential for developing muscle memory and increasing your speed and accuracy. Keep your fingers on the home row keys: A, S, D, F for your left hand, and J, K, L, ; for your right hand. Look at the screen, not at your keyboard, to improve your touch typing skills.`, nil
	}
}
