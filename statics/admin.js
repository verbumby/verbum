/******/ (function(modules) { // webpackBootstrap
/******/ 	// The module cache
/******/ 	var installedModules = {};
/******/
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/
/******/ 		// Check if module is in cache
/******/ 		if(installedModules[moduleId]) {
/******/ 			return installedModules[moduleId].exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = installedModules[moduleId] = {
/******/ 			i: moduleId,
/******/ 			l: false,
/******/ 			exports: {}
/******/ 		};
/******/
/******/ 		// Execute the module function
/******/ 		modules[moduleId].call(module.exports, module, module.exports, __webpack_require__);
/******/
/******/ 		// Flag the module as loaded
/******/ 		module.l = true;
/******/
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/
/******/
/******/ 	// expose the modules object (__webpack_modules__)
/******/ 	__webpack_require__.m = modules;
/******/
/******/ 	// expose the module cache
/******/ 	__webpack_require__.c = installedModules;
/******/
/******/ 	// define getter function for harmony exports
/******/ 	__webpack_require__.d = function(exports, name, getter) {
/******/ 		if(!__webpack_require__.o(exports, name)) {
/******/ 			Object.defineProperty(exports, name, {
/******/ 				configurable: false,
/******/ 				enumerable: true,
/******/ 				get: getter
/******/ 			});
/******/ 		}
/******/ 	};
/******/
/******/ 	// getDefaultExport function for compatibility with non-harmony modules
/******/ 	__webpack_require__.n = function(module) {
/******/ 		var getter = module && module.__esModule ?
/******/ 			function getDefault() { return module['default']; } :
/******/ 			function getModuleExports() { return module; };
/******/ 		__webpack_require__.d(getter, 'a', getter);
/******/ 		return getter;
/******/ 	};
/******/
/******/ 	// Object.prototype.hasOwnProperty.call
/******/ 	__webpack_require__.o = function(object, property) { return Object.prototype.hasOwnProperty.call(object, property); };
/******/
/******/ 	// __webpack_public_path__
/******/ 	__webpack_require__.p = "";
/******/
/******/ 	// Load entry module and return exports
/******/ 	return __webpack_require__(__webpack_require__.s = 8);
/******/ })
/************************************************************************/
/******/ ([
/* 0 */
/***/ (function(module, exports) {

module.exports = React;

/***/ }),
/* 1 */
/***/ (function(module, exports) {

module.exports = ReactRedux;

/***/ }),
/* 2 */
/***/ (function(module, exports) {

module.exports = Redux;

/***/ }),
/* 3 */
/***/ (function(module, exports) {

module.exports = ReactRouterDOM;

/***/ }),
/* 4 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__messages_actions__ = __webpack_require__(14);
var _extends = Object.assign || function (target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i]; for (var key in source) { if (Object.prototype.hasOwnProperty.call(source, key)) { target[key] = source[key]; } } } return target; };



const ifOK = f => ok => ok ? f(ok) : false;
/* harmony export (immutable) */ __webpack_exports__["b"] = ifOK;


const req = (url, { options, actionPrefix, errorMessagePrefix, successMessage }) => dispatch => {
    dispatch({ type: `${actionPrefix}/PENDING` });
    return fetch(url, _extends({}, options, { credentials: 'include' })).then(response => new Promise((resolve, reject) => {
        if (response.ok) {
            response.json().then(json => resolve(json)).catch(() => resolve({ data: {} }));
        } else {
            response.text().then(text => {
                reject(new Error(text || response.statusText));
            });
        }
    })).then(json => {
        dispatch(_extends({ type: `${actionPrefix}/FULFILLED` }, json));
        if (successMessage) {
            dispatch(Object(__WEBPACK_IMPORTED_MODULE_0__messages_actions__["b" /* showSuccessMessage */])(successMessage));
        }
        return json;
    }).catch(error => {
        dispatch({ type: `${actionPrefix}/REJECT` });
        if (errorMessagePrefix) {
            dispatch(Object(__WEBPACK_IMPORTED_MODULE_0__messages_actions__["a" /* showDangerMessage */])(`${errorMessagePrefix}: ${error.message}`));
        }
        console.error(error);
        return false;
    });
};
/* harmony export (immutable) */ __webpack_exports__["c"] = req;


const parseURLSearchParams = search => {
    const u = new URLSearchParams(search);
    let result = {};
    for (var pair of u.entries()) {
        result[pair[0]] = pair[1];
    }
    return result;
};
/* unused harmony export parseURLSearchParams */


const assembleURLQuery = params => {
    const u = new URLSearchParams();
    const defaults = params._defaults || {};
    delete params._defaults;
    for (let key of Object.keys(params)) {
        if (key in defaults && defaults[key] === params[key]) {
            continue;
        }
        u.set(key, params[key]);
    }
    let result = u.toString();
    if (result.length > 0) {
        result = `?${result}`;
    }
    return result;
};
/* harmony export (immutable) */ __webpack_exports__["a"] = assembleURLQuery;


/***/ }),
/* 5 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__utils__ = __webpack_require__(4);


const fetchList = urlQuery => Object(__WEBPACK_IMPORTED_MODULE_0__utils__["c" /* req */])(`/admin/api/articles${Object(__WEBPACK_IMPORTED_MODULE_0__utils__["a" /* assembleURLQuery */])(urlQuery)}`, {
    actionPrefix: 'ARTICLES/LIST/FETCH',
    errorMessagePrefix: 'Failed to fetch Articles list'
});
/* harmony export (immutable) */ __webpack_exports__["b"] = fetchList;


