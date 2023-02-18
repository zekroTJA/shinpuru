import { Heading } from '../Heading';
import { PropsWithChildren } from 'react';
import { ReactComponent as SPBrand } from '../../assets/sp-brand.svg';
import SPIcon from '../../assets/sp-icon.png';
import { Styled } from '../props';
import styled from 'styled-components';

type Props = PropsWithChildren & Styled & {};

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

  @media (orientation: portrait) {
    width: fit-content;

    ${Brand} > svg {
      display: none;
    }

    ${Heading} {
      display: none;
    }
  }
`;

export const Navbar: React.FC<Props> = ({ children }) => {
  return (
    <StyledNav>
      <Brand>
        <img src={SPIcon} alt="shinpuru Heading" />
        <SPBrand />
      </Brand>
      {children}
    </StyledNav>
  );
};
