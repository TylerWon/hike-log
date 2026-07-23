import type { Photo } from "../schemas/photo";

interface ThumbnailProps {
  photo: null | Photo;
  trailName: string;
}

export default function Thumbnail({ photo, trailName }: ThumbnailProps) {
  if (photo) {
    return (
      <div className="shrink-0 rounded overflow-hidden bg-forest-800 w-[96px] h-[72px]">
        <img
          alt={photo.caption ?? `${trailName} cover photo`}
          className="w-full h-full object-cover transition-transform duration-300 group-hover/hike-card:scale-105 group-hover/hike-card:brightness-110"
          loading="lazy"
          src={photo.srcUrl}
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
