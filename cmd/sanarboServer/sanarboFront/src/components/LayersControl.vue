<template>
  <div d-flex>
    <v-container fluid class="ol-custom layers-control">
      <v-tooltip top>
        <template #activator="{ props }">
          <v-btn v-bind="props" :class="{ 'btn-showlayers-on': showLayers, 'btn-showlayers-off': !showLayers }" icon="mdi-layers-outline" density="default" @click="showLayersOnClick"></v-btn>
        </template>
        <slot name="tooltip">
          <span>Sélection couche de base</span>
        </slot>
      </v-tooltip>
    </v-container>
    <v-container v-show="showLayers" class="layers-selection">
      <v-row>
        <v-col class="v-col-xs-12 v-col-sm-6 offset-sm-3 v-col-md-6 offset-md-4 v-col-lg-5 offset-lg-4 v-col-xl-4">
          <v-card>
            <v-card-item>
              <v-card-title primary-title>
                Couche de base
              </v-card-title>
              <v-card-subtitle>
                Sélection
              </v-card-subtitle>
            </v-card-item>
            <v-divider></v-divider>
            <v-card-text style="height: 340px;">
              <template v-for="(layer, key) in layers" :key="key">
                <v-container style="height: 3.5em;">
                  <v-btn append-icon="mdi-check" :block="true" class="btn-layer" @click="layerOnClick(layer.layer.toLowerCase())">
                    {{ layer.title }}
                    <template #append>
                      <v-icon :class="{ 'layer-icon-selected': layer.layer.toLowerCase() === currentLayer, 'layer-icon-notselected': layer.layer !== currentLayer }"></v-icon>
                    </template>
                  </v-btn>
                  <br/>
                </v-container>
              </template>
            </v-card-text>
            <v-divider></v-divider>
            <v-card-text>
              <v-row class="v-col-xs-12 v-col-sm-12 v-col-md-12 v-col-lg-12 v-col-xl-9">
                <v-col>
                  <v-btn color="success" type="submit" @click="showLayersOnOK">OK</v-btn>
                </v-col>
                <v-col class="v-col-xs-6 v-col-sm-9 v-col-md-9 v-col-lg-9 v-col-xl-9">
                  <v-btn color="error" type="button" @click="showLayersOnCancel">Annuler</v-btn>
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
import { computed, ref } from 'vue';

const props = defineProps({
  showLayers: {
    type: Boolean,
    required: false,
    default: false
  },
  layers: {
    type: Array,
    required: true
  },
  currentLayer: {
    type: String,
    required: true
  }
});

const layers = ref(props.layers);
const currentLayer = ref(props.currentLayer);
let oldLayer = currentLayer.value;
const showLayers = computed({
  get() {
    return props.showLayers;
  },
  set(value) {
    emit('show-changed', value);
  }
});

const emit = defineEmits(['selected-layer', 'show-changed']);

const showLayersOnClick = () => {
  showLayers.value = !showLayers.value;
}

const showLayersOnOK = () => {
  showLayers.value = false;
  oldLayer = currentLayer.value;
  emit('selected-layer', currentLayer.value); 
}

const showLayersOnCancel = () => {
  showLayers.value = false;
  currentLayer.value = oldLayer;
  emit('selected-layer', currentLayer.value);
}

const layerOnClick = (layer) => {
  currentLayer.value = layer;
  emit('selected-layer', currentLayer.value);
}

</script>

<style scoped>
.btn-showlayers-on {
  background-color: white;
  color: darkcyan;
}

.btn-showlayers-off {
  background-color: white;
  color: black;
}

.layers-selection {
  position: fixed;
  z-index: 1000;
  top: 1em;
  left: 50%;
  -webkit-transform: translateX(-50%);
  -ms-transform: translateX(-50%);
  transform: translateX(-50%);
}

.btn-layer {
  position: relative;
  left: 0%;
  -webkit-transform: translateX(0%);
  -ms-transform: translateX(0%);
  transform: translateX(0%);
}

.layer-icon-selected {
  position: absolute;
  right: -moz-calc(30%);
  right: -webkit-calc(30%);
  right: calc(15%);
  color: green;
}

.layer-icon-notselected {
  color: grey;
  display: none;
}
</style>
