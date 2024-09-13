import { loadData } from './utils'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as back from '../../wailsjs/go/main/App'
import { db, main } from '../../wailsjs/go/models'
import * as types from '../components/types'
import { useCompanyStore } from './company'

type Item = db.ListCompanyTypeRow
const ItemClass = db.ListCompanyTypeRow

export const useCompanyTypeStore = defineStore("CompanyType", () => {
    const companyStore = useCompanyStore()

    const items = ref<Item[]>([])

    const columns: types.Columns<Item> = [
        { title: 'Name', key: 'name', type: "string", requiered: true },
        { title: 'Companies', key: 'companies', type: "int", readOnly: true },
    ]

    async function select() {
        return await back.ListCompanyTypes()
    }

    async function syncItems() {
        await loadData(async () => {
            items.value = await select()
        }, ItemClass)
    }
    async function syncWithChildrens() {
        await Promise.all([syncItems(), companyStore.syncItems()])

    }
    async function syncWithParents() {
        await syncItems()
    }

    function add(): void {
        items.value = [Object.assign({}, {
            id: 0,
            name: '',
            companies: 0
        }), ...items.value]
    }
    function eq(i1: Item, i2: Item): boolean {
        return i1.name === i2.name
    }
    async function insert(item: Item): Promise<Item> {
        return await back.InsertCompanyType(item)
    }
    async function delete_(item: Item): Promise<void> {
        return await back.DeleteCompanyType(item)
    }
    async function update(item: Item): Promise<void> {
        return await back.UpdateCompanyType(item)
    }

    function findNamesFromCompany(company: main.Company): string[] {
        return items.value
            .map((e) => { return e.name })
    }


    return {
        items, columns, syncItems, syncWithChildrens, syncWithParents, add, eq, select, insert, delete_, update, findNamesFromCompany
    }
})
