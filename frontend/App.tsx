import type * as React from 'react'
import LoadingBarContainer from 'react-redux-loading-bar'
import { Route, Routes } from 'react-router'
import 'bootstrap/dist/css/bootstrap.css'
import './styles.css'

import { Helmet } from 'react-helmet'
import { Footer } from './common/Footer'
import { Logo } from './common/Logo'
import { routes } from './routes'

const App: React.FC = () => (
    <>
        <Helmet>
            <meta
                property="og:image"
                content="/statics/favicon_squared_200.png"
            />
        </Helmet>
        <LoadingBarContainer />
        <div className="content">
            <Logo />
            <Routes>
                {routes.map((r) => (
                    <Route key={r.path} {...r} />
                ))}
            </Routes>
            <Footer />
        </div>
    </>
)

export { App }
