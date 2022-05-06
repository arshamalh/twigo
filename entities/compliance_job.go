package entities

import "time"

type ComplianceType string

const (
	ComplianceTypeTweets ComplianceType = "tweets"
	ComplianceTypeUsers  ComplianceType = "users"
)

type ComplianceJob struct {
	ID                string         `json:"id"`
	Status            string         `json:"status"`
	Name              string         `json:"name"`
	Resumable         bool           `json:"resumable"`
	UploadURL         string         `json:"upload_url"`
	DownloadURL       string         `json:"download_url"`
	Type              ComplianceType `json:"type"`
	CreatedAt         time.Time      `json:"created_at"`
	UploadExpiresAt   time.Time      `json:"upload_expires_at"`
	DownloadExpiresAt time.Time      `json:"download_expires_at"`
}
