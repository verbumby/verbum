import { combineReducers } from "redux";

const list = (state = [], action) => {
    switch (action.type) {
        case 'ARTICLES/LIST/FETCH/FULFILLED':
            return action.data
        case 'ARTICLES/LIST/LEAVE':
            return []
        default:
            return state
    }
}

export default combineReducers({
    list,
})