const leaveList = () => ({ type: 'ARTICLES/LIST/LEAVE' });
/* harmony export (immutable) */ __webpack_exports__["d"] = leaveList;


const createRecord = ({ formData }) => Object(__WEBPACK_IMPORTED_MODULE_0__utils__["c" /* req */])('/admin/api/articles', {
    options: {
        method: 'post',
        body: JSON.stringify(formData)
    },
    actionPrefix: 'ARTICLES/RECORD/CREATE',
    errorMessagePrefix: 'Failed to create article',
    successMessage: 'Article has been created'
});
/* harmony export (immutable) */ __webpack_exports__["a"] = createRecord;


const leaveRecord = () => ({ type: 'ARTICLES/RECORD/LEAVE' });
/* harmony export (immutable) */ __webpack_exports__["e"] = leaveRecord;


const fetchRecord = articleID => Object(__WEBPACK_IMPORTED_MODULE_0__utils__["c" /* req */])(`/admin/api/articles/${articleID}`, {
    actionPrefix: 'ARTICLES/RECORD/FETCH',
    errorMessagePrefix: 'Failed to fetch article'
});
/* harmony export (immutable) */ __webpack_exports__["c"] = fetchRecord;


const updateRecord = ({ formData }) => Object(__WEBPACK_IMPORTED_MODULE_0__utils__["c" /* req */])(`/admin/api/articles`, {
    options: {
        method: 'post',
        body: JSON.stringify(formData)
    },
    actionPrefix: 'ARTICLES/RECORD/UPDATE',
    errorMessagePrefix: 'Failed to update article',
    successMessage: 'Article has been updated'
});
/* harmony export (immutable) */ __webpack_exports__["f"] = updateRecord;


/***/ }),
/* 6 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__utils__ = __webpack_require__(4);


const fetchList = () => Object(__WEBPACK_IMPORTED_MODULE_0__utils__["c" /* req */])('/admin/api/dictionaries', {
    actionPrefix: 'DICTIONARIES/LIST/FETCH',
    errorMessagePrefix: 'Failed to fetch Dictionaries list'
});
/* harmony export (immutable) */ __webpack_exports__["b"] = fetchList;


const leaveList = () => ({ type: 'DICTIONARIES/LIST/LEAVE' });
/* harmony export (immutable) */ __webpack_exports__["c"] = leaveList;


const createDictionary = ({ formData }) => Object(__WEBPACK_IMPORTED_MODULE_0__utils__["c" /* req */])('/admin/api/dictionaries', {
    options: {
        method: 'post',
        body: JSON.stringify(formData)
    },
    actionPrefix: 'DICTIONARIES/DICTIONARY/CREATE',
    errorMessagePrefix: 'Failed to create the dictionary',
    successMessage: 'Dictionary has been created successfully'
});
/* harmony export (immutable) */ __webpack_exports__["a"] = createDictionary;


/***/ }),
/* 7 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0_react__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_redux__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_redux___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_1_react_redux__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_react_router_dom__ = __webpack_require__(3);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_react_router_dom___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_react_router_dom__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__textarea__ = __webpack_require__(20);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__form_task__ = __webpack_require__(28);
var _extends = Object.assign || function (target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i]; for (var key in source) { if (Object.prototype.hasOwnProperty.call(source, key)) { target[key] = source[key]; } } } return target; };








class Form extends __WEBPACK_IMPORTED_MODULE_0_react___default.a.Component {
    constructor(props) {
        super(props);
        this.state = props.formData;

        this.handleInputChange = this.handleInputChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.toggleTask = this.toggleTask.bind(this);
    }

    handleInputChange(event) {
        const target = event.target;
        let value = target.type === 'checkbox' ? target.checked : target.value;
        const name = target.name;

        switch (name) {
            case "DictID":
                value = parseInt(value);
                break;
        }

        this.setState({
            [name]: value
        });
    }

    toggleTask(n) {
        let task = this.state.Tasks[n];
        task = _extends({}, task, {
            Status: task.Status == 'PENDING' ? 'DONE' : 'PENDING'
        });

        const Tasks = [...this.state.Tasks];
        Tasks[n] = task;

        this.setState({ Tasks });
    }

    handleSubmit(event) {
        event.preventDefault();
        this.props.onSave({ formData: this.state });
    }

    render() {
        return __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
            'form',
            { onSubmit: this.handleSubmit },
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                'div',
                { 'class': 'columns' },
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'div',
                    { 'class': 'column' },
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                        'div',
                        { 'class': 'field' },
                        __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                            'label',
                            { 'class': 'label' },
                            'Dictionary'
                        ),
                        __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                            'div',
                            { 'class': 'control' },
                            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                                'div',
                                { 'class': 'select' },
                                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                                    'select',
                                    { name: 'DictID', value: this.state.DictID, onChange: this.handleInputChange, required: true },
                                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement('option', null),
                                    this.props.dicts.map(d => __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                                        'option',
                                        { value: d.ID },
                                        d.Title
                                    ))
                                )
                            )
                        )
                    )
                ),
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'div',
                    { 'class': 'column' },
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                        'div',
                        { 'class': 'field' },
                        __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                            'label',
                            { 'class': 'label' },
                            'Tasks'
                        ),
                        this.state.Tasks.map((task, index) => {
                            return __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(__WEBPACK_IMPORTED_MODULE_4__form_task__["a" /* default */], { onToggle: () => {
                                    this.toggleTask(index);
                                }, task: task, index: index + 1 });
                        })
                    )
                )
            ),
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                'div',
                { className: 'field' },
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'div',
                    { className: 'control' },
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                        'label',
                        { 'class': 'label' },
                        'Content'
                    ),
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(__WEBPACK_IMPORTED_MODULE_3__textarea__["a" /* default */], {
                        className: 'textarea',
                        type: 'text',
                        name: 'Content',
                        rows: '15',
                        value: this.state.Content,
                        onChange: this.handleInputChange,
                        onSave: () => {
                            this.submitButton.click();
                        }
                    })
                )
            ),
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                'div',
                { className: 'field is-grouped' },
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'p',
                    { 'class': 'control' },
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                        'button',
                        { className: 'button is-link', type: 'submit',
                            ref: button => this.submitButton = button },
                        'Save'
                    )
                ),
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'p',
                    { 'class': 'control' },
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                        __WEBPACK_IMPORTED_MODULE_2_react_router_dom__["Link"],
                        { className: 'button', to: this.props.onCancelRedirectTo },
                        'Cancel'
                    )
                )
            )
        );
    }
}

