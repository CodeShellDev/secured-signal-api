import DocSidebarItem from "@theme-original/DocSidebarItem"
import useBaseUrl from "@docusaurus/useBaseUrl"

export default function DocSidebarItemWrapper(props) {
	const { item } = props

	let icon = item?.customProps?.icon

	if (icon) {
		// Convert to full absolute URL for CSS
		icon = useBaseUrl(icon, { absolute: true })
	}

	if (icon) {
		return (
			<DocSidebarItem
				{...props}
				style={{ "--var-custom-icon-url": `url(${icon})` }}
				data-custom-icon={icon}
			></DocSidebarItem>
		)
	}

	return <DocSidebarItem {...props}></DocSidebarItem>
}
