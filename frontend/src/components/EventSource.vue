<script lang="ts" setup>
import { db } from '../../wailsjs/go/models'
import * as back from '../../wailsjs/go/main/App'
import DataTable from './DataTable.vue'
import * as types from './types'
import { useStore } from '../store'


type Item = db.ListEventSourceRow

const columns: types.Columns<Item> = [
    { title: 'Name', key: 'name', type: "string", requiered: true },
    { title: 'Events', key: 'events', type: "int", readOnly: true },
]

const emptyItem = {
    id: 0,
    name: '',
    events: 0
} as Item


const store = useStore()
await store.init()
</script>

<template>
    <data-table :items="store.eventSource" :delete="back.DeleteEventSource" :update="back.UpdateEventSource"
        :insert="back.InsertEventSource" :empty-item="emptyItem" :columns="columns">
    </data-table>
</template>