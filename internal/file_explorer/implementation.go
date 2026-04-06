package fileexplorer

import (
	"context"

	fileexplorer "github.com/arunshankar19/home-server-common-utils/proto/file_explorer/gen"
	"github.com/arunshankar19/home-server-upload-service/domain"
	"google.golang.org/grpc"
)

type fileExplorerClient struct {
	client fileexplorer.FileExplorerClient
}

// NewFileExplorerClient returns an implementation of file explorer
func NewFileExplorerClient(conn *grpc.ClientConn) FileExplorerClient {
	return &fileExplorerClient{
		client: fileexplorer.NewFileExplorerClient(conn),
	}
}

// NotifyFileExplorer notifies file server of the file upload
func (f *fileExplorerClient) NotifyFileExplorer(ctx context.Context, file domain.UploadEvent) error {
	_, err := f.client.NotifyUpload(ctx, &fileexplorer.FileUploadRequest{
		FileId:     file.ID.String(),
		FileName:   file.FileName,
		FileSize:   int64(file.FileSize),
		FileType:   file.FileType,
		FileExt:    file.FileExt,
		FileOwner:  file.CreatedBy.String(),
		UploadedAt: file.CreatedAt.Unix(),
	})
	return err
}
