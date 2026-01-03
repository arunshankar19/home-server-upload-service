package domain

// InitiateUploadReqDTO used to represent request details of upload request
type InitiateUploadReqDTO struct {
	FileName      string `json:"file_name"`
	FileSize      int    `json:"file_size"` // in bytes
	FileType      string `json:"file_type"`
	FileExtension string `json:"file_extension"`
	ChunkSize     int    `json:"chunk_size"`
	NoOfChunks    int    `json:"no_of_chunks"`
}

// InitiateUploadResponseDTO represents init upload response
type InitiateUploadResponseDTO struct {
	UploadID string `json:"upload_id"`
}
