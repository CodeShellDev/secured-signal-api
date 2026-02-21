import {
	useVersions,
	useActiveDocContext,
} from "@docusaurus/plugin-content-docs/client"
import DocsVersionDropdownNavbarItem from "@theme-original/NavbarItem/DocsVersionDropdownNavbarItem"

export default function DocsVersionDropdownNavbarItemWrapper(props) {
	const { docsPluginId } = props

	const activeDocContext = useActiveDocContext(docsPluginId)
	const versions = useVersions(docsPluginId)

	if (!activeDocContext.activeDoc) {
		return null
	}

	if (versions.length == 0) {
		return null
	}

	if (versions.length == 1 && versions[0].name === "current") {
		return null
	}

	return <DocsVersionDropdownNavbarItem {...props} />
}
