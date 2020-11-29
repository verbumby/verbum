import { useLocation } from 'react-router-dom'

export function useURLSearch(): URLSearchParams {
    return new URLSearchParams(useLocation().search)
}
