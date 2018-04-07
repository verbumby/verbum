import { req, assembleURLQuery } from '../utils'

export const fetchList = (urlQuery) => req(`/admin/api/tasks${assembleURLQuery(urlQuery)}`, {
    actionPrefix: 'TASKS/LIST/FETCH',
    errorMessagePrefix: 'Failed to fetch Tasks list',
})
