<script lang="ts" setup>
import { main } from '../../wailsjs/go/models'
import * as back from '../../wailsjs/go/main/App'
import DataTable from './DataTable.vue'
import * as types from './types'
import { useStore } from '../store';

type Item = main.Company

const columns: types.Columns<Item> = [
    { title: 'Name', key: 'name', type: "string", requiered: true },
    { title: 'Types', key: 'company_types', type: "listrel", requiered: true },
    { title: 'Notes', key: 'notes', type: "multiline" },
]

const emptyItem = {
    id: 0,
    name: '',
    notes: '',
    company_types: [],
} as Item


const store = useStore()
await store.init()
</script>

<template>
    <data-table :items="store.companies" :delete="back.DeleteCompany" :update="back.UpdateCompany"
        :insert="back.InsertCompany" :empty-item="emptyItem" :columns="columns"
        :list-relations="{ 'company_types': store.findCompanyTypeNamesFromCompany }">
    </data-table>
</template>