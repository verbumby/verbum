import * as React from 'react'
import { IconExclamationTriangle } from '../icons'
import { Dict } from './dict'
import { useBSTooltips } from './useBSTooltips'

const warningText = 'Аўтарскі слоўнік — у ім словы і тлумачэнні пададзены паводле асабістых поглядаў укладальнікаў. Магчымыя няправільныя націскі, а таксама іншыя памылкі і недакладнасці.'

export const AuthorsDictWarning: React.FC = () => <>{warningText}</>

export const AuthorsDictWarningIcon: React.FC = () => {
	const el = useBSTooltips()
	return <span ref={el} data-bs-toggle="tooltip" data-bs-title={warningText} tabIndex={0}>
		<IconExclamationTriangle /></span>
}

export const DictTitle: React.FC<{ d: Dict }> = ({ d }) => <>
	{d.Title}
	{d.Authors && <>&nbsp;<AuthorsDictWarningIcon /></>}
</>

