import { Member, Report, ReportRequest, ReportType } from '../../lib/shinpuru-ts/src';
import { ModalContainer, ModalTextArea } from './modalParts';
import { Trans, useTranslation } from 'react-i18next';
import { useEffect, useState } from 'react';

import { Button } from '../Button';
import { ControlProps } from '../Modal/Modal';
import { DurationPicker } from '../DurationPicker';
import { Filedrop } from '../Filedrop';
import { Heading } from '../Heading';
import { Modal } from '../Modal';
import { TextArea } from '../TextArea';
import { parseToDateString } from '../../util/date';
import { readToBase64 } from '../../util/files';
import styled from 'styled-components';
import { useApi } from '../../hooks/useApi';
import { useNavigate } from 'react-router';
import { useNotifications } from '../../hooks/useNotifications';

const ALLOWED_ATTACHMENT_TYPES = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif'];

export type ReportActionType = 'report' | 'kick' | 'ban' | 'mute';

type Props = ControlProps & {
  type: ReportActionType;
  member: Member;
  onSubmitted?: (report: Report) => void;
};

const ACTION_TEXT = {
  kick: 'routes.member:moderation.kick',
  ban: 'routes.member:moderation.ban',
  mute: 'routes.member:moderation.mute',
  report: 'routes.member:moderation.report',
};

export const ModalCreateReport: React.FC<Props> = ({
  show,
  type,
  member,
  onClose = () => {},
  onSubmitted = () => {},
}) => {
  const { t } = useTranslation();
  const fetch = useApi();
  const nav = useNavigate();
  const { pushNotification } = useNotifications();
  const [reason, setReason] = useState('');
  const [attachment, setAttachment] = useState<File>();
  const [timeout, setTimeout] = useState(203535);

  useEffect(() => {
    if (show) {
      setReason('');
      setAttachment(undefined);
      setTimeout(0);
    }
  }, [show]);

  const _setAttachment = (f: File) => {
    if (!ALLOWED_ATTACHMENT_TYPES.includes(f.type))
      throw t('errors.reports.disallowed-attachment-type');
    if (f.size > 50 * 1024 * 1024) throw t('errors.reports.file-too-big');
    setAttachment(f);
  };

  const _submit = async () => {
    const rep = {
      reason,
    } as ReportRequest;

    if (attachment) {
      try {
        rep.attachment_data = await readToBase64(attachment);
      } catch (e) {
        pushNotification({
          message: t('components:modalcreatereport.errors.attachment-convert-failed'),
          type: 'ERROR',
        });
      }
    }

    if (timeout > 0) {
      rep.timeout = parseToDateString(new Date(Date.now() + timeout * 1000));
    }

    let req;
    let goBack = false;
    switch (type) {
      case 'report':
        rep.type = ReportType.WARN;
        req = fetch((c) => c.guilds.member(member.guild_id, member.user.id).report(rep));
        break;
      case 'kick':
        goBack = true;
        req = fetch((c) => c.guilds.member(member.guild_id, member.user.id).kick(rep));
        break;
      case 'ban':
        goBack = true;
        req = fetch((c) => c.guilds.member(member.guild_id, member.user.id).ban(rep));
        break;
      case 'mute':
        req = fetch((c) => c.guilds.member(member.guild_id, member.user.id).mute(rep));
        break;
    }

    try {
      const res = await req;
      onSubmitted(res);
      onClose();
      if (goBack) nav(-1);
      pushNotification({
        message: t('components:modalcreatereport.successful'),
        type: 'SUCCESS',
      });
    } catch (e) {}
  };

  const action = t(ACTION_TEXT[type]);

  return (
    <Modal
      show={show}
      onClose={onClose}
      heading={action}
      controls={
        <>
          <Button disabled={!reason} onClick={_submit}>
            <Trans ns="components" i18nKey="modalcreatereport.controls.execute">
              {{ action }}
            </Trans>
          </Button>
          <Button variant="gray" onClick={onClose}>
            {t('components:modalcreatereport.controls.cancel')}
          </Button>
        </>
      }>
      <ModalContainer>
        <section>
          <Heading>{t('components:modalcreatereport.reason')}</Heading>
          <ModalTextArea value={reason} onInput={(e) => setReason(e.currentTarget.value)} />
        </section>
        {(type === 'ban' || type === 'mute') && (
          <section>
            <Heading>{t('components:modalcreatereport.timeout')}</Heading>
            <DurationPicker value={timeout} onDurationInput={(v) => setTimeout(v)} />
          </section>
        )}
        <section>
          <Heading>{t('components:modalcreatereport.attachment')}</Heading>
          <Filedrop file={attachment} onFileInput={_setAttachment} />
        </section>
      </ModalContainer>
    </Modal>
  );
};
