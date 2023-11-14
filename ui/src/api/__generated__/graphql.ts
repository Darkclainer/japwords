/* eslint-disable */
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = {
  [_ in K]?: never;
};
export type Incremental<T> =
  | T
  | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string };
  String: { input: string; output: string };
  Boolean: { input: boolean; output: boolean };
  Int: { input: number; output: number };
  Float: { input: number; output: number };
};

export type Anki = {
  __typename?: 'Anki';
  decks: Array<Scalars['String']['output']>;
  noteFields: Array<Scalars['String']['output']>;
  notes: Array<Scalars['String']['output']>;
};

export type AnkiNoteFieldsArgs = {
  name: Scalars['String']['input'];
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

export type AnkiError = AnkiConnectionError | AnkiPermissionError | AnkiUnauthorizedError;

export type AnkiMappingElement = {
  __typename?: 'AnkiMappingElement';
  key: Scalars['String']['output'];
  value: Scalars['String']['output'];
};

export type AnkiPermissionError = Error & {
  __typename?: 'AnkiPermissionError';
  message: Scalars['String']['output'];
  version: Scalars['Int']['output'];
};

export type AnkiResult = {
  __typename?: 'AnkiResult';
  anki?: Maybe<Anki>;
  error?: Maybe<AnkiError>;
};

export type AnkiUnauthorizedError = Error & {
  __typename?: 'AnkiUnauthorizedError';
  message: Scalars['String']['output'];
  version: Scalars['Int']['output'];
};

export type Audio = {
  __typename?: 'Audio';
  source: Scalars['String']['output'];
  type: Scalars['String']['output'];
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
  createAnkiDeck: CreateAnkiDeckResult;
  createDefaultAnkiNote: CreateDefaultAnkiNoteResult;
  setAnkiConfigConnection: SetAnkiConfigConnectionResult;
  setAnkiConfigDeck: SetAnkiConfigDeckResult;
  setAnkiConfigMapping: SetAnkiConfigMappingResult;
  setAnkiConfigNote: SetAnkiConfigNoteResult;
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

export type Pitch = {
  __typename?: 'Pitch';
  hiragana: Scalars['String']['output'];
  pitch: Array<PitchType>;
};

export enum PitchType {
  Down = 'DOWN',
  Left = 'LEFT',
  Right = 'RIGHT',
  Up = 'UP',
}

export type Query = {
  __typename?: 'Query';
  Anki: AnkiResult;
  AnkiConfig: AnkiConfig;
  AnkiConfigState: AnkiConfigStateResult;
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
  pitch: Array<Pitch>;
  word: Scalars['String']['output'];
};

export type GetHealthStatusQueryVariables = Exact<{ [key: string]: never }>;

export type GetHealthStatusQuery = {
  __typename?: 'Query';
  AnkiConfigState: {
    __typename?: 'AnkiConfigStateResult';
    ankiConfigState?: {
      __typename?: 'AnkiConfigState';
      version: number;
      deckExists: boolean;
      noteTypeExists: boolean;
      noteHasAllFields: boolean;
    } | null;
    error?:
      | { __typename?: 'AnkiConnectionError'; message: string }
      | { __typename?: 'AnkiPermissionError'; version: number; message: string }
      | { __typename?: 'AnkiUnauthorizedError'; version: number; message: string }
      | null;
  };
};

export type GetConnectionConfigQueryVariables = Exact<{ [key: string]: never }>;

export type GetConnectionConfigQuery = {
  __typename?: 'Query';
  AnkiConfig: { __typename?: 'AnkiConfig'; addr: string; apiKey: string };
};

export type UpdateConnectionConfigMutationVariables = Exact<{
  addr: Scalars['String']['input'];
  apiKey: Scalars['String']['input'];
}>;

export type UpdateConnectionConfigMutation = {
  __typename?: 'Mutation';
  setAnkiConfigConnection: {
    __typename?: 'SetAnkiConfigConnectionResult';
    error?: { __typename?: 'ValidationError'; paths: Array<string>; message: string } | null;
  };
};

export type GetAnkiConfigCurrentDeckQueryVariables = Exact<{ [key: string]: never }>;

export type GetAnkiConfigCurrentDeckQuery = {
  __typename?: 'Query';
  AnkiConfig: { __typename?: 'AnkiConfig'; deck: string };
};

export type GetAnkiDecksQueryVariables = Exact<{ [key: string]: never }>;

export type GetAnkiDecksQuery = {
  __typename?: 'Query';
  Anki: {
    __typename?: 'AnkiResult';
    anki?: { __typename?: 'Anki'; decks: Array<string> } | null;
    error?:
      | { __typename: 'AnkiConnectionError' }
      | { __typename: 'AnkiPermissionError' }
      | { __typename: 'AnkiUnauthorizedError' }
      | null;
  };
};

export type SetAnkiConfigCurrentDeckMutationVariables = Exact<{
  name: Scalars['String']['input'];
}>;

export type SetAnkiConfigCurrentDeckMutation = {
  __typename?: 'Mutation';
  setAnkiConfigDeck: {
    __typename?: 'SetAnkiConfigDeckResult';
    error?: { __typename?: 'ValidationError'; message: string } | null;
  };
};

export type CreateAnkiDeckMutationVariables = Exact<{
  name: Scalars['String']['input'];
}>;

export type CreateAnkiDeckMutation = {
  __typename?: 'Mutation';
  createAnkiDeck: {
    __typename?: 'CreateAnkiDeckResult';
    ankiError?:
      | { __typename?: 'AnkiConnectionError'; message: string }
      | { __typename?: 'AnkiPermissionError'; message: string }
      | { __typename?: 'AnkiUnauthorizedError'; message: string }
      | null;
    error?:
      | { __typename?: 'CreateAnkiDeckAlreadyExists'; message: string }
      | { __typename?: 'ValidationError'; message: string }
      | null;
  };
};

export type GetAnkiConfigCurrentNoteQueryVariables = Exact<{ [key: string]: never }>;

export type GetAnkiConfigCurrentNoteQuery = {
  __typename?: 'Query';
  AnkiConfig: { __typename?: 'AnkiConfig'; noteType: string };
};

export type GetAnkiNoteFieldsAndMappingQueryVariables = Exact<{
  noteName: Scalars['String']['input'];
}>;

export type GetAnkiNoteFieldsAndMappingQuery = {
  __typename?: 'Query';
  AnkiConfig: {
    __typename?: 'AnkiConfig';
    mapping: Array<{ __typename?: 'AnkiMappingElement'; key: string; value: string }>;
  };
  Anki: {
    __typename?: 'AnkiResult';
    anki?: { __typename?: 'Anki'; noteFields: Array<string> } | null;
    error?:
      | { __typename?: 'AnkiConnectionError'; message: string }
      | { __typename?: 'AnkiPermissionError'; message: string }
      | { __typename?: 'AnkiUnauthorizedError'; message: string }
      | null;
  };
};

export type GetAnkiNotesQueryVariables = Exact<{ [key: string]: never }>;

export type GetAnkiNotesQuery = {
  __typename?: 'Query';
  Anki: {
    __typename?: 'AnkiResult';
    anki?: { __typename?: 'Anki'; notes: Array<string> } | null;
    error?:
      | { __typename: 'AnkiConnectionError' }
      | { __typename: 'AnkiPermissionError' }
      | { __typename: 'AnkiUnauthorizedError' }
      | null;
  };
};

export type SetAnkiConfigCurrentNoteMutationVariables = Exact<{
  name: Scalars['String']['input'];
}>;

export type SetAnkiConfigCurrentNoteMutation = {
  __typename?: 'Mutation';
  setAnkiConfigNote: {
    __typename?: 'SetAnkiConfigNoteResult';
    error?: { __typename?: 'ValidationError'; message: string } | null;
  };
};

export type CreateDefaultAnkiNoteMutationVariables = Exact<{
  name: Scalars['String']['input'];
}>;

export type CreateDefaultAnkiNoteMutation = {
  __typename?: 'Mutation';
  createDefaultAnkiNote: {
    __typename?: 'CreateDefaultAnkiNoteResult';
    ankiError?:
      | { __typename?: 'AnkiConnectionError'; message: string }
      | { __typename?: 'AnkiPermissionError'; message: string }
      | { __typename?: 'AnkiUnauthorizedError'; message: string }
      | null;
    error?:
      | { __typename?: 'CreateDefaultAnkiNoteAlreadyExists'; message: string }
      | { __typename?: 'ValidationError'; message: string }
      | null;
  };
};

export type GetLemmasQueryVariables = Exact<{
  query: Scalars['String']['input'];
}>;

export type GetLemmasQuery = {
  __typename?: 'Query';
  Lemmas?: {
    __typename?: 'Lemmas';
    lemmas: Array<{
      __typename?: 'Lemma';
      tags: Array<string>;
      slug: {
        __typename?: 'Word';
        word: string;
        hiragana: string;
        furigana: Array<{ __typename?: 'Furigana'; kanji: string; hiragana: string }>;
        pitch: Array<{ __typename?: 'Pitch'; hiragana: string; pitch: Array<PitchType> }>;
      };
      forms: Array<{
        __typename?: 'Word';
        word: string;
        hiragana: string;
        furigana: Array<{ __typename?: 'Furigana'; kanji: string; hiragana: string }>;
        pitch: Array<{ __typename?: 'Pitch'; hiragana: string; pitch: Array<PitchType> }>;
      }>;
      senses: Array<{
        __typename?: 'Sense';
        definition: Array<string>;
        partOfSpeech: Array<string>;
        tags: Array<string>;
      }>;
      audio: Array<{ __typename?: 'Audio'; type: string; source: string }>;
    }>;
  } | null;
};

export const GetHealthStatusDocument = {
  kind: 'Document',
  definitions: [
    {
      kind: 'OperationDefinition',
      operation: 'query',
      name: { kind: 'Name', value: 'GetHealthStatus' },
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'AnkiConfigState' },
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'ankiConfigState' },
                  selectionSet: {
                    kind: 'SelectionSet',
                    selections: [
                      { kind: 'Field', name: { kind: 'Name', value: 'version' } },
                      { kind: 'Field', name: { kind: 'Name', value: 'deckExists' } },
                      { kind: 'Field', name: { kind: 'Name', value: 'noteTypeExists' } },
                      { kind: 'Field', name: { kind: 'Name', value: 'noteHasAllFields' } },
                    ],
                  },
                },
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'error' },
                  selectionSet: {
                    kind: 'SelectionSet',
                    selections: [
                      {
                        kind: 'InlineFragment',
                        typeCondition: {
                          kind: 'NamedType',
                          name: { kind: 'Name', value: 'AnkiConnectionError' },
                        },
                        selectionSet: {
                          kind: 'SelectionSet',
                          selections: [{ kind: 'Field', name: { kind: 'Name', value: 'message' } }],
                        },
                      },
                      {
                        kind: 'InlineFragment',
                        typeCondition: {
                          kind: 'NamedType',
                          name: { kind: 'Name', value: 'AnkiPermissionError' },
                        },
                        selectionSet: {
                          kind: 'SelectionSet',
                          selections: [
                            { kind: 'Field', name: { kind: 'Name', value: 'version' } },
                            { kind: 'Field', name: { kind: 'Name', value: 'message' } },
                          ],
                        },
                      },
                      {
                        kind: 'InlineFragment',
                        typeCondition: {
                          kind: 'NamedType',
                          name: { kind: 'Name', value: 'AnkiUnauthorizedError' },
                        },
                        selectionSet: {
                          kind: 'SelectionSet',
                          selections: [
                            { kind: 'Field', name: { kind: 'Name', value: 'version' } },
                            { kind: 'Field', name: { kind: 'Name', value: 'message' } },
                          ],
                        },
                      },
                      {
                        kind: 'InlineFragment',
                        typeCondition: {
                          kind: 'NamedType',
                          name: { kind: 'Name', value: 'Error' },
                        },
                        selectionSet: {
                          kind: 'SelectionSet',
                          selections: [{ kind: 'Field', name: { kind: 'Name', value: 'message' } }],
                        },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<GetHealthStatusQuery, GetHealthStatusQueryVariables>;
export const GetConnectionConfigDocument = {
  kind: 'Document',
  definitions: [
    {
      kind: 'OperationDefinition',
      operation: 'query',
      name: { kind: 'Name', value: 'GetConnectionConfig' },
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'AnkiConfig' },
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                { kind: 'Field', name: { kind: 'Name', value: 'addr' } },
                { kind: 'Field', name: { kind: 'Name', value: 'apiKey' } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<GetConnectionConfigQuery, GetConnectionConfigQueryVariables>;
export const UpdateConnectionConfigDocument = {
  kind: 'Document',
  definitions: [
    {
      kind: 'OperationDefinition',
      operation: 'mutation',
      name: { kind: 'Name', value: 'UpdateConnectionConfig' },
      variableDefinitions: [
        {
          kind: 'VariableDefinition',
          variable: { kind: 'Variable', name: { kind: 'Name', value: 'addr' } },
          type: {
            kind: 'NonNullType',
            type: { kind: 'NamedType', name: { kind: 'Name', value: 'String' } },
          },
        },
        {
          kind: 'VariableDefinition',
          variable: { kind: 'Variable', name: { kind: 'Name', value: 'apiKey' } },
          type: {
            kind: 'NonNullType',
            type: { kind: 'NamedType', name: { kind: 'Name', value: 'String' } },
          },
        },
      ],
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'setAnkiConfigConnection' },
            arguments: [
              {
                kind: 'Argument',
                name: { kind: 'Name', value: 'input' },
                value: {
                  kind: 'ObjectValue',
                  fields: [
                    {
                      kind: 'ObjectField',
                      name: { kind: 'Name', value: 'addr' },
                      value: { kind: 'Variable', name: { kind: 'Name', value: 'addr' } },
                    },
                    {
                      kind: 'ObjectField',
                      name: { kind: 'Name', value: 'apiKey' },
                      value: { kind: 'Variable', name: { kind: 'Name', value: 'apiKey' } },
                    },
                  ],
                },
              },
            ],
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'error' },
                  selectionSet: {
                    kind: 'SelectionSet',
                    selections: [
                      {
                        kind: 'InlineFragment',
                        typeCondition: {
                          kind: 'NamedType',
                          name: { kind: 'Name', value: 'ValidationError' },
                        },
                        selectionSet: {
                          kind: 'SelectionSet',
                          selections: [
                            { kind: 'Field', name: { kind: 'Name', value: 'paths' } },
                            { kind: 'Field', name: { kind: 'Name', value: 'message' } },
                          ],
                        },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  UpdateConnectionConfigMutation,
  UpdateConnectionConfigMutationVariables
>;
export const GetAnkiConfigCurrentDeckDocument = {
  kind: 'Document',
  definitions: [
    {
      kind: 'OperationDefinition',
      operation: 'query',
      name: { kind: 'Name', value: 'GetAnkiConfigCurrentDeck' },
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'AnkiConfig' },
            selectionSet: {
              kind: 'SelectionSet',
              selections: [{ kind: 'Field', name: { kind: 'Name', value: 'deck' } }],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<GetAnkiConfigCurrentDeckQuery, GetAnkiConfigCurrentDeckQueryVariables>;
export const GetAnkiDecksDocument = {
  kind: 'Document',
  definitions: [
    {
      kind: 'OperationDefinition',
      operation: 'query',
      name: { kind: 'Name', value: 'GetAnkiDecks' },
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'Anki' },
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'anki' },
                  selectionSet: {
                    kind: 'SelectionSet',
                    selections: [{ kind: 'Field', name: { kind: 'Name', value: 'decks' } }],
                  },
                },
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'error' },
                  selectionSet: {
                    kind: 'SelectionSet',
                    selections: [{ kind: 'Field', name: { kind: 'Name', value: '__typename' } }],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<GetAnkiDecksQuery, GetAnkiDecksQueryVariables>;
export const SetAnkiConfigCurrentDeckDocument = {
  kind: 'Document',
  definitions: [
    {
      kind: 'OperationDefinition',
      operation: 'mutation',
      name: { kind: 'Name', value: 'SetAnkiConfigCurrentDeck' },
      variableDefinitions: [
        {
          kind: 'VariableDefinition',
          variable: { kind: 'Variable', name: { kind: 'Name', value: 'name' } },
          type: {
            kind: 'NonNullType',
            type: { kind: 'NamedType', name: { kind: 'Name', value: 'String' } },
          },
        },
      ],
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'setAnkiConfigDeck' },
            arguments: [
              {
                kind: 'Argument',
                name: { kind: 'Name', value: 'input' },
                value: {
                  kind: 'ObjectValue',
                  fields: [
                    {
                      kind: 'ObjectField',
                      name: { kind: 'Name', value: 'name' },
                      value: { kind: 'Variable', name: { kind: 'Name', value: 'name' } },
                    },
                  ],
                },
              },
            ],
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'error' },
                  selectionSet: {
                    kind: 'SelectionSet',
                    selections: [{ kind: 'Field', name: { kind: 'Name', value: 'message' } }],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  SetAnkiConfigCurrentDeckMutation,
  SetAnkiConfigCurrentDeckMutationVariables
>;
export const CreateAnkiDeckDocument = {
  kind: 'Document',
  definitions: [
    {
      kind: 'OperationDefinition',
      operation: 'mutation',
      name: { kind: 'Name', value: 'CreateAnkiDeck' },
      variableDefinitions: [
        {
          kind: 'VariableDefinition',
          variable: { kind: 'Variable', name: { kind: 'Name', value: 'name' } },
          type: {
            kind: 'NonNullType',
            type: { kind: 'NamedType', name: { kind: 'Name', value: 'String' } },
          },
        },
      ],
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'createAnkiDeck' },
            arguments: [
              {
                kind: 'Argument',
                name: { kind: 'Name', value: 'input' },
                value: {
                  kind: 'ObjectValue',
                  fields: [
                    {
                      kind: 'ObjectField',
                      name: { kind: 'Name', value: 'name' },
                      value: { kind: 'Variable', name: { kind: 'Name', value: 'name' } },
                    },
                  ],
                },
              },
            ],
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'ankiError' },
                  selectionSet: {
                    kind: 'SelectionSet',
                    selections: [
                      {
                        kind: 'InlineFragment',
                        typeCondition: {
                          kind: 'NamedType',
                          name: { kind: 'Name', value: 'Error' },
                        },
                        selectionSet: {
                          kind: 'SelectionSet',
                          selections: [{ kind: 'Field', name: { kind: 'Name', value: 'message' } }],
                        },
                      },
                    ],
                  },
                },
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'error' },
                  selectionSet: {
                    kind: 'SelectionSet',
                    selections: [
                      {
                        kind: 'InlineFragment',
                        typeCondition: {
                          kind: 'NamedType',
                          name: { kind: 'Name', value: 'CreateAnkiDeckAlreadyExists' },
                        },
                        selectionSet: {
                          kind: 'SelectionSet',
                          selections: [{ kind: 'Field', name: { kind: 'Name', value: 'message' } }],
                        },
                      },
                      {
                        kind: 'InlineFragment',
                        typeCondition: {
                          kind: 'NamedType',
                          name: { kind: 'Name', value: 'ValidationError' },
                        },
                        selectionSet: {
                          kind: 'SelectionSet',
                          selections: [{ kind: 'Field', name: { kind: 'Name', value: 'message' } }],
                        },
                      },
                      {
                        kind: 'InlineFragment',
                        typeCondition: {
                          kind: 'NamedType',
                          name: { kind: 'Name', value: 'Error' },
                        },
                        selectionSet: {
                          kind: 'SelectionSet',
                          selections: [{ kind: 'Field', name: { kind: 'Name', value: 'message' } }],
                        },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<CreateAnkiDeckMutation, CreateAnkiDeckMutationVariables>;
export const GetAnkiConfigCurrentNoteDocument = {
  kind: 'Document',
  definitions: [
    {
      kind: 'OperationDefinition',
      operation: 'query',
      name: { kind: 'Name', value: 'GetAnkiConfigCurrentNote' },
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'AnkiConfig' },
            selectionSet: {
              kind: 'SelectionSet',
              selections: [{ kind: 'Field', name: { kind: 'Name', value: 'noteType' } }],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<GetAnkiConfigCurrentNoteQuery, GetAnkiConfigCurrentNoteQueryVariables>;
export const GetAnkiNoteFieldsAndMappingDocument = {
  kind: 'Document',
  definitions: [
    {
      kind: 'OperationDefinition',
      operation: 'query',
      name: { kind: 'Name', value: 'GetAnkiNoteFieldsAndMapping' },
      variableDefinitions: [
        {
          kind: 'VariableDefinition',
          variable: { kind: 'Variable', name: { kind: 'Name', value: 'noteName' } },
          type: {
            kind: 'NonNullType',
            type: { kind: 'NamedType', name: { kind: 'Name', value: 'String' } },
          },
        },
      ],
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'AnkiConfig' },
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'mapping' },
                  selectionSet: {
                    kind: 'SelectionSet',
                    selections: [
                      { kind: 'Field', name: { kind: 'Name', value: 'key' } },
                      { kind: 'Field', name: { kind: 'Name', value: 'value' } },
                    ],
                  },
                },
              ],
            },
          },
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'Anki' },
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'anki' },
                  selectionSet: {
                    kind: 'SelectionSet',
                    selections: [
                      {
                        kind: 'Field',
                        name: { kind: 'Name', value: 'noteFields' },
                        arguments: [
                          {
                            kind: 'Argument',
                            name: { kind: 'Name', value: 'name' },
                            value: { kind: 'Variable', name: { kind: 'Name', value: 'noteName' } },
                          },
                        ],
                      },
                    ],
                  },
                },
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'error' },
                  selectionSet: {
                    kind: 'SelectionSet',
                    selections: [
                      {
                        kind: 'InlineFragment',
                        typeCondition: {
                          kind: 'NamedType',
                          name: { kind: 'Name', value: 'Error' },
                        },
                        selectionSet: {
                          kind: 'SelectionSet',
                          selections: [{ kind: 'Field', name: { kind: 'Name', value: 'message' } }],
                        },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  GetAnkiNoteFieldsAndMappingQuery,
  GetAnkiNoteFieldsAndMappingQueryVariables
>;
export const GetAnkiNotesDocument = {
  kind: 'Document',
  definitions: [
    {
      kind: 'OperationDefinition',
      operation: 'query',
      name: { kind: 'Name', value: 'GetAnkiNotes' },
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'Anki' },
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'anki' },
                  selectionSet: {
                    kind: 'SelectionSet',
                    selections: [{ kind: 'Field', name: { kind: 'Name', value: 'notes' } }],
                  },
                },
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'error' },
                  selectionSet: {
                    kind: 'SelectionSet',
                    selections: [{ kind: 'Field', name: { kind: 'Name', value: '__typename' } }],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<GetAnkiNotesQuery, GetAnkiNotesQueryVariables>;
export const SetAnkiConfigCurrentNoteDocument = {
  kind: 'Document',
  definitions: [
    {
      kind: 'OperationDefinition',
      operation: 'mutation',
      name: { kind: 'Name', value: 'SetAnkiConfigCurrentNote' },
      variableDefinitions: [
        {
          kind: 'VariableDefinition',
          variable: { kind: 'Variable', name: { kind: 'Name', value: 'name' } },
          type: {
            kind: 'NonNullType',
            type: { kind: 'NamedType', name: { kind: 'Name', value: 'String' } },
          },
        },
      ],
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'setAnkiConfigNote' },
            arguments: [
              {
                kind: 'Argument',
                name: { kind: 'Name', value: 'input' },
                value: {
                  kind: 'ObjectValue',
                  fields: [
                    {
                      kind: 'ObjectField',
                      name: { kind: 'Name', value: 'name' },
                      value: { kind: 'Variable', name: { kind: 'Name', value: 'name' } },
                    },
                  ],
                },
              },
            ],
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'error' },
                  selectionSet: {
                    kind: 'SelectionSet',
                    selections: [{ kind: 'Field', name: { kind: 'Name', value: 'message' } }],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  SetAnkiConfigCurrentNoteMutation,
  SetAnkiConfigCurrentNoteMutationVariables
>;
export const CreateDefaultAnkiNoteDocument = {
  kind: 'Document',
  definitions: [
    {
      kind: 'OperationDefinition',
      operation: 'mutation',
      name: { kind: 'Name', value: 'CreateDefaultAnkiNote' },
      variableDefinitions: [
        {
          kind: 'VariableDefinition',
          variable: { kind: 'Variable', name: { kind: 'Name', value: 'name' } },
          type: {
            kind: 'NonNullType',
            type: { kind: 'NamedType', name: { kind: 'Name', value: 'String' } },
          },
        },
      ],
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'createDefaultAnkiNote' },
            arguments: [
              {
                kind: 'Argument',
                name: { kind: 'Name', value: 'input' },
                value: {
                  kind: 'ObjectValue',
                  fields: [
                    {
                      kind: 'ObjectField',
                      name: { kind: 'Name', value: 'name' },
                      value: { kind: 'Variable', name: { kind: 'Name', value: 'name' } },
                    },
                  ],
                },
              },
            ],
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'ankiError' },
                  selectionSet: {
                    kind: 'SelectionSet',
                    selections: [
                      {
                        kind: 'InlineFragment',
                        typeCondition: {
                          kind: 'NamedType',
                          name: { kind: 'Name', value: 'Error' },
                        },
                        selectionSet: {
                          kind: 'SelectionSet',
                          selections: [{ kind: 'Field', name: { kind: 'Name', value: 'message' } }],
                        },
                      },
                    ],
                  },
                },
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'error' },
                  selectionSet: {
                    kind: 'SelectionSet',
                    selections: [
                      {
                        kind: 'InlineFragment',
                        typeCondition: {
                          kind: 'NamedType',
                          name: { kind: 'Name', value: 'CreateDefaultAnkiNoteAlreadyExists' },
                        },
                        selectionSet: {
                          kind: 'SelectionSet',
                          selections: [{ kind: 'Field', name: { kind: 'Name', value: 'message' } }],
                        },
                      },
                      {
                        kind: 'InlineFragment',
                        typeCondition: {
                          kind: 'NamedType',
                          name: { kind: 'Name', value: 'ValidationError' },
                        },
                        selectionSet: {
                          kind: 'SelectionSet',
                          selections: [{ kind: 'Field', name: { kind: 'Name', value: 'message' } }],
                        },
                      },
                      {
                        kind: 'InlineFragment',
                        typeCondition: {
                          kind: 'NamedType',
                          name: { kind: 'Name', value: 'Error' },
                        },
                        selectionSet: {
                          kind: 'SelectionSet',
                          selections: [{ kind: 'Field', name: { kind: 'Name', value: 'message' } }],
                        },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<CreateDefaultAnkiNoteMutation, CreateDefaultAnkiNoteMutationVariables>;
export const GetLemmasDocument = {
  kind: 'Document',
  definitions: [
    {
      kind: 'OperationDefinition',
      operation: 'query',
      name: { kind: 'Name', value: 'GetLemmas' },
      variableDefinitions: [
        {
          kind: 'VariableDefinition',
          variable: { kind: 'Variable', name: { kind: 'Name', value: 'query' } },
          type: {
            kind: 'NonNullType',
            type: { kind: 'NamedType', name: { kind: 'Name', value: 'String' } },
          },
        },
      ],
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'Lemmas' },
            arguments: [
              {
                kind: 'Argument',
                name: { kind: 'Name', value: 'query' },
                value: { kind: 'Variable', name: { kind: 'Name', value: 'query' } },
              },
            ],
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'lemmas' },
                  selectionSet: {
                    kind: 'SelectionSet',
                    selections: [
                      {
                        kind: 'Field',
                        name: { kind: 'Name', value: 'slug' },
                        selectionSet: {
                          kind: 'SelectionSet',
                          selections: [
                            { kind: 'Field', name: { kind: 'Name', value: 'word' } },
                            { kind: 'Field', name: { kind: 'Name', value: 'hiragana' } },
                            {
                              kind: 'Field',
                              name: { kind: 'Name', value: 'furigana' },
                              selectionSet: {
                                kind: 'SelectionSet',
                                selections: [
                                  { kind: 'Field', name: { kind: 'Name', value: 'kanji' } },
                                  { kind: 'Field', name: { kind: 'Name', value: 'hiragana' } },
                                ],
                              },
                            },
                            {
                              kind: 'Field',
                              name: { kind: 'Name', value: 'pitch' },
                              selectionSet: {
                                kind: 'SelectionSet',
                                selections: [
                                  { kind: 'Field', name: { kind: 'Name', value: 'hiragana' } },
                                  { kind: 'Field', name: { kind: 'Name', value: 'pitch' } },
                                ],
                              },
                            },
                          ],
                        },
                      },
                      { kind: 'Field', name: { kind: 'Name', value: 'tags' } },
                      {
                        kind: 'Field',
                        name: { kind: 'Name', value: 'forms' },
                        selectionSet: {
                          kind: 'SelectionSet',
                          selections: [
                            { kind: 'Field', name: { kind: 'Name', value: 'word' } },
                            { kind: 'Field', name: { kind: 'Name', value: 'hiragana' } },
                            {
                              kind: 'Field',
                              name: { kind: 'Name', value: 'furigana' },
                              selectionSet: {
                                kind: 'SelectionSet',
                                selections: [
                                  { kind: 'Field', name: { kind: 'Name', value: 'kanji' } },
                                  { kind: 'Field', name: { kind: 'Name', value: 'hiragana' } },
                                ],
                              },
                            },
                            {
                              kind: 'Field',
                              name: { kind: 'Name', value: 'pitch' },
                              selectionSet: {
                                kind: 'SelectionSet',
                                selections: [
                                  { kind: 'Field', name: { kind: 'Name', value: 'hiragana' } },
                                  { kind: 'Field', name: { kind: 'Name', value: 'pitch' } },
                                ],
                              },
                            },
                          ],
                        },
                      },
                      {
                        kind: 'Field',
                        name: { kind: 'Name', value: 'senses' },
                        selectionSet: {
                          kind: 'SelectionSet',
                          selections: [
                            { kind: 'Field', name: { kind: 'Name', value: 'definition' } },
                            { kind: 'Field', name: { kind: 'Name', value: 'partOfSpeech' } },
                            { kind: 'Field', name: { kind: 'Name', value: 'tags' } },
                          ],
                        },
                      },
                      {
                        kind: 'Field',
                        name: { kind: 'Name', value: 'audio' },
                        selectionSet: {
                          kind: 'SelectionSet',
                          selections: [
                            { kind: 'Field', name: { kind: 'Name', value: 'type' } },
                            { kind: 'Field', name: { kind: 'Name', value: 'source' } },
                          ],
                        },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<GetLemmasQuery, GetLemmasQueryVariables>;
