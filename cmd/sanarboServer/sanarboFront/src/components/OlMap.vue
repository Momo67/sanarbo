<script setup>
import {onMounted, ref, reactive} from "vue";
import Map from 'ol/Map.js';
import View from 'ol/View.js';
import TileLayer from 'ol/layer/Tile.js';
import VectorLayer from 'ol/layer/Vector.js';
import VectorSource from 'ol/source/Vector.js';
import OlSourceWMTS from 'ol/source/WMTS';
import OlTileGridWMTS from 'ol/tilegrid/WMTS';
import {WKT} from "ol/format.js";
import OlProjection from 'ol/proj/Projection'
import {register} from 'ol/proj/proj4';
import {useFetch} from "../composables/FetchData.js";
import {Select} from "ol/interaction.js";
import {Fill, Stroke, Style, Circle as CircleStyle, Text as TextStyle} from 'ol/style.js';
import Control from 'ol/control/Control.js';
import proj4 from 'proj4';
import 'ol/ol.css'
import TreeForm from "./TreeForm.vue";
import TrackingControl from "./TrackingControl.vue";
import LayersControl from "./LayersControl.vue";
import FeaturesControl from "./FeaturesControl.vue";
import SearchTreeControlVue from "./SearchTreeControl.vue";
import { tile_layers, default_tile_grid} from "./layers.js"
import { getValidationColor } from './features.js';
import { DEFAULT_BASE_LAYER } from '../config.js';


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

const showControlLayers = ref(false);
const showControlFeatures = ref(false);
const showSearchTrees = ref(false);
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
});

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

const layers = ref([]);
const selectedLayer = ref(DEFAULT_BASE_LAYER);

tile_layers.forEach((layer) => {
  let new_layer = new TileLayer({
    title: layer.title,
    type: layer.type,
    source: new OlSourceWMTS({
      layer: layer.layer,
      url: layer.url,
      tileGrid: new OlTileGridWMTS(default_tile_grid),
      requestEncoding: layer.requestEncoding
    }),
    visible: layer.visible
  });
  layers.value.push(new_layer);
});

const hiddenFeatureSource = new VectorSource();
const hiddenFeatureLayer = new VectorLayer({
  source: hiddenFeatureSource,
  visible: false
});
layers.value.push(hiddenFeatureLayer);

const arbreStyle = (feature, resolution) => {
  const color = getValidationColor(feature.get('idvalidation'));
  return new Style({
    image: new CircleStyle({
      radius: 5 / (resolution + 0.5),
      fill: new Fill({ color: color }),
      stroke: new Stroke({ width: 1, color: color }),
    }),
  });
}
const featureSource = new VectorSource();
const vectorLayer = new VectorLayer({
  id: 'arbre_layer',
  source: featureSource,
  style: arbreStyle,
  visible: true
});
layers.value.push(vectorLayer);

const arbreIdStyle = (feature, resolution) => {
  return new Style({
    text: new TextStyle({
      text: String(feature.get('idthing')),
      font: '10px Arial',
      offsetY: -25 / (resolution + 1),
      fill: new Fill({ color: 'rgb(255, 255, 255)' }),
      //stroke: new Stroke({color: 'rgb(255, 255, 255)', width: 1}),
      scale: 1 / (resolution + 0.5)
    })
  });
}
const textLayer = new VectorLayer({
  id: 'arbre_id_layer',
  source: featureSource,
  style: arbreIdStyle,
  maxResolution: 0.2,
  visible: true
});
layers.value.push(textLayer);

const displayed_features = ref([1, 5, 6, 7, 8, 9, 10, 11]);

const filterFeatures = (featuresToShow) => {
  featureSource.clear();

  const filter = (query) => {
    return hiddenFeatureSource.getFeatures().filter(feature => query.includes(feature.get('idvalidation')));
  }

  featureSource.addFeatures(filter(featuresToShow));
}

const chooseFeatures = (selected) => {
  displayed_features.value = selected;
  filterFeatures(displayed_features.value);
}

const controls = [
  {
    name: 'layers',
    state: null,
    displayed: showControlLayers
  },
  {
    name: 'features',
    state: null,
    displayed: showControlFeatures
  },
  {
    name: 'search-tree',
    state: null,
    displayed: showSearchTrees
  },
];

const switchOffControls = (exception) => {
  controls.forEach(control => {
    if (control.name !== exception) {
      control.displayed.value = false;
    }
  });
}

const controlFeaturesOnClick = (state) => {
  showControlFeatures.value = state;
  if (showControlFeatures.value) {
    switchOffControls('features');
  }
}

const controlSearchTreeOnClick = (state) => {
  showSearchTrees.value = state;
  if (showSearchTrees.value) {
    switchOffControls('search-tree');
  }
}

const controlLayersOnClick = (state) => {
  showControlLayers.value = state;
  if (showControlLayers.value) {
    switchOffControls('layers');
  }
}

