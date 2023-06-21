<script setup>
import {onMounted, reactive, ref} from 'vue';
import {useFetch} from "../composables/FetchData.js";


const backendUrl = import.meta.env.VITE_BACKEND_API_URL;
const urlTrees = backendUrl + "trees"


const emit = defineEmits(['formSubmitted', 'formCanceled'])
const props = defineProps({
  showForm: Boolean,
  treeId: Number,
  newTreeGeometry: String
})


const Tree = reactive({
  create_time: '',
  creator: '',
  description: '',
  id: '',
  name: '',
  geom: props.newTreeGeometry,
  tree_attributes: {},
});


// Get session storage token
const token = sessionStorage.getItem('token');
const headers = {
  'Authorization': 'Bearer ' + token,
  'Content-Type': 'application/json',
}

const options = {
  headers: headers
}


onMounted(async () => {

  const {data} = await useFetch(urlTrees + '/' + props.treeId, options)
  Tree.create_time = data.value.create_time;
  Tree.creator = data.value.creator;
  Tree.description = data.value.description;
  Tree.id = data.value.id;
  Tree.name = data.value.name;
  Tree.tree_attributes = data.value.tree_attributes;
  Tree.geom = data.value.geom;


})

const submitForm = async (event) => {
  const options = {
    headers: headers,
    method: 'PUT',
    body: JSON.stringify(Tree)
  }

  await useFetch(urlTrees + '/' + props.treeId, options)
  // Emit a custom event to notify the parent component and pass the token
  emit('formSubmitted');
};

const handleFormCanceled = () => {
  emit('formCanceled')
}


</script>


<template>
  <v-form @submit.prevent="submitForm">
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

      <v-row>
        <v-col cols="2" md="1">
          <v-btn type="submit" color="primary" @click="submitForm">Sauvegarder</v-btn>
        </v-col>
        <v-col cols="2" md="1">
          <v-btn type="button" color="secondary" @click="handleFormCanceled">Annuler</v-btn>
        </v-col>

      </v-row>

    </v-container>
  </v-form>
</template>
