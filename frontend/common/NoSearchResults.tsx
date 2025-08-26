import * as React from 'react'
import { Link } from "react-router-dom"
import { LocationDescriptor } from 'history'

type NoSearchResultsProps = {
	q: string
	suggestions: string[]
	calculateSuggestionURL: (q: string) => LocationDescriptor
}

export const NoSearchResults: React.FC<NoSearchResultsProps> = ({ q, suggestions, calculateSuggestionURL }) => {
	return (
		<div className="no-results">
			<p>Па запыце <strong>{q}</strong> нічога не знойдзена.</p>
			{suggestions.length == 1 && (
				<p>Магчыма вы шукалі&nbsp;
					<Link to={calculateSuggestionURL(suggestions[0])}>
						{suggestions[0]}
					</Link>.
				</p>
			)}
			{suggestions.length > 1 && (
				<p>Магчыма вы шукалі:
					<ul>
						{suggestions.map(s => (
							<li key={s}>
								<Link to={calculateSuggestionURL(s)}>{s}</Link>
							</li>
						))}
					</ul>
				</p>
			)}
		</div>
	)
}
