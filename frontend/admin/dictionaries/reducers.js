import { combineReducers } from "redux";

const list = (state = [], action) => {
    switch (action.type) {
        case 'DICTIONARIES/LIST/FETCH/FULFILLED':
            return action.data
        case 'DICTIONARIES/LIST/LEAVE':
            return []
        default:
            return state
    }
}

export default combineReducers({
    list,
})
