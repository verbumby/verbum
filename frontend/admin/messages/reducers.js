const messages = (state = [], action) => {
    switch (action.type) {
        case 'MESSAGES/SHOW':
            return [
                ...state,
                {
                    id: action.message.id,
                    level: action.message.level,
                    message: action.message.message,
                },
            ]
        case 'MESSAGES/DISMISS':
            return state.filter(item => item.id !== action.message.id)
        default:
            return state
    }
}

export default messages
