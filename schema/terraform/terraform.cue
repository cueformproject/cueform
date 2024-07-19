package terraform

#makeSchema: {
	providers: [...{#resources: [string]: _}]

	schema: {
		terraform: #terraformConfig & {
			required_providers: {for p in providers {
				(p.#metadata.name): {
					source:  *p.#metadata.source | string
					version: *p.#metadata.version | string
				}
			}}
		}

		provider: {
			for p in providers {
				(p.#metadata.name): p.#provider
			}
		}

		resource: {
			for p in providers
			for k, v in p.#resources {(k)?: [ID=string]: v & {
				#name: string
				#id:   ID

				#ref: "\(#name).\(#id)"

				// TODO: not working
				// for kk, _ in v {
				// 	#ref: (kk): "${\(#name).\(#id).\(kk)}"
				// }
			}}
		}

		data: {
			for p in providers
			for k, v in p.#data_sources {(k)?: [ID=string]: v & {
				#name: string
				#id:   ID
				#ref:  "data.\(#name).\(#id)"
			}}
		}
	}
}

// Alternative to [resource].#ref
// #makeRef: {
// 	target: _
// 	field:  string

// 	ref: "${\(target.#ref).\(field)}"
// }

#terraformConfig: {
	required_providers: [string]: {
		source:  string
		version: string
	}
}
