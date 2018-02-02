import { req } from '../utils'

export const fetchList = () => req('/admin/api/dictionaries', {
    actionPrefix: 'DICTIONARIES/LIST/FETCH',
    errorMessagePrefix: 'Failed to fetch Dictionaries list',
})

export const leaveList = () => ({ type: 'DICTIONARIES/LIST/LEAVE' })

export const createDictionary = ({ formData }) => req('/admin/api/dictionaries', {
    options: {
        method: 'post',
        body: JSON.stringify(formData),
    },
    actionPrefix: 'DICTIONARIES/DICTIONARY/CREATE',
    errorMessagePrefix: 'Failed to create the dictionary',
    successMessage: 'Dictionary has been created successfully',
})
