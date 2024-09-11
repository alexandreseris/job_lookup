import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as back from '../wailsjs/go/main/App'
import { main, db } from '../wailsjs/go/models'
import * as utils from './components/utils'
import * as types from './components/types'

export const useStore = defineStore('store', () => {
    let isInit = false

    const companies = ref<main.Company[]>([])
    const companyTypes = ref<db.ListCompanyTypeRow[]>([])
    const events = ref<main.Event[]>([])
    const contacts = ref<main.Contact[]>([])
    const eventSource = ref<db.ListEventSourceRow[]>([])
    const applications = ref<main.JobApplication[]>([])
    const applicationStatus = ref<db.ListJobApplicationStatusRow[]>([])

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
            let companiesConv = await back.ListCompanies()
            for (const e of companiesConv) {
                e.last_event = utils.parseBackendDateOpt(e.last_event)
                e.next_event = utils.parseBackendDateOpt(e.next_event)
            }
            companies.value = companiesConv
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
            let contactConv = await back.ListContact()
            for (const e of contactConv) {
                e.last_event = utils.parseBackendDateOpt(e.last_event)
                e.next_event = utils.parseBackendDateOpt(e.next_event)
            }
            contacts.value = contactConv
        }, "failled to load contacts")
    }
    async function loadEventSource() {
        await loadData(async () => {
            eventSource.value = await back.ListEventSource()
        }, "failled to load events sources")
    }
    async function loadJobApplication() {
        await loadData(async () => {
            let applicationConv = await back.ListJobApplication()
            for (const e of applicationConv) {
                e.last_event = utils.parseBackendDateOpt(e.last_event)
                e.next_event = utils.parseBackendDateOpt(e.next_event)
            }
            applications.value = applicationConv
        }, "failled to load job applications")
    }
    async function loadJobApplicationStatus() {
        await loadData(async () => {
            applicationStatus.value = await back.ListJobApplicationStatus()
        }, "failled to load job applications status")
    }

    async function _init() {
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
    }

    async function init() {
        if (isInit) {
            return
        }
        await _init()
        isInit = true
    }

    async function forceInit() {
        await _init()
    }

    function findCompanyTypeNamesFromCompany(company: main.Company): string[] {
        return companyTypes.value
            .map((e) => { return e.name })
    }
    function findCompanyNamesFromContact(contact: main.Contact): string[] {
        return companies.value
            .map((e) => { return e.name })
    }

    function findCompanyNamesFromApplication(application: main.JobApplication): string[] {
        return companies.value
            .map((e) => { return e.name })
    }

    function findStatusNamesFromApplication(application: main.JobApplication): string[] {
        return applicationStatus.value
            .map((e) => { return e.name })
    }

    function findCompanyNamesFromEvent(event: main.Event): string[] {
        return companies.value
            .map((e) => { return e.name })
    }

    function findSourceNamesFromEvent(event: main.Event): string[] {
        return eventSource.value
            .map((e) => { return e.name })
    }

    function findApplicationNamesFromEvent(event: main.Event): string[] {
        return applications.value
            .filter((e) => { return e.company_name === event.company_name })
            .map((e) => { return e.job_title })
    }

    function findContactNamesFromEvent(event: main.Event): string[] {
        return contacts.value
            .filter((e) => { return e.company_name === event.company_name })
            .map((e) => { return `${e.fist_name}, ${e.last_name}` })
    }

    const calendarEvents = computed(function () {
        return events.value.map(function (e) {
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
        init,
        forceInit,
        companies,
        companyTypes,
        events,
        contacts,
        eventSource,
        applications,
        applicationStatus,
        calendarEvents,

        findCompanyTypeNamesFromCompany,
        findCompanyNamesFromContact,
        findCompanyNamesFromApplication,
        findStatusNamesFromApplication,
        findCompanyNamesFromEvent,
        findSourceNamesFromEvent,
        findContactNamesFromEvent,
        findApplicationNamesFromEvent,
    }
})