/* harmony default export */ __webpack_exports__["a"] = (Object(__WEBPACK_IMPORTED_MODULE_1_react_redux__["connect"])(state => ({
    dicts: state.config.Dicts
}))(Form));

/***/ }),
/* 8 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
Object.defineProperty(__webpack_exports__, "__esModule", { value: true });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0_react__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_dom__ = __webpack_require__(9);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_dom___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_1_react_dom__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_react_redux__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_react_redux___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_react_redux__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_redux__ = __webpack_require__(2);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_redux___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_redux__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_redux_thunk__ = __webpack_require__(10);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_redux_thunk___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_4_redux_thunk__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5_react_router_dom__ = __webpack_require__(3);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5_react_router_dom___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_5_react_router_dom__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_6__core_components_top_nav__ = __webpack_require__(11);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_7__dictionaries_components_index__ = __webpack_require__(12);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_8__articles_components_index__ = __webpack_require__(17);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_9__home_components_home_page__ = __webpack_require__(23);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_10__messages_messages__ = __webpack_require__(24);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_11__messages_reducers__ = __webpack_require__(25);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_12__dictionaries_reducers__ = __webpack_require__(26);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_13__articles_reducers__ = __webpack_require__(27);

















const rootReducer = Object(__WEBPACK_IMPORTED_MODULE_3_redux__["combineReducers"])({
    user: (state = {}) => state,
    config: (state = {}) => state,
    messages: __WEBPACK_IMPORTED_MODULE_11__messages_reducers__["a" /* default */],
    dictionaries: __WEBPACK_IMPORTED_MODULE_12__dictionaries_reducers__["a" /* default */],
    articles: __WEBPACK_IMPORTED_MODULE_13__articles_reducers__["a" /* default */]
});

const store = Object(__WEBPACK_IMPORTED_MODULE_3_redux__["createStore"])(rootReducer, window.verbumInitData, Object(__WEBPACK_IMPORTED_MODULE_3_redux__["applyMiddleware"])(__WEBPACK_IMPORTED_MODULE_4_redux_thunk___default.a));

__WEBPACK_IMPORTED_MODULE_1_react_dom___default.a.render(__WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
    __WEBPACK_IMPORTED_MODULE_2_react_redux__["Provider"],
    { store: store },
    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
        __WEBPACK_IMPORTED_MODULE_5_react_router_dom__["BrowserRouter"],
        { basename: '/admin' },
        __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
            'div',
            null,
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                'div',
                null,
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(__WEBPACK_IMPORTED_MODULE_6__core_components_top_nav__["a" /* default */], null),
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(__WEBPACK_IMPORTED_MODULE_10__messages_messages__["a" /* default */], null)
            ),
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                'section',
                { className: 'section' },
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    __WEBPACK_IMPORTED_MODULE_5_react_router_dom__["Switch"],
                    null,
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(__WEBPACK_IMPORTED_MODULE_5_react_router_dom__["Route"], { exact: true, path: '/', component: __WEBPACK_IMPORTED_MODULE_9__home_components_home_page__["a" /* default */] }),
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(__WEBPACK_IMPORTED_MODULE_5_react_router_dom__["Route"], { path: '/dictionaries', component: __WEBPACK_IMPORTED_MODULE_7__dictionaries_components_index__["a" /* default */] }),
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(__WEBPACK_IMPORTED_MODULE_5_react_router_dom__["Route"], { path: '/articles', component: __WEBPACK_IMPORTED_MODULE_8__articles_components_index__["a" /* default */] })
                )
            )
        )
    )
), document.getElementById('root'));

