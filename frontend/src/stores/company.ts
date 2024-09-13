import { loadData } from './utils'
import * as utils from '../components/utils'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as back from '../../wailsjs/go/main/App'
import { main } from '../../wailsjs/go/models'
import * as types from '../components/types'
import { useCompanyTypeStore } from './company_type'
import { useContactStore } from './contact'
import { useJobApplicationStore } from './job_application'
import { useEventStore } from './event'

type Item = main.Company
const ItemClass = main.Company

export const useCompanyStore = defineStore("Company", () => {
    const companyTypeStore = useCompanyTypeStore()
    const contactStore = useContactStore()
    const jobApplicationStore = useJobApplicationStore()
    const eventStore = useEventStore()

    const items = ref<Item[]>([])

    const columns: types.Columns<Item> = [
        { key: 'name', title: 'Name', type: "string", requiered: true },
        { key: 'company_types', title: 'Types', type: "listrel", requiered: true, relations: companyTypeStore.findNamesFromCompany },
        { key: 'application_cnt', title: 'Number of jobs', type: "int", readOnly: true },
        { key: "last_event", title: "Last event", type: "date", readOnly: true },
        { key: "next_event", title: "Next event", type: "date", readOnly: true },
        { key: 'notes', title: 'Notes', type: "multiline" },
    ]

    async function select() {
        let companiesConv = await back.ListCompanies()
        for (const e of companiesConv) {
            e.last_event = utils.parseBackendDateOpt(e.last_event)
            e.next_event = utils.parseBackendDateOpt(e.next_event)
        }
        return companiesConv

    }

    async function syncItems() {
        await loadData(async () => {
            items.value = await select()
        }, ItemClass)

    }
    async function syncWithChildrens() {
        await Promise.all([syncItems(), contactStore.syncItems(), jobApplicationStore.syncItems(), eventStore.syncItems()])
    }
    async function syncWithParents() {
        await Promise.all([companyTypeStore.syncItems(), syncItems()])
    }

    function add(): void {
        items.value = [Object.assign({}, {
            id: 0,
            name: '',
            notes: '',
            company_types: [],
            application_cnt: 0,
            last_event: null,
            next_event: null,
            convertValues: ItemClass.prototype.convertValues
        }), ...items.value]
    }
    function eq(i1: Item, i2: Item): boolean {
        return i1.name === i2.name

    }
    async function insert(item: Item): Promise<Item> {
        return await back.InsertCompany(item)
    }
    async function delete_(item: Item): Promise<void> {
        return await back.DeleteCompany(item)
    }
    async function update(item: Item): Promise<void> {
        return await back.UpdateCompany(item)
    }

    function findNamesFromContact(contact: main.Contact): string[] {
        return items.value
            .map((e) => { return e.name })
    }

    function findNamesFromApplication(application: main.JobApplication): string[] {
        return items.value
            .map((e) => { return e.name })
    }

    function findNamesFromEvent(event: main.Event): string[] {
        return items.value
            .map((e) => { return e.name })
    }

    return {
        items, columns, syncItems, syncWithChildrens, syncWithParents, add, eq, select, insert, delete_, update, findNamesFromContact, findNamesFromApplication, findNamesFromEvent
    }
})
