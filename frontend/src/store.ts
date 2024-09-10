import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as back from '../wailsjs/go/main/App'
import { main, db } from '../wailsjs/go/models'
import * as utils from './components/utils'
import * as types from './components/types'

export const useStore = defineStore('store', () => {
    let isInit = false

    const companies = ref<main.Company[]>([])
    const companyTypes = ref<db.CompanyType[]>([])
    const events = ref<main.Event[]>([])
    const contacts = ref<db.ListContactRow[]>([])
    const eventSource = ref<db.EventSource[]>([])
    const applications = ref<db.ListJobApplicationRow[]>([])
    const applicationStatus = ref<db.JobApplicationStatus[]>([])

    async function loadData(callback: () => Promise<void>, errMessage: string) {
        try {
            await callback()
        } catch (e) {
            console.error(errMessage, e)
            if (e instanceof Error) {
                console.error(e.stack)
            }
            throw e
        }

    }

    async function loadCompanies() {
        await loadData(async () => {
            companies.value = await back.ListCompanies()
        }, "failled to load companies")
    }
    async function loadCompanyTypes() {
        await loadData(async () => {
            companyTypes.value = await back.ListCompanyTypes()
        }, "failled to load companies types")
    }
    async function loadEvents() {
        await loadData(async () => {
            let eventsConv = await back.ListEvents()
            for (const e of eventsConv) {
                e.date = utils.parseBackendDate(e.date)
            }
            events.value = eventsConv
        }, "failled to load events")
    }
    async function loadContact() {
        await loadData(async () => {
            contacts.value = await back.ListContact()
        }, "failled to load contacts")
    }
    async function loadEventSource() {
        await loadData(async () => {
            eventSource.value = await back.ListEventSource()
        }, "failled to load events sources")
    }
    async function loadJobApplication() {
        await loadData(async () => {
            applications.value = await back.ListJobApplication()
        }, "failled to load job applications")
    }
    async function loadJobApplicationStatus() {
        await loadData(async () => {
            applicationStatus.value = await back.ListJobApplicationStatus()
        }, "failled to load job applications status")
    }

    async function init() {
        if (isInit) {
            return
        }
        console.log("init data")
        await Promise.all([
            loadCompanies(),
            loadCompanyTypes(),
            loadEvents(),
            loadContact(),
            loadEventSource(),
            loadJobApplication(),
            loadJobApplicationStatus(),
        ])
        isInit = true
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