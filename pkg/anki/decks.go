package anki

import (
	"context"
)

func (a *Anki) DeckNames(ctx context.Context) ([]string, error) {
	var result []string
	err := a.request(ctx, "deckNames", nil, &result)
	return result, err
}

func (a *Anki) CreateDeck(ctx context.Context, deckName string) (int64, error) {
	request := struct {
		Deck string `json:"deck"`
	}{
		Deck: deckName,
	}
	var result int64
	err := a.request(ctx, "createDeck", &request, &result)
	return result, err
}

func (a *Anki) DeleteDecks(ctx context.Context, deckNames []string) error {
	request := struct {
		Decks    []string `json:"decks"`
		CardsToo bool     `json:"cardsToo"`
	}{
		Decks: deckNames,
		// because Anki since 2.1.28 forbid to delete decks without cards
		CardsToo: true,
	}
	err := a.request(ctx, "deleteDecks", &request, nil)
	return err
}
