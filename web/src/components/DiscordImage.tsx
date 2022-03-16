import styled, { css } from 'styled-components';
import DCLogoURL from '../assets/dc-logo.svg';

type ImgProps = {
  round?: boolean;
};

type Props = React.ImgHTMLAttributes<any> & ImgProps;

const StyledImg = styled.img<ImgProps>`
  ${(p) =>
    p.round &&
    css`
      border-radius: 100%;
    `}
`;

export const DiscordImage: React.FC<Props> = ({ src, ...props }) => {
  return <StyledImg src={!!src ? src : DCLogoURL} {...props} />;
};
