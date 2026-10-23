declare module '*.css'

// Bootstrap's ES module sources bundle smaller than js/dist, which @types/bootstrap covers.
declare module 'bootstrap/js/src/tooltip' {
    export { default } from 'bootstrap/js/dist/tooltip'
}
declare module 'bootstrap/js/src/dropdown' {
    export { default } from 'bootstrap/js/dist/dropdown'
}
