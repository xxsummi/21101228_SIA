// const { ApolloServer, gql } = require('apollo-server');
// const { PrismaClient } = require('@prisma/client');

// const prisma = new PrismaClient();

// const typeDefs = gql`
//   type Post {
//     id: ID!
//     title: String!
//     content: String!
//   }

//   type Query {
//     posts: [Post]
//     post(id: ID!): Post
//   }

//   type Mutation {
//     createPost(title: String!, content: String!): Post
//     updatePost(id: ID!, title: String, content: String): Post
//     deletePost(id: ID!): Post
//   }
// `;

// const resolvers = {
//   Query: {
//     posts: () => prisma.post.findMany(),
//     post: (_, { id }) => prisma.post.findUnique({ where: { id: parseInt(id) } })
//   },
//   Mutation: {
//     createPost: (_, { title, content }) => prisma.post.create({ data: { title, content } }),
//     updatePost: (_, { id, title, content }) => prisma.post.update({
//       where: { id: parseInt(id) },
//       data: { title, content }
//     }),
//     deletePost: (_, { id }) => prisma.post.delete({ where: { id: parseInt(id) } })
//   }
// };

// const server = new ApolloServer({ typeDefs, resolvers });

// server.listen({ port: 4002 }).then(({ url }) => {
//   console.log(`🚀 Posts service running at ${url}`);
// });


// const { ApolloServer, gql } = require('apollo-server');
// const { PubSub } = require('graphql-subscriptions');
// const { PrismaClient } = require('@prisma/client');

// const prisma = new PrismaClient();
// const pubsub = new PubSub();

// const POST_ADDED = 'POST_ADDED';

// const typeDefs = gql`
//   type Post {
//     id: ID!
//     title: String!
//     content: String!
//   }

//   type Query {
//     posts: [Post]
//     post(id: ID!): Post
//   }

//   type Mutation {
//     createPost(title: String!, content: String!): Post
//     updatePost(id: ID!, title: String, content: String): Post
//     deletePost(id: ID!): Post
//   }

//   type Subscription {
//     postAdded: Post
//   }
// `;

// const resolvers = {
//   Query: {
//     posts: () => prisma.post.findMany(),
//     post: (_, { id }) => prisma.post.findUnique({ where: { id: parseInt(id) } })
//   },
//   Mutation: {
//     createPost: async (_, { title, content }) => {
//       const post = await prisma.post.create({ data: { title, content } });
//       pubsub.publish(POST_ADDED, { postAdded: post });
//       return post;
//     },
//     updatePost: (_, { id, title, content }) => prisma.post.update({
//       where: { id: parseInt(id) },
//       data: { title, content }
//     }),
//     deletePost: (_, { id }) => prisma.post.delete({ where: { id: parseInt(id) } })
//   },
//   Subscription: {
//     postAdded: {
//       subscribe: () => pubsub.asyncIterator([POST_ADDED])
//     }
//   }
// };

// const server = new ApolloServer({
//   typeDefs,
//   resolvers
// });

// server.listen({ port: 4002 }).then(({ url }) => {
//   console.log(`🚀 Posts service running at ${url}`);
// });


// Updated version
const { ApolloServer, gql } = require('apollo-server');
const { PubSub } = require('graphql-subscriptions');
const { PrismaClient } = require('@prisma/client');

const prisma = new PrismaClient();
const pubsub = new PubSub();

const POST_ADDED = 'POST_ADDED';

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
    createPost(title: String!, content: String!, userId: String!): Post
    updatePost(id: ID!, title: String, content: String): Post
    deletePost(id: ID!): Post
  }

  type Subscription {
    postAdded: Post
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
      subscribe: () => pubsub.asyncIterator([POST_ADDED])
    }
  }
};

const server = new ApolloServer({
  typeDefs,
  resolvers
});

server.listen({ port: 4002 }).then(({ url }) => {
  console.log(`🚀 Posts service running at ${url}`);
});
