import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useParams } from 'react-router';
import styled from 'styled-components';
import { Flex } from '../../../components/Flex';
import { Heading } from '../../../components/Heading';
import { Loader } from '../../../components/Loader';
import {
  ModalProcessUnbanRequest,
  UnbanRequestWrapper,
} from '../../../components/Modals/ModalProcessUnbanRequest';
import { ReportsList } from '../../../components/Report';
import { UnbanRequestTile } from '../../../components/UnbanRequestTile';
import { useApi } from '../../../hooks/useApi';
import { usePerms } from '../../../hooks/usePerms';
import { Report, UnbanRequest } from '../../../lib/shinpuru-ts/src';

type Props = {};

const StyledReprtList = styled(ReportsList)``;

const Container = styled.div`
  width: 100%;
  display: flex;
  gap: 1em;

  > section {
    width: 100%;
  }

  @media (orientation: portrait) {
    flex-direction: column;
  }
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

  useEffect(() => {
    if (guildid) {
      fetch((c) => c.guilds.reports(guildid, 100))
        .then((r) => setReports(r.data))
        .catch();
    }
  }, [guildid]);

  useEffect(() => {
    if (!!guildid && !!allowedPerms && isAllowed('sp.guild.mod.unbanrequests')) {
      fetch((c) => c.guilds.unbanrequests(guildid))
        .then((r) => setUnbanRequests(r.data))
        .catch();
    }
  }, [allowedPerms, guildid]);

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
    <Container>
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
    </Container>
  );
};

export default GuildModlogRoute;
