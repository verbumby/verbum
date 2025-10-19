import * as React from 'react'
import { useState, useEffect } from 'react'
import { OverlayTrigger, Tooltip } from "react-bootstrap";
import { Article } from './article'
import { useDict } from '../store'
import { IconClipboard, IconExternal } from '../icons'
import { OverlayInjectedProps } from 'react-bootstrap/esm/Overlay';
import { OverlayDelay } from 'react-bootstrap/esm/OverlayTrigger';
import { DictTitle } from './AuthorsDict';

type ArticleViewProps = {
    a: Article
    showExternalButton: boolean
    showSource: boolean
}

const defaultIconTooltipDelayConfig: OverlayDelay = { show: 1000, hide: 20 }

const IconExternalController: React.FC<{ a: Article }> = ({ a }) => {
    const renderOpenInNewTabTooltip = (props: OverlayInjectedProps) => (
        <Tooltip
            id={`tooltip-open-article-in-new-tab-${a.DictionaryID}-${a.ID}`}
            {...props}
        >Адчыніць артыкул асобна</Tooltip>
    )
    return (<OverlayTrigger overlay={renderOpenInNewTabTooltip} delay={defaultIconTooltipDelayConfig}>
        <a href={`/${a.DictionaryID}/${a.ID}`} className="btn btn-link ms-2" target="_blank">
            <IconExternal />
        </a>
    </OverlayTrigger>)
}

const IconCopyLinkController: React.FC<{ a: Article }> = ({ a }) => {
    const [activated, setActivated] = useState<boolean>(false)
    const renderCopyLinkTooltip = (props: OverlayInjectedProps) => (
        <Tooltip
            id={`tooltip-copy-article-link-${a.DictionaryID}-${a.ID}`}
            {...props}
        >Капіраваць простую спасылку на артыкул</Tooltip>
    )

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

    return (<OverlayTrigger overlay={renderCopyLinkTooltip} delay={defaultIconTooltipDelayConfig}>
        <button type="button" className="btn btn-link ms-2" style={iconStyles} onClick={onClick}>
            <IconClipboard type={activated ? 'check' : ''} />
        </button>
    </OverlayTrigger>)
}

export const ArticleView: React.FC<ArticleViewProps> = ({ a, showExternalButton, showSource }) => {
    const [dict, _] = useDict(a.DictionaryID)
    const [articleRoot, setArticleRoot] = useState(null)
    const [bootstrapAPI, setBootstrapAPI] = useState(null)

    useEffect(() => {
        import('bootstrap').then(setBootstrapAPI)
    }, [])

    useEffect(() => {
        if (!bootstrapAPI || !articleRoot) {
            return
        }

        let ts = new Array()
        for (let e of articleRoot.querySelectorAll('[data-bs-toggle="tooltip"]')) {
            ts.push(new bootstrapAPI.Tooltip(e))
        }

        return () => {
            for (let t of ts) {
                t.dispose()
            }
        }
    }, [bootstrapAPI, articleRoot])

    return (
        <div className={`article ${a.DictionaryID}`} ref={setArticleRoot}>
            <div className="buttons" >
                {showExternalButton && <IconExternalController a={a} />}
                <IconCopyLinkController a={a} />
            </div>
            <div dangerouslySetInnerHTML={{ __html: a.Content }} />
            {showSource && (<div className="source"><p> <DictTitle d={dict} /> </p></div>)}
        </div>
    )
}
