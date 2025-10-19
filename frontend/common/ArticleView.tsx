import * as React from 'react'
import { useState } from 'react'
import { Article } from './article'
import { useDict } from '../store'
import { IconClipboard, IconExternal } from '../icons'
import { DictTitle } from './AuthorsDict';
import { useTooltips } from './useTooltips';

type ArticleViewProps = {
    a: Article
    showExternalButton: boolean
    showSource: boolean
}

const delayConfig = `{"show": 1000, "hide": 20}`

const IconExternalController: React.FC<{ a: Article }> = ({ a }) => {
    const el = useTooltips<HTMLAnchorElement>()

    return <a ref={el} href={`/${a.DictionaryID}/${a.ID}`} className="btn btn-link ms-2" target="_blank"
        data-bs-toggle="tooltip" data-bs-title="Адчыніць артыкул асобна" data-bs-delay={delayConfig}>
        <IconExternal />
    </a>
}

const IconCopyLinkController: React.FC<{ a: Article }> = ({ a }) => {
    const [activated, setActivated] = useState<boolean>(false)
    const el = useTooltips<HTMLButtonElement>()

    const onClick = () => {
        const { protocol, host } = window.location
        navigator.clipboard.writeText(`${protocol}//${host}/${a.DictionaryID}/${encodeURIComponent(a.ID)}`)
        setActivated(true)
        window.setTimeout(() => { setActivated(false) }, 1500)
    }

    const iconStyles: React.CSSProperties = {}
    if (activated) {
        iconStyles.color = 'red'
    }

    return <button ref={el} type="button" className="btn btn-link ms-2" style={iconStyles} onClick={onClick}
        data-bs-toggle="tooltip" data-bs-title="Капіраваць простую спасылку на артыкул" data-bs-delay={delayConfig}>
        <IconClipboard type={activated ? 'check' : ''} />
    </button>
}

export const ArticleView: React.FC<ArticleViewProps> = ({ a, showExternalButton, showSource }) => {
    const [dict, _] = useDict(a.DictionaryID)
    const articleRoot = useTooltips<HTMLDivElement>()

    return (
        <div className={`article ${a.DictionaryID}`} ref={articleRoot}>
            <div className="buttons" >
                {showExternalButton && <IconExternalController a={a} />}
                <IconCopyLinkController a={a} />
            </div>
            <div dangerouslySetInnerHTML={{ __html: a.Content }} />
            {showSource && (<div className="source"><p> <DictTitle d={dict} /> </p></div>)}
        </div>
    )
}
