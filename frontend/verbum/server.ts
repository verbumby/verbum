import { VerbumAPIClientImpl } from './client'
import fetch from 'node-fetch'

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
        const resp = await fetch(this.apiURL + path)
        return resp.json() as Promise<T>
    }
}
