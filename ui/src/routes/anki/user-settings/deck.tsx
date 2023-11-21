import { useMutation, useSuspenseQuery } from '@apollo/client';
import { Label } from '@radix-ui/react-label';
import { useId, useMemo } from 'react';
import { err, ok } from 'true-myth/result';

import { gql } from '../../../api/__generated__';
import { GET_HEALTH_STATUS } from '../../../api/health-status';
import SelectCreate from '../../../components/SelectCreate';
import SuspenseLoading from '../../../components/SuspenseLoading';
import { useToastify } from '../../../hooks/toastify';
import { validateDeck } from '../../../lib/validate';

const GET_CURRENT_DECK = gql(`
  query GetAnkiConfigCurrentDeck {
    AnkiConfig {
      deck
    }
  }
`);

const GET_ANKI_DECKS = gql(`
  query GetAnkiDecks {
    Anki {
      anki {
        decks
      }
      error {
        __typename
      }
    }
  }
`);

const SET_CURRENT_DECK = gql(`
  mutation SetAnkiConfigCurrentDeck($name: String!) {
    setAnkiConfigDeck(input: { name: $name }) {
      error {
          message
      }
    }
  }
`);

const CREATE_DECK = gql(`
  mutation CreateAnkiDeck($name: String!) {
    createAnkiDeck(input: { name: $name }) {
      ankiError {
        ... on Error {
          message
        }
      }
      error {
        ... on CreateAnkiDeckAlreadyExists {
          message
        }
        ... on ValidationError {
          message
        }
        ... on Error {
          message
        }
      
      }
    }
  }
`);

export function DeckSelect() {
  const deckTriggerId = useId();
  return (
    <div className="flex flex-col gap-2.5">
      <Label className="text-2xl" htmlFor={deckTriggerId}>
        Choose a deck:
      </Label>
      <SuspenseLoading>
        <DeckSelectBody triggerId={deckTriggerId} />
      </SuspenseLoading>
    </div>
  );
}

function DeckSelectBody({ triggerId }: { triggerId: string }) {
  const [setCurrentDeck] = useMutation(SET_CURRENT_DECK, {
    refetchQueries: [GET_CURRENT_DECK, GET_HEALTH_STATUS],
    awaitRefetchQueries: true,
  });
  const [createDeck] = useMutation(CREATE_DECK, {
    refetchQueries: [GET_ANKI_DECKS],
    awaitRefetchQueries: true,
  });
  const { data: currentDeckResp } = useSuspenseQuery(GET_CURRENT_DECK, {
    fetchPolicy: 'network-only',
  });
  const currentDeck = currentDeckResp.AnkiConfig.deck;
  const { data: decksResp } = useSuspenseQuery(GET_ANKI_DECKS, {
    fetchPolicy: 'no-cache',
  });
  const decks = useMemo(() => {
    if (!decksResp.Anki.anki) {
      return null;
    }
    const decks = [...decksResp.Anki.anki.decks];
    return decks.sort().map((item) => {
      return {
        value: item,
      };
    });
  }, [decksResp]);
  const toast = useToastify({
    type: 'success',
  });
  if (!decks) {
    // this error handled in parent components
    return null;
  }
  const currentDeckExists = !decks.find((e) => e.value == currentDeck);
  return (
    <>
      <SelectCreate
        id={triggerId}
        triggerClassName="max-w-md shrink"
        hasError={currentDeckExists}
        items={decks}
        selectedValue={currentDeck}
        onValueChange={async (value: string) => {
          const resp = await setCurrentDeck({
            variables: {
              name: value,
            },
          });
          if (!resp.data || resp.data.setAnkiConfigDeck.error) {
            toast('Deck change failed!', { type: 'error' });
          } else {
            toast('Deck successfully changed.');
          }
        }}
        handleCreate={async (value: string) => {
          const resp = await createDeck({
            variables: {
              name: value,
            },
          });
          if (!resp.data) {
            toast('Deck creation failed!', { type: 'error' });
            return err('request failed');
          }
          if (resp.data.createAnkiDeck.ankiError) {
            toast('Deck creation failed! No anki connection', { type: 'error' });
            return err('request failed');
          }
          if (resp.data.createAnkiDeck.error) {
            const error = resp.data.createAnkiDeck.error;
            switch (error.__typename) {
              case 'ValidationError':
                return err(error.message);
              case 'CreateAnkiDeckAlreadyExists':
                return err('deck with specified name already exists');
              default:
                return err('uknown error');
            }
          }
          return ok(value);
        }}
        validateValue={validateDeck}
        placeholderLabel="Select or create..."
        createLabel="Create new deck"
        createDefaultValue="Japwords"
        dialogTitle="Create deck"
        dialogInputLabel="Input new deck name"
      />
      {currentDeckExists && <p className="text-error-red text-lg">Selected deck does not exists</p>}
    </>
  );
}
