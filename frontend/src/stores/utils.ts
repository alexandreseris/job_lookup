export async function loadData(callback: () => Promise<void>, class_: { name: string }) {
    try {
        await callback()
    } catch (e) {
        console.error(`failled to load ${class_.name}`, e)
        if (e instanceof Error) {
            console.error(e.stack)
        }
        throw e
    }
}