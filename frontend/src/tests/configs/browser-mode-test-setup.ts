// Browser Mode tests render components in isolation, so the app's global stylesheet (normally only imported by
// `main.tsx`) is never loaded. Import it here so Tailwind utility classes actually take effect during tests.
import "../../index.css";
