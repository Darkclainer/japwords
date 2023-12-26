import { gql } from '../../../api/__generated__';

export const GET_CURRENT_NOTE = gql(`
  query GetAnkiConfigCurrentNote {
    AnkiConfig {
      noteType
    }
  }
`);

export const GET_NOTE_FIELDS_AND_MAPPING = gql(`
  query GetAnkiNoteFieldsAndMapping {
    AnkiConfig {
      mapping {
        key
        value
      }
    }
    Anki {
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
