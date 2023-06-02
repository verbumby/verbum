import { VerbumAPIClientImpl } from './client'

type VerbumAPIClientOptionsV8Bridge = {
	bridge: (url: string) => any
}

export class VerbumAPIClientV8Bridge extends VerbumAPIClientImpl {
	private bridge: (url: string) => any

    constructor(options: VerbumAPIClientOptionsV8Bridge) {
        super()
		this.bridge = options.bridge
    }

    async call<T>(path: string): Promise<T> {
		let result: T = this.bridge(path) as T
		return Promise.resolve(result)
    }
}
