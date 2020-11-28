import * as React from 'react'
import { Article as A, useDicts } from '../reducers'
import { IconExternal } from '../icons'

type ArticleProps = {
    a: A
}

export const Article: React.VFC<ArticleProps> = ({ a }) => {
    const dicts = useDicts()

    return (
        <div className={`article ${a.DictionaryID}`}>
            <div className="float-right">
                <a href={`/${a.DictionaryID}/${a.ID}`} target="_blank" style={{ color: 'darkgray' }}>
                    <IconExternal />
                </a>
            </div>
            <div dangerouslySetInnerHTML={{ __html: a.Content }} />
            <div className="source">
                {dicts.find(d => d.ID === a.DictionaryID).Title}
            </div>
        </div>
    )
}
