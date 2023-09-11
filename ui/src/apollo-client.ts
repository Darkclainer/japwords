import { ApolloClient, HttpLink, InMemoryCache } from '@apollo/client';

const cache = new InMemoryCache({
  addTypename: true,
  typePolicies: {
    AnkiState: {
      // AnkiState is singleton
      keyFields: [],
    },
    AnkiConfig: {
      // AnkiConfig is singleton
      keyFields: [],
    },
  },
});

function getQueryPath(): string {
  let queryPath = '/api/query';
  if (import.meta.env.DEV) {
    queryPath = 'http://' + import.meta.env.VITE_HOST + queryPath;
  }
  return queryPath;
}

const httpLink = new HttpLink({
  uri: getQueryPath(),
});

export default new ApolloClient({
  cache,
  link: httpLink,
});
