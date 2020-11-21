import { VerbumAPIClientImpl } from './client'

export class VerbumAPIClientBrowser extends VerbumAPIClientImpl {
    async call<T>(path: string): Promise<T> {
        const resp = await fetch(path)
        return resp.json() as Promise<T>
    }
}
