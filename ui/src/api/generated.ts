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

export type JapaneseWord = {
  __typename?: 'JapaneseWord';
  acents: Scalars['String'];
  audio?: Maybe<Array<Scalars['String']>>;
  examples?: Maybe<Array<Scalars['String']>>;
  furigana?: Maybe<Scalars['String']>;
  hiragana: Scalars['String'];
  kanji: Scalars['String'];
  meaning: Scalars['String'];
};

export type JapaneseWords = {
  __typename?: 'JapaneseWords';
  words?: Maybe<Array<JapaneseWord>>;
};

export type Query = {
  __typename?: 'Query';
  japaneseWords?: Maybe<JapaneseWords>;
};


export type QueryJapaneseWordsArgs = {
  query: Scalars['String'];
};

export type GetJapaneseWordsQueryVariables = Exact<{
  query: Scalars['String'];
}>;


export type GetJapaneseWordsQuery = { __typename?: 'Query', japaneseWords?: { __typename?: 'JapaneseWords', words?: Array<{ __typename?: 'JapaneseWord', kanji: string, furigana?: string | null, hiragana: string, acents: string, meaning: string, audio?: Array<string> | null, examples?: Array<string> | null }> | null } | null };


export const GetJapaneseWordsDoc = gql`
    query GetJapaneseWords($query: String!) {
  japaneseWords(query: $query) {
    words {
      kanji
      furigana
      hiragana
      acents
      meaning
      audio
      examples
    }
  }
}
    `;
export const GetJapaneseWords = (
            options: Omit<
              WatchQueryOptions<GetJapaneseWordsQueryVariables>, 
              "query"
            >
          ): Readable<
            ApolloQueryResult<GetJapaneseWordsQuery> & {
              query: ObservableQuery<
                GetJapaneseWordsQuery,
                GetJapaneseWordsQueryVariables
              >;
            }
          > => {
            const q = client.watchQuery({
              query: GetJapaneseWordsDoc,
              ...options,
            });
            var result = readable<
              ApolloQueryResult<GetJapaneseWordsQuery> & {
                query: ObservableQuery<
                  GetJapaneseWordsQuery,
                  GetJapaneseWordsQueryVariables
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
        