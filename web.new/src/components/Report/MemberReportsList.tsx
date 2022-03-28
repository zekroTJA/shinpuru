import { t } from 'i18next';
import { useTranslation } from 'react-i18next';
import styled from 'styled-components';
import { Report } from '../../lib/shinpuru-ts/src';
import { Loader } from '../Loader';
import { ReportTile } from './ReportTile';

type Props = {
  reports?: Report[];
  revokeAllowed?: boolean;
  onRevoke?: (rep: Report) => void;
};

const ReportsContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1em;
`;

export const MemberReportsList: React.FC<Props> = ({
  reports,
  revokeAllowed,
  onRevoke = () => {},
}) => {
  const { t } = useTranslation('components');
  return !!reports ? (
    <ReportsContainer>
      {(reports.length === 0 && <i>{t('memberreportlist.noreports')}</i>) ||
        reports.map((r) => (
          <ReportTile
            report={r}
            revokeAllowed={revokeAllowed}
            onRevoke={() => onRevoke(r)}
            key={r.id}
          />
        ))}
    </ReportsContainer>
  ) : (
    <Loader width="100%" height="4em" />
  );
};
