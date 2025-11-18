module.exports = () => {
	return (tree) => {
		const visit = require("unist-util-visit").visit

		visit(tree, "code", (node) => {
			if (node.value) {
				// Remove any line that contains "# prettier-ignore"
				node.value = node.value
					.split("\n")
					.filter((line) => !line.includes("# prettier-ignore"))
					.join("\n")
			}
		})
	}
}
