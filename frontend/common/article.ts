export type Article = {
    ID: string
    Title: string
    Headword: string[]
    Content: string
    DictionaryID: string
}

export type ArticleList = {
    DictID: string
    Prefix: string
    Articles: Article[]
    Pagination: {
        Current: number
        Total: number
    }
}
