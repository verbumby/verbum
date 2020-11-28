import { Dictionary, Article } from '../reducers'

export interface VerbumAPIClient {
    getDictionaries(): Promise<Dictionary[]>
    search(q: string): Promise<Article[]>
}

declare global {
    var verbumClient: VerbumAPIClient
}

export abstract class VerbumAPIClientImpl implements VerbumAPIClient {
    abstract call<T>(path: string): Promise<T>

    async getDictionaries(): Promise<Dictionary[]> {
        return this.call<Dictionary[]>('/api/dictionaries')
    }

    async search(q: string): Promise<Article[]> {
        return this.call<Article[]>('/api/search?q=' + encodeURIComponent(q))
    }
}
