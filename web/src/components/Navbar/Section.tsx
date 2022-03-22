import styled from 'styled-components';
import { Heading } from '../Heading';

interface Props {
  title: string;
}

const StyledHeading = styled(Heading)`
  font-size: 0.7em;
  margin-top: 2em;
`;

export const Section: React.FC<Props> = ({ title, children }) => {
  return (
    <section>
      <StyledHeading>{title}</StyledHeading>
      {children}
    </section>
  );
};
