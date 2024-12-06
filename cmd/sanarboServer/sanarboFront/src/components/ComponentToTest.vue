// ComponentToTest.vue
<template>
  <v-btn @click="handleClick">Click me</v-btn>
</template>

<script setup>
const emit = defineEmits(['button-clicked']);

const handleClick = () => {
  emit('button-clicked');
}

</script>
