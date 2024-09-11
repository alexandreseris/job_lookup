import * as types from './types'
import {
    toRaw,
    isRef,
    isReactive,
    isProxy,
} from 'vue';
import momentTz from 'moment-timezone'


export function areObjsEq<T>(item1: T, item2: T): boolean {
    for (const key in item1) {
        const item1Prop = item1[key]
        const item1PropIsArray = Array.isArray(item1Prop)
        const item2Prop = item2[key]
        const item2PropIsArray = Array.isArray(item2Prop)
        if (
            item1PropIsArray && !item2PropIsArray ||
            !item1PropIsArray && item2PropIsArray ||
            (item1PropIsArray && item2PropIsArray && item1Prop.length != item2Prop.length)) {
            return false
        }
        if (item1PropIsArray && item2PropIsArray) {
            for (let index = 0; index < item1Prop.length; index++) {
                if (item1Prop[index] !== item2Prop[index]) {
                    return false
                }
            }
            continue
        }
        if (item1Prop instanceof Date && item2Prop instanceof Date) {
            if (item1Prop.getTime() !== item2Prop.getTime()) {
                return false
            }
            continue
        }
        if (item1Prop !== item2Prop) {
            return false
        }
    }
    return true
}


const dateBackendFormat = "yyyy-MM-DDTHH:mm:ssZ"
const dateInputFormat = "yyyy-MM-DD"
const tz = momentTz.tz.guess()


export function parseBackendDate(date: string | Date): Date {
    if (date instanceof Date) {
        return date
    }
    return momentTz.parseZone(date, dateBackendFormat, true).toDate()
}

export function parseBackendDateOpt(date: string | Date): Date | null {
    if (!date) {
        return null
    }
    return parseBackendDate(date)
}

export function parseInputDate(dateStr: string): Date | null {
    if (!dateStr) {
        return null
    }
    return momentTz.parseZone(dateStr, dateInputFormat, true).tz(tz, true).toDate()
}

export function formatDateForBackend(date: Date | null): string | null {
    if (!date) {
        return null
    }
    return momentTz(date).format(dateBackendFormat)
}


export function formatDateToInput(date: Date | null): string {
    if (!date) {
        return ""
    }
    return momentTz(date).format(dateInputFormat)
}

export function formatDateToLocale(date: Date | null): string {
    if (!date) {
        return ""
    }
    return momentTz(date).toDate().toLocaleDateString()
}

export function itemsMap<T extends types.WithId>(items_: T[]): Map<number, T> {
    let map: Map<number, T> = new Map()
    for (const item of items_) {
        map.set(item.id, item)
    }
    return map
}


export function deepToRaw<T extends Record<string, any>>(sourceObj: T): T {
    const objectIterator = (input: any): any => {
        if (Array.isArray(input)) {
            return input.map((item) => objectIterator(item));
        }
        if (isRef(input) || isReactive(input) || isProxy(input)) {
            return objectIterator(toRaw(input));
        }
        if (input && input instanceof Date) {
            return new Date(input.getTime())
        }
        if (input && typeof input === 'object') {
            return Object.keys(input).reduce((acc, key) => {
                acc[key as keyof typeof acc] = objectIterator(input[key]);
                return acc;
            }, {} as T);
        }
        return input;
    };

    return objectIterator(sourceObj);
}
