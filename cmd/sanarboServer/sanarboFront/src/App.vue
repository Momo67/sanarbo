<script setup>
import Login from "./components/Login.vue";
import OlMap from "./components/OlMap.vue";
import {reactive} from "vue";

const token = sessionStorage.getItem('token')

const handleUserLoggedIn = (receivedToken) => {
    authState.isLoggedIn = true;
    authState.token = receivedToken;
    sessionStorage.setItem('token', authState.token)
}


const authState = reactive({
  isLoggedIn : token ? true : false,
  token : sessionStorage.getItem('token') ? sessionStorage.getItem('token') : null
})

</script>

<template>
  <div>
    <template  v-if="authState.isLoggedIn">
      <!-- Render Map -->
      <OlMap/>
    </template>
    <template v-else>

      <v-container fill-height fluid>
        <v-row align="center"
               justify="center">
          <v-col md="4">
            <Login @userLoggedIn="handleUserLoggedIn" />
          </v-col>
        </v-row>
      </v-container>


    </template>
  </div>
</template>
