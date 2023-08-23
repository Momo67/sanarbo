<script setup>
import { ref, onMounted } from 'vue';

const props = defineProps({
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
const showLayers = ref(false);

const emit = defineEmits(['selected-layer']);

const showLayersOnClick = () => {
  showLayers.value = !showLayers.value;
}

const showLayersOnOK = () => {
  emit('selected-layer', currentLayer.value);
  showLayers.value = false;
}

const showLayersOnCancel = () => {
  showLayers.value = false;
}

const layerOnClick = (layer) => {
  currentLayer.value = layer;
}

onMounted(() => {
});
</script>

<template>
  <div>
    <v-container fluid class="ol-custom layers-control">
      <v-tooltip top>
        <template #activator="{ props }">
          <v-btn v-bind="props" :class="{ 'btn-showlayers-on': showLayers, 'btn-showlayers-off': !showLayers }" icon="mdi-layers-outline" density="default" @click="showLayersOnClick"></v-btn>
        </template>
        <span>Sélection couche de base</span>
      </v-tooltip>
    </v-container>
    <v-container v-show="showLayers" class="layers-selection">
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
        <v-card-text style="height: 300px;">
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
        <v-card-actions>
          <v-btn color="info" @click="showLayersOnOK">OK</v-btn>
          <v-btn color="info" @click="showLayersOnCancel">Annuler</v-btn>
        </v-card-actions>
      </v-card>
    </v-container>
  </div>
</template>

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
  width: auto;
  position: fixed;
  z-index: 1000;
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
  position: relative;
  left: -moz-calc(100% - 8px);
  left: -webkit-calc(100% - 8px);
  left: calc(100% - 8px);
  color: green;
}

.layer-icon-notselected {
  color: grey;
  display: none;
}
</style>