<script setup>
import {onMounted, reactive, ref} from 'vue';
import {useFetch} from "../composables/FetchData.js";


const backendUrl = import.meta.env.VITE_BACKEND_API_URL;
const urlTrees = backendUrl + "trees"


const emit = defineEmits(['formSubmitted', 'formCanceled'])
const props = defineProps({
  showForm: Boolean,
  treeId: Number,
})


const Tree = reactive({
  create_time: '',
  creator: '',
  description: '',
  id: '',
  name: '',
  geom: '',
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

  // Emit a custom event to notify the parent component
  emit('formSubmitted');
};


const handleFormCanceled = () => {
  emit('formCanceled')
}


</script>


<template>
  <div>


    <v-form @submit.prevent="submitForm">
      <v-container>
        <h2>{{Tree.name}}</h2>
        <v-row>

          <v-col cols="12" md="4">
            <v-text-field
                v-model="Tree.tree_attributes.idtobechecked"
                label="À contrôler"
            ></v-text-field>
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field
                v-model="Tree.tree_attributes.idvalidation"
                label="Statut"
            ></v-text-field>
          </v-col>

          <v-col cols="12" md="4">
            <v-text-field
                v-model="Tree.tree_attributes.idnote"
                label="Note"
            ></v-text-field>
          </v-col>
        </v-row>

        <h3>Environnement</h3>
        <v-row>
          <v-col cols="12" md="12">
            <v-text-field
                v-model="Tree.tree_attributes.circonference"
                label="Circonférence">
            </v-text-field>
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field
                v-model="Tree.tree_attributes.identourage"
                label="Entourage / cadre"
            ></v-text-field>
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field v-model="Tree.tree_attributes.idchkentourage" label="ID chk entourage"></v-text-field>
          </v-col>
        </v-row>

        <h3>État sanitaire</h3>
        <v-row>
          <v-col cols="12" md="4">
            <v-text-field
                v-model="Tree.tree_attributes.idetatsanitairepied"
                label="Pied"></v-text-field>
          </v-col>
          <v-col cols="12" md="4">
          <v-text-field
              v-model="Tree.tree_attributes.idetatsanitairetronc"
              label="Tronc"></v-text-field>
        </v-col>
          <v-col cols="12" md="4">
            <v-text-field
                v-model="Tree.tree_attributes.idetatsanitairecouronne"
                label="Couronne"></v-text-field>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12" md="4">
            <v-text-field v-model="Tree.tree_attributes.entouragerem" label="Entourage remarque"></v-text-field>
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field
                v-model="Tree.tree_attributes.entouragerem"
                label="Remarque entourage"
            ></v-text-field>
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field
                v-model="Tree.tree_attributes.idchkrevsurface"
                label="Remarque"
            ></v-text-field>
          </v-col>
        </v-row>


        <v-row>
          <v-col cols="12" md="12">
            <v-text-field
                v-model="Tree.tree_attributes.etatsanitairerem"
                label="Remarque état sanitaire">
            </v-text-field>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12" md="2">
            <v-btn type="submit" color="primary" @click="submitForm">Sauvegarder</v-btn>
          </v-col>
          <v-col cols="12" md="2">
            <v-btn type="button" color="secondary" @click="handleFormCanceled">Annuler</v-btn>
          </v-col>
        </v-row>

      </v-container>
    </v-form>
  </div>
</template>
