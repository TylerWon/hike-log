/// <reference types="vite/client" />

interface ImportMeta {
  readonly env: ImportMetaEnv;
}

// Add environment variables here for TypeScript IntelliSense to work
interface ImportMetaEnv {
  readonly VITE_API_URL: string;
}

interface ViteTypeOptions {
  strictImportMetaEnv: unknown;
}