import { ApolloError } from '@apollo/client';
import { GraphQLError } from 'graphql';
import { createElement, ReactNode } from 'react';
import { IconProps, toast, ToastContent, ToastOptions } from 'react-toastify';

import { ErrorIcon, OKIcon, WarningIcon } from '../components/Icons/StatusIcon';

type FunctionType<T> = T extends (...a: infer U) => infer R ? (...a: U) => R : never;
export type ToastFunction<T = unknown> = FunctionType<typeof toast<T>>;

// styledToast calls toast from react-toastify with overrided styles
export function styledToast<T = unknown>(content: ToastContent<T>, opts: ToastOptions) {
  return toast(content, {
    icon: getToastIcon,
    position: 'bottom-right',
    autoClose: 2000,
    ...opts,
  });
}

function getToastIcon({ type }: IconProps): ReactNode {
  const props = { className: 'w-5 h-5' };
  switch (type) {
    case 'success':
      return createElement(OKIcon, props);
    case 'warning':
      return createElement(WarningIcon, props);
    case 'error':
    case 'default':
      return createElement(ErrorIcon, props);
  }
  return null;
}

// apolloErrorToast tries to make sensible error message for notification.
export function apolloErrorToast<T = unknown>(
  error: ApolloError | readonly GraphQLError[],
  action: string,
  opts?: { toast?: ToastFunction<T>; toastOptions?: ToastOptions },
) {
  const message = getApolloErrorMessage(error);
  const toast = opts?.toast ?? styledToast<T>;
  if (message) {
    action = `${action} ${message}.`;
  }

  return toast(action, { type: 'error', ...opts?.toastOptions });
}

function getApolloErrorMessage(error: ApolloError | readonly GraphQLError[]): string | undefined {
  if (error instanceof ApolloError) {
    if (error.networkError && !('statusCode' in error.networkError)) {
      return 'There is no connection to server';
    }
    return error.message;
  } else {
    for (const err of error) {
      return err.message;
    }
  }
}
