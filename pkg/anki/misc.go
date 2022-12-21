package anki

import "context"

func (a *Anki) Version(ctx context.Context) (int, error) {
	var result int
	err := a.request(ctx, "version", nil, &result)
	if err != nil {
		return 0, err
	}
	return result, nil
}
