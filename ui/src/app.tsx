import { createBrowserRouter, RouterProvider, Navigate } from 'react-router-dom';

import Root from './routes/root';
import ErrorPage from './error-page';
import Search, { loader as searchLoader } from './routes/search';
import { ApolloProvider } from '@apollo/client';
import apolloClient from './apollo-client';
import { ToastContainer } from 'react-toastify';
import AnkiRoot from './routes/anki/root';
import AnkiConnectionSettings from './routes/anki/connection-settings';
import AnkiUserSettings from './routes/anki/user-settings';

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
            element: <AnkiRoot />,
            children: [
              {
                index: true,
                element: <Navigate to="connection-settings" replace={true} />,
              },
              {
                path: 'connection-settings',
                element: <AnkiConnectionSettings />,
              },
              {
                path: 'user-settings',
                element: <AnkiUserSettings />,
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
      <div className="flex flex-col md:max-w-4xl mx-auto min-h-screen space-y-4 p-4 bg-white">
        <RouterProvider router={router} />
      </div>
      <ToastContainer />
    </ApolloProvider>
  );
}
