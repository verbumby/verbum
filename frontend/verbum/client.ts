import { Article, Dict, Suggestion } from '../common'

export interface VerbumAPIClient {
    getDictionaries(): Promise<Dict[]>
    search(q: string): Promise<Article[]>
    suggest(q: string): Promise<Suggestion[]>
    getArticle(dictID: string, articleID: string): Promise<Article>
}

declare global {
    var verbumClient: VerbumAPIClient
}

export abstract class VerbumAPIClientImpl implements VerbumAPIClient {
    abstract call<T>(path: string): Promise<T>

    async getDictionaries(): Promise<Dict[]> {
        return this.call<Dict[]>('/api/dictionaries')
    }

    async search(q: string): Promise<Article[]> {
        return this.call<Article[]>('/api/search?q=' + encodeURIComponent(q))
    }

    async suggest(q: string): Promise<Suggestion[]> {
        return this.call<Suggestion[]>('/api/suggest?q=' + encodeURIComponent(q))
    }

    async getArticle(dictID: string, articleID: string): Promise<Article> {
        return this.call<Article>(`/api/dictionaries/${encodeURIComponent(dictID)}/articles/${encodeURIComponent(articleID)}`)
    }
}
