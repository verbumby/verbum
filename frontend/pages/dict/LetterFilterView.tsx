import * as React from 'react'
import { Link } from 'react-router-dom'
import { LocationDescriptor } from 'history'
import { LetterFilter } from '../../common'

type LetterFilterViewProps = {
    letterFilter: LetterFilter
    prefixToURL: (prefix: string) => LocationDescriptor
}
export const LetterFilterView: React.VFC<LetterFilterViewProps> = ({ letterFilter, prefixToURL }) => {
    return (
        <div className="letter-filter">
            {letterFilter.map((lfl, i) => (
                <div key={i} className="letter-filter-level mt-1 mb-1">
                    {lfl.map(l => (
                        <Link
                            key={l.Text}
                            className={`btn ml-1 mb-1 btn-sm btn-light ${l.Active ? 'active' : ''}`}
                            title={l.Title}
                            to={prefixToURL(l.URL)}
                        >{l.Text === ' ' ? <>&nbsp;</> : <>{l.Text}</>}</Link>
                    ))}
                </div>
            ))}
        </div>
    )
}
