# Contributing

Want to contribute to our documentation? 🎉

Well then you have come to the right place! 😁

## Tips & Tricks

> Why spend 5 mins doing something when you can automate it in an hour 🤡

### Templates

As you may have seen we use templates almost everywhere in the documentation pages.

And so can you, from simple `+{{{ read "file.txt" }}}` to complex reusable functions like in [functions.inc.gtmpl](./docs/templates/functions.inc.gtmpl).

We use [**goplater**](https://github.com/codeshelldev/goplater) <img align="center" width="32" height="32" src="https://github.com/codeshelldev/goplater/raw/refs/heads/main/logo/goplater.png"> for this, since it includes almost every possible function you can think of
and was developed for exactly this use case (by @CodeShellDev).

**Example applications:**

- including config files
- processing documents
- reusing markdown snippets

### Code Formatters

Formatting your code is always a good idea even for markdown documents.

So it _is_ a good idea to have one at hand, my **personal favorite** and **popular choice** is [**Prettier**](https://prettier.io/),
because it is configurable, lightweight, has many extensions like for [**VSCode**](https://marketplace.visualstudio.com/items?itemName=esbenp.prettier-vscode),
but any other will also do.

## Documentation Types

### Integrations

For [integrations](https://codeshelldev.github.io/secured-signal-api/integrations), there are some additional things you should keep in mind when creating documentation.

#### Show don't tell

Images! Images or screenshots of the individual steps are always welcome 🤗 and help the user identify if they correctly understood the step by just comparing.

As always too much is also not good, so keep it **minimal**.

Of course this isn't supposed to be a comic, so we still want to **describe to the user the step**.

_Also rather_ **no** _screenshots than_ **crappy screenshots**, as a tip try **zooming** in via the browser as close as possible before screenshotting the web interface (if applicable),
this way the resulting images have a greater resolution and therefor better quality.

#### Iconic icons

Almost every programm has some sort of logo, so try to include it in the sidebar as an icon, for better recognition.

We achieve this with this frontmatter:

```yaml
---
title: Title goes here
sidebar_custom_props:
  icon: https://example.com/path/to/image
---
```

### API

For our API which mostly consists of modified **Signal CLI REST API** endpoints is configured via the [**Open API** standard](https://www.openapis.org/) in a yaml file.

```
api
├── generated
│   ├── cancel-scheduled-request.api.mdx
│   ├── cancel-scheduled-request.ParamsDetails.json
│   ├── cancel-scheduled-request.RequestSchema.json
│   ├── cancel-scheduled-request.StatusCodes.json
│   ├── general.tag.mdx
│   ├── get-scheduled-request.api.mdx
│   ├── get-scheduled-request.ParamsDetails.json
│   ├── get-scheduled-request.RequestSchema.json
│   ├── get-scheduled-request.StatusCodes.json
│   ├── list-api-information.api.mdx
│   ├── list-api-information.RequestSchema.json
│   ├── list-api-information.StatusCodes.json
│   ├── messages.tag.mdx
│   ├── send-message.api.mdx
│   ├── send-message.RequestSchema.json
│   ├── send-message.StatusCodes.json
│   └── sidebar.ts
├── openapi-v1.5.1.yaml
├── openapi.yaml
├── sidebars.js
└── versions.json
```

These files are located under [`api/`](./api).

> [!NOTE]
> For versioning, see [**Versioning API**](#api-1).

## Standards

> How boring… 🥱

Besides all of the basics like code formatting via code blocks, admonitions or text styling we try to keep some standards that may be unique to this project.

### Make it colorful!

Now this might just be me personally, but I am much more intruiged to read a page if it is colorful and to be honest in markdown we do not have that many options
to do that, so we have to work with what we have!

#### Admonitions

Apart from the coloring I think these are a great tool for quickly communicating important informations to the user.

**Example:**

> [!TIP]
> I ❤️ admonitions

I am still personally split on when to add punctuation and when not to,
I think for really small (one sentence) admonitions a trailing dot should be ommited.

#### Code Blocks

For anything code-like code blocks are perfect, but there is a big difference between this 🥱:

```
data:
 key: value
```

And this 😎:

```yaml
data:
  key: value
```

Please be 😎 and use the correct **syntax highlighter**, there are highlighters for almost any type of code or document 🎉.

### Gotta Link 'em All

> `configure auth […]`, huh? 🤔  
> _Search:_ "How to configure auth"
>
> `setup account […]` — and now? 🤨  
> _Search:_ "How to set up account"
>
> **No results.**  
> _Damn it._ 😑

To prevent this frustrating interaction, please be so kind and link important keywords to their respective pages, thanks 😄!

**Quick Reminder:**

```md
Bla bla bla… **[Setup Account](https://example.com/path/to/setup-account)!**
```

### Structuring is key!

Try to introduce sections by using markdown headings, this can help the user navigate and link certain parts of your documentation page.

#### Lists

Sometimes a short list can be worlds better than some boring long text 🫤.

**Quick Reminder:**

```md
- Apple
  - green
  - red
- Banana
  - green
  - yellow
  - brown
- Strawberry
  - red
```

#### Tables

Tables just kick differently, I mean they provide great overview of items, so they kind of act like lists, but allow for longer items and more uniform structure.

**Example:**

```md
| Value Type | Match Type |                                      | Notes                                                                                          |
| ---------- | ---------- | :----------------------------------: | ---------------------------------------------------------------------------------------------- |
| string     | `equals`   |         `pattern ~= string`          | case-incensitive                                                                               |
| string     | `contains` |      `pattern.Contains(string)`      | case-incensitive                                                                               |
| string     | `prefix`   |     `string.StartsWith(pattern)`     | case-incensitive                                                                               |
| string     | `suffix`   |      `string.EndsWith(pattern)`      | case-incensitive                                                                               |
| string     | `regex`    | example: `[^\S]` only non-whitespace | [regex](https://regex101.com)                                                                  |
| string     | `glob`     |   example: `[abc]` only `a\|b\|c`    | [glob-style pattern](https://www.gnu.org/software/bash/manual/html_node/Pattern-Matching.html) |
```

But I swear to God if I see someone merging this kind of abomination into the docs branch 😵:

```
| Value Type | Match Type |      | Notes   |
| ---------- | ---------- | :----------------------------------: | ---------------------------------------------------- |
| string | `equals`   |     `pattern ~= string`     | case-incensitive                                              |
| string  | `contains` |      `pattern.Contains(string)`    | case-incensitive                                                                  |
| string    | `prefix`   |    `string.StartsWith(pattern)`   | case-incensitive                                                                    |
| string    | `suffix`   |     `string.EndsWith(pattern)`    | case-incensitive                                              |
| string   | `regex`    | example: `[^\S]` only non-whitespace | [regex](https://regex101.com)                                             |
| string    | `glob`     |   example: `[abc]` only `a\|b\|c`  | [glob-style pattern](https://www.gnu.org/software/bash/manual/html_node/Pattern-Matching.html) |
```

Don't get me wrong I am not expecting anyone to manually format these by hand, but that's what code formatters like [**Prettier**](#code-formatters) are for.

### Keeping content up-to-date

In general you won't have to version anything since this will be probably primarily done by the maintainers, but if you would
then here's how:

```
npm run version-xyz vX.Y.Z
```

#### Regular Documentation

For the _standard_ documentation we like to stick with `vX.Y`, since patches often do not change config structure or similar,
in the rare event of a documentation change due to a path version it is okay to just update the current version instead of creating a new one.

Since we expect users to follow through with bugfixes and therefor the latest patch version always matches the latest version in the documentation.

**Example:**

```
npm run version-docs vX.Y
```

To remove a version delete the corresponding version folder from the `[…]versioned_docs` folder and delete the version entry from the `[…]versions.json` file.

#### API

Since we use yaml **OpenAPI** files we cannot use docusaurus' builtin versioning, instead we rely on a custom solution with build-time converting into markdown files and then saving into
the standard `versioned_api` folders.

**Example:**

```
npm run version-api vX.Y.Z
```

This generates `openapi-vX.Y.Z.yaml` file in `api/`, updates the `versions.json` file (**do not modify!**) and adds the new version to docusaurus' registry in `api_versions.json`.

To remove a version delete the corresponding yaml file, delete the version entry from both version files and you're good to go!

### The Ugly

*README*s are great, but after having moved to a hosted documentation solution maintaining both at equal quality is really hard.

But since in the end the first thing people see is the README it does need to include _something_…

This is where we **recycle** content from the documentation pages for the README with some **preprocessing**.

```
+{{{ funcCallArgs "parseDocs" "path/to/doc.md" "## Title" }}}
```

With this snippet you are almost good to go.

This way we can reuse preexisting content and also keep a single source of truth.

If you want the details on what the [template](#templates) `parseDocs` func does, here are the most important steps:

- removes frontmatter:
  ```
  ---
  title: New Title 1
  ---
  ```
- removes the first introduction line:
  ```
  In this section we'll be taking a look at how to use **Secured Signal API**.
  ```
- bumps headings depending on given root title (`## Title`)
- normalizes links:
  from `./endpoints.md` to `https://codeshelldev.github.io/secured-signal-api/docs/configuration/endpoints`
- prepends title from argument to content
