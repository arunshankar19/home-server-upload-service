package fileexplorer

import (
	"context"

	"github.com/arunshankar19/home-server-upload-service/domain"
)

// FileExplorerClient represents file explorer service funcs
type FileExplorerClient interface {
	NotifyFileExplorer(ctx context.Context, file domain.UploadEvent) error
}
