import { NotFoundError, VerbumAPIClientImpl } from './client'

export class VerbumAPIClientBrowser extends VerbumAPIClientImpl {
    async call<T>(path: string): Promise<T> {
        const resp = await fetch(path, { signal: this.signal })
        if (resp.status === 404) {
            throw new NotFoundError()
        }
        return resp.json() as Promise<T>
    }

    async callString(path: string): Promise<string> {
        const resp = await fetch(path, { signal: this.signal })
        if (resp.status === 404) {
            throw new NotFoundError()
        }
        return resp.text()
    }
}
