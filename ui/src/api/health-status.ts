import { gql } from '../api/__generated__/gql';
import apolloClient from '../apollo-client';
import { GetHealthStatusQuery } from './__generated__/graphql';

export const GET_HEALTH_STATUS = gql(`
  query GetHealthStatus {
    AnkiConfigState {
      ankiConfigState {
        version
        deckExists
        noteTypeExists
        noteHasAllFields
      }
      error {
        ... on AnkiConnectionError {
          message
        }
        ... on AnkiPermissionError {
          version
          message
        }
        ... on AnkiUnauthorizedError {
          version
          message
        }
        ... on Error {
          message
        }
      }
    }
  }
`);

export function cachedHealthStatus(): GetHealthStatusQuery | null {
  return apolloClient.readQuery({ query: GET_HEALTH_STATUS });
}
