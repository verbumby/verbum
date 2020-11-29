import * as React from 'react'
import { Article } from './article'
import { useDicts } from '../store'
import { IconExternal } from '../icons'

type ArticleViewProps = {
    a: Article
}

export const ArticleView: React.VFC<ArticleViewProps> = ({ a }) => {
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
