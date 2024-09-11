import { Log } from './log/index';

export const APP = 'sanarboFront';
export const APP_TITLE = 'Goéland';
export const VERSION = '0.1.0';
export const BUILD_DATE = '2024-08-13';
// eslint-disable-next-line no-undef
export const DEV = process.env.NODE_ENV === 'development';
export const HOME = DEV ? 'http://localhost:5173/' : '/';
// eslint-disable-next-line no-restricted-globals
const url = new URL(location.toString());
export const BACKEND_URL = DEV ? 'http://localhost:9090' : url.origin;
export const apiRestrictedUrl = `goapi/v1`;
export const getLog = (ModuleName, verbosityDev, verbosityProd) => (
  (DEV) ? new Log(ModuleName, verbosityDev) : new Log(ModuleName, verbosityProd)
);
export const DEFAULT_BASE_LAYER = 'orthophotos_ortho_lidar_2016';