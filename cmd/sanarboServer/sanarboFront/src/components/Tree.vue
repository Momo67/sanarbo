<script setup>
import {onMounted, reactive} from 'vue';
import {useFetch} from "../composables/FetchData.js";


const backendUrl = import.meta.env.VITE_BACKEND_API_URL;
const urlTrees = backendUrl + "trees"

const props = defineProps({
  newTree: Boolean,
  treeId: Number
})


const Tree = reactive({
  create_time: '',
  creator: '',
  description: '',
  id: '',
  name:'',
});

onMounted(async () => {

  if (!props.newTree) {
    const token = localStorage.getItem('token');
    const headers = {'Authorization': 'Bearer ' + token}
    const options = {
      headers: headers
    }

    const {data} = await useFetch(urlTrees + '/' + props.treeId, options)
    Tree.create_time = data.value.create_time;
    Tree.creator = data.value.creator;
    Tree.description = data.value.description;
    Tree.id = data.value.id;
    Tree.name = data.value.name;
  }
})





const submitForm = () => {
  // Handle form submission logic here
  console.log(Tree);
};



</script>


<template>
  <v-form @submit="submitForm">
    <v-container>
      <v-row>
        <v-col cols="12" md="6">
          <v-text-field v-model="Tree.id" label="ID"></v-text-field>
        </v-col>
        <v-col cols="12" md="6">
          <v-text-field
              v-model="Tree.create_time"
              label="Date de création"
              type="datetime"
          ></v-text-field>
        </v-col>
        <v-col cols="12" md="6">
          <v-text-field v-model="Tree.creator" label="Créateur"></v-text-field>
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12" md="6">
          <v-text-field v-model="Tree.description" label="Description"></v-text-field>
        </v-col>

      </v-row>
      <v-row>
        <v-col cols="12" md="6">
          <v-text-field v-model="Tree.name" label="Nom"></v-text-field>
        </v-col>
      </v-row>
      <v-btn type="submit" color="primary">Sauvegarder</v-btn>
    </v-container>
  </v-form>
</template>
