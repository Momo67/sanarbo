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
        <v-col class="v-col-xs-12 v-col-sm-6 offset-sm-3 v-col-md-6 offset-md-4 v-col-lg-5 offset-lg-4 v-col-xl-4">
          <v-card>
            <v-card-item>
              <v-card-title primary-title>
                Recherche d'arbres
              </v-card-title>
            </v-card-item>
            <v-divider></v-divider>
            <v-card-text>
              <v-form ref="form" @submit.prevent="submitForm">
                <v-container>

                  <v-row class="py-1">
                    <v-col cols="12" md="12">
                      <v-text-field v-model="idThing" clearable label="Identifiant de l'arbre" :rules="[rules.valid]" @click:clear="onClear"></v-text-field>
                    </v-col>
                  </v-row>

                  <v-row class="py-1">
                    <v-col cols="4" md="4">
                      <v-select
                        v-model="secteurName"
                        :items="gestion_com.secteurs.data"
                        item-title="value"
                        item-value="value"
                        label="Secteur"
                      >
                      </v-select>
                    </v-col>
                    <v-col cols="8" md="8">
                      <v-autocomplete
                        v-model.number="idEmplacement"
                        auto-select-first
                        :items="gestion_com.emplacements.data"
                        no-data-text="Aucune donnée"
                        item-title="value"
                        item-value="id"
                        label="Emplacement"
                      >
                      </v-autocomplete>
                    </v-col>
                  </v-row>

                  <v-row class="py-1">
                    <v-col cols="8" md="8">
                      <v-autocomplete
                        v-model="idStreet"
                        auto-select-first
                        :items="things.streets.data"
                        no-data-text="Aucune donnée"
                        item-title="value"
                        item-value="id"
                        label="Rue"
                      >
                       <template #item="{ props, item }">
                        <v-list-item v-bind="props" :subtitle="item.raw.subtitle"></v-list-item>
                       </template> 
                      </v-autocomplete>
                    </v-col>
                    <v-col cols="4" md="4">
                      <v-select
                        v-model.number="idAddress"
                        :items="things.buildings.data"
                        no-data-text="Aucune donnée"
                        item-title="value"
                        item-value="id"
                        label="N°"
                      >
                      </v-select>
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
            <v-card-text>
              <v-row class="v-col-xs-12 v-col-sm-12 v-col-md-12 v-col-lg-12 v-col-xl-9">
                <v-col>
                  <v-btn color="primary" type="submit" :disabled="submitBtnDisabled" @click="submitForm">OK</v-btn>
                </v-col>
                <v-col class="v-col-xs-6 v-col-sm-9 v-col-md-9 v-col-lg-9 v-col-xl-9">
                  <v-btn color="secondary" type="button" @click="searchTreeOnCancel">Annuler</v-btn>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
    
  </div>
</template>

<script setup>
import { computed, ref, reactive, watch, onMounted } from 'vue';
import proj4 from 'proj4';
import { WKT } from "ol/format.js";
import OlProjection from 'ol/proj/Projection'
import { register } from 'ol/proj/proj4';
import { useFetch } from "../composables/FetchData.js";
import { BACKEND_URL, apiRestrictedUrl } from '../config.js';
import { getLocalJwtTokenAuth } from './Login.js';

const backendUrl = `${BACKEND_URL}/${apiRestrictedUrl}/`;
const urlGestionCom = backendUrl + "gestion_com";
const urlStreets = backendUrl + "thing/streets";
const urlBuildings = backendUrl + "thing/buildings";

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

const secteurName = ref('');

const idEmplacement = ref();

const idStreet = ref();

const idAddress = ref();

const showAlert = ref(false);

const textAlert = ref('');

const form = ref(null);

let res = {
  hasError: null, 
  errorMessage: null, 
  isLoading: null, 
  data: null,
};

let secteurs = {data: []};
let emplacements = {data: []};

const gestion_com = ref({
  secteurs: secteurs,
  emplacements: emplacements,
});

let streets = {data: []};
let buildingsNumbers = {data: []};

const things = ref({
  streets: streets,
  buildings: buildingsNumbers,
});

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

