package response

// Envelope is the common JSON response shape used across the API.
type Envelope struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Success(message string, data any) Envelope {
	return Envelope{Status: "success", Message: message, Data: data}
}

func Error(message string) Envelope {
	return Envelope{Status: "error", Message: message}
}
