<script lang="ts" setup>
import DataTable from './DataTable.vue'
import { useEventStore } from '../stores/event';

const store = useEventStore()
</script>

<template>
    <data-table :store="store"> </data-table>
</template>