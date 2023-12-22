/* eslint-disable */
import * as types from './graphql';
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';

/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel or swc plugin for production.
 */
const documents = {
    "\n  query GetHealthStatus {\n    AnkiConfigState {\n      ankiConfigState {\n        version\n        deckExists\n        noteTypeExists\n        noteHasAllFields\n        orderDefined\n        audioFieldExists\n      }\n      error {\n        ... on AnkiConnectionError {\n          message\n        }\n        ... on AnkiInvalidAPIKey {\n          version\n          message\n        }\n        ... on AnkiCollectionUnavailable {\n          version\n          message\n        }\n        ... on AnkiForbiddenOrigin {\n          message\n        }\n        ... on Error {\n          message\n        }\n      }\n    }\n  }\n": types.GetHealthStatusDocument,
    "\n  query GetConnectionConfig {\n    AnkiConfig {\n      addr\n      apiKey\n    }\n  }\n": types.GetConnectionConfigDocument,
    "\n  mutation UpdateConnectionConfig($addr: String!, $apiKey: String!) {\n    setAnkiConfigConnection(input: { addr: $addr, apiKey: $apiKey }) {\n      error {\n          ... on ValidationError {\n            paths\n            message\n          }\n        }\n      }\n  }\n": types.UpdateConnectionConfigDocument,
    "\n  query GetAnkiConfigCurrentNote {\n    AnkiConfig {\n      noteType\n    }\n  }\n": types.GetAnkiConfigCurrentNoteDocument,
    "\n  query GetAnkiNoteFieldsAndMapping($noteName: String!) {\n    AnkiConfig {\n      mapping {\n        key\n        value\n      }\n    }\n    Anki {\n      noteFields(name: $noteName) {\n        noteFields \n        error {\n          ... on Error {\n            message\n          }\n        }\n      }\n    }\n  }\n": types.GetAnkiNoteFieldsAndMappingDocument,
    "\n  query GetAnkiConfigCurrentDeck {\n    AnkiConfig {\n      deck\n    }\n  }\n": types.GetAnkiConfigCurrentDeckDocument,
    "\n  query GetAnkiDecks {\n    Anki {\n      decks {\n        decks\n        error {\n          __typename\n        }\n      }\n    }\n  }\n": types.GetAnkiDecksDocument,
    "\n  mutation SetAnkiConfigCurrentDeck($name: String!) {\n    setAnkiConfigDeck(input: { name: $name }) {\n      error {\n          message\n      }\n    }\n  }\n": types.SetAnkiConfigCurrentDeckDocument,
    "\n  mutation CreateAnkiDeck($name: String!) {\n    createAnkiDeck(input: { name: $name }) {\n      ankiError {\n        ... on Error {\n          message\n        }\n      }\n      error {\n        ... on CreateAnkiDeckAlreadyExists {\n          message\n        }\n        ... on ValidationError {\n          message\n        }\n        ... on Error {\n          message\n        }\n      \n      }\n    }\n  }\n": types.CreateAnkiDeckDocument,
    "\n  query RenderFields($fields: [String!]!) {\n    RenderFields(fields: $fields) {\n      fields {\n        result\n        error\n      }\n    }\n  }\n": types.RenderFieldsDocument,
    "\n  mutation UpdateMapping($fields: [AnkiConfigMappingElementInput!]!) {\n    setAnkiConfigMapping(input: { mapping: $fields }) {\n      error {\n        fieldErrors {\n          key\n        }\n        valueErrors {\n          key\n        }\n        message\n      }\n    }\n  }\n": types.UpdateMappingDocument,
    "\n  query GetAnkiNotes {\n    Anki {\n      notes {\n        notes\n        error {\n          __typename\n        }\n      }\n    }\n  }\n": types.GetAnkiNotesDocument,
    "\n  mutation SetAnkiConfigCurrentNote($name: String!) {\n    setAnkiConfigNote(input: { name: $name }) {\n      error {\n          message\n      }\n    }\n  }\n": types.SetAnkiConfigCurrentNoteDocument,
    "\n  mutation CreateDefaultAnkiNote($name: String!) {\n    createDefaultAnkiNote(input: { name: $name }) {\n      ankiError {\n        ... on Error {\n          message\n        }\n      }\n      error {\n        ... on CreateDefaultAnkiNoteAlreadyExists {\n          message\n        }\n        ... on ValidationError {\n          message\n        }\n        ... on Error {\n          message\n        }\n      \n      }\n    }\n  }\n": types.CreateDefaultAnkiNoteDocument,
    "\nmutation AddAnkiNote($note: AddNoteRequestInput!) {\n  addAnkiNote(request: $note) {\n    noteID\n    error {\n      ... on AnkiIncompleteConfiguration {\n        message\n      }\n      ... on AnkiAddNoteDuplicateFound {\n        message\n      }\n    }\n    ankiError {\n      ... on Error {\n        message\n      }\n    }\n  }\n}\n": types.AddAnkiNoteDocument,
    "\nquery PrepareLemma($lemma: LemmaInput) {\n  PrepareLemma(lemma: $lemma) {\n    request {\n      fields {\n        name\n        value\n      }\n      tags\n      audioAssets {\n        field\n        filename\n        url\n        data\n      }\n    }\n    error {\n      ... on AnkiIncompleteConfiguration {\n        message\n      }\n      ... on Error {\n        message\n      }\n    }\n    ankiError {\n      __typename\n      ... on Error {\n        message\n      }\n    }\n  }\n}\n": types.PrepareLemmaDocument,
    "\n  query GetLemmas($query: String!) {\n    Lemmas(query: $query) {\n      lemmas {\n        noteID\n        lemma {\n          slug {\n            word\n            hiragana\n            furigana {\n              kanji\n              hiragana\n            }\n            pitchShapes {\n              hiragana\n              directions\n            } \n          }\n          tags\n          forms {\n            word\n            hiragana\n            furigana {\n              kanji\n              hiragana\n            }\n            pitchShapes {\n              hiragana\n              directions\n            } \n          }\n          definitions\n          partsOfSpeech\n          senseTags\n          audio {\n            mediaType\n            source\n          }\n        } \n      }\n    }\n  }\n": types.GetLemmasDocument,
};

