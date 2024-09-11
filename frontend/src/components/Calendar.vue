<script lang="ts" setup>
import { VCard, VCardText, VChip } from 'vuetify/components'
import { VCalendar } from 'vuetify/labs/VCalendar'
import { useStore } from '../store';
import * as types from './types'

const store = useStore()
await store.init()
function renderEvent(event: types.CalendarEvent) {
    return {
        title: `${event.company} - ${event.job}`,
        subtitle: event.title,
    }
}
</script>

<template>
    <v-calendar :events="store.calendarEvents" :first-day-of-week="1" show-adjacent-months>
        <template v-slot:event="{ event }">
            <v-card v-bind="renderEvent(event as types.CalendarEvent)" variant="outlined" color="secondary"
                density="compact">
                <v-card-text>
                    <template v-for="c in event.contacts">
                        <v-chip density="comfortable">
                            {{ c }}
                        </v-chip>
                    </template>
                </v-card-text>
            </v-card>
        </template>
    </v-calendar>
</template>