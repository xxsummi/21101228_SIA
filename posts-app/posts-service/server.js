const { createServer } = require('http');
const { ApolloServer , gql} = require('apollo-server-express');
const express = require('express');
const { makeExecutableSchema } = require('@graphql-tools/schema');
const { WebSocketServer } = require('ws');
const { useServer } = require('graphql-ws/lib/use/ws');
const { PubSub } = require('graphql-subscriptions');
const { PrismaClient } = require('@prisma/client');

const prisma = new PrismaClient();
const pubsub = new PubSub();
const POST_ADDED = "POST_ADDED";
const typeDefs = gql`
  type Post {
    id: ID!
    title: String!
    content: String!
    userId: String!
  }

  type Query {
    posts: [Post]
    post(id: ID!): Post
  }

  type Mutation {
    createPost(title: String!, content: String!, userId: String!): Post!
    updatePost(id: ID!, title: String, content: String): Post!
    deletePost(id: ID!): Post!
  }

  type Subscription {
    postAdded: Post!
  }
`;

const resolvers = {
  Query: {
    posts: () => prisma.post.findMany(),
    post: (_, { id }) => prisma.post.findUnique({ where: { id: parseInt(id) } })
  },
  Mutation: {
    createPost: async (_, { title, content, userId }) => {
      const post = await prisma.post.create({ data: { title, content, userId } });
      pubsub.publish(POST_ADDED, { postAdded: post });
      return post;
    },
    updatePost: (_, { id, title, content }) => prisma.post.update({
      where: { id: parseInt(id) },
      data: { title, content }
    }),
    deletePost: (_, { id }) => prisma.post.delete({ where: { id: parseInt(id) } })
  },
  Subscription: {
    postAdded: {
      subscribe: () => pubsub.asyncIterableIterator([POST_ADDED])
    }
  }
};

const app = express();
const httpServer = createServer(app);

const schema = makeExecutableSchema({ typeDefs, resolvers });

const wsServer = new WebSocketServer({
  server: httpServer,
  path: '/graphql',
});

useServer({ schema, context: () => ({ pubsub, prisma }) }, wsServer);

const server = new ApolloServer({
  schema,
  context: () => ({ pubsub, prisma }),
  plugins: [
    {
      async serverWillStart() {
        return {
          async drainServer() {
            wsServer.close();
          },
        };
      },
    },
  ],
});

(async () => {
  await server.start();
  server.applyMiddleware({ app });

  httpServer.listen(4002, () => {
    console.log(`🚀 Posts service running at http://localhost:4002${server.graphqlPath}`);
    console.log(`📡 Subscriptions available at ws://localhost:4002/graphql`);
  });
})();
