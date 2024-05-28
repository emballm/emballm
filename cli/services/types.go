package services

var Supported = supported{
	Ollama: "ollama",
	Vertex: "vertex",
}

type supported struct {
	Ollama string
	Vertex string
}

type Prompt struct {
	Messages []Message
}

type Message struct {
	Role    string
	Content string
}
