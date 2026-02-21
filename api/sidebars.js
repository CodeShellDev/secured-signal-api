// @ts-check

/**
 @type {import('@docusaurus/plugin-content-docs').SidebarsConfig}
 */
const sidebars = {
  "apiSidebar": [
    {
      "type": "category",
      "label": "General",
      "link": {
        "type": "doc",
        "id": "general"
      },
      "items": [
        {
          "type": "doc",
          "id": "list-api-information",
          "label": "List API information",
          "className": "api-method get"
        }
      ]
    },
    {
      "type": "category",
      "label": "Messages",
      "link": {
        "type": "doc",
        "id": "messages"
      },
      "items": [
        {
          "type": "doc",
          "id": "send-message",
          "label": "Send message",
          "className": "api-method post"
        },
        {
          "type": "doc",
          "id": "get-scheduled-request",
          "label": "Get scheduled request",
          "className": "api-method get"
        },
        {
          "type": "doc",
          "id": "cancel-scheduled-request",
          "label": "Cancel scheduled request",
          "className": "api-method delete"
        }
      ]
    }
  ]
};

export default sidebars;
