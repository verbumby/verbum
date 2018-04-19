export const showMessage = (message, level = 'info') => (dispatch) => {
    const id = Date.now()
    dispatch({
        type: 'MESSAGES/SHOW',
        message: { id, message, level },
    })
    setTimeout(() => {
        dispatch({
            type: 'MESSAGES/DISMISS',
            message: { id },
        })
    }, 1000)
}

export const showInfoMessage = message => showMessage(message, 'info')
export const showWarningMessage = message => showMessage(message, 'warning')
export const showDangerMessage = message => showMessage(message, 'danger')
export const showSuccessMessage = message => showMessage(message, 'success')
export const showPrimaryMessage = message => showMessage(message, 'primary')
