import { loadData } from './utils'
import * as utils from '../components/utils'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import * as back from '../../wailsjs/go/main/App'
import { main } from '../../wailsjs/go/models'
import * as types from '../components/types'
import { useCompanyStore } from './company'
import { useContactStore } from './contact'
import { useEventSourceStore } from './event_source'
import { useJobApplicationStore } from './job_application'

type Item = main.Event
const ItemClass = main.Event

export const useEventStore = defineStore("Event", () => {
    const companyStore = useCompanyStore()
    const contactStore = useContactStore()
    const eventSourceStore = useEventSourceStore()
    const jobApplicationStore = useJobApplicationStore()

    const items = ref<Item[]>([])

    const columns: types.Columns<Item> = [
        { key: "company_name", title: "Company", type: "rel", requiered: true, relations: companyStore.findNamesFromEvent },
        { key: "job_title", title: "Job", type: "rel", requiered: true, relations: jobApplicationStore.findNamesFromEvent },
        { key: "title", title: "Title", type: "string", requiered: true },
        { key: "date", title: "Date", type: "date", requiered: true },
        { key: "source", title: "Source", type: "rel", requiered: true, relations: eventSourceStore.findNamesFromEvent },
        { key: "contacts", title: "Contacts", type: "listrel", relations: contactStore.findNamesFromEvent },
        { key: "notes", title: "Notes", type: "multiline" },

    ]

    async function select() {
        let eventsConv = await back.ListEvents()
        for (const e of eventsConv) {
            e.date = utils.parseBackendDate(e.date)
        }
        return eventsConv
    }

    async function syncItems() {
        await loadData(async () => {
            items.value = await select()
        }, ItemClass)
    }
    async function syncWithChildrens() {
        await syncItems()
    }
    async function syncWithParents() {
        await Promise.all([companyStore.syncItems(), contactStore.syncItems(), eventSourceStore.syncItems(), jobApplicationStore.syncItems(), syncItems()])
    }

    function add(): void {
        items.value = [Object.assign({}, {
            company_name: "",
            contacts: [],
            date: new Date(),
            id: 0,
            job_title: "",
            notes: "",
            source: "",
            title: "",
            convertValues: ItemClass.prototype.convertValues
        }), ...items.value]
    }
    function eq(i1: Item, i2: Item): boolean {
        return i1.job_title === i2.job_title && i1.title === i2.title
    }
    async function insert(item: Item): Promise<Item> {
        let newitem = await back.InsertEvent(item)
        newitem.date = utils.parseBackendDate(newitem.date)
        return newitem
    }
    async function delete_(item: Item): Promise<void> {
        return await back.DeleteEvent(item)
    }
    async function update(item: Item): Promise<void> {
        return await back.UpdateEvent(item)
    }

    const calendarEvents = computed(function () {
        return items.value.map(function (e) {
            const eventDate = utils.parseBackendDate(e.date)
            const calendarEvent: types.CalendarEvent = {
                title: e.title,
                start: eventDate,
                end: eventDate,
                allDay: true,
                eventId: e.id,
                company: e.company_name,
                job: e.job_title,
                contacts: e.contacts,
            }
            return calendarEvent
        })
    })

    return {
        items, columns, syncItems, syncWithChildrens, syncWithParents, add, eq, select, insert, delete_, update, calendarEvents
    }
})
