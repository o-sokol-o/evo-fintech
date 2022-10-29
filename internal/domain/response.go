package domain

type Status string

const (
	Unknown       Status = "unknown"
	Start         Status = "start"
	Processing    Status = "processing"
	DownloadOk    Status = "download ok"
	DownloadError Status = "download error"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type StatusResponse struct {
	Status Status `json:"last_download_status" swaggertype:"string" enums:"unknown,start,processing,downloadOk,downloadError" example:"unknown"`
}
