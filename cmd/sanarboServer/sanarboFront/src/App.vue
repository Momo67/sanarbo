<template>
  <header>
    <Toolbar>
      <template #start>
        <span class="pl-2 text-white">{{ `${APP_TITLE} v.${VERSION}` }}</span>
      </template>
      <template #end>
        <template v-if="isUserAuthenticated">
          <Button icon="pi pi-sign-out" title="Logout" class="mr-2" @click="logout" />
        </template>
        <Button icon="pi pi-info-circle" title="A propos..." @click="aboutInfo" />
      </template>
    </Toolbar>
  </header>
  <main>
    <FeedBack ref="feedback" :msg="feedbackMsg" :msg-type="feedbackType" :visible="feedbackVisible" />
    <div class="flex card">
      <div class="col-12">
        <div class="justify-content-center">
          <Toast position="top-center" />
          <template v-if="isUserAuthenticated">
            <Tabs v-model:value="activeTab">
              <TabList>
              <!--
                <template v-if="isUserAdmin">
              -->
                  <Tab value="0" :disabled="!isUserAdmin">Utilisateurs</Tab>
                  <Tab value="1" :disabled="!isUserAdmin">Groupes</Tab>
              <!--
                </template>
              -->
                <Tab value="2">Carte</Tab>
                <Tab value="3">Aide</Tab>
              </TabList>
              <TabPanels>
              <!--
                <template v-if="isUserAdmin">
              -->
                  <TabPanel value="0">
                    <ListUsers :display="isUserAuthenticated" @user-invalid-session="logout" />
                  </TabPanel>
                  <TabPanel value="1">
                    <ListGroups :display="isUserAuthenticated" @user-invalid-session="logout" />
                  </TabPanel>
              <!--
                </template>
              -->
                <TabPanel value="2">
                  <OlMap />
                </TabPanel>
                <TabPanel value="3">
                  <div class="help-content">
                    <h3>Aide</h3>
                    <p>Bienvenue dans la section d'aide. Ici, vous trouverez des informations utiles pour vous guider à
                      travers l'application.</p>
                    <ul>
                      <li><strong>Utilisateurs:</strong> Gérer les utilisateurs de l'application.</li>
                      <li><strong>Groupes:</strong> Gérer les groupes d'utilisateurs.</li>
                      <li><strong>Carte:</strong> Visualiser et interagir avec la carte.</li>
                    </ul>
                  </div>
                </TabPanel>
              </TabPanels>
            </Tabs>
            <h4>Connexion réussie de {{ getUserLogin() }} [{{ getUserEmail() }}]</h4>
          </template>
          <template v-else>
            <LoginUser 
              :msg="`Authentification ${APP_TITLE}:`" 
              :backend="BACKEND_URL" 
              :disabled="!isNetworkOk"
              @login-ok="loginSuccess" @login-error="loginFailure" />
          </template>
        </div>
      </div>
    </div>
  </main>
</template>

<script setup>
import Button from 'primevue/button';
import Tabs from 'primevue/tabs';
import TabList from 'primevue/tablist';
import Tab from 'primevue/tab';
import TabPanels from 'primevue/tabpanels';
import TabPanel from 'primevue/tabpanel';
import Toast from 'primevue/toast';
import Toolbar from 'primevue/toolbar';
import { onMounted, ref, computed } from 'vue';
import LoginUser from './components/LoginUser.vue';
import FeedBack from './components/FeedBack.vue';
import {
  getUserId,
  getUserEmail, getUserIsAdmin, getTokenStatus, clearSessionStorage,
  doesCurrentSessionExist, getUserLogin,
} from './components/Login';
import {
  APP, APP_TITLE, BACKEND_URL, BUILD_DATE, VERSION, getLog, HOME,
} from './config';
import { isNullOrUndefined } from './tools/utils';
import ListUsers from './components/ListUsers.vue';
import ListGroups from './components/ListGroups.vue';
import OlMap from './components/OlMap.vue';

const log = getLog(APP, 4, 2);
//const activeIndex = ref('0');
const isUserAuthenticated = ref(false);
const isUserAdmin = ref(false);
const isNetworkOk = ref(true);
const feedback = ref(null);

