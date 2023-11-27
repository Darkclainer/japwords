import { ApolloError } from '@apollo/client';

import { AnkiConfigStateResult, GetHealthStatusQuery } from '../api/__generated__/graphql';

export type HealthStatus =
  | HealthStatusDisconnected
  | HealthStatusLoading
  | HealthStatusError
  | HealthStatusOk;

export interface HealthStatusDisconnected {
  kind: 'Disconnected';
}

export interface HealthStatusLoading {
  kind: 'Loading';
}

// HealthStatusError it's unknown error, when we connected to server, but get unknown error response
export interface HealthStatusError {
  kind: 'Error';
  message?: string;
}

export interface HealthStatusOk {
  kind: 'Ok';
  anki: AnkiState;
}

export interface HealthStatusThrownError {
  kind: 'Error' | 'Disconnected';
  __typename: 'HealthStatusThrownError';
}

// eslint-disable-next-line  @typescript-eslint/no-explicit-any
export function isHealthStatusThrownError(err: any): err is HealthStatusThrownError {
  return (
    typeof err == 'object' &&
    typeof err.kind == 'string' &&
    typeof err.__typename == 'string' &&
    err.__typename == 'HealthStatusThrownError'
  );
}

export function throwErrorHealthStatus(
  status: HealthStatus,
): asserts status is HealthStatusLoading | HealthStatusOk {
  switch (status.kind) {
    case 'Ok':
    case 'Loading':
      return;
    case 'Error':
    case 'Disconnected': {
      const err: HealthStatusThrownError = {
        kind: status.kind,
        __typename: 'HealthStatusThrownError',
      };
      throw err;
    }
    default: {
      const _exhaustiveCheck: never = status;
      return _exhaustiveCheck;
    }
  }
}

export type AnkiState =
  | AnkiStateConnectionError
  | AnkiStateUnknownError
  | AnkiStateForbiddenOrigin
  | AnkiStateCollectionUnavailable
  | AnkiStateInvalidAPIKey
  | AnkiStateOk;

// AnkiStateConnectionError means that backend tried to connect to anki-connect but failed
export interface AnkiStateConnectionError {
  kind: 'ConnectionError';
  message?: string;
}

export interface AnkiStateUnknownError {
  kind: 'UnknownError';
  message?: string;
}

export interface AnkiStateForbiddenOrigin {
  kind: 'ForbiddenOrigin';
}

export interface AnkiStateCollectionUnavailable {
  kind: 'CollectionUnavailable';
  version: number;
}

// AnkiStatePermissionError
export interface AnkiStateInvalidAPIKey {
  kind: 'InvalidAPIKey';
  version: number;
}

export interface AnkiStateOk {
  kind: 'Ok' | 'UserError';
  version: number;
  deckExists: boolean;
  noteTypeExists: boolean;
  noteHasAllFields: boolean;
}

export function healthStatusFromGql(
  data?: GetHealthStatusQuery | null,
  error?: ApolloError | null,
): HealthStatus {
  // it's best to add some "uknown error" in case we got error from server unexpectedly
  if (error || !data) {
    if (error?.networkError && !('statusCode' in error.networkError)) {
      return {
        kind: 'Disconnected',
      };
    } else {
      return {
        kind: 'Error',
        message: error?.message,
      };
    }
  }
  const ankiState = ankiStateFromGql(data.AnkiConfigState);
  return {
    kind: 'Ok',
    anki: ankiState,
  };
}

function ankiStateFromGql(state: AnkiConfigStateResult): AnkiState {
  if (state.error) {
    switch (state.error.__typename) {
      case 'AnkiConnectionError':
        return {
          kind: 'ConnectionError',
          message: state.error.message,
        };
      case 'AnkiForbiddenOrigin':
        return {
          kind: 'ForbiddenOrigin',
        };
      case 'AnkiInvalidAPIKey':
        return {
          kind: 'InvalidAPIKey',
          version: state.error.version,
        };
      case 'AnkiCollectionUnavailable':
        return {
          kind: 'CollectionUnavailable',
          version: state.error.version,
        };
      case 'AnkiUnknownError':
        return {
          kind: 'UnknownError',
        };
      case undefined:
        throw 'unreachable';
      default: {
        const _exhaustiveCheck: never = state.error;
        return _exhaustiveCheck;
      }
    }
  }
  if (state.ankiConfigState) {
    const confState = state.ankiConfigState;
    return {
      kind:
        confState.deckExists && confState.noteTypeExists && confState.noteHasAllFields
          ? 'Ok'
          : 'UserError',
      version: confState.version,
      deckExists: confState.deckExists,
      noteTypeExists: confState.noteTypeExists,
      noteHasAllFields: confState.noteHasAllFields,
    };
  }
  throw 'unreachable';
}
