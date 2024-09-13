<script lang="ts" setup>
import DataTable from './DataTable.vue'
import { useContactStore } from '../stores/contact';


const store = useContactStore()
</script>

<template>
    <data-table :store="store"> </data-table>
</template>