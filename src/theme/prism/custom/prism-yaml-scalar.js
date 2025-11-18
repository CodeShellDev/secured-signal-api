;(function (Prism) {
	// https://yaml.org/spec/1.2/spec.html#c-ns-anchor-property
	// https://yaml.org/spec/1.2/spec.html#c-ns-alias-node
	var anchorOrAlias = /[*&][^\s[\]{},]+/
	// https://yaml.org/spec/1.2/spec.html#c-ns-tag-property
	var tag =
		/!(?:<[\w\-%#;/?:@&=+$,.!~*'()[\]]+>|(?:[a-zA-Z\d-]*!)?[\w\-%#;/?:@&=+$.~*'()]+)?/
	// https://yaml.org/spec/1.2/spec.html#c-ns-properties(n,c)
	var properties =
		"(?:" +
		tag.source +
		"(?:[ \t]+" +
		anchorOrAlias.source +
		")?|" +
		anchorOrAlias.source +
		"(?:[ \t]+" +
		tag.source +
		")?)"

	globalThis.Prism.languages.yaml["scalar"] = {
		pattern: RegExp(
			/([\-:]\s*(?:\s<<prop>>[ \t]+)?[|>])[ \t]*([\s\S]*?)(?=(\r?\n[\S\-]|$))/.source.replace(
				/<<prop>>/g,
				function () {
					return properties
				}
			)
		),
		lookbehind: true,
		alias: "string",
	}
})(Prism)
