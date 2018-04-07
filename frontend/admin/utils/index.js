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

export const parseURLSearchParams = (search) => {
    const u = new URLSearchParams(search)
    let result = {}
    for(var pair of u.entries()) {
        result[pair[0]] = pair[1]
    }
    return result
}

export const assembleURLQuery = (params) => {
    const u = new URLSearchParams()
    const defaults = params._defaults || {}
    delete params._defaults
    for (let key of Object.keys(params)) {
        if (key in defaults && defaults[key] === params[key]) {
            continue;
        }
        u.set(key, params[key])
    }
    let result = u.toString()
    if (result.length > 0) {
        result = `?${result}`
    }
    return result
}
