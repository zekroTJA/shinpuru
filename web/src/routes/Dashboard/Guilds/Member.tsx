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
  margin-bottom: 1em;
`;

const Section = styled.section`
  margin-top: 2em;
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
    </MemberContainer>
  ) : (
    <Loaders />
  );
};
