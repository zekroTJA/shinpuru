import { isAllowed, isDisallowed } from './util';

import { Loader } from '../Loader';
import { PermTile } from './PermTile';
import styled from 'styled-components';
import { uid } from 'react-uid';

type Props = {
  perms?: string[];
};

const PermsContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.7em;

  > section {
    display: flex;
    flex-wrap: wrap;
    max-width: 30em;
    gap: 0.5em;
  }
`;

export const PermsSimpleList: React.FC<Props> = ({ perms }) => {
  return (
    <PermsContainer>
      {(perms && (
        <>
          <section>
            {perms.filter(isAllowed).map((p) => (
              <PermTile perm={p} key={uid(p)} />
            ))}
          </section>
          <section>
            {perms.filter(isDisallowed).map((p) => (
              <PermTile perm={p} key={uid(p)} />
            ))}
          </section>
        </>
      )) || (
        <>
          <section>
            <Loader width="6ch" height="1.2em" borderRadius="3px" />
            <Loader width="4ch" height="1.2em" borderRadius="3px" />
            <Loader width="5ch" height="1.2em" borderRadius="3px" />
          </section>
          <section>
            <Loader width="5ch" height="1.2em" borderRadius="3px" />
            <Loader width="6ch" height="1.2em" borderRadius="3px" />
            <Loader width="4ch" height="1.2em" borderRadius="3px" />
          </section>
        </>
      )}
    </PermsContainer>
  );
};