// Define projection
proj4.defs(
    'EPSG:2056',
    '+proj=somerc +lat_0=46.95240555555556 +lon_0=7.439583333333333 +k_0=1 +x_0=2600000 +y_0=1200000 +ellps=bessel +towgs84=674.374,15.056,405.346,0,0,0,0 +units=m +no_defs');

register(proj4);

const swissProjection = reactive(new OlProjection({
  code: 'EPSG:2056',
  units: 'm',
}));

// wkt format
const wktFormat = new WKT();

const searchTreeById = (idTree) => {
  const features = props.featureSource.getFeatures();
  const length = features.length;
  let tree = null;
  for (let count = 0; count < length; count++) {
    if (features[count].get('idthing') === parseInt(idTree)) {
      tree = features[count];
      break;
    }
  }
  return { feature: tree, zoom: 12 };
}

const getZoomLevel = (surface) => {
  let zoom = null;
  zoom = 9 - Math.floor(surface / 20000);

  return zoom;
}

const searchEmplacementCenter = async (idEmplacement) => {
  const centroid = await useFetch(urlGestionCom + '/emplacements/centroid/' + idEmplacement, options);

  let feature = wktFormat.readFeature(centroid.data.value.geometry, {
    featureProjection: swissProjection,
  });
  let surface = centroid.data.value.surface;

  return { feature: feature, zoom: getZoomLevel(surface)};
}

const searchBuildingCenter = async (idAddress) => {
  const center = await useFetch(urlBuildings + '/center/' + idAddress, options);

  let feature = wktFormat.readFeature(center.data.value.geometry, {
    featureProjection: swissProjection,
  });

  return { feature: feature, zoom: 10};
}

const resetFields = () => {
  showSearchTrees.value = false;
  idThing.value = null;
  secteurName.value = null;
  idEmplacement.value = null;
  idStreet.value = null;
  idAddress.value = null;
  showAlert.value = false;
  gestion_com.value = {
    secteurs: secteurs,
    emplacements: {data: []},
  }
}

const submitForm = async () => {
  let center = null;

  if (idThing.value != null) {
    center = searchTreeById(idThing.value);
  }
  else if (idEmplacement.value != null) {
    center = await searchEmplacementCenter(idEmplacement.value);
  } else if (idAddress.value != null) {
    center = await searchBuildingCenter(idAddress.value);
  }

  if (center.feature !== null) {
    emit('coords-found', {
      coords: center.feature.getGeometry().getCoordinates(),
      zoom: center.zoom
    });
    resetFields();
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
  resetFields();
}

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

watch(secteurName, async () => {
  if ((secteurName.value != '') && (secteurName.value != null)) {
    idEmplacement.value = null;

    emplacements = await useFetch(urlGestionCom + '/emplacements' + (secteurName.value != '' ? ('/' + secteurName.value) : ''), options);
    gestion_com.value = {
      secteurs: secteurs,
      emplacements: emplacements
    };
  }
});

watch(idStreet, async () => {
  textAlert.value = '';
  showAlert.value = false;
  things.value = {
    streets: streets,
    buildings: {data: []},
  };

  if (idStreet.value != null) {
    idAddress.value = null;

    res = await useFetch(urlBuildings + '/numbers/' + idStreet.value, options);
    if (!res.hasError.value) {
      buildingsNumbers = res;
      things.value = {
        streets: streets,
        buildings: buildingsNumbers,
      };
    } else {
      textAlert.value = 'Il n\'y a pas de numéros d\'immeuble pour cette rue!';
      showAlert.value = true;
    }
  }
});

onMounted(async () => {
  secteurs = await useFetch(urlGestionCom + '/secteurs', options);
  gestion_com.value = {
    secteurs: secteurs,
    emplacements: {data: []},
  };

  streets = await useFetch(urlStreets, options);
  things.value = {
    streets: streets,
    buildings: {data: []},
  };
})

</script>

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
  top: 1em;
  left: 50%;
  -webkit-transform: translateX(-50%);
  -ms-transform: translateX(-50%);
  transform: translateX(-50%);
}
</style>