/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = gql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function gql(source: string): unknown;

/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query GetHealthStatus {\n    AnkiConfigState {\n      ankiConfigState {\n        version\n        deckExists\n        noteTypeExists\n        noteHasAllFields\n        orderDefined\n        audioFieldExists\n      }\n      error {\n        ... on AnkiConnectionError {\n          message\n        }\n        ... on AnkiInvalidAPIKey {\n          version\n          message\n        }\n        ... on AnkiCollectionUnavailable {\n          version\n          message\n        }\n        ... on AnkiForbiddenOrigin {\n          message\n        }\n        ... on Error {\n          message\n        }\n      }\n    }\n  }\n"): (typeof documents)["\n  query GetHealthStatus {\n    AnkiConfigState {\n      ankiConfigState {\n        version\n        deckExists\n        noteTypeExists\n        noteHasAllFields\n        orderDefined\n        audioFieldExists\n      }\n      error {\n        ... on AnkiConnectionError {\n          message\n        }\n        ... on AnkiInvalidAPIKey {\n          version\n          message\n        }\n        ... on AnkiCollectionUnavailable {\n          version\n          message\n        }\n        ... on AnkiForbiddenOrigin {\n          message\n        }\n        ... on Error {\n          message\n        }\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query GetConnectionConfig {\n    AnkiConfig {\n      addr\n      apiKey\n    }\n  }\n"): (typeof documents)["\n  query GetConnectionConfig {\n    AnkiConfig {\n      addr\n      apiKey\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation UpdateConnectionConfig($addr: String!, $apiKey: String!) {\n    setAnkiConfigConnection(input: { addr: $addr, apiKey: $apiKey }) {\n      error {\n          ... on ValidationError {\n            paths\n            message\n          }\n        }\n      }\n  }\n"): (typeof documents)["\n  mutation UpdateConnectionConfig($addr: String!, $apiKey: String!) {\n    setAnkiConfigConnection(input: { addr: $addr, apiKey: $apiKey }) {\n      error {\n          ... on ValidationError {\n            paths\n            message\n          }\n        }\n      }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query GetAnkiConfigCurrentNote {\n    AnkiConfig {\n      noteType\n    }\n  }\n"): (typeof documents)["\n  query GetAnkiConfigCurrentNote {\n    AnkiConfig {\n      noteType\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query GetAnkiNoteFieldsAndMapping($noteName: String!) {\n    AnkiConfig {\n      mapping {\n        key\n        value\n      }\n    }\n    Anki {\n      noteFields(name: $noteName) {\n        noteFields \n        error {\n          ... on Error {\n            message\n          }\n        }\n      }\n    }\n  }\n"): (typeof documents)["\n  query GetAnkiNoteFieldsAndMapping($noteName: String!) {\n    AnkiConfig {\n      mapping {\n        key\n        value\n      }\n    }\n    Anki {\n      noteFields(name: $noteName) {\n        noteFields \n        error {\n          ... on Error {\n            message\n          }\n        }\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query GetAnkiConfigCurrentDeck {\n    AnkiConfig {\n      deck\n    }\n  }\n"): (typeof documents)["\n  query GetAnkiConfigCurrentDeck {\n    AnkiConfig {\n      deck\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query GetAnkiDecks {\n    Anki {\n      decks {\n        decks\n        error {\n          __typename\n        }\n      }\n    }\n  }\n"): (typeof documents)["\n  query GetAnkiDecks {\n    Anki {\n      decks {\n        decks\n        error {\n          __typename\n        }\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation SetAnkiConfigCurrentDeck($name: String!) {\n    setAnkiConfigDeck(input: { name: $name }) {\n      error {\n          message\n      }\n    }\n  }\n"): (typeof documents)["\n  mutation SetAnkiConfigCurrentDeck($name: String!) {\n    setAnkiConfigDeck(input: { name: $name }) {\n      error {\n          message\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation CreateAnkiDeck($name: String!) {\n    createAnkiDeck(input: { name: $name }) {\n      ankiError {\n        ... on Error {\n          message\n        }\n      }\n      error {\n        ... on CreateAnkiDeckAlreadyExists {\n          message\n        }\n        ... on ValidationError {\n          message\n        }\n        ... on Error {\n          message\n        }\n      \n      }\n    }\n  }\n"): (typeof documents)["\n  mutation CreateAnkiDeck($name: String!) {\n    createAnkiDeck(input: { name: $name }) {\n      ankiError {\n        ... on Error {\n          message\n        }\n      }\n      error {\n        ... on CreateAnkiDeckAlreadyExists {\n          message\n        }\n        ... on ValidationError {\n          message\n        }\n        ... on Error {\n          message\n        }\n      \n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query RenderFields($fields: [String!]!) {\n    RenderFields(fields: $fields) {\n      fields {\n        result\n        error\n      }\n    }\n  }\n"): (typeof documents)["\n  query RenderFields($fields: [String!]!) {\n    RenderFields(fields: $fields) {\n      fields {\n        result\n        error\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation UpdateMapping($fields: [AnkiConfigMappingElementInput!]!) {\n    setAnkiConfigMapping(input: { mapping: $fields }) {\n      error {\n        fieldErrors {\n          key\n        }\n        valueErrors {\n          key\n        }\n        message\n      }\n    }\n  }\n"): (typeof documents)["\n  mutation UpdateMapping($fields: [AnkiConfigMappingElementInput!]!) {\n    setAnkiConfigMapping(input: { mapping: $fields }) {\n      error {\n        fieldErrors {\n          key\n        }\n        valueErrors {\n          key\n        }\n        message\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query GetAnkiNotes {\n    Anki {\n      notes {\n        notes\n        error {\n          __typename\n        }\n      }\n    }\n  }\n"): (typeof documents)["\n  query GetAnkiNotes {\n    Anki {\n      notes {\n        notes\n        error {\n          __typename\n        }\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation SetAnkiConfigCurrentNote($name: String!) {\n    setAnkiConfigNote(input: { name: $name }) {\n      error {\n          message\n      }\n    }\n  }\n"): (typeof documents)["\n  mutation SetAnkiConfigCurrentNote($name: String!) {\n    setAnkiConfigNote(input: { name: $name }) {\n      error {\n          message\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation CreateDefaultAnkiNote($name: String!) {\n    createDefaultAnkiNote(input: { name: $name }) {\n      ankiError {\n        ... on Error {\n          message\n        }\n      }\n      error {\n        ... on CreateDefaultAnkiNoteAlreadyExists {\n          message\n        }\n        ... on ValidationError {\n          message\n        }\n        ... on Error {\n          message\n        }\n      \n      }\n    }\n  }\n"): (typeof documents)["\n  mutation CreateDefaultAnkiNote($name: String!) {\n    createDefaultAnkiNote(input: { name: $name }) {\n      ankiError {\n        ... on Error {\n          message\n        }\n      }\n      error {\n        ... on CreateDefaultAnkiNoteAlreadyExists {\n          message\n        }\n        ... on ValidationError {\n          message\n        }\n        ... on Error {\n          message\n        }\n      \n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\nmutation AddAnkiNote($note: AddNoteRequestInput!) {\n  addAnkiNote(request: $note) {\n    noteID\n    error {\n      ... on AnkiIncompleteConfiguration {\n        message\n      }\n      ... on AnkiAddNoteDuplicateFound {\n        message\n      }\n    }\n    ankiError {\n      ... on Error {\n        message\n      }\n    }\n  }\n}\n"): (typeof documents)["\nmutation AddAnkiNote($note: AddNoteRequestInput!) {\n  addAnkiNote(request: $note) {\n    noteID\n    error {\n      ... on AnkiIncompleteConfiguration {\n        message\n      }\n      ... on AnkiAddNoteDuplicateFound {\n        message\n      }\n    }\n    ankiError {\n      ... on Error {\n        message\n      }\n    }\n  }\n}\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\nquery PrepareLemma($lemma: LemmaInput) {\n  PrepareLemma(lemma: $lemma) {\n    request {\n      fields {\n        name\n        value\n      }\n      tags\n      audioAssets {\n        field\n        filename\n        url\n        data\n      }\n    }\n    error {\n      ... on AnkiIncompleteConfiguration {\n        message\n      }\n      ... on Error {\n        message\n      }\n    }\n    ankiError {\n      __typename\n      ... on Error {\n        message\n      }\n    }\n  }\n}\n"): (typeof documents)["\nquery PrepareLemma($lemma: LemmaInput) {\n  PrepareLemma(lemma: $lemma) {\n    request {\n      fields {\n        name\n        value\n      }\n      tags\n      audioAssets {\n        field\n        filename\n        url\n        data\n      }\n    }\n    error {\n      ... on AnkiIncompleteConfiguration {\n        message\n      }\n      ... on Error {\n        message\n      }\n    }\n    ankiError {\n      __typename\n      ... on Error {\n        message\n      }\n    }\n  }\n}\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query GetLemmas($query: String!) {\n    Lemmas(query: $query) {\n      lemmas {\n        noteID\n        lemma {\n          slug {\n            word\n            hiragana\n            furigana {\n              kanji\n              hiragana\n            }\n            pitchShapes {\n              hiragana\n              directions\n            } \n          }\n          tags\n          forms {\n            word\n            hiragana\n            furigana {\n              kanji\n              hiragana\n            }\n            pitchShapes {\n              hiragana\n              directions\n            } \n          }\n          definitions\n          partsOfSpeech\n          senseTags\n          audio {\n            mediaType\n            source\n          }\n        } \n      }\n    }\n  }\n"): (typeof documents)["\n  query GetLemmas($query: String!) {\n    Lemmas(query: $query) {\n      lemmas {\n        noteID\n        lemma {\n          slug {\n            word\n            hiragana\n            furigana {\n              kanji\n              hiragana\n            }\n            pitchShapes {\n              hiragana\n              directions\n            } \n          }\n          tags\n          forms {\n            word\n            hiragana\n            furigana {\n              kanji\n              hiragana\n            }\n            pitchShapes {\n              hiragana\n              directions\n            } \n          }\n          definitions\n          partsOfSpeech\n          senseTags\n          audio {\n            mediaType\n            source\n          }\n        } \n      }\n    }\n  }\n"];

export function gql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;