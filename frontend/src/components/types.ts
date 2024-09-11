import { Ref } from 'vue'
import { VTextField, VTextarea, VSelect } from 'vuetify/components'
import { VNumberInput } from 'vuetify/labs/VNumberInput'
import { VDateInput } from 'vuetify/labs/VDateInput'
import DateInput from './DateInput.vue'

export interface WithId {
    id: number
}
type KeyofItem<T> = keyof T
type HeaderType = "string" | "multiline" | "int" | "date" | "action" | "listrel" | "rel"
export type Column<T> = {
    key: KeyofItem<T>;
    title: string;
    sortable?: boolean | undefined;
    type: HeaderType;
    requiered?: boolean | undefined;
    readOnly?: boolean | undefined;
}
export type Columns<T> = Column<T>[]
export type SortBy<T> = { key: KeyofItem<T>; order?: 'asc' | 'desc'; }[]

export type RelationFinder<T> = (current_obj: T) => string[]
export type Relations<T> = { [K in keyof T]?: RelationFinder<T> }
export type ListRelations<T> = { [K in keyof T]?: RelationFinder<T> }

export type LineRender<T> = (item: T) => string[]
export type Getter<T> = () => Ref<T[]>
export type Delete<T> = (current_obj: T) => Promise<void>
export type Update<T> = (current_obj: T) => Promise<void>
export type Insert<T> = (current_obj: T) => Promise<T>

export type Inputs = VTextField | VTextarea | VNumberInput | VDateInput | VSelect | typeof DateInput


export type AlertType = "error" | "success" | "warning" | "info"
export type AlertMessage = {
    text: string
    title: string
    type: AlertType
    displayed: boolean
}

export type VuetifyHeaders = {
    key: string;
    title: string;
    sortable?: boolean | undefined;
}[]

export type CalendarEvent = {
    title: string;
    start: Date;
    end: Date;
    allDay: boolean;
    eventId: number;
    company: string;
    job: string;
    contacts: string[];
}
