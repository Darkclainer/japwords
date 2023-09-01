import { ApolloClient, InMemoryCache, HttpLink } from '@apollo/client';

const cache = new InMemoryCache({
  addTypename: true,
});

function getQueryPath(): string {
  let queryPath = '/api/query';
  if (import.meta.env.DEV) {
    queryPath = 'http://' + import.meta.env.VITE_HOST + queryPath;
  }
  console.log('query path', queryPath, import.meta.env);
  return queryPath;
}

const httpLink = new HttpLink({
  uri: getQueryPath(),
});

export default new ApolloClient({
  cache,
  link: httpLink,
});
