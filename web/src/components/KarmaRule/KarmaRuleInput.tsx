import { Guild, KarmaRule, KarmaRuleAction, KarmaRuleTrigger } from '../../lib/shinpuru-ts/src';
import { Trans, useTranslation } from 'react-i18next';
import { getActionOptions, getRoleOptions, getTriggerOptions } from './shared';
import { useMemo, useReducer } from 'react';

import { ActionButton } from '../ActionButton';
import { Button } from '../Button';
import { Input } from '../Input';
import { Select } from '../Select';
import styled from 'styled-components';

type Props = {
  guild: Guild;
  onApply: (r: KarmaRule) => Promise<any> | undefined;
};

const ruleReducer = (
  state: KarmaRule,
  [type, payload]:
    | ['set', KarmaRule]
    | ['set_trigger', KarmaRuleTrigger]
    | ['set_vaule', number]
    | ['set_action', KarmaRuleAction]
    | ['set_argument', string],
) => {
  switch (type) {
    case 'set':
      return payload;
    case 'set_trigger':
      return { ...state, trigger: payload };
    case 'set_vaule':
      return { ...state, value: payload };
    case 'set_action':
      return { ...state, action: payload, argument: '' };
    case 'set_argument':
      return { ...state, argument: payload };
    default:
      return state;
  }
};

const StyledSelect = styled(Select)``;

const RuleContainer = styled.div`
  ${StyledSelect}, ${Input} {
    display: inline-block;
    width: 10em;
    margin: 0.4em 0.2em;
  }

  ${Button} {
    margin: 1em 0;
    width: 100%;
  }
`;

export const KarmaRuleInput: React.FC<Props> = ({ guild, onApply }) => {
  const { t } = useTranslation('components');
  const [rule, dispatchRule] = useReducer(ruleReducer, {} as KarmaRule);

  const triggerOptions = useMemo(() => getTriggerOptions(t), [t]);
  const actionOptions = useMemo(() => getActionOptions(t), [t]);
  const roleOptions = useMemo(() => getRoleOptions(t, guild), [t, guild]);

  const triggerSelect = (
    <StyledSelect
      options={triggerOptions}
      value={triggerOptions.find((o) => o.value === rule.trigger)}
      onElementSelect={(e) => dispatchRule(['set_trigger', e.value as KarmaRuleTrigger])}
    />
  );
  const valueInput = (
    <Input
      type="number"
      value={rule.value}
      onInput={(e) => dispatchRule(['set_vaule', parseInt(e.currentTarget.value)])}
    />
  );
  const actionSelect = (
    <StyledSelect
      options={actionOptions}
      value={actionOptions.find((o) => o.value === rule.action)}
      onElementSelect={(e) => dispatchRule(['set_action', e.value as KarmaRuleAction])}
    />
  );
  const argumentInput = (() => {
    switch (rule.action) {
      case 'SEND_MESSAGE':
        return (
          <Input
            value={rule.argument}
            onInput={(e) => dispatchRule(['set_argument', e.currentTarget.value])}
          />
        );
      case 'TOGGLE_ROLE':
        return (
          <StyledSelect
            options={roleOptions}
            value={roleOptions.find((r) => r.id === rule.argument)}
            onElementSelect={(e) => dispatchRule(['set_argument', e.value as string])}
          />
        );
      default:
        return <></>;
    }
  })();

  return (
    <RuleContainer>
      <Trans
        ns="components"
        i18nKey="karmarule.text"
        components={{
          '1': triggerSelect,
          '2': valueInput,
          '3': actionSelect,
          '4': argumentInput,
        }}
      />
      <ActionButton onClick={() => onApply(rule)}>{t('karmarule.apply')}</ActionButton>
    </RuleContainer>
  );
};
