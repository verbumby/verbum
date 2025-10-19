import * as React from 'react'
import { IconExclamationTriangle } from '../icons'
import { Dict } from './dict'
import { useEffect, useRef, useState } from 'react'

const warningText = 'Аўтарскі слоўнік — у ім словы і тлумачэнні пададзены паводле асабістых поглядаў укладальнікаў. Магчымыя няправільныя націскі, а таксама іншыя памылкі і недакладнасці.'

export const AuthorsDictWarning: React.FC = () => <>{warningText}</>

export const AuthorsDictWarningIcon: React.FC = () => {
	const el = useRef()
    const [bootstrapAPI, setBootstrapAPI] = useState(null)
    useEffect(() => {
        import('bootstrap').then(setBootstrapAPI)
    }, [])

	useEffect(() => {
		if (!bootstrapAPI || !el.current) {
			return
		}

		const tt = new bootstrapAPI.Tooltip(el.current)

		return () => tt.dispose()
	}, [bootstrapAPI, el])

	return <span ref={el} data-bs-title={warningText} tabIndex={0}><IconExclamationTriangle /></span>
}

export const DictTitle: React.FC<{d: Dict}> = ({ d }) => <>
	{d.Title}
	{d.Authors && <>&nbsp;<AuthorsDictWarningIcon /></>}
</>

