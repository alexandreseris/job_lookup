import { loadData } from './utils'
import * as utils from '../components/utils'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as back from '../../wailsjs/go/main/App'
import { main } from '../../wailsjs/go/models'
import * as types from '../components/types'
import { useCompanyStore } from './company'
import { useEventStore } from './event'

type Item = main.Contact
const ItemClass = main.Company

export const useContactStore = defineStore("Contact", () => {
    const companyStore = useCompanyStore()

    const items = ref<Item[]>([])


    function getColumns() {
        const columns: types.Columns<Item> = [
            { key: 'company_name', title: "Company", type: "rel", requiered: true, relations: companyStore.findNamesFromContact },
            { key: 'fist_name', title: "Fist name", type: "string", requiered: true },
            { key: 'last_name', title: "Last name", type: "string", requiered: true },
            { key: 'job_position', title: "Position", type: "string", requiered: true },
            { key: 'email', title: "Email", type: "string" },
            { key: 'phone_number', title: "Phone", type: "string" },
            { key: "last_event", title: "Last event", type: "date", readOnly: true },
            { key: "next_event", title: "Next event", type: "date", readOnly: true },
            { key: 'notes', title: "Notes", type: "multiline" },
        ]
        return columns
    }

    async function select() {
        let contactConv = await back.ListContact()
        for (const e of contactConv) {
            e.last_event = utils.parseBackendDateOpt(e.last_event)
            e.next_event = utils.parseBackendDateOpt(e.next_event)
        }
        return contactConv
    }

    async function syncItems() {
        await loadData(async () => {
            items.value = await select()
        }, ItemClass)
    }
    async function syncWithChildrens() {
        const eventStore = useEventStore()
        await Promise.all([syncItems(), eventStore.syncItems()])
    }
    async function syncWithParents() {
        await Promise.all([companyStore.syncItems(), syncItems()])
    }

    function add(): void {
        items.value = [Object.assign({}, {
            company_id: 0,
            company_name: "",
            email: "",
            fist_name: "",
            id: 0,
            job_position: "",
            last_name: "",
            notes: "",
            phone_number: "",
            last_event: null,
            next_event: null,
            convertValues: ItemClass.prototype.convertValues
        }), ...items.value]
    }
    function eq(i1: Item, i2: Item): boolean {
        return i1.company_name === i2.company_name && i1.fist_name === i2.fist_name && i1.last_name === i2.last_name

    }
    async function insert(item: Item): Promise<Item> {
        return await back.InsertContact(item)
    }
    async function delete_(item: Item): Promise<void> {
        return await back.DeleteContact(item)
    }
    async function update(item: Item): Promise<void> {
        return await back.UpdateContact(item)
    }

    function findNamesFromEvent(event: main.Event): string[] {
        return items.value
            .filter((e) => { return e.company_name === event.company_name })
            .map((e) => { return `${e.fist_name}, ${e.last_name}` })
    }

    return {
        items, getColumns, syncItems, syncWithChildrens, syncWithParents, add, eq, select, insert, delete_, update, findNamesFromEvent
    }
})
