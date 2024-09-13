<script lang="ts" setup>
import DataTable from './DataTable.vue'
import { useJobApplicationStatusStore } from '../stores/job_application_status'


const store = useJobApplicationStatusStore()
</script>

<template>
    <data-table :store="store"> </data-table>
</template>