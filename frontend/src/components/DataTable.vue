<script lang="ts" setup generic="T extends { id: number }">
import { computed, ref } from 'vue'
import {
    VAlert,
    VCard,
    VCardItem,
    VCardActions,
    VDataTable,
    VChip,
    VTextField,
    VTextarea,
    VAutocomplete,
    VTooltip,
    VBtn,
    VListItem,
    VIcon,
    VDialog,
    VContainer,
    VCol,
    VRow,
    VCheckbox

} from 'vuetify/components'
import { VNumberInput } from 'vuetify/labs/VNumberInput'
import DateInput from './DateInput.vue'
import * as types from './types'
import * as utils from './utils'
import { Store } from '../stores/interface'



const props = defineProps<{ store: Store<T> }>()

await props.store.syncWithParents()

const columns = props.store.getColumns()
const columnsWithDelete: types.VuetifyHeaders = (
    [{ title: 'Delete', key: 'del', sortable: false }] as types.VuetifyHeaders
).concat(columns.filter((e) => { return !e.readOnly }) as types.VuetifyHeaders)
const inputs = ref<types.Inputs[]>([])
const alertMessage = ref<types.AlertMessage>({
    text: "",
    title: "",
    type: "info",
    displayed: false,
})
const edit = ref(false)

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


async function toogleEdit() {
    const isEdit = edit.value
    if (isEdit) {
        await props.store.syncItems()
        edit.value = false
    } else {
        edit.value = true
    }
}

function formatDateBeforeSend(item: T): T {
    let cp: T = Object.assign({}, item)
    for (const col of columns) {
        if (col.type === "date") {
            if (item[col.key]) {
                cp[col.key] = utils.formatDateForBackend(item[col.key] as Date) as T[keyof T]
            }
        }
    }
    return cp
}

