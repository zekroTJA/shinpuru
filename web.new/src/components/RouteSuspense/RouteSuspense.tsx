import { Suspense } from 'react';

type Props = {};

export const RouteSuspense: React.FC<Props> = ({ children }) => {
  // TODO: Use better and fancier fallback
  return <Suspense fallback={<>loading ...</>}>{children}</Suspense>;
};
