import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useParams } from 'react-router';
import styled from 'styled-components';
import { Heading } from '../../../components/Heading';
import { ReportsList } from '../../../components/Report';
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
`;

const GuildModlogRoute: React.FC<Props> = () => {
  const { t } = useTranslation('routes.guildmodlog');
  const fetch = useApi();
  const { guildid } = useParams();
  const { isAllowed } = usePerms(guildid);

  const [reports, setReports] = useState<Report[]>();
  const [unbanRequests, setUnbanRequests] = useState<UnbanRequest[]>();

  useEffect(() => {
    if (guildid) {
      fetch((c) => c.guilds.reports(guildid, 100))
        .then((r) => setReports(r.data))
        .catch();

      fetch((c) => c.guilds.unbanrequests(guildid))
        .then((r) => setUnbanRequests(r.data))
        .catch();
    }
  }, [guildid]);

  return (
    <Container>
      <section>
        <Heading>{t('reports')}</Heading>
        <StyledReprtList
          reports={reports}
          onReportsUpdated={setReports}
          revokeAllowed={isAllowed('sp.guild.mod.report.revoke')}
          emptyPlaceholder={<i>{t('noreports')}</i>}
        />
      </section>
      <section>
        <Heading>{t('unbanrequests')}</Heading>
        {/* {unbanRequests?.map(r => r.)} */}
      </section>
    </Container>
  );
};

export default GuildModlogRoute;
