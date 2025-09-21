import * as React from 'react'
import { renderToString } from 'react-dom/server'
import Koa from 'koa'
import { matchPath, StaticRouter } from 'react-router'
import { configureStore } from '@reduxjs/toolkit'
import { Provider } from 'react-redux'
import { Helmet } from "react-helmet";

import { VerbumAPIClientServer } from './verbum/server'
import { App } from './App'
import { rootReducer } from './store'
import { routes } from './routes'
import { DictsMetadata, dictsSet, SetRedirectContext, SetStatusCodeContext } from './common'
import { sectionsSet } from './common/sections'
import { createPath, To } from 'react-router'

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
        const match = matchPath(route.path, ctx.URL.pathname)
        if (match) {
            for (const dataLoader of route.dataLoaders) {
                promises.push(store.dispatch(dataLoader(match.params, ctx.URL.searchParams)))
            }
        }
        return match
    })
    await Promise.all(promises)

    const preloadedState = store.getState()

    let statusCode: number | undefined
    let to: To | undefined

    const reactRendered = renderToString(
        <Provider store={store}>
            <SetStatusCodeContext.Provider value={sc => statusCode = sc}>
                <SetRedirectContext.Provider value={t => to = t}>
                    <StaticRouter location={ctx.url}>
                        <App />
                    </StaticRouter>
                </SetRedirectContext.Provider>
            </SetStatusCodeContext.Provider>
        </Provider>
    )
    const helmet = Helmet.renderStatic()

    if (to) {
        ctx.status = 301
        ctx.redirect(typeof to === "string" ? to : createPath(to))
        return
    }

    if (statusCode) {
        ctx.response.status = statusCode
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
