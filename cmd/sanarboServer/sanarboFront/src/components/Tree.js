import axios from 'axios';
import { functionExist, getErrorMessage } from '../tools/utils';
import { apiRestrictedUrl, BACKEND_URL, getLog } from '../config';
import { getLocalJwtTokenAuth } from './Login';

const log = getLog('Tree', 4, 1);

// User Singleton and stateless class to get and persist data to backend
const tree = {
  getList: (callbackLoaded) => {
    const method = 'getList';
    log.t(`## IN ${method}`);
    axios.defaults.headers.common.Authorization = getLocalJwtTokenAuth();
    axios.get(`${BACKEND_URL}/${apiRestrictedUrl}/trees`)
      .then((resp) => {
        log.t(`## IN ${method} axios get success resp.data :`, resp.data);
        if (functionExist(callbackLoaded)) {
          callbackLoaded(resp.data, 'SUCCESS');
        }
      })
      .catch((err) => {
        const errMessage = getErrorMessage(method, `## ERREUR RESEAU DANS ${method} PENDANT UN APPEL DISTANT axios.get`, err, log);
        if (functionExist(callbackLoaded)) callbackLoaded(err, errMessage);
      });
  },

  getTree: (idTree, callbackLoaded) => {
    const method = 'getTree';
    log.t(`## IN ${method}`);
    axios.defaults.headers.common.Authorization = getLocalJwtTokenAuth();
    axios.get(`${BACKEND_URL}/${apiRestrictedUrl}/trees/${idTree}`)
      .then((resp) => {
        log.t(`## IN ${method} axios get success resp.data :`, resp.data);
        if (functionExist(callbackLoaded)) {
          callbackLoaded(resp.data, 'SUCCESS');
        }
      })
      .catch((err) => {
        const errMessage = getErrorMessage(method, `## ERREUR RESEAU DANS ${method} PENDANT UN APPEL DISTANT axios.get`, err, log);
        if (functionExist(callbackLoaded)) callbackLoaded(err, errMessage);
      });
  },

  treesToValidate: (secteur, epmplacement, callbackLoaded) => {
    const method = 'treesToValidate';
    log.t(`## IN ${method}`);
    axios.defaults.headers.common.Authorization = getLocalJwtTokenAuth();
    axios.get(`${BACKEND_URL}/${apiRestrictedUrl}/validation?secteur=${secteur}&emplacement=${epmplacement}`)
      .then((resp) => {
        log.t(`## IN ${method} axios put success resp.data :`, resp.data);
        if (functionExist(callbackLoaded)) {
          callbackLoaded(resp.data, 'SUCCESS');
        }
      })
      .catch((err) => {
        const errMessage = getErrorMessage(method, `## ERREUR RESEAU DANS ${method} PENDANT UN APPEL DISTANT axios.put`, err, log);
        if (functionExist(callbackLoaded)) callbackLoaded(err, errMessage);
      });
  },

  validateTrees: (data, callbackLoaded) => {
    const method = 'validateTrees';
    log.t(`## IN ${method}`);
    axios.defaults.headers.common.Authorization = getLocalJwtTokenAuth();
    axios.post(`${BACKEND_URL}/${apiRestrictedUrl}/validation`, data)
      .then((resp) => {
        log.t(`## IN ${method} axios put success resp.data :`, resp.data);
        if (functionExist(callbackLoaded)) {
          callbackLoaded(resp.data, 'SUCCESS');
        }
      })
      .catch((err) => {
        const errMessage = getErrorMessage(method, `## ERREUR RESEAU DANS ${method} PENDANT UN APPEL DISTANT axios.put`, err, log);
        if (functionExist(callbackLoaded)) callbackLoaded(err, errMessage);
      });
  },

};
// prevents modification to properties and values of the user singleton
Object.freeze(tree);
export default tree;
