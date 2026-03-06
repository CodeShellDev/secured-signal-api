import React, { useEffect } from "react"
import useDocusaurusContext from "@docusaurus/useDocusaurusContext"

export default function IssueReport({ templateUrl, labels = [], title = "" }) {
	const { siteConfig } = useDocusaurusContext()
	const { organizationName, projectName } = siteConfig

	const REPO_ISSUE_URL = `https://github.com/${organizationName}/${projectName}/issues/new`

	useEffect(() => {
		async function buildIssue() {
			let template = await fetch(templateUrl).then((r) => r.text())

			const params = new URLSearchParams(window.location.search)

			const version = params.get("version")

			if (version) {
				labels.push(version)
			}

			params.forEach((value, key) => {
				template = template.replaceAll(`\${[{${key.toUpperCase()}}]}`, value)
			})

			const regex = /\$\{\[\{[A-Z0-9_]+\}\]\}/g

			template = template.replaceAll(regex, "")

			const body = encodeURIComponent(template)

			const labelString = Array.isArray(labels) ? labels.join(",") : labels

			const url = `${REPO_ISSUE_URL}?labels=${encodeURIComponent(
				labelString,
			)}&title=${title}&body=${body}${version ? `&milestone=${version}` : ""}`

			window.location.href = url
		}

		buildIssue()
	}, [templateUrl, REPO_ISSUE_URL, labels])

	return ""
}
