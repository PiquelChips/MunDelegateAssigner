const countries: {
    P5: string[];
    High: string[];
    Medium: string[];
    Standard: string[];
} = {
    P5: [
        "United States",
        "China",
        "Russia",
        "United Kingdom",
        "France",
    ],
    High: [
        "Germany",
        "Japan",
        "India",
        "Brazil",
        "Canada",
        "Australia",
        "South Korea",
        "Italy",
        "Spain",
        "Saudi Arabia",
    ],
    Medium: [
        "Mexico",
        "Indonesia",
        "Turkey",
        "Netherlands",
        "Switzerland",
        "Sweden",
        "Poland",
        "Argentina",
        "Nigeria",
        "South Africa",
        "Egypt",
        "Pakistan",
        "Vietnam",
        "UAE",
        "Israel",
    ],
    Standard: [
        "Malaysia",
        "Singapore",
        "Thailand",
        "Philipines",
        "Chile",
        "Peru",
        "Colombia",
        "Morocco",
        "Kenya",
        "Ghana",
        "Ethiopia",
        "Iraq",
        "Iran",
        "Kuwait",
        "Qatar",
        "New Zealand",
        "Portugal",
        "Ireland",
        "Greece",
        "Luxemburg",
        "Denmark",
        "Norway",
        "Finland",
        "Hungary",
        "Czech Republic",
        "Romania",
        "Austria",
        "Jordan",
        "Lebanon",
        "Syria",
        "Algeria",
        "Tunisia",
        "Bangladesh",
        "Venezuela",
        "Cuba",
        "Chile",
        "Sri Lanka",
        "Ecuador",
        "Bolivia",
        "Zambia",
        "Malawi",
        "Namibia",
        "Kiribati",
        "Eritria",
        "Gambia",
        "Jamaica",
        "Yemen",
        "Myanmar",
        "Belize",
        "Iceland",
        "Kazakstan",
        "Rwanda",
        "Laos",
        "Uganda",
        "Afghanistan",
        "Sudan",
        "North Korea",
        "Central African Republic",
        "DRC",
        "Haiti",
        "South Sudan",
        "Somalia",
        "Libya",
        "Chad",
        "Niger",
        "Papua New Guinea",
        "Timor-Leste",
        "Belarus",
        "Uzbekistan",
        "Djibouti",
        "Comoros",
        "Vanuatu",
        "Tuvalu",
        "Lesotho",
        "Burundi",
        "Togo",
        "Eswatini",
        "Lichtenstein",
        "Cambodia",
        "Estonia",
        "Georgia",
        "Gabon",
        "Botswana",
        "Panama",
    ],
};

const important_countries = countries.P5.concat(countries.High).sort(() =>
    0.5 - Math.random()
);
const other_countries = countries.Medium.concat(countries.Standard).sort(() =>
    0.5 - Math.random()
);

// get delegates
let delegates: string[] = [];
// get chairs
const chairs: string[] = [];

// the final assignments to display
const assignments: { delegate: string; country: string | undefined }[] = [];

const assign_delegates = () => {
    delegates = delegates.sort(() => Math.random() - 0.5);

    for (let index = 0; index < delegates.length; ++index) {
        const delegate = delegates[index];

        if (chairs.includes(delegate)) {
            assignments.push({ delegate, country: "Chair" });
            continue;
        }

        if (important_countries.length > 0) {
            assignments.push({ delegate, country: important_countries.pop() });
            continue;
        }

        if (other_countries.length > 0) {
            assignments.push({ delegate, country: other_countries.pop() });
            continue;
        }
    }
};
