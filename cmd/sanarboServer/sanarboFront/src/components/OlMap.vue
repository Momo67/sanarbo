<script setup>
import Map from 'ol/Map.js';
import View from 'ol/View.js';
import TileLayer from 'ol/layer/Tile.js';
import VectorLayer from 'ol/layer/Vector.js';
import VectorSource from 'ol/source/Vector.js';
import {WKT} from "ol/format.js";
import {onMounted, ref} from "vue";
import {OSM} from 'ol/source';
import proj4 from 'proj4'
import OlProjection from 'ol/proj/Projection'
import {register} from 'ol/proj/proj4';
import 'ol/ol.css'
import {useFetch} from "../composables/FetchData.js";
import Tree from "./Tree.vue";
import {Select} from "ol/interaction.js";
import {click} from "ol/events/condition.js";


// Fetch data
const backendUrl = import.meta.env.VITE_BACKEND_API_URL;
const urlTrees = backendUrl + "trees"
const token = sessionStorage.getItem('token');

const headers = {'Authorization': 'Bearer ' + token}

const options = {
  headers: headers
}

const errorFetch = ref(false);
const errorFetchMessage = ref('');
const fetchIsLoading = ref(true);

const showForm = ref(false);
const treeId = ref(null);

// wkt format
const wktFormat = new WKT();

// Interactions
const selectInteraction = new Select({
});

selectInteraction.on('select', (event) => {
      if (event.selected.length > 0) {
        showForm.value = true;
        const selectedFeature = event.selected[0];
        treeId.value = selectedFeature.get('id');

      } else {
        showForm.value = false;
      }
}
)

// Handle form submission / cancel
const formSubmitted = ref(false);
const handleFormSubmitted = () => {
  formSubmitted.value = true;
  showForm.value = false;
  selectInteraction.getFeatures().clear();

}

const handleFormCanceled = () => {
  showForm.value = false;
  selectInteraction.getFeatures().clear();
}


onMounted(   async () => {

  const {hasError, errorMessage, isLoading, data} = await useFetch(urlTrees, options);
  errorFetch.value = hasError.value;
  fetchIsLoading.value = isLoading.value;
  errorFetchMessage.value = errorMessage.value;


  // Define projection
  proj4.defs(
      'EPSG:2056',
      '+proj=somerc +lat_0=46.95240555555556 +lon_0=7.439583333333333 +k_0=1 +x_0=2600000 +y_0=1200000 +ellps=bessel +towgs84=674.374,15.056,405.346,0,0,0,0 +units=m +no_defs');

  register(proj4);

  const swissProjection = new OlProjection({
    code: 'EPSG:2056',
    units: 'm',
  });



const features = data.value.map((d) => {

  let feature = wktFormat.readFeature(d.geom, {
    featureProjection: swissProjection,
  })

  feature.set('id', d.id)

  return feature
});


// Define vector layer
const vectorLayer = new VectorLayer({
    source: new VectorSource({
      features: features
    })
  });


const map = new Map({
    view: new View({
      center: 	[2537850.0, 1152445.0],
      zoom: 12,
      projection: swissProjection
    }),
    layers: [
      new TileLayer({
        source: new OSM(),
      }),
      vectorLayer
    ],
    target: 'map',
  });

  map.addInteraction(selectInteraction)

});
</script>


<template>
  <div id="map">
    <div v-if="fetchIsLoading">Loading...</div>
    <div v-else-if="errorFetch">Error: {{ errorFetchMessage }}</div>
  </div>

  <v-dialog
      v-model="showForm"
      scrollable
      width="auto"
  >
    <v-card>
      <v-card-text>
        <Tree :showForm='showForm' @formSubmitted='handleFormSubmitted' @formCanceled="handleFormCanceled" :tree-id="treeId" ></Tree>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>


<style scoped>
#map {
  width: 100vw;
  height: 100vh;
}
</style>
