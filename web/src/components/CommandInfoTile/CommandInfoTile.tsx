import { CommandInfo, CommandOptionType } from '../../lib/shinpuru-ts/src';

import { ArgsTable } from './ArgsTable';
import { Embed } from '../Embed';
import { Flex } from '../Flex';
import { ReactComponent as MailIcon } from '../../assets/mail.svg';
import styled from 'styled-components';
import { uid } from 'react-uid';
import { useTranslation } from 'react-i18next';

type Props = {
  cmd: CommandInfo;
};

const CommandContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1em;
  padding: 1em;
  border-radius: 12px;
  background-color: ${(p) => p.theme.background2};

  > header {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 0.5em;
  }

  h3,
  h4,
  h5 {
    margin: 0;
  }

  ul {
    margin: 0;
  }

  h4 {
    text-transform: uppercase;
    opacity: 0.6;
    font-size: 0.7rem;
    margin-bottom: 0.5em;
  }

  h5 {
    font-weight: bold;
  }
`;

const SubCommandContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.5em;
  padding: 0.5em;
  border-radius: 8px;

  background-color: ${(p) => p.theme.background3};
`;

const DMCapable = styled.div`
  color: ${(p) => p.theme.lime};
  border-radius: 3px;
  padding: 0.2em;
  display: flex;
  align-items: center;

  &:hover {
    background-color: ${(p) => p.theme.background3};
  }
`;

export const CommandInfoTile: React.FC<Props> = ({ cmd }) => {
  const { t } = useTranslation('components', { keyPrefix: 'commandtile' });

  const options = cmd.options?.filter((o) => o.type !== CommandOptionType.SUBCOMMAND);
  const subCommands = cmd.options?.filter((o) => o.type === CommandOptionType.SUBCOMMAND);

  return (
    <CommandContainer>
      <header>
        <h3 id={cmd.name}>{cmd.name}</h3>
        <Embed>{cmd.domain}</Embed>
        <Embed>{cmd.version}</Embed>
        {cmd.dm_capable && (
          <DMCapable title={t('dmcapable')}>
            <MailIcon />
          </DMCapable>
        )}
      </header>
      {cmd.description && (
        <section>
          <h4>{t('description')}</h4>
          <span>{cmd.description}</span>
        </section>
      )}
      {cmd.subdomains?.length > 0 && (
        <section>
          <h4>{t('subdomains')}</h4>
          <ul>
            {cmd.subdomains.map((sd) => (
              <li key={uid(sd)}>
                <Embed>
                  {cmd.domain}.{sd.term}
                </Embed>
              </li>
            ))}
          </ul>
        </section>
      )}
      {options?.length > 0 && (
        <section>
          <h4>{t('arguments')}</h4>
          <ArgsTable options={options} />
        </section>
      )}
      {subCommands?.length > 0 && (
        <section>
          <h4>{t('subcommands')}</h4>
          <Flex direction="column" gap="1em">
            {subCommands.map((sc) => (
              <SubCommandContainer key={uid(sc)}>
                <h5>{sc.name}</h5>
                <span>{sc.description}</span>
                {sc.options && <ArgsTable options={sc.options} />}
              </SubCommandContainer>
            ))}
          </Flex>
        </section>
      )}
    </CommandContainer>
  );
};
