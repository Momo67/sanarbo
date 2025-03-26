<template>
  <div class="field">
    <div class="form-row">
    </div>
    <div class="form-row">
      <div class="form-group">
        <FloatLabel>
          <Select
            id="secteurs"
            v-model="secteurName"
            :options="gestion_com.secteurs.data"
            option-label="value"
            option-value="value"
            placeholder="Choisissez un secteur"
          />
          <label for="secteurs">Secteur</label>
        </FloatLabel>
      </div>
      <div class="form-group">
        <FloatLabel>
          <Select
            v-model="idEmplacement"
            label-id="emplacements"
            :options="gestion_com.emplacements.data"
            option-label="value"
            option-value="id"
            placeholder="Choisissez un emplacement"
          />
          <label for="emplacements">Emplacement</label>
        </FloatLabel>
      </div>
      <Button label="Rechercher" icon="pi pi-search" class="p-button-text" @click="getList" />
    </div>
  </div>
  <div class="card">
    <DataTable
      ref="dt"
      v-model:filters="filters"
      :value="dataTrees"
      removable-sort
      responsive-layout="scroll"
      striped-rows
      filter-display="menu"
      paginator
      :rows="10"
      :rows-per-page-options="[10,20,50,100]"
    >
      <template #header>
        <div class="flex justify-end">
          <Button type="button" icon="pi pi-filter-slash" label="Clear" outlined @click="clearFilter()" />
          <IconField>
            <InputIcon>
              <i class="pi pi-search" />
            </InputIcon>
            <InputText v-model="filters['global'].value" placeholder="Recherche mot-clé" />
          </IconField>
        </div>
      </template>
      <template #empty>
        <div class="p-text-center">
          <span class="p-text-bold">Aucun arbre à valider</span>
        </div>
      </template>
      <template #loading>
        <div class="p-text-center">
          <ProgressSpinner />
        </div>
      </template> 
      <ColumnGroup type="header">
        <Row>
          <Column field="id" header="id" :sortable="true" />
          <Column field="name" header="Nom" :sortable="true" />
          <Column field="last_modification_user" header="Modifié par" :sortable="true" />
          <Column field="last_modification_time" header="Date de modification" :sortable="true" />
          <Column header="Actions" />
        </Row>
      </ColumnGroup>
      <Column field="external_id" />
      <Column field="name" class="align-left" />
      <Column field="last_modification_user" class="align-left" />
      <Column field="last_modification_time" class="align-left" />
      <Column :exportable="false" style="min-width:8rem">
        <template #body="slotProps">
          <Button v-tooltip="'Voir'" icon="pi pi-eye" class="p-button-rounded  mr-2" @click="viewTree(slotProps.data.id)" />
          <Button 
                  v-tooltip="slotProps.data.is_validated ? 'Dévalider' : 'Valider'"
                  :icon="slotProps.data.is_validated ? 'pi pi-check' : 'pi pi-question'"
                  class="p-button-rounded mr-2" 
                  :class="slotProps.data.is_validated ? 'p-button-success' : 'p-button-danger'"
                  @click="toggleValidation(slotProps.index, slotProps.data)" />

        </template>
      </Column>
    </DataTable>
  </div>
  <div class="field">
    <div class="form-row">
      <div class="form-group button-container">
        <Button label="Sauver" icon="pi pi-check" class="p-button-success" @click="save" />
      </div>
    </div>
  </div>
  <Dialog 
    v-model:visible="treeDialog" 
    :modal="true" 
    :style="{ width: '50vw' }"
  >
    <template #header>
      <div class="dialog-header">
        <h3>Détails de l'arbre (ID: {{ dataTree.external_id }})</h3>
      </div>
    </template>
    <div class="dialog-content">
      <h3>{{ dataTree.name }}</h3>
      <p><strong>Description:</strong> {{ dataTree.description }}</p>
      <p><strong>Id Goéland:</strong> {{ dataTree.external_id }}</p>
      <p><strong>Date de création:</strong> {{ formatDate(dataTree.create_time) }}</p>
      <p><strong>Dernière modification:</strong> {{ formatDate(dataTree.last_modification_time) }}</p>
      
      <h4>📌 Attributs de l'arbre</h4>
      <ul>
        <li><strong>Circonférence:</strong> {{ dataTree.tree_attributes.circonference }} mm</li>
        <li><strong>Entourage:</strong> {{ dataTree.tree_attributes.entouragerem }}</li>
        <li><strong>État sanitaire:</strong> {{ dataTree.tree_attributes.etatsanitairerem }}</li>
        <li><strong>Public:</strong> {{ dataTree.tree_attributes.ispublic ? 'Oui' : 'Non' }}</li>
      </ul>
    </div>
    
    <!-- Bouton pour fermer -->
    <template #footer>
      <Button label="Fermer" icon="pi pi-times" class="p-button-text" @click="treeDialog = false" />
    </template>
  </Dialog>  
