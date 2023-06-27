<script setup>
import {onMounted, reactive, ref} from 'vue';
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

  const {data} = await useFetch(urlTrees + '/' + props.treeId, options)
  Tree.create_time = data.value.create_time;
  Tree.creator = data.value.creator;
  Tree.description = data.value.description;
  Tree.id = data.value.id;
  Tree.name = data.value.name;
  Tree.tree_attributes = data.value.tree_attributes;
  Tree.geom = data.value.geom;
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
                <v-text-field
                    v-model.number="Tree.tree_attributes.idtobechecked"
                    label="À contrôler"
                    type="number"
                >

                  <template v-slot:prepend>
                    <p>À contrôler</p>
                  </template>

                </v-text-field>
              </v-col>
              <v-col cols="12" md="12">
                <v-text-field
                    v-model.number="Tree.tree_attributes.idvalidation"
                    label="Statut"
                    type="number"
                >

                  <template v-slot:prepend>
                    <p>Statut</p>
                  </template>


                </v-text-field>
              </v-col>

              <v-col cols="12" md="12">
                <v-text-field
                    v-model.number="Tree.tree_attributes.idnote"
                    label="Note"
                    type="number"
                >

                  <template v-slot:prepend>
                    <p>Note</p>
                  </template>


                </v-text-field>
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
                <p>Circonférence</p>
              </template>

            </v-text-field>
          </v-col>
          <v-col cols="12" md="4">

            <v-text-field
                v-model.number="Tree.tree_attributes.identourage"
                label="Entourage / cadre"
                type="number"
            >

              <template v-slot:prepend>
                <p>Entourage / cadre</p>
              </template>

            </v-text-field>
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field
                v-model.number="Tree.tree_attributes.idchkentourage"
                label="ID chk entourage"
                type="number"
            ></v-text-field>

          </v-col>

          <v-col cols="12" md="4">

          <v-text-field
              v-model="Tree.tree_attributes.entouragerem"
              label="Remarque entourage"
              type="string"
          ></v-text-field>
          </v-col>


        <v-col cols="12" md="4">
          <v-text-field
              v-model.number="Tree.tree_attributes.idrevsurface"
              label="Revêtement de surface"
              type="number"
          >
            <template v-slot:prepend>
              <p>Revêtement de surface</p>
            </template>

          </v-text-field>
        </v-col>
        <v-col cols="12" md="4">
          <v-text-field
              v-model.number="Tree.tree_attributes.idchkrevsurface"
              label="ID chk revsurface"
              type="number"
          ></v-text-field>

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
            <v-text-field
                v-model.number="Tree.tree_attributes.idetatsanitairepied"
                label="Pied">


              <template v-slot:prepend>
                <p>Pied</p>
              </template>


            </v-text-field>
          </v-col>






          <v-col cols="12" md="4">
          <v-text-field
              v-model.number="Tree.tree_attributes.idetatsanitairetronc"
              label="Tronc">
            <template v-slot:prepend>
              <p>Tronc</p>
            </template>

          </v-text-field>
        </v-col>
          <v-col cols="12" md="4">
            <v-text-field
                v-model.number="Tree.tree_attributes.idetatsanitairecouronne"
                label="Couronne">

              <template v-slot:prepend>
                <p>Couronne</p>
              </template>

            </v-text-field>
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
