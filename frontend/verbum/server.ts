import { VerbumAPIClientImpl } from './client'
import fetch from 'node-fetch'

type VerbumAPIClientServerOptions = {
    host: string
}

export class VerbumAPIClientServer extends VerbumAPIClientImpl {
    private host: string

    constructor(options: VerbumAPIClientServerOptions) {
        super()
        this.host = options.host
    }

    async call<T>(path: string): Promise<T> {
        const resp = await fetch(this.host + path)
        return resp.json() as Promise<T>
    }
}
