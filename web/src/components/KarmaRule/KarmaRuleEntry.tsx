import { Guild, KarmaRule } from '../../lib/shinpuru-ts/src';
import { Trans, useTranslation } from 'react-i18next';
import { getActionOptions, getTriggerOptions } from './shared';

import { Container } from '../Container';
import { EmbedWrapper } from '../Embed';
import { ReactComponent as IconDelete } from '../../assets/delete.svg';
import styled from 'styled-components';
import { useMemo } from 'react';

type Props = {
  guild: Guild;
  rule: KarmaRule;
  onRemove: () => void;
};

const RuleContainer = styled(Container)`
  margin: 1em 0;
  display: flex;
  align-items: center;
  justify-content: space-between;

  &:hover > button {
    opacity: 1;
  }

  > button {
    background: none;
    border: none;
    padding: 0;
    cursor: pointer;
    opacity: 0;
    transition: all 0.25s ease;
  }
`;

export const KarmaRuleEntry: React.FC<Props> = ({ guild, rule, onRemove }) => {
  const { t } = useTranslation('components');

  const triggerOptions = useMemo(() => getTriggerOptions(t), [t]);
  const actionOptions = useMemo(() => getActionOptions(t), [t]);

  const triggerElem = (
    <EmbedWrapper value={triggerOptions.find((o) => o.value === rule.trigger)?.display} />
  );
  const valueElem = <EmbedWrapper value={rule.value} />;
  const actionElem = (
    <EmbedWrapper value={actionOptions.find((o) => o.value === rule.action)?.display} />
  );
  const argumentElem = (() => {
    switch (rule.action) {
      case 'SEND_MESSAGE':
        return <EmbedWrapper value={rule.argument} />;
      case 'TOGGLE_ROLE':
        return <EmbedWrapper value={guild.roles?.find((r) => r.id === rule.argument)?.name} />;
      default:
        return <></>;
    }
  })();

  return (
    <RuleContainer>
      <div>
        <Trans
          ns="components"
          i18nKey="karmarule.text"
          components={{
            '1': triggerElem,
            '2': valueElem,
            '3': actionElem,
            '4': argumentElem,
          }}
        />
      </div>
      <button onClick={onRemove}>
        <IconDelete />
      </button>
    </RuleContainer>
  );
};
