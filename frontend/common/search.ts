import { Article } from '.'
import { Pagination } from './pagination'

export type SearchResult = {
    Articles: Article[]
    TermSuggestions: string[]
    Pagination: Pagination
}