/***/ }),
/* 9 */
/***/ (function(module, exports) {

module.exports = ReactDOM;

/***/ }),
/* 10 */
/***/ (function(module, exports) {

module.exports = ReduxThunk;

/***/ }),
/* 11 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0_react__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_redux__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_redux___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_1_react_redux__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_react_router_dom__ = __webpack_require__(3);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_react_router_dom___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_react_router_dom__);




class TopNav extends __WEBPACK_IMPORTED_MODULE_0_react___default.a.Component {
    render() {
        return __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
            'nav',
            { className: 'navbar is-link' },
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                'div',
                { className: 'navbar-brand' },
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    __WEBPACK_IMPORTED_MODULE_2_react_router_dom__["Link"],
                    { className: 'navbar-item has-text-weight-bold is-size-4', to: '/' },
                    'Verbum'
                )
            ),
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                'div',
                { className: 'navbar-menu' },
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'div',
                    { className: 'navbar-start' },
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                        NavLink,
                        { path: '/dictionaries' },
                        'Dictionaries'
                    ),
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                        NavLink,
                        { path: '/articles' },
                        'Articles'
                    )
                ),
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'div',
                    { className: 'navbar-end' },
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                        'div',
                        { className: 'navbar-item' },
                        this.props.user.name
                    )
                )
            )
        );
    }
}

const NavLink = ({ path, children }) => __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(__WEBPACK_IMPORTED_MODULE_2_react_router_dom__["Route"], { path: path, children: ({ match }) => __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
        __WEBPACK_IMPORTED_MODULE_2_react_router_dom__["Link"],
        { className: `navbar-item ${match ? 'is-active' : ''}`, to: path },
        children
    ) });

/* harmony default export */ __webpack_exports__["a"] = (Object(__WEBPACK_IMPORTED_MODULE_2_react_router_dom__["withRouter"])(Object(__WEBPACK_IMPORTED_MODULE_1_react_redux__["connect"])(store => ({
    user: store.user
}))(TopNav)));

/***/ }),
/* 12 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0_react__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_router_dom__ = __webpack_require__(3);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_router_dom___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_1_react_router_dom__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__list_page__ = __webpack_require__(13);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__new_page__ = __webpack_require__(15);






class Index extends __WEBPACK_IMPORTED_MODULE_0_react___default.a.Component {
    render() {
        const { path } = this.props.match;
        return __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
            'div',
            { className: 'container is-fluid' },
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                __WEBPACK_IMPORTED_MODULE_1_react_router_dom__["Switch"],
                null,
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(__WEBPACK_IMPORTED_MODULE_1_react_router_dom__["Route"], { exact: true, path: path, component: __WEBPACK_IMPORTED_MODULE_2__list_page__["a" /* default */] }),
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(__WEBPACK_IMPORTED_MODULE_1_react_router_dom__["Route"], { path: `${path}/new`, component: __WEBPACK_IMPORTED_MODULE_3__new_page__["a" /* default */] })
            )
        );
    }
}
/* harmony export (immutable) */ __webpack_exports__["a"] = Index;


/***/ }),
/* 13 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0_react__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_router_dom__ = __webpack_require__(3);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_router_dom___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_1_react_router_dom__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_react_redux__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_react_redux___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_react_redux__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_redux__ = __webpack_require__(2);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_redux___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_redux__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__actions__ = __webpack_require__(6);







class ListPage extends __WEBPACK_IMPORTED_MODULE_0_react___default.a.Component {
    componentWillMount() {
        this.props.fetchList();
    }

    componentWillUnmount() {
        this.props.leaveList();
    }

    render() {
        const { url } = this.props.match;
        return __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
            'div',
            null,
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                'div',
                { className: 'level' },
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'div',
                    { className: 'level-left' },
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                        'h1',
                        { className: 'level-item title' },
                        'Dictionaries'
                    )
                ),
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'div',
                    { className: 'level-right' },
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                        'div',
                        { className: 'field is-grouped level-item' },
                        __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                            'p',
                            { 'class': 'control' },
                            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                                __WEBPACK_IMPORTED_MODULE_1_react_router_dom__["Link"],
                                { to: `${url}/new`, className: 'button is-link' },
                                'New Dictionary'
                            )
                        )
                    )
                )
            ),
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement('hr', null),
            this.props.data && __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                'table',
                { className: 'table is-hoverable is-fullwidth' },
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'thead',
                    null,
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                        'tr',
                        null,
                        __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                            'th',
                            null,
                            'Title'
                        )
                    )
                ),
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'tbody',
                    null,
                    this.props.data.map(item => __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                        'tr',
                        null,
                        __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                            'td',
                            null,
                            item.Title
                        )
                    ))
                )
            )
        );
    }
}

/* harmony default export */ __webpack_exports__["a"] = (Object(__WEBPACK_IMPORTED_MODULE_2_react_redux__["connect"])(state => ({
    data: state.dictionaries.list
}), dispatch => Object(__WEBPACK_IMPORTED_MODULE_3_redux__["bindActionCreators"])({ fetchList: __WEBPACK_IMPORTED_MODULE_4__actions__["b" /* fetchList */], leaveList: __WEBPACK_IMPORTED_MODULE_4__actions__["c" /* leaveList */] }, dispatch))(ListPage));

/***/ }),
/* 14 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
const showMessage = (message, level = 'info') => dispatch => {
    const id = Date.now();
    dispatch({
        type: 'MESSAGES/SHOW',
        message: { id, message, level }
    });
    setTimeout(() => {
        dispatch({
            type: 'MESSAGES/DISMISS',
            message: { id }
        });
    }, 5000);
};
/* unused harmony export showMessage */


