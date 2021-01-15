export type LetterFilterEntry = {
    URL: string
    Text: string
    Active: boolean
    Title: string
}

export type LetterFilter = {
    DictID: string
    Prefix: string
    Entries: LetterFilterEntry[][]
}
