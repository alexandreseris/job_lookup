<script lang="ts" setup>
import DataTable from './DataTable.vue'
import { useCompanyTypeStore } from '../stores/company_type'

const store = useCompanyTypeStore()
</script>

<template>
    <data-table :store="store"> </data-table>
</template>