<template>
  <div>
    <v-form @submit.prevent="submitForm">
      <v-container>
        <h2> Arbre - {{ Tree.external_id }}</h2>
        <div>{{ Tree.name }}</div>
        <div style="font-style: italic">{{ Tree.description }}</div>
        <v-row class="py-5">
          <v-col cols="12" md="12">
          <!--
            <v-select
                v-model.number="Tree.tree_attributes.idvalidation"
                :items="dictionaries.validation.data"
                item-title="value"
                item-value="id"
                label="Statut"
                disabled
            >
            </v-select>
            -->
            <v-text-field
                v-model.number="statut"
                label="Statut"
                outlined
                readonly
                disabled
            >
            </v-text-field>
          </v-col>
          <v-col cols="12" md="12">
            <v-select
                v-model.number="Tree.tree_attributes.idtobechecked"
                :items="dictionaries.to_be_checked.data"
                item-title="value"
                item-value="id"
                label="À contrôler"
            >
            </v-select>
          </v-col>
          <v-col cols="12" md="12">
            <v-select
                v-model.number="Tree.tree_attributes.idnote"
                :items="dictionaries.note.data"
                item-title="value"
                item-value="id"
                label="Note"
            >
            </v-select>
          </v-col>
        </v-row>
      </v-container>

      <v-container>
        <h2>Environnement</h2>
        <v-row class="py-5">
          <v-col cols="12" md="12">
            <v-text-field
                v-model.number="Tree.tree_attributes.circonference"
                label="Circonférence [cm]"
                type="number"
            >
            </v-text-field>
          </v-col>

          <v-card
              class="mx-auto"
              prepend-icon="mdi-circle-box-outline"
              width="100%"
          >
            <template #title>
              Entourage / Cadre
            </template>
            <v-col cols="12" md="12">
              <v-select
                  v-model.number="Tree.tree_attributes.identourage"
                  :items="dictionaries.entourage.data"
                  item-title="value"
                  item-value="id"
                  label="Type"
              >
              </v-select>
            </v-col>
            <v-col cols="12" md="12">
              <v-select
                  v-model.number="Tree.tree_attributes.idchkentourage"
                  :items="dictionaries.check.data"
                  item-title="value"
                  item-value="id"
                  label="Statut"
              >
              </v-select>
            </v-col>
            <v-col cols="12" md="12">
              <v-textarea 
                  v-model="Tree.tree_attributes.entouragerem"
                  rows="3"
                  auto-grow
                  label="Remarque entourage"
                  type="string"
              ></v-textarea>
            </v-col>
          </v-card>


          <v-card
              class="mx-auto"
              prepend-icon="mdi-texture-box"
              width="100%"
          >
            <template #title>
              Revêtement
            </template>
            <v-col cols="12" md="12">
              <v-select
                  v-model.number="Tree.tree_attributes.idrevsurface"
                  :items="dictionaries.rev_surface.data"
                  item-title="value"
                  item-value="id"
                  label="Type"
              >
              </v-select>
            </v-col>
            <v-col cols="12" md="12">
              <v-select
                  v-model.number="Tree.tree_attributes.idchkrevsurface"
                  :items="dictionaries.check.data"
                  item-title="value"
                  item-value="id"
                  label="Statut"
              >
              </v-select>
            </v-col>
            <v-col cols="12" md="12">
              <v-textarea 
                  v-model="Tree.tree_attributes.revsurfacerem"
                  rows="3"
                  auto-grow
                  label="Remarque revêtement"
                  type="string"
              ></v-textarea>
            </v-col>
          </v-card>
        </v-row>


        <h2 class="pt-10">État sanitaire</h2>
        <v-row class="py-5">
          <v-col cols="12" md="12">
            <v-select
                v-model.number="Tree.tree_attributes.idetatsanitairepied"
                :items="dictionaries.etat_sanitaire.data"
                item-title="value"
                item-value="id"
                label="Pied"
            >
            </v-select>
          </v-col>
          <v-col cols="12" md="12">
            <v-select
                v-model.number="Tree.tree_attributes.idetatsanitairetronc"
                :items="dictionaries.etat_sanitaire.data"
                item-title="value"
                item-value="id"
                label="Tronc"
            >
            </v-select>
          </v-col>
          <v-col cols="12" md="12">
            <v-select
                v-model.number="Tree.tree_attributes.idetatsanitairecouronne"
                :items="dictionaries.etat_sanitaire.data"
                item-title="value"
                item-value="id"
                label="Couronne"
            >
            </v-select>
          </v-col>
          <v-col cols="12" md="12">
            <v-textarea 
                v-model="Tree.tree_attributes.etatsanitairerem"
                rows="3"
                auto-grow
                label="Remarque état sanitaire">
            </v-textarea>
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12" md="2">
            <v-btn color="primary" type="submit" @click="submitForm">Sauver</v-btn>
          </v-col>
          <v-col cols="12" md="2">
            <v-btn color="secondary" type="button" @click="handleFormCanceled">Annuler</v-btn>
          </v-col>
        </v-row>
      </v-container>
    </v-form>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue';
import {useFetch} from "../composables/FetchData.js";
import { BACKEND_URL, apiRestrictedUrl } from '../config.js';
import { getLocalJwtTokenAuth } from './Login.js';


const backendUrl = `${BACKEND_URL}/${apiRestrictedUrl}/`;
const urlTrees = backendUrl + "trees";


const emit = defineEmits(['formSubmitted', 'formCanceled'])
const props = defineProps({
  showForm: {type: Boolean, required: false, default: false},
  treeId: {type: Number, required: false, default: 0},
  dictionaries: {type: Object, required: true, default: null}
})


const Tree = reactive({
  external_id: '',
  is_active: '',
  is_validated: '',
  create_time: '',
  creator: '',
  description: '',
  id: '',
  name: '',
  geom: '',
  tree_attributes: {},
});

const statut = ref('');

// Get session storage token
//const token = sessionStorage.getItem('token');
const token = getLocalJwtTokenAuth();
const headers = {
  'Authorization': token,
  'Content-Type': 'application/json',
}

const options = {
  headers: headers
}


onMounted(async () => {

  const tree = await useFetch(urlTrees + '/' + props.treeId, options)
  Tree.external_id = tree.data.value.external_id;
  Tree.is_active = tree.data.value.is_active;
  Tree.is_validated = false;
  Tree.create_time = tree.data.value.create_time;
  Tree.creator = tree.data.value.creator;
  Tree.description = tree.data.value.description;
  Tree.id = tree.data.value.id;
  Tree.name = tree.data.value.name;
  Tree.tree_attributes = tree.data.value.tree_attributes;
  Tree.geom = tree.data.value.geom;
  statut.value = props.dictionaries.validation.data.find(x => x.id === Tree.tree_attributes.idvalidation).value;
})


// eslint-disable-next-line no-unused-vars
const submitForm = async () => {

  const options = {
    headers: headers,
    method: 'PUT',
    body: JSON.stringify(Tree)
  }

  await useFetch(urlTrees + '/' + props.treeId, options)

  emit('formSubmitted', JSON.stringify(Tree));
};


const handleFormCanceled = () => {
  emit('formCanceled')
}
</script>
