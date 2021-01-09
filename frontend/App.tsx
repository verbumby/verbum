import * as React from "react"
import { Switch, Route } from "react-router-dom"
import LoadingBarContainer from "react-redux-loading-bar"
import './styles.css'

import { Footer, Logo } from './common'
import { routes } from './routes'

const App: React.VFC = () => (
    <>
        <LoadingBarContainer />
        <div className="content">
            <Logo />
            <Switch>
                {routes.map(r => <Route key={r.path} {...r} />)}
            </Switch>
            <Footer />
        </div>
    </>
)

export { App }
