import { implement } from "./prism-go-template"
;(function (Prism) {
	const yaml = Prism.languages.yaml
	if (!yaml) return

	implement(Prism, yaml.scalar)
	implement(Prism, yaml.string)
})(Prism)
