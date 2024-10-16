package gemini

type Request struct {
	SystemInstruction Instruction     `json:"system_instruction"`
	Contents          Content         `json:"contents"`
	SafetySettings    []SafetySetting `json:"safetySettings"`
}

type Instruction struct {
	Parts TextPart `json:"parts"`
}

type Content struct {
	Parts TextPart `json:"parts"`
}

type SafetySetting struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

type Response struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content ContentParts `json:"content"`
}

type ContentParts struct {
	Parts []TextPart `json:"parts"`
}

type TextPart struct {
	Text string `json:"text"`
}
