import { useEffect, useRef } from 'react';
import { Id as ToastId, toast, ToastContent, ToastOptions } from 'react-toastify';

import { styledToast, ToastFunction } from '../lib/styled-toast';

interface Id {
  prefix: number;
  index: number;
}

function idToToastId(id: Id): ToastId {
  return `${id.prefix}-${id.index}`;
}

let prefixCounter = 0;
function newId(): Id {
  const id = {
    prefix: prefixCounter,
    index: 0,
  };
  prefixCounter++;
  return id;
}

// provide wrapper around toast function that ties notification with component, and
// limits number of notifications to one
export function useToastify<T = unknown>(defaultOptions?: ToastOptions): ToastFunction<T> {
  const state = useRef<Id>();
  if (!state.current) {
    state.current = newId();
  }
  useEffect(() => {
    return () => {
      toast.dismiss(idToToastId(state.current!));
    };
  }, []);
  return (content: ToastContent<T>, options?: ToastOptions) => {
    toast.dismiss(idToToastId(state.current!));
    state.current!.index++;
    const currentOptions = { toastId: idToToastId(state.current!), ...defaultOptions, ...options };
    return styledToast(content, currentOptions);
  };
}
