import clsx from "clsx"
import Heading from "@theme/Heading"
import styles from "./styles.module.css"

import { useColorMode } from "@docusaurus/theme-common"

const FeatureList = [
	{
		title: "Secure Layer",
		lightSvg: require("@site/static/img/light/shield.svg").default,
		darkSvg: require("@site/static/img/dark/shield.svg").default,
		description: (
			<>
				Secured Signal API main focus is and was to be a secure layer for
				signal-cli-rest-api using Bearer, Basic and Query Auth.
			</>
		),
	},
	{
		title: "Quality of Life",
		lightSvg: require("@site/static/img/light/heart.svg").default,
		darkSvg: require("@site/static/img/dark/heart.svg").default,
		description: (
			<>
				Implements many Quality of Life features, to elevate the Developer and
				User Experience.
			</>
		),
	},
	{
		title: "Compatibility in Mind",
		lightSvg: require("@site/static/img/light/chain.svg").default,
		darkSvg: require("@site/static/img/dark/chain.svg").default,
		description: (
			<>
				Secured Signal API was built with Compatibility in Mind. And supports
				almost any signal-cli-rest-api-compatible Programm.
			</>
		),
	},
]

function Feature({ title, description, lightSvg, darkSvg }) {
	const { colorMode } = useColorMode()
	const Svg = colorMode === "dark" ? darkSvg : lightSvg

	return (
		<div className={clsx("col col--4")}>
			<div className="text--center">
				<Svg className={styles.featureSvg} role="img" />
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
