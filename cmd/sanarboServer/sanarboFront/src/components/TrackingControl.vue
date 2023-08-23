<script setup>
import { ref, reactive, onMounted } from 'vue';
import Geolocation from 'ol/Geolocation.js';


const props = defineProps({
  trackingEnabled: {
    type: Boolean,
    required: true
  },
  projection: {
    type: Object,
    required: true
  }
});

const trackingEnabled = ref(props.trackingEnabled);
const projection = reactive(props.projection);

const emit = defineEmits(['toggle-tracking', 'position-changed']);

const geolocation = new Geolocation({
  // enableHighAccuracy must be set to true to have the heading value.
  trackingOptions: {
    enableHighAccuracy: true,
  },
  projection: projection,
});

const trackingOnClick = () => {
  trackingEnabled.value = !trackingEnabled.value;
  geolocation.setTracking(trackingEnabled.value);
  emit('toggle-tracking', trackingEnabled.value);
}


onMounted(() => {
  geolocation.setTracking(trackingEnabled.value);
  geolocation.on('change:position', () => {
    emit('position-changed', geolocation.getPosition());
  });
});
</script>

<template>
  <div>
    <v-container fluid class="ol-custom tracking-control">
      <v-tooltip top>
        <template #activator="{ props }">
          <v-btn v-bind="props" :class="{ 'btn-tracking-on': trackingEnabled, 'btn-tracking-off': !trackingEnabled }" :icon="trackingEnabled ? 'mdi-crosshairs-gps' : 'mdi-crosshairs'" density="default" @click="trackingOnClick"></v-btn>
        </template>
        <slot name="tooltip" v-bind="{ trackingEnabled }">
          <span>{{  trackingEnabled ? 'Suvi GPS activé' : 'Suvi GPS désactivé' }}</span>
        </slot>
      </v-tooltip>
    </v-container>
  </div>
</template>

<style scoped>
.btn-tracking-on {
  background-color: red;
  color: white;
}

.btn-tracking-off {
  background-color: white;
  color: black;
}
</style>