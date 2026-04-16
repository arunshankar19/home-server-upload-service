package domain

// InitiateUploadReqDTO used to represent request details of upload request
type InitiateUploadReqDTO struct {
	FileName      string `json:"file_name"`
	FileSize      int    `json:"file_size"` // in bytes
	FileType      string `json:"file_type"`
	FileExtension string `json:"file_extension"`
	IsMultipart   bool   `json:"is_multipart"`
	ChunkSize     int    `json:"chunk_size"`
	NoOfChunks    int    `json:"no_of_chunks"`
}

// InitiateUploadResponseDTO represents init upload response
type InitiateUploadResponseDTO struct {
	UploadID      string `json:"upload_id"`
	UploadEventID string `json:"upload_event_id"`
}

// PresignedMultipartURLReqDTO represents presigned multipart request
type PresignedMultipartURLReqDTO struct {
	UploadEventID string `json:"upload_event_id"`
	IsMultipart   bool   `json:"is_multipart"`
	PartNumber    int    `json:"part_number"`
}

// PresignedMultipartResponseDTO represent presigned multipart response
type PresignedMultipartResponseDTO struct {
	PresignedURL string `json:"presigned_url"`
}

// CompleteUploadReqDTO represents part numbers and etag for multipart upload
type CompleteUploadReqDTO struct {
	UploadEventID string                  `json:"upload_event_id"`
	Etags         []PartNumberEtagMapping `json:"e_tags"`
}

type PartNumberEtagMapping struct {
	PartNumber int    `json:"part_number"`
	Etag       string `json:"e_tag"`
}
