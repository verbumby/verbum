import * as React from 'react'
import { Link, To } from 'react-router'

type PaginationProps = {
    current: number
    total: number
    pageToURL: (p: number) => To
}

type PageLink = {
    Key: number
    Active: boolean
    Disabled: boolean
    URL: To
    Text: string
}

export const PaginationView: React.FC<PaginationProps> = ({ current, total, pageToURL }) => {
    const links: PageLink[] = []

    const d = 3
    if (Math.max(current - d, 1) > 1) {
        links.push({
            Key: 1,
            URL: pageToURL(1),
            Text: '1',
            Active: false,
            Disabled: false,
        })
    }

    if (Math.max(current - d, 1) > 2) {
        links.push({
            Key: -1,
            URL: {},
            Text: '...',
            Active: false,
            Disabled: true,
        })
    }

    for (let i = Math.max(current - d, 1); i < current; i++) {
        links.push({
            Key: i,
            URL: pageToURL(i),
            Text: `${i}`,
            Active: false,
            Disabled: false,
        })
    }

    links.push({
        Key: current,
        URL: pageToURL(current),
        Text: `${current}`,
        Active: true,
        Disabled: false,
    })

    for (let i = current + 1; i <= Math.min(current + d, total); i++) {
        links.push({
            Key: i,
            URL: pageToURL(i),
            Text: `${i}`,
            Active: false,
            Disabled: false,
        })
    }

    if (Math.min(current + d, total) < total - 1) {
        links.push({
            Key: -2,
            URL: {},
            Text: '...',
            Active: false,
            Disabled: true,
        })
    }

    if (Math.min(current + d, total) < total) {
        links.push({
            Key: total,
            URL: pageToURL(total),
            Text: `${total}`,
            Active: false,
            Disabled: false,
        })
    }

    return (
        <>
            <p />
            <nav aria-label="pagination" >
                <ul className="pagination justify-content-center">
                    {links.map(l => (
                        <li key={l.Key} className={`page-item ${l.Active ? 'active' : ''} ${l.Disabled ? 'disabled' : ''}`}>
                            <Link className="page-link" to={l.URL}>{l.Text}</Link>
                        </li>
                    ))}
                </ul>
            </nav>
        </>
    )
}
