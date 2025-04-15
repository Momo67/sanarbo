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
        <v-col class="v-col-xs-12 v-col-sm-6 offset-sm-3 v-col-md-6 offset-md-4 v-col-lg-5 offset-lg-4 v-col-xl-4">
          <v-card>
            <v-card-item>
              <v-card-title primary-title>
                Choix d'affichage
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
              <v-container grid-list-xs class="radio-btn-show">
              <v-row>
                <v-col>
                  <v-radio-group v-model="showOnlyValidated" inline>
                    <v-radio label="Tous" :value="false"></v-radio>
                    <v-radio label="Modifiés" :value="true"></v-radio>
                  </v-radio-group>
                </v-col>
                <v-col>
                  <v-radio-group v-model="showOnlyPublic" inline>
                    <v-radio label="Publics" :value="false"></v-radio>
                    <v-radio label="Privés" :value="true"></v-radio>
                  </v-radio-group>
                </v-col>
              </v-row>
              </v-container>
            </v-card-text>
            <v-divider></v-divider>
            <v-card-text>
              <v-row class="v-col-xs-12 v-col-sm-12 v-col-md-12 v-col-lg-12 v-col-xl-9">
                <v-col>
                  <v-btn color="primary" type="submit" @click="selectFeaturesOnOK">OK</v-btn>
                </v-col>
                <v-col class="v-col-xs-6 v-col-sm-9 v-col-md-9 v-col-lg-9 v-col-xl-9">
                  <v-btn color="secondary" type="button" @click="selectFeaturesOnCancel">Annuler</v-btn>
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
  validationToShow: {
    type: Array,
    required: true
  },
});

const validationToShow = ref(props.validationToShow);

const showFeatures = computed({
  get() {
    return props.showFeatures;
  },
  set(value) {
    emit('show-changed', value);
  }
});

const showOnlyValidated = ref(false);
const showOnlyPublic = ref(false);

const emit = defineEmits(['selected-validation', 'show-changed']);

const showFeaturesOnClick = () => {
  showFeatures.value = !showFeatures.value;
}

const selectFeaturesOnOK = () => {
  showFeatures.value = false;
  emit('selected-validation', {validationToShow: validationToShow.value, showOnlyValidated: showOnlyValidated.value, showOnlyPublic: showOnlyPublic.value});
}
const selectFeaturesOnCancel = () => {
  showFeatures.value = false;
}

</script>

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
  top: 1em;
  left: 50%;
  -webkit-transform: translateX(-50%);
  -ms-transform: translateX(-50%);
  transform: translateX(-50%);
}

.card-text {
  height: 36em;
  transform: translateY(-1em);
  -webkit-transform: translateY(-1em);
  -ms-transform: translateY(-1em);
}

.radio-btn-show {
  transform: translateY(1em);
  -webkit-transform: translateY(1em);
  -ms-transform: translateY(1em);
}

.radio-btn-public {
  transform: translateY(-2em);
  -webkit-transform: translateY(-2em);
  -ms-transform: translateY(-2em);
}
</style>