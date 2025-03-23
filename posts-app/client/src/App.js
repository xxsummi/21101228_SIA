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
      userId
      title
      content
    }
  }
`;

const POST_SUBSCRIPTION = gql`
  subscription OnPostAdded {
    postAdded {
      id
      userId
      title
      content
    }
  }
`;

const Posts = () => {
  const { data, loading, error, refetch } = useQuery(GET_POSTS);
  
  // Subscription to update table in real-time
  useSubscription(POST_SUBSCRIPTION, {
    onData: ({ data }) => {
      console.log('New post added:', data.data.postAdded);
      refetch(); // Refresh data table
    }
  });

  if (loading) return (
    <div className="loading-container">
      <div className="loading-spinner"></div>
      <p>Loading posts...</p>
    </div>
  );
  
  if (error) return (
    <div className="error-container">
      <p>Error: {error.message}</p>
    </div>
  );

  return (
    <div className="posts-container">
      <h1 className="page-title">Posts</h1>
      
      <div className="table-container">
        <table className="posts-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>User ID</th>
              <th>Title</th>
              <th>Content</th>
            </tr>
          </thead>
          <tbody>
            {data?.posts.map(post => (
              <tr key={post.id}>
                <td>{post.id}</td>
                <td>{post.userId || 'Anonymous'}</td>
                <td>{post.title}</td>
                <td>{post.content}</td>
              </tr>
            ))}
            {data?.posts.length === 0 && (
              <tr>
                <td colSpan="4" className="no-posts">No posts available</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};

function App() {
  return (
    <ApolloProvider client={client}>
      <div className="app">
        <header className="app-header">
          <h1> Posts Management System</h1>
        </header>
        <main>
          <Posts />
        </main>
        <footer className="app-footer">
          <p>© {new Date().getFullYear()} Posts Management System</p>
        </footer>
      </div>
      
      <style jsx global>{`
        @import url('https://fonts.googleapis.com/css2?family=Lora:wght@400;700&display=swap');
        
        :root {
          --pink-light: #ffb6c1;
          --pink-medium: #f8a5c2;
          --purple-light: #c8a2c8;
          --yellow-light: #fffacd;
          --blue-light: #add8e6;
          --text-dark: #5a4a5a;
          --text-light: #ffffff;
        }
        
        * {
          box-sizing: border-box;
          margin: 0;
          padding: 0;
          font-family: 'Lora', serif;
        }
        
        body {
          background: linear-gradient(135deg, var(--pink-light), var(--purple-light), var(--blue-light), var(--yellow-light));
          background-size: 400% 400%;
          animation: gradient 15s ease infinite;
          color: var(--text-dark);
          min-height: 100vh;
        }
        
        @keyframes gradient {
          0% { background-position: 0% 50%; }
          50% { background-position: 100% 50%; }
          100% { background-position: 0% 50%; }
        }
        
        .app {
          min-height: 100vh;
          display: flex;
          flex-direction: column;
        }
        
        .app-header {
          background: rgba(255, 255, 255, 0.2);
          backdrop-filter: blur(10px);
          color: var(--text-dark);
          padding: 1.5rem 2rem;
          text-align: center;
          font-family: 'Lora', serif;
          font-size: 1.2rem;
          box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
        }
        
        .app-header h1 {
          font-family: 'Lora', serif;
          font-weight: 700;
          font-size: 2.5rem;
          color: var(--text-dark);
          text-shadow: 1px 1px 2px rgba(255, 255, 255, 0.8);
        }
        
        main {
          flex: 1;
          padding: 2rem;
          max-width: 1200px;
          margin: 0 auto;
          width: 100%;
          display: flex;
          justify-content: center;
          align-items: center;
        }
        
        .app-footer {
          background: rgba(255, 255, 255, 0.2);
          backdrop-filter: blur(10px);
          color: var(--text-dark);
          text-align: center;
          padding: 1rem;
          margin-top: auto;
          font-size: 0.9rem;
        }
        
        .page-title {
          font-family: 'Lora', serif;
          color: var(--text-dark);
          margin-bottom: 2rem;
          font-size: 3.5rem;
          text-align: center;
          font-weight: 700;
          text-shadow: 2px 2px 4px rgba(255, 255, 255, 0.6);
        }
        
        .posts-container {
          background: rgba(255, 255, 255, 0.6);
          backdrop-filter: blur(10px);
          border-radius: 20px;
          padding: 2.5rem;
          box-shadow: 0 8px 30px rgba(0, 0, 0, 0.1);
          width: 100%;
        }
        
        .table-container {
          overflow-x: auto;
          border-radius: 12px;
          box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
        }
        
        .posts-table {
          width: 100%;
          border-collapse: collapse;
          overflow: hidden;
          border-radius: 12px;
        }
        
        .posts-table thead {
          background: linear-gradient(90deg, var(--pink-medium), var(--purple-light));
          color: var(--text-light);
        }
        
        .posts-table th {
          text-align: left;
          padding: 1.2rem;
          font-weight: 600;
          letter-spacing: 0.5px;
        }
        
        .posts-table td {
          padding: 1.2rem;
          border-bottom: 1px solid rgba(255, 255, 255, 0.5);
        }
        
        .posts-table tbody tr {
          background-color: rgba(255, 255, 255, 0.4);
          transition: all 0.2s ease;
        }
        
        .posts-table tbody tr:nth-child(even) {
          background-color: rgba(255, 255, 255, 0.2);
        }
        
        .posts-table tbody tr:hover {
          background-color: rgba(255, 255, 255, 0.7);
        }
        
        .loading-container {
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          padding: 3rem;
          color: var(--text-dark);
        }
        
        .loading-spinner {
          border: 4px solid rgba(200, 162, 200, 0.3);
          border-top: 4px solid var(--purple-light);
          border-radius: 50%;
          width: 50px;
          height: 50px;
          animation: spin 1s linear infinite;
          margin-bottom: 1.5rem;
        }
        
        @keyframes spin {
          0% { transform: rotate(0deg); }
          100% { transform: rotate(360deg); }
        }
        
        .error-container {
          background: rgba(255, 235, 238, 0.7);
          color: #d81b60;
          padding: 1.5rem;
          border-radius: 12px;
          text-align: center;
          box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
        }
        
        .no-posts {
          text-align: center;
          color: #777;
          font-style: italic;
          padding: 2rem;
        }
      `}</style>
    </ApolloProvider>
  );
}

export default App;
