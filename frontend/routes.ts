import * as React from 'react'
import { match } from 'react-router-dom'
import { articleFetchServer, ArticlePage } from './pages/article'
import { IndexPage, searchServer } from './pages/index/index'
import { AppThunkAction } from './store'
import { dictsFetch } from './common'
import { DictPage, letterFilterFetchServer, dictArticlesFetchServer, abbrFetchServer } from './pages/dict'

type DataLoader = (match: match, urlSearch: URLSearchParams) => AppThunkAction

type Route = {
    path: string,
    component: React.ComponentType<any>,
    dataLoaders: DataLoader[],
}

export const routes: Route[] = [
    {
        path: '/:dictID/:articleID',
        component: ArticlePage,
        dataLoaders: [dictsFetch, articleFetchServer],
    },
    {
        path: '/:dictID',
        component: DictPage,
        dataLoaders: [dictsFetch, letterFilterFetchServer, dictArticlesFetchServer, abbrFetchServer],
    },
    {
        path: '/',
        component: IndexPage,
        dataLoaders: [dictsFetch, searchServer],
    },
]
