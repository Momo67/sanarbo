<script setup>
import Map from 'ol/Map.js';
import View from 'ol/View.js';
import TileLayer from 'ol/layer/Tile.js';
import VectorLayer from 'ol/layer/Vector.js';
import VectorSource from 'ol/source/Vector.js';
import {WKT} from "ol/format.js";
import {onMounted, defineProps} from "vue";
import { OSM } from 'ol/source';
import proj4 from 'proj4'
import OlProjection from 'ol/proj/Projection'
import {register} from 'ol/proj/proj4';
import 'ol/ol.css'


// Define props
const props = defineProps(['data'])


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

const features = props.data.map(d=> format.readFeature(d.geom, {
  featureProjection: 'EPSG:2056'
}))



// Define vector layer
const vector = new VectorLayer({
  source: new VectorSource({
    features: features
})
});




onMounted(() => {

 
// Create map
const map = new Map({
    view: new View({
      center: 	[	2537850.0, 1152445.0],
      zoom: 12,
      projection: swissProjection
    }),
    layers: [
      new TileLayer({
        source: new OSM(),
      }), 
      vector
      
    ],
    target: 'map',
  });
})


</script>


<template>
  <div id="map">
  </div>
</template>


<style scoped>
#map {
  width: 100vw;
  height: 100vh;
}
</style>
