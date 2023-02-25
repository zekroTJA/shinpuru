import React, { useEffect, useRef, useState } from 'react';

import { CommandInfo } from '../../lib/shinpuru-ts/src';
import { CommandInfoTile } from '../../components/CommandInfoTile';
import { Flex } from '../../components/Flex';
import Fuse from 'fuse.js';
import { MaxWidthContainer } from '../../components/MaxWidthContainer';
import { SearchBar } from '../../components/SearchBar';
import { useApi } from '../../hooks/useApi';
import { useLocation } from 'react-router';
import { useTranslation } from 'react-i18next';

type Props = {};

type CommandMap = { [key: string]: CommandInfo[] };

const CommandsRoute: React.FC<Props> = () => {
  const { t } = useTranslation('routes.info.commands');
  const fetch = useApi();
  const location = useLocation();

  const fuseRef = useRef<Fuse<CommandInfo>>();

  const [commands, setCommands] = useState<CommandInfo[]>();
  const commandsRef = useRef<CommandInfo[]>();

  const _onSearch = (search: string) => {
    if (!search || !fuseRef.current) setCommands(commandsRef.current);
    else {
      const results = fuseRef.current.search(search);
      setCommands(results.filter((r) => r.score! < 0.5).map((r) => r.item));
    }
  };

  useEffect(() => {
    fetch((c) => c.util.slashcommands())
      .then((r) => {
        setCommands(r.data);
        commandsRef.current = r.data;
      })
      .catch();
  }, []);

  useEffect(() => {
    if (!commandsRef.current) return;
    fuseRef.current = new Fuse(commandsRef.current, {
      keys: [
        { name: 'name', weight: 1 },
        { name: 'domain', weight: 0.8 },
        { name: 'description', weight: 0.5 },
      ],
      includeScore: true,
      shouldSort: true,
    });
  }, [commands]);

  useEffect(() => {
    if (!commands || !location.hash) return;
    document.getElementById(location.hash.substring(1))?.scrollIntoView({
      behavior: 'smooth',
      block: 'start',
    });
  }, [location, commands]);

  const commandMap = (() => {
    if (!commands) return undefined;
    const commandMap: CommandMap = {};
    commands.forEach((c) => {
      if (!commandMap[c.group]) commandMap[c.group] = [c];
      else commandMap[c.group].push(c);
    });
    return commandMap;
  })();

  return (
    <MaxWidthContainer>
      <SearchBar onValueChange={_onSearch} placeholder={t('search')} />
      {commandMap &&
        Object.keys(commandMap).map((k) => (
          <>
            <h2>{k}</h2>
            <Flex direction="column" gap="1em">
              {commandMap[k].map((c) => (
                <CommandInfoTile cmd={c} />
              ))}
            </Flex>
          </>
        ))}
    </MaxWidthContainer>
  );
};

export default CommandsRoute;
