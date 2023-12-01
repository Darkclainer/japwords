package ankiconnect

import "context"

func (a *Anki) FindNotes(ctx context.Context, query string) ([]int64, error) {
	request := struct {
		Query string `json:"query"`
	}{
		Query: query,
	}
	var result []int64
	err := a.request(ctx, "findNotes", &request, &result)
	return result, err
}

type DuplicateScope int

const (
	DuplicateScopeDeck DuplicateScope = iota
	DuplicateScopeEverywhere
)

type DuplicateFlags int

const (
	DuplicateFlagsCheck DuplicateFlags = 1 << iota
	DuplicateFlagsWithChildren
	DuplicateFlagsWithModels
)

type AddNoteOptions struct {
	Deck           string
	Model          string
	DuplicateScope DuplicateScope
	DuplicateFlags DuplicateFlags
	DuplicateDeck  string
}

type AddNoteAsset struct {
	Asset MediaAssetRequest
	Type  MediaType
}

type AddNoteParams struct {
	Tags   []string
	Fields map[string]string
	Assets []*AddNoteAsset
}

// AddNote creates new notes. Params specify note data and opts specifies where it belongs and how it should be
// added. So likely, that many notes can be added with different params and same opts.
func (a *Anki) AddNote(ctx context.Context, params *AddNoteParams, opts *AddNoteOptions) (int64, error) {
	request := addNoteRequestNote{
		DeckName:  opts.Deck,
		ModelName: opts.Model,
		Fields:    params.Fields,
		Options: addNoteRequestOptions{
			AllowDuplicate: opts.DuplicateFlags&DuplicateFlagsCheck == 0,
			DuplicateScope: "",
			DuplicateScopeOptions: addNoteRequestDuplicateOptions{
				DeckName:       opts.DuplicateDeck,
				CheckChildren:  opts.DuplicateFlags&DuplicateFlagsWithChildren != 0,
				CheckAllModels: opts.DuplicateFlags&DuplicateFlagsWithModels != 0,
			},
		},
		Tags: params.Tags,
	}
	switch opts.DuplicateScope {
	case DuplicateScopeDeck:
		request.Options.DuplicateScope = "deck"
	case DuplicateScopeEverywhere:
		// actually this field should be anything except "deck"
		request.Options.DuplicateScope = "all"
	default:
		panic("uknown DuplicateScope value")
	}
	for _, asset := range params.Assets {
		switch asset.Type {
		case MediaTypeAudio:
			request.Audio = append(request.Audio, &asset.Asset)
		case MediaTypeVideo:
			request.Video = append(request.Video, &asset.Asset)
		case MediaTypePicture:
			request.Picture = append(request.Picture, &asset.Asset)
		default:
			panic("unknown asset type")
		}
	}
	realRequest := addNoteRequest{
		Note: request,
	}
	var response int64

	err := a.request(ctx, "addNote", &realRequest, &response)
	if err != nil {
		return 0, err
	}
	// anki-connect specify that response can be null (converted as 0, so we have to cases here, but make no distinction)
	// which means that note wasn't created.
	if response == 0 {
		return 0, newServerError("note creation failed")
	}
	return response, nil
}

type addNoteRequest struct {
	Note addNoteRequestNote `json:"note"`
}

type addNoteRequestNote struct {
	DeckName  string                `json:"deckName,omitempty"`
	ModelName string                `json:"modelName,omitempty"`
	Fields    map[string]string     `json:"fields,omitempty"`
	Options   addNoteRequestOptions `json:"options,omitempty"`
	Tags      []string              `json:"tags,omitempty"`
	Audio     []*MediaAssetRequest  `json:"audio,omitempty"`
	Video     []*MediaAssetRequest  `json:"video,omitempty"`
	Picture   []*MediaAssetRequest  `json:"picture,omitempty"`
}

type addNoteRequestOptions struct {
	AllowDuplicate        bool                           `json:"allowDuplicate,omitempty"`
	DuplicateScope        string                         `json:"duplicateScope,omitempty"`
	DuplicateScopeOptions addNoteRequestDuplicateOptions `json:"duplicateScopeOptions,omitempty"`
}

type addNoteRequestDuplicateOptions struct {
	DeckName       string `json:"deckName,omitempty"`
	CheckChildren  bool   `json:"checkChildren,omitempty"`
	CheckAllModels bool   `json:"checkAllModels,omitempty"`
}

func (a *Anki) NotesInfo(ctx context.Context, ids []int64) ([]*NoteInfo, error) {
	request := struct {
		Notes []int64 `json:"notes"`
	}{
		Notes: ids,
	}
	var result []*NoteInfo
	err := a.request(ctx, "notesInfo", &request, &result)
	return result, err
}

type NoteInfo struct {
	NoteID    int64
	ModelName string
	Tags      []string
	Fields    map[string]*NoteInfoField
}

type NoteInfoField struct {
	Value string
	Order int
}

func (a *Anki) DeleteNotes(ctx context.Context, ids []int64) error {
	request := struct {
		Notes []int64 `json:"notes"`
	}{
		Notes: ids,
	}
	return a.request(ctx, "deleteNotes", &request, nil)
}
