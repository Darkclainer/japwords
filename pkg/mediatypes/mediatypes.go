package mediatypes

// some types from https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
var mediaTypeToExtension = map[string]string{
	"audio/3gpp":  ".3gp",
	"audio/3gpp2": ".3g2",
	"audio/aac":   ".aac",
	"audio/mp4":   ".mp4",
	"audio/mpeg":  ".mp3",
	"audio/ogg":   ".oga",
	"audio/opus":  ".opus",
	"audio/wav":   ".wav",
	"audio/webm":  ".weba",
}

func GetExtensionByMediaType(mediaType string) string {
	return mediaTypeToExtension[mediaType]
}
