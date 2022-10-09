import { ApolloClient, InMemoryCache, HttpLink, ApolloLink } from '@apollo/client';

const cache = new InMemoryCache({
	addTypename: true
});

const httpLink = new HttpLink({
	uri: 'http://' + (process.env.HOST || '') + '/api/query'
});

export default new ApolloClient({
	cache,
	link: httpLink
});
