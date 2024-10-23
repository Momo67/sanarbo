<script setup>
import { onMounted, ref, reactive } from "vue";
import VectorLayer from 'ol/layer/Vector.js';
import VectorSource from 'ol/source/Vector.js';
import { WKT } from "ol/format.js";
import OlProjection from 'ol/proj/Projection'
import { register } from 'ol/proj/proj4';
import { useFetch } from "../composables/FetchData.js";
import { Select } from "ol/interaction.js";
import { Fill, Stroke, Style, Circle as CircleStyle, Text as TextStyle, RegularShape } from 'ol/style.js';
import Control from 'ol/control/Control.js';
import proj4 from 'proj4';
import 'ol/ol.css'
import TreeForm from "./TreeForm.vue";
import TrackingControl from "./TrackingControl.vue";
import LayersControl from "./LayersControl.vue";
import FeaturesControl from "./FeaturesControl.vue";
import SearchTreeControlVue from "./SearchTreeControl.vue";
import { createLausanneMap } from "./layers.js"
import { getValidationColor } from './features.js';
import { DEFAULT_BASE_LAYER, BACKEND_URL, apiRestrictedUrl } from '../config.js';
import { getLocalJwtTokenAuth } from './Login';


// Fetch data
const backendUrl = `${BACKEND_URL}/${apiRestrictedUrl}/`;
const urlTrees = backendUrl + "trees?offset=0&limit=1000000"
//const token = sessionStorage.getItem('token');
const token = getLocalJwtTokenAuth();

const headers = {'Authorization': token}

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

let selectedFeature = null

selectInteraction.on('select', (event) => {
  if (event.selected.length > 0) {
    showForm.value = true;
    selectedFeature = event.selected[0];
    treeId.value = selectedFeature.get('id');

  } else {
    showForm.value = false;
  }
});

