/* eslint-disable */
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
};

export enum AccentDirection {
  Down = 'DOWN',
  Left = 'LEFT',
  Right = 'RIGHT',
  Up = 'UP'
}

export type AddNoteField = {
  __typename?: 'AddNoteField';
  name: Scalars['String']['output'];
  value: Scalars['String']['output'];
};

export type AddNoteFieldInput = {
  name: Scalars['String']['input'];
  value: Scalars['String']['input'];
};

export type AddNoteRequest = {
  __typename?: 'AddNoteRequest';
  audioURL: Scalars['String']['output'];
  fields: Array<AddNoteField>;
  tags: Array<Scalars['String']['output']>;
};

export type AddNoteRequestInput = {
  audioURL: Scalars['String']['input'];
  fields: Array<AddNoteFieldInput>;
  tags: Array<Scalars['String']['input']>;
};

export type Anki = {
  __typename?: 'Anki';
  decks: AnkiDecksResult;
  noteFields: AnkiNoteFieldsResult;
  notes: AnkiNotesResult;
};


export type AnkiNoteFieldsArgs = {
  name: Scalars['String']['input'];
};

export type AnkiAddNoteDuplicateFound = Error & {
  __typename?: 'AnkiAddNoteDuplicateFound';
  message: Scalars['String']['output'];
};

export type AnkiAddNoteError = AnkiAddNoteDuplicateFound | AnkiIncompleteConfiguration;

export type AnkiAddNoteResult = {
  __typename?: 'AnkiAddNoteResult';
  ankiError?: Maybe<AnkiError>;
  error?: Maybe<AnkiAddNoteError>;
};

export type AnkiCollectionUnavailable = Error & {
  __typename?: 'AnkiCollectionUnavailable';
  message: Scalars['String']['output'];
  version: Scalars['Int']['output'];
};

export type AnkiConfig = {
  __typename?: 'AnkiConfig';
  addr: Scalars['String']['output'];
  apiKey: Scalars['String']['output'];
  deck: Scalars['String']['output'];
  mapping: Array<AnkiMappingElement>;
  noteType: Scalars['String']['output'];
};

export type AnkiConfigMappingElementError = {
  __typename?: 'AnkiConfigMappingElementError';
  key: Scalars['String']['output'];
  message: Scalars['String']['output'];
};

export type AnkiConfigMappingElementInput = {
  key: Scalars['String']['input'];
  value: Scalars['String']['input'];
};

export type AnkiConfigMappingError = Error & {
  __typename?: 'AnkiConfigMappingError';
  fieldErrors?: Maybe<Array<AnkiConfigMappingElementError>>;
  message: Scalars['String']['output'];
  valueErrors?: Maybe<Array<AnkiConfigMappingElementError>>;
};

export type AnkiConfigState = {
  __typename?: 'AnkiConfigState';
  deckExists: Scalars['Boolean']['output'];
  noteHasAllFields: Scalars['Boolean']['output'];
  noteTypeExists: Scalars['Boolean']['output'];
  orderDefined: Scalars['Boolean']['output'];
  version: Scalars['Int']['output'];
};

export type AnkiConfigStateResult = {
  __typename?: 'AnkiConfigStateResult';
  ankiConfigState?: Maybe<AnkiConfigState>;
  error?: Maybe<AnkiError>;
};

export type AnkiConnectionError = Error & {
  __typename?: 'AnkiConnectionError';
  message: Scalars['String']['output'];
};

export type AnkiDecksResult = {
  __typename?: 'AnkiDecksResult';
  decks?: Maybe<Array<Scalars['String']['output']>>;
  error?: Maybe<AnkiError>;
};

export type AnkiError = AnkiCollectionUnavailable | AnkiConnectionError | AnkiForbiddenOrigin | AnkiInvalidApiKey | AnkiUnknownError;

export type AnkiForbiddenOrigin = Error & {
  __typename?: 'AnkiForbiddenOrigin';
  message: Scalars['String']['output'];
};

export type AnkiIncompleteConfiguration = Error & {
  __typename?: 'AnkiIncompleteConfiguration';
  message: Scalars['String']['output'];
};

