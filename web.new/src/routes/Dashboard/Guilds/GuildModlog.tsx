import {
  ModalProcessUnbanRequest,
  UnbanRequestWrapper,
} from '../../../components/Modals/ModalProcessUnbanRequest';
import { Report, UnbanRequest } from '../../../lib/shinpuru-ts/src';
import { useEffect, useRef, useState } from 'react';

import { Button } from '../../../components/Button';
import { Flex } from '../../../components/Flex';
import { Heading } from '../../../components/Heading';
import { Loader } from '../../../components/Loader';
import { ReportsList } from '../../../components/Report';
import { SplitContainer } from '../../../components/SplitContainer';
import { UnbanRequestTile } from '../../../components/UnbanRequestTile';
import styled from 'styled-components';
import { useApi } from '../../../hooks/useApi';
import { useParams } from 'react-router';
import { usePerms } from '../../../hooks/usePerms';
import { useTranslation } from 'react-i18next';

const REPORTS_LIMIT = 10;

type Props = {};

const StyledReprtList = styled(ReportsList)``;

const LoadButtonContainer = styled.div`
  display: flex;
  justify-content: center;
  margin-top: 1em;
`;

const GuildModlogRoute: React.FC<Props> = () => {
  const { t } = useTranslation('routes.guildmodlog');
  const fetch = useApi();
  const { guildid } = useParams();
  const { allowedPerms, isAllowed } = usePerms(guildid);

  const [reports, setReports] = useState<Report[]>();
  const [unbanRequests, setUnbanRequests] = useState<UnbanRequest[]>();
  const [showUnabnModal, setShowUnbanModal] = useState(false);
  const [selectedUnban, setSelectedUnban] = useState<UnbanRequestWrapper>();

  const reportsOffsetRef = useRef(0);
  const reportsTotalCountRef = useRef(0);
  const unbanOffsetRef = useRef(0);
  const unbanTotalCountRef = useRef(0);

  useEffect(() => {
    if (guildid) {
      fetch((c) => c.guilds.reportsCount(guildid))
        .then((r) => (reportsTotalCountRef.current = r.count))
        .catch();
      fetch((c) => c.guilds.reports(guildid, REPORTS_LIMIT))
        .then((r) => setReports(r.data))
        .catch();
    }
  }, [guildid]);

  useEffect(() => {
    if (!!guildid && !!allowedPerms && isAllowed('sp.guild.mod.unbanrequests')) {
      fetch((c) => c.guilds.unbanrequestsCount(guildid))
        .then((r) => (unbanTotalCountRef.current = r.count))
        .catch();
      fetch((c) => c.guilds.unbanrequests(guildid))
        .then((r) => setUnbanRequests(r.data))
        .catch();
    }
  }, [allowedPerms, guildid]);

  const _fetchMoreReports = () => {
    if (guildid && !!reports) {
      reportsOffsetRef.current += REPORTS_LIMIT;
      fetch((c) => c.guilds.reports(guildid, REPORTS_LIMIT, reportsOffsetRef.current))
        .then((r) => setReports([...reports, ...r.data]))
        .catch();
    }
  };

  const _fetchMoreUnban = () => {
    if (guildid && !!unbanRequests) {
      unbanOffsetRef.current += REPORTS_LIMIT;
      fetch((c) => c.guilds.unbanrequests(guildid, REPORTS_LIMIT, unbanOffsetRef.current))
        .then((r) => setUnbanRequests([...unbanRequests, ...r.data]))
        .catch();
    }
  };

  const _onUnbanReqeustProcess = (request: UnbanRequest, isAccept: boolean) => {
    setSelectedUnban({ ...request, isAccept });
    setShowUnbanModal(true);
  };

  const _onUnbanRequestProcessed = (request: UnbanRequest) => {
    if (!unbanRequests) return;
    const i = unbanRequests.findIndex((r) => request.id === r.id);
    if (i !== -1) unbanRequests[i] = request;
    setUnbanRequests([...unbanRequests]);
  };

  return (
    <SplitContainer>
      <ModalProcessUnbanRequest
        show={showUnabnModal}
        request={selectedUnban}
        onClose={() => setShowUnbanModal(false)}
        onProcessed={_onUnbanRequestProcessed}
      />

      <section>
        <Heading>{t('reports')}</Heading>
        <StyledReprtList
          reports={reports}
          onReportsUpdated={setReports}
          revokeAllowed={isAllowed('sp.guild.mod.report.revoke')}
          emptyPlaceholder={<i>{t('noreports')}</i>}
        />
        {reports && reports.length < reportsTotalCountRef.current && (
          <LoadButtonContainer>
            <Button onClick={_fetchMoreReports}>{t('loadmore')}</Button>
          </LoadButtonContainer>
        )}
      </section>
      {isAllowed('sp.guild.mod.unbanrequests') && (
        <section>
          <Heading>{t('unbanrequests')}</Heading>
          <Flex direction="column" gap="1em">
            {(unbanRequests && (
              <>
                {(unbanRequests.length === 0 && <i>{t('nounbanrequests')}</i>) ||
                  unbanRequests.map((r) => (
                    <UnbanRequestTile
                      key={r.id}
                      request={r}
                      showControls
                      onProcess={(a) => _onUnbanReqeustProcess(r, a)}
                    />
                  ))}
                {unbanRequests && unbanRequests.length < unbanTotalCountRef.current && (
                  <LoadButtonContainer>
                    <Button onClick={_fetchMoreUnban}>{t('loadmore')}</Button>
                  </LoadButtonContainer>
                )}
              </>
            )) || (
              <>
                <Loader height="10em" />
                <Loader height="10em" />
                <Loader height="10em" />
              </>
            )}
          </Flex>
        </section>
      )}
    </SplitContainer>
  );
};

export default GuildModlogRoute;
