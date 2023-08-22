package ankiconnect

import (
	"encoding/base64"
	"fmt"
	"slices"
)

type MediaType int

const (
	MediaTypeAudio MediaType = iota
	MediaTypePicture
	MediaTypeVideo
)

// MediaAssetOptions represents common options about how asset should be saved.
type MediaAssetOptions struct {
	Filename       string
	DeleteExisting bool
	Fields         []string
	SkipHash       []byte
}

func (o *MediaAssetOptions) applyToRequest(r *MediaAssetRequest) {
	r.Filename = o.Filename
	r.DeleteExisting = o.DeleteExisting
	r.Fields = slices.Clone(o.Fields)
	r.SkipHash = fmt.Sprintf("%x", o.SkipHash)
	r.DeleteExisting = o.DeleteExisting
}

func NewMediaFile(filename string, opts *MediaAssetOptions) *MediaAssetRequest {
	asset := &MediaAssetRequest{
		Path: filename,
	}
	opts.applyToRequest(asset)
	return asset
}

func NewMediaURL(url string, opts *MediaAssetOptions) *MediaAssetRequest {
	asset := &MediaAssetRequest{
		URL: url,
	}
	opts.applyToRequest(asset)
	return asset
}

func NewMediaBlob(data []byte, opts *MediaAssetOptions) *MediaAssetRequest {
	encodedLen := base64.StdEncoding.EncodedLen(len(data))
	buffer := make([]byte, encodedLen)
	base64.StdEncoding.Encode(buffer, data)
	asset := &MediaAssetRequest{
		Data: string(buffer),
	}
	opts.applyToRequest(asset)
	return asset
}

// MediaAssetRequest represents actual request that will go to anki-connect
type MediaAssetRequest struct {
	Data           string   `json:"data,omitempty"`
	Path           string   `json:"path,omitempty"`
	URL            string   `json:"url,omitempty"`
	Filename       string   `json:"filename,omitempty"`
	Fields         []string `json:"fields,omitempty"`
	SkipHash       string   `json:"skipHash,omitempty"`
	DeleteExisting bool     `json:"deleteExisting,omitempty"`
}
