import * as React from "react"
import './styles.css'

import { Logo } from './Logo'
import { Footer } from './Footer'
import { IndexPage } from './IndexPage'

type AppProps = { message: string }
const App: React.VoidFunctionComponent<AppProps> = ({ message }: AppProps) => (
    <>
        <Logo />
        <div>{message}</div>
        <IndexPage />
        <Footer />
    </>
)

export { App }