</template>

<script setup>
import { ref, watch, onMounted } from 'vue';
import { FilterMatchMode } from '@primevue/core/api'
import { useToast } from 'primevue/usetoast';
import Button from 'primevue/button';
import Column from 'primevue/column';
import ColumnGroup from 'primevue/columngroup';
import DataTable from 'primevue/datatable';
import InputIcon from 'primevue/inputicon';
import ProgressSpinner from 'primevue/progressspinner';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Select from 'primevue/select';
import FloatLabel from 'primevue/floatlabel';
import Row from 'primevue/row';
import tree from './Tree.js';
import { BACKEND_URL, apiRestrictedUrl, getLog } from '../config';
import { isNullOrUndefined } from '../tools/utils.js';
import { useFetch } from "../composables/FetchData.js";
import { getLocalJwtTokenAuth } from './Login.js';
import IconField from 'primevue/iconfield';

const backendUrl = `${BACKEND_URL}/${apiRestrictedUrl}/`;
const urlGestionCom = backendUrl + "gestion_com";

const moduleName = 'TreeValidation';
const timeToDisplayError = 7000;
const timeToDisplaySucces = 4000;
const toast = useToast();
const dt = ref();

const log = getLog(moduleName, 4, 2);
const dataTrees = ref(null);

const secteurName = ref('');
const idEmplacement = ref();

let secteurs = {data: []};
let emplacements = {data: []};
const gestion_com = ref({
  secteurs: secteurs,
  emplacements: emplacements,
});

const treeDialog = ref(false);
const dataTree = ref(null);

const filters = ref();
const initFilters = () => {
  filters.value = {
    global: { value: null, matchMode: FilterMatchMode.CONTAINS },
  };
};
initFilters();
const clearFilter = () => {
  initFilters();
  dt.value.filter(filters.value);
};

const emit = defineEmits(['user-invalid-session']);

// Get session storage token
//const token = sessionStorage.getItem('token');
const token = getLocalJwtTokenAuth();
const headers = {
  'Authorization': token,
  'Content-Type': 'application/json',
}

const options = {
  headers: headers
}

const defaultSecteur = {id: -1, value: 'Aucun'};

const getSecteurs = async () => {
  const method = 'getSecteurs';
  log.t(`##-->${moduleName}::${method}()`);
  secteurs = await useFetch(urlGestionCom + '/secteurs', options);
  gestion_com.value = {
    secteurs: secteurs,
    emplacements: {data: []},
  };
  gestion_com.value.secteurs.data.unshift(defaultSecteur);
};

const defaultEmplacement = {id: -1, value: 'Aucun'};

watch(secteurName, async () => {
  const method = 'watch(secteurName)';
  log.t(`##-->${moduleName}::${method}()`);
  if ((secteurName.value != defaultSecteur.value) && (secteurName.value != null)) {
    idEmplacement.value = null;
    
    emplacements = await useFetch(urlGestionCom + '/emplacements' + (secteurName.value != '' ? ('/' + secteurName.value) : ''), options);
    gestion_com.value = {
      secteurs: secteurs,
      emplacements: emplacements
    };
    gestion_com.value.emplacements.data.unshift(defaultEmplacement);
  } else {
    gestion_com.value = {
      secteurs: secteurs,
      emplacements: {data: [defaultEmplacement]},
    };
  }
});

