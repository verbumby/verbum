import { combineReducers } from "redux";

const list = (state = null, action) => {
    switch (action.type) {
        case 'ARTICLES/LIST/FETCH/FULFILLED':
            return action.Data
        case 'ARTICLES/LIST/LEAVE':
            return null
        default:
            return state
    }
}

const record = (state = {}, action) => {
    switch (action.type) {
        case 'ARTICLES/RECORD/FETCH/FULFILLED':
            return { ...state, data: action.Data }
        case 'ARTICLES/RECORD/LEAVE':
            return {}
        default:
            return state
    }
}

export default combineReducers({
    list,
    record,
})
