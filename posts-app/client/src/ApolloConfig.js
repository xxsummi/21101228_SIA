import { split, HttpLink, from } from '@apollo/client';
import { getMainDefinition } from '@apollo/client/utilities';
import { GraphQLWsLink } from '@apollo/client/link/subscriptions';
import { createClient } from 'graphql-ws';
import { ApolloClient, InMemoryCache, } from '@apollo/client'

const httpLink = new HttpLink({ uri: 'http://localhost:4002/graphql' });

const wsLink = new GraphQLWsLink(createClient({ url: 'ws://localhost:4002/graphql' }));

const link = split(
  ({ query }) => {
    const definition = getMainDefinition(query);
    return (
      definition.kind === 'OperationDefinition' &&
      definition.operation === 'subscription'
    );
  },
  wsLink,
  httpLink
);

const client = new ApolloClient({ link, cache: new InMemoryCache() });

export default client;
