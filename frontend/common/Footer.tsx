import * as React from 'react'
import { IconGitHub } from '../icons'

export const Footer: React.VoidFunctionComponent = () => (
    <footer className="text-center" style={{ marginTop: 'auto' }}>
        <a className="btn btn-link btn-sm text-secondary" target="_blank" href="https://github.com/verbumby/verbum">
            <IconGitHub />
        </a>
        {' '}
        <a className="btn btn-link btn-sm text-secondary" href="mailto:vramanenka@gmail.com">
            vramanenka@gmail.com
        </a>
        {' '}
        <a className="btn btn-link btn-sm text-secondary" target="_blank" href="https://daviedka.bnkorpus.info">
            Моўная даведка Інстытута мовазнаўства
        </a>
    </footer>
)
