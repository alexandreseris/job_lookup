import { loadData } from './utils'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as back from '../../wailsjs/go/main/App'
import { db, main } from '../../wailsjs/go/models'
import * as types from '../components/types'
import { useEventStore } from './event'

type Item = db.ListEventSourceRow
const ItemClass = db.ListEventSourceRow

export const useEventSourceStore = defineStore("EventSource", () => {
    const items = ref<Item[]>([])

    function getColumns() {
        const columns: types.Columns<Item> = [
            { title: 'Name', key: 'name', type: "string", requiered: true },
            { title: 'Events', key: 'events', type: "int", readOnly: true },
        ]
        return columns
    }

    async function select() {
        return await back.ListEventSource()
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
        await syncItems()
    }

    function add(): void {
        items.value = [Object.assign({}, {
            id: 0,
            name: '',
            events: 0
        }), ...items.value]
    }
    function eq(i1: Item, i2: Item): boolean {
        return i1.name === i2.name
    }
    async function insert(item: Item): Promise<Item> {
        return await back.InsertEventSource(item)
    }
    async function delete_(item: Item): Promise<void> {
        return await back.DeleteEventSource(item)
    }
    async function update(item: Item): Promise<void> {
        return await back.UpdateEventSource(item)
    }

    function findNamesFromEvent(event: main.Event): string[] {
        return items.value
            .map((e) => { return e.name })
    }

    return {
        items, getColumns, syncItems, syncWithChildrens, syncWithParents, add, eq, select, insert, delete_, update, findNamesFromEvent
    }
})
