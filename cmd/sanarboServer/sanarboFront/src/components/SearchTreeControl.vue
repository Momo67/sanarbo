<script setup>
import { computed, ref } from 'vue';
/*
import {useFetch} from "../composables/FetchData.js";

const backendUrl = import.meta.env.VITE_BACKEND_API_URL;
const urlTrees = backendUrl + "trees";
*/

const props = defineProps({
  featureSource: {
    type: Object,
    required: true
  },
  showSearchTrees: {
    type: Boolean,
    required: false,
    default: false
  }
});

const showSearchTrees = computed({
  get() {
    return props.showSearchTrees;
  },
  set(value) {
    emit('show-changed', value);
  }
});

const submitBtnDisabled = computed({
  get: () => {
    if (form.value != null)
      return (form.value.errors.length != 0);
    else
      return false;
  }
});

const idThing = ref();

const showAlert = ref(false);

const textAlert = ref('');

const form = ref(null);

const emit = defineEmits(['coords-found', 'show-changed']);

const showTreesOnClick = () => {
  showSearchTrees.value = !showSearchTrees.value;
}

const rules = ref({
  valid: (value) => isValid(value) || 'Valeur invalide'
});

const isValid = (value) => {
  if ((value != undefined) &&(value.length != 0)) {
    const reg = /^[1-9][0-9]*$/;
    return reg.test(value);
  } else {
    return true;
  }
}
/*
const token = sessionStorage.getItem('token');
const headers = {
  'Authorization': 'Bearer ' + token,
  'Content-Type': 'application/json',
}
*/

const submitForm = () => {
  console.log('### validate:', form.value.errors.length);
  const features = props.featureSource.getFeatures();
  const length = features.length;
  let found = null;
  for (let count = 0; count < length; count++) {
    if (features[count].get('idthing') === parseInt(idThing.value)) {
      found = features[count];
      break;
    }
  }

  if (found !== null) {
    emit('coords-found', found.getGeometry().getCoordinates());
    showSearchTrees.value = false;
    idThing.value = null;
    showAlert.value = false;
  } else {
    textAlert.value = 'Aucun arbre trouvé!';
    showAlert.value = true;
  }
};

const onClear = () => {
  showAlert.value = false;
}

const alertOnClose = () => {
  showAlert.value = false;
}

const searchTreeOnCancel = () => {
  showSearchTrees.value = false;
  showAlert.value = false;
  idThing.value = null;
}
</script>

<template>
  <div d-flex>
    <v-container fluid class="ol-custom sesrch-control">
      <v-tooltip top>
        <template #activator="{ props }">
          <v-btn v-bind="props" :class="{ 'btn-treesearch-on': showSearchTrees, 'btn-treesearch-off': !showSearchTrees }" icon="mdi-magnify" density="default" @click="showTreesOnClick"></v-btn>
        </template>
        <slot name="tooltip">
          <span>Recherche d'un arbre</span>
        </slot>
      </v-tooltip>
    </v-container>
    <v-container v-show="showSearchTrees" class="tree-search">
      <v-row>
        <v-col class="v-col-xs-12 v-col-sm-6 offset-sm-3 v-col-md-4 offset-md-4 v-col-lg-4 offset-lg-4">
          <v-card>
            <v-card-item>
              <v-card-title primary-title>
                Recherche d'arbres
              </v-card-title>
            </v-card-item>
            <v-divider></v-divider>
            <v-card-text style="height: 300px;">
              <v-form ref="form" @submit.prevent="submitForm">
                <v-container>

                  <v-row class="py-5">
                    <v-col cols="12" md="12">
                      <v-text-field v-model="idThing" clearable label="Identifiant de l'arbre" :rules="[rules.valid]" @click:clear="onClear"></v-text-field>
                    </v-col>
                  </v-row>

                </v-container>
              </v-form>
            </v-card-text>
            <v-card-text>

              <v-alert v-model="showAlert" type="warning" :text="textAlert" closable close-label="Fermer" @click:close="alertOnClose">
              </v-alert>
            
            </v-card-text>
            <v-divider></v-divider>
            <v-card-actions>
              <v-btn color="primary" :disabled="submitBtnDisabled" @click="submitForm">OK</v-btn>
              <v-btn color="secondary" @click="searchTreeOnCancel">Annuler</v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
    
  </div>
</template>

<style scoped>
.btn-treesearch-on {
  background-color: white;
  color: darkcyan;
}

.btn-treesearch-off {
  background-color: white;
  color: black;
}

.tree-search {
  position: fixed;
  z-index: 1000;
  top: 10em;
  left: 50%;
  -webkit-transform: translateX(-50%);
  -ms-transform: translateX(-50%);
  transform: translateX(-50%);
}
</style>