import { GuildStarboardEntry, StarboardSortOrder } from '../../../lib/shinpuru-ts/src';
import React, { useEffect, useRef, useState } from 'react';

import { Button } from '../../../components/Button';
import { Flex } from '../../../components/Flex';
import { MaxWidthContainer } from '../../../components/MaxWidthContainer';
import { ReactComponent as RefreshIcon } from '../../../assets/refresh.svg';
import { StarboardEntry } from '../../../components/StarboardEntry';
import styled from 'styled-components';
import { useApi } from '../../../hooks/useApi';
import { useNotifications } from '../../../hooks/useNotifications';
import { useParams } from 'react-router';
import { useTranslation } from 'react-i18next';

type Props = {};

const PAGE_SIZE = 2;

const Header = styled.div`
  display: flex;
  gap: 1em;
  margin-bottom: 1em;
  align-items: center;
  background-color: ${(p) => p.theme.background3};
  border-radius: 12px;
  padding: 0.5em 0.8em;

  > h1 {
    margin: 0 auto 0 0;
  }

  > ${Button} {
    padding: 0.5em;
    border-radius: 8px;
  }
`;

const StarboardRoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.guildstarboard');
  const { pushNotification } = useNotifications();
  const { guildid } = useParams();
  const fetch = useApi();

  const [entries, setEntries] = useState<GuildStarboardEntry[]>();
  const [sortBy, setSortBy] = useState<StarboardSortOrder>('latest');

  const offsetRef = useRef(0);
  const totalCountRef = useRef(0);

  const _refresh = () => {
    if (!guildid) return;
    fetch((c) => c.guilds.starboardCount(guildid))
      .then((r) => (totalCountRef.current = r.count))
      .catch();
    fetch((c) => c.guilds.starboard(guildid, sortBy, PAGE_SIZE, 0))
      .then((r) => setEntries(r.data))
      .catch();
  };

  const _loadMore = () => {
    if (!guildid || entries === undefined) return;
    offsetRef.current += PAGE_SIZE;
    fetch((c) => c.guilds.starboard(guildid, sortBy, PAGE_SIZE, offsetRef.current))
      .then((r) => setEntries([...entries, ...r.data]))
      .catch();
  };

  useEffect(() => _refresh(), [guildid, sortBy]);

  return (
    <MaxWidthContainer>
      <Header>
        <h1>{t('heading')}</h1>
        <Button onClick={() => setSortBy(sortBy === 'latest' ? 'top' : 'latest')}>
          {t(sortBy === 'latest' ? 'sort.latest' : 'sort.top')}
        </Button>
        <Button onClick={_refresh}>
          <RefreshIcon />
        </Button>
      </Header>
      {entries && (
        <Flex direction="column" gap="1em">
          {entries.map((e) => (
            <StarboardEntry entry={e} />
          ))}
          {entries.length < totalCountRef.current && (
            <Button onClick={_loadMore}>{t('loadmore')}</Button>
          )}
        </Flex>
      )}
    </MaxWidthContainer>
  );
};

export default StarboardRoute;
