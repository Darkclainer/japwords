import { ApolloClient, InMemoryCache, HttpLink } from '@apollo/client';

const cache = new InMemoryCache({
	addTypename: true
});

function getQueryPath(): string {
	let queryPath = '/api/query';
	if (!process.env.PRODUCTION) {
		queryPath = 'http://' + process.env.HOST + queryPath;
	}
	return queryPath;
}

const httpLink = new HttpLink({
	uri: getQueryPath()
});

export default new ApolloClient({
	cache,
	link: httpLink
});
