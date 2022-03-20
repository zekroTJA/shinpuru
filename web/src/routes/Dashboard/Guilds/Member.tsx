import { useTranslation } from 'react-i18next';
import { useParams } from 'react-router';
import styled from 'styled-components';
import { useMember } from '../../../hooks/useMember';
import { DiscordImage } from '../../../components/DiscordImage';
import { Loader } from '../../../components/Loader';
import { memberName } from '../../../util/users';
import { RoleList } from '../../../components/RoleList';
import { useGuild } from '../../../hooks/useGuild';
import { Flex } from '../../../components/Flex';
import { Embed } from '../../../components/Embed';
import { ReactComponent as BotIcon } from '../../../assets/bot.svg';
import { LinearGradient } from '../../../components/styleParts';
import { Heading } from '../../../components/Heading';
import { format } from 'date-fns';
import { formatDate, formatSince } from '../../../util/date';
import { SinceDate } from '../../../components/SinceDate';

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

const BotContainer = styled(Flex)`
  ${(p) => LinearGradient(p.theme.blurple)};

  border-radius: 12px;
  padding: 0.5em;
  margin-bottom: 1em;

  > svg {
    height: 1.5em;
    width: auto;
    margin-right: 0.5em;
  }
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
  const { t, i18n } = useTranslation('routes.member');
  const { guildid, memberid } = useParams();
  const guild = useGuild(guildid);
  const [member, memberReq] = useMember(guildid, memberid);

  return member ? (
    <MemberContainer>
      {member.user.bot && (
        <BotContainer>
          <BotIcon />
          <span>{t('isbot')}</span>
        </BotContainer>
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
