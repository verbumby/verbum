import * as React from "react"
import './styles.css'

type AppProps = { message: string }
const App = ({ message }: AppProps) => <div>{message}</div>

export { App }
