import styled from 'styled-components';
import { Flex } from '../Flex';
import { LinearGradient } from '../styleParts';

enum KarmaType {
  VERY_LOW,
  LOW,
  ZERO,
  HIGH,
  VERY_HIGH,
}

const getKarmaType = (v: number) => {
  if (v > 100) return KarmaType.VERY_HIGH;
  if (v > 0) return KarmaType.HIGH;
  if (v === 0) return KarmaType.ZERO;
  if (v < 0) return KarmaType.LOW;
  return KarmaType.VERY_LOW;
};

type Props = React.HTMLAttributes<HTMLDivElement> & {
  heading?: string;
  points: number;
};

const StyledDiv = styled(Flex)<{ type: KarmaType }>`
  flex-direction: column;
  padding: 1em;
  border-radius: 8px;

  ${(p) => {
    switch (p.type) {
      case KarmaType.VERY_HIGH:
        return LinearGradient(p.theme.blurple);
      case KarmaType.HIGH:
        return LinearGradient(p.theme.green);
      case KarmaType.ZERO:
        return LinearGradient(p.theme.yellow);
      case KarmaType.LOW:
        return LinearGradient(p.theme.orange);
      default:
        return LinearGradient(p.theme.red);
    }
  }};

  > h6 {
    margin: 0 0 1em 0;
    text-transform: uppercase;
    letter-spacing: 0.2ch;
    font-weight: 500;
    opacity: 0.9;
  }
`;

export const KarmaTile: React.FC<Props> = ({ points, heading, ...props }) => {
  return (
    <StyledDiv type={getKarmaType(points)} {...props}>
      <h6>{heading}</h6>
      <span>{points}</span>
    </StyledDiv>
  );
};
