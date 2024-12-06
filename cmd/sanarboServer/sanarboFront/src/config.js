import { Log } from './log/index';

export const APP = 'sanarboFront';
export const APP_TITLE = 'Goéland';
export const VERSION = '0.1.6';
export const BUILD_DATE = '2024-12-06';
// eslint-disable-next-line no-undef
export const DEV = process.env.NODE_ENV === 'development';
// eslint-disable-next-line no-restricted-globals
export const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;
export const HOME = DEV ? 'http://localhost:5173/' : `${BACKEND_URL}/`;
export const apiRestrictedUrl = `goapi/v1`;
export const getLog = (ModuleName, verbosityDev, verbosityProd) => (
  (DEV) ? new Log(ModuleName, verbosityDev) : new Log(ModuleName, verbosityProd)
);
export const DEFAULT_BASE_LAYER = 'orthophotos_ortho_lidar_2016';