import type * as React from 'react'
import { useContext } from 'react'
import { SetStatusCodeContext } from './StatusCodeContext'

export const NotFound: React.FC = () => {
    const setStatusCode = useContext(SetStatusCodeContext)
    setStatusCode(404)
    return <div>Такой старонкі не існуе.</div>
}
