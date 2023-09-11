import { isApolloError } from '@apollo/client';
import { useContext, useEffect } from 'react';
import { isRouteErrorResponse, Link, useNavigate, useRouteError } from 'react-router-dom';

import { HealthStatusContext } from './contexts/health-status';
import { isHealthStatusThrownError } from './model/health-status';

export default function ErrorPage() {
  const navigate = useNavigate();

  return (
    <div className="flex flex-col place-items-center">
      <h1 className="text-3xl text-dark-blue font-bold p-4 pt-12">Error has occured!</h1>
      <ErrorDetails />
      <h2 className="text-xl p-4">
        <button onClick={() => navigate(0)} className="text-blue">
          Refresh page
        </button>{' '}
        or{' '}
        <Link to="/" className="text-blue">
          return to home page
        </Link>
      </h2>
    </div>
  );
}

function ErrorDetails() {
  const error: unknown = useRouteError();
  console.log(error);
  let errorMessage: string;
  if (isRouteErrorResponse(error)) {
    errorMessage = error.error?.message || error.statusText;
  } else if (isHealthStatusThrownError(error)) {
    if (error.kind == 'Disconnected') {
      return <ApolloNetworkError />;
    }
    errorMessage = 'Unknown error response from server';
  } else if (error instanceof Error) {
    if (isApolloError(error) && error.networkError && !('statusCode' in error.networkError)) {
      return <ApolloNetworkError />;
    } else {
      errorMessage = error.message;
    }
  } else if (typeof error === 'string') {
    errorMessage = error;
  } else {
    errorMessage = 'Unknown error';
  }
  return (
    <>
      <p>Sorry, an unexpected error has occurred.</p>
      <p>
        <i>{errorMessage}</i>
      </p>
    </>
  );
}

// error can be useful later for more rich erorr messages
// eslint-disable-next-line @typescript-eslint/no-unused-vars
function ApolloNetworkError() {
  const healthStatus = useContext(HealthStatusContext);
  const connectionOk = healthStatus.kind === 'Ok';
  const navigate = useNavigate();
  useEffect(() => {
    if (connectionOk) {
      const timeoutId = setTimeout(() => {
        navigate(0);
      }, 1000);
      return () => clearTimeout(timeoutId);
    }
  }, [connectionOk]);
  return (
    <>
      <p>Request to server failed!</p>
      <p>
        {connectionOk
          ? 'Connection to server was established, try to refresh the page.'
          : 'Currently there seems to be no connection to server, try to enable server and refresh the page.'}
      </p>
    </>
  );
}
