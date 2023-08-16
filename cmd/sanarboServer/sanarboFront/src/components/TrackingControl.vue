<script setup>
import { ref, reactive, onMounted } from 'vue';
import Geolocation from 'ol/Geolocation.js';

const props = defineProps({
  trackingEnabled: {
    type: Boolean,
    required: true
  },
  olView: {
    type: Object,
    required: true
  }
});

const trackingEnabled = ref(props.trackingEnabled);
const olView = reactive(props.olView);

let geolocation = null;

onMounted(() => {
  geolocation = new Geolocation({
    // enableHighAccuracy must be set to true to have the heading value.
    trackingOptions: {
      enableHighAccuracy: true,
    },
    projection: olView.getProjection(),
  });
  geolocation.setTracking(trackingEnabled.value);
  geolocation.on('change:position', () => {
    olView.animate({
      center: geolocation.getPosition(),
      //duration: 2000,
    });
  });
});
</script>

<template>
  <div>
    <v-container fluid class="ol-custom tracking-control">
      <v-btn :class="{ 'btn-tracking-on': trackingEnabled, 'btn-tracking-off': !trackingEnabled }" :icon="trackingEnabled ? 'mdi-crosshairs-gps' : 'mdi-crosshairs'" density="default" @click="trackingOnClick"></v-btn>
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