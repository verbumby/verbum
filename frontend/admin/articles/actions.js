import { req, assembleURLQuery } from '../utils'

export const fetchList = (urlQuery) => req(`/admin/api/articles${assembleURLQuery(urlQuery)}`, {
    actionPrefix: 'ARTICLES/LIST/FETCH',
    errorMessagePrefix: 'Failed to fetch Articles list',
})

export const leaveList = () => ({ type: 'ARTICLES/LIST/LEAVE' })

export const createRecord = ({ formData }) => req('/admin/api/articles', {
    options: {
        method: 'post',
        body: JSON.stringify(formData),
    },
    actionPrefix: 'ARTICLES/RECORD/CREATE',
    errorMessagePrefix: 'Failed to create article',
    successMessage: 'Article has been created',
})

export const leaveRecord = () => ({ type: 'ARTICLES/RECORD/LEAVE' })

export const fetchRecord = (articleID)  => req(`/admin/api/articles/${articleID}`, {
    actionPrefix: 'ARTICLES/RECORD/FETCH',
    errorMessagePrefix: 'Failed to fetch article',
})

export const updateRecord = ({ formData }) => req(`/admin/api/articles`, {
    options: {
        method: 'post',
        body: JSON.stringify(formData),
    },
    actionPrefix: 'ARTICLES/RECORD/UPDATE',
    errorMessagePrefix: 'Failed to update article',
    successMessage: 'Article has been updated',
})
