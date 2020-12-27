export type Article = {
    ID: string
    Title: string
    Headword: string[]
    Content: string
    DictionaryID: string
}

export type ArticleList = {
    Articles: Article[]
    Pagination: {
        Current: number
        Total: number
    }
}
