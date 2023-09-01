import { createBrowserRouter, RouterProvider, Navigate } from 'react-router-dom';

import Root from './routes/root';
import ErrorPage from './error-page';
import Search, { loader as searchLoader } from './routes/search';
import Anki from './routes/anki';
import { ApolloProvider } from '@apollo/client';
import apolloClient from './apollo-client';
import { ToastContainer } from 'react-toastify';
import AnkiConfig from './routes/anki-config';
import AnkiState from './routes/anki-state';

const router = createBrowserRouter([
  {
    path: '/',
    element: <Root />,
    errorElement: <ErrorPage />,
    children: [
      {
        errorElement: <ErrorPage />,
        children: [
          {
            index: true,
            element: <Navigate to="search" replace={true} />,
          },
          {
            path: 'search/:query?',
            element: <Search />,
            loader: searchLoader,
          },
          {
            path: 'anki',
            element: <Anki />,
            children: [
              {
                path: 'config',
                element: <AnkiConfig />,
              },
              {
                path: 'state',
                element: <AnkiState />,
              },
            ],
          },
        ],
      },
    ],
  },
]);

export default function App() {
  return (
    <ApolloProvider client={apolloClient}>
      <div className="flex flex-col md:max-w-4xl mx-auto min-h-screen space-y-4 py-4 bg-white">
        <RouterProvider router={router} />
      </div>
      <ToastContainer />
    </ApolloProvider>
  );
}
