<script lang="ts" setup generic="T extends { id: number }">
import { ref } from 'vue'
import {
    VAlert,
    VCard,
    VCardActions,
    VSpacer,
    VDataTable,
    VChip,
    VTextField,
    VTextarea,
    VSelect,
    VTooltip,
    VBtn,

} from 'vuetify/components'
import { VNumberInput } from 'vuetify/labs/VNumberInput'
import DateInput from './DateInput.vue'
import * as types from './types'
import * as utils from './utils'

type KeyofItem = keyof T

const props = defineProps<{
    items: T[]
    delete: types.Delete<T>
    update: types.Update<T>
    insert: types.Insert<T>
    emptyItem: T
    columns: types.Columns<T>
    relations?: types.Relations<T>
    listRelations?: types.ListRelations<T>
    title?: string
}>()


const columnsWithDelete: types.VuetifyHeaders = (props.columns as types.VuetifyHeaders).concat([{ title: 'Delete', key: 'del', sortable: false }])
const inputs = ref<types.Inputs[]>([])
const alertMessage = ref<types.AlertMessage>({
    text: "",
    title: "",
    type: "info",
    displayed: false,
})
const edit = ref(false)
let itemsBackup = utils.deepToRaw(props.items)

function setAlert(type: types.AlertType, title: string, message: string) {
    alertMessage.value.displayed = true
    alertMessage.value.type = type
    alertMessage.value.title = title
    alertMessage.value.text = message
    if (type !== 'error') {
        setTimeout(function () {
            alertMessage.value.displayed = false
        }, 10000)
    }
}

function setError(exception: any) {
    console.trace("Error when saving:", exception)
    let errStr = ""
    if (exception instanceof Error) {
        errStr = exception.toString()
    } else {
        errStr = JSON.stringify(exception)
    }
    setAlert('error', "Error", errStr)
    throw exception
}

function replaceArrayContent(current: T[], new_: T[]) {
    current.length = 0
    for (const item of new_) {
        current.push(item)
    }
}

function toogleEdit() {
    const isEdit = edit.value
    if (isEdit) {
        replaceArrayContent(props.items, itemsBackup)
        edit.value = false
    } else {
        itemsBackup = utils.deepToRaw(props.items)
        edit.value = true
    }
}

function formatDateBeforeSend(item: T): T {
    let cp: T = Object.assign({}, item)
    for (const col of props.columns) {
        if (col.type === "date") {
            cp[col.key] = utils.formatDateForBackend(item[col.key] as Date) as T[keyof T]
        }
    }
    return cp
}

async function save() {
    for (const inp of inputs.value) {
        if (Array.isArray(inp.errorMessages) && inp.errorMessages.length > 0) {
            setAlert("error", "Error", "Some part of the forms are invalid, please fix the error to validate")
            return
        }
    }
    let insertCnt = 0
    let updateCnt = 0
    let deleteCnt = 0
    let modifiedMap = utils.itemsMap(props.items)
    let backupMap = utils.itemsMap(itemsBackup)
    for (const oldItem of itemsBackup) {
        let delItem = modifiedMap.get(oldItem.id)
        if (!delItem) {
            // delete
            try {
                await props.delete(formatDateBeforeSend(oldItem))
                deleteCnt += 1
            } catch (e) {
                setError(e)
            }
        }
    }
    for (const modifiedItem of props.items) {
        let oldItem = backupMap.get(modifiedItem.id)
        let rawModifiedItem = utils.deepToRaw(modifiedItem)
        if (!oldItem) {
            // insert
            try {
                let itemWithId = await props.insert(formatDateBeforeSend(modifiedItem))
                modifiedItem.id = itemWithId.id
                itemsBackup.push(rawModifiedItem)
                insertCnt += 1
            } catch (e) {
                setError(e)
            }
        } else if (!utils.areObjsEq(oldItem, rawModifiedItem)) {
            // update
            try {
                await props.update(formatDateBeforeSend(modifiedItem))
                updateCnt += 1
            } catch (e) {
                setError(e)
            }
        }
    }
    if (insertCnt +
        updateCnt +
        deleteCnt === 0) {
        setAlert('warning', "No change", "They were no changes to save")
    } else {
        setAlert('success', "Saved", `Changes successfully saved (${insertCnt} inserts, ${updateCnt}, updates, ${deleteCnt} deletes)`)
    }
    edit.value = false
}

function deleteItem(item: T) {
    replaceArrayContent(
        props.items,
        props.items.filter(function (e) { return e.id != item.id })
    )
}

function newItem() {
    edit.value = true
    props.items.push(Object.assign({}, props.emptyItem))
}

