import * as React from 'react'
import { hydrate } from 'react-dom'
import { BrowserRouter } from 'react-router-dom'
import { configureStore } from '@reduxjs/toolkit'
import { Provider } from 'react-redux'

import { App } from './App'
import { rootReducer, RootState } from './reducers'
import { VerbumAPIClientBrowser } from './verbum/browser'

window.verbumClient = new VerbumAPIClientBrowser()

declare global {
    interface Window {
        __PRELOADED_STATE__: RootState
    }
}

const preloadedState = window.__PRELOADED_STATE__
delete window.__PRELOADED_STATE__

const store = configureStore({
    reducer: rootReducer,
    preloadedState: preloadedState,
})

hydrate(
    (
        <Provider store={store}>
            <BrowserRouter>
                <App />
            </BrowserRouter>
        </Provider>
    ),
    document.querySelector('body .content'),
)
