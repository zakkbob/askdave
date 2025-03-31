<script setup lang="ts">
const props = defineProps({
    siteName: String,
    title: String,
    description: String,
    url: String,
})
</script>

<template>
<a class="bg-zinc-700 w-full rounded-lg p-4" :href="url">
    
    <div class="flex flex-col mb-1"><span>{{ siteName }}</span> <span class="text-xs text-zinc-200">{{ url }}</span></div>

    <div class="text-blue-300 text-lg">{{ title }}</div>
    <div class="text-sm">{{ description }}</div>
</a>
</template>