const showInfoMessage = message => showMessage(message, 'info');
/* unused harmony export showInfoMessage */

const showWarningMessage = message => showMessage(message, 'warning');
/* unused harmony export showWarningMessage */

const showDangerMessage = message => showMessage(message, 'danger');
/* harmony export (immutable) */ __webpack_exports__["a"] = showDangerMessage;

const showSuccessMessage = message => showMessage(message, 'success');
/* harmony export (immutable) */ __webpack_exports__["b"] = showSuccessMessage;

const showPrimaryMessage = message => showMessage(message, 'primary');
/* unused harmony export showPrimaryMessage */


/***/ }),
/* 15 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0_react__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_redux__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_redux___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_1_react_redux__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_redux__ = __webpack_require__(2);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_redux___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_redux__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__form__ = __webpack_require__(16);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__actions__ = __webpack_require__(6);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5__utils__ = __webpack_require__(4);








class NewPage extends __WEBPACK_IMPORTED_MODULE_0_react___default.a.Component {
    render() {
        const onSave = data => {
            this.props.createDictionary(data).then(Object(__WEBPACK_IMPORTED_MODULE_5__utils__["b" /* ifOK */])(data => this.props.history.push('/dictionaries')));
        };
        return __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
            'div',
            null,
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                'h1',
                { className: 'title' },
                'Create New Dictionary'
            ),
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(__WEBPACK_IMPORTED_MODULE_3__form__["a" /* default */], {
                formData: {},
                onSave: onSave,
                onCancelRedirectTo: '/dictionaries'
            })
        );
    }
}

/* harmony default export */ __webpack_exports__["a"] = (Object(__WEBPACK_IMPORTED_MODULE_1_react_redux__["connect"])(state => {}, dispatch => Object(__WEBPACK_IMPORTED_MODULE_2_redux__["bindActionCreators"])({ createDictionary: __WEBPACK_IMPORTED_MODULE_4__actions__["a" /* createDictionary */] }, dispatch))(NewPage));

/***/ }),
/* 16 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0_react__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_router_dom__ = __webpack_require__(3);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_router_dom___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_1_react_router_dom__);



class Form extends __WEBPACK_IMPORTED_MODULE_0_react___default.a.Component {
    constructor(props) {
        super(props);
        this.state = props.formData;

        this.handleInputChange = this.handleInputChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleInputChange(event) {
        const target = event.target;
        let value = target.type === 'checkbox' ? target.checked : target.value;
        const name = target.name;

        this.setState({
            [name]: value
        });
    }

    handleSubmit(event) {
        event.preventDefault();
        this.props.onSave({ formData: this.state });
    }

    render() {
        return __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
            'form',
            { onSubmit: this.handleSubmit },
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                'div',
                { 'class': 'field' },
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'label',
                    { 'class': 'label' },
                    'Title'
                ),
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'div',
                    { 'class': 'control' },
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement('input', { className: 'input',
                        type: 'text',
                        name: 'Title',
                        value: this.state.Title,
                        onChange: this.handleInputChange,
                        required: true
                    })
                )
            ),
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                'div',
                { className: 'field is-grouped' },
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'p',
                    { 'class': 'control' },
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                        'button',
                        { className: 'button is-link', type: 'submit' },
                        'Save'
                    )
                ),
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'p',
                    { 'class': 'control' },
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                        __WEBPACK_IMPORTED_MODULE_1_react_router_dom__["Link"],
                        { className: 'button', to: this.props.onCancelRedirectTo },
                        'Cancel'
                    )
                )
            )
        );
    }
}
/* harmony export (immutable) */ __webpack_exports__["a"] = Form;


/***/ }),
/* 17 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0_react__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_router_dom__ = __webpack_require__(3);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_router_dom___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_1_react_router_dom__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__list_page__ = __webpack_require__(18);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__new_page__ = __webpack_require__(19);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__edit_page__ = __webpack_require__(22);







class Index extends __WEBPACK_IMPORTED_MODULE_0_react___default.a.Component {
    render() {
        const { path } = this.props.match;
        return __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
            'div',
            { className: 'container is-fluid' },
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                __WEBPACK_IMPORTED_MODULE_1_react_router_dom__["Switch"],
                null,
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(__WEBPACK_IMPORTED_MODULE_1_react_router_dom__["Route"], { exact: true, path: path, component: __WEBPACK_IMPORTED_MODULE_2__list_page__["a" /* default */] }),
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(__WEBPACK_IMPORTED_MODULE_1_react_router_dom__["Route"], { path: `${path}/new`, component: __WEBPACK_IMPORTED_MODULE_3__new_page__["a" /* default */] }),
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(__WEBPACK_IMPORTED_MODULE_1_react_router_dom__["Route"], { path: `${path}/:articleID/edit`, component: __WEBPACK_IMPORTED_MODULE_4__edit_page__["a" /* default */] })
            )
        );
    }
}
/* harmony export (immutable) */ __webpack_exports__["a"] = Index;


