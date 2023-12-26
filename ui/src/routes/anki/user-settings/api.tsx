import { gql } from '../../../api/__generated__';

export const GET_ANKI_STATE = gql(`
  query GetAnkiState {
    Anki {
      decks {
        decks
        error {
          __typename
        }
      }
      notes {
        notes
        error {
          __typename
        }
      }
      noteFields {
        noteFields 
        error {
          ... on Error {
            message
          }
        }
      }
    }
  }
`);

export const GET_ANKI_CONFIG = gql(`
  query GetAnkiConfig {
    AnkiConfig {
      deck
      noteType
      mapping {
        key
        value
      }
      audioField
      audioPreferredType
    }
  }
`);
