<script lang="ts" setup>
import DataTable from './DataTable.vue'
import { useEventSourceStore } from '../stores/event_source'


const store = useEventSourceStore()
</script>

<template>
    <data-table :store="store"> </data-table>
</template>