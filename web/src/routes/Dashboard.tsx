import styled from 'styled-components';
import { Navbar } from '../components/Navbar';
import { useSelfUser } from '../hooks/useSelfUser';

interface Props {}

const RouteContainer = styled.div`
  display: flex;
  height: 100%;
`;

export const DashboardRoute: React.FC<Props> = ({}) => {
  const selfUser = useSelfUser();
  return (
    <RouteContainer>
      <Navbar />
    </RouteContainer>
  );
};
