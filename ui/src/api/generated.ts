import client from "../../src/apollo-client";
import type {
        ApolloQueryResult, ObservableQuery, WatchQueryOptions
      } from "@apollo/client";
import { readable } from "svelte/store";
import type { Readable } from "svelte/store";
import gql from "graphql-tag"
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
};

export type Audio = {
  __typename?: 'Audio';
  source: Scalars['String'];
  type: Scalars['String'];
};

export type Furigana = {
  __typename?: 'Furigana';
  hiragana: Scalars['String'];
  kanji: Scalars['String'];
};

export type Lemma = {
  __typename?: 'Lemma';
  audio: Array<Audio>;
  forms: Array<Word>;
  senses: Array<Sense>;
  slug: Word;
  tags: Array<Scalars['String']>;
};

export type Lemmas = {
  __typename?: 'Lemmas';
  lemmas: Array<Lemma>;
};

export type Pitch = {
  __typename?: 'Pitch';
  hiragana: Scalars['String'];
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
  Lemmas?: Maybe<Lemmas>;
};


export type QueryLemmasArgs = {
  query: Scalars['String'];
};

export type Sense = {
  __typename?: 'Sense';
  definition: Array<Scalars['String']>;
  partOfSpeech: Array<Scalars['String']>;
  tags: Array<Scalars['String']>;
};

export type Word = {
  __typename?: 'Word';
  furigana: Array<Furigana>;
  hiragana: Scalars['String'];
  pitch: Array<Pitch>;
  word: Scalars['String'];
};

export type GetLemmasQueryVariables = Exact<{
  query: Scalars['String'];
}>;


export type GetLemmasQuery = { __typename?: 'Query', Lemmas?: { __typename?: 'Lemmas', lemmas: Array<{ __typename?: 'Lemma', tags: Array<string>, slug: { __typename?: 'Word', word: string, hiragana: string, furigana: Array<{ __typename?: 'Furigana', kanji: string, hiragana: string }>, pitch: Array<{ __typename?: 'Pitch', hiragana: string, pitch: Array<PitchType> }> }, forms: Array<{ __typename?: 'Word', word: string, hiragana: string, furigana: Array<{ __typename?: 'Furigana', kanji: string, hiragana: string }>, pitch: Array<{ __typename?: 'Pitch', hiragana: string, pitch: Array<PitchType> }> }>, senses: Array<{ __typename?: 'Sense', definition: Array<string>, partOfSpeech: Array<string>, tags: Array<string> }>, audio: Array<{ __typename?: 'Audio', type: string, source: string }> }> } | null };


export const GetLemmasDoc = gql`
    query GetLemmas($query: String!) {
  Lemmas(query: $query) {
    lemmas {
      slug {
        word
        hiragana
        furigana {
          kanji
          hiragana
        }
        pitch {
          hiragana
          pitch
        }
      }
      tags
      forms {
        word
        hiragana
        furigana {
          kanji
          hiragana
        }
        pitch {
          hiragana
          pitch
        }
      }
      senses {
        definition
        partOfSpeech
        tags
      }
      audio {
        type
        source
      }
    }
  }
}
    `;
export const GetLemmas = (
            options: Omit<
              WatchQueryOptions<GetLemmasQueryVariables>, 
              "query"
            >
          ): Readable<
            ApolloQueryResult<GetLemmasQuery> & {
              query: ObservableQuery<
                GetLemmasQuery,
                GetLemmasQueryVariables
              >;
            }
          > => {
            const q = client.watchQuery({
              query: GetLemmasDoc,
              ...options,
            });
            var result = readable<
              ApolloQueryResult<GetLemmasQuery> & {
                query: ObservableQuery<
                  GetLemmasQuery,
                  GetLemmasQueryVariables
                >;
              }
            >(
              { data: {} as any, loading: true, error: undefined, networkStatus: 1, query: q },
              (set) => {
                q.subscribe((v: any) => {
                  set({ ...v, query: q });
                });
              }
            );
            return result;
          }
        