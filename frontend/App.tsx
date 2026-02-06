import * as React from "react"
import { Route, Routes } from "react-router"
import LoadingBarContainer from "react-redux-loading-bar"
import 'bootstrap/dist/css/bootstrap.css'
import './styles.css'

import { Footer, Logo, useURLSearch } from './common'
import { routes } from './routes'
import { Helmet } from "react-helmet"

const App: React.FC = () => {
    const urlSearch = useURLSearch({og_preview: ''})
    const ogPreview = !!urlSearch.get('og_preview')
    return (<>
        <Helmet>
            <meta property="og:image" content="/statics/favicon_squared_200.png" />
        </Helmet>
        <LoadingBarContainer />
        <div className={`content ${ogPreview ? 'og-preview': ''} `}>
            <Logo />
            <Routes>
                {routes.map(r => <Route key={r.path} {...r} />)}
            </Routes>
            {!ogPreview && <Footer />}
        </div>
    </>)
}

export { App }
