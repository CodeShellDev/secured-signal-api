import Heading from "@theme/Heading"
import styles from "./styles.module.css"

import { useColorMode } from "@docusaurus/theme-common"

const FeatureList = [
	{
		title: "Secure Layer",
		Svg: require("@site/static/img/features/shield.svg").default,
		description: (
			<>
				The main focus of Secured Signal API is to provide a secure layer for
				signal-cli-rest-api, supporting <a href="docs/usage#auth">Bearer</a>,{" "}
				<a href="docs/usage#auth">Basic</a>,{" "}
				<a href="docs/usage#auth">Query Auth</a> and{" "}
				<a href="docs/usage#auth">more</a>.
			</>
		),
	},
	{
		title: "Quality of Life",
		Svg: require("@site/static/img/features/heart.svg").default,
		description: (
			<>
				Implements many <a href="docs/features">Quality-of-Life features</a>, to
				enhance the developer and user experience.
			</>
		),
	},
	{
		title: "Compatibility in Mind",
		Svg: require("@site/static/img/features/chain.svg").default,
		description: (
			<>
				Secured Signal API was built with{" "}
				<a href="docs/integrations#the-solution">compatibility in mind</a>, and
				it supports almost any signal-cli-rest-api-compatible program.
			</>
		),
	},
]

function Feature({ title, description, Svg }) {
	const { colorMode } = useColorMode()

	const svgStyle =
		colorMode === "dark" ? { filter: "brightness(0) invert(1)" } : {}

	return (
		<div className="col col--4">
			<div className="text--center">
				<Svg className={styles.featureSvg} role="img" style={svgStyle} />
			</div>
			<div className="text--center padding-horiz--md">
				<Heading as="h3">{title}</Heading>
				<p>{description}</p>
			</div>
		</div>
	)
}
export default function HomepageFeatures() {
	return (
		<section className={styles.features}>
			<div className="container">
				<div className="row">
					{FeatureList.map((props, idx) => (
						<Feature key={idx} {...props} />
					))}
				</div>
			</div>
		</section>
	)
}
