<script setup>
import Login from "./components/Login.vue";
import OlMap from "./components/OlMap.vue";
import {reactive} from "vue";

const authState = reactive({
  isLoggedIn : localStorage.getItem('token') ? true : false,
  token : localStorage.getItem('token') ? localStorage.getItem('token') : null
})


const handleUserLoggedIn = (receivedToken) => {
    authState.isLoggedIn = true;
    authState.token = receivedToken;
    localStorage.setItem('token', authState.token)
}
</script>



<template>
  <div>
    <template  v-if="authState.isLoggedIn">
      <!-- Render Map -->
      <OlMap :token="authState.token"/>
    </template>
    <template v-else>
      <!-- Render the login page -->
      <Login @userLoggedIn="handleUserLoggedIn" />
    </template>
  </div>
</template>
