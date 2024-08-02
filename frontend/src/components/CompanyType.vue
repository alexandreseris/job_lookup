<script lang="ts" setup>
import { db } from '../../wailsjs/go/models'
import * as back from '../../wailsjs/go/main/App'
import DataTable from './DataTable.vue'
import * as types from './types'
import { useStore } from '../store'


type Item = db.CompanyType

const columns: types.Columns<Item> = [
    { title: 'Name', key: 'name', type: "string", requiered: true },
]

const emptyItem = {
    id: 0,
    name: '',
} as Item


const store = useStore()
</script>

<template>
    <data-table :items="store.companyTypes" :delete="back.DeleteCompanyType" :update="back.UpdateCompanyType"
        :insert="back.InsertCompanyType" :empty-item="emptyItem" :columns="columns">
    </data-table>
</template>