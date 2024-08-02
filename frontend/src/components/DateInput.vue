<script lang="ts" setup>
import { VTextField } from 'vuetify/components'
import { computed } from 'vue'
import * as utils from './utils'


const props = defineProps<{
    rules: ((val: any) => string | boolean)[]
}>()

const model = defineModel<Date, string>({ required: true })
// dbl dbl binding :O
const dateStr = computed<string>({
    get() {
        return utils.formatDateToInput(model.value)
    },
    set(value) {
        model.value = utils.parseInputDate(value)
    }

})

</script>

<template>
    <v-text-field type="date" v-model="dateStr" :rules="props.rules" density="compact"></v-text-field>
</template>