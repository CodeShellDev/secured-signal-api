const lang = "go-template"

;(function (Prism) {
	Prism.languages[lang] = {
		// Go template comments {{/* ... */}}
		"go-template-comment": {
			pattern: /\{\{[\s\-]*\/\*[\s\S]*?\*\/[\s\-]*\}\}/,
			greedy: true,
			alias: "comment",
			inside: {
				delimiter: /^\{\{\/\*|\*\/\}\}/,
				content: /[\s\S]+/,
			},
		},

		// Regular Go template expressions {{ ... }}
		"go-template-variable": {
			pattern: /\{\{[\s\-]*(?!\/\*)[\s\S]+?[\s\-]*\}\}/,
			greedy: true,
			alias: "variable",
			inside: {
				delimiter: /^\{\{|\}\}$/,
				expression: /[\s\S]+/,
			},
		},
	}
})(Prism)

export function implement(Prism, object) {
	if (!object) return

	const base = Prism.languages[lang]

	object.inside = object.inside || {}

	Object.keys(base).forEach((attribute) => {
		object.inside[attribute] = base[attribute]
	})
}
