import { visit } from "unist-util-visit"

export default function replacePrettierIgnore() {
	return (tree) => {
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
