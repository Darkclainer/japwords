import { ApolloProvider } from '@apollo/client';
import { createBrowserRouter, Navigate, RouterProvider } from 'react-router-dom';
import { ToastContainer } from 'react-toastify';

import apolloClient from './apollo-client';
import { HealthStatusProvider } from './contexts/health-status';
import ErrorPage from './error-page';
import AnkiConnectionSettings from './routes/anki/connection-settings';
import AnkiRoot from './routes/anki/root';
import AnkiUserSettings from './routes/anki/user-settings';
import HealthDashboard from './routes/health-dashboard';
import Root from './routes/root';
import Search, { loader as searchLoader } from './routes/search';

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
            path: 'health-dashboard',
            element: <HealthDashboard />,
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
    <>
      <ApolloProvider client={apolloClient}>
        <HealthStatusProvider>
          <div className="mx-auto flex min-h-screen flex-col bg-white px-4 md:max-w-screen-xl">
            <RouterProvider router={router} />
          </div>
        </HealthStatusProvider>
      </ApolloProvider>
      <ToastContainer />
    </>
  );
}
