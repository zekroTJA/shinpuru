import React, { useEffect, useState } from 'react';

import { ReactComponent as APIIcon } from '../../assets/api.svg';
import { ReactComponent as BookIcon } from '../../assets/book.svg';
import { ReactComponent as BugIcon } from '../../assets/bug.svg';
import { ReactComponent as CubeIcon } from '../../assets/dashed-cube.svg';
import { Flex } from '../../components/Flex';
import { ReactComponent as GithubIcon } from '../../assets/github.svg';
import { ReactComponent as IDIcon } from '../../assets/id.svg';
import { LinearGradient } from '../../components/styleParts';
import { Loader } from '../../components/Loader';
import { ReactComponent as LockIcon } from '../../assets/lock-open.svg';
import { MaxWidthContainer } from '../../components/MaxWidthContainer';
import { ReactComponent as MedalIcon } from '../../assets/karma.svg';
import { PrivacyInfo } from '../../lib/shinpuru-ts/src';
import { range } from '../../util/utils';
import styled from 'styled-components';
import { uid } from 'react-uid';
import { useApi } from '../../hooks/useApi';
import { useNavigate } from 'react-router';
import { useTranslation } from 'react-i18next';

type Props = {};

const Link = styled.a`
  color: ${(p) => p.theme.text};
  display: flex;
  gap: 0.5em;
  align-items: center;
  padding: 1em;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.2s ease;

  ${(p) => LinearGradient(p.theme.background3)};

  &:enabled:hover {
    transform: translateY(-3px);
  }
`;

const DetailsList = styled.ul`
  li {
    margin: 0 0 0.3em 0;
  }
`;

const GeneralRoute: React.FC<Props> = () => {
  const { t } = useTranslation('routes.info.general');
  const fetch = useApi();
  const nav = useNavigate();

  const [privacy, setPrivacy] = useState<PrivacyInfo>();

  useEffect(() => {
    fetch((c) => c.etc.privacyInfo())
      .then((r) => setPrivacy(r))
      .catch();
  }, []);

  return (
    <MaxWidthContainer>
      <h2>{t('documentation')}</h2>
      <Flex gap="1em" wrap>
        <Link href="https://github.com/zekroTJA/shinpuru/wiki" target="_blank">
          <BookIcon />
          <span>{t('links.wiki')}</span>
        </Link>
        <Link href="https://github.com/zekroTJA/shinpuru/wiki/Permissions-Guide" target="_blank">
          <LockIcon />
          <span>{t('links.permissions')}</span>
        </Link>
        <Link href="https://github.com/zekroTJA/shinpuru/wiki/REST-API-Docs" target="_blank">
          <APIIcon />
          <span>{t('links.api')}</span>
        </Link>
        <Link href="https://github.com/zekroTJA/shinpuru/wiki/Self-Hosting" target="_blank">
          <CubeIcon />
          <span>{t('links.selfhosting')}</span>
        </Link>
      </Flex>

      <h2>{t('about')}</h2>
      <Flex gap="1em" wrap>
        <Link href="https://github.com/zekroTJA/shinpuru" target="_blank">
          <GithubIcon />
          <span>{t('links.githubrepo')}</span>
        </Link>
        <Link href="https://github.com/zekroTJA/shinpuru/issues" target="_blank">
          <BugIcon />
          <span>{t('links.issues')}</span>
        </Link>
        <Link href="https://github.com/zekroTJA/shinpuru/blob/master/bughunters.md" target="_blank">
          <MedalIcon />
          <span>{t('links.bughunters')}</span>
        </Link>
      </Flex>

      <h2>{t('privacy')}</h2>
      <Flex gap="1em" wrap>
        {(privacy && (
          <>
            <Link href={privacy?.noticeurl} target="_blank">
              <IDIcon />
              <span>{t('links.privacynotice')}</span>
            </Link>
            <Link href="#" onClick={() => nav('/usersettings/privacy')}>
              <IDIcon />
              <span>{t('links.clearuserdata')}</span>
            </Link>
          </>
        )) || (
          <>
            <Loader width="10em" height="3.5em" />
            <Loader width="12em" height="3.5em" />
          </>
        )}
      </Flex>
      <h3>{t('contact')}</h3>
      <DetailsList>
        {(privacy && (
          <>
            {privacy.contact.map((c) => (
              <li key={uid(c)}>
                {(c.url && (
                  <a href={c.url} target="_blank" rel="noreferrer">
                    {c.title}: {c.value}
                  </a>
                )) || (
                  <span>
                    {c.title}: {c.value}
                  </span>
                )}
              </li>
            ))}
          </>
        )) ||
          range(3).map((i) => (
            <li key={uid(i)}>
              <Loader height="1em" width="12em" />
            </li>
          ))}
      </DetailsList>
    </MaxWidthContainer>
  );
};

export default GeneralRoute;
