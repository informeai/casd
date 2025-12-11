package dto

import "time"

// Chunk represents a data chunk stored in the system.
type Chunk struct {
	Data []byte `json:"data"`
}

// Formula represents a formula that consists of multiple data chunks.
type Formula struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Hashs     []string  `json:"hashs"`
}
