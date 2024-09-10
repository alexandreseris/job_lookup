<script lang="ts" setup>
import { main } from '../../wailsjs/go/models'
import * as back from '../../wailsjs/go/main/App'
import DataTable from './DataTable.vue'
import * as types from './types'
import { useStore } from '../store';
import * as utils from './utils'


type Item = main.Event

const columns: types.Columns<Item> = [
    { key: "company_name", title: "Company", type: "rel", requiered: true },
    { key: "job_title", title: "Job", type: "string", requiered: true },
    { key: "title", title: "Title", type: "string", requiered: true },
    { key: "date", title: "Date", type: "date", requiered: true },
    { key: "source", title: "Source", type: "rel", requiered: true },
    { key: "contacts", title: "Contacts", type: "listrel" },
    { key: "notes", title: "Notes", type: "multiline" },

]

const emptyItem = {
    company_name: "",
    contacts: [],
    date: new Date(),
    id: 0,
    job_title: "",
    notes: "",
    source: "",
    title: "",
    convertValues(a: any, classs: any, asMap: boolean = false): any { }
} as Item

const store = useStore()
await store.init()

async function insertEvent(item: Item) {
    let newitem = await back.InsertEvent(item)
    newitem.date = utils.parseBackendDate(newitem.date)
    return newitem
}

</script>

<template>
    <data-table :items="store.events" :delete="back.DeleteEvent" :update="back.UpdateEvent" :insert="insertEvent"
        :empty-item="emptyItem" :columns="columns" :relations="{
            company_name: store.findCompanyNamesFromEvent,
            source: store.findSourceNamesFromEvent,
        }" :list-relations="{
            contacts: store.findContactNamesFromEvent
        }">
    </data-table>
</template>