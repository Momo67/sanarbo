<script setup>
import {reactive, ref} from "vue";



const loginUrl = import.meta.env.VITE_BACKEND_LOGIN_URL;
const emit = defineEmits(['userLoggedIn'])
const token = ref(null);

const userName = reactive({
  value : '',
  rules : [value => value ? true : 'Veuillez rentrer votre compte utilisateur']
})

const password = reactive({
  value : '',
  rules : [value => value ? true : 'Veuillez rentrer votre mot de passe'],
  showPassword : false
})

const authenticateUser = async () => {
  // Create the form data with URL-encoded key-value pairs
  const formData = new URLSearchParams();
  formData.append('login', userName.value);
  formData.append('pass', password.value);


  // Perform the authentication request
  const response = await fetch(loginUrl, {
    method: 'POST',
    body: formData.toString(),
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' }
  });

  // Store the JWT token in local storage
  const data = await response.json();

  //localStorage.setItem('token', data.token);

  // Return the JWT token
  return data.token;
}

const submitForm = async () => {

    // Authenticate user
    const response = authenticateUser();

    try {
      const token = await response;

      // Emit a custom event to notify the parent component and pass the token
      emit('userLoggedIn', token);
    } catch (error) {
      // Handle error if token retrieval fails
      console.error('Failed to retrieve token:', error);
    }

}


</script>

<template>
    <v-form @submit.prevent @submit="submitForm">
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
                :rules="password.rules"
                label="Mot de passe"
                :append-icon="password.showPassword ? 'mdi-eye' : 'mdi-eye-off'"
                required
                :type="password.showPassword ? 'text' : 'password'"
                @click:append="password.showPassword = !password.showPassword"
            ></v-text-field>
          </v-col>

          <v-col
              cols="12"
              md="12"
          >
          </v-col>

        </v-row>
        <v-btn type="submit" block class="mt-2">Se connecter</v-btn>
      </v-container>
    </v-form>
</template>



