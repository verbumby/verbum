import * as React from 'react'
import { hydrate } from 'react-dom'
import { BrowserRouter } from 'react-router-dom'
import { createStore } from 'redux'
import { Provider } from 'react-redux'
import { devToolsEnhancer } from 'redux-devtools-extension'

import { App } from './app/App'
import { rootReducer, RootState } from './reducers'

declare global {
    interface Window {
        __PRELOADED_STATE__: RootState
        f1: string
    }
}

const preloadedState = window.__PRELOADED_STATE__
delete window.__PRELOADED_STATE__

const store = createStore(rootReducer, preloadedState, devToolsEnhancer({}))

hydrate(
    (
        <Provider store={store}>
            <BrowserRouter>
                <App message="ololo" />
            </BrowserRouter>
        </Provider>
    ),
    document.querySelector('body .content'),
)
