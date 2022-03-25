import { useNavigate } from 'react-router';
import styled from 'styled-components';
import { Button } from '../components/Button';

interface Props {}

const StartContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1em;
  width: 100%;
  height: 100%;
  justify-content: center;
  align-items: center;
`;

export const StartRoute: React.FC<Props> = () => {
  const nav = useNavigate();

  return (
    <StartContainer>
      <p>This is only a placeholder start page until a shiny start page is presented here. 😉</p>
      <Button onClick={() => nav('/login')}>Go to Login</Button>
    </StartContainer>
  );
};
