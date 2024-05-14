<script setup>
import * as jose from "jose";
import Login from "./components/Login.vue";
import OlMap from "./components/OlMap.vue";
import { ref, reactive, onMounted } from "vue";

const token = sessionStorage.getItem('token')

const handleUserLoggedIn = (receivedToken) => {
    authState.isLoggedIn = true;
    authState.token = receivedToken;
    sessionStorage.setItem('token', authState.token)
}

const mustReconnect = ref(false);
const reconnectMsg = ref(`Votre session a expiré. 
                          Veuillez vous reconnecter.`);

const isTokenExpired = (token) => {
  let { exp } = jose.decodeJwt(token);
  return (Date.now() >= exp * 1000);
}

const authState = reactive({
  isLoggedIn : token ? true : false,
  token : sessionStorage.getItem('token') ? sessionStorage.getItem('token') : null
})

const reconnectCallback = () => {
  mustReconnect.value = false;
}

onMounted(() => {
  window.setInterval(() => {
    let token = sessionStorage.getItem('token')
    if (isTokenExpired(token)) {
      mustReconnect.value = true;
      authState.isLoggedIn = false;
    }
  }, 60000);
});

</script>

<template>
  <div>
    <template  v-if="authState.isLoggedIn">
      <!-- Render Map -->
      <OlMap/>
    </template>
    <template v-else>

      <v-container>
        <v-dialog
          v-model="mustReconnect"
          max-width="xs6"
          width="550"
          hide-overlay
        >
          <v-alert 
            v-model="mustReconnect"
            style="margin-bottom: 0px;"
            type="warning"
            closable
            prominent
            dismissible
            :text="reconnectMsg"
            class="alert"
            @click:close="reconnectCallback"
          >
          <!--
            <template #text>
              <div>{{ reconnectMsg }}</div>
            </template>
          -->
          </v-alert>
        </v-dialog>
      </v-container>
      <v-container v-if="!mustReconnect" fill-height fluid>
        <v-row 
          align="center"
          justify="center">
          <v-col md="4">
            <Login @userLoggedIn="handleUserLoggedIn" />
          </v-col>
        </v-row>
      </v-container>
    </template>
  </div>
</template>

<style lang="css" scoped>
.alert {
 white-space: pre-line; 
}
</style>