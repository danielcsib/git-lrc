package network

type ReleasePlatformArtifact struct {
	Binary     string `json:"binary"`
	SHA256Sums string `json:"sha256sums"`
	SHA256     string `json:"sha256"`
}

type ReleaseManifestVersion struct {
	Platforms map[string]ReleasePlatformArtifact `json:"platforms"`
}

type ReleaseManifest struct {
	SchemaVersion int                               `json:"schema_version"`
	GeneratedAt   string                            `json:"generated_at"`
	LatestVersion string                            `json:"latest_version"`
	Bucket        string                            `json:"bucket"`
	Prefix        string                            `json:"prefix"`
	DownloadBase  string                            `json:"download_base"`
	Releases      map[string]ReleaseManifestVersion `json:"releases"`
}