export type AnkiInvalidApiKey = Error & {
  __typename?: 'AnkiInvalidAPIKey';
  message: Scalars['String']['output'];
  version: Scalars['Int']['output'];
};

export type AnkiMappingElement = {
  __typename?: 'AnkiMappingElement';
  key: Scalars['String']['output'];
  value: Scalars['String']['output'];
};

export type AnkiNoteFieldsResult = {
  __typename?: 'AnkiNoteFieldsResult';
  error?: Maybe<AnkiError>;
  noteFields?: Maybe<Array<Scalars['String']['output']>>;
};

export type AnkiNotesResult = {
  __typename?: 'AnkiNotesResult';
  error?: Maybe<AnkiError>;
  notes?: Maybe<Array<Scalars['String']['output']>>;
};

export type AnkiUnknownError = Error & {
  __typename?: 'AnkiUnknownError';
  message: Scalars['String']['output'];
};

export type Audio = {
  __typename?: 'Audio';
  source: Scalars['String']['output'];
  type: Scalars['String']['output'];
};

export type AudioInput = {
  source: Scalars['String']['input'];
  type: Scalars['String']['input'];
};

export type CreateAnkiDeckAlreadyExists = Error & {
  __typename?: 'CreateAnkiDeckAlreadyExists';
  message: Scalars['String']['output'];
};

export type CreateAnkiDeckError = CreateAnkiDeckAlreadyExists | ValidationError;

export type CreateAnkiDeckInput = {
  name: Scalars['String']['input'];
};

export type CreateAnkiDeckResult = {
  __typename?: 'CreateAnkiDeckResult';
  ankiError?: Maybe<AnkiError>;
  error?: Maybe<CreateAnkiDeckError>;
};

export type CreateDefaultAnkiNoteAlreadyExists = Error & {
  __typename?: 'CreateDefaultAnkiNoteAlreadyExists';
  message: Scalars['String']['output'];
};

export type CreateDefaultAnkiNoteError = CreateDefaultAnkiNoteAlreadyExists | ValidationError;

export type CreateDefaultAnkiNoteInput = {
  name: Scalars['String']['input'];
};

export type CreateDefaultAnkiNoteResult = {
  __typename?: 'CreateDefaultAnkiNoteResult';
  ankiError?: Maybe<AnkiError>;
  error?: Maybe<CreateDefaultAnkiNoteError>;
};

export type Error = {
  message: Scalars['String']['output'];
};

export type Furigana = {
  __typename?: 'Furigana';
  hiragana: Scalars['String']['output'];
  kanji: Scalars['String']['output'];
};

export type FuriganaInput = {
  hiragana: Scalars['String']['input'];
  kanji: Scalars['String']['input'];
};

export type Lemma = {
  __typename?: 'Lemma';
  audio: Array<Audio>;
  definitions: Array<Scalars['String']['output']>;
  forms: Array<Word>;
  partsOfSpeech: Array<Scalars['String']['output']>;
  senseTags: Array<Scalars['String']['output']>;
  slug: Word;
  tags: Array<Scalars['String']['output']>;
};

export type LemmaInput = {
  audio: Array<AudioInput>;
  definitions: Array<Scalars['String']['input']>;
  forms: Array<WordInput>;
  partsOfSpeech: Array<Scalars['String']['input']>;
  senseTags: Array<Scalars['String']['input']>;
  slug: WordInput;
  tags: Array<Scalars['String']['input']>;
};

export type LemmaNoteInfo = {
  __typename?: 'LemmaNoteInfo';
  lemma: Lemma;
  noteID: Scalars['String']['output'];
};

export type LemmasResult = {
  __typename?: 'LemmasResult';
  lemmas: Array<LemmaNoteInfo>;
};

export type Mutation = {
  __typename?: 'Mutation';
  addAnkiNote: AnkiAddNoteResult;
  createAnkiDeck: CreateAnkiDeckResult;
  createDefaultAnkiNote: CreateDefaultAnkiNoteResult;
  setAnkiConfigConnection: SetAnkiConfigConnectionResult;
  setAnkiConfigDeck: SetAnkiConfigDeckResult;
  setAnkiConfigMapping: SetAnkiConfigMappingResult;
  setAnkiConfigNote: SetAnkiConfigNoteResult;
};


