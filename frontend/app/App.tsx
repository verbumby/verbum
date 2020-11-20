import * as React from "react"
import './styles.css'

import { Logo } from './Logo'
import { Footer } from './Footer'
import { IndexPage } from './IndexPage'

const App: React.VoidFunctionComponent = () => (
    <>
        <Logo />
        <IndexPage />
        <Footer />
    </>
)

export { App }
