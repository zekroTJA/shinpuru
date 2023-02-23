import { KarmaType, getKarmaType } from '../../../util/karma';
import React, { useEffect, useState } from 'react';
import { Trans, useTranslation } from 'react-i18next';
import styled, { css } from 'styled-components';
import { useNavigate, useParams } from 'react-router';

import { ReactComponent as ArrowIcon } from '../../../assets/arrow.svg';
import { ReactComponent as CrownIcon } from '../../../assets/crown.svg';
import { EmptyPlaceholder } from '../../../components/EmptyPlaceholder';
import { Flex } from '../../../components/Flex';
import { GuildScoreboardEntry } from '../../../lib/shinpuru-ts/src';
import { Link } from 'react-router-dom';
import { Loader } from '../../../components/Loader';
import { MaxWidthContainer } from '../../../components/MaxWidthContainer';
import { range } from '../../../util/utils';
import { uid } from 'react-uid';
import { useApi } from '../../../hooks/useApi';
import { useSelfMember } from '../../../hooks/useSelfMember';

type Props = {};

const LIMIT = 50;

const rankIcon = (v: number, i: number) => {
  if (i === 0 && v > 0) return <CrownIcon color="#db9d18" />;
  if (i === 1 && v > 0) return <CrownIcon color="#cecece" />;
  if (i === 2 && v > 0) return <CrownIcon color="#db4f18" />;
  return <>{i + 1}</>;
};

const KarmaEntry = styled.div<{ type: KarmaType; self: boolean }>`
  display: flex;
  gap: 1em;
  align-items: center;
  background-color: ${(p) => p.theme.background2};
  padding: 0.5em 1.5em 0.5em 1em;
  border-radius: 12px;
  cursor: pointer;

  ${(p) =>
    p.self &&
    css`
      outline: solid 2px ${p.theme.accent};
    `};

  > div {
    width: 3ch;
    display: flex;
    align-items: center;
    justify-content: center;

    > svg {
      width: 1.3em;
      height: 1.3em;
    }
  }

  > img {
    width: 3em;
    height: 3em;
    border-radius: 100%;
  }

  > span:last-child {
    margin-left: auto;
    color: ${(p) => {
      switch (p.type) {
        case KarmaType.VERY_HIGH:
          return p.theme.blurple;
        case KarmaType.HIGH:
          return p.theme.green;
        case KarmaType.ZERO:
          return p.theme.yellow;
        case KarmaType.LOW:
          return p.theme.orange;
        default:
          return p.theme.red;
      }
    }};
  }

  transition: transform 0.2s ease;
  &:hover {
    transform: scale(1.01);
  }
`;

const SelfKarmaEntry = styled(KarmaEntry)`
  box-shadow: 0 0.3em 1.5em rgba(0 0 0 / 30%);
  margin-bottom: 0.5em;
`;

const KarmaRoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.guildkarma');
  const { guildid } = useParams();
  const fetch = useApi();
  const nav = useNavigate();
  const selfMember = useSelfMember(guildid);

  const [entries, setEntries] = useState<GuildScoreboardEntry[]>();

  useEffect(() => {
    if (!guildid) return;
    fetch((c) => c.guilds.scoreboard(guildid, LIMIT))
      .then((r) => setEntries(r.data))
      .catch();
  }, [guildid]);

  const containsSelf = !!entries?.find((e) => e.member.user.id === selfMember?.user.id);

  return (
    <MaxWidthContainer>
      <Flex direction="column" gap="1em">
        {(entries && selfMember && (
          <>
            {(entries.length > 0 && (
              <>
                {!containsSelf && (
                  <SelfKarmaEntry
                    key={uid(0)}
                    type={getKarmaType(selfMember.karma)}
                    self={true}
                    onClick={() => nav(`/db/guilds/${guildid}/members/${selfMember.user.id}`)}>
                    <div>
                      <ArrowIcon />
                    </div>
                    <img src={selfMember.avatar_url} alt="" />
                    <span>
                      {selfMember.user.username}#{selfMember.user.discriminator}
                    </span>
                    <span>{selfMember.karma}</span>
                  </SelfKarmaEntry>
                )}
                {entries.map((e, i) => (
                  <KarmaEntry
                    key={uid(e)}
                    type={getKarmaType(e.value)}
                    self={e.member.user.id === selfMember?.user.id}
                    onClick={() => nav(`/db/guilds/${guildid}/members/${e.member.user.id}`)}>
                    <div>{rankIcon(e.value, i)}</div>
                    <img src={e.member.avatar_url} alt="" />
                    <span>
                      {e.member.user.username}#{e.member.user.discriminator}
                    </span>
                    <span>{e.value}</span>
                  </KarmaEntry>
                ))}
              </>
            )) || (
              <EmptyPlaceholder icon={<CrownIcon />}>
                <Trans
                  ns="routes.guildkarma"
                  i18nKey="empty"
                  components={{
                    1: <Link to={`/db/guilds/${guildid}/settings/karma`}>_</Link>,
                    br: <br />,
                  }}
                />
              </EmptyPlaceholder>
            )}
          </>
        )) || (
          <>
            {range(10).map((i) => (
              <Loader key={uid(i)} height="4em" />
            ))}
          </>
        )}
      </Flex>
    </MaxWidthContainer>
  );
};

export default KarmaRoute;
