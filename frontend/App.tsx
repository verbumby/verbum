import * as React from "react"
import { Switch, Route } from "react-router-dom"
import './styles.css'

import { Footer, Logo } from './common'
import { routes } from './routes'

const App: React.VFC = () => (
    <>
        <Logo />
        <Switch>
            {routes.map(r => <Route key={r.path} {...r} />)}
        </Switch>
        <Footer />
    </>
)

export { App }
