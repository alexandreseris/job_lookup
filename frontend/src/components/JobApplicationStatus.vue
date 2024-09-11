<script lang="ts" setup>
import { db } from '../../wailsjs/go/models'
import * as back from '../../wailsjs/go/main/App'
import DataTable from './DataTable.vue'
import * as types from './types'
import { useStore } from '../store'


type Item = db.ListJobApplicationStatusRow

const columns: types.Columns<Item> = [
    { title: 'Name', key: 'name', type: "string", requiered: true },
    { title: 'Jobs', key: 'applications', type: "int", readOnly: true },
]

const emptyItem = {
    id: 0,
    name: '',
    applications: 0
} as Item


const store = useStore()
await store.init()
</script>

<template>
    <data-table :items="store.applicationStatus" :delete="back.DeleteJobApplicationStatus"
        :update="back.UpdateJobApplicationStatus" :insert="back.InsertJobApplicationStatus" :empty-item="emptyItem"
        :columns="columns">
    </data-table>
</template>