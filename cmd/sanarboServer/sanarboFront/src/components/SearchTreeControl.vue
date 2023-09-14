<script setup>
import { computed, ref } from 'vue';

const props = defineProps({
  showSearchTrees: {
    type: Boolean,
    required: false,
    default: false
  }
});

const showSearchTrees = computed({
  get() {
    return props.showSearchTrees;
  },
  set(value) {
    emit('show-changed', value);
  }
});

const emit = defineEmits(['tree-found', 'show-changed']);

const showTreesOnClick = () => {
  showSearchTrees.value = !showSearchTrees.value;
}

const treeId = ref(0);

const searchTreeOnOK = () => {
  showSearchTrees.value = false;
  emit('tree-found', treeId.value);
}

const searchTreeOnCancel = () => {
  showSearchTrees.value = false;
}
</script>

<template>
  <div d-flex>
    <v-container fluid class="ol-custom sesrch-control">
      <v-tooltip top>
        <template #activator="{ props }">
          <v-btn v-bind="props" :class="{ 'btn-treesearch-on': showSearchTrees, 'btn-treesearch-off': !showSearchTrees }" icon="mdi-magnify" density="default" @click="showTreesOnClick"></v-btn>
        </template>
        <slot name="tooltip">
          <span>Recherche d'un arbre</span>
        </slot>
      </v-tooltip>
    </v-container>
    <v-container v-show="showSearchTrees" class="tree-search">
      <v-row>
        <v-col class="v-col-xs-12 v-col-sm-6 offset-sm-3 v-col-md-4 offset-md-4 v-col-lg-4 offset-lg-4">
          <v-card>
            <v-card-item>
              <v-card-title primary-title>
                Recherche d'arbres
              </v-card-title>
              <!--
              <v-card-subtitle>
                Sélection
              </v-card-subtitle>
              -->
            </v-card-item>
            <v-divider></v-divider>
            <v-card-text style="height: 300px;">
              <v-form @submit.prevent="submitForm">
                <v-container>

                  <v-row class="py-5">
                    <v-col cols="12" md="12">
                      <v-text-field label="Identifiant de l'arbre"></v-text-field>
                    </v-col>
                  </v-row>

                </v-container>
              </v-form>
            </v-card-text>
            <v-divider></v-divider>
            <v-card-actions>
              <v-btn color="info" @click="searchTreeOnOK">OK</v-btn>
              <v-btn color="info" @click="searchTreeOnCancel">Annuler</v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
    
  </div>
</template>

<style scoped>
.btn-treesearch-on {
  background-color: white;
  color: darkcyan;
}

.btn-treesearch-off {
  background-color: white;
  color: black;
}

.tree-search {
  position: fixed;
  z-index: 1000;
  top: 10em;
  left: 50%;
  -webkit-transform: translateX(-50%);
  -ms-transform: translateX(-50%);
  transform: translateX(-50%);
}
</style>