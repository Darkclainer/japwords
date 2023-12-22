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
        orderDefined
        audioFieldExists
      }
      error {
        ... on AnkiConnectionError {
          message
        }
        ... on AnkiInvalidAPIKey {
          version
          message
        }
        ... on AnkiCollectionUnavailable {
          version
          message
        }
        ... on AnkiForbiddenOrigin {
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
