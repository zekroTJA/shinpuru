import { Suspense } from 'react';

interface Props {}

export const RouteSuspense: React.FC<Props> = ({ children }) => {
  // TODO: Use better and fancier fallback
  return <Suspense fallback={<>loading ...</>}>{children}</Suspense>;
};
