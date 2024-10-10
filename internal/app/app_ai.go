package app

import "time"

func (a *appImpl) CreateSnippetWithAI() {
	if ok, text := a.tui.ShowAiPrompt(); ok {
		stopChan := make(chan bool)

		// Run the spinner in a separate goroutine
		go a.tui.ShowSpinner(text, stopChan)

		// Simulate doing some work in the main goroutine
		time.Sleep(2 * time.Second)

		// Send stop signal to stop the spinner
		stopChan <- true

		//nolint:mnd // Wait briefly to ensure spinner quits cleanly
		time.Sleep(100 * time.Millisecond)
	}
}
