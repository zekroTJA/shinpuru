import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useParams } from 'react-router';
import styled, { useTheme } from 'styled-components';
import { ReactComponent as BotIcon } from '../../../assets/bot.svg';
import { ReactComponent as InfoIcon } from '../../../assets/info.svg';
import { Button } from '../../../components/Button';
import { Container } from '../../../components/Container';
import { DiscordImage } from '../../../components/DiscordImage';
import { Embed } from '../../../components/Embed';
import { Flex } from '../../../components/Flex';
import { Heading } from '../../../components/Heading';
import { Hint } from '../../../components/Hint';
import { KarmaTile } from '../../../components/KarmaTile';
import { Loader } from '../../../components/Loader';
import { ModalCreateReport, ReportActionType } from '../../../components/Modals/ModalCreateReport';
import { ModalRevokeReport } from '../../../components/Modals/ModalRevokeReport';
import { NotificationType } from '../../../components/Notifications';
import { PermsSimpleList } from '../../../components/Permissions';
import { MemberReportsList } from '../../../components/Report';
import { RoleList } from '../../../components/RoleList';
import { SinceDate } from '../../../components/SinceDate';
import { useApi } from '../../../hooks/useApi';
import { useGuild } from '../../../hooks/useGuild';
import { useMember } from '../../../hooks/useMember';
import { useNotifications } from '../../../hooks/useNotifications';
import { usePerms } from '../../../hooks/usePerms';
import { useSelfUser } from '../../../hooks/useSelfUser';
import { Report } from '../../../lib/shinpuru-ts/src';
import { formatDate } from '../../../util/date';
import { memberName } from '../../../util/users';

type Props = {};

const MemberContainer = styled.div``;

const Header = styled.div`
  display: flex;
`;

const HeaderName = styled(Flex)`
  margin-bottom: 1em;

  > * {
    margin: 0 0.5em 0 0;
  }
`;

const StyledDiscordImage = styled(DiscordImage)`
  width: 7em;
  height: 7em;
  margin-right: 1em;
`;

const MarginHint = styled(Hint)`
  margin-bottom: 1.5em;
`;

const Section = styled(Container)<{ hide?: boolean; fw?: boolean }>`
  ${(p) => (p.hide ? 'display: none;' : '')}
  width: ${(p) => (p.fw ? '100%' : 'fit-content')};
`;

const Table = styled.table`
  text-align: left;

  th {
    vertical-align: baseline;
    padding-right: 1.5em;
  }

  td {
    padding-bottom: 1em;
  }
`;

const Secondary = styled.span`
  opacity: 0.6;
`;

const DetailsContainer = styled(Flex)`
  margin-top: 1.5em;
  gap: 1em;
  flex-wrap: wrap;
  align-items: stretch;
`;

const Loaders = () => (
  <>
    <Flex>
      <Loader width="7em" height="7em" margin="0 1em 0 0" />
      <div style={{ width: '100%' }}>
        <Loader width="100%" height="2em" />
        <Loader width="100%" height="4em" margin="1em 0 0 0" />
      </div>
    </Flex>
  </>
);