/***/ }),
/* 18 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0_react__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_router_dom__ = __webpack_require__(3);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_router_dom___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_1_react_router_dom__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_react_redux__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_react_redux___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_react_redux__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_redux__ = __webpack_require__(2);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_redux___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_redux__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__list_page_filter_dictionary__ = __webpack_require__(29);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5__actions__ = __webpack_require__(5);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_6__utils__ = __webpack_require__(4);









class ListPage extends __WEBPACK_IMPORTED_MODULE_0_react___default.a.Component {
    constructor(props) {
        super(props);
        this.state = {
            offset: 0,
            limit: 20,
            filter$DictID: -1,
            filter$TitlePrefix: ''
        };
    }

    componentWillMount() {
        this.fetchList();
    }

    componentWillUnmount() {
        this.props.leaveList();
    }

    fetchList() {
        this.props.fetchList({
            offset: this.state.offset,
            limit: this.state.limit,
            filter$DictID: this.state.filter$DictID,
            filter$TitlePrefix: this.state.filter$TitlePrefix,
            _defaults: {
                offset: 0,
                limit: 20,
                filter$DictID: -1,
                filter$TitlePrefix: ''
            }
        });
    }

    setFilterState(state) {
        this.setState(state, () => {
            this.fetchList();
        });
    }

    render() {
        const { url } = this.props.match;
        const onPrevPageClick = () => {
            const offset = Math.max(0, this.state.offset - this.state.limit);
            this.setState({ offset }, () => {
                this.fetchList();
            });
        };

        const onNextPageClick = () => {
            const offset = this.state.offset + this.state.limit;
            this.setState({ offset }, () => {
                this.fetchList();
            });
        };

        return __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
            'div',
            null,
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                'div',
                { className: 'level' },
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'div',
                    { className: 'level-left' },
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                        'h1',
                        { className: 'level-item title' },
                        'Articles'
                    )
                ),
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'div',
                    { className: 'level-right' },
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                        'div',
                        { className: 'field is-grouped level-item' },
                        __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                            'p',
                            { 'class': 'control' },
                            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                                __WEBPACK_IMPORTED_MODULE_1_react_router_dom__["Link"],
                                { to: `${url}/new`, className: 'button is-link' },
                                'New Article'
                            )
                        )
                    )
                )
            ),
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement('hr', null),
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                'div',
                { 'class': 'field is-grouped' },
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'p',
                    { 'class': 'control' },
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement('input', {
                        'class': 'input',
                        type: 'text',
                        value: this.state.filter$TitlePrefix,
                        onChange: ev => this.setFilterState({ filter$TitlePrefix: ev.target.value })
                    })
                ),
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'p',
                    { 'class': 'control' },
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(__WEBPACK_IMPORTED_MODULE_4__list_page_filter_dictionary__["a" /* default */], { value: this.state.filter$DictID, onChange: filter$DictID => this.setFilterState({ filter$DictID }) })
                )
            ),
            this.props.data && __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                'div',
                null,
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'table',
                    { className: 'table is-hoverable is-fullwidth' },
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                        'thead',
                        null,
                        __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                            'tr',
                            null,
                            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                                'th',
                                { width: '1px' },
                                'ID'
                            ),
                            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                                'th',
                                null,
                                'Title'
                            ),
                            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement('th', { width: '1px' })
                        )
                    ),
                    __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                        'tbody',
                        null,
                        this.props.data.map(item => __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                            'tr',
                            null,
                            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                                'td',
                                null,
                                item.ID
                            ),
                            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                                'td',
                                null,
                                item.Title
                            ),
                            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                                'td',
                                null,
                                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                                    __WEBPACK_IMPORTED_MODULE_1_react_router_dom__["Link"],
                                    { to: `${url}/${item.ID}/edit`, className: 'button' },
                                    'Edit'
                                )
                            )
                        ))
                    )
                )
            ),
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                'nav',
                { 'class': 'pagination', role: 'navigation', 'aria-label': 'pagination' },
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'a',
                    { 'class': 'pagination-previous', onClick: onPrevPageClick },
                    'Previous'
                ),
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'a',
                    { 'class': 'pagination-next', onClick: onNextPageClick },
                    'Next page'
                )
            )
        );
    }
}

/* harmony default export */ __webpack_exports__["a"] = (Object(__WEBPACK_IMPORTED_MODULE_2_react_redux__["connect"])(state => ({
    data: state.articles.list
}), dispatch => Object(__WEBPACK_IMPORTED_MODULE_3_redux__["bindActionCreators"])({ fetchList: __WEBPACK_IMPORTED_MODULE_5__actions__["b" /* fetchList */], leaveList: __WEBPACK_IMPORTED_MODULE_5__actions__["d" /* leaveList */] }, dispatch))(ListPage));

/***/ }),
/* 19 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0_react__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_redux__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_redux___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_1_react_redux__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_redux__ = __webpack_require__(2);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_redux___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_redux__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__form__ = __webpack_require__(7);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__actions__ = __webpack_require__(5);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5__utils__ = __webpack_require__(4);








class NewPage extends __WEBPACK_IMPORTED_MODULE_0_react___default.a.Component {
    render() {
        const onSave = data => {
            this.props.createRecord(data).then(Object(__WEBPACK_IMPORTED_MODULE_5__utils__["b" /* ifOK */])(data => this.props.history.push('/articles')));
        };
        return __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
            'div',
            null,
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                'h1',
                { className: 'title' },
                'Create New Article'
            ),
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(__WEBPACK_IMPORTED_MODULE_3__form__["a" /* default */], {
                formData: {},
                onSave: onSave,
                onCancelRedirectTo: '/articles'
            })
        );
    }
}

