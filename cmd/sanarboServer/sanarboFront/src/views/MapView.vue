<script setup>
import {useFetch} from "../composables/FetchData.js";
import OlMap from "../components/OlMap.vue";
import {onBeforeMount, ref, defineProps, reactive} from "vue";

// OpenLayers
import Map from 'ol/Map.js';
import View from 'ol/View.js';
import TileLayer from 'ol/layer/Tile.js';
import VectorLayer from 'ol/layer/Vector.js';
import VectorSource from 'ol/source/Vector.js';
import {WKT} from "ol/format.js";
import { OSM } from 'ol/source';
import proj4 from 'proj4'
import OlProjection from 'ol/proj/Projection'
import {register} from 'ol/proj/proj4';
import 'ol/ol.css'



// Define projection
proj4.defs(
    'EPSG:2056',
    '+proj=somerc +lat_0=46.95240555555556 +lon_0=7.439583333333333 +k_0=1 +x_0=2600000 +y_0=1200000 +ellps=bessel +towgs84=674.374,15.056,405.346,0,0,0,0 +units=m +no_defs');

register(proj4);

const swissProjection = new OlProjection({
  code: 'EPSG:2056',
  units: 'm',
});


// Read features from props
const format = new WKT();


// Fetch data
// env variables
const backendUrl = import.meta.env.VITE_BACKEND_API_URL;
const urlTrees = backendUrl + "trees"

const props = defineProps({
  token: String
})

const headers = {'Authorization': 'Bearer ' + props.token}
const options = {
  headers: headers
}

let trees = ref([])



// const features = trees.value.map(d=> format.readFeature(d.geom, {
//   featureProjection: 'EPSG:2056'
// }))

onBeforeMount( async () => {


  const { data } = await useFetch(urlTrees, options);
  trees.value = data;
//
// // Define vector layer
//   const vector = new VectorLayer({
//     source: new VectorSource({
//       features: features
//     })
//   });

});









</script>

<template>
<ul>
  <li v-for="item in trees" :key="item.id">{{ item.name }}</li>
</ul>

</template>


