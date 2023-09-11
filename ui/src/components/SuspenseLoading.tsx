import { Suspense } from 'react';

export type SuspenseLoadingProps = {
  className?: string;
  children: React.ReactNode;
};
export default function SuspenseLoading({ className, children }: SuspenseLoadingProps) {
  return (
    <Suspense fallback={<p className={className ?? 'text-2xl'}>Loading...</p>}>{children}</Suspense>
  );
}
