import * as React from 'react'
import { createContext, useContext } from 'react'
import { Navigate, To } from "react-router";

export const SetRedirectContext = createContext((_: To): void => {})

export const Redirect: React.FC<{ to: To }> = ({ to }) => {
	const setRedirect = useContext(SetRedirectContext)
	setRedirect(to)
	return <Navigate to={to} />
}