export type MutationAddAnkiNoteArgs = {
  request?: InputMaybe<AddNoteRequestInput>;
};


export type MutationCreateAnkiDeckArgs = {
  input?: InputMaybe<CreateAnkiDeckInput>;
};


export type MutationCreateDefaultAnkiNoteArgs = {
  input?: InputMaybe<CreateDefaultAnkiNoteInput>;
};


export type MutationSetAnkiConfigConnectionArgs = {
  input: SetAnkiConfigConnectionInput;
};


export type MutationSetAnkiConfigDeckArgs = {
  input: SetAnkiConfigDeckInput;
};


export type MutationSetAnkiConfigMappingArgs = {
  input: SetAnkiConfigMappingInput;
};


export type MutationSetAnkiConfigNoteArgs = {
  input: SetAnkiConfigNote;
};

export type PitchShape = {
  __typename?: 'PitchShape';
  directions: Array<AccentDirection>;
  hiragana: Scalars['String']['output'];
};

export type PitchShapeInput = {
  directions: Array<AccentDirection>;
  hiragana: Scalars['String']['input'];
};

export type PrepareLemmaError = AnkiIncompleteConfiguration;

export type PrepareLemmaResult = {
  __typename?: 'PrepareLemmaResult';
  ankiError?: Maybe<AnkiError>;
  error?: Maybe<PrepareLemmaError>;
  request?: Maybe<AddNoteRequest>;
};

export type Query = {
  __typename?: 'Query';
  Anki: Anki;
  AnkiConfig: AnkiConfig;
  AnkiConfigState: AnkiConfigStateResult;
  Lemmas?: Maybe<LemmasResult>;
  PrepareLemma: PrepareLemmaResult;
  RenderFields: RenderedFields;
};


export type QueryLemmasArgs = {
  query: Scalars['String']['input'];
};


export type QueryPrepareLemmaArgs = {
  lemma?: InputMaybe<LemmaInput>;
};


export type QueryRenderFieldsArgs = {
  fields?: InputMaybe<Array<Scalars['String']['input']>>;
  template?: InputMaybe<Scalars['String']['input']>;
};

export type RenderedField = {
  __typename?: 'RenderedField';
  error?: Maybe<Scalars['String']['output']>;
  field: Scalars['String']['output'];
  result: Scalars['String']['output'];
};

export type RenderedFields = {
  __typename?: 'RenderedFields';
  fields: Array<RenderedField>;
  template: Scalars['String']['output'];
  templateError?: Maybe<Scalars['String']['output']>;
};

export type SetAnkiConfigConnectionInput = {
  addr: Scalars['String']['input'];
  apiKey: Scalars['String']['input'];
};

export type SetAnkiConfigConnectionResult = {
  __typename?: 'SetAnkiConfigConnectionResult';
  error?: Maybe<ValidationError>;
};

export type SetAnkiConfigDeckInput = {
  name: Scalars['String']['input'];
};

export type SetAnkiConfigDeckResult = {
  __typename?: 'SetAnkiConfigDeckResult';
  error?: Maybe<ValidationError>;
};

export type SetAnkiConfigMappingInput = {
  mapping: Array<AnkiConfigMappingElementInput>;
};

export type SetAnkiConfigMappingResult = {
  __typename?: 'SetAnkiConfigMappingResult';
  error?: Maybe<AnkiConfigMappingError>;
};

export type SetAnkiConfigNote = {
  name: Scalars['String']['input'];
};

export type SetAnkiConfigNoteResult = {
  __typename?: 'SetAnkiConfigNoteResult';
  error?: Maybe<ValidationError>;
};

export type ValidationError = Error & {
  __typename?: 'ValidationError';
  message: Scalars['String']['output'];
  paths: Array<Scalars['String']['output']>;
};

export type Word = {
  __typename?: 'Word';
  furigana: Array<Furigana>;
  hiragana: Scalars['String']['output'];
  pitchShapes: Array<PitchShape>;
  word: Scalars['String']['output'];
};

