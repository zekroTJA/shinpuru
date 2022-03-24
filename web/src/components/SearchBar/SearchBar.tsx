import styled from 'styled-components';
import { ReactComponent as SearchIcon } from '../../assets/search.svg';
import { Input } from '../Input';

type Props = React.HTMLAttributes<HTMLInputElement> & {
  value?: string;
  onValueChange?: (v: string) => void;
};

const SearchContainer = styled.div`
  margin-left: 2px;
  position: relative;
  width: 100%;

  > svg {
    position: absolute;
    width: 1.5em;
    height: 100%;
    margin-left: 0.3em;
  }

  > input {
    padding-left: 2.2em;
    width: 100%;
  }
`;

export const SearchBar: React.FC<Props> = ({ value, onValueChange = () => {}, ...props }) => {
  return (
    <SearchContainer>
      <SearchIcon />
      <Input value={value} onInput={(e) => onValueChange(e.currentTarget.value)} {...props} />
    </SearchContainer>
  );
};