/* harmony default export */ __webpack_exports__["a"] = (Object(__WEBPACK_IMPORTED_MODULE_1_react_redux__["connect"])(state => {}, dispatch => Object(__WEBPACK_IMPORTED_MODULE_2_redux__["bindActionCreators"])({ createRecord: __WEBPACK_IMPORTED_MODULE_4__actions__["a" /* createRecord */] }, dispatch))(NewPage));

/***/ }),
/* 20 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0_react__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_simplemde__ = __webpack_require__(21);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_simplemde___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_1_simplemde__);
var _extends = Object.assign || function (target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i]; for (var key in source) { if (Object.prototype.hasOwnProperty.call(source, key)) { target[key] = source[key]; } } } return target; };




class Textarea extends __WEBPACK_IMPORTED_MODULE_0_react___default.a.Component {
    shouldComponentUpdate(nextProps, nextState) {
        return nextProps.value != this.simplemde.value();
    }

    setElement(element) {
        this.element = element;
        if (element) {
            this.simplemde = this.newSimpleMDE(element);
            this.simplemde.codemirror.on('change', () => {
                this.props.onChange({
                    target: {
                        name: this.props.name,
                        value: this.simplemde.value()
                    }
                });
            });
        } else {
            this.simplemde.toTextArea();
            this.simplemde = null;
        }
    }

    render() {
        return __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement('textarea', _extends({}, this.props, { ref: el => {
                this.setElement(el);
            } }));
    }

    newSimpleMDE(element) {
        const smde = new __WEBPACK_IMPORTED_MODULE_1_simplemde___default.a({
            element,
            autofocus: true,
            spellChecker: false,
            toolbar: [{
                name: "bold",
                action: __WEBPACK_IMPORTED_MODULE_1_simplemde___default.a.toggleBold,
                title: "Bold",
                className: "fa fa-bold"
            }, {
                name: "italic",
                action: __WEBPACK_IMPORTED_MODULE_1_simplemde___default.a.toggleItalic,
                title: "Italic",
                className: "fa fa-italic"
            }, {
                name: "strikethrough",
                action: __WEBPACK_IMPORTED_MODULE_1_simplemde___default.a.toggleStrikethrough,
                title: "Strikethrough",
                className: "fa fa-strikethrough"
            }, {
                name: "quote",
                action: __WEBPACK_IMPORTED_MODULE_1_simplemde___default.a.toggleBlockquote,
                title: "Quote",
                className: "fa fa-quote-left"
            }, "|", {
                name: "unordered-list",
                action: __WEBPACK_IMPORTED_MODULE_1_simplemde___default.a.toggleUnorderedList,
                title: "Generic List",
                className: "fa fa-list-ul"
            }, {
                name: "ordered-list",
                action: __WEBPACK_IMPORTED_MODULE_1_simplemde___default.a.toggleOrderedList,
                title: "Numbered List",
                className: "fa fa-list-ol"
            }, "|", {
                name: "headword",
                action: ({ codemirror }) => {
                    this.headwordAction(codemirror);
                },
                title: "Headword",
                className: "fa fa-header color-tag"
            }]
        });
        smde.codemirror.addKeyMap({
            'Shift-Alt-W': this.headwordAction,
            'Cmd-S': this.props.onSave
        }, true);
        return smde;
    }

    headwordAction(codemirror) {
        const selection = codemirror.getSelection();
        codemirror.replaceSelection(`<v-hw>${selection}</v-hw>`, 'around');
        codemirror.focus();
    }
}
/* harmony export (immutable) */ __webpack_exports__["a"] = Textarea;


/***/ }),
/* 21 */
/***/ (function(module, exports) {

module.exports = SimpleMDE;

/***/ }),
/* 22 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0_react__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_redux__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_redux___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_1_react_redux__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_redux__ = __webpack_require__(2);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_redux___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_redux__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__actions__ = __webpack_require__(5);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__form__ = __webpack_require__(7);







class EditPage extends __WEBPACK_IMPORTED_MODULE_0_react___default.a.Component {
    componentWillMount() {
        this.props.fetchRecord(this.props.match.params.articleID);
    }

    componentWillUnmount() {
        this.props.leaveRecord();
    }

    render() {
        if (!this.props.data) {
            return null;
        }

        const onSave = data => {
            this.props.updateRecord(data);
        };

        return __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
            'div',
            null,
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                'h1',
                { className: 'title' },
                'Edit #',
                this.props.data.ID,
                ' `',
                this.props.data.Title,
                '`'
            ),
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(__WEBPACK_IMPORTED_MODULE_4__form__["a" /* default */], {
                formData: this.props.data,
                onSave: onSave,
                onCancelRedirectTo: '/articles'
            })
        );
    }
}

/* harmony default export */ __webpack_exports__["a"] = (Object(__WEBPACK_IMPORTED_MODULE_1_react_redux__["connect"])(state => ({
    data: state.articles.record.data
}), dispatch => Object(__WEBPACK_IMPORTED_MODULE_2_redux__["bindActionCreators"])({ fetchRecord: __WEBPACK_IMPORTED_MODULE_3__actions__["c" /* fetchRecord */], leaveRecord: __WEBPACK_IMPORTED_MODULE_3__actions__["e" /* leaveRecord */], updateRecord: __WEBPACK_IMPORTED_MODULE_3__actions__["f" /* updateRecord */] }, dispatch))(EditPage));

