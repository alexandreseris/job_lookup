import * as types from '../components/types'

export type Store<Item> = {
    items: Item[];
    columns: types.Columns<Item>;
    syncItems: () => Promise<void>;
    syncWithChildrens: () => Promise<void>;
    syncWithParents: () => Promise<void>;
    select: () => Promise<Item[]>;
    insert: (item: Item) => Promise<Item>;
    update: (item: Item) => Promise<void>;
    delete_: (item: Item) => Promise<void>;
    add: () => void;
    eq: (i1: Item, i2: Item) => boolean;
}