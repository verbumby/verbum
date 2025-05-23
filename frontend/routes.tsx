import * as React from 'react'
import { match } from 'react-router-dom'
import { articleFetchServer, ArticlePage } from './pages/article'
import { IndexPage, searchServer } from './pages/index/index'
import { AppThunkAction } from './store'
import { dictsFetch } from './common'
import { DictPage, letterFilterFetchServer, dictArticlesFetchServer, abbrFetchServer, prefaceFetchServer } from './pages/dict'
import { SupportPage } from './pages/support/SupportPage'

type DataLoader = (match: match, urlSearch: URLSearchParams) => AppThunkAction

type Route = {
    path: string,
    children: React.ReactElement,
    dataLoaders: DataLoader[],
}

export const routes: Route[] = [
    {
        path: '/:dictID/:articleID',
        children: <ArticlePage />,
        dataLoaders: [dictsFetch, articleFetchServer],
    },
    {
        path: '/support',
        children: <SupportPage />,
        dataLoaders: [dictsFetch],
    },
    {
        path: '/:dictID',
        children: <DictPage />,
        dataLoaders: [dictsFetch, letterFilterFetchServer, dictArticlesFetchServer, abbrFetchServer, prefaceFetchServer],
    },
    {
        path: '/',
        children: <IndexPage />,
        dataLoaders: [dictsFetch, searchServer],
    },
]