const getList = () => {
  const method = 'getList';
  log.t(`##-->${moduleName}::${method}()`);
  const secteur = ((secteurName.value != null) && secteurName.value != defaultSecteur.value) ? secteurName.value : '';
  const epmplacement = (idEmplacement.value != null) ? idEmplacement.value : -1;
  tree.treesToValidate(secteur, epmplacement, (retval, statusMessage) => {
    if (statusMessage === 'SUCCESS') {
      dataTrees.value = retval;
      log.l(`# IN loadData -> dataTrees.value.length : ${dataTrees.value.length}`);
      if (dataTrees.value.length === 0) {
        toast.add({
          severity: 'info', summary: 'Info', detail: '⚠⚠ Aucun arbre à valider pour ce secteur et emplacement !', life: timeToDisplayError,
        })
      }
    } else {
      log.e(`# ERROR in getList tree.treesToValidate callback: ${statusMessage} \n error:`, retval);
      toast.add({
        severity: 'error', summary: 'Error', detail: `⚡⚡⚠ Unable to retrieve list of trees to validate for secteur: ${secteurName.value} and emplacement: ${idEmplacement.value} from DB ! error: ${statusMessage}`, life: timeToDisplayError,
      });
      checkNetworkError(retval);
      log.e(`# GOT ERROR calling tree.treesToValidate : ${statusMessage}, \n error:`, retval);
    }
  });
};

const viewTree = (id) => {
  const method = 'viewTree';
  log.t(`## IN ${method}`);
  log.l(`# IN viewTree -> id : ${id}`);
  tree.getTree(id, (retval, statusMessage) => {
    if (statusMessage === 'SUCCESS') {
      dataTree.value = retval;
      treeDialog.value = true;
    } else {
      log.e(`# ERROR in viewTree tree.getTree callback: ${statusMessage} \n error:`, retval);
      toast.add({
        severity: 'error', summary: 'Error', detail: `⚡⚡⚠ Unable to retrieve tree from DB ! error: ${statusMessage}`, life: timeToDisplayError,  // 7s
      });
      checkNetworkError(retval);
      log.e(`# GOT ERROR calling tree.getTree : ${statusMessage}, \n error:`, retval);
    }
  })
}

const toggleValidation = (index, data) => {
  const method = 'setValidation';
  log.t(`## IN ${method}`);
  log.l(`# IN setValidation -> index, idTree : ${data}, ${data.external_id}`);
  console.log('###index:', index);
  dataTrees.value[index].is_validated = !data.is_validated;
}

const save = () => {
  const method = 'save';
  log.t(`##-->${moduleName}::${method}()`);
  const validatedTrees = dataTrees.value.filter(tree => tree.is_validated === true);
  log.l(`# IN save -> validatedTrees.length : ${validatedTrees.length}`);
  const treesToSave = validatedTrees.map(tree => {
    return {
      external_id: tree.external_id,
      is_validated: tree.is_validated,
    };
  });
  if (treesToSave.length === 0) {
    toast.add({
      severity: 'info', summary: 'Info', detail: '⚠⚠ Aucune validation à sauver !', life: timeToDisplayError,
    });
    return;
  }
  tree.validateTrees(treesToSave, (retval, statusMessage) => {
    if (statusMessage === 'SUCCESS') {
      getList();
      toast.add({
        severity: 'success', summary: 'Success', detail: `✔✔ ${validatedTrees.length} trees have been validated successfully !`, life: timeToDisplaySucces,
      });
    } else {
      log.e(`# ERROR in save tree.saveValidatedTrees callback: ${statusMessage} \n error:`, retval);
      toast.add({
        severity: 'error', summary: 'Error', detail: `⚡⚡⚠ Unable to save validated trees ! error: ${statusMessage}`, life: timeToDisplayError,
      });
      checkNetworkError(retval);
      log.e(`# GOT ERROR calling tree.saveValidatedTrees : ${statusMessage}, \n error:`, retval);
    }
  });
}

function checkNetworkError(err) {
  if (!isNullOrUndefined(err.response)) {
    log.w('retval.response', err.response, err.response.status, (err.response.status === 401));
    if (err.response.status === 401) {
      emit('user-invalid-session', 'User session is invalid', err.response);
    }
  }
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString("fr-FR", {
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit"
  });
}

onMounted(async () => {
  const method = 'onMounted';
  log.t(`##-->${moduleName}::${method}`);
  getSecteurs();
});

</script>

<style scoped>
.field {
  display: flex;
  flex-direction: column;
  gap: 1rem; /* Espacement constant entre chaque élément */
}

.form-row {
  display: flex;
  gap: 2rem; /* Espacement entre les deux groupes */
  align-items: flex-end; /* Aligner proprement avec le plus grand élément */
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem; /* Espacement entre label et Select */
}

.button-container {
  display: flex;
  margin-left: auto; /* Aligner à droite */
  margin-top: 1rem; /* Espacement du haut */
}
</style>