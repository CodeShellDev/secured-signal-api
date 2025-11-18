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
	// https://yaml.org/spec/1.2/spec.html#ns-plain(n,c)
	// This is a simplified version that doesn't support "#" and multiline keys
	// All these long scarry character classes are simplified versions of YAML's characters
	var plainKey =
		/(?:[^\s\x00-\x08\x0e-\x1f!"#%&'*,\-:>?@[\]`{|}\x7f-\x84\x86-\x9f\ud800-\udfff\ufffe\uffff]|[?:-]<PLAIN>)(?:[ \t]*(?:(?![#:])<PLAIN>|:<PLAIN>))*/.source.replace(
			/<PLAIN>/g,
			function () {
				return /[^\s\x00-\x08\x0e-\x1f,[\]{}\x7f-\x84\x86-\x9f\ud800-\udfff\ufffe\uffff]/
					.source
			}
		)
	var string = /"(?:[^"\\\r\n]|\\.)*"|'(?:[^'\\\r\n]|\\.)*'/.source

	globalThis.Prism.languages.yaml["key"] = {
		pattern: RegExp(
			/((?:^|[:\-,[{\r\n?])[ \t]*(?:<<prop>>[ \t]+)?)<<key>>(?=\s*:(?:\s*|$))/.source
				.replace(/<<prop>>/g, function () {
					return properties
				})
				.replace(/<<key>>/g, function () {
					return "(?:" + plainKey + "|" + string + ")"
				})
		),
		greedy: true,
		lookbehind: true,
		alias: "atrule",
	}
})(Prism)
