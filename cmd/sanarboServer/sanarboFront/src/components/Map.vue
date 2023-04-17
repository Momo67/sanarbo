<script setup>
import {useFetch} from "../composables/FetchData.js";

// env variables
const backendUrl = import.meta.env.VITE_BACKEND_URL;
const urlTrees = backendUrl + "/goapi/v1/trees"
const auth = 'Bearer ' + import.meta.env.VITE_JWT_SECRET

// fetch options
const options = {
  credentials: 'include',
  headers: {
    'Authorization': auth
  }
}

// fetch data
const {data} = useFetch(urlTrees, options);




</script>

<template>
  <ul>
    <li v-for="tree in data">
      {{tree.geom}}
    </li>
  </ul>
</template>

<style scoped>

</style>
