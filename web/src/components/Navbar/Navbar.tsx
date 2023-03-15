import { EntryContainer } from './EntryContainer';
import { Heading } from '../Heading';
import { MAX_WIDTH } from '../MaxWidthContainer';
import { PropsWithChildren } from 'react';
import { PropsWithStyle } from '../props';
import { ReactComponent as SPBrand } from '../../assets/sp-brand.svg';
import SPIcon from '../../assets/sp-icon.png';
import styled from 'styled-components';

type Props = PropsWithChildren & PropsWithStyle & {};

export const NAVBAR_WIDTH = '15rem';

const Brand = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  table-layout: fixed;

  > img {
    width: 38px;
    height: 38px;
  }

  > svg {
    width: 100%;
    height: 38px;
    justify-content: flex-start;
  }
`;

const StyledNav = styled.nav`
  display: flex;
  flex-direction: column;
  background-color: ${(p) => p.theme.background2};
  margin: 1rem 0 1rem 1rem;
  padding: 1rem;
  border-radius: 12px;
  width: 30vw;
  max-width: 15rem;

  @media (orientation: portrait) and (max-width: calc(${NAVBAR_WIDTH} + ${MAX_WIDTH})) {
    width: fit-content;

    ${Brand} > svg {
      display: none;
    }

    ${Heading} {
      display: none;
    }
  }

  > ${EntryContainer} {
    &::-webkit-scrollbar {
      width: 5px;
    }
  }
`;

export const Navbar: React.FC<Props> = ({ children, ...props }) => {
  return (
    <StyledNav {...props}>
      <Brand>
        <img src={SPIcon} alt="shinpuru Heading" />
        <SPBrand />
      </Brand>
      {children}
    </StyledNav>
  );
};
