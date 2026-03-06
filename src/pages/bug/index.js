import React from "react"
import Heading from "@theme/Heading"
import Layout from "@theme/Layout"

import useDocusaurusContext from "@docusaurus/useDocusaurusContext"

import IssueReport from "@site/src/components/IssueReport"

import styles from "./index.module.css"

export default function ReportBugPage() {
	const { siteConfig } = useDocusaurusContext()

	return (
		<Layout
			title={`${siteConfig.title}`}
			description="Official Secured Signal API Documentation"
		>
			<IssueReport
				templateUrl={siteConfig.customFields.bug_report}
				labels={["bug"]}
			/>
			<div className={styles.centered}>
				<Heading as="h3">Creating Issue</Heading>
				<Heading as="h4">preparing bug report…</Heading>
			</div>
		</Layout>
	)
}
