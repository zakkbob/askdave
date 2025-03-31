<script setup lang="ts">
import Result from '@/components/Result.vue';
import { useFetch } from '@vueuse/core'
import { computed, ref } from 'vue';
import router from '@/router';

const props = defineProps({
    query: {
      type: String,
      required: true
    },
})

const query = ref("e")
query.value = props.query


const url = computed(() => {
  return `http://localhost:8080/api/v1/search/?q=${encodeURIComponent(props.query)}`;
})

const { data } = useFetch<String>(url, { refetch: true })

const results = computed(() => {
  return JSON.parse(data.value)
})

function search() {
  router.push({ name: 'search', query: { q: query.value }})
}

</script>

<template>
    <main class="flex flex-col px-12 py-8 items-left h-fit w-full">
      <div class="w-full h-fit mb-12"><input type="text" class="w-full bg-slate-300 rounded-full text-slate-900 px-4 py-2 max-w-140 h-10" @keyup.enter="search" v-model="query"></div>
      <div class="ml-2 flex flex-col gap-4 max-w-240">
        <Result v-for="result in results" :url="result.url" :description="result.description" :site-name="result.site_name" :title="result.title">
        </Result>
      </div>
    </main>
</template>