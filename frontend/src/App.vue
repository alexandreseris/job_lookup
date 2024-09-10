<script lang="ts" setup>
import { VCard, VTabs, VTab, VCardText, VTabsWindow, VTabsWindowItem } from 'vuetify/lib/components/index.mjs';
import { ref, Suspense } from 'vue'
import Company from './components/Company.vue'
import CompanyType from './components/CompanyType.vue'
import Contact from './components/Contact.vue'
import Event from './components/Event.vue'
import EventSource from './components/EventSource.vue'
import JobApplication from './components/JobApplication.vue';
import JobApplicationStatus from './components/JobApplicationStatus.vue';
import Calandar from './components/Calandar.vue';

const tab = ref("companies")
</script>

<template>
  <v-card>
    <v-tabs v-model="tab" bg-color="primary">
      <v-tab value="companies">Companies</v-tab>
      <v-tab value="contacts">Contacts</v-tab>
      <v-tab value="applications">Applications</v-tab>
      <v-tab value="events">Events</v-tab>
      <v-tab value="calandar">Calandar</v-tab>
      <v-tab value="company_types">Company types</v-tab>
      <v-tab value="event_sources">Event sources</v-tab>
      <v-tab value="application_statuses">Application statuses</v-tab>
    </v-tabs>
  </v-card>

  <v-card-text>
    <v-tabs-window v-model="tab">
      <v-tabs-window-item value="companies">
        <Suspense>
          <company></company>
        </Suspense>
      </v-tabs-window-item>
      <v-tabs-window-item value="contacts">
        <Suspense>
          <contact></contact>
        </Suspense>
      </v-tabs-window-item>
      <v-tabs-window-item value="applications">
        <Suspense>
          <job-application></job-application>
        </Suspense>
      </v-tabs-window-item>
      <v-tabs-window-item value="events">
        <Suspense>
          <event></event>
        </Suspense>
      </v-tabs-window-item>
      <v-tabs-window-item value="calandar">
        <Suspense>
          <calandar></calandar>
        </Suspense>
      </v-tabs-window-item>
      <v-tabs-window-item value="company_types">
        <Suspense>
          <company-type></company-type>
        </Suspense>
      </v-tabs-window-item>
      <v-tabs-window-item value="event_sources">
        <Suspense>
          <event-source></event-source>
        </Suspense>
      </v-tabs-window-item>
      <v-tabs-window-item value="application_statuses">
        <Suspense>
          <job-application-status></job-application-status>
        </Suspense>
      </v-tabs-window-item>
    </v-tabs-window>
  </v-card-text>
</template>

<style lang="css">
body {
  font-family: Tahoma, Verdana, Arial;
}
</style>