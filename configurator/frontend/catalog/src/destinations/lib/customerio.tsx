import { filteringExpressionDocumentation, modeParameter, tableName } from "./common"
import { stringType } from "../../sources/types"

const icon = (
    <svg height="100%"
        width="100%" viewBox="0 0 157 111" version="1.1" xmlns="http://www.w3.org/2000/svg">
        <title>customer-io-logo-vector</title>
        <g id="Page-1" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
            <g id="customer-io-logo-vector" transform="translate(0.700000, 0.800000)" fill-rule="nonzero">
                <path d="M78,62.6 C95.3,62.6 109.3,48.6 109.3,31.3 C109.3,14 95.3,0 78,0 C60.7,0 46.7,14 46.7,31.3 C46.7,31.3 46.7,31.3 46.7,31.3 C46.8,48.6 60.8,62.6 78,62.6 C78,62.6 78,62.6 78,62.6 Z" id="Path" fill="#FFCD00"></path>
                <path d="M78.2,78.2 L78.1,78.2 C56.7,78.2 38,63.7 32.7,42.9 C31.1,36.7 26,31.4 19.6,31.4 L0,31.4 C0,74.5 35,109.5 78.1,109.5 L78.1,109.5 L78.2,109.5 L78.2,78.2 Z" id="Path" fill="#00ECBB"></path>
                <path d="M78,78.2 L78,78.2 C99.5,78.2 118.2,63.7 123.5,42.9 C125.1,36.7 130.2,31.4 136.6,31.4 L156.2,31.4 C156.2,74.5 121.2,109.5 78.1,109.5 L78,109.5 L78,78.2 Z" id="Path" fill="#AF64FF"></path>
                <path d="M133.4,86.6 C102.9,117.1 53.4,117.1 22.9,86.6 C22.9,86.6 22.9,86.6 22.9,86.6 L45,64.5 C63.3,82.8 93,82.8 111.3,64.5 L133.4,86.6 Z" id="Path" fill="#7131FF"></path>
            </g>
        </g>
    </svg>
)

const customerIODestination = {
    description: (
        <>
            Jitsu can send events from JS SDK or Events API to{" "}
            <a target="_blank" href="https://customer.io/docs/api">
                CustomerIO API
            </a>
            . CustomerIO is an real-time analytics platform for marketers that can build dashboards to filter new users by
            country, user activity, retention rate and funnel audiences by custom events
        </>
    ),
    syncFromSourcesStatus: "not_supported",
    id: "customerio",
    type: "other",
    displayName: "Customer IO",
    defaultTransform: `// Code of CustomerIO transform:
// https://github.com/jitsucom/jitsu/blob/master/server/storages/transform/customerio.js
return toCustomerIO($)`,
    hidden: false,
    deprecated: false,
    ui: {
        icon,
        title: cfg => `Site ID: ${cfg._formData.siteID} API Key: ${cfg._formData.apiKey.substr(0, cfg._formData.apiKey.length / 2)}*****`,
        connectCmd: _ => null,
    },
    parameters: [
        modeParameter("stream"),
        tableName(filteringExpressionDocumentation),
        {
            id: "_formData.siteID",
            displayName: "Site ID",
            required: true,
            type: stringType,
            documentation: (
                <>
                    Your CustomerIO Site ID from{" "}
                    <a target="_blank" href="https://fly.customer.io/">
                        Project Settings
                    </a>{" "}
                    page.
                </>
            ),
        },
        {
            id: "_formData.apiKey",
            displayName: "API Key",
            required: true,
            type: stringType,
            documentation: (
                <>
                    Your CustomerIO API Key from{" "}
                    <a target="_blank" href="https://fly.customer.io/">
                        Project Settings
                    </a>{" "}
                    page.
                </>
            ),
        },
    ],
} as const

export default customerIODestination
