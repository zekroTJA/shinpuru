import styled from 'styled-components';
import { Element } from './Select';

type Props<T> = {
  value: Element<T>;
};

const StyledDiv = styled.div``;

export const Option = <T extends unknown>({ value }: Props<T>) => {
  return <StyledDiv key={value.id}>{value.display}</StyledDiv>;
};
