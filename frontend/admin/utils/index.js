import { showSuccessMessage, showDangerMessage } from '../messages/actions'

export const ifOK = f => ok => (ok ? f(ok) : false)

export const req = (url, { options, actionPrefix, errorMessagePrefix, successMessage }) => (dispatch) => {
    dispatch({ type: `${actionPrefix}/PENDING` })
    return fetch(url, { ...options, credentials: 'include' })
        .then(response => new Promise((resolve, reject) => {
            if (response.ok) {
                response.json()
                    .then(json => resolve(json))
                    .catch(() => resolve({ data: {} }))
            } else {
                response.text().then((text) => {
                    reject(new Error(text || response.statusText))
                })
            }
        }))
        .then((json) => {
            dispatch({ type: `${actionPrefix}/FULFILLED`, ...json })
            if (successMessage) {
                dispatch(showSuccessMessage(successMessage))
            }
            return json
        })
        .catch((error) => {
            dispatch({ type: `${actionPrefix}/REJECT` })
            if (errorMessagePrefix) {
                dispatch(showDangerMessage(`${errorMessagePrefix}: ${error.message}`))
            }
            console.error(error)
            return false
        })
}
