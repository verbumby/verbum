import * as React from 'react'
import { renderToString } from 'react-dom/server'
import Koa from 'koa'
import { matchPath, StaticRouter, StaticRouterContext } from 'react-router'
import { configureStore } from '@reduxjs/toolkit'
import { Provider } from 'react-redux'
import { Helmet } from "react-helmet";

import { VerbumAPIClientServer } from './verbum/server'
import { App } from './App'
import { rootReducer } from './store'
import { routes } from './routes'
import { DictsMetadata, dictsSet } from './common'
import { sectionsSet } from './common/sections'

globalThis.verbumClient = new VerbumAPIClientServer({ apiURL: 'http://127.0.0.1:8080' })

let indexhtml: string
async function getIndexHTML(): Promise<string> {
    if (indexhtml) {
        return Promise.resolve(indexhtml)
    }

    const r = await verbumClient.getIndexHTML()
    indexhtml = r
    return Promise.resolve(indexhtml)
}

let dictsMetadata: DictsMetadata
async function getDictsMetadata(): Promise<DictsMetadata> {
    if (dictsMetadata) {
        return Promise.resolve(dictsMetadata)
    }

    const r = await verbumClient.getDictionaries()
    dictsMetadata = r
    return Promise.resolve(dictsMetadata)
}

const k = new Koa()
k.use(async ctx => {
    const store = configureStore({
        reducer: rootReducer,
    })
    const dm = await getDictsMetadata()
    store.dispatch(dictsSet(dm.Dicts))
    store.dispatch(sectionsSet(dm.Sections))

    const promises: Promise<void>[] = []
    routes.some(route => {
        const match = matchPath(ctx.URL.pathname, route)
        if (match) {
            for (const dataLoader of route.dataLoaders) {
                promises.push(store.dispatch(dataLoader(match, ctx.URL.searchParams)))
            }
        }
        return match
    })
    await Promise.all(promises)

    const preloadedState = store.getState()

    const routerContext: StaticRouterContext = {}
    const reactRendered = renderToString(
        <Provider store={store}>
            <StaticRouter location={ctx.url} context={routerContext}>
                <App />
            </StaticRouter>
        </Provider>
    )
    const helmet = Helmet.renderStatic()

    if (routerContext.url) {
        ctx.status = 301
        ctx.redirect(routerContext.url)
        return
    }
    if (routerContext.statusCode) {
        ctx.response.status = routerContext.statusCode
    }

    let body = await getIndexHTML()
    body = body.replace('HEAD_TITLE_PLACEHOLDER', helmet.title.toString())
    body = body.replace('HEAD_META_PLACEHOLDER', helmet.meta.toString())
    body = body.replace('PRELOADED_STATE_PLACEHOLDER', JSON.stringify(preloadedState).replace(/</g, '\\u003c'))
    body = body.replace('BODY_PLACEHOLDER', reactRendered)
    ctx.body = body
})

console.log('listening on 127.0.0.1:8079')
k.listen(8079, '127.0.0.1')
