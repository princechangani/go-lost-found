/// <reference types="vite/client" />
interface ImportMetaEnv {
    readonly VITE_API_URL: string;
    // add more variables here if needed
    readonly VITE_APP_NAME?: string;
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}