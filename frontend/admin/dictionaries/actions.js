import { req } from '../utils'

export const createDictionary = ({ formData }) => req('/admin/api/dictionaries', {
    options: {
        method: 'post',
        body: JSON.stringify(formData)
    },
    errorMessagePrefix: 'Failed to create the dictionary',
    successMessage: 'Dictionary has been created successfully'
})
