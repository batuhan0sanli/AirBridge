package pkg

// SmallFilePayload represents a payload for a small file.
type SmallFilePayload struct {
	Key      string       `json:"key"`
	Data     string       `json:"data"`
	Nonce    string       `json:"nonce"`
	Metadata FileMetadata `json:"metadata"`
}

// LargeFilePayload represents a payload for a large file chunk.
type LargeFilePayload struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Hash     string `json:"hash"`
	ChunkNum int    `json:"chunk_num"`
}
