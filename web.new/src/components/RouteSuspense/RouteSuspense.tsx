import { PropsWithChildren, Suspense } from 'react';

type Props = PropsWithChildren & {};

export const RouteSuspense: React.FC<Props> = ({ children }) => {
  // TODO: Use better and fancier fallback
  return <Suspense fallback={<>loading ...</>}>{children}</Suspense>;
};
