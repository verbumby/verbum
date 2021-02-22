import { Article } from '.'

export type SearchResult = {
    Articles: Article[]
    TermSuggestions: string[]
}
