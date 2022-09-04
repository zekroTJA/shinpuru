import { Guild, KarmaRuleAction, KarmaRuleTrigger } from '../../lib/shinpuru-ts/src';
import { Element } from '../Select';

type TranslateFunc = (v: string) => string;

export const getTriggerOptions = (t: TranslateFunc) =>
  [
    {
      id: 'below',
      value: KarmaRuleTrigger.BELOW,
      display: t('karmarule.trigger.dropsbelow'),
    },
    {
      id: 'above',
      value: KarmaRuleTrigger.ABOVE,
      display: t('karmarule.trigger.risesabove'),
    },
  ] as Element<KarmaRuleTrigger>[];

export const getActionOptions = (t: TranslateFunc) =>
  (['TOGGLE_ROLE', 'KICK', 'BAN', 'SEND_MESSAGE'] as KarmaRuleAction[]).map((a) => ({
    id: a,
    display: t(`karmarule.action.${a.toLowerCase().replaceAll('_', '')}`),
    value: a,
  })) as Element<KarmaRuleAction>[];

export const getRoleOptions = (t: TranslateFunc, g: Guild) =>
  (g?.roles ?? []).map((r) => ({
    id: r.id,
    display: r.name,
    value: r.id,
  })) as Element<string>[];
