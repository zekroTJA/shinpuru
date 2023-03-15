import { ModalCreateReport, ReportActionType } from '../../../components/Modals/ModalCreateReport';
import styled, { useTheme } from 'styled-components';
import { useEffect, useState } from 'react';

import { ReactComponent as BotIcon } from '../../../assets/bot.svg';
import { Button } from '../../../components/Button';
import { Container } from '../../../components/Container';
import { DiscordImage } from '../../../components/DiscordImage';
import { Embed } from '../../../components/Embed';
import { Flex } from '../../../components/Flex';
import { Heading } from '../../../components/Heading';
import { Hint } from '../../../components/Hint';
import { ReactComponent as InfoIcon } from '../../../assets/info.svg';
import { KarmaTile } from '../../../components/KarmaTile';
import { Loader } from '../../../components/Loader';
import { PermsSimpleList } from '../../../components/Permissions';
import { Report } from '../../../lib/shinpuru-ts/src';
import { ReportsList } from '../../../components/Report';
import { RoleList } from '../../../components/RoleList';
import { SinceDate } from '../../../components/SinceDate';
import { formatDate } from '../../../util/date';
import { memberName } from '../../../util/users';
import { useGuild } from '../../../hooks/useGuild';
import { useMember } from '../../../hooks/useMember';
import { useParams } from 'react-router';
import { usePerms } from '../../../hooks/usePerms';
import { useSelfUser } from '../../../hooks/useSelfUser';
import { useTranslation } from 'react-i18next';

type Props = {};

const MemberContainer = styled.div``;

const Header = styled.div`
  display: flex;
`;

const HeaderName = styled(Flex)`
  margin-bottom: 1em;

  gap: 1em;
  align-items: center;

  > h1 {
    margin: 0;
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

  const selfUser = useSelfUser();
  const guild = useGuild(guildid);
  const [member, memberReq] = useMember(guildid, memberid);
  const [perms, setPerms] = useState<string[]>();
  const [reports, setReports] = useState<Report[]>();
  const { isAllowed } = usePerms(guild?.id);

  const [reportModalType, setReportModalType] = useState<ReportActionType>('report');
  const [showReportModal, setShowReportModal] = useState(false);

  const _report = (type: ReportActionType) => {
    setReportModalType(type);
    setShowReportModal(true);
  };

  const _onReportSubmitted = (rep: Report) => {
    setReports([rep, ...(reports ?? [])]);
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
        type={reportModalType}
        member={member}
        onClose={() => setShowReportModal(false)}
        onSubmitted={_onReportSubmitted}
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
          <ReportsList
            reports={reports}
            revokeAllowed={isAllowed('sp.guild.mod.report.revoke')}
            emptyPlaceholder={<i>{t('reports.noreports')}</i>}
            onReportsUpdated={setReports}
          />
        </Section>
      </DetailsContainer>
    </MemberContainer>
  ) : (
    <Loaders />
  );
};

export default MemberRoute;
