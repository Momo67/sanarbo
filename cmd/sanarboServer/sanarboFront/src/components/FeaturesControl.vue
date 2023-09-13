<script setup>
import { computed, ref } from 'vue';
import { getValidationColor } from './features.js';


const props = defineProps({
  showFeatures: {
    type: Boolean,
    required: false,
    default: false,
  },
  validations: {
    type: Object,
    required: true
  },
  //dictionaries: {type: Object, required: true, default: null},
  validationToShow: {
    type: Array,
    required: true
  }
});

const validationToShow = ref(props.validationToShow);
//const showFeatures = ref(props.showFeatures);
const showFeatures = computed({
  get() {
    return props.showFeatures;
  },
  set(value) {
    emit('show-changed', value);
  }
});

const emit = defineEmits(['selected-validation', 'show-changed']);

const showFeaturesOnClick = () => {
  showFeatures.value = !showFeatures.value;
  //emit('show-changed', showFeatures.value);
}

const selectFeaturesOnOK = () => {
  showFeatures.value = false;
  emit('selected-validation', validationToShow.value);
  //emit('show-changed', showFeatures.value);
}
const selectFeaturesOnCancel = () => {
  showFeatures.value = false;
  //emit('show-changed', showFeatures.value);
}

</script>

<template>
  <div d-flex>
    <v-container fluid class="ol-custom features-control">
      <v-tooltip top>
        <template #activator="{ props }">
          <v-btn v-bind="props" :class="{ 'btn-showfeatures-on': showFeatures, 'btn-showfeatures-off': !showFeatures }" icon="mdi-pine-tree" density="default" @click="showFeaturesOnClick"></v-btn>
        </template>
        <slot name="tooltip">
          <span>Sélection types de validation</span>
        </slot>
      </v-tooltip>
    </v-container>
    <v-container v-show="showFeatures" class="features-selection">
      <v-row>
        <v-col class="v-col-xs-12 v-col-sm-6 offset-sm-3 v-col-md-4 offset-md-4 v-col-lg-4 offset-lg-4">
          <v-card>
            <v-card-item>
              <v-card-title primary-title>
                Types de validation
              </v-card-title>
              <v-card-subtitle>
                Sélection
              </v-card-subtitle>
            </v-card-item>
            <v-divider></v-divider>
            <v-card-text class="card-text">
              <template v-for="(validation, key) in validations.data" :key="key">
                <v-container v-if="validation.id != 2" style="height: 3.5em;">
                  <v-checkbox v-model="validationToShow" :color="getValidationColor(validation.id)" :label="validation.value" :value="validation.id" class="chk-box-validation"></v-checkbox>
                </v-container>
              </template>
              <br/>
            </v-card-text>
            <v-divider></v-divider>
            <v-card-actions>
              <v-btn color="info" @click="selectFeaturesOnOK">OK</v-btn>
              <v-btn color="info" @click="selectFeaturesOnCancel">Annuler</v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<style scoped>
.btn-showfeatures-on {
  background-color: white;
  color: darkcyan;
}

.btn-showfeatures-off {
  background-color: white;
  color: black;
}

.features-selection {
  position: fixed;
  z-index: 1000;
  top: 10em;
  left: 50%;
  -webkit-transform: translateX(-50%);
  -ms-transform: translateX(-50%);
  transform: translateX(-50%);
}

.card-text {
  transform: translateY(-1em);
  -webkit-transform: translateY(-1em);
  -ms-transform: translateY(-1em);
}
</style>