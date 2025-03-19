import React, { useEffect, useState } from 'react';
import { ApolloClient, InMemoryCache, ApolloProvider, useQuery, useSubscription, gql, useMutation } from '@apollo/client';

const client = new ApolloClient({
  uri: 'http://localhost:4002', // Connect to posts service
  cache: new InMemoryCache(),
});

const GET_POSTS = gql`
  query GetPosts {
    posts {
      id
      title
      content
    }
  }
`;

const POST_SUBSCRIPTION = gql`
  subscription OnPostAdded {
    postAdded {
      id
      title
      content
    }
  }
`;

const CREATE_POST = gql`
  mutation CreatePost($title: String!, $content: String!) {
    createPost(title: $title, content: $content) {
      id
      title
      content
    }
  }
`;

const Posts = () => {
  const { data, loading, error, refetch } = useQuery(GET_POSTS);
  const [createPost] = useMutation(CREATE_POST);

  // Subscription to update table in real-time
  useSubscription(POST_SUBSCRIPTION, {
    onData: ({ data }) => {
      console.log('New post added:', data.data.postAdded);
      refetch(); // Refresh data table
    }
  });

  const handleCreatePost = async () => {
    await createPost({
      variables: {
        title: `New Post ${Date.now()}`,
        content: 'This is a new post.'
      }
    });
  };

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error: {error.message}</p>;

  return (
    <div>
      <h2>Posts</h2>
      <button onClick={handleCreatePost}>Create Post</button>
      <table border="1" style={{ marginTop: 10 }}>
        <thead>
          <tr>
            <th>ID</th>
            <th>Title</th>
            <th>Content</th>
          </tr>
        </thead>
        <tbody>
          {data?.posts.map(post => (
            <tr key={post.id}>
              <td>{post.id}</td>
              <td>{post.title}</td>
              <td>{post.content}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

function App() {
  return (
    <ApolloProvider client={client}>
      <Posts />
    </ApolloProvider>
  );
}

export default App;
