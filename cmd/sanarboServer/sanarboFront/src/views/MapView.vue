<script setup>
import {useFetch} from "../composables/FetchData.js";
import OlMap from "../components/OlMap.vue";
import { onMounted, ref } from "vue";

// env variables
const backendUrl = import.meta.env.VITE_BACKEND_URL;
const urlTrees = backendUrl + "/goapi/v1/trees"
const auth = 'Bearer ' + import.meta.env.VITE_JWT_SECRET

// Fetch options
const options = {
  credentials: 'include',
  headers: {
    'Authorization': auth
  }
}

let trees = ref([])


// Fetch data
onMounted(() => {
  let {data} = useFetch(urlTrees, options);
  trees = data

})


</script>

<template>


<OlMap :data="trees">
</OlMap>


</template>

<style scoped>


</style>
