type Props = {
  url: string;
};

export const MediaTile: React.FC<Props> = ({ url }) => {
  if (url.match(/.*\.(?:jpe?g|png|webp|gif|tiff|svg)/g)) return <img src={url} alt="" />;
  return <video controls muted src={url} />;
};
