import { Flex } from '../Flex';
import { LinearGradient } from '../styleParts';
import styled from 'styled-components';

type HintProps = {
  color?: string;
  icon?: string | JSX.Element;
};

type Props = HintProps & React.HTMLAttributes<HTMLDivElement>;

const Div = styled(Flex)<HintProps>`
  ${(p) => LinearGradient(p.color ?? p.theme.blurple)};

  border-radius: 12px;
  padding: 0.7em;

  > svg,
  > img {
    height: 1.5em;
    width: auto;
    margin-right: 0.5em;
  }
`;

export const Hint: React.FC<Props> = ({ icon, children, ...props }) => {
  const _icon = typeof icon === 'string' ? <img src={icon} alt="" /> : icon;
  return (
    <Div {...props}>
      {_icon ?? <></>}
      <span>{children}</span>
    </Div>
  );
};
