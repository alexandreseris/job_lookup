<script lang="ts" setup>
import DataTable from './DataTable.vue'

import { useJobApplicationStore } from '../stores/job_application';

const store = useJobApplicationStore()
</script>

<template>
    <data-table :store="store"> </data-table>
</template>