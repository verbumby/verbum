import * as React from 'react'
import { Link, To } from 'react-router'
import { LetterFilter } from '../../common'

type LetterFilterViewProps = {
    letterFilter: LetterFilter
    prefixToURL: (prefix: string) => To
}
export const LetterFilterView: React.FC<LetterFilterViewProps> = ({ letterFilter, prefixToURL }) => {
    return (
        <div className="letter-filter">
            {letterFilter.Entries.map((lfl, i) => (
                <div key={i} className="letter-filter-level mt-1 mb-1">
                    {lfl.map(l => (
                        <Link
                            key={l.Text}
                            className={`btn ms-1 mb-1 btn-sm btn-light ${l.Active ? 'active' : ''}`}
                            title={l.Title}
                            to={prefixToURL(l.URL)}
                        >{l.Text === ' ' ? <>&nbsp;</> : <>{l.Text}</>}</Link>
                    ))}
                </div>
            ))}
        </div>
    )
}
