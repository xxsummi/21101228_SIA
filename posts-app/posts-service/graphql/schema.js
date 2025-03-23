const { gql } = require('apollo-server');

const typeDefs = gql`
  type Post {
    id: ID!
    title: String!
    content: String!
    userId: String!
  }

  type Query {
    posts: [Post!]
  }

  type Mutation {
    createPost(title: String!, content: String!, userId: String!): Post!
  }

  type Subscription {
    postCreated: Post!
  }
`;

module.exports = typeDefs;
