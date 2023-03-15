import { Heading } from '../Heading';
import { PropsWithChildren } from 'react';
import styled from 'styled-components';

type Props = PropsWithChildren & {
  title: string;
};

const StyledSection = styled.section`
  > ${Heading} {
    font-size: 0.7em;
  }
`;

export const Section: React.FC<Props> = ({ title, children }) => {
  return (
    <StyledSection>
      <Heading>{title}</Heading>
      {children}
    </StyledSection>
  );
};