export type WordInput = {
  furigana: Array<FuriganaInput>;
  hiragana: Scalars['String']['input'];
  pitchShapes: Array<PitchShapeInput>;
  word: Scalars['String']['input'];
};

export type GetHealthStatusQueryVariables = Exact<{ [key: string]: never; }>;


export type GetHealthStatusQuery = { __typename?: 'Query', AnkiConfigState: { __typename?: 'AnkiConfigStateResult', ankiConfigState?: { __typename?: 'AnkiConfigState', version: number, deckExists: boolean, noteTypeExists: boolean, noteHasAllFields: boolean, orderDefined: boolean } | null, error?: { __typename?: 'AnkiCollectionUnavailable', version: number, message: string } | { __typename?: 'AnkiConnectionError', message: string } | { __typename?: 'AnkiForbiddenOrigin', message: string } | { __typename?: 'AnkiInvalidAPIKey', version: number, message: string } | { __typename?: 'AnkiUnknownError', message: string } | null } };

export type PrepareLemmaQueryVariables = Exact<{
  lemma?: InputMaybe<LemmaInput>;
}>;


export type PrepareLemmaQuery = { __typename?: 'Query', PrepareLemma: { __typename?: 'PrepareLemmaResult', request?: { __typename?: 'AddNoteRequest', tags: Array<string>, audioURL: string, fields: Array<{ __typename?: 'AddNoteField', name: string, value: string }> } | null, error?: { __typename?: 'AnkiIncompleteConfiguration', message: string } | null, ankiError?: { __typename: 'AnkiCollectionUnavailable', message: string } | { __typename: 'AnkiConnectionError', message: string } | { __typename: 'AnkiForbiddenOrigin', message: string } | { __typename: 'AnkiInvalidAPIKey', message: string } | { __typename: 'AnkiUnknownError', message: string } | null } };

export type AddAnkiNoteMutationVariables = Exact<{
  note: AddNoteRequestInput;
}>;


export type AddAnkiNoteMutation = { __typename?: 'Mutation', addAnkiNote: { __typename?: 'AnkiAddNoteResult', error?: { __typename?: 'AnkiAddNoteDuplicateFound', message: string } | { __typename?: 'AnkiIncompleteConfiguration', message: string } | null, ankiError?: { __typename?: 'AnkiCollectionUnavailable', message: string } | { __typename?: 'AnkiConnectionError', message: string } | { __typename?: 'AnkiForbiddenOrigin', message: string } | { __typename?: 'AnkiInvalidAPIKey', message: string } | { __typename?: 'AnkiUnknownError', message: string } | null } };

export type GetConnectionConfigQueryVariables = Exact<{ [key: string]: never; }>;


export type GetConnectionConfigQuery = { __typename?: 'Query', AnkiConfig: { __typename?: 'AnkiConfig', addr: string, apiKey: string } };

export type UpdateConnectionConfigMutationVariables = Exact<{
  addr: Scalars['String']['input'];
  apiKey: Scalars['String']['input'];
}>;


export type UpdateConnectionConfigMutation = { __typename?: 'Mutation', setAnkiConfigConnection: { __typename?: 'SetAnkiConfigConnectionResult', error?: { __typename?: 'ValidationError', paths: Array<string>, message: string } | null } };

export type GetAnkiConfigCurrentNoteQueryVariables = Exact<{ [key: string]: never; }>;


export type GetAnkiConfigCurrentNoteQuery = { __typename?: 'Query', AnkiConfig: { __typename?: 'AnkiConfig', noteType: string } };

export type GetAnkiConfigCurrentDeckQueryVariables = Exact<{ [key: string]: never; }>;


export type GetAnkiConfigCurrentDeckQuery = { __typename?: 'Query', AnkiConfig: { __typename?: 'AnkiConfig', deck: string } };

export type GetAnkiDecksQueryVariables = Exact<{ [key: string]: never; }>;