async function save() {
    for (const inp of inputs.value) {
        if (Array.isArray(inp.errorMessages) && inp.errorMessages.length > 0) {
            setAlert("error", "Error", "Some part of the forms are invalid, please fix the error before saving (focus especially each inputs on new lines)")
            return
        }
    }
    let insertCnt = 0
    let updateCnt = 0
    let deleteCnt = 0
    let modifiedMap = utils.itemsMap(props.store.items)
    const itemsBackup = await props.store.select()
    let backupMap = utils.itemsMap(await props.store.select())
    for (const oldItem of itemsBackup) {
        let delItem = modifiedMap.get(oldItem.id)
        if (!delItem) {
            // delete
            try {
                await props.store.delete_(formatDateBeforeSend(oldItem))
                deleteCnt += 1
            } catch (e) {
                setError(e)
            }
        }
    }
    for (const modifiedItem of props.store.items) {
        let oldItem = backupMap.get(modifiedItem.id)
        let rawModifiedItem = utils.deepToRaw(modifiedItem)
        if (!oldItem) {
            // insert
            try {
                let itemWithId = await props.store.insert(formatDateBeforeSend(modifiedItem))
                modifiedItem.id = itemWithId.id
                itemsBackup.push(rawModifiedItem)
                insertCnt += 1
            } catch (e) {
                setError(e)
            }
        } else if (!utils.areObjsEq(oldItem, rawModifiedItem)) {
            // update
            try {
                await props.store.update(formatDateBeforeSend(modifiedItem))
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
    await props.store.syncWithChildrens()
}

function deleteItem(index: number) {
    props.store.items.splice(index, 1)
}

function newItem() {
    edit.value = true
    props.store.add()
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

const DATE_WIDTH_NUMBER = 12;
const TIME_WIDTH_NUMBER = 8;

const ACTION_WIDTH = "4em";
const DATE_WIDTH = `${DATE_WIDTH_NUMBER}em`;
const TIME_WIDTH = `${TIME_WIDTH_NUMBER}em`;
const DATETIME_WIDTH = "12em";
const STRING_WIDTH = "12em";
const MULTILINE_WIDTH = "30em";
const INT_WIDTH = "9em";
const REL_WIDTH = "10em";
const LISTREL_WIDTH = "15em";

function getCellStyle(column: types.Column<T>) {
    let width = ""
    switch (column.type) {
        case 'date':
            width = DATETIME_WIDTH
            break
        case 'string':
            width = STRING_WIDTH
            break
        case 'multiline':
            width = MULTILINE_WIDTH
            break
        case 'int':
            width = INT_WIDTH
            break
        case 'listrel':
            width = LISTREL_WIDTH
            break
        case 'rel':
            width = REL_WIDTH
            break
        case 'action':
            width = ACTION_WIDTH
            break
    }
    return {
        width: width,
        maxWidth: width,
    }
}

const ACTION_CELL_STYLE = {
    width: ACTION_WIDTH,
    maxWidth: ACTION_WIDTH,
}
const CELL_CLASSES = [
    "cell",
    "v-data-table__td",
    "v-data-table-column--align-start"
]

const search = ref("")
const searchReg = ref(false)
const invalidRegErrs = ref<string[]>([])

const searchRegIsInvalid = computed(() => {
    return invalidRegErrs.value.length > 0
})

function compareString(value: string): boolean {
    const lowerSearch = search.value.toLowerCase()
    return value.toLowerCase().includes(lowerSearch)
}

function compareStringReg(value: string): boolean {
    try {
        const regSearch = new RegExp(search.value)
        invalidRegErrs.value = []
        return regSearch.test(value)
    } catch (e) {
        invalidRegErrs.value = [`invalid regex: ${e}`]
        return true
    }
}

function compareValueToSearch(value: any): boolean {
    let compMethod = compareString
    if (searchReg.value) {
        compMethod = compareStringReg
    }
    if (Array.isArray(value)) {
        for (const e of value) {
            if (compareValueToSearch(e)) {
                return true
            }
        }
    } else if (value instanceof Date && compMethod(utils.formatDateToLocale(value))) {
        return true
    } else if (value && compMethod(String(value))) {
        return true
    }
    return false
}


const searchFilter = computed(() => {
    if (search.value) {
        return props.store.items.filter(
            function (e) {
                for (const key in e) {
                    let prop = e[key]
                    let eq = compareValueToSearch(prop)
                    if (eq) {
                        return true
                    }
                }
                return false
            }
        )
    }
    invalidRegErrs.value = []
    return props.store.items
})

</script>

<template>
    <v-alert :text="alertMessage.text" :title="alertMessage.title" :type="alertMessage.type" closable
        v-model="alertMessage.displayed" elevation="24" class="elevated"></v-alert>

    <v-card variant="tonal" color="secondary">
        <v-card-actions>
            <v-container>
                <v-row>
                    <v-col>
                        <v-text-field v-model="search" label="Search" density="compact" prepend-icon="mdi-file-search"
                            :error-messages="invalidRegErrs" :error="searchRegIsInvalid" width="30em"
                            clearable></v-text-field>
                    </v-col>
                    <v-col>
                        <v-checkbox v-model="searchReg" label="Regex"></v-checkbox>
                    </v-col>
                    <v-col>
                        <v-tooltip :text="edit ? 'Cancel modifications' : 'Edit'" location="top">
                            <template v-slot:activator="{ props }">
                                <v-btn v-bind="props" color="secondary" variant="flat"
                                    :icon="edit ? 'mdi-cancel' : 'mdi-pencil'" density="comfortable"
                                    @click="toogleEdit"></v-btn>
                            </template>
                        </v-tooltip>
                    </v-col>
                    <v-col>
                        <v-tooltip text="Add" location="top">
                            <template v-slot:activator="{ props }">
                                <v-btn v-bind="props" color="secondary" variant="flat" icon="mdi-plus"
                                    density="comfortable" @click="newItem" v-show="edit"></v-btn>
                            </template>
                        </v-tooltip>
                    </v-col>
                    <v-col>
                        <v-tooltip text="Save" location="top">
                            <template v-slot:activator="{ props }">
                                <v-btn v-bind="props" color="secondary" variant="flat" icon="mdi-content-save-all"
                                    density="comfortable" @click="save" v-show="edit"></v-btn>
                            </template>
                        </v-tooltip>
                    </v-col>
                </v-row>
                <v-row no-gutters>
                </v-row>
            </v-container>
        </v-card-actions>
    </v-card>

    <v-data-table :headers="edit ? columnsWithDelete : columns as types.VuetifyHeaders" :items="searchFilter"
        density="compact" height="65vh" items-per-page="10">

        <template v-slot:headers="{ columns, isSorted, getSortIcon, toggleSort }">
            <tr>
                <template v-for="c in columns" :key="c.key">
                    <td class="cell" :style="getCellStyle(c as unknown as types.Column<T>)">
                        <span class="mr-2 cursor-pointer header" @click="() => toggleSort(c)">{{ c.title }}</span>
                        <template v-if="isSorted(c)">
                            <v-icon :icon="getSortIcon(c)"></v-icon>
                        </template>
                    </td>
                </template>
            </tr>
        </template>

        <template v-slot:item="{ item, index }">
            <tr v-show="!edit" class="v-data-table__tr">
                <td :class="CELL_CLASSES" :style="getCellStyle(c)" v-for="c in columns">
                    <template v-if="Array.isArray(item[c.key])">
                        <v-chip v-for="subitem in item[c.key]">
                            {{ subitem }}
                        </v-chip>
                    </template>
                    <template v-else-if="c.type === 'multiline'">
                        <v-tooltip location="start" :disabled="formatCol(c, item[c.key]).length === 0">
                            <template v-slot:activator="{ props }">
                                <div v-bind="props" class="multiline multilinecontent">{{ formatCol(c, item[c.key]) }}
                                </div>
                            </template>
                            <div class="multiline">{{ formatCol(c, item[c.key]) }}</div>
                        </v-tooltip>
                    </template>
                    <template v-else>
                        {{ formatCol(c, item[c.key]) }}
                    </template>
                </td>
            </tr>
            <tr v-show="edit" class="v-data-table__tr">
                <td :class="CELL_CLASSES" :style="ACTION_CELL_STYLE">
                    <v-tooltip text="Delete line" location="top">
                        <template v-slot:activator="{ props }">
                            <div class="buttonActionContainer">
                                <v-btn v-bind="props" color="secondary" variant="plain" icon="mdi-delete"
                                    density="comfortable" @click="deleteItem(index)"></v-btn>
                            </div>
                        </template>
                    </v-tooltip>
                </td>
                <td :class="CELL_CLASSES" v-for="c in columns.filter((e) => !e.readOnly)" :style="getCellStyle(c)">
                    <v-text-field v-if="c.type === 'string'" v-model="item[c.key]" :rules="buildRules(c)" ref="inputs"
                        density="compact" :width="STRING_WIDTH"></v-text-field>
                    <v-dialog v-else-if="c.type === 'multiline'" scrollable>
                        <template v-slot:activator="{ props: activatorProps }">
                            <v-tooltip location="start" :disabled="formatCol(c, item[c.key]).length === 0">
                                <template v-slot:activator="{ props: multiLineActivator }">
                                    <div v-bind="multiLineActivator" class="buttonActionContainer">
                                        <v-btn v-bind="activatorProps" color="secondary" variant="plain"
                                            icon="mdi-arrow-expand" density="comfortable"></v-btn>
                                    </div>
                                </template>
                                <div class="multiline">{{ formatCol(c, item[c.key]) }}</div>
                            </v-tooltip>
                        </template>
                        <template v-slot:default="{ isActive }">
                            <v-card>
                                <v-card-item>
                                    <v-textarea v-model="item[c.key]" :rules="buildRules(c)" ref="inputs"
                                        :label="c.title" rows="40"></v-textarea>
                                </v-card-item>
                                <v-card-actions style="justify-content: flex-start;">
                                    <v-btn text="Close" color="primary" variant="text" density="comfortable"
                                        @click="isActive.value = false"></v-btn>
                                </v-card-actions>
                            </v-card>
                        </template>
                    </v-dialog>
                    <v-number-input v-else-if="c.type === 'int'" v-model="item[c.key]" :rules="buildRules(c)"
                        ref="inputs" density="compact" :width="INT_WIDTH"></v-number-input>
                    <date-input v-else-if="c.type === 'date'" v-model="item[c.key] as Date" :rules="buildRules(c)"
                        :date-width="DATE_WIDTH" :time-width="TIME_WIDTH" ref="inputs"></date-input>
                    <v-autocomplete v-else-if="c.type === 'listrel' && c.relations" v-model="item[c.key] as string[]"
                        :items="c.relations(item)" :rules="buildRules(c)" ref="inputs" chips multiple clearable
                        density="compact" :width="LISTREL_WIDTH">
                        <template v-slot:item="slotItem">
                            <v-list-item v-bind="slotItem.props" :active="false" :title="slotItem.item.props.value"
                                :append-icon="(item[c.key] as string[]).indexOf(slotItem.item.props.value) != -1 ? 'mdi-checkbox-multiple-marked-circle' : 'mdi-checkbox-multiple-blank-circle-outline'"></v-list-item>
                        </template>
                    </v-autocomplete>
                    <v-autocomplete v-else-if="c.type === 'rel' && c.relations" v-model="item[c.key] as string"
                        :items="c.relations(item)" :rules="buildRules(c)" ref="inputs" density="compact"
                        :width="REL_WIDTH">
                        <template v-slot:item="slotItem">
                            <v-list-item v-bind="slotItem.props" :active="false" :title="slotItem.item.props.value"
                                :append-icon="(item[c.key] as string) == slotItem.item.props.value ? 'mdi-checkbox-marked-circle' : 'mdi-checkbox-blank-circle-outline'"></v-list-item>
                        </template>
                    </v-autocomplete>
                </td>
            </tr>
        </template>
    </v-data-table>
</template>

<style lang="css" scoped>
.header {
    cursor: pointer;
    user-select: none;
}

.elevated {
    z-index: 100;
    position: absolute;
}

.multiline {
    white-space: pre-wrap;
}

.multilinecontent {
    height: 3em;
}

.cell {
    white-space: pre-wrap;
    overflow: hidden;
    padding: 0 !important;
}

.buttonActionContainer {
    display: flex;
    justify-content: flex-start;
    align-items: center;
}
</style>