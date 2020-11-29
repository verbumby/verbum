import * as React from "react"
import './styles.css'

import { Footer, Logo } from './common'
import { IndexPage } from './pages/index'

const App: React.VFC = () => (
    <>
        <Logo />
        <IndexPage />
        <Footer />
    </>
)

export { App }
