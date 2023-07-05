<script setup>
import {onBeforeMount, onMounted, reactive, ref} from 'vue';
import {useFetch} from "../composables/FetchData.js";


const backendUrl = import.meta.env.VITE_BACKEND_API_URL;
const urlTrees = backendUrl + "trees";



const emit = defineEmits(['formSubmitted', 'formCanceled'])
const props = defineProps({
  showForm: Boolean,
  treeId: Number,
})


const Tree = reactive({
  create_time: '',
  creator: '',
  description: '',
  id: '',
  name: '',
  geom: '',
  tree_attributes: {},
});

const Dict = ref({
  "validation" : {},
  "to_be_checked": {},
  "note": {},
  "check": {},
  "entourage" : {},
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

onBeforeMount(async () => {

})



onMounted(async () => {

  const dict_validation = await useFetch(backendUrl + 'dico/validation', options);
  const dict_to_be_checked = await useFetch(backendUrl + 'dico/to_be_checked', options);
  const dict_note = await useFetch(backendUrl + 'dico/note', options);
  const dict_entourage = await useFetch(backendUrl + 'dico/entourage', options);
  const dict_check = await useFetch(backendUrl + 'dico/check', options);
  const rev_surface = await useFetch(backendUrl + 'dico/rev_surface', options);
  const rev_etat_sanitaire = await useFetch(backendUrl + 'dico/etat_sanitaire', options);
  const rev_etat_sanitaire_rem = await useFetch(backendUrl + 'dico/etat_sanitaire', options);


  Dict.value = {
    "validation" : dict_validation,
    "to_be_checked": dict_to_be_checked,
    "note": dict_note,
    "check": dict_check,
    "entourage" : dict_entourage,
    "rev_surface": rev_surface,
    "etat_sanitaire": rev_etat_sanitaire,
    "etat_sanitaire_rem": rev_etat_sanitaire_rem
  }

  console.log(Dict.value.validation.data)

  const tree = await useFetch(urlTrees + '/' + props.treeId, options)
  Tree.create_time = tree.data.value.create_time;
  Tree.creator = tree.data.value.creator;
  Tree.description = tree.data.value.description;
  Tree.id = tree.data.value.id;
  Tree.name = tree.data.value.name;
  Tree.tree_attributes = tree.data.value.tree_attributes;
  Tree.geom = tree.data.value.geom;
})


const submitForm = async (event) => {

  const options = {
    headers: headers,
    method: 'PUT',
    body: JSON.stringify(Tree)
  }

  await useFetch(urlTrees + '/' + props.treeId, options)

  // Emit a custom event to notify the parent component
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
        <h2>Arbre - {{Tree.name}}</h2>
        <v-row class="py-5">

              <v-col cols="12" md="12">
                <v-select
                    label="À contrôler"
                    :items="Dict.to_be_checked.data"
                    v-model.number="Tree.tree_attributes.idtobechecked"
                    item-title="value"
                    item-value="id"
                >
                  <template v-slot:prepend>
                    <p>À contrôler</p>
                  </template>
                </v-select>




              </v-col>
              <v-col cols="12" md="12">

                <v-select
                    label="Statut"
                    :items="Dict.validation.data"
                    v-model.number="Tree.tree_attributes.idvalidation"
                    item-title="value"
                    item-value="id"
                >
                  <template v-slot:prepend>
                    <p>Statut</p>
                  </template>
                </v-select>



              </v-col>

              <v-col cols="12" md="12">

                <v-select
                    label="Note"
                    :items="Dict.note.data"
                    v-model.number="Tree.tree_attributes.idnote"
                    item-title="value"
                    item-value="id"
                >
                  <template v-slot:prepend>
                    <p>Note</p>
                  </template>
                </v-select>


              </v-col>
            </v-row>
      </v-container>




      <v-container>
        <h2 >Environnement</h2>
        <v-row class="py-5">
          <v-col cols="12" md="12">
            <v-text-field
                v-model.number="Tree.tree_attributes.circonference"
                label="Circonférence"
                type="number"
            >
              <template v-slot:prepend>
                <p>Circonférence [cm]</p>
              </template>

            </v-text-field>
          </v-col>
          <v-col cols="12" md="4">



            <v-select
                label="Entourage / Cadre"
                :items="Dict.entourage.data"
                v-model.number="Tree.tree_attributes.identourage"
                item-title="value"
                item-value="id"
            >
              <template v-slot:prepend>
                <p>Entourage / Cadre</p>
              </template>
            </v-select>



          </v-col>
          <v-col cols="12" md="4">



            <v-select
                label="Statut"
                :items="Dict.check.data"
                v-model.number="Tree.tree_attributes.idchkentourage"
                item-title="value"
                item-value="id"
            >
            </v-select>


          </v-col>





          <v-col cols="12" md="4">

          <v-text-field
              v-model="Tree.tree_attributes.entouragerem"
              label="Remarque entourage"
              type="string"
          ></v-text-field>


          </v-col>


        <v-col cols="12" md="4">

          <v-select
              label="Revêtement de surface"
              :items="Dict.rev_surface.data"
              v-model.number="Tree.tree_attributes.idrevsurface"
              item-title="value"
              item-value="id"
          >
            <template v-slot:prepend>
              <p>Revêtement de surface</p>
            </template>
          </v-select>
        </v-col>
        <v-col cols="12" md="4">
          <v-select
              label="Statut"
              :items="Dict.check.data"
              v-model.number="Tree.tree_attributes.idchkrevsurface"
              item-title="value"
              item-value="id"
          >
          </v-select>
        </v-col>

        <v-col cols="12" md="4">

          <v-text-field
              v-model="Tree.tree_attributes.revsurfacerem"
              label="Remarque revêtement"
              type="string"
          ></v-text-field>
        </v-col>

        </v-row>

        <h2 class="pt-10">État sanitaire</h2>
        <v-row class="py-5">
          <v-col cols="12" md="4">


            <v-select
                label="Statut"
                :items="Dict.etat_sanitaire.data"
                v-model.number="Tree.tree_attributes.idetatsanitairepied"
                item-title="value"
                item-value="id"
            >

              <template v-slot:prepend>
                <p>Pied</p>
              </template>
            </v-select>


          </v-col>


          <v-col cols="12" md="4">
            <v-select
                label="Statut"
                :items="Dict.etat_sanitaire.data"
                v-model.number="Tree.tree_attributes.idetatsanitairetronc"
                item-title="value"
                item-value="id"
            >

              <template v-slot:prepend>
                <p>Tronc</p>
              </template>
            </v-select>
        </v-col>
          <v-col cols="12" md="4">
            <v-select
                label="Statut"
                :items="Dict.etat_sanitaire.data"
                v-model.number="Tree.tree_attributes.idetatsanitairecouronne"
                item-title="value"
                item-value="id"
            >

              <template v-slot:prepend>
                <p>Couronne</p>
              </template>
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
            <v-btn type="submit" color="primary" @click="submitForm">Sauver</v-btn>
          </v-col>
          <v-col cols="12" md="2">
            <v-btn type="button" color="secondary" @click="handleFormCanceled">Annuler</v-btn>
          </v-col>
        </v-row>

      </v-container>
    </v-form>

  </div>
</template>
