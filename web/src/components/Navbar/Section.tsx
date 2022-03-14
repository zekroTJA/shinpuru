import styled from 'styled-components';

interface Props {
  title: string;
}

const StyledSection = styled.section`
  > h4 {
    display: block;
    font-family: 'Cantarell', sans-serif;
    font-size: 0.7em;
    text-transform: uppercase;
    opacity: 0.6;
  }
`;

export const Section: React.FC<Props> = ({ title, children }) => {
  return (
    <StyledSection>
      <h4>{title}</h4>
      {children}
    </StyledSection>
  );
};
