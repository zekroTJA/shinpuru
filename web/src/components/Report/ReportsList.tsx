import { Loader } from '../Loader';
import { ModalRevokeReport } from '../Modals/ModalRevokeReport';
import { Report } from '../../lib/shinpuru-ts/src';
import { ReportTile } from './ReportTile';
import styled from 'styled-components';
import { useApi } from '../../hooks/useApi';
import { useNotifications } from '../../hooks/useNotifications';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';

type Props = React.HTMLAttributes<HTMLDivElement> & {
  reports?: Report[];
  onReportsUpdated?: (reports: Report[]) => void;
  revokeAllowed?: boolean;
  emptyPlaceholder?: JSX.Element;
};

const ReportsContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1em;
`;

export const ReportsList: React.FC<Props> = ({
  reports,
  revokeAllowed,
  onReportsUpdated = () => {},
  emptyPlaceholder,
  ...props
}) => {
  const { t } = useTranslation('components', { keyPrefix: 'reportlist' });
  const fetch = useApi();
  const { pushNotification } = useNotifications();

  const [revokeReport, setRevokeReport] = useState<Report>();

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
        onReportsUpdated(reports);
        pushNotification({
          heading: t<string>('notifications.reportrevoked.heading'),
          message: t<string>('notifications.reportrevoked.message'),
          type: 'SUCCESS',
        });
      })
      .catch();
  };

  return (
    <>
      <ModalRevokeReport
        report={revokeReport}
        onConfirm={_revokeReportConfirm}
        onClose={() => setRevokeReport(undefined)}
      />

      <ReportsContainer {...props}>
        {(reports &&
          ((reports.length === 0 && emptyPlaceholder) ||
            reports.map((r) => (
              <ReportTile
                report={r}
                revokeAllowed={revokeAllowed}
                onRevoke={() => _revokeReport(r)}
                key={r.id}
              />
            )))) || <Loader width="100%" height="4em" />}
      </ReportsContainer>
    </>
  );
};