/***/ }),
/* 23 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0_react__);


class HomePage extends __WEBPACK_IMPORTED_MODULE_0_react___default.a.Component {
    render() {
        return __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
            'div',
            null,
            'Verbum Admin Home Page'
        );
    }
}
/* harmony export (immutable) */ __webpack_exports__["a"] = HomePage;


/***/ }),
/* 24 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0_react__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_redux__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_redux___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_1_react_redux__);



class Messages extends __WEBPACK_IMPORTED_MODULE_0_react___default.a.Component {
    render() {
        return __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
            'div',
            { className: 'messages' },
            this.props.messages.map(message => __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                'div',
                { key: message.id, className: `notification is-${message.level}` },
                message.message
            ))
        );
    }
}

/* harmony default export */ __webpack_exports__["a"] = (Object(__WEBPACK_IMPORTED_MODULE_1_react_redux__["connect"])(state => ({
    messages: state.messages
}))(Messages));

/***/ }),
/* 25 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
const messages = (state = [], action) => {
    switch (action.type) {
        case 'MESSAGES/SHOW':
            return [...state, {
                id: action.message.id,
                level: action.message.level,
                message: action.message.message
            }];
        case 'MESSAGES/DISMISS':
            return state.filter(item => item.id !== action.message.id);
        default:
            return state;
    }
};

/* harmony default export */ __webpack_exports__["a"] = (messages);

/***/ }),
/* 26 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_redux__ = __webpack_require__(2);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_redux___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0_redux__);


const list = (state = null, action) => {
    switch (action.type) {
        case 'DICTIONARIES/LIST/FETCH/FULFILLED':
            return action.Data;
        case 'DICTIONARIES/LIST/LEAVE':
            return null;
        default:
            return state;
    }
};

/* harmony default export */ __webpack_exports__["a"] = (Object(__WEBPACK_IMPORTED_MODULE_0_redux__["combineReducers"])({
    list
}));

/***/ }),
/* 27 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_redux__ = __webpack_require__(2);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_redux___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0_redux__);
var _extends = Object.assign || function (target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i]; for (var key in source) { if (Object.prototype.hasOwnProperty.call(source, key)) { target[key] = source[key]; } } } return target; };



const list = (state = null, action) => {
    switch (action.type) {
        case 'ARTICLES/LIST/FETCH/FULFILLED':
            return action.Data;
        case 'ARTICLES/LIST/LEAVE':
            return null;
        default:
            return state;
    }
};

const record = (state = {}, action) => {
    switch (action.type) {
        case 'ARTICLES/RECORD/FETCH/FULFILLED':
            return _extends({}, state, { data: action.Data });
        case 'ARTICLES/RECORD/LEAVE':
            return {};
        default:
            return state;
    }
};

/* harmony default export */ __webpack_exports__["a"] = (Object(__WEBPACK_IMPORTED_MODULE_0_redux__["combineReducers"])({
    list,
    record
}));

/***/ }),
/* 28 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0_react__);


class Task extends __WEBPACK_IMPORTED_MODULE_0_react___default.a.Component {
    constructor(props) {
        super(props);
        this.keyUpHandler = this.keyUpHandler.bind(this);
    }

    componentWillMount() {
        document.addEventListener('keyup', this.keyUpHandler);
    }

    componentWillUnmount() {
        document.removeEventListener('keyup', this.keyUpHandler);
    }

    keyUpHandler(e) {
        const { index, onToggle } = this.props;
        if (e.ctrlKey && e.keyCode - 48 == index) {
            e.preventDefault();
            e.stopPropagation();
            onToggle();
        }
    }

    render() {
        const { onToggle, task: it, index: i } = this.props;

        const style = it.Status == 'PENDING' ? 'is-info' : 'is-success';
        const icon = it.Status == 'PENDING' ? __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement('i', { 'class': 'fa fa-circle-o', 'aria-hidden': 'true' }) : __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement('i', { 'class': 'fa fa-check-circle', 'aria-hidden': 'true' });
        return __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
            'div',
            null,
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                'a',
                { 'class': `button ${style}`, onClick: onToggle },
                icon,
                '\xA0(',
                i,
                ')\xA0',
                it.Task.Title
            )
        );
    }
}
/* harmony export (immutable) */ __webpack_exports__["a"] = Task;


/***/ }),
/* 29 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_react___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0_react__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_redux__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_react_redux___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_1_react_redux__);



class FilterDictionary extends __WEBPACK_IMPORTED_MODULE_0_react___default.a.Component {
    render() {
        const { value, dicts, onChange } = this.props;
        return __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
            'div',
            { 'class': 'select' },
            __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                'select',
                { value: value, onChange: ev => onChange(parseInt(ev.target.value)) },
                __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'option',
                    { value: '-1' },
                    '- Filter by Dictionary -'
                ),
                dicts.map(d => __WEBPACK_IMPORTED_MODULE_0_react___default.a.createElement(
                    'option',
                    { value: d.ID },
                    d.Title
                ))
            )
        );
    }
}

/* harmony default export */ __webpack_exports__["a"] = (Object(__WEBPACK_IMPORTED_MODULE_1_react_redux__["connect"])(state => ({
    dicts: state.config.Dicts
}))(FilterDictionary));

/***/ })
/******/ ]);