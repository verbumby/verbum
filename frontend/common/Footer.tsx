import * as React from 'react'
import { IconGitHub } from '../icons'
import ThemeSelector from './ThemeSelector'
import { Link } from 'react-router'

export const Footer: React.FunctionComponent = () => (
    <footer className="text-center" style={{ marginTop: 'auto' }}>
        <ThemeSelector />
        {' '}
        <a className="btn btn-link btn-sm text-secondary" target="_blank" href="https://github.com/verbumby">
            <IconGitHub />
        </a>
        {' '}
        <a className="btn btn-link btn-sm text-secondary" href="mailto:vramanenka@gmail.com">
            vramanenka@gmail.com
        </a>
        {' '}
        <Link className='btn btn-link btn-sm text-secondary' to={'/support'}>Падтрымаць</Link>
        {' '}
        <a className="btn btn-link btn-sm text-secondary" target="_blank" href="https://daviedka.bnkorpus.info">
            Моўная даведка
        </a>
    </footer>
)
