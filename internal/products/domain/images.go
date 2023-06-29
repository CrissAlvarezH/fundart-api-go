package domain

type ImageID int
type ImageTag string

type Image struct {
	ID            ImageID
	Path          string
	PromptEnglish string
	PromptSpanish string
	ThumbnailPath string
	OrderPriority int
	Tags          []ImageTag
}
