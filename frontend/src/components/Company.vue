<script lang="ts" setup>
import { main } from '../../wailsjs/go/models'
import * as back from '../../wailsjs/go/main/App'
import DataTable from './DataTable.vue'
import * as types from './types'
import { useStore } from '../store';

type Item = main.Company

const columns: types.Columns<Item> = [
    { key: 'name', title: 'Name', type: "string", requiered: true },
    { key: 'company_types', title: 'Types', type: "listrel", requiered: true },
    { key: 'application_cnt', title: 'Number of jobs', type: "int", readOnly: true },
    { key: "last_event", title: "Last event", type: "date", readOnly: true },
    { key: "next_event", title: "Next event", type: "date", readOnly: true },
    { key: 'notes', title: 'Notes', type: "multiline" },
]

const emptyItem = {
    id: 0,
    name: '',
    notes: '',
    company_types: [],
    application_cnt: 0,
    last_event: null,
    next_event: null,
    convertValues(a: any, classs: any, asMap: boolean = false): any { },
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