import { Dictionary } from '../reducers'

export interface VerbumAPIClient {
    getDictionaries(): Promise<Dictionary[]>
}

declare global {
    var verbumClient: VerbumAPIClient
}

export abstract class VerbumAPIClientImpl implements VerbumAPIClient {
    abstract call<T>(path: string): Promise<T>

    async getDictionaries(): Promise<Dictionary[]> {
        return this.call<Dictionary[]>('/api/dictionaries')
    }
}
