<script setup lang="ts">
import router from '@/router';
import { computed, ref } from 'vue';



const query = ref("e")


const url = computed(() => {
  return `http://localhost:8080/api/v1/search/?q=${encodeURIComponent(props.query)}`;
})

function search() {
  router.push({ name: 'search', query: { q: query.value }})
}
</script>

<template>
  <main class="flex flex-row justify-center items-center h-screen w-full">
    <div class="flex flex-col items-center gap-8 w-full m-4">
    <div class="text-6xl">Ask Dave</div>
    <input type="text" class="w-full bg-slate-300 rounded-full text-slate-900 px-4 py-2 max-w-140 h-10" @keyup.enter="search" v-model="query">
  </div>
    
  </main>
</template>
