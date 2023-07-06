<script setup>
import {onMounted, ref} from 'vue';
import {useFetch} from "../composables/FetchData.js";


const backendUrl = import.meta.env.VITE_BACKEND_API_URL;
const urlTrees = backendUrl + "trees";


const emit = defineEmits(['formSubmitted', 'formCanceled'])
const props = defineProps({
  showForm: {type: Boolean, required: false, default: false},
  treeId: {type: String, required: false, default: ''},
})


const Tree = ref({
  create_time: '',
  creator: '',
  description: '',
  id: '',
  name: '',
  geom: '',
  tree_attributes: {},
});


const Dict = ref({
  "validation": {},
  "to_be_checked": {},
  "note": {},
  "check": {},
  "entourage": {},
  "rev_surface": {},
  "etat_sanitaire": {},
  "etat_sanitaire_rem": {}
})


// Get session storage token
const token = sessionStorage.getItem('token');
const headers = {
  'Authorization': 'Bearer ' + token,
  'Content-Type': 'application/json',
}

const options = {
  headers: headers
}


onMounted(async () => {

  const dict_validation = await useFetch(backendUrl + 'dico/validation', options);
  const dict_to_be_checked = await useFetch(backendUrl + 'dico/to_be_checked', options);
  const dict_note = await useFetch(backendUrl + 'dico/note', options);
  const dict_entourage = await useFetch(backendUrl + 'dico/entourage', options);
  const dict_check = await useFetch(backendUrl + 'dico/check', options);
  const rev_surface = await useFetch(backendUrl + 'dico/rev_surface', options);
  const etat_sanitaire = await useFetch(backendUrl + 'dico/etat_sanitaire', options);
  const etat_sanitaire_rem = await useFetch(backendUrl + 'dico/etat_sanitaire_rem', options);


  Dict.value = {
    "validation": dict_validation,
    "to_be_checked": dict_to_be_checked,
    "note": dict_note,
    "check": dict_check,
    "entourage": dict_entourage,
    "rev_surface": rev_surface,
    "etat_sanitaire": etat_sanitaire,
    "etat_sanitaire_rem": etat_sanitaire_rem
  }


  const tree = await useFetch(urlTrees + '/' + props.treeId, options)
  Tree.value.create_time = tree.data.value.create_time;
  Tree.value.creator = tree.data.value.creator;
  Tree.value.description = tree.data.value.description;
  Tree.value.id = tree.data.value.id;
  Tree.value.name = tree.data.value.name;
  Tree.value.tree_attributes = tree.data.value.tree_attributes;
  Tree.value.geom = tree.data.value.geom;
})


// eslint-disable-next-line no-unused-vars
const submitForm = async () => {

  const options = {
    headers: headers,
    method: 'PUT',
    body: JSON.stringify(Tree)
  }

  await useFetch(urlTrees + '/' + props.treeId, options)

  emit('formSubmitted');
};


const handleFormCanceled = () => {
  emit('formCanceled')
}


</script>


<template>
  <div>
    <v-form @submit.prevent="submitForm">
      <v-container>
        <h2> Arbre - {{ Tree.name }}</h2>
        <v-row class="py-5">
          <v-col cols="12" md="12">
            <v-select
                v-model.number="Tree.tree_attributes.idtobechecked"
                :items="Dict.to_be_checked.data"
                item-title="value"
                item-value="id"
                label="À contrôler"
            >
            </v-select>
          </v-col>
          <v-col cols="12" md="12">
            <v-select
                v-model.number="Tree.tree_attributes.idvalidation"
                :items="Dict.validation.data"
                item-title="value"
                item-value="id"
                label="Statut"
            >
            </v-select>
          </v-col>
          <v-col cols="12" md="12">
            <v-select
                v-model.number="Tree.tree_attributes.idnote"
                :items="Dict.note.data"
                item-title="value"
                item-value="id"
                label="Note"
            >
            </v-select>
          </v-col>
        </v-row>
      </v-container>

      <v-container>
        <h2>Environnement</h2>
        <v-row class="py-5">
          <v-col cols="12" md="12">
            <v-text-field
                v-model.number="Tree.tree_attributes.circonference"
                label="Circonférence [cm]"
                type="number"
            >
            </v-text-field>
          </v-col>

          <v-card
              class="mx-auto"
              prepend-icon="mdi-circle-box-outline"
              width="100%"
          >
            <template #title>
              Entourage / Cadre
            </template>
            <v-col cols="12" md="12">
              <v-select
                  v-model.number="Tree.tree_attributes.identourage"
                  :items="Dict.entourage.data"
                  item-title="value"
                  item-value="id"
                  label="Type"
              >
              </v-select>
            </v-col>
            <v-col cols="12" md="12">
              <v-select
                  v-model.number="Tree.tree_attributes.idchkentourage"
                  :items="Dict.check.data"
                  item-title="value"
                  item-value="id"
                  label="Statut"
              >
              </v-select>
            </v-col>
            <v-col cols="12" md="12">
              <v-text-field
                  v-model="Tree.tree_attributes.entouragerem"
                  label="Remarque entourage"
                  type="string"
              ></v-text-field>
            </v-col>
          </v-card>


          <v-card
              class="mx-auto"
              prepend-icon="mdi-texture-box"
              width="100%"
          >
            <template #title>
              Revêtement
            </template>
            <v-col cols="12" md="12">
              <v-select
                  v-model.number="Tree.tree_attributes.idrevsurface"
                  :items="Dict.rev_surface.data"
                  item-title="value"
                  item-value="id"
                  label="Type"
              >
              </v-select>
            </v-col>
            <v-col cols="12" md="12">
              <v-select
                  v-model.number="Tree.tree_attributes.idchkrevsurface"
                  :items="Dict.check.data"
                  item-title="value"
                  item-value="id"
                  label="Statut"
              >
              </v-select>
            </v-col>
            <v-col cols="12" md="12">
              <v-text-field
                  v-model="Tree.tree_attributes.revsurfacerem"
                  label="Remarque revêtement"
                  type="string"
              ></v-text-field>
            </v-col>
          </v-card>
        </v-row>


        <h2 class="pt-10">État sanitaire</h2>
        <v-row class="py-5">
          <v-col cols="12" md="12">
            <v-select
                v-model.number="Tree.tree_attributes.idetatsanitairepied"
                :items="Dict.etat_sanitaire.data"
                item-title="value"
                item-value="id"
                label="Pied"
            >
            </v-select>
          </v-col>
          <v-col cols="12" md="12">
            <v-select
                v-model.number="Tree.tree_attributes.idetatsanitairetronc"
                :items="Dict.etat_sanitaire.data"
                item-title="value"
                item-value="id"
                label="Tronc"
            >
            </v-select>
          </v-col>
          <v-col cols="12" md="12">
            <v-select
                v-model.number="Tree.tree_attributes.idetatsanitairecouronne"
                :items="Dict.etat_sanitaire.data"
                item-title="value"
                item-value="id"
                label="Couronne"
            >
            </v-select>
          </v-col>
          <v-col cols="12" md="12">
            <v-text-field
                v-model="Tree.tree_attributes.etatsanitairerem"
                label="Remarque état sanitaire">
            </v-text-field>
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12" md="2">
            <v-btn color="primary" type="submit" @click="submitForm">Sauver</v-btn>
          </v-col>
          <v-col cols="12" md="2">
            <v-btn color="secondary" type="button" @click="handleFormCanceled">Annuler</v-btn>
          </v-col>
        </v-row>
      </v-container>
    </v-form>
  </div>
</template>


<style scoped>

</style>
