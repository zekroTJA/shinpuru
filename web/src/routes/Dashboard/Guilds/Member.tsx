import { useTranslation } from 'react-i18next';
import { useParams } from 'react-router';
import styled, { useTheme } from 'styled-components';
import { useMember } from '../../../hooks/useMember';
import { DiscordImage } from '../../../components/DiscordImage';
import { Loader } from '../../../components/Loader';
import { memberName } from '../../../util/users';
import { RoleList } from '../../../components/RoleList';
import { useGuild } from '../../../hooks/useGuild';
import { Flex } from '../../../components/Flex';
import { Embed } from '../../../components/Embed';
import { Heading } from '../../../components/Heading';
import { formatDate } from '../../../util/date';
import { SinceDate } from '../../../components/SinceDate';
import { Hint } from '../../../components/Hint';
import { useSelfUser } from '../../../hooks/useSelfUser';
import { ReactComponent as BotIcon } from '../../../assets/bot.svg';
import { ReactComponent as InfoIcon } from '../../../assets/info.svg';
import { Container } from '../../../components/Container';
import { KarmaTile } from '../../../components/KarmaTile';
import { PermsSimpleList } from '../../../components/Permissions';
import { useEffect, useState } from 'react';
import { Report } from '../../../lib/shinpuru-ts/src';
import { MemberReportsList } from '../../../components/Report';
import { usePerms } from '../../../hooks/usePerms';

interface Props {}

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

const ReportsContainer = styled.div``;

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

export const MemebrRoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.member');
  const { guildid, memberid } = useParams();
  const theme = useTheme();

  const selfUser = useSelfUser();
  const guild = useGuild(guildid);
  const [member, memberReq] = useMember(guildid, memberid);
  const [perms, setPerms] = useState<string[]>();
  const [reports, setReports] = useState<Report[]>();
  const { isAllowed } = usePerms(guild?.id);

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
      {member.user.bot && (
        <MarginHint icon={<BotIcon />}>{t('isbot')}</MarginHint>
      )}
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
          {guild && (
            <RoleList guildroles={guild.roles!} roleids={member.roles} />
          )}
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
        <Section hide={member.user.bot} fw>
          <Heading>{t('reports.heading')}</Heading>
          <MemberReportsList
            reports={reports}
            revokeAllowed={isAllowed('sp.guild.mod.report.revoke')}
          />
        </Section>
      </DetailsContainer>
    </MemberContainer>
  ) : (
    <Loaders />
  );
};
