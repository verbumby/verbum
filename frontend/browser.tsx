import { configureStore } from '@reduxjs/toolkit'
import * as React from 'react'
import { hydrate } from 'react-dom'
import { Provider } from 'react-redux'
import { loadingBarMiddleware } from 'react-redux-loading-bar'
import { BrowserRouter } from 'react-router'

import { App } from './App'
import { type RootState, rootReducer } from './store'
import { VerbumAPIClientBrowser } from './verbum/browser'

window.verbumClient = new VerbumAPIClientBrowser()

declare global {
    interface Window {
        __PRELOADED_STATE__?: RootState
    }
}

const preloadedState = window.__PRELOADED_STATE__
if (!preloadedState) {
    throw new Error('__PRELOADED_STATE__ missing')
}
delete window.__PRELOADED_STATE__

const store = configureStore({
    reducer: rootReducer,
    preloadedState: preloadedState,
    middleware: (getDefaultMiddleware) =>
        getDefaultMiddleware().concat(
            loadingBarMiddleware({
                promiseTypeSuffixes: ['KickOff', 'Success', 'Failure'],
            }),
        ),
})

export type AppDispatch = typeof store.dispatch

hydrate(
    <Provider store={store}>
        <BrowserRouter>
            <App />
        </BrowserRouter>
    </Provider>,
    document.querySelector('body .root'),
)
