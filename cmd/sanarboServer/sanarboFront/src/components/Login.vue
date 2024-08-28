<script setup>
defineOptions({
  name: 'UserLogin',
  inheritAttrs: false
});

import {reactive} from "vue";


const loginUrl = import.meta.env.VITE_BACKEND_LOGIN_URL;
const emit = defineEmits(['userLoggedIn', 'authenticationFailed'])


const userName = reactive({
  value: '',
  rules: [value => value ? true : 'Veuillez rentrer votre compte utilisateur']
})

const password = reactive({
  value: '',
  rules: [value => value ? true : 'Veuillez rentrer votre mot de passe'],
  showPassword: false
})


const authenticateUser = async () => {
  // Create the form data with URL-encoded key-value pairs
  const formData = new URLSearchParams();
  formData.append('login', userName.value);
  formData.append('pass', password.value);

  /*

  const login = {
    password_hash: `${password.value}`,
    usename: `${userName.value}`
  };
  */


  // Perform the authentication request
  const response = await fetch(loginUrl, {
    method: 'POST',
    body: formData.toString(),
    //body: login,
    headers: {'Content-Type': 'application/x-www-form-urlencoded'}
  });

  // Store the JWT token in local storage
  const data = await response.json();

  return data
}

const submitForm = async () => {

  // Authenticate user
  const response = authenticateUser();

  try {
    const data = await response;

    // Emit a custom event to notify the parent component and pass the token
    if (data.token) {
      emit('userLoggedIn', data.token);
    }
  } catch (error) {
    // Handle error if token retrieval fails
    emit('authenticationFailed');
  }
}


</script>

<template>
  <v-form @submit="submitForm" @submit.prevent>
    <v-container>
      <v-row>
        <v-col
            cols="12"
            md="12"
        >
          <v-text-field
              v-model="userName.value"
              :rules="userName.rules"
              label="Compte utilisateur"
              required
          ></v-text-field>
        </v-col>

        <v-col
            cols="12"
            md="12"
        >
          <v-text-field
              v-model="password.value"
              :append-icon="password.showPassword ? 'mdi-eye' : 'mdi-eye-off'"
              :rules="password.rules"
              :type="password.showPassword ? 'text' : 'password'"
              label="Mot de passe"
              required
              @click:append="password.showPassword = !password.showPassword"
          ></v-text-field>
        </v-col>

        <v-col
            cols="12"
            md="12"
        >
        </v-col>

      </v-row>
      <v-btn block class="mt-2" color="secondary" type="submit">Se connecter</v-btn>
    </v-container>
  </v-form>
</template>



