import { VerbumAPIClientImpl } from './client'

export class VerbumAPIClientBrowser extends VerbumAPIClientImpl {
    async call<T>(path: string): Promise<T> {
        const resp = await fetch(path, { signal: this.signal })
        return resp.json() as Promise<T>
    }
}
