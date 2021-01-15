export type URLSearchEntries = Record<string, string|number>

export class URLSearch<Entries extends URLSearchEntries = {}> {

    private defaults: Entries

    private values: Entries

    constructor(defaults: Entries, init?: URLSearchParams) {
        this.defaults = defaults
        this.values = {} as any

        for (const [k, dflt] of Object.entries(defaults)) {
            if (typeof dflt === 'string') {
                const v = (init && init.has(k) ? init.get(k) : dflt) as Entries[Extract<keyof Entries, string>]
                this.values[k as Extract<keyof Entries, string>] = v
            } else if (typeof dflt === 'number') {
                const v = (init && init.has(k) ? parseInt(init.get(k)) : dflt) as Entries[Extract<keyof Entries, string>]
                this.values[k as Extract<keyof Entries, string>] = v
            }
        }
    }

    get<T extends keyof Entries>(k: T): Entries[T] {
        return this.values[k]
    }

    set<T extends keyof Entries>(k: T, v: Entries[T]): URLSearch<Entries> {
        this.values[k] = v
        return this
    }

    reset<T extends keyof Entries>(k: T): URLSearch<Entries> {
        this.values[k] = this.defaults[k]
        return this
    }

    clone(): URLSearch<Entries> {
        const result: URLSearch<Entries> = new URLSearch<Entries>(this.defaults)
        result.values = {...this.values}
        return result
    }

    encode(): string {
        const pairs: string[] = []
        for (const k in this.defaults) {
            if (this.values[k] === this.defaults[k]) {
                continue
            }

            let v: string
            if (typeof this.defaults[k] === 'string') {
                v = this.values[k] as string
            } else if (typeof this.defaults[k] === 'number') {
                v = `${this.values[k]}`
            }
            pairs.push(`${encodeURIComponent(k)}=${encodeURIComponent(v)}`)
        }

        if (pairs.length > 0) {
            return `?${pairs.join('&')}`
        }

        return ''
    }
}
