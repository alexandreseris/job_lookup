import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as back from '../wailsjs/go/main/App'
import { main, db } from '../wailsjs/go/models'
import * as utils from './components/utils'
import * as types from './components/types'

export const useStore = defineStore('store', () => {
    const companies = ref<main.Company[]>([])
    const companyTypes = ref<db.CompanyType[]>([])
    const events = ref<main.Event[]>([])
    const contacts = ref<db.ListContactRow[]>([])
    const eventSource = ref<db.EventSource[]>([])
    const applications = ref<db.ListJobApplicationRow[]>([])
    const applicationStatus = ref<db.JobApplicationStatus[]>([])

    async function init() {
        await Promise.all([
            async function () { companies.value = await back.ListCompanies() }(),
            async function () { companyTypes.value = await back.ListCompanyTypes() }(),
            async function () {
                let eventsConv = await back.ListEvents()
                for (const e of eventsConv) {
                    e.date = utils.parseBackendDate(e.date)
                }
                events.value = eventsConv
            }(),
            async function () { contacts.value = await back.ListContact() }(),
            async function () { eventSource.value = await back.ListEventSource() }(),
            async function () { applications.value = await back.ListJobApplication() }(),
            async function () { applicationStatus.value = await back.ListJobApplicationStatus() }(),
        ])
    }

    const companyTypeNames = computed(function () {
        return companyTypes.value.map(function (e) { return e.name })
    })

    const companyNames = computed(function () {
        return companies.value.map(function (e) { return e.name })
    })

    const eventSourceNames = computed(function () {
        return eventSource.value.map(function (e) { return e.name })
    })
    const contactNames = computed(function () {
        return contacts.value.map(function (e) { return `${e.fist_name}, ${e.last_name}` })
    })

    const applicationStatusNames = computed(function () {
        return applicationStatus.value.map(function (e) { return e.name })
    })

    const calendarEvents = computed(function () {
        return events.value.map(function (e) {
            const eventDate = utils.parseBackendDate(e.date)
            return {
                title: e.title,
                start: eventDate,
                end: eventDate,
                allDay: true,
                eventId: e.id,
                company: e.company_name,
                job: e.job_title,
                contacts: e.contacts,
            } as types.CalendarEvent
        })
    })

    return {
        init,
        companies,
        companyTypes,
        events,
        contacts,
        eventSource,
        applications,
        applicationStatus,
        companyTypeNames,
        companyNames,
        eventSourceNames,
        contactNames,
        applicationStatusNames,
        calendarEvents,
    }
})