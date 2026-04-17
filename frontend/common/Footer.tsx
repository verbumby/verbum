import type * as React from 'react'
import { Link } from 'react-router'
import { IconGitHub } from '../icons/IconGitHub'
import ThemeSelector from './ThemeSelector'

export const Footer: React.FunctionComponent = () => (
    <footer className="text-center" style={{ marginTop: 'auto' }}>
        <ThemeSelector />{' '}
        <a
            className="btn btn-link btn-sm text-secondary"
            target="_blank"
            href="https://github.com/verbumby"
            rel="noopener"
        >
            <IconGitHub />
        </a>{' '}
        <a
            className="btn btn-link btn-sm text-secondary"
            href="mailto:vramanenka@gmail.com"
        >
            vramanenka@gmail.com
        </a>{' '}
        <Link className="btn btn-link btn-sm text-secondary" to={'/support'}>
            Падтрымаць
        </Link>{' '}
        <a
            className="btn btn-link btn-sm text-secondary"
            target="_blank"
            href="https://daviedka.bnkorpus.info"
            rel="noopener"
        >
            Моўная даведка
        </a>
    </footer>
)
