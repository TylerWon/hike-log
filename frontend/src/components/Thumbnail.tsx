import { type Photo } from "../data/hikes";

interface ThumbnailProps {
  photo?: Photo;
  trail: string;
}

export default function Thumbnail({ photo, trail }: ThumbnailProps) {
  if (photo) {
    return (
      <div className="shrink-0 rounded overflow-hidden bg-forest-800 w-[96px] h-[72px]">
        <img
          src={photo.url}
          alt={photo.caption ?? `${trail} cover photo`}
          className="w-full h-full object-cover transition-transform duration-300 group-hover/hike-card:scale-105"
          loading="lazy"
        />
      </div>
    );
  }

  return (
    <div className="shrink-0 rounded bg-olive-850 flex items-center justify-center w-[96px] h-[72px]">
      <svg className="stroke-forest-700" width="24" height="24" strokeWidth="1.5" viewBox="0 0 24 24" fill="none">
        <path d="M3 17L9 9l4 5 3-3 5 6H3z" />
      </svg>
    </div>
  );
}
