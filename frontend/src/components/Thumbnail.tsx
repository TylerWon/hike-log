import { type Photo } from "../data/hikes";

interface ThumbnailProps {
  photo: null | Photo;
  trail: string;
}

export default function Thumbnail({ photo, trail }: ThumbnailProps) {
  if (photo) {
    return (
      <div className="shrink-0 rounded overflow-hidden bg-forest-800 w-[96px] h-[72px]">
        <img
          alt={photo.caption ?? `${trail} cover photo`}
          className="w-full h-full object-cover transition-transform duration-300 group-hover/hike-card:scale-105 group-hover/hike-card:brightness-110"
          loading="lazy"
          src={photo.url}
        />
      </div>
    );
  }

  return (
    <div className="shrink-0 rounded bg-sage-950 flex items-center justify-center w-[96px] h-[72px]">
      <svg className="stroke-forest-700" fill="none" height="24" strokeWidth="1.5" viewBox="0 0 24 24" width="24">
        <path d="M3 17L9 9l4 5 3-3 5 6H3z" />
      </svg>
    </div>
  );
}