const feedbackMsg = ref(`${APP_TITLE}, v.${VERSION}`);
const feedbackType = ref('info');
const feedbackVisible = ref(false);
let autoLogoutTimer = null;
const displayFeedBack = (text, type) => {
  log.t(`displayFeedBack() text:'${text}' type:'${type}'`);
  feedbackType.value = type;
  feedbackMsg.value = text;
  feedbackVisible.value = true;
  feedback.value.displayFeedBack(feedbackMsg, feedbackType);
};
const aboutInfo = () => {
  const appInfo = `${APP_TITLE}, v.${VERSION} ${BUILD_DATE}`;
  if (isUserAuthenticated.value) {
    const userInfo = `${getUserLogin()} id[${getUserId()}] Admin:${getUserIsAdmin()}`;
    displayFeedBack(`${appInfo} ⇒ 😊 vous êtes authentifié comme ${userInfo}`, 'info');
  } else {
    displayFeedBack(`${appInfo} ⇒ vous n'êtes pas encore authentifié`, 'info');
  }
  feedbackVisible.value = true;
};

const logout = () => {
  log.t('# IN logout()');
  clearSessionStorage();
  isUserAuthenticated.value = false;
  isUserAdmin.value = false;
  displayFeedBack('Vous vous êtes déconnecté de l\'application avec succès !', 'success');
  if (isNullOrUndefined(autoLogoutTimer)) {
    clearInterval(autoLogoutTimer);
  }
  setTimeout(() => {
    window.location.href = HOME;
  }, 2000); // after 2 sec redirect to home page just in case
};

const checkIsSessionTokenValid = () => {
  log.t('# IN checkIsSessionTokenValid()');
  if (doesCurrentSessionExist()) {
    getTokenStatus()
      .then((val) => {
        if (val instanceof Error) {
          log.e('# getTokenStatus() ERROR err: ', val);
          if (val.message === 'Network Error') {
            displayFeedBack(`Il semble qu'il y a un problème de réseau !${val}`, 'error');
          }
          log.e('# getTokenStatus() ERROR err.response: ', val.response);
          log.w('# getTokenStatus() ERROR err.response.data: ', val.response.data);
          if (!isNullOrUndefined(val.response)) {
            let reason = val.response.data;
            if (!isNullOrUndefined(val.response.data.message)) {
              reason = val.response.data.message;
            }
            log.w(`# getTokenStatus() SERVER SAYS REASON : ${reason}`);
          }
        } else {
          log.l('# getTokenStatus() SUCCESS res: ', val);
          if (isNullOrUndefined(val.err) && (val.status === 200)) {
            // everything is okay, session is still valid
            isUserAuthenticated.value = true;
            isUserAdmin.value = getUserIsAdmin();
            return;
          }
          if (val.status === 401) {
            // jwt token is no more valid
            isUserAuthenticated.value = false;
            isUserAdmin.value = false;
            displayFeedBack('Votre session a expiré !', 'warn');
            logout();
          }
          displayFeedBack(`Un problème est survenu avec votre session erreur: ${val.err}`, 'err');
        }
      })
      .catch((err) => {
        log.e('# getJwtToken() in catch ERROR err: ', err);
        displayFeedBack(`Il semble qu'il y a eu un problème réseau ! erreur: ${err}`, 'error');
      });
  } else {
    log.w('SESSION DOES NOT EXIST OR HAS EXPIRED !');
  }
};

const loginSuccess = (v) => {
  log.t(' loginSuccess()', v);
  isUserAuthenticated.value = true;
  isUserAdmin.value = getUserIsAdmin();
  feedbackVisible.value = false;
  if (isNullOrUndefined(autoLogoutTimer)) {
    // check every 60 seconds(60'000 milliseconds) if jwt is still valid
    autoLogoutTimer = setInterval(checkIsSessionTokenValid, 60000);
  }
};

const loginFailure = (v) => {
  log.w('loginFailure()', v);
  isUserAuthenticated.value = false;
  isUserAdmin.value = false;
};

const activeTab = computed({
  get() {
    return isUserAdmin.value ? '0' : '2';
  },
});

onMounted(() => {
  log.t('mounted()');
  log.w(`${APP} - ${VERSION}, du ${BUILD_DATE}`);

  window.addEventListener('online', () => {
    log.w('ONLINE AGAIN :)');
    isNetworkOk.value = true;
    displayFeedBack('⚡⚡🚀  CONNEXION RESEAU RETABLIE :  😊 vous êtes "ONLINE"  ', 'success');
  });
  window.addEventListener('offline', () => {
    log.w('OFFLINE :((');
    isNetworkOk.value = false;
    displayFeedBack('⚡⚡⚠ PAS DE RESEAU ! ☹ vous êtes "OFFLINE" ', 'error');
  });
});

</script>

<style>
html, body {
  padding: 0;
  margin: 0;
  min-width: 240px;
  font-family: Arial, sans-serif;
}
li {
  margin-left: 30px;
}
</style>
