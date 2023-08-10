package models

type Mode string

const (
	ModeInMemory Mode = "in-memory"
	ModeMongo    Mode = "mongo"
	ModeCached   Mode = "cached"
)

type UrlItem struct {
	Key string `bson:"_id"`
	URL string `bson:"url"`
}
