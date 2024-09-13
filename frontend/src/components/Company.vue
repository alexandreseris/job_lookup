<script lang="ts" setup>
import DataTable from './DataTable.vue'
import { useCompanyStore } from '../stores/company';

const store = useCompanyStore()
</script>

<template>
    <data-table :store="store"></data-table>
</template>