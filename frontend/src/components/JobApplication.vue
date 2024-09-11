<script lang="ts" setup>
import { main } from '../../wailsjs/go/models'
import * as back from '../../wailsjs/go/main/App'
import DataTable from './DataTable.vue'
import * as types from './types'
import { useStore } from '../store';


type Item = main.JobApplication

const columns: types.Columns<Item> = [
    { key: "company_name", title: "Company", type: "rel", requiered: true },
    { key: "job_title", title: "Title", type: "string", requiered: true },
    { key: "status_name", title: "Status", type: "rel", requiered: true },
    { key: "event_cnt", title: "Number of events", type: "int", readOnly: true },
    { key: "last_event", title: "Last event", type: "date", readOnly: true },
    { key: "next_event", title: "Next event", type: "date", readOnly: true },
    { key: "notes", title: "Notes", type: "multiline" },
]
const emptyItem = {
    company_id: 0,
    company_name: "",
    id: 0,
    job_title: "",
    notes: "",
    status_id: 0,
    status_name: "",
    event_cnt: 0,
    last_event: null,
    next_event: null,
} as Item

const store = useStore()
await store.init()
</script>

<template>
    <data-table :items="store.applications" :delete="back.DeleteJobApplication" :update="back.UpdateJobApplication"
        :insert="back.InsertJobApplication" :empty-item="emptyItem" :columns="columns" :relations="{
            company_name: store.findCompanyNamesFromApplication,
            status_name: store.findStatusNamesFromApplication,
        }">
    </data-table>
</template>