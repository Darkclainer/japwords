import { ApolloError, useQuery } from '@apollo/client';
import { createContext, ReactNode, useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetHealthStatusQuery } from '../api/__generated__/graphql';
import { GET_HEALTH_STATUS } from '../api/health-status';
import { HealthStatus, healthStatusFromGql } from '../model/health-status';

const MIN_LOAD_TIMEOUT = 200;
const STATUS_REFRESH_TIMEOUT = 5000;

export const HealthStatusContext = createContext<HealthStatus>({
  kind: 'Loading',
});

export const HealthStatusReloadContext = createContext<() => void>(() => {});

export function HealthStatusProvider({ children }: { children: ReactNode }) {
  const [reloading, setReloading] = useState(true);
  const { data, error, loading, refetch } = useQuery(GET_HEALTH_STATUS, {
    pollInterval: STATUS_REFRESH_TIMEOUT,
    errorPolicy: 'none',
    fetchPolicy: 'network-only',
    notifyOnNetworkStatusChange: true,
  });
  const [lastData, lastError] = usePersistLoaded(data, error, loading);
  const reload = useCallback(() => {
    setReloading(true);
    refetch();
  }, []);
  const loadStartTime = useRef<number>();
  useEffect(() => {
    if (loading) {
      loadStartTime.current = new Date().getTime();
    }
    if (!loading && reloading) {
      const now = new Date().getTime();
      const timePassed = now - (loadStartTime.current ?? now);
      const delay = Math.max(MIN_LOAD_TIMEOUT - timePassed, 0);
      const timeoutid = setTimeout(() => {
        setReloading(false);
      }, delay);
      return () => clearTimeout(timeoutid);
    }
  }, [loading]);
  const healthStatus: HealthStatus = useMemo(() => {
    return reloading ? { kind: 'Loading' } : healthStatusFromGql(lastData, lastError);
  }, [lastData, lastError, reloading]);
  console.log(healthStatus);
  return (
    <HealthStatusReloadContext.Provider value={reload}>
      <HealthStatusContext.Provider value={healthStatus}>{children}</HealthStatusContext.Provider>
    </HealthStatusReloadContext.Provider>
  );
}

function usePersistLoaded(
  data: GetHealthStatusQuery | undefined,
  error: ApolloError | undefined,
  loading: boolean,
): [GetHealthStatusQuery | undefined, ApolloError | undefined] {
  const lastState = useRef<[GetHealthStatusQuery | undefined, ApolloError | undefined]>();
  if (!loading) {
    lastState.current = [data, error];
  } else if (lastState.current) {
    return lastState.current;
  }
  return [data, error];
}
