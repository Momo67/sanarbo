<script setup>
import Map from 'ol/Map.js';
import View from 'ol/View.js';
import TileLayer from 'ol/layer/Tile.js';
import VectorLayer from 'ol/layer/Vector.js';
import VectorSource from 'ol/source/Vector.js';
import OlSourceWMTS from 'ol/source/WMTS';
import OlTileGridWMTS from 'ol/tilegrid/WMTS';
import {WKT} from "ol/format.js";
import {onMounted, ref} from "vue";
import {OSM} from 'ol/source';
import proj4 from 'proj4'
import OlProjection from 'ol/proj/Projection'
import {register} from 'ol/proj/proj4';
import 'ol/ol.css'
import {useFetch} from "../composables/FetchData.js";
import TreeForm from "./TreeForm.vue";
import {Select} from "ol/interaction.js";
import {Fill, Stroke, Style, Circle as CircleStyle, Text as TextStyle} from 'ol/style.js';


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
const selectInteraction = new Select({});

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

const dictionaries = ref({
  "validation": {},
  "to_be_checked": {},
  "note": {},
  "check": {},
  "entourage": {},
  "rev_surface": {},
  "etat_sanitaire": {},
  "etat_sanitaire_rem": {}
})

const fetchDictionaries = async () => {
  const dict_validation = await useFetch(backendUrl + 'dico/validation', options);
  const dict_to_be_checked = await useFetch(backendUrl + 'dico/to_be_checked', options);
  const dict_note = await useFetch(backendUrl + 'dico/note', options);
  const dict_entourage = await useFetch(backendUrl + 'dico/entourage', options);
  const dict_check = await useFetch(backendUrl + 'dico/check', options);
  const rev_surface = await useFetch(backendUrl + 'dico/rev_surface', options);
  const etat_sanitaire = await useFetch(backendUrl + 'dico/etat_sanitaire', options);
  const etat_sanitaire_rem = await useFetch(backendUrl + 'dico/etat_sanitaire_rem', options);


  dictionaries.value = {
    "validation": dict_validation,
    "to_be_checked": dict_to_be_checked,
    "note": dict_note,
    "check": dict_check,
    "entourage": dict_entourage,
    "rev_surface": rev_surface,
    "etat_sanitaire": etat_sanitaire,
    "etat_sanitaire_rem": etat_sanitaire_rem
  }
}


onMounted(async () => {

  const {hasError, errorMessage, isLoading, data} = await useFetch(urlTrees, options);
  errorFetch.value = hasError.value;
  fetchIsLoading.value = isLoading.value;
  errorFetchMessage.value = errorMessage.value;

  fetchDictionaries();

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
    feature.set('idthing', d.external_id)
    feature.set('idvalidation', d.tree_att_light.idvalidation)

    return feature
  });

  const getValidationColor = (idvalidation) => {
    let color = '';
    switch (idvalidation) {
      case 1:   //Existant
        color = '#00FF00';
        break;
      case 5:   //En attente de soins
        color = '#FF00FF';
        break;
      case 6:   //En attente d'abattage
        color = '#FFFF00';
        break;
      case 7:   //En attente de remplacement
        color = '#00FFFF';
        break;
      case 8:   //En attente de tomographie
        color = '#0000FF';
        break;
      case 9:   //A surveiller
        color = '#FF0000';
        break;
      case 10:  //En demande d'abattage
        color = '#FF7D00';
        break;
      case 11:  //En attente de projet
        color = '#009696';
        break;
      default:
        color = 'white';
        break;
    };
    return color;
  };

// Define vector layer
  const vectorLayer = new VectorLayer({
    source: new VectorSource({
      features: features
    }),
    style: function(feature, resolution) {
      const color = getValidationColor(feature.get('idvalidation'));
      return new Style({
        image: new CircleStyle({
          radius: 5/(resolution+0.5),
          fill: new Fill({color: color}),
          stroke: new Stroke({width: 1, color: color}),
        }),
      });
    }
  });

  const textLayer = new VectorLayer({
    source: new VectorSource({features: features}),
    style: function(feature, resolution) {
      return new Style({
        text: new TextStyle({
          text: String(feature.get('idthing')),
          font: '10px Arial',
          offsetY: -25/(resolution+1),
          fill: new Fill({color: 'rgb(255, 255, 255)'}),
          stroke: new Stroke({color: 'rgb(255, 255, 255)', width: 1}),
          scale: 1/(resolution+0.5)
        })
      });
    },
    maxResolution: 0.2,
    visible: true
  });

  const ortho2016 = new TileLayer({
    source: new OlSourceWMTS({
      layer: 'orthophotos_ortho_lidar_2016',
      url: `https://tiles01.lausanne.ch/tiles/1.0.0/{Layer}/default/2016/swissgrid_05/{TileMatrix}/{TileRow}/{TileCol}.png`,
      tileGrid: new OlTileGridWMTS({
        origin: [2420000, 1350000],
        resolutions: [50, 20, 10, 5, 2.5, 1, 0.5, 0.25, 0.1, 0.05],
        matrixIds: [0, 1, 2, 3, 4, 5, 6, 7, 8, 9],
      }),
      requestEncoding: 'REST'
    }),
    visible: true
  });

  const ortho2012 = new TileLayer({
    source: new OlSourceWMTS({
      layer: 'orthophotos_ortho_lidar_2012',
      url: `https://tiles01.lausanne.ch/tiles/1.0.0/{Layer}/default/2012/swissgrid_05/{TileMatrix}/{TileRow}/{TileCol}.png`,
      tileGrid: new OlTileGridWMTS({
        origin: [2420000, 1350000],
        resolutions: [50, 20, 10, 5, 2.5, 1, 0.5, 0.25, 0.1, 0.05],
        matrixIds: [0, 1, 2, 3, 4, 5, 6, 7, 8, 9],
      }),
      requestEncoding: 'REST'
    }),
    visible: false
  });

  const swissLayer = new TileLayer({
    source: new OlSourceWMTS({
      layer: 'fonds_geo_carte_nationale_msgroup',
      url: `https://tiles01.lausanne.ch/tiles/1.0.0/{Layer}/default/2014/swissgrid_05/{TileMatrix}/{TileRow}/{TileCol}.png`,
      tileGrid: new OlTileGridWMTS({
        origin: [2420000, 1350000],
        resolutions: [50, 20, 10, 5, 2.5, 1, 0.5, 0.25, 0.1, 0.05],
        matrixIds: [0, 1, 2, 3, 4, 5, 6, 7, 8, 9],
      }),
      requestEncoding: 'REST'
    }),
    visible: false
  });

  const map = new Map({
    controls: [],
    view: new View({
      center: [2537850.0, 1152445.0],
      zoom: 18,
      projection: swissProjection
    }),
    layers: [
      /*
      new TileLayer({
        source: new OSM(),
      }),
      */
      ortho2012,
      ortho2016,
      swissLayer,
      vectorLayer,
      textLayer
    ],
    target: 'map',
  });

  map.addInteraction(selectInteraction)

});
</script>


<template>
  <div id="map" ref="map">
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
        <TreeForm :showForm='showForm' :tree-id="treeId" :dictionaries="dictionaries" @formCanceled="handleFormCanceled"
                  @formSubmitted='handleFormSubmitted'></TreeForm>
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
