import React from "react"
import DefaultAdmonitionTypes from "@theme-original/Admonition/Types"
import AdmonitionTypeImportant from "./Type/Important.js"

const AdmonitionTypes = {
	...DefaultAdmonitionTypes,
	important: AdmonitionTypeImportant,
}

export default AdmonitionTypes
