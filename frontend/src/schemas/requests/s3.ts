interface UploadPhotosRequestItem {
  index: number;
  photo: File;
  presignedUrl: string;
}

export type UploadPhotosRequest = UploadPhotosRequestItem[];