export type GetAnkiDecksQuery = { __typename?: 'Query', Anki: { __typename?: 'Anki', decks: { __typename?: 'AnkiDecksResult', decks?: Array<string> | null, error?: { __typename: 'AnkiCollectionUnavailable' } | { __typename: 'AnkiConnectionError' } | { __typename: 'AnkiForbiddenOrigin' } | { __typename: 'AnkiInvalidAPIKey' } | { __typename: 'AnkiUnknownError' } | null } } };

export type SetAnkiConfigCurrentDeckMutationVariables = Exact<{
  name: Scalars['String']['input'];
}>;


export type SetAnkiConfigCurrentDeckMutation = { __typename?: 'Mutation', setAnkiConfigDeck: { __typename?: 'SetAnkiConfigDeckResult', error?: { __typename?: 'ValidationError', message: string } | null } };

export type CreateAnkiDeckMutationVariables = Exact<{
  name: Scalars['String']['input'];
}>;


export type CreateAnkiDeckMutation = { __typename?: 'Mutation', createAnkiDeck: { __typename?: 'CreateAnkiDeckResult', ankiError?: { __typename?: 'AnkiCollectionUnavailable', message: string } | { __typename?: 'AnkiConnectionError', message: string } | { __typename?: 'AnkiForbiddenOrigin', message: string } | { __typename?: 'AnkiInvalidAPIKey', message: string } | { __typename?: 'AnkiUnknownError', message: string } | null, error?: { __typename?: 'CreateAnkiDeckAlreadyExists', message: string } | { __typename?: 'ValidationError', message: string } | null } };

export type GetAnkiNoteFieldsAndMappingQueryVariables = Exact<{
  noteName: Scalars['String']['input'];
}>;


export type GetAnkiNoteFieldsAndMappingQuery = { __typename?: 'Query', AnkiConfig: { __typename?: 'AnkiConfig', mapping: Array<{ __typename?: 'AnkiMappingElement', key: string, value: string }> }, Anki: { __typename?: 'Anki', noteFields: { __typename?: 'AnkiNoteFieldsResult', noteFields?: Array<string> | null, error?: { __typename?: 'AnkiCollectionUnavailable', message: string } | { __typename?: 'AnkiConnectionError', message: string } | { __typename?: 'AnkiForbiddenOrigin', message: string } | { __typename?: 'AnkiInvalidAPIKey', message: string } | { __typename?: 'AnkiUnknownError', message: string } | null } } };

export type RenderFieldsQueryVariables = Exact<{
  fields: Array<Scalars['String']['input']> | Scalars['String']['input'];
}>;


export type RenderFieldsQuery = { __typename?: 'Query', RenderFields: { __typename?: 'RenderedFields', fields: Array<{ __typename?: 'RenderedField', result: string, error?: string | null }> } };

export type UpdateMappingMutationVariables = Exact<{
  fields: Array<AnkiConfigMappingElementInput> | AnkiConfigMappingElementInput;
}>;


export type UpdateMappingMutation = { __typename?: 'Mutation', setAnkiConfigMapping: { __typename?: 'SetAnkiConfigMappingResult', error?: { __typename?: 'AnkiConfigMappingError', message: string, fieldErrors?: Array<{ __typename?: 'AnkiConfigMappingElementError', key: string }> | null, valueErrors?: Array<{ __typename?: 'AnkiConfigMappingElementError', key: string }> | null } | null } };

export type GetAnkiNotesQueryVariables = Exact<{ [key: string]: never; }>;


export type GetAnkiNotesQuery = { __typename?: 'Query', Anki: { __typename?: 'Anki', notes: { __typename?: 'AnkiNotesResult', notes?: Array<string> | null, error?: { __typename: 'AnkiCollectionUnavailable' } | { __typename: 'AnkiConnectionError' } | { __typename: 'AnkiForbiddenOrigin' } | { __typename: 'AnkiInvalidAPIKey' } | { __typename: 'AnkiUnknownError' } | null } } };

export type SetAnkiConfigCurrentNoteMutationVariables = Exact<{
  name: Scalars['String']['input'];
}>;


