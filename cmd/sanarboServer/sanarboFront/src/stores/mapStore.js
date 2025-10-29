import { defineStore} from 'pinia';
import { ref } from 'vue';

export const useMapStore = defineStore('map', () => {
  const zoomTarget = ref({
    idthing: null,
    zoom: null,
    timestamp: 0,
  });

  const setZoomTarget = (idthing, zoom) => {
    zoomTarget.value = {
      idthing: idthing,
      zoom: zoom,
      timestamp: Date.now(),
    };
  };

  return {
    zoomTarget,
    setZoomTarget,
  };
});