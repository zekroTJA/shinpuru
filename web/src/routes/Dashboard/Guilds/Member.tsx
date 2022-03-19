import { useTranslation } from 'react-i18next';
import { useParams } from 'react-router';
import { format } from 'date-fns';
import styled from 'styled-components';
import { useMember } from '../../../hooks/useMember';

interface Props {}

const MemberContainer = styled.div``;

export const MemebrRoute: React.FC<Props> = ({}) => {
  // const { t } = useTranslation('routes.member');
  const { guildid, memberid } = useParams();
  // const [member, memberReq] = useMember(guildid, memberid);

  return (
    <>
      {guildid}/{memberid}
    </>
  );
};
