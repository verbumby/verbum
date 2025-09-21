import * as React from "react"
import { Route, Routes } from "react-router"
import LoadingBarContainer from "react-redux-loading-bar"
import 'bootstrap/dist/css/bootstrap.css'
import './styles.css'

import { Footer, Logo } from './common'
import { routes } from './routes'
import { Helmet } from "react-helmet"

const App: React.FC = () => (
    <>
        <Helmet>
            <meta property="og:image" content="/statics/favicon_squared.png" />
        </Helmet>
        <LoadingBarContainer />
        <div className="content">
            <Logo />
            <Routes>
                {routes.map(r => <Route key={r.path} {...r} />)}
            </Routes>
            <Footer />
        </div>
    </>
)

export { App }
