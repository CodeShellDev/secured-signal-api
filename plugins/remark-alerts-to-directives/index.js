/**
 * MIT License
 *
 * Copyright (c) 2024 Incentro Nederland B.V.
 * Modifications copyright (c) 2026 CodeShellDev
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

import { visit } from "unist-util-visit"

const GITHUB_ALERT_DECLARATION_REGEX = /^\s*\[\!(?<type>\w+)\]\s*$/

export default function remarkGithubAdmonitionsToDirectives(options) {
	const { mapping = {} } = options ?? {}

	return (tree) => {
		visit(tree, "blockquote", (node, index, parent) => {
			const githubAlert = parseGithubAlertBlockquote(node)

			if (githubAlert === null) return

			const directive = {
				type: "containerDirective",
				name: mapping[githubAlert.type] ?? githubAlert.type.toLowerCase(),
				children: githubAlert.children,
			}

			if (parent === undefined || index === undefined) return
			parent.children[index] = directive
		})
	}
}

function parseGithubAlertDeclaration(text) {
	const match = text.match(GITHUB_ALERT_DECLARATION_REGEX)

	const type = match?.groups?.type

	return type ?? null
}

function parseGithubAlertBlockquote(node) {
	const [firstChild, ...blockQuoteChildren] = node.children

	if (firstChild?.type !== "paragraph") return null

	const [firstParagraphChild, ...paragraphChildren] = firstChild.children

	if (firstParagraphChild?.type !== "text") return null

	const [possibleTypeDeclaration, ...textNodes] =
		firstParagraphChild.value.split("\n")

	if (possibleTypeDeclaration === undefined) return null

	const type = parseGithubAlertDeclaration(possibleTypeDeclaration)

	if (type === null) return null

	const textNodeChildren =
		textNodes.length > 0 ? [{ type: "text", value: textNodes.join("\n") }] : []

	const hasParagraphChildren =
		textNodeChildren.length > 0 || paragraphChildren.length > 0

	const alertParagraphChildren = hasParagraphChildren
		? [
				{
					type: "paragraph",
					children: [...textNodeChildren, ...paragraphChildren],
				},
			]
		: []

	return {
		type,
		children: [...alertParagraphChildren, ...blockQuoteChildren],
	}
}
