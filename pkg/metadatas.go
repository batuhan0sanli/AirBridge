package pkg

// FileMetadata contains metadata about a file.
type FileMetadata struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Hash string `json:"hash"`
}
