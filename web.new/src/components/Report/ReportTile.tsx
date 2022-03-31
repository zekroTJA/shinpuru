import { useTranslation } from 'react-i18next';
import styled from 'styled-components';
import { ReactComponent as HammerIcon } from '../../assets/hammer.svg';
import { ReactComponent as TargetIcon } from '../../assets/target.svg';
import { Report, ReportType } from '../../lib/shinpuru-ts/src';
import { formatDate } from '../../util/date';
import { Container } from '../Container';
import { Embed } from '../Embed';
import { Heading } from '../Heading';
import { LinkButton } from '../LinkButton';
import { UserTileSmall } from '../UserTileSmall';
import { LinearGradient } from '../styleParts';

type Props = React.HTMLAttributes<HTMLDivElement> & {
  report: Report;
  revokeAllowed?: boolean;
  onRevoke?: () => void;
};

const ReportTileContainer = styled(Container)`
  background-color: ${(p) => p.theme.background3};
  position: relative;
  overflow: hidden;
  border-radius: 8px;
`;

const TypeHead = styled.div<{ type: number }>`
  font-size: 0.8rem;
  letter-spacing: 0.2ch;
  text-transform: uppercase;
  position: absolute;
  width: 100%;
  top: 0;
  left: 0;
  text-align: center;
  padding: 0.4em;
  border-radius: 8px;
  color: ${(p) => p.theme.background2};

  ${(p) => {
    switch (p.type) {
      case ReportType.KICK:
      case ReportType.BAN:
        return LinearGradient(p.theme.red);
      case ReportType.MUTE:
        return LinearGradient(p.theme.pink);
      default:
        return LinearGradient(p.theme.orange);
    }
  }};
`;

const Section = styled.section`
  margin-top: 1em;

  > ${Heading} {
    font-size: 0.8rem;
  }

  > img {
    width: 100%;
    height: auto;
    margin-top: 1em;
  }
`;

const ReportUsers = styled.div`
  display: flex;
  justify-content: space-between;
  margin-top: 1.5em;
  gap: 1.5em;
`;

const Footer = styled(Section)`
  display: flex;
  align-items: center;

  font-size: 0.8rem;

  ${Embed} {
    font-size: 0.8em;
  }

  > a {
    text-decoration: underline;
    color: ${(p) => p.theme.accent};
    cursor: pointer;
  }
`;

const Spacer = styled.div`
  height: 1em;
  width: 1px;
  background-color: ${(p) => p.theme.text};
  margin: 0 0.5em;
`;

export const ReportTile: React.FC<Props> = ({
  report,
  revokeAllowed,
  onRevoke = () => {},
  ...props
}) => {
  const { t } = useTranslation('components', { keyPrefix: 'report' });
  return (
    <ReportTileContainer {...props}>
      <TypeHead type={report.type}>{report.type_name}</TypeHead>
      <ReportUsers>
        <UserTileSmall
          fallbackId={report.executor_id}
          user={report.executor}
          icon={<HammerIcon />}
        />
        <UserTileSmall fallbackId={report.victim_id} user={report.victim} icon={<TargetIcon />} />
      </ReportUsers>
      <Section>
        <Heading>{t('reason')}</Heading>
        <span>{report.message}</span>
        {report.attachment_url && <img src={report.attachment_url} alt="Report Attachment" />}
      </Section>
      <hr />
      <Footer>
        <span>
          ID: <Embed>{report.id}</Embed>
        </span>
        <Spacer />
        <span>{formatDate(report.created)}</span>
        {revokeAllowed && (
          <>
            <Spacer />
            <LinkButton onClick={onRevoke}>{t('revoke')}</LinkButton>
          </>
        )}
      </Footer>
    </ReportTileContainer>
  );
};
