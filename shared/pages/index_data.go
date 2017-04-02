package pages

// IndexData DTO
type IndexData struct {
	PageData
	ShowWelcome     bool
	FollowerCount   int
	StarCount       int
	RepositoryCount int
}