// Handle form submission / cancel
const formSubmitted = ref(false);
const handleFormSubmitted = () => {
  getFeatures();
  
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

let textStyles = {};
let fill = null;
let stroke = null;

const tile_layers = ref([]);


const hiddenFeatureSource = new VectorSource();
const hiddenFeatureLayer = new VectorLayer({
  source: hiddenFeatureSource,
  visible: false
});
layers.value.push(hiddenFeatureLayer);

const arbreStyle = (feature, resolution) => {
  let color = getValidationColor(feature.get('idvalidation'));
  let is_validated = feature.get('is_validated');
  let style = [];

  if (is_validated == false) {
    style.push(new Style({
      image: new RegularShape({
        fill: new Fill({ color: color }),
        stroke: new Stroke({ width: 1, color: color }),
        points: 4,
        radius: 6 / (resolution + 0.5),
        //angle: Math.PI / 4,
      }),
    }));
  } else {
    style.push(new Style({
        image: new CircleStyle({
          radius: 5 / (resolution + 0.5),
          fill: new Fill({ color: color }),
          stroke: new Stroke({ width: 1, color: color }),
        }),
      }),
    );
  }

  return style;
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
      //stroke: new Stroke({color: 'rgb(0, 0, 0)', width: 0.5}),
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

const filterFeatures = (selected, showOnlyValidated, showOnlyPublic) => {
  featureSource.clear();

  const filter = (query) => {
    /*
    if (showOnlyValidated === true)
      return hiddenFeatureSource.getFeatures().filter(feature => query.includes(feature.get('idvalidation')) && feature.get('is_validated') === false);
    else
      return hiddenFeatureSource.getFeatures().filter(feature => query.includes(feature.get('idvalidation')));
    */

    return hiddenFeatureSource.getFeatures().filter(feature => {
      if (showOnlyValidated === true)
      {
        if (showOnlyPublic === true) {
          return query.includes(feature.get('idvalidation')) && (feature.get('is_validated') === false) && (feature.get('ispublic') === false);
        } else {
          return query.includes(feature.get('idvalidation')) && (feature.get('is_validated') === false) && (feature.get('ispublic') === true);
        }
      } else 
        if (showOnlyPublic === true) 
          return query.includes(feature.get('idvalidation')) && (feature.get('ispublic') === false);
        else {
          return query.includes(feature.get('idvalidation')) && (feature.get('ispublic') === true);
        }
    });
  }

  featureSource.addFeatures(filter(selected));
}

const chooseFeatures = (featuresToShow) => {
  let selected = featuresToShow.validationToShow;
  let showOnlyValidated = featuresToShow.showOnlyValidated;
  let showOnlyPublic = featuresToShow.showOnlyPublic;
  displayed_features.value = selected;
  filterFeatures(selected, showOnlyValidated, showOnlyPublic);
}

const coordsFound = (geom) => {
  setPosition(geom);
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
  let textStyle = null;
  selectedLayer.value = selected;
  const map_layers = map.getLayers();
  map_layers.forEach((layer) => {
    const layerName = layer.get('source').layer_;
    if (layer.get('type') === 'base') {
      if (layerName === selected) {
        layer.setVisible(true);
        textStyle = textStyles[selected];
        if (textStyle != null) {
          fill = textStyle.fill;
          stroke = textStyle.stroke;
          textLayer.setStyle(function(feature, resolution) {
            return new Style({
              text: new TextStyle({
                text: String(feature.get('idthing')),
                font: '10px Arial',
                offsetY: -25 / (resolution + 1),
                fill: fill ? new Fill({ color: fill.color ? fill.color : null }) : null,
                stroke: stroke ? new Stroke({color: stroke.color ? stroke.color : null, width: stroke.width ? stroke.width : null}) : null,
                scale: 1 / (resolution + 0.5)
              })
            });
          });
        }
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
      fill = textStyles[layerName].fill;
      stroke = textStyles[layerName].stroke;
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

const setPosition = (position) => {
  let coords = position.coords;
  let zoom = position.zoom;
  console.log('zoom:', zoom);
  if ((parseInt(coords.length) == 2) && (parseInt(coords[0]) > 2000000) && (parseInt(coords[0]) < 2900000) && (parseInt(coords[1]) > 1000000) && (parseInt(coords[1]) < 1300000)) {
    map.getView().animate({
      center: coords,
      duration: 2000,
      zoom: zoom
    });
  }
}
//Tracking
const trackingEnabled = ref(false);

const getFeatures = async () => {
  const {hasError, errorMessage, isLoading, data} = await useFetch(urlTrees, options);
  console.log('### data:', data.value);
  errorFetch.value = hasError.value;
  fetchIsLoading.value = isLoading.value;
  errorFetchMessage.value = errorMessage.value;

  const features = data.value.map((d) => {
    let feature = wktFormat.readFeature(d.geom, {
      featureProjection: swissProjection,
    })

    feature.set('id', d.id);
    feature.set('idthing', d.external_id);
    feature.set('is_validated', d.is_validated);
    feature.set('idvalidation', d.tree_att_light.idvalidation);
    feature.set('ispublic', d.tree_att_light.ispublic);

    return feature;
  });

  hiddenFeatureSource.clear();
  hiddenFeatureSource.addFeatures(features);

  filterFeatures(displayed_features.value)
}

onMounted(async () => {
  getFeatures();

  fetchDictionaries();
  
  (async () => {
		const placeStFrancoisM95 = [2538202, 1152364];
		const myOlMap = await createLausanneMap('map', placeStFrancoisM95, 8, DEFAULT_BASE_LAYER);

    map = myOlMap.map;
    layers.value.forEach((layer) => {
      map.addLayer(layer);
    });

    console.log('### layers:', myOlMap.map.getLayers());
    myOlMap.map.getLayers()
		       .forEach((layer) => {
				     const type = layer.get('type');
				     const source = layer.getSource();
				     if (type === 'base') {
	  			      const currentBaseLayer = source.getLayer();
                console.log(`currentBaseLayer : ${currentBaseLayer}`)
                tile_layers.value.push({title: layer.getProperties().title , layer: source.getLayer()});
				     }
             console.log('### layer:', layer);
		       });

    textStyles = myOlMap.textStyles;
		console.log("myOlMap contains a ref to your OpenLayers Map Object : ", myOlMap);
    
    const myControl = new Control({
      element: document.getElementById("expandCustomControl")
    });
    map.addControl(myControl);
  
    map.addInteraction(selectInteraction)
    setDefaultBaseLayer();
  })();

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
      :feature-source="featureSource"
      class="ol-custom search-control"
      @show-changed="controlSearchTreeOnClick"
      @coords-found="coordsFound">
    </SearchTreeControlVue>

  </div>  

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
  position: relative;
  max-width: 100px;
  max-height: auto;
  margin: 0px; /* important to ensure the custom control is not centered since the container has margin: auto by default */
  left: calc(100% - 120px);
}

.ol-custom.tracking-control {
  position: relative;
  z-index: 1000;
  top: 1.0em;
}

.ol-custom.layers-control {
  position: relative;
  z-index: 1000;
  top: 0.5em;
}

.ol-custom.features-control {
  position: relative;
  z-index: 1000;
  top: 0.5em;
}

.ol-custom.search-control {
  position: relative;
  z-index: 1000;
  top: 0.5em;
}
</style>
