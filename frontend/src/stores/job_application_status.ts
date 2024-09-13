import { loadData } from './utils'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as back from '../../wailsjs/go/main/App'
import { db, main } from '../../wailsjs/go/models'
import * as types from '../components/types'
import { useJobApplicationStore } from './job_application'

type Item = db.ListJobApplicationStatusRow
const ItemClass = db.ListJobApplicationStatusRow


export const useJobApplicationStatusStore = defineStore("JobApplicationStatus", () => {
    const items = ref<Item[]>([])

    function getColumns() {
        const columns: types.Columns<Item> = [
            { title: 'Name', key: 'name', type: "string", requiered: true },
            { title: 'Jobs', key: 'applications', type: "int", readOnly: true },
        ]
        return columns
    }
    async function select() {
        return await back.ListJobApplicationStatus()
    }

    async function syncItems() {
        await loadData(async () => {
            items.value = await select()
        }, ItemClass)
    }
    async function syncWithChildrens() {
        const jobApplicationStore = useJobApplicationStore()
        await Promise.all([syncItems(), jobApplicationStore.syncItems()])
    }
    async function syncWithParents() {
        await syncItems()
    }

    function add(): void {
        items.value = [Object.assign({}, {
            id: 0,
            name: '',
            applications: 0
        }), ...items.value]
    }
    function eq(i1: Item, i2: Item): boolean {
        return i1.name === i2.name
    }
    async function insert(item: Item): Promise<Item> {
        return await back.InsertJobApplicationStatus(item)
    }
    async function delete_(item: Item): Promise<void> {
        return await back.DeleteJobApplicationStatus(item)
    }
    async function update(item: Item): Promise<void> {
        return await back.UpdateJobApplicationStatus(item)
    }

    function findNamesFromApplication(application: main.JobApplication): string[] {
        return items.value
            .map((e) => { return e.name })
    }

    return {
        items, getColumns, syncItems, syncWithChildrens, syncWithParents, add, eq, select, insert, delete_, update, findNamesFromApplication
    }
})
