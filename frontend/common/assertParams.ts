export function assertParams<T extends Record<string, string>, const K extends keyof T & string>(
    params: Partial<T>,
    keys: [keyof T] extends [K] ? K[] : never
): asserts params is T {
    for (const key of keys) {
        if (params[key] === undefined) {
            throw new Error(`Route param "${key}" is undefined`)
        }
    }
}
