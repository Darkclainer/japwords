package anki

import "context"

func (a *Anki) ModelNames(ctx context.Context) ([]string, error) {
	var result []string
	err := a.request(ctx, "modelNames", nil, &result)
	return result, err
}

func (a *Anki) ModelFieldNames(ctx context.Context, modelName string) ([]string, error) {
	request := struct {
		ModelName string `json:"modelName"`
	}{
		ModelName: modelName,
	}
	var result []string
	err := a.request(ctx, "modelFieldNames", &request, &result)
	return result, err
}

type CreateModelRequest struct {
	ModelName     string                    `json:"modelName,omitempty"`
	Fields        []string                  `json:"inOrderFields,omitempty"`
	CSS           string                    `json:"css,omitempty"`
	CardTemplates []CreateModelCardTemplate `json:"cardTemplates,omitempty"`
}

type CreateModelCardTemplate struct {
	Name  string `json:"Name,omitempty"`
	Front string `json:"Front,omitempty"`
	Back  string `json:"Back,omitempty"`
}

func (a *Anki) CreateModel(ctx context.Context, parameters *CreateModelRequest) (int64, error) {
	// actually there are a lot more of fields in response, but I don't need them
	// and they are poorly documented.
	response := struct {
		ID int64 `json:"id"`
	}{}
	err := a.request(ctx, "createModel", &parameters, &response)
	return response.ID, err
}
