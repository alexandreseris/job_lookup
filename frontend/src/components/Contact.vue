<script lang="ts" setup>
import { main } from '../../wailsjs/go/models'
import * as back from '../../wailsjs/go/main/App'
import DataTable from './DataTable.vue'
import * as types from './types'
import { useStore } from '../store';


type Item = main.Contact

const columns: types.Columns<Item> = [
    { key: 'company_name', title: "Company", type: "rel", requiered: true },
    { key: 'fist_name', title: "Fist name", type: "string", requiered: true },
    { key: 'last_name', title: "Last name", type: "string", requiered: true },
    { key: 'job_position', title: "Position", type: "string", requiered: true },
    { key: 'email', title: "Email", type: "string" },
    { key: 'phone_number', title: "Phone", type: "string" },
    { key: "last_event", title: "Last event", type: "date", readOnly: true },
    { key: "next_event", title: "Next event", type: "date", readOnly: true },
    { key: 'notes', title: "Notes", type: "multiline" },
]

const emptyItem = {
    company_id: 0,
    company_name: "",
    email: "",
    fist_name: "",
    id: 0,
    job_position: "",
    last_name: "",
    notes: "",
    phone_number: "",
    last_event: new Date(),
    next_event: new Date(),
} as Item

const store = useStore()
await store.init()
</script>

<template>
    <data-table :items="store.contacts" :delete="back.DeleteContact" :update="back.UpdateContact"
        :insert="back.InsertContact" :empty-item="emptyItem" :columns="columns"
        :relations="{ company_name: store.findCompanyNamesFromContact }">
    </data-table>
</template>