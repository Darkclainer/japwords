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

export type AnkiConfig = {
  __typename?: 'AnkiConfig';
  addr: Scalars['String']['output'];
  apiKey: Scalars['String']['output'];
  deck: Scalars['String']['output'];
  mapping: Array<AnkiMappingElement>;
  noteType: Scalars['String']['output'];
};

export type AnkiConnectionError = Error & {
  __typename?: 'AnkiConnectionError';
  message: Scalars['String']['output'];
};

export type AnkiConnectionInput = {
  addr: Scalars['String']['input'];
  apiKey: Scalars['String']['input'];
};

export type AnkiConnectionPayload = ValidationError;

export type AnkiDeckInput = {
  name: Scalars['String']['input'];
};

export type AnkiDeckPayload = ValidationError;

export type AnkiMappingElement = {
  __typename?: 'AnkiMappingElement';
  key: Scalars['String']['output'];
  value: Scalars['String']['output'];
};

export type AnkiMappingElementError = {
  __typename?: 'AnkiMappingElementError';
  key: Scalars['String']['output'];
  message: Scalars['String']['output'];
};

export type AnkiMappingElementInput = {
  key: Scalars['String']['input'];
  value: Scalars['String']['input'];
};

export type AnkiMappingError = Error & {
  __typename?: 'AnkiMappingError';
  fieldErrors?: Maybe<Array<AnkiMappingElementError>>;
  message: Scalars['String']['output'];
  valueErrors?: Maybe<Array<AnkiMappingElementError>>;
};

export type AnkiMappingInput = {
  mapping: Array<AnkiMappingElementInput>;
};

export type AnkiMappingPayload = AnkiMappingError;

export type AnkiNoteTypeInput = {
  name: Scalars['String']['input'];
};

export type AnkiNoteTypePayload = ValidationError;

export type AnkiState = {
  __typename?: 'AnkiState';
  apiKeyRequired: Scalars['Boolean']['output'];
  connected: Scalars['Boolean']['output'];
  deckExists: Scalars['Boolean']['output'];
  noteMissingFields: Array<Scalars['String']['output']>;
  noteTypeExists: Scalars['Boolean']['output'];
  permissionGranted: Scalars['Boolean']['output'];
  version: Scalars['Int']['output'];
};

export type AnkiStatePayload = AnkiConnectionError | AnkiState;

export type Audio = {
  __typename?: 'Audio';
  source: Scalars['String']['output'];
  type: Scalars['String']['output'];
};

export type Error = {
  message: Scalars['String']['output'];
};

export type Furigana = {
  __typename?: 'Furigana';
  hiragana: Scalars['String']['output'];
  kanji: Scalars['String']['output'];
};

export type Lemma = {
  __typename?: 'Lemma';
  audio: Array<Audio>;
  forms: Array<Word>;
  senses: Array<Sense>;
  slug: Word;
  tags: Array<Scalars['String']['output']>;
};

export type Lemmas = {
  __typename?: 'Lemmas';
  lemmas: Array<Lemma>;
};

export type Mutation = {
  __typename?: 'Mutation';
  setAnkiConnection?: Maybe<AnkiConnectionPayload>;
  setAnkiDeck?: Maybe<AnkiDeckPayload>;
  setAnkiMapping?: Maybe<AnkiMappingPayload>;
  setAnkiNoteType?: Maybe<AnkiNoteTypePayload>;
};


export type MutationSetAnkiConnectionArgs = {
  input: AnkiConnectionInput;
};


export type MutationSetAnkiDeckArgs = {
  input: AnkiDeckInput;
};


export type MutationSetAnkiMappingArgs = {
  input: AnkiMappingInput;
};


export type MutationSetAnkiNoteTypeArgs = {
  input: AnkiNoteTypeInput;
};

export type Pitch = {
  __typename?: 'Pitch';
  hiragana: Scalars['String']['output'];
  pitch: Array<PitchType>;
};

export enum PitchType {
  Down = 'DOWN',
  Left = 'LEFT',
  Right = 'RIGHT',
  Up = 'UP'
}

export type Query = {
  __typename?: 'Query';
  AnkiConfig: AnkiConfig;
  AnkiState: AnkiStatePayload;
  Lemmas?: Maybe<Lemmas>;
};


export type QueryLemmasArgs = {
  query: Scalars['String']['input'];
};

export type Sense = {
  __typename?: 'Sense';
  definition: Array<Scalars['String']['output']>;
  partOfSpeech: Array<Scalars['String']['output']>;
  tags: Array<Scalars['String']['output']>;
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
  pitch: Array<Pitch>;
  word: Scalars['String']['output'];
};

export type GetLemmasQueryVariables = Exact<{
  query: Scalars['String']['input'];
}>;


export type GetLemmasQuery = { __typename?: 'Query', Lemmas?: { __typename?: 'Lemmas', lemmas: Array<{ __typename?: 'Lemma', tags: Array<string>, slug: { __typename?: 'Word', word: string, hiragana: string, furigana: Array<{ __typename?: 'Furigana', kanji: string, hiragana: string }>, pitch: Array<{ __typename?: 'Pitch', hiragana: string, pitch: Array<PitchType> }> }, forms: Array<{ __typename?: 'Word', word: string, hiragana: string, furigana: Array<{ __typename?: 'Furigana', kanji: string, hiragana: string }>, pitch: Array<{ __typename?: 'Pitch', hiragana: string, pitch: Array<PitchType> }> }>, senses: Array<{ __typename?: 'Sense', definition: Array<string>, partOfSpeech: Array<string>, tags: Array<string> }>, audio: Array<{ __typename?: 'Audio', type: string, source: string }> }> } | null };


export const GetLemmasDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetLemmas"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"query"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"Lemmas"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"query"},"value":{"kind":"Variable","name":{"kind":"Name","value":"query"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"lemmas"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"slug"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"word"}},{"kind":"Field","name":{"kind":"Name","value":"hiragana"}},{"kind":"Field","name":{"kind":"Name","value":"furigana"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"kanji"}},{"kind":"Field","name":{"kind":"Name","value":"hiragana"}}]}},{"kind":"Field","name":{"kind":"Name","value":"pitch"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"hiragana"}},{"kind":"Field","name":{"kind":"Name","value":"pitch"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"tags"}},{"kind":"Field","name":{"kind":"Name","value":"forms"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"word"}},{"kind":"Field","name":{"kind":"Name","value":"hiragana"}},{"kind":"Field","name":{"kind":"Name","value":"furigana"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"kanji"}},{"kind":"Field","name":{"kind":"Name","value":"hiragana"}}]}},{"kind":"Field","name":{"kind":"Name","value":"pitch"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"hiragana"}},{"kind":"Field","name":{"kind":"Name","value":"pitch"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"senses"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"definition"}},{"kind":"Field","name":{"kind":"Name","value":"partOfSpeech"}},{"kind":"Field","name":{"kind":"Name","value":"tags"}}]}},{"kind":"Field","name":{"kind":"Name","value":"audio"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"type"}},{"kind":"Field","name":{"kind":"Name","value":"source"}}]}}]}}]}}]}}]} as unknown as DocumentNode<GetLemmasQuery, GetLemmasQueryVariables>;