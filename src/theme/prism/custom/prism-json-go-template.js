import { implement } from "./prism-go-template"
;(function (Prism) {
	const json = Prism.languages.json
	if (!json) return

	implement(Prism, json.string)
})(Prism)
