<script lang="ts" setup>
import { VTextField } from 'vuetify/components'
import { computed, ref, watch } from 'vue'
import * as utils from './utils'


const props = defineProps<{
    rules: ((val: any) => string | boolean)[],
    dateWidth: string,
    timeWidth: string
}>()

const dateModel = defineModel<Date | null, string>({ required: true })
const timeStr = ref(utils.formatDateToInput(dateModel.value).time)
// dbl dbl binding :O
const dateStr = computed<string>({
    get() {
        return utils.formatDateToInput(dateModel.value).date
    },
    set(value) {
        dateModel.value = utils.parseInputDate(value, timeStr.value)
    }

})

watch(timeStr, function (value: string) {
    dateModel.value = utils.parseInputDate(dateStr.value, value)
})

</script>

<template>
    <div style="display: flex; flex-direction: row;">
        <v-text-field type="date" v-model="dateStr" :rules="props.rules" density="compact" :width="dateWidth">
        </v-text-field>
        <v-text-field type="time" v-model="timeStr" :rules="props.rules" density="compact"
            :width="timeWidth"></v-text-field>
    </div>
</template>