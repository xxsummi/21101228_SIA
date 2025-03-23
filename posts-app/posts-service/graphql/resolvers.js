const { PrismaClient } = require('@prisma/client');
const { PubSub } = require('apollo-server');

const prisma = new PrismaClient();
const pubsub = new PubSub();

const POST_CREATED = "POST_CREATED";

const resolvers = {
  Query: {
    posts: async () => {
      return await prisma.post.findMany();
    },
  },

  Mutation: {
    createPost: async (_, { title, content, userId }) => {
      const post = await prisma.post.create({
        data: {
          title,
          content,
          userId,
        },
      });

      pubsub.publish(POST_CREATED, { postCreated: post }); // Publish event

      return post;
    },
  },

  Subscription: {
    postCreated: {
      subscribe: () => pubsub.asyncIterator([POST_CREATED]), // Subscribe to new posts
    },
  },
};

module.exports = resolvers;