const MemberRoute: React.FC<Props> = () => {
  const { t } = useTranslation('routes.member');
  const { guildid, memberid } = useParams();
  const theme = useTheme();
  const fetch = useApi();
  const { pushNotification } = useNotifications();

  const selfUser = useSelfUser();
  const guild = useGuild(guildid);
  const [member, memberReq] = useMember(guildid, memberid);
  const [perms, setPerms] = useState<string[]>();
  const [reports, setReports] = useState<Report[]>();
  const { isAllowed } = usePerms(guild?.id);

  const [reportModalType, setReportModalType] = useState<ReportActionType>('report');
  const [showReportModal, setShowReportModal] = useState(false);
  const [revokeReport, setRevokeReport] = useState<Report>();

  const _report = (type: ReportActionType) => {
    setReportModalType(type);
    setShowReportModal(true);
  };

  const _revokeReport = (rep: Report) => {
    setRevokeReport(rep);
  };

  const _revokeReportConfirm = (rep: Report, reason: string) => {
    if (!revokeReport || !reports) return;

    fetch((c) =>
      c.reports.revoke(revokeReport.id, {
        reason,
      }),
    )
      .then(() => {
        const i = reports.findIndex((r) => r.id === rep.id);
        if (i !== -1) reports.splice(i, 1);
        console.log(rep, i);
        setReports([...reports]);
        pushNotification({
          heading: t('notifications.reportrevoked.heading'),
          message: t('notifications.reportrevoked.message'),
          type: NotificationType.SUCCESS,
        });
      })
      .catch();
  };

  useEffect(() => {
    memberReq((c) => c.permissions())
      .then((res) => setPerms(res.permissions))
      .catch();
  }, []);

  useEffect(() => {
    if (member && !member.user.bot) {
      memberReq((c) => c.reports())
        .then((res) => setReports(res.data))
        .catch();
    }
  }, [member]);

  return member ? (
    <MemberContainer>
      <ModalCreateReport
        show={showReportModal}
        onClose={() => setShowReportModal(false)}
        type={reportModalType}
        member={member!}
      />
      <ModalRevokeReport
        report={revokeReport}
        onConfirm={_revokeReportConfirm}
        onClose={() => setRevokeReport(undefined)}
      />

      {member.user.bot && <MarginHint icon={<BotIcon />}>{t('isbot')}</MarginHint>}
      {member.user.id === selfUser?.id && (
        <MarginHint icon={<InfoIcon />} color={theme.green}>
          {t('isyou')}
        </MarginHint>
      )}

      <Header>
        <StyledDiscordImage src={member.avatar_url} />
        <div>
          <HeaderName>
            <h1>{memberName(member)}</h1>
            <small>
              {member.user.username}#{member.user.discriminator}
            </small>
            <Embed>{member.user.id}</Embed>
          </HeaderName>
          {guild && <RoleList guildroles={guild.roles!} roleids={member.roles} />}
        </div>
      </Header>

      <DetailsContainer>
        <Section>
          <Heading>{t('info.heading')}</Heading>
          <Table>
            <tbody>
              <tr>
                <th>{t('info.guild-joined')}</th>
                <td>
                  {formatDate(member.joined_at)}
                  <br />
                  <Secondary>
                    <SinceDate date={member.joined_at} />
                  </Secondary>
                </td>
              </tr>
              <tr>
                <th>{t('info.account-created')}</th>
                <td>
                  {formatDate(member.created_at)}
                  <br />
                  <Secondary>
                    <SinceDate date={member.created_at} />
                  </Secondary>
                </td>
              </tr>
            </tbody>
          </Table>
        </Section>
        <Section hide={member.user.bot}>
          <Heading>{t('karma.heading')}</Heading>
          <Flex gap="1em">
            <KarmaTile heading={t('karma.guild')} points={member.karma} />
            <KarmaTile heading={t('karma.total')} points={member.karma_total} />
          </Flex>
        </Section>
        <Section hide={member.user.bot}>
          <Heading>{t('permissions.heading')}</Heading>
          <PermsSimpleList perms={perms} />
        </Section>
        <Section hide={member.user.bot || !isAllowed('sp.guild.mod.')} fw>
          <Heading>{t('moderation.heading')}</Heading>
          <Flex gap="1em">
            {isAllowed('sp.guild.mod.report') && (
              <Button onClick={() => _report('report')} variant="blue">
                {t('moderation.report')}
              </Button>
            )}
            {isAllowed('sp.guild.mod.kick') && (
              <Button onClick={() => _report('kick')} variant="orange">
                {t('moderation.kick')}
              </Button>
            )}
            {isAllowed('sp.guild.mod.ban') && (
              <Button onClick={() => _report('ban')} variant="red">
                {t('moderation.ban')}
              </Button>
            )}
            {isAllowed('sp.guild.mod.mute') && (
              <Button onClick={() => _report('mute')} variant="pink">
                {t('moderation.mute')}
              </Button>
            )}
          </Flex>
        </Section>
        <Section hide={member.user.bot} fw>
          <Heading>{t('reports.heading')}</Heading>
          <MemberReportsList
            reports={reports}
            revokeAllowed={isAllowed('sp.guild.mod.report.revoke')}
            onRevoke={_revokeReport}
          />
        </Section>
      </DetailsContainer>
    </MemberContainer>
  ) : (
    <Loaders />
  );
};

export default MemberRoute;
