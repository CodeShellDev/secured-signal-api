// @ts-check
// `@type` JSDoc annotations allow editor autocompletion and type checking
// (when paired with `@ts-check`).
// There are various equivalent ways to declare your Docusaurus config.
// See: https://docusaurus.io/docs/api/docusaurus-config

import { themes as prismThemes } from "prism-react-renderer"

// This runs in Node.js - Don't use client-side code here (browser APIs, JSX...)

/** @type {import('@docusaurus/types').Config} */
const config = {
	title: "Secured Signal API",
	tagline: "Secure Proxy for Signal Messenger API",
	favicon: "img/favicon.ico",

	// Future flags, see https://docusaurus.io/docs/api/docusaurus-config#future
	future: {
		v4: true, // Improve compatibility with the upcoming Docusaurus v4
	},

	// Set the production url of your site here
	url: "https://codeshelldev.github.io",
	// Set the /<baseUrl>/ pathname under which your site is served
	// For GitHub pages deployment, it is often '/<projectName>/'
	baseUrl: "/secured-signal-api/",
	deploymentBranch: "docs-build",
	trailingSlash: false,

	// GitHub pages deployment config.
	// If you aren't using GitHub pages, you don't need these.^
	organizationName: "CodeShell", // Usually your GitHub org/user name.
	projectName: "secured-signal-api", // Usually your repo name.

	onBrokenLinks: "throw",

	// Even if you don't use internationalization, you can use this field to set
	// useful metadata like html lang. For example, if your site is Chinese, you
	// may want to replace "en" with "zh-Hans".
	i18n: {
		defaultLocale: "en",
		locales: ["en"],
	},

	markdown: {
		mermaid: true,
	},

	themes: ["@docusaurus/theme-mermaid"],

	presets: [
		[
			"classic",
			/** @type {import('@docusaurus/preset-classic').Options} */
			({
				docs: {
					sidebarPath: "./sidebars.js",
					// Please change this to your repo.
					// Remove this to remove the "edit this page" links.
					editUrl:
						"https://github.com/codeshelldev/secured-signal-api/tree/docs",
					beforeDefaultRemarkPlugins: [
						require("remark-github-admonitions-to-directives"),
					],
					remarkPlugins: [require("remark-gfm"), require("remark-directive")],
				},
				theme: {
					customCss: ["./src/css/custom.css"],
				},
        		sitemap: {
          			lastmod: 'date',
          			changefreq: 'weekly',
          			priority: 0.5,
          			filename: 'sitemap.xml',
				},
			}),
		],
	],

	themeConfig:
		/** @type {import('@docusaurus/preset-classic').ThemeConfig} */
		({
			// Replace with your project's social card
			image: "img/social-card.png",
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
				theme: prismThemes.github,
				darkTheme: prismThemes.oneDark,
				additionalLanguages: ["bash", "apacheconf"],
				plugins: ["line-numbers", "command-line"],
			},
			metadata: [
    			{ name: "google-site-verification", content: "g8d_0UGQgwAYseQGOOqRvsTPup3xawCbb-i2jT9HyVc" },
				{ name: "og:site_name", content: "Secured Signal API" },
  			]
	}),
}

export default config
