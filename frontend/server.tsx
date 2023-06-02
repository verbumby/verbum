import 'fastestsmallesttextencoderdecoder'
import 'core-js/actual/url'
import 'core-js/actual/url-search-params'

import * as React from 'react'
import { renderToString } from 'react-dom/server'
import { matchPath, StaticRouter, StaticRouterContext } from 'react-router'
import { configureStore } from '@reduxjs/toolkit'
import { Provider } from 'react-redux'
import { Helmet } from "react-helmet";

import { VerbumAPIClientV8Bridge } from './verbum/v8_bridge'
import { App } from './App'
import { rootReducer } from './store'
import { routes } from './routes'

declare global {
    var verbumV8Bridge: (url: string) => any
}

global.verbumClient = new VerbumAPIClientV8Bridge({ bridge: verbumV8Bridge })

export async function render(rawUrl: string) {
    let url = new URL(rawUrl)

    const store = configureStore({
        reducer: rootReducer,
    })

    const promises: Promise<void>[] = []
    routes.some(route => {
        const match = matchPath(url.pathname, route)
        if (match) {
            for (const dataLoader of route.dataLoaders) {
                promises.push(store.dispatch(dataLoader(match, url.searchParams)))
            }
        }
        return match
    })
    await Promise.all(promises)

    const preloadedState = store.getState()

    const routerContext: StaticRouterContext = {}
    const reactRendered = renderToString(
        <Provider store={store}>
            <StaticRouter location={url} context={routerContext}>
                <App />
            </StaticRouter>
        </Provider>
    )
    const helmet = Helmet.renderStatic()

    if (routerContext.url) {
		return {
			Location: routerContext.url
		}
    }

	return{
		StatusCode: routerContext.statusCode ? routerContext.statusCode : -1,
		Title: helmet.title.toString(),
		Meta: helmet.meta.toString(),
		State: JSON.stringify(preloadedState).replace(/</g, '\\u003c'),
		Body: reactRendered,
	}
}

declare global {
    var render: any
}

globalThis.render = render