export type SetAnkiConfigCurrentNoteMutation = { __typename?: 'Mutation', setAnkiConfigNote: { __typename?: 'SetAnkiConfigNoteResult', error?: { __typename?: 'ValidationError', message: string } | null } };

export type CreateDefaultAnkiNoteMutationVariables = Exact<{
  name: Scalars['String']['input'];
}>;


export type CreateDefaultAnkiNoteMutation = { __typename?: 'Mutation', createDefaultAnkiNote: { __typename?: 'CreateDefaultAnkiNoteResult', ankiError?: { __typename?: 'AnkiCollectionUnavailable', message: string } | { __typename?: 'AnkiConnectionError', message: string } | { __typename?: 'AnkiForbiddenOrigin', message: string } | { __typename?: 'AnkiInvalidAPIKey', message: string } | { __typename?: 'AnkiUnknownError', message: string } | null, error?: { __typename?: 'CreateDefaultAnkiNoteAlreadyExists', message: string } | { __typename?: 'ValidationError', message: string } | null } };

export type GetLemmasQueryVariables = Exact<{
  query: Scalars['String']['input'];
}>;


export type GetLemmasQuery = { __typename?: 'Query', Lemmas?: { __typename?: 'LemmasResult', lemmas: Array<{ __typename?: 'LemmaNoteInfo', noteID: string, lemma: { __typename?: 'Lemma', tags: Array<string>, definitions: Array<string>, partsOfSpeech: Array<string>, senseTags: Array<string>, slug: { __typename?: 'Word', word: string, hiragana: string, furigana: Array<{ __typename?: 'Furigana', kanji: string, hiragana: string }>, pitchShapes: Array<{ __typename?: 'PitchShape', hiragana: string, directions: Array<AccentDirection> }> }, forms: Array<{ __typename?: 'Word', word: string, hiragana: string, furigana: Array<{ __typename?: 'Furigana', kanji: string, hiragana: string }>, pitchShapes: Array<{ __typename?: 'PitchShape', hiragana: string, directions: Array<AccentDirection> }> }>, audio: Array<{ __typename?: 'Audio', type: string, source: string }> } }> } | null };


