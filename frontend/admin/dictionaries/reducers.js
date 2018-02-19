import { combineReducers } from "redux";

const list = (state = null, action) => {
    switch (action.type) {
        case 'DICTIONARIES/LIST/FETCH/FULFILLED':
            return action.data
        case 'DICTIONARIES/LIST/LEAVE':
            return null
        default:
            return state
    }
}

export default combineReducers({
    list,
})