function requiredCheck(value: any): string | boolean {
    if (typeof value === "string" && value === "") {
        return "please provide at a value"
    }
    if (Array.isArray(value) && value.length === 0) {
        return "please provide at least one value"
    }
    if (value === undefined || value === null) {
        return "please provide at least one value"
    }
    return true
}

function buildRules(col: types.Column<T>): ((val: any) => string | boolean)[] {
    if (col.requiered) {
        return [requiredCheck]
    }
    return []
}

function formatCol(col: types.Column<T>, value: any): string {
    if (col.type === 'date') {
        return utils.formatDateToLocale(value)
    }
    if (value === undefined || value === null) {
        return ""
    }
    if (typeof value === "string") {
        return value
    }
    if (typeof value === "number") {
        return value.toString()
    }
    throw new Error(`unknown type ${typeof value} for column ${String(col.key)}, value is ${value}`)
}

</script>

<template>
    <v-alert :text="alertMessage.text" :title="alertMessage.title" :type="alertMessage.type" closable
        v-model="alertMessage.displayed" elevation="24"></v-alert>

    <v-card variant="tonal" color="secondary" :title="props.title" style="text-align: center;">
        <v-card-actions>
            <v-spacer />
            <v-tooltip :text="edit ? 'Cancel modifications' : 'Edit'" location="top">
                <template v-slot:activator="{ props }">
                    <v-btn v-bind="props" color="secondary" variant="flat" :icon="edit ? 'mdi-cancel' : 'mdi-pencil'"
                        density="comfortable" @click="toogleEdit()"></v-btn>
                </template>
            </v-tooltip>
            <v-tooltip text="Add" location="top">
                <template v-slot:activator="{ props }">
                    <v-btn v-bind="props" color="secondary" variant="flat" icon="mdi-plus" density="comfortable"
                        @click="newItem()" v-show="edit"></v-btn>
                </template>
            </v-tooltip>
            <v-tooltip text="Save" location="top">
                <template v-slot:activator="{ props }">
                    <v-btn v-bind="props" color="secondary" variant="flat" icon="mdi-content-save-all"
                        density="comfortable" @click="save()" v-show="edit"></v-btn>
                </template>
            </v-tooltip>
            <v-spacer />
        </v-card-actions>
    </v-card>

    <v-data-table :headers="edit ? columnsWithDelete : props.columns as types.VuetifyHeaders" :items="props.items"
        density="compact" hide-default-footer>
        <template v-slot:item="{ item }">
            <tr v-show="!edit" class="v-data-table__tr">
                <td class="v-data-table__td v-data-table-column--align-start" v-for="c in props.columns">
                    <template v-if="Array.isArray(item[c.key])">
                        <v-chip v-for="subitem in item[c.key]">
                            {{ subitem }}
                        </v-chip>
                    </template>
                    <template v-else>
                        {{ formatCol(c, item[c.key]) }}
                    </template>
                </td>
            </tr>
            <tr v-show="edit" class="v-data-table__tr">
                <td class="v-data-table__td v-data-table-column--align-start" v-for="c in props.columns">
                    <v-text-field v-if="c.type === 'string'" v-model="item[c.key]" :rules="buildRules(c)" ref="inputs"
                        density="compact"></v-text-field>
                    <v-textarea v-else-if="c.type === 'multiline'" v-model="item[c.key]" :rules="buildRules(c)"
                        ref="inputs" density="compact" no-resize></v-textarea>
                    <v-number-input v-else-if="c.type === 'int'" v-model="item[c.key]" :rules="buildRules(c)"
                        ref="inputs" density="compact"></v-number-input>
                    <date-input v-else-if="c.type === 'date'" v-model="item[c.key] as Date" :rules="buildRules(c)"
                        ref="inputs"></date-input>
                    <v-select v-else-if="c.type === 'listrel' && props.listRelations" v-model="item[c.key] as string[]"
                        :items="(props.listRelations[c.key] as types.RelationFinder<T>)(item)" :rules="buildRules(c)"
                        ref="inputs" chips multiple clearable density="compact">
                    </v-select>
                    <v-select v-else-if="c.type === 'rel' && props.relations" v-model="item[c.key] as string"
                        :items="(props.relations[c.key] as types.RelationFinder<T>)(item)" :rules="buildRules(c)"
                        ref="inputs" density="compact"></v-select>
                </td>
                <td class="v-data-table__td v-data-table-column--align-start">
                    <v-tooltip text="Delete line" location="top">
                        <template v-slot:activator="{ props }">
                            <v-btn v-bind="props" color="secondary" variant="plain" icon="mdi-delete"
                                density="comfortable" @click="deleteItem(item)"></v-btn>
                        </template>
                    </v-tooltip>
                </td>
            </tr>
        </template>
    </v-data-table>
</template>