export const GetHealthStatusDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetHealthStatus"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"AnkiConfigState"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"ankiConfigState"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"version"}},{"kind":"Field","name":{"kind":"Name","value":"deckExists"}},{"kind":"Field","name":{"kind":"Name","value":"noteTypeExists"}},{"kind":"Field","name":{"kind":"Name","value":"noteHasAllFields"}},{"kind":"Field","name":{"kind":"Name","value":"orderDefined"}}]}},{"kind":"Field","name":{"kind":"Name","value":"error"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"AnkiConnectionError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"AnkiInvalidAPIKey"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"version"}},{"kind":"Field","name":{"kind":"Name","value":"message"}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"AnkiCollectionUnavailable"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"version"}},{"kind":"Field","name":{"kind":"Name","value":"message"}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"AnkiForbiddenOrigin"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Error"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}}]}}]}}]}}]}}]} as unknown as DocumentNode<GetHealthStatusQuery, GetHealthStatusQueryVariables>;
export const PrepareLemmaDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"PrepareLemma"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"lemma"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"LemmaInput"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"PrepareLemma"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"lemma"},"value":{"kind":"Variable","name":{"kind":"Name","value":"lemma"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"request"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"fields"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"value"}}]}},{"kind":"Field","name":{"kind":"Name","value":"tags"}},{"kind":"Field","name":{"kind":"Name","value":"audioURL"}}]}},{"kind":"Field","name":{"kind":"Name","value":"error"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"AnkiIncompleteConfiguration"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Error"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"ankiError"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Error"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}}]}}]}}]}}]}}]} as unknown as DocumentNode<PrepareLemmaQuery, PrepareLemmaQueryVariables>;
export const AddAnkiNoteDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"AddAnkiNote"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"note"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"AddNoteRequestInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"addAnkiNote"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"request"},"value":{"kind":"Variable","name":{"kind":"Name","value":"note"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"error"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"AnkiIncompleteConfiguration"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"AnkiAddNoteDuplicateFound"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"ankiError"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Error"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}}]}}]}}]}}]}}]} as unknown as DocumentNode<AddAnkiNoteMutation, AddAnkiNoteMutationVariables>;
export const GetConnectionConfigDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetConnectionConfig"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"AnkiConfig"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"addr"}},{"kind":"Field","name":{"kind":"Name","value":"apiKey"}}]}}]}}]} as unknown as DocumentNode<GetConnectionConfigQuery, GetConnectionConfigQueryVariables>;
export const UpdateConnectionConfigDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpdateConnectionConfig"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"addr"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"apiKey"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"setAnkiConfigConnection"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"ObjectValue","fields":[{"kind":"ObjectField","name":{"kind":"Name","value":"addr"},"value":{"kind":"Variable","name":{"kind":"Name","value":"addr"}}},{"kind":"ObjectField","name":{"kind":"Name","value":"apiKey"},"value":{"kind":"Variable","name":{"kind":"Name","value":"apiKey"}}}]}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"error"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ValidationError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"paths"}},{"kind":"Field","name":{"kind":"Name","value":"message"}}]}}]}}]}}]}}]} as unknown as DocumentNode<UpdateConnectionConfigMutation, UpdateConnectionConfigMutationVariables>;
export const GetAnkiConfigCurrentNoteDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetAnkiConfigCurrentNote"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"AnkiConfig"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"noteType"}}]}}]}}]} as unknown as DocumentNode<GetAnkiConfigCurrentNoteQuery, GetAnkiConfigCurrentNoteQueryVariables>;
export const GetAnkiConfigCurrentDeckDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetAnkiConfigCurrentDeck"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"AnkiConfig"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"deck"}}]}}]}}]} as unknown as DocumentNode<GetAnkiConfigCurrentDeckQuery, GetAnkiConfigCurrentDeckQueryVariables>;
export const GetAnkiDecksDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetAnkiDecks"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"Anki"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"decks"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"decks"}},{"kind":"Field","name":{"kind":"Name","value":"error"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}}]}}]}}]}}]}}]} as unknown as DocumentNode<GetAnkiDecksQuery, GetAnkiDecksQueryVariables>;
export const SetAnkiConfigCurrentDeckDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"SetAnkiConfigCurrentDeck"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"name"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"setAnkiConfigDeck"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"ObjectValue","fields":[{"kind":"ObjectField","name":{"kind":"Name","value":"name"},"value":{"kind":"Variable","name":{"kind":"Name","value":"name"}}}]}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"error"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}}]}}]}}]}}]} as unknown as DocumentNode<SetAnkiConfigCurrentDeckMutation, SetAnkiConfigCurrentDeckMutationVariables>;
export const CreateAnkiDeckDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"CreateAnkiDeck"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"name"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createAnkiDeck"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"ObjectValue","fields":[{"kind":"ObjectField","name":{"kind":"Name","value":"name"},"value":{"kind":"Variable","name":{"kind":"Name","value":"name"}}}]}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"ankiError"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Error"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"error"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"CreateAnkiDeckAlreadyExists"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ValidationError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Error"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}}]}}]}}]}}]}}]} as unknown as DocumentNode<CreateAnkiDeckMutation, CreateAnkiDeckMutationVariables>;
export const GetAnkiNoteFieldsAndMappingDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetAnkiNoteFieldsAndMapping"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"noteName"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"AnkiConfig"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"mapping"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"key"}},{"kind":"Field","name":{"kind":"Name","value":"value"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"Anki"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"noteFields"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"name"},"value":{"kind":"Variable","name":{"kind":"Name","value":"noteName"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"noteFields"}},{"kind":"Field","name":{"kind":"Name","value":"error"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Error"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}}]}}]}}]}}]}}]}}]} as unknown as DocumentNode<GetAnkiNoteFieldsAndMappingQuery, GetAnkiNoteFieldsAndMappingQueryVariables>;
export const RenderFieldsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"RenderFields"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"fields"}},"type":{"kind":"NonNullType","type":{"kind":"ListType","type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"RenderFields"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"fields"},"value":{"kind":"Variable","name":{"kind":"Name","value":"fields"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"fields"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"result"}},{"kind":"Field","name":{"kind":"Name","value":"error"}}]}}]}}]}}]} as unknown as DocumentNode<RenderFieldsQuery, RenderFieldsQueryVariables>;
export const UpdateMappingDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpdateMapping"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"fields"}},"type":{"kind":"NonNullType","type":{"kind":"ListType","type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"AnkiConfigMappingElementInput"}}}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"setAnkiConfigMapping"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"ObjectValue","fields":[{"kind":"ObjectField","name":{"kind":"Name","value":"mapping"},"value":{"kind":"Variable","name":{"kind":"Name","value":"fields"}}}]}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"error"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"fieldErrors"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"key"}}]}},{"kind":"Field","name":{"kind":"Name","value":"valueErrors"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"key"}}]}},{"kind":"Field","name":{"kind":"Name","value":"message"}}]}}]}}]}}]} as unknown as DocumentNode<UpdateMappingMutation, UpdateMappingMutationVariables>;
export const GetAnkiNotesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetAnkiNotes"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"Anki"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"notes"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"notes"}},{"kind":"Field","name":{"kind":"Name","value":"error"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}}]}}]}}]}}]}}]} as unknown as DocumentNode<GetAnkiNotesQuery, GetAnkiNotesQueryVariables>;
export const SetAnkiConfigCurrentNoteDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"SetAnkiConfigCurrentNote"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"name"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"setAnkiConfigNote"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"ObjectValue","fields":[{"kind":"ObjectField","name":{"kind":"Name","value":"name"},"value":{"kind":"Variable","name":{"kind":"Name","value":"name"}}}]}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"error"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}}]}}]}}]}}]} as unknown as DocumentNode<SetAnkiConfigCurrentNoteMutation, SetAnkiConfigCurrentNoteMutationVariables>;
export const CreateDefaultAnkiNoteDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"CreateDefaultAnkiNote"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"name"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createDefaultAnkiNote"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"ObjectValue","fields":[{"kind":"ObjectField","name":{"kind":"Name","value":"name"},"value":{"kind":"Variable","name":{"kind":"Name","value":"name"}}}]}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"ankiError"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Error"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"error"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"CreateDefaultAnkiNoteAlreadyExists"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ValidationError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Error"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}}]}}]}}]}}]}}]} as unknown as DocumentNode<CreateDefaultAnkiNoteMutation, CreateDefaultAnkiNoteMutationVariables>;
export const GetLemmasDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetLemmas"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"query"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"Lemmas"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"query"},"value":{"kind":"Variable","name":{"kind":"Name","value":"query"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"lemmas"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"noteID"}},{"kind":"Field","name":{"kind":"Name","value":"lemma"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"slug"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"word"}},{"kind":"Field","name":{"kind":"Name","value":"hiragana"}},{"kind":"Field","name":{"kind":"Name","value":"furigana"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"kanji"}},{"kind":"Field","name":{"kind":"Name","value":"hiragana"}}]}},{"kind":"Field","name":{"kind":"Name","value":"pitchShapes"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"hiragana"}},{"kind":"Field","name":{"kind":"Name","value":"directions"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"tags"}},{"kind":"Field","name":{"kind":"Name","value":"forms"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"word"}},{"kind":"Field","name":{"kind":"Name","value":"hiragana"}},{"kind":"Field","name":{"kind":"Name","value":"furigana"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"kanji"}},{"kind":"Field","name":{"kind":"Name","value":"hiragana"}}]}},{"kind":"Field","name":{"kind":"Name","value":"pitchShapes"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"hiragana"}},{"kind":"Field","name":{"kind":"Name","value":"directions"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"definitions"}},{"kind":"Field","name":{"kind":"Name","value":"partsOfSpeech"}},{"kind":"Field","name":{"kind":"Name","value":"senseTags"}},{"kind":"Field","name":{"kind":"Name","value":"audio"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"type"}},{"kind":"Field","name":{"kind":"Name","value":"source"}}]}}]}}]}}]}}]}}]} as unknown as DocumentNode<GetLemmasQuery, GetLemmasQueryVariables>;