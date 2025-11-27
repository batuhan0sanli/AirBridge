package pkg

type SmallFilePayload struct {
	Key      string       `json:"key"`
	Data     string       `json:"data"`
	IV       string       `json:"iv"`
	Metadata FileMetadata `json:"metadata"`
}

type largeFilePayload struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Hash     string `json:"hash"`
	ChunkNum int    `json:"chunk_num"`
}
