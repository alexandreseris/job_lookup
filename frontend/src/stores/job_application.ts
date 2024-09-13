import { loadData } from './utils'
import * as utils from '../components/utils'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as back from '../../wailsjs/go/main/App'
import { main } from '../../wailsjs/go/models'
import * as types from '../components/types'
import { useCompanyStore } from './company'
import { useJobApplicationStatusStore } from './job_application_status'
import { useEventStore } from './event'

type Item = main.JobApplication
const ItemClass = main.JobApplication

export const useJobApplicationStore = defineStore("JobApplication", () => {
    const companyStore = useCompanyStore()
    const jobApplicationStatusStore = useJobApplicationStatusStore()

    const items = ref<Item[]>([])

    function getColumns() {
        const columns: types.Columns<Item> = [
            { key: "company_name", title: "Company", type: "rel", requiered: true, relations: companyStore.findNamesFromApplication },
            { key: "job_title", title: "Title", type: "string", requiered: true },
            { key: "status_name", title: "Status", type: "rel", requiered: true, relations: jobApplicationStatusStore.findNamesFromApplication },
            { key: "event_cnt", title: "Number of events", type: "int", readOnly: true },
            { key: "last_event", title: "Last event", type: "date", readOnly: true },
            { key: "next_event", title: "Next event", type: "date", readOnly: true },
            { key: "notes", title: "Notes", type: "multiline" },
        ]
        return columns
    }

    async function select() {
        let applicationConv = await back.ListJobApplication()
        for (const e of applicationConv) {
            e.last_event = utils.parseBackendDateOpt(e.last_event)
            e.next_event = utils.parseBackendDateOpt(e.next_event)
        }
        return applicationConv
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
        await Promise.all([companyStore.syncItems(), jobApplicationStatusStore.syncItems(), syncItems()])
    }

    function add(): void {
        items.value = [Object.assign({}, {
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
            convertValues: ItemClass.prototype.convertValues
        }), ...items.value]
    }
    function eq(i1: Item, i2: Item): boolean {

        return i1.company_name === i2.company_name && i1.job_title === i2.job_title
    }
    async function insert(item: Item): Promise<Item> {
        return await back.InsertJobApplication(item)
    }
    async function delete_(item: Item): Promise<void> {
        return await back.DeleteJobApplication(item)
    }
    async function update(item: Item): Promise<void> {
        return await back.UpdateJobApplication(item)
    }

    function findNamesFromEvent(event: main.Event): string[] {
        return items.value
            .filter((e) => { return e.company_name === event.company_name })
            .map((e) => { return e.job_title })
    }

    return {
        items, getColumns, syncItems, syncWithChildrens, syncWithParents, add, eq, select, insert, delete_, update, findNamesFromEvent
    }
})
