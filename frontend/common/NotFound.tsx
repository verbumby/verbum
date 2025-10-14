import * as React from 'react'
import { SetStatusCodeContext } from './StatusCodeContext';
import { useContext } from 'react';

export const NotFound: React.FC = () => {
    const setStatusCode = useContext(SetStatusCodeContext)
    setStatusCode(404)
    return (<div>Такой старонкі не існуе.</div>)
}
