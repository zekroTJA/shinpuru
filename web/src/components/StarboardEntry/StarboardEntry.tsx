import { GuildStarboardEntry } from '../../lib/shinpuru-ts/src';
import { MediaTile } from './MediaTile';
import { ReactComponent as StarIcon } from '../../assets/starfilled.svg';
import styled from 'styled-components';
import { uid } from 'react-uid';

type Props = {
  entry: GuildStarboardEntry;
};

const StyledContainer = styled.div`
  padding: 1em;
  background-color: ${(p) => p.theme.background2};
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  gap: 1em;
  cursor: pointer;
  transition: transform 0.2s ease;

  &:hover {
    transform: scale(1.01);
  }

  > header {
    display: flex;
    gap: 1em;
    align-items: center;

    > img {
      width: 2em;
      aspect-ratio: 1;
      border-radius: 100%;
    }
  }

  > main {
    display: flex;
    flex-direction: column;
    gap: 0.5em;
    background-color: ${(p) => p.theme.background3};
    border-radius: 8px;
    padding: 0.7em;
  }
`;

const ImageCotnainer = styled.div`
  display: flex;
  gap: 0.5em;

  img,
  video {
    width: 100%;
    max-height: 20em;
    object-fit: contain;
    object-position: 0 0;
  }
`;

const StarCount = styled.div`
  display: flex;
  color: ${(p) => p.theme.yellow};
  display: flex;
  gap: 0.4em;
  align-items: center;
`;

export const StarboardEntry: React.FC<Props> = ({ entry }) => {
  const _onClick = () => {
    window.open(entry.message_url);
  };

  return (
    <StyledContainer onClick={_onClick}>
      <header>
        <img src={entry.author_avatar_url} alt="avatar" />
        <span>{entry.author_username}</span>
        <StarCount>
          <StarIcon />
          <span>{entry.score}</span>
        </StarCount>
      </header>
      <main>
        <span>{entry.content}</span>
        {entry.media_urls.length > 0 && (
          <ImageCotnainer>
            {entry.media_urls.map((e) => (
              <MediaTile key={uid(e)} url={e} />
            ))}
          </ImageCotnainer>
        )}
      </main>
    </StyledContainer>
  );
};
