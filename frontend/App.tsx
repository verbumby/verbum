import * as React from "react"
import './styles.css'

import { Logo } from './common/Logo'
import { Footer } from './common/Footer'
import { IndexPage } from './pages/index/IndexPage'

const App: React.VFC = () => (
    <>
        <Logo />
        <IndexPage />
        <Footer />
    </>
)

export { App }
