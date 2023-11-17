import { gql } from '../../../api/__generated__';

export const GET_CURRENT_NOTE = gql(`
  query GetAnkiConfigCurrentNote {
    AnkiConfig {
      noteType
    }
  }
`);