const chooseLayer = (selected) => {
  selectedLayer.value = selected;
  const map_layers = map.getLayers();
  map_layers.forEach((layer) => {
    const layerName = layer.get('source').layer_;
    if (layer.get('type') === 'base') {
      if (layerName === selected) {
        layer.setVisible(true);
      } else {
        layer.setVisible(false);
      }
    }
  });
}

const setDefaultBaseLayer = () => {
  const map_layers = map.getLayers();
  map_layers.forEach((layer) => {
    const layerName = layer.get('source').layer_;
    if ((layerName === DEFAULT_BASE_LAYER) && (layer.get('type') === 'base')) {
      layer.setVisible(true);
    }
  });
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

let map = null;
const view = new View({
  center: [2537633.0, 1152618.0],
  zoom: 18,
  projection: swissProjection
});

const setPosition = (position) => {
  view.animate({
    center: position,
    duration: 2000,
  });
}
//Tracking
const trackingEnabled = ref(false);

onMounted(async () => {
  const {hasError, errorMessage, isLoading, data} = await useFetch(urlTrees, options);
  errorFetch.value = hasError.value;
  fetchIsLoading.value = isLoading.value;
  errorFetchMessage.value = errorMessage.value;
  
  fetchDictionaries();
  
  const features = data.value.map((d) => {
    let feature = wktFormat.readFeature(d.geom, {
      featureProjection: swissProjection,
    })

    feature.set('id', d.id)
    feature.set('idthing', d.external_id)
    feature.set('idvalidation', d.tree_att_light.idvalidation)

    return feature
  });

  hiddenFeatureSource.addFeatures(features);

  filterFeatures(displayed_features.value)

  map = new Map({
    controls: [],
    view: view,
    layers: layers.value,
    target: 'map',
  });

  const myControl = new Control({
    element: document.getElementById("expandCustomControl")
  });
  map.addControl(myControl);

  map.addInteraction(selectInteraction)
  setDefaultBaseLayer();
});
</script>


<template>

  <div id="expandCustomControl" >
    <TrackingControl 
      :tracking-enabled="trackingEnabled" 
      :projection="swissProjection" 
      class="ol-custom tracking-control" 
      @position-changed="setPosition">
    </TrackingControl>
    <LayersControl 
      :show-layers="showControlLayers" 
      :layers="tile_layers" 
      :current-layer="selectedLayer" 
      class="ol-custom layers-control" 
      @show-changed="controlLayersOnClick" 
      @selected-layer="chooseLayer">
    </LayersControl>
    <FeaturesControl 
      :show-features="showControlFeatures" 
      :validations="dictionaries.validation" 
      :validation-to-show="displayed_features" 
      class="ol-custom features-control" 
      @show-changed="controlFeaturesOnClick" 
      @selected-validation="chooseFeatures">
    </FeaturesControl>
    <SearchTreeControlVue
      :show-search-trees="showSearchTrees"
      class="ol-custom search-control"
      @show-changed="controlSearchTreeOnClick">
    </SearchTreeControlVue>
  </div>  

  <v-select
    v-model="selectedLayer"
    :items="tile_layers"
    item-title="title"
    item-value="layer"
    label="Choix des couches"
    width="100%"
    @update:model-value="chooseLayer"
  ></v-select>
  <div>showControlLayers:{{ showControlLayers }},&nbsp;showControlFeatures:{{ showControlFeatures }}</div>

  <div id="map" ref="mymap">
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
        <TreeForm 
          :show-form='showForm' 
          :tree-id="treeId" 
          :dictionaries="dictionaries" 
          @form-canceled="handleFormCanceled"
          @form-submitted='handleFormSubmitted'>
        </TreeForm>
      </v-card-text>
    </v-card>
  </v-dialog>

</template>


<style scoped>
#map {
  width: 100vw;
  height: 100vh;
}

#expandCustomControl {
  max-width: auto;
  max-height: auto;
  margin: 0px; /* important to ensure the custom control is not centered since the container has margin: auto by default */
}

.ol-custom.tracking-control {
  position: relative;
  z-index: 1000;
  top: 1.0em;
  left: -moz-calc(100% - 32px);
  left: -webkit-calc(100% - 32px);
  left: calc(100% - 100px);
}

.ol-custom.layers-control {
  position: relative;
  z-index: 1000;
  top: 0.5em;
  left: -moz-calc(100% - 32px);
  left: -webkit-calc(100% - 32px);
  left: calc(100% - 100px);
}

.ol-custom.features-control {
  position: relative;
  z-index: 1000;
  top: 0.5em;
  left: -moz-calc(100% - 32px);
  left: -webkit-calc(100% - 32px);
  left: calc(100% - 100px);
}

.ol-custom.search-control {
  position: relative;
  z-index: 1000;
  top: 0.5em;
  left: -moz-calc(100% - 32px);
  left: -webkit-calc(100% - 32px);
  left: calc(100% - 100px);
}
</style>
