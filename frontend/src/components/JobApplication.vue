<script lang="ts" setup>
import { db } from '../../wailsjs/go/models'
import * as back from '../../wailsjs/go/main/App'
import DataTable from './DataTable.vue'
import * as types from './types'
import { useStore } from '../store';


type Item = db.ListJobApplicationRow

const columns: types.Columns<Item> = [
    { key: "company_name", title: "Company", type: "rel", requiered: true },
    { key: "job_title", title: "Title", type: "string", requiered: true },
    { key: "status_name", title: "Status", type: "rel", requiered: true },
    { key: "notes", title: "Notes", type: "string" },
]
const emptyItem = {
    company_id: 0,
    company_name: "",
    id: 0,
    job_title: "",
    notes: "",
    status_id: 0,
    status_name: "",
} as Item

const store = useStore()
await store.init()
</script>

<template>
    <data-table :items="store.applications" :delete="back.DeleteJobApplication" :update="back.UpdateJobApplication"
        :insert="back.InsertJobApplication" :empty-item="emptyItem" :columns="columns" :relations="{
            company_name: store.companyNames,
            status_name: store.applicationStatusNames,
        }">
    </data-table>
</template>