import { VerbumAPIClientImpl } from './client'

export class VerbumAPIClientBrowser extends VerbumAPIClientImpl {
    async call<T>(path: string): Promise<T> {
        const resp = await fetch(path, { signal: this.signal })
        if (resp.status === 404) {
            return Promise.resolve(null)
        }
        return resp.json() as Promise<T>
    }
}
