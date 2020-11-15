import * as React from 'react'
import { hydrate } from 'react-dom'
import { BrowserRouter } from 'react-router-dom'

import { App } from './app/App'

hydrate(
    (
        <BrowserRouter>
            <App message="ololo" />
        </BrowserRouter>
    ),
    document.querySelector('body .content'),
)
