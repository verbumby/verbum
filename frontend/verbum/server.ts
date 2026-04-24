import fetch from 'node-fetch'
import { NotFoundError, VerbumAPIClientImpl } from './client'

type VerbumAPIClientServerOptions = {
    apiURL: string
}

export class VerbumAPIClientServer extends VerbumAPIClientImpl {
    private apiURL: string

    constructor(options: VerbumAPIClientServerOptions) {
        super()
        this.apiURL = options.apiURL
    }

    async call<T>(path: string): Promise<T> {
        const resp = await fetch(this.apiURL + path, { signal: this.signal })
        if (resp.status === 404) {
            throw new NotFoundError()
        }
        return resp.json() as Promise<T>
    }

    async callString(path: string): Promise<string> {
        const resp = await fetch(this.apiURL + path, { signal: this.signal })
        if (resp.status === 404) {
            throw new NotFoundError()
        }
        return resp.text()
    }
}
