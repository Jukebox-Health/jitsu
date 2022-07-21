import { filteringExpressionDocumentation, modeParameter, tableName } from "./common"
import { stringType } from "../../sources/types"

const icon = (
    <svg
        version="1.0"
        xmlns="http://www.w3.org/2000/svg"
        width="200.000000pt"
        height="200.000000pt"
        viewBox="0 0 200.000000 200.000000"
        preserveAspectRatio="xMidYMid meet">
        <g transform="translate(0.000000,200.000000) scale(0.100000,-0.100000)"
            fill="#000000" stroke="none">
            <path d="M1475 1179 c-22 -20 -40 -48 -49 -78 -20 -69 -60 -141 -98 -180 -159
    -161 -402 -187 -583 -60 l-58 41 -99 -99 -99 -98 52 -47 c61 -55 186 -120 279
    -145 44 -11 103 -17 180 -17 137 0 208 16 330 76 184 91 341 303 379 512 6 33
    11 75 11 93 l0 33 -105 0 c-102 0 -105 -1 -140 -31z"/>
        </g>
    </svg>

)

const customerIODestination = {
    description: (
        <>
            Jitsu can send events from JS SDK or Events API to{" "}
            <a target="_blank" href="https://developers.amplitude.com/docs/http-api-v2>">
                CustomerIO API
            </a>
            . CustomerIO is an real-time analytics platform for marketers that can build dashboards to filter new users by
            country, user activity, retention rate and funnel audiences by custom events
        </>
    ),
    syncFromSourcesStatus: "not_supported",
    id: "customerio",
    type: "other",
    displayName: "CustomerIO",
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
                    <a target="_blank" href="https://analytics.amplitude.com/">
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
                    <a target="_blank" href="https://analytics.amplitude.com/">
                        Project Settings
                    </a>{" "}
                    page.
                </>
            ),
        },
    ],
} as const

export default customerIODestination
