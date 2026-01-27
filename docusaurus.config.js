import { themes as prismThemes } from "prism-react-renderer"

import remark_githubAlertsToDirectives from "./plugins/remark-alerts-to-directives/index.js"
import remark_gfm from "remark-gfm"
import remark_directive from "remark-directive"
import remark_prettierIgnore from "./plugins/remark-prettier-ignore/index.js"

import goplater from "./plugins/goplater/index.js"

const baseUrl = "/secured-signal-api/"

/** @type {import('@docusaurus/types').Config} */
const config = {
	title: "Secured Signal API",
	tagline: "Secure Proxy for Signal CLI REST API",
	favicon: "favicon/favicon.ico",

	future: {
		v4: true,
	},

	url: "https://codeshelldev.github.io",

	baseUrl: baseUrl,
	deploymentBranch: "docs-build",
	trailingSlash: false,

	organizationName: "codeshelldev",
	projectName: "secured-signal-api",

	onBrokenLinks: "throw",

	i18n: {
		defaultLocale: "en",
		locales: ["en"],
	},

	markdown: {
		mermaid: true,
		preprocessor: ({ filePath, fileContent }) => {
			return goplater({ path: filePath, content: fileContent })
		},
	},

	themes: ["@docusaurus/theme-mermaid"],

	presets: [
		[
			"classic",
			/** @type {import('@docusaurus/preset-classic').Options} */
			({
				docs: {
					sidebarPath: "./sidebars.js",
					editUrl:
						"https://github.com/codeshelldev/secured-signal-api/tree/docs",
					beforeDefaultRemarkPlugins: [
						[
							remark_githubAlertsToDirectives,
							{
								mapping: {
									CAUTION: "danger",
								},
							},
						],
					],
					remarkPlugins: [remark_gfm, remark_directive, remark_prettierIgnore],
				},
				theme: {
					customCss: ["./src/css/custom.css", "./src/css/alerts.css"],
				},
				sitemap: {
					lastmod: "date",
					changefreq: "weekly",
					priority: 0.5,
					filename: "sitemap.xml",
				},
			}),
		],
	],

	themeConfig:
		/** @type {import('@docusaurus/preset-classic').ThemeConfig} */
		({
			image: "img/banner.png",
			colorMode: {
				respectPrefersColorScheme: true,
			},
			navbar: {
				title: "Secured Signal",
				logo: {
					alt: "Logo",
					src: "img/logo_background.svg",
				},
				items: [
					{
						type: "docSidebar",
						sidebarId: "documentationSidebar",
						position: "left",
						label: "Documentation",
					},
					{
						type: "docsVersionDropdown",
						position: "right",
						dropdownActiveClassDisabled: true,
					},
					{
						href: "https://github.com/codeshelldev/secured-signal-api",
						label: "GitHub",
						position: "right",
					},
				],
			},
			footer: {
				style: "dark",
				links: [
					{
						title: "Docs",
						items: [
							{
								label: "Documentation",
								to: "/docs/about",
							},
						],
					},
					{
						title: "Community",
						items: [
							{
								label: "Github Discussions",
								href: "https://github.com/codeshelldev/secured-signal-api/discussions",
							},
						],
					},
					{
						title: "More",
						items: [
							{
								label: "GitHub",
								href: "https://github.com/codeshelldev/secured-signal-api",
							},
						],
					},
				],
				copyright: `Copyright © ${new Date().getFullYear()} CodeShellDev. Built with Docusaurus.`,
			},
			prism: {
				theme: prismThemes.oneLight,
				darkTheme: prismThemes.oneDark,
				additionalLanguages: [
					"bash",
					"apacheconf",
					"custom-go-template",
					"custom-yaml-go-template",
					"custom-yaml-empty-key",
					"custom-json-go-template",
				],
			},
			metadata: [
				{
					name: "google-site-verification",
					content: "g8d_0UGQgwAYseQGOOqRvsTPup3xawCbb-i2jT9HyVc",
				},

				{
					property: "og:site_name",
					content: "Secured Signal API",
				},
			],
		}),

	headTags: [
		{
			tagName: "link",
			attributes: {
				rel: "icon",
				type: "image/png",
				sizes: "32x32",
				href: `${baseUrl}/favicon/favicon-32x32.png`,
			},
		},
		{
			tagName: "link",
			attributes: {
				rel: "icon",
				type: "image/png",
				sizes: "16x16",
				href: `${baseUrl}favicon/favicon-16x16.png`,
			},
		},
		{
			tagName: "link",
			attributes: {
				rel: "apple-touch-icon",
				sizes: "180x180",
				href: `${baseUrl}favicon/apple-touch-icon.png`,
			},
		},
		{
			tagName: "link",
			attributes: {
				rel: "manifest",
				href: `${baseUrl}favicon/site.webmanifest`,
			},
		},
		{
			tagName: "link",
			attributes: {
				rel: "icon",
				type: "image/png",
				sizes: "192x192",
				href: `${baseUrl}favicon/android-chrome-192x192.png`,
			},
		},
		{
			tagName: "link",
			attributes: {
				rel: "icon",
				type: "image/png",
				sizes: "512x512",
				href: `${baseUrl}favicon/android-chrome-512x512.png`,
			},
		},
	],
}

export default config
